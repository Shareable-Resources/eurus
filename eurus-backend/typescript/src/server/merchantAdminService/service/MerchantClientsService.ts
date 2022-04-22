import { ResponseBase } from '../../../foundation/server/ApiMessage';
import * as DBModel from '../model/dbModel/0_index';
import * as SeqModel from '../model/seqModel/0_index';
import seq from '../sequelize';
import CommonService from '../../../foundation/server/CommonService';
import * as ReqModel from '../model/reqModel/0_index';
import { Transaction } from 'sequelize/types';
import logger from '../util/ServiceLogger';
import { ServerReturnCode } from '../../../foundation/server/ServerReturnCode';
const modelModule = seq.sequelize.models;
export default class Service implements CommonService {
  async importWallet(obj: ReqModel.ImportWalletReq): Promise<ResponseBase> {
    const msg = `[merchant_clients_service][importWallet](obj.walletAddr : ${obj.walletAddress})`;
    obj.walletAddress = obj.walletAddress.toLowerCase();
    const resp = new ResponseBase();
    try {
      let t: Transaction = await seq.sequelize.transaction();
      const result: any = await modelModule[
        SeqModel.name.MerchantClient
      ].findOne({
        where: {
          walletAddress: obj.walletAddress,
        },
        transaction: t,
      });

      //Found user in merchant client table
      if (result) {
        if (result.username != obj.username) {
          resp.success = false;
          resp.respType = 'warning';
          resp.returnCode = ServerReturnCode.UniqueViolationError;
          resp.msg = 'Wallet address is already imported with another username';
          logger.info(msg + ' - fail ', { message: resp.msg });
        } else {
          resp.success = true;
          resp.respType = 'success';
          resp.msg = 'User found';
          resp.data = result;
          logger.info(msg + ' - success ', { message: resp.msg });
        }
      } else {
        //User cannot be found in merchant client table
        const result: any = await modelModule[
          SeqModel.name.MerchantClient
        ].findOne({
          where: {
            username: obj.username,
          },
          transaction: t,
        });
        if (result) {
          resp.success = false;
          resp.respType = 'warning';
          resp.msg = 'Username is taken, please retry another username';
          logger.info(msg + ' - fail ', { message: resp.msg });
        } else {
          const newUser = new DBModel.MerchantClient();
          newUser.username = obj.username;
          newUser.createdDate = new Date();
          newUser.lastModifiedDate = new Date();
          newUser.status = 1;
          newUser.walletAddress = obj.walletAddress;
          newUser.balance = BigInt(0);
          const insertResult: any = await modelModule[
            SeqModel.name.MerchantClient
          ].create(newUser, {
            returning: true,
            transaction: t,
          });
          resp.success = true;
          resp.respType = 'success';
          resp.msg = `User not found, new user is created. New record id :(${insertResult.id})`;
          resp.data = insertResult;
          logger.info(msg + ' - success ', { message: resp.msg });
        }
      }
      await t.commit();
    } catch (e) {
      logger.error(msg + ' - fail ');
      logger.error(e);
      resp.success = false;
      resp.respType = 'error';
      resp.data = null;
      resp.msg = 'Import wallet fail :' + e;
    }
    return resp;
  }

  async getAll(searchParams?: any): Promise<any> {
    const result: any = await modelModule[
      SeqModel.name.MerchantClient
    ].findAll();
    return result;
  }
  async getById(id: BigInt): Promise<any> {
    const result: any = await modelModule[
      SeqModel.name.MerchantClient
    ].findByPk(id.toString(), {});
    return result;
  }
}
