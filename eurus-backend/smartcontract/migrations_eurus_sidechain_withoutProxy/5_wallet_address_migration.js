let deployHelper = require('./deploy_helper');

var walletAddressSmartContract = artifacts.require("WalletAddressMap")
var internalSmartContractConfig = artifacts.require("InternalSmartContractConfig")
const ownedProxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();
    
    await deployHelper.deploy(deployer, walletAddressSmartContract);
    let accounts = await web3.eth.getAccounts();

    var walletAddressMap = await walletAddressSmartContract.deployed();
    var internalSC = await internalSmartContractConfig.deployed()
    
    walletAddressMapProxy = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<WalletAddressMap>")
    console.log("Wallet address proxy address: " + walletAddressMapProxy.address);


    await internalSC.setWalletAddressMap(walletAddressMapProxy.address);
    let ownedProxyInstance = new web3.eth.Contract(ownedProxyJson.abi, walletAddressMapProxy.address);
    await ownedProxyInstance.methods.upgradeTo(walletAddressMap.address).send(await deployHelper.callParams(accounts, deployer));

    deployHelper.writeDeployLog();
};
