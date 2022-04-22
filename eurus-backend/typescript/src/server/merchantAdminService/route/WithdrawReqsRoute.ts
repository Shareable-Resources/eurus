import { WithdrawReqsController } from '../controller/0_index';
import CommonRoute from '../../../foundation/server/CommonRoute';
import logger from '../util/ServiceLogger';
import { makeHandlerAwareOfAsyncErrors } from '../../../foundation/server/Middlewares';
class Route extends CommonRoute {
  controller: WithdrawReqsController = new WithdrawReqsController();
  constructor() {
    super(logger);
    this.prefix = 'withdraw-reqs';
    this.setRoutes();
  }

  protected setRoutes() {
    this.setDefaultRoutes();
  }
}

export default Route;
