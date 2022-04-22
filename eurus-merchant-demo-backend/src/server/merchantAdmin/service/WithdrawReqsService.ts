import seq from '../sequelize';
import CommonService from '../../../foundation/server/CommonService';
import * as DBModel from '../model/dbModel/0_index';
import * as SeqModel from '../model/seqModel/0_index';
import { Transaction } from 'sequelize/types';
import { ResponseBase } from '../../../foundation/server/ApiMessage';
import logger from '../util/ServiceLogger';
import { WithdrawReqStatus } from '../model/dbModel/WithdrawReq';
const modelModule = seq.sequelize.models;
export default class Service implements CommonService {
  async getAll(searchParams?: any): Promise<any> {
    // Raw SQL Query
    const [results, metadata]: any = await seq.sequelize.query(
      `select req.id as "id", req.last_modified_date as "lastModifiedDate", req.from_token_amt as "fromTokenAmt", req.to_wallet_addr as "fromWalletAddr", req.to_asset_id as "toAssetId", req.to_asset_addr as "toAssetAddr", req.user_id as "userId", mc.username as "username", req.status as "status" from "merchant_demo"."withdraw_reqs" as req left join "merchant_demo"."merchant_clients" mc on req.user_id =mc.id`
    );
    return results;
  }

  async create(obj: DBModel.WithdrawReq): Promise<ResponseBase> {
    const resp = new ResponseBase();
    const msg = `[withdraw_reqs_service][create](userId : ${obj.userId})`;
    obj.status = WithdrawReqStatus.StatusPendingApproval;
    obj.rate = BigInt(1);
    let t: Transaction = await seq.sequelize.transaction();
    try {
      const insertResult: any = await modelModule[
        SeqModel.name.WithdrawReqs
      ].create(obj); //result would be the created Object
      resp.data = insertResult;
      resp.msg = `New withdraw request is created. New record id :(${insertResult.id})`;
      resp.respType = 'success';
      resp.success = true;
      logger.info(msg + ' - success ', { message: resp.msg });
      await t.commit();
    } catch (e) {
      logger.error(msg + ' - fail ');
      logger.error(e);
      resp.success = false;
      resp.respType = 'error';
      resp.data = null;
      resp.msg = 'Insert withdraw request fail :' + e;
      await t.rollback();
    }
    return resp;
  }
}
