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

module.exports.deployWithProxy = async function(deployer, contract, contractNameOrNull, proxyContract, params = null){
    return new Promise(async(res,rej)=>{

        ! (deployer.network_id in this.smartContractMap) && (this.smartContractMap[deployer.network_id] = {'smartContract':{}});

        let contractName = (contractNameOrNull == null)? contract.contractName:contractNameOrNull;
        
        let contractData = this.smartContractMap[deployer.network_id]['smartContract'][contractName]
        
        this.smartContractMap[deployer.network_id]['smartContract'][contractName] = { address: contractData.address};

        let receipt1 = await deployer.deploy(proxyContract);
        //console.log(receipt1)

        let proxyInstance = await proxyContract.deployed();
        let proxyContractName = proxyContract.contractName + "<" + contractName + ">";
        this.smartContractMap[deployer.network_id]['smartContract'][proxyContractName] = {address: proxyInstance.address};

        console.log(contractName + ' ' + contractData.address);
        console.log(proxyContractName + ' ' + proxyInstance.address);
        
        let receipt2 = await proxyInstance.upgradeTo(contractData.address);
        console.log(receipt2)

        return res(0);
    })
}

module.exports.deployWithProxyWithFixedProxyContractName = async function(deployer, contract, proxyContract, proxyContractName) {
    return new Promise(async (res, rej) => {
        !(deployer.network_id in this.smartContractMap) && (this.smartContractMap[deployer.network_id] = {'smartContract':{}});


        let contractData = this.smartContractMap[deployer.network_id]['smartContract'][contract.contractName];

        let receipt1 = await deployer.deploy(proxyContract);
        console.log(receipt1)
        let proxyInstance = await proxyContract.deployed();
        this.smartContractMap[deployer.network_id]['smartContract'][proxyContractName] = {address: proxyInstance.address};

        console.log(contract.contractName + ' ' + contractData.address);
        console.log(proxyContractName + ' ' + proxyInstance.address);
        let receipt2 = await proxyInstance.upgradeTo(contract.address);
        console.log(receipt2);

        return res(0);
    });
}

module.exports.deploy = async function(deployer, contract){
    return new Promise(async(res,rej)=>{
        ! (deployer.network_id in this.smartContractMap) && (this.smartContractMap[deployer.network_id] = {'smartContract':{}});


        await deployer.deploy(contract);
        let instance = await contract.deployed();
        console.log(contract.contractName + ' ' + instance.address);
        this.smartContractMap[deployer.network_id]['smartContract'][contract.contractName] = { address: instance.address};
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
