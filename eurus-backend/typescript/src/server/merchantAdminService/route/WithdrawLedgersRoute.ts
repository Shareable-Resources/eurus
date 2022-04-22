import { WithdrawLedgersController } from '../controller/0_index';
import CommonRoute from '../../../foundation/server/CommonRoute';
import logger from '../util/ServiceLogger';
import { makeHandlerAwareOfAsyncErrors } from '../../../foundation/server/Middlewares';
class Route extends CommonRoute {
  controller: WithdrawLedgersController = new WithdrawLedgersController();
  constructor() {
    super(logger);
    this.prefix = 'withdraw-ledgers';
    this.setRoutes();
  }

  protected setRoutes() {
    this.setDefaultRoutes();
    this.router.post(
      `/${this.prefix}/insertWithdrawLedger`,
      makeHandlerAwareOfAsyncErrors(
        this.controller.insertWithdrawLedger,
        logger,
      ),
    );
  }
}

export default Route;
