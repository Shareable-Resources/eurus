package test

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
)

func Test_CalculateAdminFeeAmount(t *testing.T) {
	ethFee := big.NewInt(15000000000000000)
	inputRate := decimal.NewFromFloat32(0.00038424)
	var dec int64 = 18

	newAmount, _ := CalculateAdminFeeAmount(ethFee, inputRate, dec)
	fmt.Println(newAmount.String())

}
func CalculateAdminFeeAmount(ethFee *big.Int, inputRate decimal.Decimal, Decimal int64) (*big.Int, error) {
	//empty := big.NewInt(0)
	//dbConn, err := me.DefaultDatabase.GetConn()

	//getDecimal := new(Assets)
	//tx1 := dbConn.Where("currency_id = ?", AssetName).Find(getDecimal)
	//err = tx1.Error
	//if err != nil   {
	//	return empty, "" ,err
	//}
	//if tx1.RowsAffected == 0 {
	//	err := errors.New("Can not find this currency in asset table")
	//	return empty, "" , err
	//}

	FormatWithdrawFee := new(big.Float)
	FormatWithdrawFee.Quo(new(big.Float).SetInt(ethFee), new(big.Float).SetFloat64(math.Pow(10, 18)))

	targetAdminFee := new(big.Float)
	targetAdminFee.Quo(FormatWithdrawFee, inputRate.BigFloat()) //admin fee / rate = new Currency amount (unit in ETH(gel))

	newDecimal := new(big.Float).SetFloat64(math.Pow(10, float64(Decimal)))

	formattedRate := new(big.Float).Mul(newDecimal, targetAdminFee)

	targetAmount, _ := formattedRate.Int(nil)

	return targetAmount, nil

}

// func TestAddAssetAndGetter(t *testing.T) {
// 	configServer := conf.NewConfigServer()
// 	configServer.InitLog("ConfigServer.log")
// 	serverConfig := new(server.ServerConfigBase)
// 	err := configServer.LoadConfig("", serverConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	configServer.ServerConfig = serverConfig
// 	configServer.InitDBFromConfig(configServer.ServerConfig)
// 	configServer.InitSC()

// 	_, err = configServer.InitEthereumClientFromConfig(configServer.ServerConfig)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	txHash, err := conf.DelCurrencyInfoFromSC(configServer, "USDT")
// 	if err != nil {
// 		t.Error(err.Error())
// 	} else {
// 		t.Logf("txHash: %v", txHash)
// 	}

// 	txHash, err = conf.AddCurrencyInfoFromSC(configServer, "0x84DfaaBF9fD8E72764E4997Fbc775758E00a82f7", "USDT")
// 	if err != nil {
// 		t.Error(err.Error())
// 	} else {
// 		t.Logf("txHash: %v", txHash)
// 	}
// 	addr, err := conf.GetCurrencySCAddrFromSC(configServer, "USDT")
// 	if err != nil {
// 		t.Error(err.Error())
// 	} else {
// 		t.Logf("addr: %v", addr)
// 	}
// 	asset, err := conf.GetCurrencyNameByAddrFromSC(configServer, "0x84DfaaBF9fD8E72764E4997Fbc775758E00a82f7")
// 	if err != nil {
// 		t.Error(err.Error())
// 	} else {
// 		t.Logf("asset: %v", asset)
// 	}
// }
