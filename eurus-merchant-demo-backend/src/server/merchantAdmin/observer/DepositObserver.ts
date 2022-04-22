import { EthClient } from '../../../foundation/utils/ethereum/0_index';
import CommonObserver from './CommonObserver';
import SmartContractDeploy from '../../../smartContract/contracts/SmartContractDeploy.json';
import { AbiItem } from 'web3-utils';
import Service from '../service/DepositLedgersService';
import * as DBModel from '../model/dbModel/0_index';
import logger from '../util/ServiceLogger';
import {
  DepositedToDApp,
  EurusERC20,
} from '../../../smartContract/typescript/EurusERC20';
import EurusERC20Contract from '../../../smartContract/contracts/EurusERC20.json';
import { DepositLedgerStatus } from '../model/dbModel/DepositLedger';
import Web3 from 'web3';
import { EventData } from 'web3-eth-contract';

export class DepositObserver extends CommonObserver {
  constructor(ethClient: EthClient, private service: Service = new Service()) {
    super(ethClient);
  }

  async startMonitoringSmartContract() {
    logger.info('[depositObserver] Starts monitoring block');
    const smartContractName = 'OwnedUpgradeabilityProxy<USDT>';
    const chainID = await this.ethClient.web3Client.eth.getChainId();
    logger.info(`ChainId ${chainID}`);
    const abiItems: AbiItem[] = EurusERC20Contract.abi as AbiItem[];
    const web3Contract = new this.ethClient.web3Client.eth.Contract(
      abiItems,
      SmartContractDeploy[chainID]['smartContract'][smartContractName].address
    );
    const erc20Contract: EurusERC20 = web3Contract as any;
    /*
        EventEmitter: The event emitter has the following events:

        "data" returns Object: Fires on each incoming event with the event object as argument.
        "changed" returns Object: Fires on each event which was removed from the blockchain. The event will have the additional property "removed: true".
        "error" returns Object: Fires when an error in the subscription occours.
        "connected" returns String: Fires once after the subscription successfully connected. Returns the subscription id.
    */
    erc20Contract.events
      .DepositedToDApp()
      .on('data', async (event: EventData) => {
        logger.info('erc20Contract [data]', event);
        await this.proccessBlock(event);
      })
      .on('changed', (event) => {
        logger.info('erc20Contract [changed]', event);
      })
      .on('error', (event) => {
        logger.error('erc20Contract [error]', event);
      })
      .on('connected', (event) => {
        logger.info('erc20Contract [connected]', event);
      });
  }

  async proccessBlock(event: EventData) {
    const dateNow = new Date();
    const obj = new DBModel.DepositLedger();
    const returnValues = event.returnValues;
    obj.txHash = event.transactionHash;
    obj.fromAssetAmt = BigInt(returnValues.amount);
    obj.toTokenAmt = BigInt(returnValues.amount);
    obj.fromAssetId = returnValues.symbol;
    obj.fromWalletAddr = returnValues.buyer;
    obj.fromAssetAddr = event.address;
    obj.toWalletAddr = returnValues.dappAddress;
    obj.rate = BigInt(1);
    obj.createdDate = dateNow;
    obj.lastModifiedDate = dateNow;
    obj.remarks = `Block number : ${event.blockNumber}`;
    obj.status = DepositLedgerStatus.DepositCompleted;
    obj.userId = BigInt(Web3.utils.hexToNumber(returnValues.extraData));
    const isSuccess = await this.service.insertDepositLedger(obj);
  }
}
