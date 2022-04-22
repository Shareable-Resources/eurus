
let deployHelper = require('./deploy_helper');
const Migrations = artifacts.require("EurusMigrations");


module.exports = function (deployer) {
  deployHelper.readDeployLog();
  
  deployer.deploy(Migrations);

  deployHelper.smartContractMap['lastUpdateTime'] =  new Date().toISOString();

  deployHelper.writeDeployLog();
};