let deployHelper = require('./deploy_helper');
var eurusInternalConfig = artifacts.require("EurusInternalConfig")

module.exports = async function (deployer) {
    deployHelper.readDeployLog();
    
    await deployHelper.deploy(deployer, eurusInternalConfig);    

    deployHelper.writeDeployLog();
};

