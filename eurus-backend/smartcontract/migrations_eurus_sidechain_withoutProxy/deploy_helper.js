const fs = require('fs');
let smartContractMap = {};

module.exports.smartContractMap = smartContractMap;

module.exports.callParams = async (accounts, deployer)=>{
    return new Promise(async(res,rej)=>{
        return res({
            from: accounts[0],
            gas: 5700000,
            gasPrice: deployer.networks[deployer.network].gasPrice,
        });
    })

}

module.exports.deployWithProxy = async function(deployer, contract, contractNameOrNull, proxyContract){
    return new Promise(async(res,rej)=>{

        ! (deployer.network_id in this.smartContractMap) && (this.smartContractMap[deployer.network_id] = {'smartContract':{}});

        let deployedContract = await deployer.deploy(contract);

        let contractName = (contractNameOrNull == null)? contract.contractName:contractNameOrNull;

        let contractInstance = await contract.deployed();
        this.smartContractMap[deployer.network_id]['smartContract'][contractName] = { address: contractInstance.address};

        let deployedProxy = await deployer.deploy(proxyContract);
        let proxyInstance = await proxyContract.deployed();
        let proxyContractName = proxyContract.contractName + "<" + contractName + ">";
        this.smartContractMap[deployer.network_id]['smartContract'][proxyContractName] = {address: proxyInstance.address};

        console.log(contractName + ' ' + contractInstance.address);
        console.log(proxyContractName + ' ' + proxyInstance.address);
        await proxyInstance.upgradeTo(contract.address);


        return res(0);
    })


}

module.exports.deploy = async function(deployer, contract,  contractNameOrNull, params = null){
    return new Promise(async(res,rej)=>{
        ! (deployer.network_id in this.smartContractMap) && (this.smartContractMap[deployer.network_id] = {'smartContract':{}});

        if(params != null){
            await deployer.deploy(contract, ...params);
        }else{
            await deployer.deploy(contract)
        }
        let instance = await contract.deployed();
        let contractName = (contractNameOrNull == null)? contract.contractName:contractNameOrNull;

        console.log(contractName + ' ' + instance.address);
        this.smartContractMap[deployer.network_id]['smartContract'][contractName] = { address: instance.address};
        return res(0);
    })

}

module.exports.getSmartContractInfoByName = function(deployer, contractName){
    return this.smartContractMap[deployer.network_id]['smartContract'][contractName];
}

let deployFileName;
module.exports.readDeployLog = function() {
    
    try{
        let suffix = config.network;
        let index = suffix.indexOf('_');
        if (index >= 0){
            suffix = suffix.substring(0, index);
        }
        deployFileName = 'SmartContractDeploy_' + suffix + '.json';

        var rawData = fs.readFileSync('SmartContractDeploy_' + suffix + '.json');
        if (rawData.length > 0) {
            this.smartContractMap = JSON.parse(rawData);
        }
    }catch (ex) {
        console.log("readDeployLog: " + ex);
    }
}

module.exports.writeDeployLog = function(){

    this.smartContractMap.lastUpdateTime = new Date().toISOString();
    let jsonStr = JSON.stringify(this.smartContractMap);

    if (deployFileName && deployFileName != '' ){
        fs.writeFileSync(deployFileName, jsonStr);
    }else {
        console.log('deployFileName is not defined, cannot write to file');
    }
}
