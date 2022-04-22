const Web3 = require("web3");
var XHR2 = require('xhr2-cookies').XMLHttpRequest; // jshint ignore: line
const EEAClient = require("web3-eea");

//var HttpHeaderProvider = require('httpheaderprovider');

Web3.providers.HttpProvider.prototype._prepareRequest = function(){
    var request;

    // the current runtime is a browser
    if (typeof XMLHttpRequest !== 'undefined') {
        request = new XMLHttpRequest();
    } else {
        request = new XHR2();
        var agents = {httpsAgent: this.httpsAgent, httpAgent: this.httpAgent, baseUrl: this.baseUrl};

        if (this.agent) {
            agents.httpsAgent = this.agent.https;
            agents.httpAgent = this.agent.http;
            agents.baseUrl = this.agent.baseUrl;
        }

        request.nodejsSet(agents);
    }

    request.open('POST', this.host, true);
    request.setRequestHeader('Content-Type','application/json');
    request.setRequestHeader('Authorization','Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJwZXJtaXNzaW9ucyI6WyIqOioiXSwiZXhwIjoxNjM3MzkwOTM2MDAwfQ.iIq1vQ_mFPM3AU-qAo_5ILZjYdVnnrWewPrQ88bvv6e-7KQhku8DG_cb5sF6oswq1MTqrh2Ie6nKkOXOj3M0_gLgRiiwoebx3a2aqGuvQkxF2tNHb7l41RIWSvqNv4nrXLikixIGr30C54MjEftRf9TFdK12hAE4FCXw2gSk1jhcbTXQec1VTEQcGJdBWFnMxfhavpSlxh6IvYQJya8XdzNCiMOk8NmbG9wmDMSu2IGsZHjDp3wNO0Ay_kaiXurCqpcWHywa_dyxcQ21feILG9Yq7UNUDa0xqMAeYjERRhzz_dEpPdN24R5TyKP5cbupuMcFsXIOzwYfLS-Cgoy3Ag');

    request.timeout = this.timeout;
    request.withCredentials = this.withCredentials;

    if(this.headers) {
        this.headers.forEach(function(header) {
            request.setRequestHeader(header.name, header.value);
        });
    }

    return request;
};

const erc20Contract = require("../../../besuuuuuuuuuuuuuuuuuu/build/contracts/TetherToken.json");
//const httpProvider = new HttpHeaderProvider('http://18.141.43.75:20000',headers)
//http://18.141.43.75:20000
//http://3.1.184.208:8545
const web3Client = new Web3('http://3.1.184.208:8545');
const eeaClient = new EEAClient(web3Client, 2018);



const createContract =  (contractArgs) => {
    return new Promise(async(res,rej)=>{
        try{
            let contract = new web3Client.eth.Contract(erc20Contract.abi);
            const result = await contract.deploy({data: erc20Contract.bytecode, arguments: contractArgs});
            
            return res(result);
        }catch(err){
            return rej(err);
        }
    });
};

const deployContract = (contract)=>{
    return new Promise(async(res,rej)=>{
        try{
            let gas = await contract.estimateGas({from: "0xfe3b557e8fb62b89f4916b721be55ceb828dbd73"});
            web3Client.eth.defaultAccount="0xfe3b557e8fb62b89f4916b721be55ceb828dbd73"
            const nonce=web3Client.eth.getTransactionCount("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73", 'pending')

            const rawTx={
                data: contract.encodeABI(),
                gas: gas,
                nonce: nonce
            }
            const signatureObj=await web3Client.eth.accounts.signTransaction(rawTx,"8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63" );
            const result = await web3Client.eth.sendSignedTransaction(signatureObj.rawTransaction);
            return res(result);
        }catch(err){
            return rej(err);
        }

    });
}

const deployPrivateContract = (contractBin)=>{
    return new Promise(async(res,rej)=>{
        const contractOptions = {
            data: contractBin,
            privateFrom: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
            //privacyGroupId: 2,
            privateFor: ["Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs=","k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="],
            privateKey: "8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"
        };
        try{
            const receipt = await eeaClient.eea.sendRawTransaction(contractOptions);
            return res(receipt);
        }catch(err){
            return rej(err);
        }
        
    })
}

const getPrivateContractAddress = transactionHash => {
  console.log("Transaction Hash ", transactionHash);
  return eeaClient.priv
    .getTransactionReceipt(transactionHash, "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=")
    .then(privateTransactionReceipt => {
      console.log("Private Transaction Receipt\n", privateTransactionReceipt);
      return privateTransactionReceipt.contractAddress;
    });
};

const main=async()=>{
    try{
        const contractObj = await createContract(["10000000000000000", "J Ethereum", "JHY","6"]);
        const receipt = await deployContract(contractObj);
        console.log(receipt);
    }catch(err){
        console.log(err);
    }
}

const privateMain=async()=>{
    try{
        const contractObj = await createContract(["10000000000000000", "J Ethereum", "JHY","6"]);
        const receipt = await deployPrivateContract(contractObj);
        await getPrivateContractAddress(receipt);
    }catch(err){
        console.log(err);
    }
}

main();
//privateMain();