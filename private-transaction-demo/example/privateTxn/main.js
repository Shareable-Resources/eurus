const fs = require("fs");
const Web3 = require("web3");
const EEAClient=require("web3-eea");

//const tmpABI = require("../../build/contracts/MyHelloWorld.json");
const erc20Contract = require("../../../besuuuuuuuuuuuuuuuuuu/build/contracts/TetherToken.json");
const web3Client = new Web3("http://3.1.184.208:8545")
const eeaClient = new EEAClient(web3Client, 2018);
const contractAddress="0x4e15c5e3897ae1e23054da64d9ffc118d2fe7099"



const  transfer=  (to, amount) => {
    return new Promise(async(res,rej)=>{
        try{
            let contractObj = new web3Client.eth.Contract(erc20Contract.abi,contractAddress);
            let tmpABI =(await contractObj.methods.transfer(to,amount)).encodeABI();
            const contractOptions = {
                to: contractAddress,
                data: tmpABI,
                privateFrom: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
                privateFor: ["Ko2bVqD+nNlNYL5EE7y3IdOnviftjiizpjRt+HTuFBs=","k2zXEin4Ip/qBGlRkJejnGWdP9cjkK+DAvKNW31L2C8="],
                privateKey: "8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"
            };
            //return res(tmpABI);
            return res(eeaClient.eea.sendRawTransaction(contractOptions));

        }catch(err){
            return rej(err);
        }
    })
};

const getPrivateContractAddress = transactionHash => {
  console.log("Transaction Hash ", transactionHash);
  return eeaClient.priv
    .getTransactionReceipt(transactionHash, "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=")
    .then(privateTransactionReceipt => {
      console.log("Private Transaction Receipt\n", privateTransactionReceipt);
      return privateTransactionReceipt.contractAddress;
    });
};

const privateMain=async()=>{
    const result=await transfer("0x99B32dAD54F630D9ED36E193Bc582bbed273d666",1);
    console.log(result);
    getPrivateContractAddress(result);
}

privateMain();