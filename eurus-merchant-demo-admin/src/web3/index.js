/* eslint-disable */
import Web3 from 'web3';
import axios from 'axios'

import {
    getSideChainNetworks,
    getEurusApiUrl,
    getDappSmartContractAddress,
} from '../network';
import {
    getOperatorId
} from '../store';

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

export async function importWallet(inputWalletAddress, username) {
    let url = getEurusApiUrl() + '/merchant-clients/importWallet';
    const ts = Date.now();
    const walletAddress = inputWalletAddress.toLowerCase();
    const message = 'username=' + username + '&timestamp=' + ts + '&walletAddress=' + walletAddress;
    let signature = await personalSign(inputWalletAddress, message);
    let body = {
        "username": username,
        "walletAddress": walletAddress,
        "message": message,
        "signature": signature
    };

    const headers = {
        "Content-Type": "application/json",
    };

    console.log("url:", url)
    console.log("body:", body)

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

export async function refund(requestID, targetAddress, refundBalance) {
    console.log("refund start")

    const url = getSideChainNetworks().url
    const chainId = getSideChainNetworks().chainId
    const net_version = await getChainId()
    if (net_version != chainId) {
        throw 9991
    }

    const accounts = await getAccount();
    const walletAddress = accounts[0]

    let web3 = new Web3(url)
    const extraDataObjStr = `${requestID}|${getOperatorId()}`;
    const extraData = web3.utils.padLeft(web3.utils.utf8ToHex(extraDataObjStr), 64)

    let txnData = web3.eth.abi.encodeFunctionCall({
        name: "refund",
        type: 'function',
        inputs:
            [
                {
                    type: 'string',
                    name: 'targetAssetName'
                },
                {
                    type: 'uint256',
                    name: 'srcAmount'
                },
                {
                    type: 'address',
                    name: 'dest'
                },
                {
                    type: 'bytes32',
                    name: 'extraData'
                }
            ]
    }, ['USDT', refundBalance.toString(), targetAddress, extraData]);

    let requestParams = {
        from: walletAddress,
        to: getDappSmartContractAddress(),
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