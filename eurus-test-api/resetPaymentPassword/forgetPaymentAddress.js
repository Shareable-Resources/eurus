
const utils = require("../utils");
const nonce = utils.uuidv4();

var request = {
    "nonce": nonce,
}

var requestStr = JSON.stringify(request);
console.log("request: ");
console.log(requestStr);