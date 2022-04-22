let deployHelper = require('./deploy_helper');
var internalSmartContractConfig = artifacts.require("InternalSmartContractConfig")
var ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    await deployHelper.deployWithProxy(deployer, internalSmartContractConfig, "InternalSmartContractConfig", ownedProxy);

    deployHelper.writeDeployLog();
};
