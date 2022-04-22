package implementation

import (
	"dapp-payment-api/config"
	"dapp-payment-api/database"
	"dapp-payment-api/oapi"
	"dapp-payment-api/signature"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type PaymentAPI struct {
	config *config.ServerConfig
	logger *logrus.Logger
	db     *gorm.DB

	mutex                              sync.RWMutex
	_networksByID                      map[int64]*database.DBNetwork
	_networksByCode                    map[string]*database.DBNetwork
	_tokensByID                        map[int64]*database.DBToken
	_tokensByNetworkCodeAndSymbol      map[string]map[string]*database.DBToken
	_merchantsByID                     map[int64]*database.DBMerchant
	_merchantsByCode                   map[string]*database.DBMerchant
	_merchantWalletsByMerchantAndToken map[int64]map[int64]*database.DBMerchantWallet
}

func NewPaymentAPI(c *config.ServerConfig, l *logrus.Logger) *PaymentAPI {
	if c == nil || l == nil {
		panic("config and logger cannot be nil")
	}

	ret := new(PaymentAPI)
	ret.config = c
	ret.logger = l

	err := ret.dbOpen()
	panicIfError(l, err)

	ret.cacheData()
	go func() {
		l.Info("Start go routine to update cache data every ", c.RefreshCacheSeconds, " second(s)")
		for {
			time.Sleep(time.Duration(c.RefreshCacheSeconds) * time.Second)
			ret.cacheData()
		}
	}()

	return ret
}

func (p *PaymentAPI) GetAllMerchants(w http.ResponseWriter, r *http.Request) {
	// This call just returns simply information, for details use the call with merchantCode specified
	merchants := []*oapi.Merchant{}
	for _, dbm := range p.merchantsByID() {
		m := new(oapi.Merchant)
		m.MerchantID = dbm.ID
		m.MerchantCode = dbm.MerchantCode
		m.MerchantName = dbm.MerchantName
		merchants = append(merchants, m)
	}

	ret := make(map[string]interface{})
	ret["merchants"] = merchants

	OK(w, r, ret)
}

func (p *PaymentAPI) GetMerchant(w http.ResponseWriter, r *http.Request, merchantCode string) {
	dbMerchant, found := p.merchantsByCode()[merchantCode]
	if !found {
		NotFound(w, r)
		return
	}

	ret := oapi.Merchant{}
	ret.MerchantID = dbMerchant.ID
	ret.MerchantCode = dbMerchant.MerchantCode
	ret.MerchantName = dbMerchant.MerchantName

	// TagDisplayName from DB is a JSON
	if dbMerchant.TagDisplayName != nil {
		tagDisplayName := new(map[string]interface{})
		err := json.Unmarshal([]byte(*dbMerchant.TagDisplayName), tagDisplayName)
		if err != nil {
			p.logger.Error(err.Error())
			InternalServerError(w, r)
			return
		}

		ret.TagDisplayName = tagDisplayName
	} else {
		ret.TagDisplayName = nil
	}
	ret.TagDescription = dbMerchant.TagDescription

	dbMerchantWallets, found := p.merchantWalletsByMerchantAndToken()[dbMerchant.ID]
	if !found {
		// No wallets are given to this merchant, so just return an empty map
		dbMerchantWallets = make(map[int64]*database.DBMerchantWallet)
	}

	ret.Wallets = &[]oapi.MerchantWallet{}
	for _, dbw := range dbMerchantWallets {
		w := oapi.MerchantWallet{}
		dbToken := p.tokensByID()[dbw.TokenID]
		w.NetworkCode = p.networksByID()[dbToken.NetworkID].NetworkCode
		w.Symbol = dbToken.Symbol
		w.Address = dbw.Address
		*ret.Wallets = append(*ret.Wallets, w)
	}

	OK(w, r, ret)
}

func (p *PaymentAPI) GetTransactionsByMerchant(w http.ResponseWriter, r *http.Request, merchantCode string, params oapi.GetTransactionsByMerchantParams) {
	validated, err := p.verifyAPIKey(merchantCode, r.Header.Get("X-API-Key"))
	if err != nil {
		p.logger.Error(err.Error())
		InternalServerError(w, r)
		return
	}

	if !validated {
		Unauthorized(w, r)
		return
	}

	ret := oapi.MerchantQueryReply{}
	ret.Data = []oapi.MerchantTransaction{}

	dbTransactions, err := p.dbGetTransactionsStartFrom(p.merchantsByCode()[merchantCode].ID, params.StartingSeqNo, params.Limit)
	if err != nil {
		p.logger.Error(err.Error())
		InternalServerError(w, r)
		return
	}

	if len(*dbTransactions) != 0 {
		for _, dbt := range *dbTransactions {
			t := oapi.MerchantTransaction{}
			t.SeqNo = dbt.MerchantSeqNo
			t.NetworkCode = p.networksByID()[dbt.NetworkID].NetworkCode
			t.Symbol = p.tokensByID()[dbt.TokenID].Symbol
			t.Amount = dbt.Amount.String()
			t.Tag = dbt.Tag
			t.TxHash = dbt.TxHash
			t.BlockHash = dbt.BlockHash
			t.BlockNo = dbt.BlockNumber
			t.OnchainStatus = dbt.OnchainStatus
			t.ConfirmStatus = dbt.ConfirmStatus
			ret.Data = append(ret.Data, t)
		}

		last := ret.Data[len(ret.Data)-1].SeqNo
		dbMerchant, err := p.dbGetMerchant(merchantCode)
		if err != nil {
			p.logger.Error(err.Error())
			InternalServerError(w, r)
			return
		}

		if last < dbMerchant.MerchantLastSeq {
			ret.Meta.HasNext = true
			ret.Meta.NextSeqNo = new(int64)
			*ret.Meta.NextSeqNo = last + 1
		}
	}

	OK(w, r, ret)
}

func (p *PaymentAPI) GetAllNetworks(w http.ResponseWriter, r *http.Request) {
	networks := []*oapi.Network{}
	for _, dbn := range p.networksByID() {
		n := new(oapi.Network)
		n.NetworkID = dbn.ID
		n.NetworkCode = dbn.NetworkCode
		n.NetworkName = dbn.NetworkName
		n.ChainID = dbn.ChainID
		networks = append(networks, n)
	}

	ret := make(map[string]interface{})
	ret["networks"] = networks

	OK(w, r, ret)
}

func (p *PaymentAPI) GetAllTokens(w http.ResponseWriter, r *http.Request, networkCode string) {
	dbNetwork, found := p.networksByCode()[networkCode]
	if !found {
		NotFound(w, r)
		return
	}

	tokens := []*oapi.Token{}
	for _, dbt := range p.tokensByID() {
		if dbt.NetworkID != dbNetwork.ID {
			continue
		}

		t := new(oapi.Token)
		t.TokenID = dbt.ID
		t.NetworkCode = networkCode
		t.Address = dbt.Address
		t.Symbol = dbt.Symbol
		t.Name = dbt.Name
		t.Decimals = dbt.Decimals
		tokens = append(tokens, t)
	}

	ret := make(map[string]interface{})
	ret["tokens"] = tokens

	OK(w, r, ret)
}

func (p *PaymentAPI) GetSubmission(w http.ResponseWriter, r *http.Request, networkCode string, txHash string) {
	// Always use lower case in DB or return value
	txHash = strings.ToLower(txHash)

	// If the network code is invalid, return HTTP 400 instead of 404
	dbNetwork, found := p.networksByCode()[networkCode]
	if !found {
		BadRequest(w, r)
		return
	}

	// Users are allowed to submit again only if the previous one is failed, so here just return the last one
	dbSubmission, err := p.dbGetLatestSubmission(dbNetwork.ID, txHash)
	if err != nil {
		p.logger.Error(err.Error())
		InternalServerError(w, r)
		return
	}

	// If no one had submitted correct signature for this transaction, return 404
	if dbSubmission == nil {
		NotFound(w, r)
		return
	}

	ret := oapi.SubmissionReply{}
	ret.TxHash = txHash
	ret.SubmitTime = dbSubmission.SubmitTime
	ret.NetworkCode = networkCode
	ret.Symbol = p.tokensByID()[dbSubmission.TokenID].Symbol
	ret.FromAddress = dbSubmission.FromAddress
	ret.Amount = dbSubmission.Amount.String()
	ret.TxStatus = dbSubmission.TxStatus
	ret.PaymentStatus = dbSubmission.PaymentStatus

	OK(w, r, ret)
}

func (p *PaymentAPI) SubmitPaymentTransaction(w http.ResponseWriter, r *http.Request, networkCode string, txHash string) {
	// Always use lower case in DB or return value
	txHash = strings.ToLower(txHash)

	// Read the body and deserialize it, the raw message will be needed at the end so do not use the wrapper from chi render
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.logger.Error(err.Error())
		InternalServerError(w, r)
		return
	}

	submission := new(oapi.Submission)
	err = json.Unmarshal(body, submission)
	if err != nil {
		p.logger.Error(err.Error())
		InternalServerError(w, r)
		return
	}

	// Validate the input if network code or merchant code is invalid
	dbNetwork, found := p.networksByCode()[networkCode]
	if !found {
		BadRequest(w, r)
		return
	}

	dbToken, found := p.tokensByNetworkCodeAndSymbol()[networkCode][submission.Coin]
	if !found {
		BadRequest(w, r)
		return
	}

	dbMerchant, found := p.merchantsByCode()[submission.Merchant]
	if !found {
		BadRequest(w, r)
		return
	}

	amount := big.NewInt(0)
	_, success := amount.SetString(submission.Amount, 10)
	if !success {
		BadRequest(w, r)
		return
	}

	sig, err := signature.ReadSignature(submission.Signature)
	if err != nil {
		p.logger.Error(err.Error())
		InternalServerError(w, r)
		return
	}

	signer, err := signature.GetSigner(
		*dbNetwork.ChainID,
		common.HexToHash(txHash),
		submission.Coin,
		amount,
		submission.Merchant,
		submission.Tag,
		sig)
	if err != nil {
		p.logger.Error(err.Error())
		InternalServerError(w, r)
		return
	}

	expect := common.HexToAddress(submission.From)
	if signer != expect {
		BadRequest(w, r)
		return
	}

	// If no submission for this transaction before, just insert it
	// However, if there are submissions before, depend on the last's status, determine if this submission should be allowed
	// Therefore here wraps the actions in transaction, and lock the table to prevent race condition
	statusCode, err := p.dbInsertSubmission(dbNetwork.ID, dbToken.ID, dbMerchant.ID, txHash, submission, amount, string(body))
	if err != nil {
		p.logger.Error(err.Error())
		InternalServerError(w, r)
		return
	}

	switch statusCode {
	case 201:
		Created(w, r)
	case 403:
		Forbidden(w, r)
	case 409:
		Conflict(w, r)
	case 410:
		Gone(w, r)
	default:
		InternalServerError(w, r)
	}
}

func panicIfError(l *logrus.Logger, err error) {
	if err != nil {
		l.Fatal(err.Error())
		panic(err)
	}
}

func (p *PaymentAPI) networksByID() map[int64]*database.DBNetwork {
	p.mutex.RLock()
	defer func() { p.mutex.RUnlock() }()
	return p._networksByID
}

func (p *PaymentAPI) networksByCode() map[string]*database.DBNetwork {
	p.mutex.RLock()
	defer func() { p.mutex.RUnlock() }()
	return p._networksByCode
}

func (p *PaymentAPI) tokensByID() map[int64]*database.DBToken {
	p.mutex.RLock()
	defer func() { p.mutex.RUnlock() }()
	return p._tokensByID
}

func (p *PaymentAPI) tokensByNetworkCodeAndSymbol() map[string]map[string]*database.DBToken {
	p.mutex.RLock()
	defer func() { p.mutex.RUnlock() }()
	return p._tokensByNetworkCodeAndSymbol
}

func (p *PaymentAPI) merchantsByID() map[int64]*database.DBMerchant {
	p.mutex.RLock()
	defer func() { p.mutex.RUnlock() }()
	return p._merchantsByID
}

func (p *PaymentAPI) merchantsByCode() map[string]*database.DBMerchant {
	p.mutex.RLock()
	defer func() { p.mutex.RUnlock() }()
	return p._merchantsByCode
}

func (p *PaymentAPI) merchantWalletsByMerchantAndToken() map[int64]map[int64]*database.DBMerchantWallet {
	p.mutex.RLock()
	defer func() { p.mutex.RUnlock() }()
	return p._merchantWalletsByMerchantAndToken
}

func (p *PaymentAPI) cacheData() {
	p.logger.Info("Update cache data")

	dbNetworks, err := p.dbGetAllNetworks()
	if err != nil {
		p.logger.Error(err.Error())
		return
	}

	dbTokens, err := p.dbGetAllTokens()
	if err != nil {
		p.logger.Error(err.Error())
		return
	}

	dbMerchants, err := p.dbGetAllMerchants()
	if err != nil {
		p.logger.Error(err.Error())
		return
	}

	dbMerchantWallets, err := p.dbGetAllMerchantWallets()
	if err != nil {
		p.logger.Error(err.Error())
		return
	}

	// Create new maps, then replace those existing at once
	networksByID := make(map[int64]*database.DBNetwork)
	networksByCode := make(map[string]*database.DBNetwork)
	for _, dbNetwork := range *dbNetworks {
		networksByID[dbNetwork.ID] = dbNetwork
		networksByCode[dbNetwork.NetworkCode] = dbNetwork
	}

	tokensByID := make(map[int64]*database.DBToken)
	tokensByNetworkCodeAndSymbol := make(map[string]map[string]*database.DBToken)
	for _, dbToken := range *dbTokens {
		networkID := dbToken.NetworkID
		networkCode := networksByID[networkID].NetworkCode
		tokenID := dbToken.ID
		symbol := dbToken.Symbol

		tokensByID[tokenID] = dbToken

		if _, found := tokensByNetworkCodeAndSymbol[networkCode]; !found {
			tokensByNetworkCodeAndSymbol[networkCode] = make(map[string]*database.DBToken)
		}
		tokensByNetworkCodeAndSymbol[networkCode][symbol] = dbToken
	}

	merchantsByID := make(map[int64]*database.DBMerchant)
	merchantsByCode := make(map[string]*database.DBMerchant)
	for _, dbMerchant := range *dbMerchants {
		merchantsByID[dbMerchant.ID] = dbMerchant
		merchantsByCode[dbMerchant.MerchantCode] = dbMerchant
	}

	merchantWalletsByMerchantAndToken := make(map[int64]map[int64]*database.DBMerchantWallet)
	for _, dbMerchantWallet := range *dbMerchantWallets {
		merchantID := dbMerchantWallet.MerchantID
		tokenID := dbMerchantWallet.ID

		if _, found := merchantWalletsByMerchantAndToken[merchantID]; !found {
			merchantWalletsByMerchantAndToken[merchantID] = make(map[int64]*database.DBMerchantWallet)
		}
		merchantWalletsByMerchantAndToken[merchantID][tokenID] = dbMerchantWallet
	}

	p.mutex.Lock()
	defer func() { p.mutex.Unlock() }()

	p._networksByID = networksByID
	p._networksByCode = networksByCode
	p._tokensByID = tokensByID
	p._tokensByNetworkCodeAndSymbol = tokensByNetworkCodeAndSymbol
	p._merchantsByID = merchantsByID
	p._merchantsByCode = merchantsByCode
	p._merchantWalletsByMerchantAndToken = merchantWalletsByMerchantAndToken
}

func (p *PaymentAPI) verifyAPIKey(merchantCode string, apiKey string) (bool, error) {
	// Merchant can have multiple api keys so check them one by one
	dbMerchantAPIKeys, err := p.dbGetAPIKeys(merchantCode)
	if err != nil {
		return false, err
	}

	for _, dbm := range *dbMerchantAPIKeys {
		// If the api key saved in DB does not match PHC string format, just skip it
		parts := strings.Split(dbm.APIKey, "$")
		if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
			continue
		}

		var version int
		if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
			continue
		}

		var memory uint32
		var time uint32
		var parallelism uint8
		if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &parallelism); err != nil {
			continue
		}

		salt := []byte(merchantCode + "." + dbm.Salt)
		if base64.RawStdEncoding.EncodeToString(salt) != parts[4] {
			continue
		}

		idKey := argon2.IDKey([]byte(apiKey), salt, time, memory, parallelism, 32)
		if base64.RawStdEncoding.EncodeToString(idKey) == parts[5] {
			return true, nil
		}
	}

	return false, nil
}
