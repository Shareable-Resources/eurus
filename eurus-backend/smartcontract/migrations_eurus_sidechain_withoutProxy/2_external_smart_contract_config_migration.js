let deployHelper = require('./deploy_helper');

var externalSmartContractConfig = artifacts.require("ExternalSmartContractConfig")

const ownedProxyJson = artifacts.require("OwnedUpgradeabilityProxy")

module.exports = async function (deployer) {

    deployHelper.readDeployLog();

    await deployHelper.deploy(deployer, externalSmartContractConfig);

    let extInstance = await externalSmartContractConfig.deployed()
    console.log("ext addr : ", extInstance.address)

    let ownedProxyObj = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<ExternalSmartContractConfig>");
    console.log("proxy addr : ", ownedProxyObj.address)

    let ownedProxyInstance = new web3.eth.Contract(ownedProxyJson.abi, ownedProxyObj.address);
    await ownedProxyInstance.methods.upgradeTo(extInstance.address).send({from: deployer.provider.addresses[0]});

    deployHelper.writeDeployLog();
};
