let deployHelper = require('./deploy_helper');

var dappSample = artifacts.require("DAppStockSample");
var ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");


module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    await deployHelper.deployWithProxy(deployer, dappSample, null, ownedProxy);

    deployHelper.writeDeployLog();
}
