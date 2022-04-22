import {
  apiResponse,
  ResponseBase,
  responseBySuccess,
  statusCode,
} from '../../../foundation/server/ApiMessage';
import CommonController from '../../../foundation/server/CommonController';
import Service from '../service/WithdrawReqsService';
import * as DBModel from '../model/dbModel/0_index';
export class Controller implements CommonController {
  /**
   * @description Creates an instance of withdraw ledgers controller.
   *
   * @constructor
   * @param {Service} service
   */
  constructor(private service: Service = new Service()) {
    this.getAll = this.getAll.bind(this);
    this.create = this.create.bind(this);
  }

  async getAll(req: any, res: any) {
    const list = await this.service.getAll();
    return apiResponse(
      res,
      responseBySuccess(list, true, 'query', 'success'),
      statusCode(true, 'query')
    );
  }

  async create(req: any, res: any) {
    const reqObj: DBModel.WithdrawReq = req.body;
    const responseObj: ResponseBase = await this.service.create(reqObj);
    return apiResponse(
      res,
      responseBySuccess(
        responseObj.data,
        responseObj.success,
        'add',
        responseObj.respType,
        responseObj.msg
      ),
      statusCode(responseObj.success, 'add')
    );
  }
}
export default Controller;
