package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"strings"
	"time"

	"eurus-backend/foundation"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"eurus-backend/sign_service/sign_api"

	"github.com/ethereum/go-ethereum"
	go_ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EthClient struct {
	Protocol   string
	IP         string
	Port       int
	Client     *EthClientCommunicator
	Subscriber *BlockSubscriber
	ChainID    *big.Int
	Logger     *logrus.Logger
}

const ETHTransferStandardGasLimit uint64 = 21000
const DefaultGasLimit uint64 = 8000000

func (me *EthClient) Close() {
	me.Client.Close()
}

func (me *EthClient) TransferETHInEIP1559(privateHexKey string, to string, amount *big.Int, isSmartContractAddress bool, txFeeFromAmount bool, estimateGasOnly bool, designatedGasFeeCap *big.Int, designatedGasTipCap *big.Int, designatedGasLimit *uint64) (*big.Int, *types.Transaction, error) {
	ethKeyPair, err := GetEthKeyPair(privateHexKey)
	if err != nil {
		return nil, nil, err
	}

	nonce, err := ethKeyPair.GetNonce(me)
	if err != nil {
		return nil, nil, err
	}

	toAddress := common.HexToAddress(to)

	gasTipCap := big.NewInt(0)
	if designatedGasTipCap != nil {
		gasTipCap.Set(designatedGasTipCap)
	} else {
		gasTipCap, err = me.GetGasTipCap()
		if err != nil {
			return nil, nil, err
		}
	}

	gasFeeCap := big.NewInt(0)
	if designatedGasFeeCap != nil {
		gasFeeCap.Set(designatedGasFeeCap)
	} else {
		header, err := me.Client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return nil, nil, err
		}

		gasFeeCap.Add(gasTipCap, new(big.Int).Mul(header.BaseFee, big.NewInt(2)))
	}

	gas := uint64(0)
	if designatedGasLimit != nil {
		gas = *designatedGasLimit
	} else if !isSmartContractAddress {
		gas = ETHTransferStandardGasLimit
	} else {
		fakeAmount := big.NewInt(0)
		fakeAmount.Sub(amount, new(big.Int).Mul(gasFeeCap, big.NewInt(int64(ETHTransferStandardGasLimit)*2)))

		msg := ethereum.CallMsg{
			From:      ethKeyPair.Address,
			To:        &toAddress,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
			Value:     fakeAmount,
			Data:      nil,
		}
		gas, err = me.Client.EstimateGas(context.Background(), msg)
		if err != nil {
			return nil, nil, err
		}
	}

	effectiveAmount := big.NewInt(0)
	if txFeeFromAmount {
		effectiveAmount.Sub(amount, new(big.Int).Mul(gasFeeCap, big.NewInt(int64(gas))))
	} else {
		effectiveAmount.Set(amount)
	}

	baseTx := types.DynamicFeeTx{
		ChainID:   me.ChainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gas,
		To:        &toAddress,
		Value:     effectiveAmount,
		Data:      nil,
	}

	txData := types.TxData(&baseTx)
	signedTx, err := types.SignNewTx(ethKeyPair.PrivateKey, types.NewLondonSigner(me.ChainID), txData)
	if err != nil {
		return nil, nil, err
	}

	if estimateGasOnly {
		return effectiveAmount, signedTx, nil
	}

	err = me.SendTransaction(signedTx)
	if err != nil {
		return nil, nil, err
	}

	return effectiveAmount, signedTx, nil
}

func (me *EthClient) TransferETHToSmartContract(privateHexKey string, to string, amount *big.Int, estimateGasOnly bool, designatedGasPrice *big.Int, designatedGasLimit *uint64) (string, *types.Transaction, error) {
	ethKeyPair, err := GetEthKeyPair(privateHexKey)
	if err != nil {
		return "", nil, err
	}
	nonce, err := ethKeyPair.GetNonce(me)
	if err != nil {
		return "", nil, err
	}
	gasPrice, err := me.GetGasPrice()
	if err != nil {
		return "", nil, err
	}

	toAddress := common.HexToAddress(to)

	LegacyTx := types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    amount,
		Gas:      ETHTransferStandardGasLimit + 100000,
		GasPrice: gasPrice,
		Data:     nil,
	}

	if designatedGasPrice != nil {
		LegacyTx.GasPrice = designatedGasPrice
	}

	if designatedGasLimit != nil {
		LegacyTx.Gas = *designatedGasLimit
	}

	toData := types.TxData(&LegacyTx)
	signedTx, err := types.SignNewTx(ethKeyPair.PrivateKey, types.NewEIP2930Signer(me.ChainID), toData)

	//signedTx, err := types.SignTx(tx, types.NewEIP155Signer(me.ChainID), ethKeyPair.PrivateKey)
	if err != nil {
		return "", nil, err
	}

	if estimateGasOnly {
		return signedTx.Hash().Hex(), signedTx, nil
	}

	err = me.SendTransaction(signedTx)
	if err != nil {
		return "", nil, err
	}
	return signedTx.Hash().Hex(), signedTx, nil
}

func (me *EthClient) TransferETH(privateHexKey string, to string, amount *big.Int) (string, *types.Transaction, error) {
	ethKeyPair, err := GetEthKeyPair(privateHexKey)
	if err != nil {
		return "", nil, err
	}
	nonce, err := ethKeyPair.GetNonce(me)
	if err != nil {
		return "", nil, err
	}
	gasPrice, err := me.GetGasPrice()

	if err != nil {
		return "", nil, err
	}
	toAddress := common.HexToAddress(to)

	LegacyTx := types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    amount,
		Gas:      ETHTransferStandardGasLimit + 25000,
		GasPrice: gasPrice,
		Data:     nil,
	}

	toData := types.TxData(&LegacyTx)
	signedTx, err := types.SignNewTx(ethKeyPair.PrivateKey, types.NewEIP2930Signer(me.ChainID), toData)

	//tx := types.NewTransaction(nonce, toAddress, amount, ETHTransferStandardGasLimit, gasPrice, nil)
	//signedTx, err := types.SignTx(tx, types.NewEIP155Signer(me.ChainID), ethKeyPair.PrivateKey)
	if err != nil {
		return "", nil, err
	}

	err = me.SendTransaction(signedTx)
	if err != nil {
		return "", nil, err
	}
	return signedTx.Hash().Hex(), signedTx, nil
}

func (me *EthClient) validateClientField() {
	var err error = nil
	if me.Protocol == "" {
		err = errors.New("EthClientProtocol should be provided!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
		panic(err)
	} else if me.IP == "" {
		err = errors.New("EthClientIP should be provided!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
		panic(err)
	} else if me.Port == 0 {
		err = errors.New("EthClientPort should be provided!")
		logger := log.GetLogger(log.Name.Root)
		logger.Error(err.Error())
		panic(err)
	}
}

func (me *EthClient) Connect() (*EthClientCommunicator, error) {
	me.validateClientField()

	urlStr := fmt.Sprintf("%s://%s", me.Protocol, me.IP)
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	url.Host = url.Host + ":" + strconv.Itoa(me.Port)
	client, err := Dial(url.String())
	if err != nil {
		return nil, err
	}
	me.Client = client
	return client, nil
}

func (me *EthClient) SendTransaction(signedTx *types.Transaction) error {
	err := me.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	return nil
}

/// functor returns transaction, isAborted, err object
/// if isAborted = true, no more retry on error
/// if error is null, no more retry and isAborted is ignored
func (me *EthClient) InvokeSmartContract(retrySetting foundation.IRetrySetting, privateKey string,
	targetGasLimit uint64,
	functor func(*EthClient, *bind.TransactOpts) (*types.Transaction, bool, error)) (*types.Transaction, error) {
	var err error
	var transOpt *bind.TransactOpts
	var useTryAndError bool = false
	var gasLimit uint64 = DefaultGasLimit
	var tx *types.Transaction
	var isAborted bool = false

	if targetGasLimit > 0 {
		useTryAndError = true
		gasLimit = targetGasLimit
	}

	for i := 0; i < retrySetting.GetRetryCount(); i, gasLimit = i+1, gasLimit+200000 {
		transOpt, err = me.GetNewTransactorFromPrivateKey(privateKey, me.ChainID)
		if useTryAndError {
			transOpt.GasLimit = gasLimit
		}
		if err != nil {
			err = errors.WithMessage(err, "GetNewTransactorFromPrivateKey error")
			time.Sleep(retrySetting.GetRetryInterval() * time.Second)
			continue
		}
		tx, isAborted, err = functor(me, transOpt)
		if err != nil {
			if strings.Contains(err.Error(), "failed to estimate gas needed") {
				useTryAndError = true
			}
			if !isAborted {
				time.Sleep(retrySetting.GetRetryInterval() * time.Second)
				continue
			} else {
				return tx, err
			}
		} else {
			break
		}
	}
	return tx, err
}

func (me *EthClient) SubscribeBlockHeader(bufferCount int) error {
	blockSubscriber, err := NewBlockSubscriber(me, bufferCount)
	if err != nil {
		return err
	} else {
		blockSubscriber.Logger = me.Logger
		me.Subscriber = blockSubscriber
	}
	return nil
}

func (me *EthClient) GetTransaction(hash common.Hash) (*types.Transaction, bool, error) {
	tx, isPending, err := me.Client.TransactionByHash(context.Background(), hash)
	return tx, isPending, err
}

func (me *EthClient) GetTransactionByFilter(hash common.Hash, filter func(*types.Transaction) bool) (*types.Transaction, bool, error) {
	tx, isPending, err := me.Client.TransactionByHash(context.Background(), hash)
	if filter(tx) {
		return tx, isPending, err
	}
	return nil, false, nil
}

func (me *EthClient) GetBlockByNumber(blockNum *big.Int) (block *types.Block, err error) {
	block, err = me.Client.BlockByNumber(context.Background(), blockNum)
	return block, err
}

func (me *EthClient) GetBlockByHash(blockHash common.Hash) (*types.Block, error) {
	block, err := me.Client.BlockByHash(context.Background(), blockHash)
	return block, err
}

func (me *EthClient) GetHeaderByHash(blockHash common.Hash) (*types.Header, error) {
	header, err := me.Client.HeaderByHash(context.Background(), blockHash)
	return header, err
}

func (me *EthClient) GetConfirmedTransactionReceipt(hash common.Hash) (*BesuReceipt, error) {
	receipt, err := me.Client.TransactionReceipt(context.Background(), hash)

	return receipt, err
}

func (me *EthClient) GetConfirmedTransactionReceiptByHashStr(hashStr string) (*BesuReceipt, error) {
	receipt, err := me.Client.TransactionReceipt(context.Background(), common.HexToHash(hashStr))

	return receipt, err
}

func (me *EthClient) GetConfirmedTransactionReceiptByFilter(hash common.Hash, filter func(*BesuReceipt) bool) (*BesuReceipt, error) {
	receipt, err := me.Client.TransactionReceipt(context.Background(), hash)
	if filter(receipt) {
		return receipt, err
	}
	return nil, nil
}

func (me *EthClient) GetTransactionByHash(hashStr string) (*types.Transaction, error) {
	tx, _, err := me.Client.TransactionByHash(context.Background(), common.HexToHash(hashStr))
	if err != nil {
		return nil, err
	} else {
		return tx, nil
	}
}

func (me *EthClient) GetNewTransactorFromPrivateKey(privateHexKey string, chainID *big.Int) (*bind.TransactOpts, error) {
	keyPair, err := GetEthKeyPair(privateHexKey)
	if err != nil {
		return nil, err
	}
	auth, err := keyPair.NewTransactor(me, chainID)
	if err != nil {
		return nil, err
	}
	return auth, err
}

func (me *EthClient) GetNewTransactorFromSignServer(authClient auth_base.IAuth, signServerUrl string, keyType sign_api.WalletKeyType) (*bind.TransactOpts, error) {
	auth, err := NewTransactorFromServer(me, authClient, signServerUrl, keyType)
	if err != nil {
		return nil, err
	}
	return auth, err
}

func (me *EthClient) GetGasPrice() (*big.Int, error) {
	gasPrice, err := me.Client.SuggestGasPrice(context.Background())
	return gasPrice, err
}

func (me *EthClient) GetGasTipCap() (*big.Int, error) {
	return me.Client.SuggestGasTipCap(context.Background())
}

func (me *EthClient) GetBaseFee() (*big.Int, error) {
	header, err := me.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return header.BaseFee, nil
}

func (me *EthClient) isWebSocketClient() bool {
	if me.Protocol == "ws" || me.Protocol == "wss" {
		return true
	}
	return false
}

func (me *EthClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	sub, err := me.Client.SubscribeFilterLogs(ctx, q, ch)
	if err == nil {
		go me.handleSubscribeError(sub, q, ch)
	}
	return sub, err
}

func (me *EthClient) handleSubscribeError(subscription ethereum.Subscription, q ethereum.FilterQuery, ch chan<- types.Log) {
	for {
		err := <-subscription.Err()
		if err != nil {
			if me.Logger != nil {
				me.Logger.Errorln("Subcription error: ", err, ". Going to resubscribe")
			}

			for {
				sub, err := me.SubscribeFilterLogs(context.Background(), q, ch)
				if err == nil {
					subscription = sub
					break
				}
				time.Sleep(time.Second)
			}

			if me.Logger != nil {
				me.Logger.Errorln("Resubscribe successfully")
			}

		} else {
			break
		}
	}

	if me.Logger != nil {
		me.Logger.Infoln("Subscription monitor stopped")
	}
}

func (me *EthClient) QueryEthReceipt(tx *types.Transaction) (*BesuReceipt, error) {
	return me.QueryEthReceiptWithSetting(tx, 3, 20)
}

func (me *EthClient) QueryEthReceiptWithSetting(tx *types.Transaction, waitSecond int, retryCount int) (*BesuReceipt, error) {
	return me.QueryEthReceiptByTransHash(tx.Hash(), waitSecond, retryCount)
}

func (me *EthClient) QueryEthReceiptByTransHash(transHash common.Hash, waitSecond int, retryCount int) (*BesuReceipt, error) {
	var receipt *BesuReceipt
	var err error
	var isForever bool = false
	if retryCount < 0 {
		isForever = true
	}
	for i := 0; i < retryCount || isForever; i++ {
		receipt, err = me.Client.TransactionReceipt(context.Background(), transHash)
		if err != nil {
			if err == go_ethereum.NotFound {
				time.Sleep(time.Duration(waitSecond) * time.Second)
				continue
			} else {
				return nil, err
			}
		} else {
			break
		}
	}

	return receipt, err
}

func (me *EthClient) GetBalance(address common.Address) (*big.Int, error) {
	return me.Client.BalanceAt(context.Background(), address, nil)
}

func (me *EthClient) SuggestGasFeeCap() (*big.Int, error) {
	head, err := me.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	if head.BaseFee != nil {
		tip, err := me.Client.SuggestGasTipCap(context.Background())
		if err != nil {
			return nil, errors.Wrap(err, "Getting gas tip cap error")
		}
		ratio := big.NewFloat(1.25)
		b := big.NewFloat(0)
		baseFeeFloat := b.SetInt(head.BaseFee)
		baseFeeFloat = baseFeeFloat.Mul(baseFeeFloat, ratio)

		baseFee, _ := baseFeeFloat.Int(nil)
		gasFeeCap := new(big.Int).Add(tip, baseFee)

		return gasFeeCap, nil
	}
	return nil, errors.New("Base fee not found at head block")
}
