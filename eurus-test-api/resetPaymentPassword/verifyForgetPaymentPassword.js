const email = "eu16@18m.dev";
const verificationCode = "407454";

const utils = require("../utils");
const nonce = utils.uuidv4();

var request = {
    "nonce": nonce,
    "email": email,
    "code": verificationCode,
}

var requestStr = JSON.stringify(request);
console.log("request: ");
console.log(requestStr);