import {
  apiResponse,
  responseBySuccess,
  statusCode,
} from '../../../foundation/server/ApiMessage';
import CommonController from '../../../foundation/server/CommonController';
import Service from '../service/DbService';
import * as DBModel from '../model/dbModel/0_index';
import logger from '../util/ServiceLogger';

import { resetTable } from '../script/ResetTableScript';

export class Controller implements CommonController {
  /**
   * @description Creates an instance of merchant admin controller.
   *
   * @constructor
   * @param {Service} service
   */
  constructor(private service: Service = new Service()) {
    this.restart = this.restart.bind(this);
  }

  async restart(req: any, res: any) {
    logger.info('restarting DB');
    await resetTable();
    return apiResponse(
      res,
      responseBySuccess(null, true, 'query', 'success', 'DB Restarted'),
      statusCode(true, 'query'),
    );
  }

  async getAll() {
    throw Error('Not Implemented');
  }
}
export default Controller;
