let deployHelper = require('./deploy_helper');

var usdt = artifacts.require("EurusERC20");
const ownedProxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");
var ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");
var internalSCConfig = artifacts.require("InternalSmartContractConfig");

const erc20Json = require("./../build/contracts/EurusERC20.json");
const externalSCConfigJson = require("./../build/contracts/ExternalSmartContractConfig.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    await deployHelper.deploy(deployer,usdt,"USDT")
    let usdtInstance = await usdt.deployed();
    let accounts = await web3.eth.getAccounts();


    let ownedProxyObj = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<USDT>");
    let ownedProxyInstance = new web3.eth.Contract(ownedProxyJson.abi,ownedProxyObj.address);
    await ownedProxyInstance.methods.upgradeTo(usdtInstance.address).send(await deployHelper.callParams(accounts, deployer));



    // let blacklistSC = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<USDT>");
    // let blacklistSCObj = new web3.eth.Contract(erc20Json.abi,blacklistSC.address)

    // let EurusUserDeposit = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<EurusUserDeposit>");
    //
    // await blacklistSCObj.methods.addBlackListDestAddress(EurusUserDeposit.address).send(await deployHelper.callParams(accounts));

    deployHelper.writeDeployLog();
};
