const email = "eu14@18m.dev";
const paymentPassword = "bbbbbbb1";
const raEncryptedString = "rLoCAJwaGU5Quf5ZxD6PVMVyCNrMGyMiBVAXAMQAY4zRJnXdLAfytmIfo6PjmUU1fCZFsI0Bqc+A9fDTuKzAXLGb8Aiw9vOxk3k9G3OXyKxf+TtbHAxbJv8pbY9n4EyxCnViwmGVGr5k2ISE4xXbjuF2UTcHhy0wtJaBPHYpizeec2lAMLeOczB5Jcgoq6P3gaF3wgsmgHxfuV1wXUJMS5QTKA+6Oo1gmUIeoDMXZgMBfAfQWVrDQXXH2p24grL2q+kkcHufHFLo1+Cmk216YHXZRWtgQEoFdtRZ7eYnciAGT5c4"
const txnReceiver = '0x4DfB6d6790054F3EB68324BC230E3104137CA8Db';
const txnToken = 'USDT';
const txnAmount = '1000000';

const utils = require("./utils");
const nonce = utils.uuidv4();

const mnemonicString = utils.decryptMnemonic(raEncryptedString);
const paymentWallet = utils.generateWallet(email + paymentPassword, mnemonicString);

var centralizedTransferParams = utils.getCentralizedTransferParams(
    txnReceiver,
    txnToken,
    txnAmount,
    paymentWallet.privateKey,
    paymentWallet.walletAddress.toLowerCase()
);

var request = {
    nonce: nonce,
    value: '0',
    gasPrice: 2500000000,
    inputFunction: centralizedTransferParams.params,
};

// console.log(request);
var requestStr = JSON.stringify(request);
console.log("request: ");
console.log(requestStr);