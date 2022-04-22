let deployHelper = require('./deploy_helper');
const eurusPlatformWalletRC = artifacts.require("EurusPlatformWallet");
var ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");

module.exports = async function (deployer) {
  deployHelper.readDeployLog();
  
  await deployHelper.deployWithProxy(deployer, eurusPlatformWalletRC, null, ownedProxy);

  deployHelper.writeDeployLog();

};