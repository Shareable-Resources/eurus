package bc_indexer

import (
	"errors"
	"eurus-backend/asset_service/asset"
	"eurus-backend/foundation/ethereum"
	"eurus-backend/user_service/user_service/user"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type ExtractedTransaction struct {
	OriginalTransaction *types.Transaction
	IsMainnetTrans      bool
	TxHash              string
	sender              string
	to                  string
	CreatedDate         time.Time
	AssetName           string
	ToUser              *user.User
	FromUser            *user.User
	Status              int16
	Amount              *big.Int
	RequestTransId      *big.Int
	ConfirmTransHash    string
	TransGasUsed        uint64
	UserGasUsed         uint64
	Remarks             string
	EffectiveGasPrice   *big.Int
	TransactionType     asset.TransType
	Block               *types.Block
	Quantity            *big.Int
	ProductId           *big.Int
	ChildObject         interface{}
}

type TopUpExtractedTransaction struct {
	ExtractedTransaction
	IsDirectTopUp bool
	TargetGas     *big.Int
}

func (me *ExtractedTransaction) GetSender() (string, error) {
	if me.sender == "" {
		if me.Block == nil {
			return "", errors.New("Missing block object")
		}
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

func (me *ExtractedTransaction) GetTo() string {
	if me.to == "" {
		if me.OriginalTransaction.To() != nil {
			me.to = strings.ToLower(me.OriginalTransaction.To().Hex())
		} else {
			me.to = "0x0"
		}
	}
	return me.to
}

func (me *ExtractedTransaction) SetTo(to string) {
	me.to = strings.ToLower(to)
}

func (me *ExtractedTransaction) SetSender(from string) {
	me.sender = strings.ToLower(from)
}

func NewTopUpExtractedTransaction() *TopUpExtractedTransaction {
	trans := new(TopUpExtractedTransaction)
	trans.ChildObject = trans
	return trans
}
