let deployHelper = require('./deploy_helper');

const marketingWallet = artifacts.require("MarketingWallet");
const ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");

const internalSCConfigJson = require("./../build/contracts/InternalSmartContractConfig.json");

module.exports = async (deployer) => {
    deployHelper.readDeployLog();

    try {
        // Deploy contracts
        await deployHelper.deployWithProxy(deployer, marketingWallet, "MarketingRegWallet", ownedProxy);
        let walletProxyInstance = await ownedProxy.deployed();

        // Get the deployed contracts' addresses
        let internalSCConfigAddr = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<InternalSmartContractConfig>");

        let accounts = await web3.eth.getAccounts();

        // Register new currency in external SC config
        let internalSCConfigObj = new web3.eth.Contract(internalSCConfigJson.abi, internalSCConfigAddr.address);

        await internalSCConfigObj.methods.setMarketingRegWalletAddress(walletProxyInstance.address).send(await deployHelper.callParams(accounts, deployer));
    } catch (err) {
        throw err;
    }

    deployHelper.writeDeployLog();
}
