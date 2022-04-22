let deployHelper = require('./deploy_helper');

var dappSample = artifacts.require("DAppSample");
const ownedProxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    await deployHelper.deploy(deployer, dappSample, "DAppSampleToken");


    let dappInstance = await dappSample.deployed()
    //
    //
    let accounts = await web3.eth.getAccounts();
    //
    //
    let ownedProxyObj = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<DAppSampleToken>");
    let ownedProxyInstance = new web3.eth.Contract(ownedProxyJson.abi,ownedProxyObj.address);
    await ownedProxyInstance.methods.upgradeTo(dappInstance.address).send(await deployHelper.callParams(accounts, deployer));

    deployHelper.writeDeployLog();
}
