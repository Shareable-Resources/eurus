import { makeHandlerAwareOfAsyncErrors } from '../../../foundation/server/Middlewares';
import { UsersController } from '../controller/0_index';
import CommonRoute from '../../../foundation/server/CommonRoute';
import logger from '../util/ServiceLogger';
class Route extends CommonRoute {
  controller: UsersController = new UsersController();
  constructor() {
    super(logger);
    this.prefix = 'merchant-clients';
    this.setRoutes();
  }

  protected setRoutes() {
    this.setDefaultRoutes();
    this.router.post(
      `/${this.prefix}/importWallet`,
      makeHandlerAwareOfAsyncErrors(this.controller.importWallet, logger),
    );
  }
}

export default Route;
