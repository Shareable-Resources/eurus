const userId = 341;
const email = "eu16@18m.dev";
const paymentPassword = "bbbbbbb1";
const raEncryptedString = "rLoCAGMY08nHf9RApMIbdObMurCfkLtWDUuns3Sy5U2Kc7YT/X4qGv97egIj+waz0sPosl6PxXIFVJtIzIBjYI8g49jbk2OkjbN4in6K0qarEgeiH9sHfqDDdsNyo54Zl80zLo7SS87BGc570CBxH0KTvejLEqKNOMyBeQFROpEvW5Kg4I1kXsZ1s7YOdtnxC1LOVsC0B9+Z7K4Zug1ktl5CjKoLgLlzIDT8W3+aWAv0w7sApTBdm0YDveSC88ZHtDTq+LY4grGqmMiZx3ApN85baRqfuUUxf3B4UWLfp+sr52aW"

const utils = require("../utils");
const nonce = utils.uuidv4();
const ts = Date.now();
const deviceId = ""
const mnemonicString = utils.decryptMnemonic(raEncryptedString);
const paymentWallet = utils.generateWallet(email+paymentPassword, mnemonicString);
const walletAddress = paymentWallet.walletAddress.toLowerCase().substring(2);
const publicKey = paymentWallet.publicKey.substring(2);
const privateKey = paymentWallet.privateKey;

const message = 'deviceId=' + deviceId + '&timestamp=' + ts + '&walletAddress=' + walletAddress;
const messageHash = utils.getKeccak256Hash(message);
const signature = utils.ethSign(privateKey, messageHash);
const signature2 = signature.substring(2, 130);

var request = {
    nonce: nonce,
    timestamp: ts,
    deviceId: deviceId,
    email: email,
    userId: userId,
    ownerWalletAddress: walletAddress,
    sign: signature2,
    publicKey: publicKey,
};

var requestStr = JSON.stringify(request);
console.log("request: ");
console.log(requestStr);