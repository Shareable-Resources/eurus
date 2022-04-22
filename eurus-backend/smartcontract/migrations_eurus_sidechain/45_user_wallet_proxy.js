let deployHelper = require('./deploy_helper');
var userWallet = artifacts.require("UserWallet");
var userWalletProxy = artifacts.require("UserWalletProxy");

web3.eth.handleRevert = true

module.exports = async function (deployer) {
    deployHelper.readDeployLog();    
    await deployHelper.deploy(deployer, userWalletProxy);
    deployHelper.writeDeployLog();
}

