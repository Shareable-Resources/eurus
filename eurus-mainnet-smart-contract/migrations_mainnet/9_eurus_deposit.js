let deployHelper = require('./deploy_helper');

const eurusUserDeposit = artifacts.require("EurusUserDeposit");
const ownedProxy = artifacts.require("EtherForwardOwnedUpgradeabilityProxy");

const eurusInternalConfigJson = require("./../build/contracts/EurusInternalConfig.json");
const eurusUserDepositJson = require("./../build/contracts/EurusUserDeposit.json");
const etherForwardProxyJson = require("./../build/contracts/EtherForwardOwnedUpgradeabilityProxy.json");

module.exports = async (deployer) => {
    deployHelper.readDeployLog();

    try {
        // Deploy contracts, because here use a special proxy contract which can forward received ETH, proxy contract name need to be hardcoded
        await deployHelper.deployWithProxyWithFixedProxyContractName(deployer, eurusUserDeposit, ownedProxy, "OwnedUpgradeabilityProxy<EurusUserDeposit>");
        let eurusUserDepositInstance = await ownedProxy.deployed();
        
        // Get the deployed contracts' addresses, later need to update their stored values
        let eurusInternalConfigAddr = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<EurusInternalConfig>");
        let eurusUserDepositAddr = eurusUserDepositInstance.address;
        let eurusPlatformWalletAddr = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<EurusPlatformWallet>");

        let networkSetting = deployer.networks[deployer.network];

        let eurusInternalConfigObj = new web3.eth.Contract(eurusInternalConfigJson.abi, eurusInternalConfigAddr);
        console.log("eurusInternalConfig.setEurusUserDepositAddress")
        let receipt1 = await eurusInternalConfigObj.methods.setEurusUserDepositAddress(eurusUserDepositAddr).send(networkSetting);
        console.log(receipt1)

        let eurusUserDepositObj = new web3.eth.Contract(eurusUserDepositJson.abi, eurusUserDepositAddr);
        let receipt2 = await eurusUserDepositObj.methods.setEurusInternalConfigAddress(eurusInternalConfigAddr).send(networkSetting);
        console.log(receipt2)

        console.log("eurusUserDeposit.setEurusPlatformAddress")
        let receipt3 = await eurusUserDepositObj.methods.setEurusPlatformAddress(eurusPlatformWalletAddr).send(networkSetting);
        console.log(receipt3)

        // Address for forwarding ETH is in proxy contract, not logic contract
        let etherForwardProxyObj = new web3.eth.Contract(etherForwardProxyJson.abi, eurusUserDepositAddr);
        console.log("etherForwardProxy.setEtherForwardAddress")
        let receipt4 = await etherForwardProxyObj.methods.setEtherForwardAddress(eurusPlatformWalletAddr).send(networkSetting);
        console.log(receipt4)
        
    } catch (err) {
        throw err;
    }

    deployHelper.writeDeployLog();
}
