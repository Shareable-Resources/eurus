let deployHelper = require('./deploy_helper');

var usdt = artifacts.require("EurusERC20");
var ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");
var internalSCConfig = artifacts.require("InternalSmartContractConfig");

const erc20Json = require("./../build/contracts/EurusERC20.json");
const externalSCConfigJson = require("./../build/contracts/ExternalSmartContractConfig.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();
    
    let usdtInstance;
    let ownedProxyInstance;

   await deployHelper.deployWithProxy(deployer, usdt, 'USDT', ownedProxy)

    usdtInstance = await usdt.deployed();
    ownedProxyInstance = await ownedProxy.deployed();

    // let accounts = await web3.eth.getAccounts();
    // let usdt20obj = new web3.eth.Contract(erc20Json.abi, ownedProxyInstance.address);
    // let err= await usdt20obj.methods.init(internalSCConfig.address, "USDT", "USDT", 13242013521397079n, 6).send(await deployHelper.callParams(accounts, deployer));
    // console.log("init function tx receipt: " + JSON.stringify(err));

    // console.log("Register USDT to ExternalSmartContractConfig" );
    // let extSCInfo = deployHelper.getSmartContractInfoByName(deployer, 'ExternalSmartContractConfig');
    // let externalScObj = new web3.eth.Contract(externalSCConfigJson.abi,extSCInfo.address)

    // await externalScObj.methods.removeCurrencyInfo("USDT").send(await deployHelper.callParams(accounts, deployer));
    // await externalScObj.methods.addCurrencyInfo(ownedProxyInstance.address, "USDT", 6, "").send(await deployHelper.callParams(accounts, deployer));

    deployHelper.writeDeployLog();
};
