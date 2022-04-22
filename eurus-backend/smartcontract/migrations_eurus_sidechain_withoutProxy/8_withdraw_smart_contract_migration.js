let deployHelper = require('./deploy_helper');

let withdrawSC = artifacts.require('WithdrawSmartContract');
let internalSCConfig = artifacts.require('InternalSmartContractConfig');
const proxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");

module.exports = async function (deployer) {

    deployHelper.readDeployLog();
    
    let accounts = await web3.eth.getAccounts();
    await deployHelper.deploy(deployer,withdrawSC)
    let withdrawSCInstance = await withdrawSC.deployed();

    let ownedProxy = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<WithdrawSmartContract>");
    let ownedProxyObj = new web3.eth.Contract(proxyJson.abi, ownedProxy.address)

    await ownedProxyObj.methods.upgradeTo(withdrawSCInstance.address).send(await deployHelper.callParams(accounts, deployer));

    deployHelper.writeDeployLog();
}
