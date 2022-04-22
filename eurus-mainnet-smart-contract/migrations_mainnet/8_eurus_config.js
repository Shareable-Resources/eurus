let deployHelper = require('./deploy_helper');

const eurusInternalConfig = artifacts.require("EurusInternalConfig");
const ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");

const addressJson = require("./../SmartContractDeploy.json");

const eurusPlatformWalletJson = require("./../build/contracts/EurusPlatformWallet.json");
const eurusInternalConfigJson = require("./../build/contracts/EurusInternalConfig.json");

module.exports = async (deployer)=>{
    deployHelper.readDeployLog();
    await deployHelper.deployWithProxy(deployer,eurusInternalConfig, null, ownedProxy);

    let accounts = await web3.eth.getAccounts()
    let eurusPlatformWalletAddr=deployHelper.getSmartContractInfoByName(deployer,"OwnedUpgradeabilityProxy<EurusPlatformWallet>");

    try{
        let eurusInternalConfigObj = new web3.eth.Contract(eurusInternalConfigJson.abi, eurusInternalConfig.address);

        console.log("Calling setPlatformAddress. SC address: " + eurusPlatformWalletAddr);
        //await eurusInternalConfigObj.methods.setPlatformWalletAddress(eurusPlatformWalletAddr).send(await deployHelper.rinkebyCallParams(accounts))
    }catch(err){
        throw err;
    }


    deployHelper.writeDeployLog();
}


