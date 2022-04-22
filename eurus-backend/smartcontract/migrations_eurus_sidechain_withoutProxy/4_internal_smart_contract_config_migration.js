let deployHelper = require('./deploy_helper');

var internalSmartContractConfig = artifacts.require("InternalSmartContractConfig")
var externalSmartContractConfig = artifacts.require("ExternalSmartContractConfig")

const ownedProxyJson = artifacts.require("OwnedUpgradeabilityProxy")

module.exports = async function (deployer) {
    deployHelper.readDeployLog();
    // let accounts = await web3.eth.getAccounts();

    await deployHelper.deploy(deployer, internalSmartContractConfig);

    let ownedProxyObj = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<InternalSmartContractConfig>");
    let ownedProxyInstance = new web3.eth.Contract(ownedProxyJson.abi,ownedProxyObj.address);
    await ownedProxyInstance.methods.upgradeTo(internalSmartContractConfig.address).send({from: deployer.provider.addresses[0]});

    deployHelper.writeDeployLog();
  };
