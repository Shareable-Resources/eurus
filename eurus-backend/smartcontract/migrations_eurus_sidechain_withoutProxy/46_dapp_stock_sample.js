// var storeTest = artifacts.require("StoreTest");


// web3.eth.handleRevert = true

// module.exports = async function (deployer) {

//     deployer.deploy(storeTest)

// }

let deployHelper = require('./deploy_helper');

var dappStockSample = artifacts.require("DAppStockSample");
const ownedProxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    await deployHelper.deploy(deployer, dappStockSample, "DAppStockSample");


    let dappInstance = await dappStockSample.deployed()
    //
    //
    let accounts = await web3.eth.getAccounts();
    //
    //
    let ownedProxyObj = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<DAppStockSample>");
    let ownedProxyInstance = new web3.eth.Contract(ownedProxyJson.abi,ownedProxyObj.address);
    await ownedProxyInstance.methods.upgradeTo(dappInstance.address).send(await deployHelper.callParams(accounts, deployer));

    deployHelper.writeDeployLog();
}


