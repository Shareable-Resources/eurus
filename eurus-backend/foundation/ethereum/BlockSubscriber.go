package ethereum

import (
	"context"
	"errors"
	"eurus-backend/foundation"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

type BlockSubscriber struct {
	Client              *EthClient
	Subscription        ethereum.Subscription
	Header              chan *types.Header
	isReConnecting      bool
	Logger              *logrus.Logger
	buffer              *BlockBuffer
	IsVerboseLog        bool
	beginingBlockNumber *big.Int
}

func NewBlockSubscriber(client *EthClient, bufferCount int) (*BlockSubscriber, error) {
	if !client.isWebSocketClient() {
		return nil, errors.New("EthClient of the Block Subscriber should be using WebSocket Protocol")
	}

	me := new(BlockSubscriber)
	headers := make(chan *types.Header)
	me.Header = headers
	me.Client = client
	if bufferCount > 0 {
		me.buffer = NewBlockBuffer(bufferCount)
	}
	sub, err := me.Client.Client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return nil, err
	} else {
		me.Subscription = sub
		return me, nil
	}
}

func (me *BlockSubscriber) GetLatestBlock(byNumber bool) (*types.Block, *foundation.ServerError) {
	var err error
	if me.isReConnecting {
		return nil, foundation.NewErrorWithMessage(foundation.NetworkError, "Reconnecting")
	}

	var block *types.Block
	var serverErr *foundation.ServerError

	if me.buffer != nil {
		count, blockBufferNode := me.buffer.GetConfirmedBlock()
		if me.IsVerboseLog {
			me.Logger.Debug("Level: ", count)
		}
		if blockBufferNode != nil {
			if me.IsVerboseLog {
				me.Logger.Debug("Confirmed block: ", blockBufferNode.Hash())
			}
			block, serverErr = me.getBlockByHeader(blockBufferNode.BlockHeader, byNumber)
			if serverErr != nil {
				return nil, serverErr
			}
			return block, nil
		}
	}
outer:
	for {
		select {
		case err = <-me.Subscription.Err():
			if !me.isReConnecting {
				me.isReConnecting = true
				go me.reSubscribeHead()
			}
			if err == nil {
				return nil, foundation.NewErrorWithMessage(foundation.NetworkError, "Reconnecting to ethereum")
			}
			return nil, foundation.NewErrorWithMessage(foundation.NetworkError, err.Error())
		case header := <-me.Header:
			if me.buffer != nil {

				isAdded := me.buffer.AddNode(header)

				if isAdded {
					if me.beginingBlockNumber == nil {
						me.beginingBlockNumber = header.Number
					}
					if me.IsVerboseLog {
						me.Logger.Debug("Added header hash: ", header.Hash(), " block number: ", header.Number.String(), " parent hash: ", header.ParentHash)
					}
				} else {
					if me.IsVerboseLog {
						me.Logger.Debug("Orphan header found! hash: ", header.Hash(), " block number: ", header.Number.String(), " parent hash: ", header.ParentHash)
					}
					var blockHeaderStack []*types.Header = make([]*types.Header, 0)

					var currentHeader *types.Header = header

					for {
						blockHeaderStack = append([]*types.Header{currentHeader}, blockHeaderStack...)

						parentHeader, err := me.Client.GetHeaderByHash(currentHeader.ParentHash)
						if err != nil {
							me.Logger.Error("Query parent header failed. ", err, " Parent hash: ", currentHeader.ParentHash)
							break
						} else {
							if me.beginingBlockNumber != nil && parentHeader.Number.Cmp(me.beginingBlockNumber) < 0 {
								me.Logger.Error("Query parent header number is smaller than the initial block number. Query number: ", parentHeader.Number.String(), " initial number: ", me.beginingBlockNumber.String())
								break
							}

							if me.IsVerboseLog {
								me.Logger.Debug("Trying to add orphan parent. header hash: ", parentHeader.Hash(), " block number: ", parentHeader.Number.String(), " parent hash: ", parentHeader.ParentHash)
							}
							isAdded = me.buffer.AddNode(parentHeader)
							if isAdded {
								for _, childNode := range blockHeaderStack {
									if me.IsVerboseLog {
										me.Logger.Debug("Adding back Orphan header hash: ", childNode.Hash(), " block number: ", childNode.Number.String(), " parent hash: ", childNode.ParentHash)
									}
									me.buffer.AddNode(childNode)
								}
								break
							} else {
								if me.IsVerboseLog {
									me.Logger.Debug("Orphan parent added failed header hash: ", parentHeader.Hash(), " block number: ", parentHeader.Number.String(), " parent hash: ", parentHeader.ParentHash)
								}
								currentHeader = parentHeader
								continue
							}
						}
					}
				}
				count, blockBufferNode := me.buffer.GetConfirmedBlock()
				if me.IsVerboseLog {
					me.Logger.Debug("Level: ", count)
				}

				if blockBufferNode != nil {
					if me.IsVerboseLog {
						me.Logger.Debug("Confirmed block: ", blockBufferNode.Hash(), " block number: ", blockBufferNode.BlockHeader.Number.String())
					}
					block, serverErr = me.getBlockByHeader(blockBufferNode.BlockHeader, byNumber)
					if serverErr != nil {
						return nil, serverErr
					}
					break outer
				}
			} else {
				block, serverErr = me.getBlockByHeader(header, byNumber)
				if serverErr != nil {
					return nil, serverErr
				}
				break outer
			}
		}
	}
	return block, nil
}

func (me *BlockSubscriber) getBlockByHeader(header *types.Header, byNumber bool) (*types.Block, *foundation.ServerError) {
	var block *types.Block
	var err error
	if !byNumber {
		block, err = me.Client.GetBlockByHash(header.Hash())
		if err != nil {
			if err.Error() == "not found" {
				return nil, foundation.NewErrorWithMessage(foundation.RecordNotFound, "Block hash: "+header.Hash().Hex()+" number: "+header.Number.String()+" not found")
			}
			return block, foundation.NewErrorWithMessage(foundation.EthereumError, err.Error())
		}
	} else {
		block, err = me.Client.GetBlockByNumber(header.Number)
		if err != nil {
			if err.Error() == "not found" {
				return nil, foundation.NewErrorWithMessage(foundation.RecordNotFound, "Block hash: "+header.Hash().Hex()+" number: "+header.Number.String()+" not found")
			}
			return block, foundation.NewErrorWithMessage(foundation.EthereumError, err.Error())
		}
	}
	return block, nil
}

func (me *BlockSubscriber) reSubscribeHead() {
	if me.Logger != nil {
		me.Logger.Infoln("Going to resubscribe block head at IP: ", me.Client.IP, " Port: ", me.Client.Port)
	}

	for {
		sub, err := me.Client.Client.SubscribeNewHead(context.Background(), me.Header)
		if err == nil {
			me.Subscription = sub
			break
		}
		time.Sleep(time.Second)
	}
	me.isReConnecting = false

	if me.Logger != nil {
		me.Logger.Infoln("Resubscribe block head successful at IP: ", me.Client.IP, " Port: ", me.Client.Port)
	}
}

// func (me *BlockSubscriber) RunSubscriber(blockHandler func(*types.Block)) {
// 	for {
// 		block, err := me.GetLatestBlock()
// 		if err != nil {
// 			if me.Logger != nil {
// 				me.Logger.Error("Unable to get block by block number: ", err.Error())
// 			}
// 		}
// 		go blockHandler(block)
// 	}
// }
