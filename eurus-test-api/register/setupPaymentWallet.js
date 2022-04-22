const userId = 0;
const email = "eu17@18m.dev";
const paymentPassword = "bbbbbbb1";
const raEncryptedString = "rLoCAMfNa2Xuh3QfSG10q+JignFnRojdxo3Ze+ZyTo9bP758FB7uKtltLU0wxoqQfEQOb0ufg37X+eSSnBP7dv/7HKNHBybXKeEEQhAxCPH6NnTZhJm+t2tCZHzY4o1lYBXBFs/i9P/2e/lCFvRZ7luOq/bvpWrePaEban6r28LjMu47vOeOagx4JYlVGRExGw+WCahyBX83SE8fgThVu3pgbaoaNy+Nk0mIwSvjDmv9cjc5je39iiCaKl+1Gg/HSAXtmFCL8Ji0KIqIoXwmOiZSp1R6BbPFhaq7yNs2kOBDAmff";

const utils = require("../utils");
const mnemonicString = utils.decryptMnemonic(raEncryptedString);
const paymentWallet = utils.generateWallet(email+paymentPassword, mnemonicString);
const walletAddress = paymentWallet.walletAddress.toLowerCase().substring(2);

var request = {
    "userId": userId,
    "address": walletAddress
}

var requestStr = JSON.stringify(request);
console.log("request: ");
console.log(requestStr);