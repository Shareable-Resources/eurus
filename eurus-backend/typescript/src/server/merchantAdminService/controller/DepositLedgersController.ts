import {
  apiResponse,
  responseBySuccess,
  statusCode,
} from '../../../foundation/server/ApiMessage';
import CommonController from '../../../foundation/server/CommonController';
import Service from '../service/DepositLedgersService';

import * as DBModel from '../model/dbModel/0_index';
import { DepositLedgerStatus } from '../model/dbModel/DepositLedger';
export class Controller implements CommonController {
  /**
   * @description Creates an instance of deposit ledgers controller.
   *
   * @constructor
   * @param {Service} service
   */
  constructor(private service: Service = new Service()) {
    this.getAll = this.getAll.bind(this);
  }

  async getAll(req: any, res: any) {
    const list = await this.service.getAll();
    return apiResponse(
      res,
      responseBySuccess(list, true, 'query', 'success'),
      statusCode(true, 'query'),
    );
  }
}
export default Controller;
