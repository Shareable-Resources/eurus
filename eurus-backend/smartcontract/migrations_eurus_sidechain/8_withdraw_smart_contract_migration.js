let deployHelper = require('./deploy_helper');

let withdrawSC = artifacts.require('WithdrawSmartContract');
let internalSCConfig = artifacts.require('InternalSmartContractConfig');
let ownedProxy = artifacts.require('OwnedUpgradeabilityProxy');
const withdrawSCJson = require("./../build/contracts/WithdrawSmartContract.json");
const usdtJson = require("./../build/contracts/EurusERC20.json");
module.exports = async function (deployer) {
    deployHelper.readDeployLog();
    
    let accounts = await web3.eth.getAccounts();
    
    await deployHelper.deployWithProxy(deployer, withdrawSC, 'WithdrawSmartContract', ownedProxy);
    let ownedProxyInstance = await ownedProxy.deployed();

    console.log("Update InternalSmartContractConfig set WithdrawSmartContract address");
    let internalSCInstance = await internalSCConfig.deployed();
    await internalSCInstance.setWithdrawSmartContract(ownedProxyInstance.address);

    let withdrawContract = new web3.eth.Contract(withdrawSCJson.abi, ownedProxyInstance.address);

    console.log("Set ApprovalWallet address as the writer to WithdrawSmartContract");
    let approvalWalletSC = deployHelper.getSmartContractInfoByName(deployer, 'OwnedUpgradeabilityProxy<ApprovalWallet>');
    await withdrawContract.methods.addWriter(approvalWalletSC.address).send(
        await deployHelper.callParams(accounts, deployer)
    );

    let requirements = 5
    console.log("Init WithdrawSmartContract, set the number of requirement to " + requirements);
    withdrawContract.methods.init(internalSCInstance.address, requirements).send(
        await deployHelper.callParams(accounts, deployer)
    );

    let usdtInfo = deployHelper.getSmartContractInfoByName(deployer, 'OwnedUpgradeabilityProxy<USDT>');
    let usdtJsonSC = new web3.eth.Contract(usdtJson.abi, usdtInfo.address);
    usdtJsonSC.methods.addOwner(ownedProxyInstance.address).send(
        await deployHelper.callParams(accounts, deployer)
    );

    deployHelper.writeDeployLog();
}
