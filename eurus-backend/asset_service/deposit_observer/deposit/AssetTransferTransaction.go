package deposit

import (
	"eurus-backend/foundation/ethereum"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionStatus int16

const (
	TransStatusUnknown TransactionStatus = iota - 1
	TransStatusFailed
	TransStatusSuccess
)

type AssetTransferTransaction struct {
	OriginalTransaction *types.Transaction
	AssetName           string
	// Status              TransactionStatus
	sender      string //Mainnet address
	Receiptant  string //Side chain address
	Block       *types.Block
	TransferLog []types.Log
}

func NewAssetTransferTransaction(block *types.Block, tx *types.Transaction, transferLog []types.Log) *AssetTransferTransaction {
	trans := new(AssetTransferTransaction)
	trans.OriginalTransaction = tx
	// trans.Status = TransStatusUnknown
	trans.Block = block
	trans.TransferLog = transferLog
	return trans
}

func (me *AssetTransferTransaction) GetSender() (string, error) {
	if me.sender == "" {
		chainConfig := ethereum.GetChainConfigFromChainId(me.OriginalTransaction.ChainId())
		signer := types.MakeSigner(chainConfig, me.Block.Number())
		sender, err := me.OriginalTransaction.AsMessage(signer, nil)
		if err != nil {
			return "", err
		}
		me.sender = strings.ToLower(sender.From().Hex())
	}
	return me.sender, nil
}

func (me *AssetTransferTransaction) Hash() common.Hash {
	return me.OriginalTransaction.Hash()
}

func (me *AssetTransferTransaction) GetTo() *common.Address {
	return me.OriginalTransaction.To()
}
