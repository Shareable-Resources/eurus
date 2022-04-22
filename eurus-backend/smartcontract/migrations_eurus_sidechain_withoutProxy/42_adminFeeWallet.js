let deployHelper = require('./deploy_helper');
const AdminFeeWallet = artifacts.require("AdminFeeWallet");
var internalSmartContractConfig = artifacts.require("InternalSmartContractConfig")
const ownedProxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");
const ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();


    await deployHelper.deploy(deployer,AdminFeeWallet);

    let adminFeeInstance = await AdminFeeWallet.deployed()
    //
    //
    let accounts = await web3.eth.getAccounts();
    //
    //
    let ownedProxyObj = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<AdminFeeWallet>");
    let ownedProxyInstance = new web3.eth.Contract(ownedProxyJson.abi,ownedProxyObj.address);
    await ownedProxyInstance.methods.upgradeTo(AdminFeeWallet.address).send(await deployHelper.callParams(accounts, deployer));
    //
    //
    //
    //
    // var internalSC = await internalSmartContractConfig.deployed()
    // internalSC.setAdminFeeWalletAddress(adminFeeInstance.address);
    //
    // let internalSCConfig = deployHelper.getSmartContractInfoByName(deployer, "InternalSmartContractConfig");
    // let proxy = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<PlatformWallet>")
    // let internalSCObj = new web3.eth.Contract(internalSCJson.abi, internalSCConfig.address);
    // let receipt = await internalSCObj.methods.setInnetWalletAddress(proxy.address).send(await deployHelper.callParams(accounts));

    deployHelper.writeDeployLog();

};