import { EthClient } from '../../../foundation/utils/ethereum/0_index';
import CommonObserver from './CommonObserver';
import DAppSampleContract from '../../../foundation/smartContract/DAppSample.json';
import SmartContractDeploy from '../../../foundation/smartContract/SmartContractDeploy.json';
import { AbiItem } from 'web3-utils';
import Service from '../service/WithdrawLedgersService';
import * as DBModel from '../model/dbModel/0_index';
import logger from '../util/ServiceLogger';
import {
  DAppSample,
  refunded,
} from '../../../../../smartcontract/build/typescript/DAppSample';
import { WithdrawReqStatus } from '../model/dbModel/WithdrawReq';
import { WithdrawLedgerStatus } from '../model/dbModel/WithdrawLedger';
import { EventData } from 'web3-eth-contract';
import Web3 from 'web3';
import { ContractEventLog } from '../../../../../smartcontract/build/typescript/types';
export class WithdrawObserver extends CommonObserver {
  constructor(ethClient: EthClient, private service: Service = new Service()) {
    super(ethClient);
  }

  async startMonitoringSmartContract() {
    logger.info('[withdrawObserver] Starts monitoring block');
    const smartContractName = 'OwnedUpgradeabilityProxy<DAppSampleToken>';
    const chainID = await this.ethClient.web3Client.eth.getChainId();
    const abiItems: AbiItem[] = DAppSampleContract.abi as AbiItem[];
    const web3Contract = new this.ethClient.web3Client.eth.Contract(
      abiItems,
      SmartContractDeploy[chainID]['smartContract'][smartContractName].address,
    );
    const dAppContract: DAppSample = web3Contract as any;

    dAppContract.events
      .refunded()
      .on('data', async (event: EventData) => {
        logger.info('dAppContract [data]', event);
        await this.proccessBlock(event);
      })
      .on('changed', (event) => {
        logger.info('dAppContract [changed]', event);
      })
      .on('error', (event) => {
        logger.error('dAppContract [error]', event);
      })
      .on('connected', (event) => {
        logger.info('dAppContract [connected]', event);
      });
  }

  async proccessBlock(event: EventData) {
    const dateNow = new Date();
    const obj = new DBModel.WithdrawLedger();
    // Gets the returnValues from refunded event
    const specificEvent:refunded = event as any;
    const returnValues= specificEvent.returnValues;
    const concatedStr = Web3.utils.hexToUtf8(returnValues.extraData);
    const splitArray = concatedStr.split('|');
    const reqId = splitArray[0];
    const approveBy = splitArray[1];
    obj.txHash = event.transactionHash;
    obj.fromWalletAddr = event.address;
    obj.fromTokenAmt = BigInt(returnValues.srcAmount);
    obj.toWalletAddr = returnValues.dest.toLowerCase();
    obj.toAssetId = returnValues.targetAssetName;
    obj.toAssetAmt = BigInt(returnValues.targetAmount);
    obj.rate = BigInt(1);
    obj.reqId = BigInt(reqId);
    obj.createdDate = dateNow;
    obj.lastModifiedDate = dateNow;
    obj.status = WithdrawLedgerStatus.StatusCompleted;
    obj.remarks = `Block number : ${event.blockNumber}`;
    const withdrawReq: DBModel.WithdrawReq = new DBModel.WithdrawReq();
    withdrawReq.approveBy = BigInt(approveBy);
    withdrawReq.lastModifiedDate = dateNow;
    withdrawReq.approveDate = dateNow;
    withdrawReq.txHash = event.transactionHash;
    withdrawReq.id = BigInt(reqId);
    withdrawReq.status = WithdrawReqStatus.StatusApproved;
    const isSuccess = await this.service.insertWithdrawLedger(obj, withdrawReq);
  }
}
