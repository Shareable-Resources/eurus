let deployHelper = require('./deploy_helper');

const marketingWallet = artifacts.require("MarketingWallet");
const ownedProxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");
const ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");



module.exports = async (deployer) => {
    deployHelper.readDeployLog();

    await deployHelper.deploy(deployer, marketingWallet, "MarketingRegWallet");


    let marketingWalletInstance = await marketingWallet.deployed()
 
 
    let accounts = await web3.eth.getAccounts();

    let ownedProxyObj = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<MarketingRegWallet>");
    let ownedProxyInstance = new web3.eth.Contract(ownedProxyJson.abi,ownedProxyObj.address);
    await ownedProxyInstance.methods.upgradeTo(marketingWalletInstance.address).send(await deployHelper.callParams(accounts, deployer));

    deployHelper.writeDeployLog();
}
