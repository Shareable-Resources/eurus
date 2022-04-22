import { ResponseBase } from '../../../foundation/server/ApiMessage';
import { ServerReturnCode } from '../../../foundation/server/ServerReturnCode';
import * as DBModel from '../model/dbModel/0_index';
import * as SeqModel from '../model/seqModel/0_index';
import seq from '../sequelize';
import logger from '../util/ServiceLogger';
import CommonService from '../../../foundation/server/CommonService';
const models = seq.sequelize.models;
export default class Service implements CommonService {
  async login(merchantAdmin: DBModel.MerchantAdmin): Promise<ResponseBase> {
    const msg = `[merchant_admins_service][importWallet](obj.username : ${merchantAdmin.username})`;
    const resp = new ResponseBase();
    const result: any = await models[SeqModel.name.MerchantAdmins].findOne({
      where: {
        merchantId: merchantAdmin.merchantId,
        username: merchantAdmin.username,
      },
    });

    if (result) {
      const merchantAdminInDb: DBModel.MerchantAdmin = result;
      if (merchantAdminInDb.passwordHash != merchantAdmin.passwordHash) {
        resp.success = false;
        resp.respType = 'warning';
        resp.returnCode = ServerReturnCode.InvalidArgument;
        resp.msg = 'Password mismatched';
        logger.info(msg + ' - fail ', { message: resp.msg });
      } else {
        const { passwordHash, ...removeSensitiveObj } = JSON.parse(
          JSON.stringify(merchantAdminInDb),
        );
        resp.success = true;
        resp.respType = 'success';
        resp.msg = 'User Found';
        resp.data = removeSensitiveObj;
        logger.info(msg + ' - success ', { message: resp.msg });
      }
    } else {
      resp.success = false;
      resp.respType = 'warning';
      resp.returnCode = ServerReturnCode.UserNotFound;
      resp.msg = 'User record not found';
      logger.info(msg + ' - fail ', { message: resp.msg });
    }
    return resp;
  }

  async getAll(searchParams: any): Promise<any> {
    throw Error('Not Implemented');
  }
}
