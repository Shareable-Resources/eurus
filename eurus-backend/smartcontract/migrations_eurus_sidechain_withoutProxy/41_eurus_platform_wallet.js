let deployHelper = require('./deploy_helper');
const eurusPlatformWalletRC = artifacts.require("EurusPlatformWallet");


module.exports = async function (deployer) {

    //NOT REQUIRED TO DEPLOY AT SIDECHAIN!
    
    // deployHelper.readDeployLog();

    // await deployHelper.deploy(deployer, eurusPlatformWalletRC);

    // deployHelper.writeDeployLog();

};