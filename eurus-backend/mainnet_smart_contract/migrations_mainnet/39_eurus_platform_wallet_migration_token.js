let deployHelper = require('./deploy_helper');
const eurusPlatformWalletRC = artifacts.require("EurusPlatformWalletMigration");
const ownedProxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");

module.exports = async function (deployer) {
  deployHelper.readDeployLog();
  
  await deployHelper.deploy(deployer, eurusPlatformWalletRC);
  
  let eurusPlatformWalletInstance = deployHelper.getSmartContractInfoByName(deployer, 'EurusPlatformWalletMigration')

  let proxyInfo = deployHelper.getSmartContractInfoByName(deployer, "Old_OwnedUpgradeabilityProxy<EurusPlatformWallet>")

  let proxy = new web3.eth.Contract(ownedProxyJson.abi, proxyInfo.address)

  let networkSetting = deployer.networks[deployer.network];

  await proxy.methods.upgradeTo(eurusPlatformWalletInstance.address).send(networkSetting).catch(err=>console.log(err))
  
  deployHelper.writeDeployLog();

};