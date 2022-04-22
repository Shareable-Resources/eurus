let deployHelper = require('./deploy_helper');
const eurusPlatformWalletRC = artifacts.require("EurusPlatformWallet");
const ownedProxyJson = require("./../build/contracts/OwnedUpgradeabilityProxy.json");

module.exports = async function (deployer) {
  deployHelper.readDeployLog();
  
  await deployHelper.deploy(deployer, eurusPlatformWalletRC);
  let eurusPlatformWalletInstance = await eurusPlatformWalletRC.deployed();

  let proxyInfo = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<EurusPlatformWallet>")
  console.log('network id: ' + deployer.network_id)
  console.log('proxy address: ' + proxyInfo.address)
  let proxy = new web3.eth.Contract(ownedProxyJson.abi, proxyInfo.address)

  await proxy.methods.upgradeTo(eurusPlatformWalletInstance.address).send({from: deployer.provider.addresses[0]}).catch(err=>console.log(err))
  
  deployHelper.writeDeployLog();

};