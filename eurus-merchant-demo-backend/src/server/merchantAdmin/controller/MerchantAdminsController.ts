import {
  apiResponse,
  responseBySuccess,
  statusCode,
} from '../../../foundation/server/ApiMessage';
import CommonController from '../../../foundation/server/CommonController';
import Service from '../service/MerchantAdminsService';
import * as DBModel from '../model/dbModel/0_index';
export class Controller implements CommonController {
  /**
   * @description Creates an instance of merchant admin controller.
   *
   * @constructor
   * @param {Service} service
   */
  constructor(private service: Service = new Service()) {
    this.login = this.login.bind(this);
  }

  async login(req: any, res: any) {
    const obj: DBModel.MerchantAdmin = req.body;
    obj.merchantId = BigInt(1);
    const responseBase = await this.service.login(obj);
    return apiResponse(
      res,
      responseBySuccess(
        responseBase.data,
        responseBase.success,
        'query',
        responseBase.respType,
        responseBase.msg,
      ),
      statusCode(responseBase.success, 'query'),
    );
  }

  async getAll() {
    throw Error('Not Implemented');
  }
}
export default Controller;
