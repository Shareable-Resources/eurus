let deployHelper = require('./deploy_helper');

var usdc = artifacts.require("EurusERC20");
var ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");
var internalSCConfig = artifacts.require("InternalSmartContractConfig");

const erc20Json = require("./../build/contracts/EurusERC20.json");
const externalSCConfigJson = require("./../build/contracts/ExternalSmartContractConfig.json");

module.exports = async function (deployer) {
    deployHelper.readDeployLog();

    let ownedProxyInstance;

    await deployHelper.deployWithProxy(deployer, usdc, 'USDC', ownedProxy)

    await usdc.deployed();
    ownedProxyInstance = await ownedProxy.deployed();

    // let accounts = await web3.eth.getAccounts();
    // let usdcobj = new web3.eth.Contract(erc20Json.abi, ownedProxyInstance.address);
    // let err= await usdcobj.methods.init(internalSCConfig.address, "USDC", "USDC", 13242013521397079n, 6).send(await deployHelper.callParams(accounts, deployer));
    // console.log("init function tx receipt: " + JSON.stringify(err));

    // console.log("Register USDC to ExternalSmartContractConfig" );
    // let extSCInfo = deployHelper.getSmartContractInfoByName(deployer, 'ExternalSmartContractConfig');
    // let externalScObj = new web3.eth.Contract(externalSCConfigJson.abi,extSCInfo.address)

    // await externalScObj.methods.removeCurrencyInfo("USDC").send(await deployHelper.callParams(accounts, deployer));
    // await externalScObj.methods.addCurrencyInfo(ownedProxyInstance.address, "USDC").send(await deployHelper.callParams(accounts, deployer));
    deployHelper.writeDeployLog();
};
