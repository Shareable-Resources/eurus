let deployHelper = require('./deploy_helper');
var userWallet = artifacts.require("UserWallet");
const internalSCConfigJson = require("./../build/contracts/InternalSmartContractConfig.json");

web3.eth.handleRevert = true

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    let accounts = await web3.eth.getAccounts();

    await deployHelper.deploy(deployer, userWallet)

    const wallet = await userWallet.deployed();

    let internalSC = deployHelper.getSmartContractInfoByName(deployer, 'OwnedUpgradeabilityProxy<InternalSmartContractConfig>');

    var iscConfigContract = new web3.eth.Contract(internalSCConfigJson.abi, internalSC.address);
    console.log(await deployHelper.callParams(accounts, deployer))
    console.log(deployer.provider.addresses[0])
    const reciept = await iscConfigContract.methods.setUserWalletAddress(wallet.address).send({from: deployer.provider.addresses[0]}).catch(err=>console.log(err))
    console.log(reciept)

    deployHelper.writeDeployLog();
}