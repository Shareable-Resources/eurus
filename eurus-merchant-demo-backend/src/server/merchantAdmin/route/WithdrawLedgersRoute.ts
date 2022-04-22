import { WithdrawLedgersController } from '../controller/0_index';
import CommonRoute from '../../../foundation/server/CommonRoute';
import logger from '../util/ServiceLogger';
class Route extends CommonRoute {
  controller: WithdrawLedgersController = new WithdrawLedgersController();
  constructor() {
    super(logger);
    this.prefix = 'withdraw-ledgers';
    this.setRoutes();
  }

  protected setRoutes() {
    this.setDefaultRoutes();
  }
}

export default Route;
