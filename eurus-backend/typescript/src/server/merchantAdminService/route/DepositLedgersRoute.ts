import { DepositLedgersController } from '../controller/0_index';
import CommonRoute from '../../../foundation/server/CommonRoute';
import logger from '../util/ServiceLogger';
import { makeHandlerAwareOfAsyncErrors } from '../../../foundation/server/Middlewares';
class Route extends CommonRoute {
  controller: DepositLedgersController = new DepositLedgersController();
  constructor() {
    super(logger);
    this.prefix = 'deposit-ledgers';
    this.setRoutes();
  }

  protected setRoutes() {
    this.setDefaultRoutes();
  }
}

export default Route;
