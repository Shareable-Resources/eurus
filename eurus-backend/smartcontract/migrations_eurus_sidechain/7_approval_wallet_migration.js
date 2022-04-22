let deployHelper = require('./deploy_helper');

let approvalWallet = artifacts.require("ApprovalWallet");
let approvalWalletReader = artifacts.require("ApprovalWalletReader");
let ownedProxyForApprovalWallet = artifacts.require("OwnedUpgradeabilityProxy");
let internalSC = artifacts.require("InternalSmartContractConfig");

const approvalWalletJson = require("./../build/contracts/ApprovalWallet.json");
const approvalWalletReaderJson = require("./../build/contracts/ApprovalWalletReader.json");

const internalSCJson = require("./../build/contracts/InternalSmartContractConfig.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();
    
    let approvalWalletInstance;
    let ownedProxyInstance;
    let approvalWalletReaderInstance;
    let accounts = await web3.eth.getAccounts();


    await deployHelper.deployWithProxy(deployer, approvalWallet, null, ownedProxyForApprovalWallet);
    approvalWalletInstance = await approvalWallet.deployed();
    let usdtProxy = deployHelper.getSmartContractInfoByName(deployer,'OwnedUpgradeabilityProxy<USDT>');
    
    console.log("Construct init param. SC address: " + internalSC.address);
    let data = await approvalWalletInstance.init.request(internalSC.address,2);
    console.log("calling init");
    
    ownedProxyInstance = await ownedProxyForApprovalWallet.deployed();
    await ownedProxyInstance.sendTransaction({data: data.data});

    console.log("Construct addWriter param. USDT proxy address: " + usdtProxy.address);
    let data1 = await approvalWalletInstance.addWriter.request(usdtProxy.address);

    console.log("Calling addWriter");
    let result = await ownedProxyInstance.sendTransaction({data: data1.data});
    console.log("Calling addWriter result: " + JSON.stringify(result))
    
    
    await deployHelper.deploy(deployer,approvalWalletReader);
    approvalWalletReaderInstance = await approvalWalletReader.deployed();
   // await InitApprovalWallet(deployer, accounts,ownedProxyInstance.address,approvalWalletReaderInstance.address);

    deployHelper.writeDeployLog();

}



let InitApprovalWallet = (deployer, accounts,approvalWalletAddress,approvalWalletReaderAddress)=>{
    return new Promise(async(res,rej)=>{
        let approvalWalletObj = new web3.eth.Contract(approvalWalletJson.abi, approvalWalletAddress);

        let requirement = 1;
        console.log("setting requirement to " + requirement)
        let receipt = await approvalWalletObj.methods.changeRequirement(requirement).send(await deployHelper.callParams(accounts, deployer));
        console.log("init function tx receipt: " + JSON.stringify(receipt));
        try{
            let requirement = await approvalWalletObj.methods.required().call({
                from:accounts[0]
            });
            console.log("Approval Wallet Requirement: "+requirement);
        }catch(err){
            return rej(err);
        }


        receipt = await approvalWalletObj.methods.setWalletOwner(accounts[0]).send(await deployHelper.callParams(accounts, deployer));
        console.log("init function tx receipt: " + JSON.stringify(receipt));

        let walletOwner = await approvalWalletObj.methods.getWalletOwner().call({from:accounts[0]});
        console.log("MultiSig Wallet Owner: "+JSON.stringify(walletOwner));

        await approvalWalletObj.methods.setFallbackAddress(approvalWalletReaderAddress).send(await deployHelper.callParams(accounts, deployer));
        let fallbackAddr = await approvalWalletObj.methods.fallbackAddr().call({from:accounts[0]});
        console.log("Fallback Address: "+JSON.stringify(fallbackAddr));

        console.log("Setting approval wallet to internal smart contract config");
        let internalSCInfo = deployHelper.getSmartContractInfoByName(deployer,'InternalSmartContractConfig');
        let internalSCObj = new web3.eth.Contract(internalSCJson.abi, internalSCInfo.address);

        await internalSCObj.methods.setApprovalWalletAddress(approvalWalletAddress).send(
            await deployHelper.callParams(accounts, deployer)
        );

        return res(0);
    });
}
