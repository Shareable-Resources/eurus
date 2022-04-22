let deployHelper = require('./deploy_helper');
let platformWallet = artifacts.require("PlatformWallet");
const proxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");

module.exports = async function (deployer) {

    deployHelper.readDeployLog();
    
    await deployHelper.deploy(deployer, platformWallet);

    let platformWalletObj = await platformWallet.deployed();

    let accounts = await web3.eth.getAccounts();

    let proxy = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<PlatformWallet>");
 
    let proxyObj = new web3.eth.Contract(proxyJson.abi, proxy.address);
    await proxyObj.methods.upgradeTo(platformWalletObj.address).send(await deployHelper.callParams(accounts, deployer));

    deployHelper.writeDeployLog();
}
