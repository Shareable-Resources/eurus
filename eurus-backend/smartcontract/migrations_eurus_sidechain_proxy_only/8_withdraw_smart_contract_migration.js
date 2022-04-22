let deployHelper = require('./deploy_helper');

let withdrawSC = artifacts.require('WithdrawSmartContract');
let internalSCConfig = artifacts.require('InternalSmartContractConfig');
let ownedProxy = artifacts.require('OwnedUpgradeabilityProxy');
const withdrawSCJson = require("./../build/contracts/WithdrawSmartContract.json");
const usdtJson = require("./../build/contracts/EurusERC20.json");
module.exports = async function (deployer) {
    deployHelper.readDeployLog();
    
    let accounts = await web3.eth.getAccounts();
    
    await deployHelper.deployWithProxy(deployer, withdrawSC, 'WithdrawSmartContract', ownedProxy);

    deployHelper.writeDeployLog();
}
