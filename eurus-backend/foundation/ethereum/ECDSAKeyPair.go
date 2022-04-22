package ethereum

import (
	"context"
	"crypto/ecdsa"
	"strconv"

	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/auth_base"
	"eurus-backend/foundation/log"
	"eurus-backend/sign_service/sign_api"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	eurus_crypto "eurus-backend/foundation/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type ECDSAKeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}

type PendingNonceResponse struct {
	Returncode int    `json:"returnCode"`
	Message    string `json:"message"`
	Nonce      string `json:"nonce"`
	Data       struct {
		PendingNonce uint64 `json:"pendingNonce"`
	} `json:"data"`
}

func (me *ECDSAKeyPair) NewEmptyTransactor() *bind.TransactOpts {
	auth := bind.NewKeyedTransactor(me.PrivateKey)
	return auth
}

func (me *ECDSAKeyPair) NewTransactor(ethClient *EthClient, chainID *big.Int) (*bind.TransactOpts, error) {
	auth, _ := bind.NewKeyedTransactorWithChainID(me.PrivateKey, chainID)
	nonce, err := me.GetNonce(ethClient)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Context = context.Background()

	return auth, nil
}

func NewTransactorFromServer(ethClient *EthClient, authClient auth_base.IAuth, signServerUrl string, walletKeyType sign_api.WalletKeyType) (*bind.TransactOpts, error) {

	auth, err := NewKeyedTransactorWithChainIDFromServer(ethClient.ChainID, authClient, signServerUrl)
	if err != nil {

		return nil, errors.Wrap(err, "NewKeyedTransactorWithChainIDFromServer error")
	}

	pendingNonceReq := sign_api.NewGetPendingNonceRequest("")
	res := &sign_api.GetPendingNonceFullResponse{}

	resRes := api.NewRequestResponse(pendingNonceReq, res)
	signUrl, err := url.Parse(signServerUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid sign server URL")
	}

	_, err = api.SendApiRequest(*signUrl, resRes, authClient)
	if err != nil {
		return nil, err
	}

	if res.ReturnCode != int64(foundation.Success) {
		return nil, errors.New("Sign server response error code: " + strconv.Itoa(int(res.ReturnCode)))
	}

	auth.Nonce = big.NewInt(int64(res.Data.PendingNonce))
	gasPrice, err := ethClient.GetGasPrice()
	if err != nil {
		return nil, err
	}
	auth.GasPrice = gasPrice

	return auth, nil
}

func NewKeyedTransactorWithChainIDFromServer(chainID *big.Int, authClient auth_base.IAuth, signServerUrl string) (*bind.TransactOpts, error) {
	if chainID == nil {
		return nil, bind.ErrNoChainID
	}

	signServerUrlObj, err := url.Parse(signServerUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid sign server URL")
	}

	// signer := types.LatestSignerForChainID(chainID)
	return &bind.TransactOpts{
		From: common.HexToAddress("0x0"),
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {

			req := sign_api.NewSignTransactionRequest()
			res := &sign_api.SignTransactionFullResponse{}
			reqRes := api.NewRequestResponse(req, res)
			txBin, err := tx.MarshalBinary()
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("MarshalBinary transaction failed. Nonce: ", req.Nonce, " Error: ", err)
				return nil, err
			}

			req.Data = txBin
			sign, err := eurus_crypto.GenerateRSASignWithRawData(authClient.GetConfig().GetPrivateKey(), txBin)
			if err != nil {
				return nil, err
			}
			req.Sign = sign

			_, err = api.SendApiRequest(*signServerUrlObj, reqRes, authClient)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Failed to sign transaction. Nonce: ", req.Nonce, " Error: ", err)
				return nil, err
			}
			if res.ReturnCode != int64(foundation.Success) {
				log.GetLogger(log.Name.Root).Errorln("Failed to sign transaction with return code: ", res.ReturnCode, " Nonce: ", req.Nonce, " Error: ", err)
				return nil, errors.New("Sign transaction failed. Return code: " + strconv.FormatInt(res.ReturnCode, 10))
			}

			err = tx.UnmarshalBinary(res.Data.SignedTrans)
			if err != nil {
				log.GetLogger(log.Name.Root).Errorln("Failed to UnmarshalBinary transaction with return code: ", res.ReturnCode, " Nonce: ", req.Nonce, " Error: ", err)
				return nil, err
			}
			return tx, nil

		},
		Context: context.Background(),
	}, nil
}

func (me *ECDSAKeyPair) GetNonce(ethClient *EthClient) (uint64, error) {
	nonce, err := ethClient.Client.PendingNonceAt(context.Background(), me.Address)
	return nonce, err
}

func GetEthKeyPair(privateHexKey string) (*ECDSAKeyPair, error) {
	ecdsaKeyPair := &ECDSAKeyPair{}
	privateKey, err := crypto.HexToECDSA(privateHexKey)
	if err != nil {
		return nil, err
	}
	ecdsaKeyPair.PrivateKey = privateKey

	publicKey := ecdsaKeyPair.PrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}
	ecdsaKeyPair.Address = crypto.PubkeyToAddress(*publicKeyECDSA)
	return ecdsaKeyPair, err
}
