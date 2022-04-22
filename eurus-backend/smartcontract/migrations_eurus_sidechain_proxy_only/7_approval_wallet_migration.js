let deployHelper = require('./deploy_helper');

let approvalWallet = artifacts.require("ApprovalWallet");
let approvalWalletReader = artifacts.require("ApprovalWalletReader");
let ownedProxyForApprovalWallet = artifacts.require("OwnedUpgradeabilityProxy");
let internalSC = artifacts.require("InternalSmartContractConfig");

const approvalWalletJson = require("./../build/contracts/ApprovalWallet.json");
const approvalWalletReaderJson = require("./../build/contracts/ApprovalWalletReader.json");

const internalSCJson = require("./../build/contracts/InternalSmartContractConfig.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();
    
    await deployHelper.deployWithProxy(deployer, approvalWallet, null, ownedProxyForApprovalWallet);
    
    deployHelper.writeDeployLog();

}


