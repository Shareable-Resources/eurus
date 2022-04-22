const email = "eu14@18m.dev";
const oldPassword = "aaaaaaa2";
const newPassword = "bbbbbbb1";
const raEncryptedString = "rLoCAJwaGU5Quf5ZxD6PVMVyCNrMGyMiBVAXAMQAY4zRJnXdLAfytmIfo6PjmUU1fCZFsI0Bqc+A9fDTuKzAXLGb8Aiw9vOxk3k9G3OXyKxf+TtbHAxbJv8pbY9n4EyxCnViwmGVGr5k2ISE4xXbjuF2UTcHhy0wtJaBPHYpizeec2lAMLeOczB5Jcgoq6P3gaF3wgsmgHxfuV1wXUJMS5QTKA+6Oo1gmUIeoDMXZgMBfAfQWVrDQXXH2p24grL2q+kkcHufHFLo1+Cmk216YHXZRWtgQEoFdtRZ7eYnciAGT5c4"

const utils = require("./utils");
const nonce = utils.uuidv4();
const ts = Date.now();
const deviceId = ""

const mnemonicString = utils.decryptMnemonic(raEncryptedString);
const oldWallet = utils.generateWallet(email + oldPassword, mnemonicString);
const newWallet = utils.generateWallet(email + newPassword, mnemonicString);

const oldMessage = 'deviceId=' + deviceId + '&timestamp=' + ts + '&walletAddress=' + oldWallet.walletAddress.toLowerCase().substring(2);
const oldMessageHash = utils.getKeccak256Hash(oldMessage);
const signatureOld = utils.ethSign(oldWallet.privateKey, oldMessageHash).substring(2, 130);

const newMessage = 'deviceId=' + deviceId + '&timestamp=' + ts + '&walletAddress=' + newWallet.walletAddress.toLowerCase().substring(2);
const newMessageHash = utils.getKeccak256Hash(newMessage);
const signatureNew = utils.ethSign(newWallet.privateKey, newMessageHash).substring(2, 130);


var request = {
    "nonce": nonce,
    "timestamp": ts,
    "deviceId": deviceId,
    "oldOwnerWalletAddress": oldWallet.walletAddress.toLowerCase().substring(2),
    "oldPublicKey": oldWallet.publicKey.substring(2),
    "oldSign": signatureOld,
    "ownerWalletAddress": newWallet.walletAddress.toLowerCase().substring(2),
    "publicKey": newWallet.publicKey.substring(2),
    "sign": signatureNew,
}
console.log(request);
var requestStr = JSON.stringify(request);
console.log("request: ");
console.log(requestStr);