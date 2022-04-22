import {
  apiResponse,
  responseBySuccess,
  statusCode,
} from '../../../foundation/server/ApiMessage';
import CommonController from '../../../foundation/server/CommonController';
import logger from '../util/ServiceLogger';

import { resetTable } from '../script/ResetTableScript';

export class Controller implements CommonController {
  /**
   * @description Creates an instance of hello world db controller.
   *
   * @constructor
   *
   */
  constructor() {
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
