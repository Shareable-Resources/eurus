let deployHelper = require('./deploy_helper');

let approvalWallet = artifacts.require("ApprovalWallet");
let approvalWalletReader = artifacts.require("ApprovalWalletReader");
let ownedProxyForApprovalWallet = artifacts.require("OwnedUpgradeabilityProxy");
let internalSC = artifacts.require("InternalSmartContractConfig");

const approvalWalletJson = require("./../build/contracts/ApprovalWallet.json");
const approvalWalletReaderJson = require("./../build/contracts/ApprovalWalletReader.json");

const proxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    await deployHelper.deploy(deployer,approvalWallet)

    approvalWalletInstance = await approvalWallet.deployed();
    let accounts = await web3.eth.getAccounts();
    let proxy = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<ApprovalWallet>")
    console.log('proxy address: ', proxy.address);
    let proxyInstance = new web3.eth.Contract(proxyJson.abi, proxy.address)
    await proxyInstance.methods.upgradeTo(approvalWalletInstance.address).send(await deployHelper.callParams(accounts, deployer));

    deployHelper.writeDeployLog();

}


