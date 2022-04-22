package implementation

import (
	"dapp-payment-api/database"
	"dapp-payment-api/oapi"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (p *PaymentAPI) dbOpen() error {
	dsn := fmt.Sprintf("host='%v' port='%v' user='%v' password='%v' dbname='%v'",
		pgEscape(p.config.Database.Host),
		p.config.Database.Port,
		pgEscape(p.config.Database.User),
		pgEscape(p.config.Database.Password),
		pgEscape(p.config.Database.Name))

	var err error
	p.db, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.New(
				p.logger,
				logger.Config{
					SlowThreshold:             200 * time.Millisecond,
					Colorful:                  false,
					IgnoreRecordNotFoundError: false,
					LogLevel:                  logger.Warn})})
	if err != nil {
		return err
	}

	dbo, err := p.db.DB()
	if err != nil {
		return err
	}

	dbo.SetMaxIdleConns(10)
	dbo.SetMaxOpenConns(100)

	return nil
}

func (p *PaymentAPI) dbGetAllNetworks() (*[]*database.DBNetwork, error) {
	var dbNetworks []*database.DBNetwork
	tx := p.db.Table("payment.t_networks").Order("id ASC").Find(&dbNetworks)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &dbNetworks, nil
}

func (p *PaymentAPI) dbGetAllTokens() (*[]*database.DBToken, error) {
	var dbTokens []*database.DBToken
	tx := p.db.Table("payment.t_tokens").Order("id ASC").Find(&dbTokens)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &dbTokens, nil
}

func (p *PaymentAPI) dbGetAllMerchants() (*[]*database.DBMerchant, error) {
	var dbMerchants []*database.DBMerchant
	tx := p.db.Table("payment.t_merchants").Order("id ASC").Find(&dbMerchants)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &dbMerchants, nil
}

func (p *PaymentAPI) dbGetMerchant(merchantCode string) (*database.DBMerchant, error) {
	var dbMerchant database.DBMerchant
	tx := p.db.Table("payment.t_merchants").Where("merchant_code = ?", merchantCode).Find(&dbMerchant)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected != 1 {
		return nil, nil
	}

	return &dbMerchant, nil
}

func (p *PaymentAPI) dbGetAllMerchantWallets() (*[]*database.DBMerchantWallet, error) {
	var dbMerchantWallets []*database.DBMerchantWallet
	tx := p.db.Table("payment.t_merchant_wallets").Order("ID ASC").Find(&dbMerchantWallets)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &dbMerchantWallets, nil
}

func (p *PaymentAPI) dbGetMerchantWallets(merchantID int64) (*[]*database.DBMerchantWallet, error) {
	var dbMerchantWallets []*database.DBMerchantWallet
	tx := p.db.Table("payment.t_merchant_wallets").Where("merchant_id = ?", merchantID).Order("ID ASC").Find(&dbMerchantWallets)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &dbMerchantWallets, nil
}

func (p *PaymentAPI) dbGetAPIKeys(merchantCode string) (*[]*database.DBMerchantAPIKey, error) {
	var dbMerchantAPIKeys []*database.DBMerchantAPIKey
	tx := p.db.Raw(`
		SELECT		U.created_at, U.updated_at, U.merchant_id, U.api_key, U.salt
		FROM		payment.t_merchants			T
		INNER JOIN	payment.t_merchant_api_keys	U
		ON			T.id			= U.merchant_id
		WHERE		T.merchant_code	= ?`, merchantCode).Find(&dbMerchantAPIKeys)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &dbMerchantAPIKeys, nil
}

func (p *PaymentAPI) dbGetLatestSubmission(networkID int64, txHash string) (*database.DBSubmission, error) {
	var dbSubmission database.DBSubmission
	tx := p.db.Table("payment.t_submissions").Where("network_id = ? AND tx_hash = ?", networkID, txHash).Order("submit_time DESC").Limit(1).Find(&dbSubmission)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, nil
	}

	return &dbSubmission, nil
}

func (p *PaymentAPI) dbInsertSubmission(networkID int64, tokenID int64, merchantID int64, txHash string, submission *oapi.Submission, amount *big.Int, message string) (int, error) {
	statusCode := 201
	err := p.db.Transaction(func(_tx *gorm.DB) error {
		tx := _tx.Exec("LOCK TABLE payment.t_submissions IN EXCLUSIVE MODE")
		if tx.Error != nil {
			return tx.Error
		}

		txStatus := -1
		var dbSubmissionLast database.DBSubmission
		tx = _tx.Table("payment.t_submissions").Where("network_id = ? AND tx_hash = ?", networkID, txHash).Order("submit_time DESC").Limit(1).Find(&dbSubmissionLast)
		if tx.Error != nil {
			return tx.Error
		}

		if tx.RowsAffected != 0 {
			// Transaction failed, no more submission for it
			if dbSubmissionLast.TxStatus == 0 {
				statusCode = 410
				return nil
			}

			txStatus = dbSubmissionLast.TxStatus

			switch dbSubmissionLast.PaymentStatus {
			case 0, 1:
				// Still not processed
				statusCode = 403
				return nil
			case 2:
				// Already verified with correct information
				statusCode = 409
				return nil
			}
		}

		dbSubmission := database.DBSubmission{}
		dbSubmission.SubmitTime = time.Now()
		dbSubmission.NetworkID = networkID
		dbSubmission.TokenID = tokenID
		dbSubmission.FromAddress = strings.ToLower(submission.From)
		dbSubmission.Amount = decimal.NewFromBigInt(amount, 0)
		dbSubmission.MerchantID = merchantID
		dbSubmission.Tag = submission.Tag
		dbSubmission.TxHash = txHash
		dbSubmission.TxStatus = txStatus
		dbSubmission.PaymentStatus = 0
		dbSubmission.Signature = strings.ToLower(submission.Signature)
		dbSubmission.MessageBody = string(message)

		tx = _tx.Table("payment.t_submissions").Create(&dbSubmission)
		if tx.Error != nil {
			return tx.Error
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return statusCode, nil
}

func (p *PaymentAPI) dbGetTransactionsStartFrom(merchantID int64, startingSeqNo int64, limit *int) (*[]*database.DBTransaction, error) {
	var dbTransactions []*database.DBTransaction
	tx := p.db.Table("payment.t_transactions").Where("merchant_id = ? AND merchant_seq_no >= ?", merchantID, startingSeqNo).Order("merchant_seq_no ASC")
	if limit != nil {
		tx = tx.Limit(*limit)
	}
	tx = tx.Find(&dbTransactions)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &dbTransactions, nil
}

func pgEscape(s string) string {
	ret := strings.ReplaceAll(s, "\\", "\\\\")
	ret = strings.ReplaceAll(ret, "'", "\\'")
	return ret
}
