package background

import (
	"encoding/json"
	"eurus-backend/config_service/conf_api"
	"eurus-backend/foundation"
	"eurus-backend/foundation/api"
	"eurus-backend/foundation/api/response"
	"eurus-backend/foundation/log"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type TokenPrice struct {
	EthPrice float64 `json:"eth"`
}

//type TokenIdResponse struct {
//	CoingeckoId string `json:"currencyId"`
//}

type TokenIdResponse struct {
	Returncode int    `json:"returnCode"`
	Message    string `json:"message"`
	Nonce      string `json:"nonce"`
	Data       []struct {
		Decimal    int    `json:"Decimal"`
		Assetname  string `json:"AssetName"`
		Currencyid string `json:"CurrencyId"`
	} `json:"data"`
}

func ConstructQueryCoingeckoUrl(ids []string, vsCurreny []string) string {
	baseUrl := "https://api.coingecko.com/api/v3/simple/price?ids="
	for _, e := range ids {
		baseUrl = baseUrl + e + "%2C"
	}
	baseUrl = baseUrl + "&vs_currencies="
	for _, e := range vsCurreny {
		baseUrl = baseUrl + e + "%2C"
	}
	log.GetLogger(log.Name.Root).Debugln("coingecko endpoint: ", baseUrl)
	return baseUrl
}

func QueryConfigServerCoingeckoId(server *BackgroundServer) ([]*conf_api.Asset, error) {

	req := conf_api.NewGetAssetRequest()
	res := conf_api.NewGetAssetResponse()
	reqRes := api.NewRequestResponse(req, res)
	_, err := server.SendConfigApiRequest(reqRes)

	if err != nil || reqRes.Res.GetReturnCode() != int64(foundation.Success) {
		if err != nil {
			log.GetLogger(log.Name.Root).Errorln("get coingecko id from config server error: ", err)
		} else {
			log.GetLogger(log.Name.Root).Errorln("get coingecko id from config server error code: ", reqRes.Res.GetReturnCode(), " message: ", reqRes.Res.GetMessage())
		}
		return nil, err
	}

	var resData []*conf_api.Asset = res.Data
	var tmp []*conf_api.Asset
	for _, e := range resData {
		if e.AutoUpdate {
			tmp = append(tmp, e)
		}
	}

	return tmp, nil
}

func QueryCoingeckoExchangeRate(url string) (map[string]*TokenPrice, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("ioutil.ReadAll error: ", err)
		return nil, err
	}
	var result map[string]*TokenPrice
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		log.GetLogger(log.Name.Root).Errorln("json unmarshal error: ", err)
		return nil, err
	}
	for currency, price := range result {
		log.GetLogger(log.Name.Root).Debugf("currency: %s, exchange rate: %v", currency, price.EthPrice)
	}

	return result, nil
}

func CallConfigServerUpdateExchangeRate(server *BackgroundServer, data map[string]*TokenPrice, assetInfoList []*conf_api.Asset) {

	for currencyId, element := range data {
		log.GetLogger(log.Name.Root).Debugln("Processing currency id: ", currencyId)
		targetAssetInfoList := GetAssetInfoByCurrencyId(currencyId, assetInfoList)

		for _, assetInfo := range targetAssetInfoList {

			for i := 0; i < server.Config.GetRetryCount(); i++ {
				exRateReq := conf_api.NewSetExchangeRate()
				exRateReq.AssetName = assetInfo.AssetName
				exRateReq.Rate = decimal.NewFromFloat(element.EthPrice)
				log.GetLogger(log.Name.Root).Debugln("Sending request to config server. Asset: ", exRateReq.AssetName, " exchange rate: ", exRateReq.Rate.String())
				res := &response.ResponseBase{}
				reqRes := api.NewRequestResponse(exRateReq, res)

				_, err := server.SendConfigApiRequest(reqRes)

				if err != nil || reqRes.Res.GetReturnCode() != int64(foundation.Success) {
					if err != nil {
						log.GetLogger(log.Name.Root).Errorln("call config server update exchange rate error for currency: ", assetInfo.AssetName, " Error: ", err)
					} else {
						log.GetLogger(log.Name.Root).Errorln("call config server update exchange rate error for currency: ", assetInfo.AssetName, " response code: ", reqRes.Res.GetReturnCode(), " Error: ", reqRes.Res.GetMessage())
					}
					log.GetLogger(log.Name.Root).Debugln("Now retry ", i, " times")
					time.Sleep(server.Config.GetRetryInterval() * time.Second)
					continue
				}

				break
			}
		}
	}
	log.GetLogger(log.Name.Root).Infoln("Update finished")
}

func GetAssetInfoByCurrencyId(currencyId string, assetInfoList []*conf_api.Asset) []*conf_api.Asset {
	var outputList []*conf_api.Asset = make([]*conf_api.Asset, 0)
	for _, assetInfo := range assetInfoList {
		if assetInfo.CurrencyId == currencyId {
			outputList = append(outputList, assetInfo)
		}
	}

	return outputList
}
