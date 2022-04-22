
import axios from 'axios'
import { getUserId, getWalletAddress } from "../store";
import {
  getAdminApiUrl,
  getUsdtSmartContractAddress,
  getDappSmartContractAddress
} from '../network';

export async function requestRefund(amount) {
  let result = null;

  try {
    var request = {
      fromWalletAddr: getDappSmartContractAddress(),
      toWalletAddr: getWalletAddress(),
      toAssetId: "USDT",
      toAssetAddr: getUsdtSmartContractAddress(),
      fromTokenAmt: amount,
      rate: "1",
      userId: getUserId()
    };

    let url = getAdminApiUrl() + "/withdraw-reqs";
    const headers = {
      "Content-Type": "application/json",
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

export async function getUserBalance() {
  let result = null;

  try {
    let url = getAdminApiUrl() + "/merchant-clients/" + getUserId();
    const headers = {
      "Content-Type": "application/json",
    };
    let response = await axios.get(url, { headers: headers });
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
