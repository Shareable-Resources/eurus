let deployHelper = require('./deploy_helper');

var walletAddressSmartContract = artifacts.require("WalletAddressMap")
let ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");
const internalSCConfigJson = require("./../build/contracts/InternalSmartContractConfig.json");

module.exports = async function (deployer) {

    deployHelper.readDeployLog();
    
    // await deployHelper.deploy(deployer, walletAddressSmartContract);
    //
    // var walletAddressMap = await walletAddressSmartContract.deployed();
    // var internalSC = await internalSmartContractConfig.deployed()
    // internalSC.setWalletAddressMap(walletAddressMap.address);

    await deployHelper.deployWithProxy(deployer, walletAddressSmartContract,null, ownedProxy)
    var proxy = await ownedProxy.deployed();
    var internalSCInfo = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<InternalSmartContractConfig>")

    let internalSC = new web3.eth.Contract(internalSCConfigJson.abi, internalSCInfo.address)

    internalSC.methods.setWalletAddressMap(proxy.address).send({from: deployer.provider.addresses[0]}).catch(err=>console.log(err))

    deployHelper.writeDeployLog();

};
