import { DbController } from '../controller/0_index';
import CommonRoute from '../../../foundation/server/CommonRoute';
import logger from '../util/ServiceLogger';
import { makeHandlerAwareOfAsyncErrors } from '../../../foundation/server/Middlewares';
class Route extends CommonRoute {
  controller: DbController = new DbController();
  constructor() {
    super(logger);
    this.prefix = 'db';
    this.setRoutes();
  }

  protected setRoutes() {
    this.router.post(
      `/${this.prefix}/restart`,
      makeHandlerAwareOfAsyncErrors(this.controller.restart, logger),
    );
  }
}

export default Route;
