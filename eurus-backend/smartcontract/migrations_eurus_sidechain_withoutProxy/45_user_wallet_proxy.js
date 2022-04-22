let deployHelper = require('./deploy_helper');
var userWalletProxy = artifacts.require("UserWalletProxy");
const internalSCConfigJson = require("./../build/contracts/InternalSmartContractConfig.json");


web3.eth.handleRevert = true

module.exports = async function (deployer) {
    deployHelper.readDeployLog();    
    await deployHelper.deploy(deployer, userWalletProxy);
    let proxy = await userWalletProxy.deployed();
    let internalSC = deployHelper.getSmartContractInfoByName(deployer, 'OwnedUpgradeabilityProxy<InternalSmartContractConfig>');

    var iscConfigContract = new web3.eth.Contract(internalSCConfigJson.abi, internalSC.address);

    const receipt = await iscConfigContract.methods.setUserWalletProxyAddress(proxy.address).send({from: deployer.provider.addresses[0]}).catch(err=>console.log(err))


    console.log(receipt)

    deployHelper.writeDeployLog();
} 

