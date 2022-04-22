package conf

import (
	"eurus-backend/config_service/conf_api"
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"
)

func getCurrencyInfoInArrayFromDB(me *ConfigServer, formatValue *big.Int) ([]string, []*big.Int) {
	dbConn, err := me.SlaveDatabase.GetConn()
	if err != nil {
		fmt.Println(err)
	}
	assetList := make([]*conf_api.Asset, 0)

	tx := dbConn.Where("decimal > 0").Find(&assetList)
	err = tx.Error

	var assetArray []string
	var assetAmount []*big.Int
	for _, e := range assetList {
		exchangeRateInDB, err := getExchangeRateFromDB(me, e.AssetName)
		if err != nil {
			continue
		}
		amount, err := CalculateAdminFeeAmount(me, formatValue, exchangeRateInDB, e.Decimal)
		if err != nil {
			continue
		}
		assetArray = append(assetArray, e.AssetName)
		assetAmount = append(assetAmount, amount)
	}
	return assetArray, assetAmount

}

func getExchangeRateFromDB(me *ConfigServer, assetName string) (decimal.Decimal, error) {
	dbConn, err := me.SlaveDatabase.GetConn()
	if err != nil {
		fmt.Println(err)
	}

	exchangeRate := new(conf_api.ExchangeRates)

	tx := dbConn.Where("asset_name = ?", assetName).Find(exchangeRate)
	err = tx.Error
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	return exchangeRate.Rate, nil
}
