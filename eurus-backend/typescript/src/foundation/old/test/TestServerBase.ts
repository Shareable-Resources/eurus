/*
import { EthClient } from '../../src/ethereum/EthClient';
import { EurusERC20 } from './../../../smartcontract/build/typescript/EurusERC20.d';
import Web3 from "web3";
import { ServerConfigBase } from '../../src/server/ServerConfigBase';
import { ServerReturnCode } from '../ServerReturnCode';
import { ServerBase } from '../../src/server/ServerBase';



export async function TestLoadConfig(){

    let config :ServerConfigBase =  new ServerConfigBase();
    let server: ServerBase = new ServerBase(config);

    let isSuccess: boolean = server.loadConfig("AdminServerConfig.json");

    if (isSuccess){
        console.log("Login authServer");
        isSuccess = server.loginAuthServer((returnCode: ServerReturnCode) => {
            if (returnCode == ServerReturnCode.Success) {
                console.log("login successfully");
                server.initDatabase();
            } else {
                console.log("Login failed. Return code: " + returnCode);
            }
        })
        
        await new Promise((resolve) => {
            setTimeout(resolve, 1000000000)
        })
    }

}

export async function TestConnectEth() {
    const abiJson = require('../../../smartcontract/build/contracts/EurusERC20.json');

    let config :ServerConfigBase =  new ServerConfigBase();
    let server: ServerBase = new ServerBase(config);

    let isSuccess: boolean = server.loadConfig("AdminServerConfig.json");

    if (isSuccess){
        console.log("Login authServer");
        isSuccess = server.loginAuthServer(async (returnCode: ServerReturnCode) =>  {
            if (returnCode == ServerReturnCode.Success) {
                console.log("login successfully");
                
                let ethClient: EthClient = server.getSideChainEthClient();
                
                let erc20Instance : EurusERC20 = ethClient.NewSmartContract<EurusERC20>("EurusERC20", '0xa54Dee79c3bB34251DEbf86C1BA7D21898FFb7AC');

                let result = await erc20Instance.methods.balanceOf('0x5c28178B275657c0B71f1513A239705605988BD4').call()
                console.log("balance: ", result);
            } else {
                console.log("Login failed. Return code: " + returnCode);
            }
        })
        
        await new Promise((resolve) => {
            setTimeout(resolve, 1000000000)
        })
    }

}*/
