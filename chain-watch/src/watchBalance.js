const Web3 = require('web3');
const ERC20 = require('./IERC20.json');

const fs = require('fs');
const path = require('path');
const SEP = ",";

//var cwd  = path.dirname(fs.realpathSync(__filename));

var cwd  = process.cwd();

console.info(`Current working directory ${cwd}`);

let checkListJsonStr = fs.readFileSync(`${cwd}/config/checkList.json`);
let checkList = JSON.parse(checkListJsonStr);

let chainJsonStr = fs.readFileSync(`${cwd}/config/chains.json`)
let chains = JSON.parse(chainJsonStr);

let web3Instances = {};
function getWeb3Instance(chainName) {
    let chain = chains[chainName];
    if (chain && chain.rpc) {
        if (!web3Instances[chainName]) {
            web3Instances[chainName] = new Web3(chain.rpc);
        }
        return web3Instances[chainName];
    }
    throw "Chain name not found " + chainName;
}

let first_write = true;
function writeResult(chain, wallet_name, address, token, balance){
    let filename = `${cwd}/watchBalance_result.csv`;
    if(first_write){
        console.info(`Write result file ${filename}`);
        let header = `chain${SEP}wallet_name${SEP}address${SEP}token${SEP}balance`;
        console.info(header);
        fs.writeFileSync(filename, header + "\n");
        first_write = false;
    } 
    let log = `${chain}${SEP}${wallet_name}${SEP}${address}${SEP}${token}${SEP}${balance}`
    console.info(log);
    fs.appendFileSync(filename, log + "\n");
    
}

async function checkBalance(chainName, checkItem) {
    let chain = chains[chainName];
    let tokenName = checkItem.token;
    let token = chain.tokens[tokenName];
    let web3 = getWeb3Instance(chainName);
    let log = "";
    if (token.address) {
        // is ERC20 token
        let contract = new web3.eth.Contract(ERC20.abi, token.address, {});
        let decimals = await contract.methods.decimals().call();
        let balance = await contract.methods.balanceOf(checkItem.address).call();
        writeResult(chainName, checkItem.name, checkItem.address, tokenName, balance / 10 ** decimals);
    } else {
        let balance = await web3.eth.getBalance(checkItem.address)
        writeResult(chainName, checkItem.name, checkItem.address, tokenName, balance / 10 ** 18);
    }
}

async function main() {
    for (let chainName in checkList) {
        let checkItems = checkList[chainName];
        for (var i = 0; i < checkItems.length; i++) {
            let checkItem = checkItems[i];
            await checkBalance(chainName, checkItem);
        }
    }
}

main();