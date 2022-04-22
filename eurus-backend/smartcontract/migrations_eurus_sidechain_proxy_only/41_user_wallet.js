let deployHelper = require('./deploy_helper');
var userWalletObj = artifacts.require("UserWallet");
const internalSCJson = require("./../build/contracts/InternalSmartContractConfig.json");

web3.eth.handleRevert = true

module.exports = async function (deployer) {

    deployHelper.readDeployLog();

    let accounts = await web3.eth.getAccounts();

    await deployHelper.deploy(deployer, userWalletObj)

    const wallet = await userWalletObj.deployed();

    let internalSC = deployHelper.getSmartContractInfoByName(deployer, 'OwnedUpgradeabilityProxy<InternalSmartContractConfig>');

    var iscConfigContract = new web3.eth.Contract(internalSCJson.abi, internalSC.address);
    // console.log(await deployHelper.callParams(accounts))
    // console.log(deployer.provider.addresses[0])
    const reciept = await iscConfigContract.methods.setUserWalletAddress(wallet.address).send({from: accounts[0]}).catch(err=>console.log(err))
    // console.log(reciept)

    deployHelper.writeDeployLog();
}