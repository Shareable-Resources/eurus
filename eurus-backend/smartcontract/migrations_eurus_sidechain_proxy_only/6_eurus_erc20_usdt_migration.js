let deployHelper = require('./deploy_helper');

var usdt = artifacts.require("EurusERC20");
var ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");
var internalSCConfig = artifacts.require("InternalSmartContractConfig");

const erc20Json = require("./../build/contracts/EurusERC20.json");
const externalSCConfigJson = require("./../build/contracts/ExternalSmartContractConfig.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    let ownedProxyInstance;
   await deployHelper.deployWithProxy(deployer, usdt, 'USDT', ownedProxy)
    ownedProxyInstance = await ownedProxy.deployed();
    deployHelper.writeDeployLog();
};
