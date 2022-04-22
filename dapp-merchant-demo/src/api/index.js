
import axios from 'axios'
// import ethCrypto from "eth-crypto";
import { uuidv4 } from "@/utils";
import {
  setEurusChainId,
  setEurusRpcDomain,
  setEurusRpcPort,
  getToken,
} from '@/utils/auth'
import { amountMultipleDecimals } from "@/web3";

const eurusApiUrl = process.env.VUE_APP_EURUS_RPC_URL

export async function getServerConfig() {
  let result = false
  try {
    const url = eurusApiUrl + '/user/serverConfig'
    const headers = {
      'Content-Type': 'application/x-www-form-urlencoded'
    }
    const response = await axios.get(url, { headers: headers })
    if (response && response.status === 200) {
      if (response.data && response.data.data) {
        const data = response.data.data
        if (
          data.eurusRPCDomain &&
          data.eurusPRCProtocol &&
          data.eurusRPCPort &&
          data.eurusChainId
        ) {
          setEurusRpcDomain(data.eurusPRCProtocol + '://' + data.eurusRPCDomain)
          setEurusRpcPort(data.eurusRPCPort)
          setEurusChainId(data.eurusChainId)
          result = true
        }
      }
    }
  } catch (error) {
    console.error(error)
    result = false
  }
  return result
}

export async function requestRefund(asset, amount, reason) {
  let result = null;

  try {
    var request = {
      nonce: uuidv4(),
      merchantId: 1,
      assetName: asset,
      amount: amountMultipleDecimals(amount, 6),
      reason: reason
    };

    let url = eurusApiUrl + "/user/merchant/requestRefund";
    const headers = {
      Authorization: "Bearer " + getToken(),
      "Content-Type": "application/x-www-form-urlencoded",
    };
    console.log("requestRefund request: ", request);
    let response = await axios.post(url, request, { headers: headers });
    if (response && response.status === 200) {
      if (response.data) {
        let data = response.data;
        result = data;
      }
    }
  } catch (error) {
    console.error(error);
  }

  return result;
}