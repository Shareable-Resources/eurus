let deployHelper = require('./deploy_helper');

const wrappedEUN = artifacts.require("WrappedEUN");
const ownedProxy = artifacts.require("ReceiveFallbackOwnedUpgradeabilityProxy");

const externalSCConfigJson = require("./../build/contracts/ExternalSmartContractConfig.json");

module.exports = async (deployer) => {
    deployHelper.readDeployLog();

    try {
        // Deploy contracts
        await deployHelper.deployWithProxyWithFixedProxyContractName(deployer, wrappedEUN, ownedProxy, "OwnedUpgradeabilityProxy<WEUN>");
        let wrappedEUNInstance = await ownedProxy.deployed();

        // Get the deployed contracts' addresses
        let wrappedEUNAddr = wrappedEUNInstance.address;
        let externalSCConfigAddr = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<ExternalSmartContractConfig>");

        let accounts = await web3.eth.getAccounts();

        // Register new currency in external SC config
        let externalSCConfigObj = new web3.eth.Contract(externalSCConfigJson.abi, externalSCConfigAddr.address);

        await externalSCConfigObj.methods.removeCurrencyInfo("WEUN").send(await deployHelper.callParams(accounts, deployer));
        await externalSCConfigObj.methods.addCurrencyInfo(wrappedEUNAddr, "WEUN", 18, "").send(await deployHelper.callParams(accounts, deployer));
    } catch (err) {
        throw err;
    }

    deployHelper.writeDeployLog();
}
