import { makeHandlerAwareOfAsyncErrors } from '../../../foundation/server/Middlewares';
import { MerchantAdminsController } from '../controller/0_index';
import CommonRoute from '../../../foundation/server/CommonRoute';
import logger from '../util/ServiceLogger';
class Route extends CommonRoute {
  controller: MerchantAdminsController = new MerchantAdminsController();
  constructor() {
    super(logger);
    this.prefix = 'merchant-admins';
    this.setRoutes();
  }

  protected setRoutes() {
    this.router.post(
      `/${this.prefix}/login`,
      makeHandlerAwareOfAsyncErrors(this.controller.login, logger),
    );
  }
}

export default Route;
