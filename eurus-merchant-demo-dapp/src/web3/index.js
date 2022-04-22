/* eslint-disable */
import Web3 from 'web3';
import EthCrypto from 'eth-crypto';
import axios from 'axios'
import {
    getSideChainNetworks,
    getEurusApiUrl,
    getAdminApiUrl,
    getUsdtSmartContractAddress,
    getDappSmartContractAddress,
} from '../network';
import { uuidv4 } from "../utils";
import { setUserId, setBalance, getUserId } from "../store";

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
    console.log('#### addEthereumChain start');
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
            await addEthereumChain();
        }
        // handle other "switch" errors
    }
    let chainId2 = await getChainId();
    if (chainId2 == chainId) {
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

export async function importWalletForEurus(inputWalletAddress, username) {
    let url = getEurusApiUrl() + '/user/importWallet';
    const nonce = uuidv4();
    const ts = Date.now();
    const deviceId = "";
    const walletAddress = inputWalletAddress.toLowerCase().substring(2);
    const message = 'deviceId=' + deviceId + '&timestamp=' + ts + '&walletAddress=' + walletAddress;
    let signature = await personalSign(inputWalletAddress, message);
    const importWalletResponse = await importWalletForAdmin(inputWalletAddress, username, message, signature)
    if (
        importWalletResponse &&
        importWalletResponse.data &&
        importWalletResponse.data.success &&
        importWalletResponse.data.data &&
        importWalletResponse.data.data.id &&
        importWalletResponse.data.data.balance
    ) {
        setBalance(amountDivideDecimals(importWalletResponse.data.data.balance, 6));
        setUserId(importWalletResponse.data.data.id);

        const rpcUrl = getSideChainNetworks().url;
        const web3Utils = (new Web3(rpcUrl)).utils;
        const initialBuffer = Buffer.from(web3Utils.hexToBytes(web3Utils.utf8ToHex(message), 'hex'));
        const preamble = Buffer.from(`\x19Ethereum Signed Message:\n${initialBuffer.length}`);
        const messageBuffer = Buffer.concat([preamble, initialBuffer]);
        const messageHash = web3Utils.sha3(messageBuffer);

        let publicKey = '';
        if (signature) {
            const recoverPublicKey = EthCrypto.recoverPublicKey(
                signature, // signature
                messageHash // message hash
            );
            publicKey = EthCrypto.publicKey.compress(recoverPublicKey)
            signature = signature.substring(2, 130);
        }
        console.log("publicKey:", publicKey);

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
    } else {
        console.log("importWalletResponse:", importWalletResponse.data)
        if (importWalletResponse.data.msg) {
            return { errorMsg: importWalletResponse.data.msg }
        }
        return {}
    }
}

export async function importWalletForAdmin(inputWalletAddress, username, message, signature) {
    let url = getAdminApiUrl() + '/merchant-clients/importWallet';
    let body = {
        "username": username,
        "walletAddress": inputWalletAddress.toLowerCase(),
        "message": message,
        "signature": signature
    };
    const headers = {
        "Content-Type": "application/json",
    };
    return axios.post(url, body, { headers: headers });
}

export function amountDivideDecimals(amount, decimals) {
    let output = 0;
    try {
        output = parseFloat(amount) / Math.pow(10, decimals);
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

export async function depositToDapp(amount) {
    console.log("depositToDapp start")

    const url = getSideChainNetworks().url
    const chainId = getSideChainNetworks().chainId
    const net_version = await getChainId()
    if (net_version != chainId) {
        throw 9991
    }

    const contractAddress = getUsdtSmartContractAddress()
    const dappAddress = getDappSmartContractAddress()

    const accounts = await getAccount();
    const walletAddress = accounts[0]

    let web3 = new Web3(url)
    let depositBalance = amountMultipleDecimals(amount, 6)
    let extraData = web3.utils.padLeft(web3.utils.numberToHex(getUserId().toString()), 64)
    let txnData = web3.eth.abi.encodeFunctionCall({
        name: "depositToDApp",
        type: 'function',
        inputs:
            [
                {
                    type: 'uint256',
                    name: 'amount'
                },
                {
                    type: 'address',
                    name: 'dappAddress'
                },
                {
                    type: 'bytes32',
                    name: 'extraData'
                }
            ]
    }, [depositBalance.toString(), dappAddress, extraData]);

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

export async function getTransaction(txnHash) {
    console.log("getTransaction start")
    let outTxn = {};
    try {
        let url = getSideChainNetworks().url;
        let web3 = new Web3(url);
        let txn = await web3.eth.getTransaction(txnHash);
        let txnReceipt = await web3.eth.getTransactionReceipt(txnHash);

        console.log("getTransaction txn:", txn)
        console.log("getTransaction txnReceipt:", txnReceipt)

        if (txnReceipt == null || txnReceipt == undefined) {
            return {};
        }
        if (txnReceipt && txnReceipt.transactionHash) {
            outTxn = {
                gas_fee: amountDivideDecimals(txn.gasPrice * txnReceipt.gasUsed, 18),
                txn_id: txnReceipt.transactionHash,
                status: txnReceipt.status
            };
            if (txnReceipt.revertReason) {
                outTxn.revertReason = web3.utils.hexToAscii(txnReceipt.revertReason)
            }
        }
    } catch (err) {
        console.error('#### getTransaction err: ', err);
    }
    return outTxn;
}