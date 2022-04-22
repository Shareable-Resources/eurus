const email = "eu16@18m.dev";
const raEncryptedString = "rLoCAGMY08nHf9RApMIbdObMurCfkLtWDUuns3Sy5U2Kc7YT/X4qGv97egIj+waz0sPosl6PxXIFVJtIzIBjYI8g49jbk2OkjbN4in6K0qarEgeiH9sHfqDDdsNyo54Zl80zLo7SS87BGc570CBxH0KTvejLEqKNOMyBeQFROpEvW5Kg4I1kXsZ1s7YOdtnxC1LOVsC0B9+Z7K4Zug1ktl5CjKoLgLlzIDT8W3+aWAv0w7sApTBdm0YDveSC88ZHtDTq+LY4grGqmMiZx3ApN85baRqfuUUxf3B4UWLfp+sr52aW"

const utils = require("./utils");

for (let index = 1; index <= 9; index++) {
    // const oldPassword = "bbbbbbb" + index;
    const oldPassword = "aaaaaaa" + index;
    const serect = email + oldPassword
    // const serect = oldPassword
    // const serect = email
    // const serect = ""

    const mnemonicString = utils.decryptMnemonic(raEncryptedString);
    const oldWallet = utils.generateWallet2(serect, mnemonicString);
    console.log(serect, oldWallet.walletAddress);
}

