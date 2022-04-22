package approval

import (
	"eurus-backend/foundation/ethereum"
	"eurus-backend/foundation/log"
	"eurus-backend/user_service/user_service/user"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type ExtractedTransaction struct {
	OriginalTransaction *types.Transaction `json:"Transaction"`
	TxHash              string             `json:"txHash"`
	sender              string
	To                  string `json:"to"`
	// CreatedDate         time.Time `json:"createdDate"`
	TransDate time.Time
	AssetName string `json:"assetName"`
	User      *user.User
	Amount    *big.Int
	AdminFee  *big.Int
}

func NewExtractedTransaction(tx *types.Transaction, asset string, senderUser *user.User, senderAddress string, destAddress string) *ExtractedTransaction {
	ext := new(ExtractedTransaction)
	ext.OriginalTransaction = tx
	ext.To = ethereum.ToLowerAddressString(destAddress)
	ext.AssetName = asset
	ext.TxHash = tx.Hash().Hex()
	ext.User = senderUser
	ext.sender = ethereum.ToLowerAddressString(senderAddress)
	return ext
}

func (me *ExtractedTransaction) GetSender() (string, error) {
	if me.sender == "" {
		sender, err := me.OriginalTransaction.AsMessage(types.NewEIP155Signer(me.OriginalTransaction.ChainId()), nil)
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("Error to get the sender. The transaction hash :", me.OriginalTransaction.Hash().Hex())
			return "", err
		}
		me.sender = sender.From().Hex()
	}
	return me.sender, nil
}
