package ethereum

import (
	"context"
	"encoding/json"
	"eurus-backend/foundation/log"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

type EthClientCommunicator struct {
	ethclient.Client
	c *rpc.Client
}

func Dial(rawurl string) (*EthClientCommunicator, error) {
	return DialContext(context.Background(), rawurl)
}

func DialContext(ctx context.Context, rawurl string) (*EthClientCommunicator, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewEthClientCommunicator(c), nil
}

func NewEthClientCommunicator(inputRpcClient *rpc.Client) *EthClientCommunicator {
	client := &EthClientCommunicator{
		Client: *ethclient.NewClient(inputRpcClient),
		c:      inputRpcClient,
	}
	return client
}

func (me *EthClientCommunicator) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	// msg.GasPrice = big.NewInt(1)
	// msg.Value = big.NewInt(0)

	estimatedGasLimit, err := me.Client.EstimateGas(ctx, msg)
	if err != nil {
		errBytes, _ := json.Marshal(err)
		detailedErr := errors.WithMessage(err, "Inner error JSON: "+string(errBytes))
		return estimatedGasLimit, detailedErr
	}
	//duncan debug
	fmt.Println("Estimated gas limit: ", estimatedGasLimit)
	log.GetLogger(log.Name.Root).Debugln("Estimated gas limit: ", estimatedGasLimit)
	return estimatedGasLimit, nil
}

func (ec *EthClientCommunicator) TransactionReceipt(ctx context.Context, txHash common.Hash) (*BesuReceipt, error) {
	var r *BesuReceipt
	err := ec.c.CallContext(ctx, &r, "eth_getTransactionReceipt", txHash)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}

	return r, err
}
