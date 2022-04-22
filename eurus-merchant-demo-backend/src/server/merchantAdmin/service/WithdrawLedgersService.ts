import { ResponseBase } from '../../../foundation/server/ApiMessage';
import * as DBModel from '../model/dbModel/0_index';
import * as SeqModel from '../model/seqModel/0_index';
import seq from '../sequelize';
import CommonService from '../../../foundation/server/CommonService';
import logger from '../util/ServiceLogger';
import { Transaction } from 'sequelize/types';
import JSONbig from 'json-bigint';
const modelModule = seq.sequelize.models;
export default class Service implements CommonService {
  async getAll(searchParams?: any): Promise<any> {
    const result: any = await modelModule[
      SeqModel.name.WithdrawLedgers
    ].findAll();
    return result;
  }
  //Non API function
  async insertWithdrawLedger(
    obj: DBModel.WithdrawLedger,
    withdrawReq: DBModel.WithdrawReq
  ): Promise<ResponseBase> {
    const msg = `[withdraw_ledgers_serveice][insertWithdrawLedger](obj.txHash : ${obj.txHash})`;
    const resp = new ResponseBase();
    let t: Transaction = await seq.sequelize.transaction();
    let result: any;
    let rowsAffected: number = 0;
    obj.reqId = withdrawReq.id;
    try {
      // 1. First check if user exists in db
      const userInDb: any = await modelModule[
        SeqModel.name.MerchantClient
      ].findOne({
        where: { walletAddress: obj.toWalletAddr },
        transaction: t,
      });
      //2. If user exists in db, update its balance
      if (userInDb) {
        logger.info('user is in database');
        obj.toUserId = userInDb.id;
        const userInDbInBigJsonObj: DBModel.MerchantClient = JSONbig.parse(
          JSONbig.stringify(userInDb)
        );
        userInDbInBigJsonObj.lastModifiedDate = new Date();
        userInDbInBigJsonObj.balance =
          BigInt(userInDbInBigJsonObj.balance) - BigInt(obj.fromTokenAmt);
        result = await modelModule[SeqModel.name.MerchantClient].update(
          userInDbInBigJsonObj,
          {
            where: {
              id: userInDbInBigJsonObj.id,
            },
            transaction: t,
          }
        );
        rowsAffected += result[0];
        logger.info(
          `[insertWithdrawLedger] MerchantClient, update balance, affected rows : (${result[0]})`
        );
      } else {
        logger.info('user is not in database');
        logger.error(
          `Scanned a block, but the user who owned this wallet(${obj.toWalletAddr}) does not exist in the database`
        );
      }

      //3. Update tx_hash in withdraw_reqs and and insert new withdraw_ledgers record
      result = await modelModule[SeqModel.name.WithdrawReqs].update(
        withdrawReq,
        {
          fields: [
            'approveBy',
            'approveDate',
            'txHash',
            'lastModifiedDate',
            'status',
          ],
          where: {
            id: withdrawReq.id,
          },
          transaction: t,
        }
      );
      rowsAffected += result[0];
      logger.info(
        `WithdrawReqs, update approve status, affected rows : (${result[0]})`
      );

      const insertResult: any = await modelModule[
        SeqModel.name.WithdrawLedgers
      ].create(obj, { transaction: t }); //result would be the created Object

      resp.msg = `New withdraw ledger is created. New record txHash :(${insertResult.txHash})`;
      logger.info(msg + ' - success ', { message: resp.msg });
      await t.commit();
    } catch (e) {
      logger.error(msg + ' - fail ');
      logger.error(e);
      resp.success = false;
      resp.respType = 'error';
      resp.data = null;
      resp.msg = 'Insert withdraw ledger fail :' + e;
      await t.rollback();
    }
    return resp;
  }
}
