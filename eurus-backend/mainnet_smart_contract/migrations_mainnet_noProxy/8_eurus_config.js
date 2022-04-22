let deployHelper = require('./deploy_helper');

const eurusInternalConfig = artifacts.require("EurusInternalConfig");
const ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");


const eurusPlatformWalletJson = require("./../build/contracts/EurusPlatformWallet.json");
const proxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");

module.exports = async (deployer)=>{
    deployHelper.readDeployLog();
    
    await deployHelper.deploy(deployer,eurusInternalConfig);

    // let accounts = await web3.eth.getAccounts()
    // let proxyAddr = deployHelper.getSmartContractInfoByName(deployer,"OwnedUpgradeabilityProxy<EurusInternalConfig>");
    // let eurusInteranlConfigInstance = await eurusInternalConfig.deployed();

    // console.log('Proxy address: ' + proxyAddr.address);
    // let proxyJsonObj = new web3.eth.Contract(proxyJson.abi, proxyAddr.address);
    // await proxyJsonObj.methods.upgradeTo(eurusInteranlConfigInstance.address).send({from: accounts[0]});

    deployHelper.writeDeployLog();
}


