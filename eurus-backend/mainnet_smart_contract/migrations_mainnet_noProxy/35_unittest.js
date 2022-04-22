let deployHelper = require('./deploy_helper');

let unitTestSC = artifacts.require('UnitTest');

module.exports = async function (deployer) {
    await deployHelper.deploy(deployer, unitTestSC);

}
