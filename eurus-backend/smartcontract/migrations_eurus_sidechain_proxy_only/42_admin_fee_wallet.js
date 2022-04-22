let deployHelper = require('./deploy_helper');
var adminFeeWallet = artifacts.require("AdminFeeWallet");
let ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");
const internalSCJson = require("./../build/contracts/InternalSmartContractConfig.json");

web3.eth.handleRevert = true

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    let accounts = await web3.eth.getAccounts();

    await deployHelper.deployWithProxy(deployer, adminFeeWallet,null, ownedProxy)

    const ownedProxyInstance = await ownedProxy.deployed();


    let internalSC = deployHelper.getSmartContractInfoByName(deployer, 'OwnedUpgradeabilityProxy<InternalSmartContractConfig>');

    var iscConfigContract = new web3.eth.Contract(internalSCJson.abi, internalSC.address);
    const reciept = await iscConfigContract.methods.setAdminFeeWalletAddress(ownedProxyInstance.address).send({from: deployer.provider.addresses[0]}).catch(err=>console.log(err))
    // console.log(reciept)

    deployHelper.writeDeployLog();
}