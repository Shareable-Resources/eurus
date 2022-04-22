let deployHelper = require('./deploy_helper');
const eurusUserDeposit = artifacts.require("EurusUserDeposit");
const ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");
const deployedAddress = require("./../SmartContractDeploy.json");
const eurusInternalConfigJson = require("./../build/contracts/EurusInternalConfig.json");
const eurusUserDepositJson = require("./../build/contracts/EurusUserDeposit.json");
module.exports = async (deployer)=>{
    let accounts = await web3.eth.getAccounts();

    deployHelper.readDeployLog();

    await deployHelper.deployWithProxy(deployer,eurusUserDeposit, null, ownedProxy);
    let eurusUserDepositInstance = await ownedProxy.deployed();

    let eurusUserDepositObj = new web3.eth.Contract(eurusUserDepositJson.abi, eurusUserDepositInstance.address);
    let eurusInternalConfigAddr = deployHelper.getSmartContractInfoByName(deployer,"OwnedUpgradeabilityProxy<EurusInternalConfig>");
    try{
        //await eurusUserDepositObj.methods.setEurusInternalConfigAddress(eurusInternalConfigAddr).send(await deployHelper.rinkebyCallParams(accounts));
        let eurusInternalConfigObj = new web3.eth.Contract(eurusInternalConfigJson.abi,eurusInternalConfigAddr)
        //await eurusInternalConfigObj.methods.setEurusUserDepositAddress(eurusUserDepositInstance.address).send(await deployHelper.rinkebyCallParams(accounts));
    }catch(err){
        throw err;
    }

    deployHelper.writeDeployLog();
}