let deployHelper = require('./deploy_helper');

var testERC = artifacts.require("TestERC20");
var ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");

const erc20Json = require("./../build/contracts/TestERC20.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    let ownedProxyInstance;

    await deployHelper.deployWithProxy(deployer, testERC, 'PLA', ownedProxy)

    await testERC.deployed();
    ownedProxyInstance = await ownedProxy.deployed();

    deployHelper.writeDeployLog();
};
