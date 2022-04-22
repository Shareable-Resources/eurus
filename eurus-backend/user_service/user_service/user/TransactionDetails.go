package user

import (
	"eurus-backend/asset_service/asset"
	"eurus-backend/marketing/reward"
)

type ChainLocation int8

const (
	Innet   ChainLocation = 1
	Mainnet ChainLocation = 2
)

type TransactionDetail struct {
	TransType     asset.TransType `json:"transType"`
	TransDate     int64           `json:"transDate"`
	TxHash        string          `json:"txHash"`
	FromAddress   string          `json:"fromAddress"`
	ToAddress     string          `json:"toAddress"`
	Amount        string          `json:"amount"`
	ChainLocation ChainLocation   `json:"chainLocation"`
	GasUsed       uint64          `json:"gasUsed"`
	GasPrice      uint64          `json:"gasPrice"`
}

type TransferDetail struct {
	TransactionDetail
	IsSend  bool                 `json:"isSend"`
	Status  asset.TransferStatus `json:"status"`
	Remarks string               `json:"remarks"`
}

type DepositDetail struct {
	TransactionDetail
	DestAddress   string              `json:"destAddress"`
	MintTransId   uint64              `json:"mintTransId"`
	DepositTxHash string              `json:"depositTxHash"`
	Remarks       string              `json:"remarks"`
	Status        asset.DepositStatus `json:"status"`
}

type WithdrawDetail struct {
	TransactionDetail
	AdminFee       float64              `json:"adminFee"`
	TargetAddress  string               `json:"targetAddress"` //mainnet
	RequestTransId uint64               `json:"requestTransId"`
	BurnTransId    uint64               `json:"burnTransId"`
	WithdrawTxHash string               `json:"withdrawTxHash"`
	Remarks        string               `json:"remarks"`
	Status         asset.WithdrawStatus `json:"status"`
}

type PurchaseDetail struct {
	TransactionDetail
	ProductId *uint64              `json:"productId"`
	Quantity  *uint64              `json:"quantity"`
	GasFee    uint64               `json:"gasFee"`
	Remarks   string               `json:"remarks"`
	Status    asset.PurchaseStatus `json:"status"`
}

type DistributedTokenDetail struct {
	TransactionDetail
	DistributedType reward.TokenDistributedType        `json:"distributedType"`
	TriggerType     reward.TokenDistributedTriggerType `json:"triggerType"`
}

type TopUpDetail struct {
	TransactionDetail
	TransferGas   uint64            `json:"transferGas"`
	TargetGas     uint64            `json:"targetGas"`
	TransGasUsed  uint64            `json:"transGasUsed"`
	IsDirectTopUp bool              `json:"isDirectTopUp"`
	Remarks       string            `json:"remarks"`
	Status        asset.TopUpStatus `json:"status"`
}
