/* eslint-disable */
import Web3 from 'web3';
import { ethers } from 'ethers';
import EthCrypto from 'eth-crypto';
import forge from 'node-forge';
import cryptoJS from 'crypto-js';
import elliptic from 'elliptic';
const secp256k1 = new elliptic.ec('secp256k1');

import EurusERC20 from '../contractsAbi/EurusERC20.json';
import DAppStockSample from '../contractsAbi/DAppStockSample.json';
import axios from 'axios'
import { uuidv4 } from '@/utils';
import {
  getUsdtAddress
} from '@/utils/auth';

import {
  getSideChainNetworks,
  getEurusApiUrl,
} from './network';

const { ethereum } = window;

export function checkDappBrowser() {
  if (ethereum) {
    return true;
  } else {
    return false;
  }
}

export async function getChainId() {
  let net_version = null;
  try {
    net_version = await ethereum.request({ method: 'eth_chainId' });
    net_version = parseInt(net_version);
    console.log('net_version: ', net_version);
  } catch (error) {
    console.log('getChainId error: ', error);
  }
  return net_version
}

export async function getAccount() {
  return ethereum.request({ method: 'eth_requestAccounts' })
}

export async function addEthereumChain() {
  let result = false;

  let rpcUrl = getSideChainNetworks().url;
  let chainId = getSideChainNetworks().chainId;
  let chainName = getSideChainNetworks().name;
  let web3 = new Web3(rpcUrl);
  let chainIdHex = web3.utils.numberToHex(chainId);
  let nativeCurrency = 'EUN';
  
  const networkObject = {
    method: 'wallet_addEthereumChain',
    params: [{
      chainId: chainIdHex, // A 0x-prefixed hexadecimal string
      chainName: chainName,
      nativeCurrency: {
        name: nativeCurrency,
        symbol: nativeCurrency, // 2-6 characters long
        decimals: 18,
      },
      rpcUrls: [rpcUrl],
    }],
  }
  console.log('#### addEthereumChain networkObject:', networkObject);

  try {
    await ethereum.request(networkObject);
    result = true;
  } catch (addError) {
    console.log('#### addEthereumChain addError:', addError);
  }
  console.log('#### addEthereumChain result:', result);
  return result;
}

export async function switchEthereumChain() {
  let result = false;
  let rpcUrl = getSideChainNetworks().url;
  let chainId = getSideChainNetworks().chainId;
  let web3 = new Web3(rpcUrl);
  let chainIdHex = web3.utils.numberToHex(chainId);
  
  console.log('#### switchEthereumChain chainId:', chainId);
  console.log('#### switchEthereumChain chainIdHex:', chainIdHex);

  try {
    await ethereum.request({
      method: 'wallet_switchEthereumChain',
      params: [{ chainId: chainIdHex }],
    });
  } catch (switchError) {
    // This error code indicates that the chain has not been added to MetaMask.
    console.log('#### switchEthereumChain switchError:', switchError);
    if (switchError.code === 4902) {
      await addEthereumChain(network);
    }
    // handle other "switch" errors
  }
  let chainId2 = await getChainId();
  if (chainId2 == chainId){
    result = true;
  }
  console.log('#### switchEthereumChain result:', result);
  return result;
}

export async function personalSign(walletAddress, message) {
    let signature = '';
    try {
      const ethResult = await ethereum.request({
        method: 'personal_sign',
        params: [walletAddress, message],
      })
      signature = ethResult;
    } catch (err) {
      console.error('#### sign err: ', err);
    }
    return signature;
  }

export async function importWallet(inputWalletAddress) {
  let url = getEurusApiUrl() + '/user/importWallet';
  const nonce = uuidv4();
  const ts = Date.now();
  const deviceId = "";
  const walletAddress = inputWalletAddress.toLowerCase().substring(2);
  const message = 'deviceId=' + deviceId + '&timestamp=' + ts + '&walletAddress=' + walletAddress;
  // const messageHash = EthCrypto.hash.keccak256(message);
  // let signature = await sign(inputWalletAddress, messageHash);

  const rpcUrl = getSideChainNetworks().url;
  const web3Utils = (new Web3(rpcUrl)).utils;
  const initialBuffer = Buffer.from(web3Utils.hexToBytes(web3Utils.utf8ToHex(message), 'hex'));
  const preamble = Buffer.from(`\x19Ethereum Signed Message:\n${initialBuffer.length}`);
  const messageBuffer = Buffer.concat([preamble, initialBuffer]);
  const messageHash = web3Utils.sha3(messageBuffer);

  let signature = await personalSign(inputWalletAddress, message);

  let publicKey = '';
  if (signature) {
    const recoverPublicKey = EthCrypto.recoverPublicKey(
      signature, // signature
      messageHash // message hash
    );
    publicKey = EthCrypto.publicKey.compress(recoverPublicKey)
    signature = signature.substring(2, 130);
  }
  let body = {
    "nonce": nonce,
    "timestamp": ts,
    "walletAddress": walletAddress,
    "publicKey": publicKey,
    "deviceId": "",
    "sign": signature,
    "isPersonalSign": true,
  };

  const headers = {
    "Content-Type": "application/x-www-form-urlencoded",
  };

  return axios.post(url, body, { headers: headers });
}

export function amountDivideDecimals(amount, decimals) {
  let output = 0;
  try {
    output = amount / Math.pow(10, decimals);
  } catch (err) {
    // amountDivideDecimals error
  }
  return output;
}

export function amountMultipleDecimals(amount, decimals) {
  let output = 0;
  try {
    output = parseFloat(amount) * Math.pow(10, decimals)
  } catch (err) {
    // amountMultipleDecimals error
  }
  return output;
}

export async function getErc20Balance(walletAddress, contractAddress) {
  let contractJson = EurusERC20
  let url = getSideChainNetworks().url
  let web3 = new Web3(url)
  var tokenContract = new web3.eth.Contract(contractJson.abi, contractAddress)
  let balance = await tokenContract.methods.balanceOf(walletAddress).call()
  let decimals = await tokenContract.methods.decimals().call()
  // let currencyCode = await tokenContract.methods.symbol().call();
  // let output = amountDivideDecimals(balance, decimals)
  return amountDivideDecimals(balance, decimals)

}

export async function depositToDapp(contractAddress, dappAddress, amount) {
  console.log("depositToDapp start")

  const url = getSideChainNetworks().url
  const chainId = getSideChainNetworks().chainId
  const net_version = await getChainId()
  if (net_version != chainId) {
    throw 9991
  }

  const accounts = await getAccount();
  const walletAddress = accounts[0]

  const contractJson = EurusERC20
  
  let web3 = new Web3(url)
  var tokenContract = new web3.eth.Contract(contractJson.abi, contractAddress)
  let decimals = await tokenContract.methods.decimals().call()
  let depositBalance = amountMultipleDecimals(amount, decimals)

  let dappAddressData = web3.utils.padLeft(dappAddress, 64)
  let amountData = web3.utils.padLeft(web3.utils.numberToHex(depositBalance.toString()), 64)
  let extraData = web3.eth.abi.encodeParameter('bytes32', '0x0');

  dappAddressData = dappAddressData.substr(2)
  amountData = amountData.substr(2)
  extraData = extraData.substr(2)
  
  
  let txnData = '0x3a6de00f' + amountData + dappAddressData + extraData

  let requestParams = {
    from: walletAddress,
    to: contractAddress,
    value: '0x00',
    // gasPrice: '0x09184e72a000',
    // gas: '0x186A0',//gas limit
    data: txnData
  }

  requestParams.gasPrice = await getSideChainGasPriceInWei();
  requestParams.gas = '0xAA690';

  console.log(requestParams)

  return ethereum.request({
    method: 'eth_sendTransaction',
    params: [requestParams],
  })

}

export async function getSideChainGasPriceInWei() {
  let gasPrice = "0x8F0D1800";//2400000000
  try {
    let url = getSideChainNetworks().url;
    let web3 = new Web3(url);
    gasPrice = await web3.eth.getGasPrice();
    gasPrice = web3.utils.numberToHex(gasPrice.toString());
  } catch (err) {
    gasPrice = "0x8F0D1800";
    console.error('#### getGasPrice err: ', err);
  }
  return gasPrice;
}

export async function getDAppProductList(contractAddress) {
  console.log("getDAppProductList start ");
  let contractJson = DAppStockSample
  let url = getSideChainNetworks().url
  let web3 = new Web3(url)
  var tokenContract = new web3.eth.Contract(contractJson.abi, contractAddress)
  let productList = await tokenContract.methods.getProductList().call()

  let contractJson2 = EurusERC20
  var tokenContract2 = new web3.eth.Contract(contractJson2.abi, getUsdtAddress())
  let decimals = await tokenContract2.methods.decimals().call()

  let productArr = []
  for (let index = 0; index < productList.length; index++) {
    const element = productList[index];
    // console.log(`${index} product: ${element}`);
    // console.log(`${index} product: ${element.name}`);
    let stock = await tokenContract.methods.productStockList(element.productId).call()
    productArr.push({
      productId: element.productId,
      name: element.name,
      price: amountDivideDecimals(element.price, decimals),
      stock: stock,
      onShelf: (element.onShelf?"Yes":"No"),
    })
  }

  console.log("productArr: ");
  console.log(productArr);

  return productArr
}

export async function getDAppProductStockList(contractAddress, productId) {
  console.log("getDAppProductStockList start ");
  console.log("productId:", productId);
  let contractJson = DAppStockSample
  let url = getSideChainNetworks().url
  let web3 = new Web3(url)
  var tokenContract = new web3.eth.Contract(contractJson.abi, contractAddress)
  let productStockList = await tokenContract.methods.productStockList(productId).call()

  console.log("productStockList: ");
  console.log(productStockList);

}

export async function purchaseToDapp(contractAddress, dappAddress, productId, quantity, price) {
  console.log("purchaseToDapp start")

  const url = getSideChainNetworks().url
  const chainId = getSideChainNetworks().chainId
  const net_version = await getChainId()
  if (net_version != chainId) {
    throw 9991
  }

  const accounts = await getAccount();
  const walletAddress = accounts[0]

  const contractJson = EurusERC20
  
  let web3 = new Web3(url)
  var tokenContract = new web3.eth.Contract(contractJson.abi, contractAddress)
  let decimals = await tokenContract.methods.decimals().call()
  price = amountMultipleDecimals(price, decimals);

  let dappAddressData = web3.utils.padLeft(dappAddress, 64)
  let productIdData = web3.utils.padLeft(web3.utils.numberToHex(productId.toString()), 64)
  
  let quantityData = web3.utils.padLeft(web3.utils.numberToHex(quantity.toString()), 64)
  let amountData = web3.utils.padLeft(web3.utils.numberToHex(price.toString()), 64)
  let extraData = web3.eth.abi.encodeParameter('bytes32', '0x0');
  dappAddressData = dappAddressData.substr(2)
  productIdData = productIdData.substr(2)
  quantityData = quantityData.substr(2)
  amountData = amountData.substr(2)
  extraData = extraData.substr(2)
  
  
  let txnData = '0x855bfa4e' + productIdData + quantityData + amountData + dappAddressData + extraData

  let requestParams = {
    from: walletAddress,
    to: contractAddress,
    value: '0x00',
    // gasPrice: '0x09184e72a000',
    // gas: '0x186A0',//gas limit
    data: txnData
  }

  requestParams.gasPrice = await getSideChainGasPriceInWei();
  requestParams.gas = '0xAA690';

  console.log(requestParams)

  return ethereum.request({
    method: 'eth_sendTransaction',
    params: [requestParams],
  })

}
