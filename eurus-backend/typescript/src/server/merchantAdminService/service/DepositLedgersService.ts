import * as SeqModel from '../model/seqModel/0_index';
import seq from '../sequelize';
import CommonService from '../../../foundation/server/CommonService';
import * as DBModel from '../model/dbModel/0_index';
import { ResponseBase } from '../../../foundation/server/ApiMessage';
import logger from '../util/ServiceLogger';
import { MerchantClient } from '../model/dbModel/0_index';
import { Transaction } from 'sequelize/types';
import JSONbig from 'json-bigint';
const modelModule = seq.sequelize.models;
export default class Service implements CommonService {
  async getAll(searchParams?: any): Promise<any> {
    const result: any = await modelModule[
      SeqModel.name.DepositLedgers
    ].findAll();
    return result;
  }
  //Non Api function
  async insertDepositLedger(obj: DBModel.DepositLedger): Promise<ResponseBase> {
    const msg = `[deposit_ledgers_service][insertDepositLedger](obj.txHash : ${obj.txHash})`;
    const resp = new ResponseBase();
    let t: Transaction = await seq.sequelize.transaction();
    try {
      //1. First check if user exists in db
      const userInDb: any = await modelModule[
        SeqModel.name.MerchantClient
      ].findByPk(obj.userId.toString(), { transaction: t });
      //2. If user exists in db, update its balance
      if (userInDb) {
        const userInDbInBigJsonObj: DBModel.MerchantClient = JSONbig.parse(
          JSONbig.stringify(userInDb),
        );
        userInDbInBigJsonObj.lastModifiedDate = new Date();
        userInDbInBigJsonObj.id = obj.userId;
        userInDbInBigJsonObj.balance =
          BigInt(userInDbInBigJsonObj.balance) + BigInt(obj.toTokenAmt);
        const result = await modelModule[SeqModel.name.MerchantClient].update(
          userInDbInBigJsonObj,
          {
            where: {
              id: userInDbInBigJsonObj.id,
            },
            transaction: t,
          },
        );

        logger.info(
          `[insertDepositLedger] MerchantClient, update balance, affected rows : (${result[0]})`,
        );
      } else {
        logger.info('user is not in database');
        logger.error(
          `Scanned a block, but the user who owned this wallet(${obj.toWalletAddr}) does not exist in the database`,
        );
      }
      //3. Insert new deposit_ledgers record
      const insertResult: any = await modelModule[
        SeqModel.name.DepositLedgers
      ].create(obj, { transaction: t }); //result would be the created Object
      resp.success = true;
      resp.respType = 'success';
      resp.data = insertResult;
      resp.msg = `New deposit ledger is created. New record txHash :(${insertResult.txHash})`;
      logger.info(msg + ' - success ', { message: resp.msg });
      await t.commit();
    } catch (e) {
      logger.error(msg + ' - fail ');
      logger.error(e);
      resp.success = false;
      resp.respType = 'error';
      resp.data = null;
      resp.msg = 'Insert deposit ledger fail :' + e;
      await t.rollback();
    }
    return resp;
  }
}
