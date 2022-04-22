let deployHelper = require('./deploy_helper');
let platformWallet = artifacts.require("PlatformWallet");
let ownedProxy = artifacts.require("OwnedUpgradeabilityProxy");

const internalSCJson = require("./../build/contracts/InternalSmartContractConfig.json");
const platformWalletJson = require("./../build/contracts/PlatformWallet.json");
const USDTJson = require("./../build/contracts/EurusERC20.json");

module.exports = async function (deployer) {

    deployHelper.readDeployLog();


   await deployHelper.deployWithProxy(deployer, platformWallet, null, ownedProxy);
    // await deployHelper.deploy(deployer, platformWallet);

    console.log("Setting OwnedUpgradeabilityProxy<PlatformWallet> address to InternalSmartContractConfig");

    let accounts = await web3.eth.getAccounts();
    let internalSCConfig = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<InternalSmartContractConfig>");
    let proxy = deployHelper.getSmartContractInfoByName(deployer, "OwnedUpgradeabilityProxy<PlatformWallet>")
    let internalSCObj = new web3.eth.Contract(internalSCJson.abi, internalSCConfig.address);
    let receipt = await internalSCObj.methods.setInnetWalletAddress(proxy.address).send(await deployHelper.callParams(accounts, deployer));

    let requirement = 5;
    console.log("Setting PlatformWallet multi signature requirement to " + requirement);
    let platformWalletObj = new web3.eth.Contract(platformWalletJson.abi, proxy.address);
    let receipt1 = await platformWalletObj.methods.changeRequirement(requirement).send(await deployHelper.callParams(accounts, deployer));


    let receipt2 = await platformWalletObj.methods.setInternalSmartContractConfig(internalSCConfig.address).send(await deployHelper.callParams(accounts, deployer));
    deployHelper.writeDeployLog();


}
