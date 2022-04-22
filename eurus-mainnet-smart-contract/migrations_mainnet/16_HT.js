let deployHelper = require('./deploy_helper');

var usdt = artifacts.require("TestERC20");
var ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");

const erc20Json = require("./../build/contracts/TestERC20.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    let ownedProxyInstance;

    await deployHelper.deployWithProxy(deployer, usdt, 'HT', ownedProxy)

    await usdt.deployed();
    ownedProxyInstance = await ownedProxy.deployed();

    let accounts = await web3.eth.getAccounts();
    let usdt20obj = new web3.eth.Contract(erc20Json.abi, ownedProxyInstance.address);
    //await usdt20obj.methods.init("HT", "HT", 13242013521397079n, 18).send(await deployHelper.rinkebyCallParams(accounts));
    deployHelper.writeDeployLog();
};
