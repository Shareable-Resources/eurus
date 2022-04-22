
const EthCrypto = require('eth-crypto');
const { ethers } = require("ethers");
const forge = require("node-forge");
const cryptoJS = require("crypto-js");
const elliptic = require('elliptic');
const secp256k1 = new elliptic.ec('secp256k1');
const Web3 = require('web3');

const ethSign = (privateKey, messageHash) => {
    return EthCrypto.sign(
        privateKey, // privateKey
        messageHash // hash of message
    );
}

const getKeccak256Hash = (message) => {
    return EthCrypto.hash.keccak256(message).substring(2);
}

const uuidv4 = () => {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

const generateWallet2 = (serect, mnemonic = '') => {
    let mnemonic2 = "carbon shuffle shoot knock alter bottom polar maple husband poet match spring";
    var path = "m/44'/60'/0'/0/0";
    if (mnemonic != '') {
        mnemonic2 = mnemonic;
        // path = "m/44'/60'/1'/0/0";
    }
    let mnemonicWallet2 = ethers.utils.HDNode.fromMnemonic(mnemonic2, serect);
    mnemonicWallet2 = mnemonicWallet2.derivePath(path);

    let wallet = {
        walletAddress: mnemonicWallet2.address,
        privateKey: mnemonicWallet2.privateKey,
        publicKey: mnemonicWallet2.publicKey
    }
    return wallet;
}

const generateWallet = (serect, mnemonic = '') => {
    let mnemonic2 = "carbon shuffle shoot knock alter bottom polar maple husband poet match spring";
    var path = "m/44'/60'/0'/0/0";
    if (mnemonic != '') {
        mnemonic2 = mnemonic;
        path = "m/44'/60'/1'/0/0";
    }
    let mnemonicWallet2 = ethers.utils.HDNode.fromMnemonic(mnemonic2, serect);
    mnemonicWallet2 = mnemonicWallet2.derivePath(path);

    let wallet = {
        walletAddress: mnemonicWallet2.address,
        privateKey: mnemonicWallet2.privateKey,
        publicKey: mnemonicWallet2.publicKey
    }
    return wallet;
}

const decryptMnemonic = (raEncryptedString) => {
    const privateKeyString = "MIICXQIBAAKBgQDU5YufylzObiWijgdmBfeZAyKrSOxq6Nrh+5Oh4crA/QQPkp8ZAcOoRvApgLLqQuUbfd7egWMOwczLID/yrA0Wi3k+Tk9Z7z4SfrPxaWu1+elHBktzLhQLkf4bj2PVAJH1zePJIy4PJCIO8gFelstEHFvWcso2ePA1nnZ6knunJQIDAQABAoGAYtUNNGjlHI/VuNjmZl5uywHBnnKEDj17H12C86u2TFEpCXGvmhRPmFcWNq4gYNAdO937EKBQNBGT2Nhn12g3yl4+4edPkf39URg3ZwCq4uxAWQ9z5rwPa1eOOAok2VNrJd1NX4WEUerXRgzT1MU449jaZkGz4m8LEiKYKEIW+CECQQDZ2hQy80HXM/gJVIaYN2khzW7Q5LigrVzGf1SJbXGb10XfSbm1r0rN6DSa2haSfChb/lj8eN1xjXMaxiZTsRfbAkEA+i1VRWtdHNMdTnXW/gU3cfbGxo7/nj7WUU4LaZ2AvZ12K/dRkwZYPEyxlId2HnJwnw8TRX8FIR8H3mkV5HXs/wJAcmOcL5Sjgch7+Qo1EkAmJ+WixnUSrOvaxy+cx/x7pwTGX5RquwesE6pV1Omm6Ivg9Uz8lLUyManAQtLA1Tkr+QJBAI5BzOUmgdHsMhP1agUTzk1dd/ZcRfoj3RZqfI7X4ubvbMzfW2FxECdprOi6hm4VwPiRR/ISokYNMRpFQw+gBt0CQQCzxWkD/jeGCI9F4qYI9qwgT8lxce702Y2D9huaQ3sxc9NiD/Bhm1Nw8P8W2tgXENh9BV9xcq/gEsHu7OaidSfF";
    const privateKeyPem = `-----BEGIN RSA PRIVATE KEY-----${privateKeyString}-----END RSA PRIVATE KEY-----`;
    const wholeEncrypted = new Buffer.from(raEncryptedString, 'base64')

    const rsaEncryptedAesKey = wholeEncrypted.slice(4, 132);
    const iv = wholeEncrypted.slice(132, 148).toString('base64');
    const aesEncryptedVal = wholeEncrypted.slice(148).toString('base64');

    const pki = forge.pki;
    const privateKey = pki.privateKeyFromPem(privateKeyPem);

    var decryptedAesKey = privateKey.decrypt(rsaEncryptedAesKey, 'RSA-OAEP');
    var decryptedAesKeyWords = cryptoJS.enc.Base64.parse(decryptedAesKey);
    var ivWords = cryptoJS.enc.Base64.parse(iv);

    var decrypted = cryptoJS.AES.decrypt(aesEncryptedVal, decryptedAesKeyWords, { iv: ivWords });
    var decryptedUtf8 = cryptoJS.enc.Utf8.stringify(decrypted);

    // console.log("### decryptMnemonic decryptedUtf8:", decryptedUtf8);
    return decryptedUtf8;
}

const getCentralizedTransferParams = (receiver, assetName, amount, privateKey, walletAddress) => {
    console.log('getCentralizedTransferParams start');
    let web3 = new Web3("http://13.228.169.25:8545");
    const encodedParametersHex = web3.eth.abi.encodeParameters(['address', 'string', 'uint256'], [receiver, assetName, amount]);
    const hashedMessage = Web3.utils.sha3(encodedParametersHex);
    const initialBuffer = Buffer.from(Web3.utils.hexToBytes(hashedMessage, 'hex'));
    const result = secp256k1EcdsaSignMessage('0x' + initialBuffer.toString('hex'), privateKey);
    let encodedParametersHex2 = web3.eth.abi.encodeParameters(['address', 'string', 'uint256', 'bytes'], [receiver, assetName, amount, result.signature]);
    return { params: '3c4cd911' + encodedParametersHex2.substring(2), walletAddress: walletAddress }
}

const secp256k1EcdsaSignMessage = (hash, privateKey) => {
    const signature = secp256k1
        .keyFromPrivate(Buffer.from(privateKey.slice(2), 'hex'))
        .sign(Buffer.from(hash.slice(2), 'hex'), { canonical: true });

    let v = Number(27 + Number(signature.recoveryParam)).toString(16);
    let r = Web3.utils.padLeft(signature.r.toString(16), 64);
    let s = Web3.utils.padLeft(signature.s.toString(16), 64);

    let output = {
        "v": '0x' + v,
        "r": '0x' + r,
        "s": '0x' + s,
        "signature": '0x' + r + s + v
    }

    return output;
}

exports.ethSign = ethSign;
exports.getKeccak256Hash = getKeccak256Hash;
exports.uuidv4 = uuidv4;
exports.generateWallet = generateWallet;
exports.generateWallet2 = generateWallet2;
exports.decryptMnemonic = decryptMnemonic;
exports.secp256k1EcdsaSignMessage = secp256k1EcdsaSignMessage;
exports.getCentralizedTransferParams = getCentralizedTransferParams;