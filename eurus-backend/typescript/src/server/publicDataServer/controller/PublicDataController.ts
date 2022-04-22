import winston  from 'winston';
import { EthClient } from '../../../foundation/utils/ethereum/EthClient';
import PublicDataServer from '../PublicDataServer';
import Big from 'big.js';
import { Response } from 'express';
import { StatusCodes } from 'http-status-codes';

export class PublicDataController {
    logger : winston.Logger;
    server : PublicDataServer;
    constructor (server : PublicDataServer){
        this.logger = server.logger!;
        this.server = server;

        this.getEUNAvailability = this.getEUNAvailability.bind(this);
    }

    public async getEUNAvailability(req: any, res: Response){
        let totalBalance: Big = new Big(0);

        try{
            let blockNumber : number | undefined = await this.server.ethClient?.web3Client.eth.getBlockNumber();
            if (blockNumber){
                let bigBlockNumber : Big = new Big(blockNumber);
                let avilEun = bigBlockNumber.mul(3000000000000000000);
                totalBalance = avilEun.plus(this.server.config.initializeEunAmount);
                totalBalance = totalBalance.mul(Math.pow(10, -18));
            }else{
                throw new Error("RPC error");
            }
        }catch (e : any ){
            let err = e as Error;
            throw new Error(err.message);
        }   
      
        return res.format({
            text:()=> {
                Big.PE = 40;
                res.status(StatusCodes.OK).send(totalBalance.toString())
            }
        });
        
    }
}