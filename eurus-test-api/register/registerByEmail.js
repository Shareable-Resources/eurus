const email = "eu17@18m.dev";
const password = "aaaaaaa1";

const utils = require("../utils");
const nonce = utils.uuidv4();
const ts = Date.now();
const deviceId = ""
const loginWallet = utils.generateWallet(email+password);
const walletAddress = loginWallet.walletAddress.toLowerCase().substring(2);
const publicKey = loginWallet.publicKey.substring(2);
const privateKey = loginWallet.privateKey;

const message = 'deviceId=' + deviceId + '&timestamp=' + ts + '&walletAddress=' + walletAddress;
const messageHash = utils.getKeccak256Hash(message);
const signature = utils.ethSign(privateKey, messageHash);
const signature2 = signature.substring(2, 130);

// console.log("message: ", message);
// console.log("messageHash: ", messageHash);
// console.log("signature: ", signature);
// console.log("signature2: ", signature2);

var request = {
    "nonce": nonce,
    "timestamp": ts,
    "deviceId": deviceId,
    "email": email,
    "loginAddress": walletAddress,
    "signature": signature2,
    "publicKey": publicKey
}

var requestStr = JSON.stringify(request);
console.log("request: ");
console.log(requestStr);