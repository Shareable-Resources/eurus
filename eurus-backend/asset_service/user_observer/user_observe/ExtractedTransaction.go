package userObserver

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/user_service/user_service/user"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type ExtractedTransaction struct {
	OriginalTransaction *types.Transaction `json:"Transaction"`
	TxHash              string             `json:"txHash"`
	sender              string             `json:"sender"`
	to                  string             `json:"to"`
	CreatedDate         time.Time          `json:"createdDate"`
	AssetName           string             `json:"assetName"`
	User                *user.User
	Status              bool `json:"status"`
	Amount              *big.Int
}

func (me *ExtractedTransaction) GetSender() (string, error) {
	if me.sender == "" {
		sender, err := me.OriginalTransaction.AsMessage(types.NewEIP155Signer(me.OriginalTransaction.ChainId()), nil)
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

func ExtractStatus(ethClient *ethereum.EthClient, ext *ExtractedTransaction) error {
	receipt, err := ethClient.GetConfirmedTransactionReceipt(ext.OriginalTransaction.Hash())
	if err != nil {
		return err
	}
	ext.Status = receipt.Status != 0
	return nil
}
