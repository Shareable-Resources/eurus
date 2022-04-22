import { LogController } from '../controller/0_index';
import CommonRoute from '../../../foundation/server/CommonRoute';
import logger from '../util/ServiceLogger';
import { makeHandlerAwareOfAsyncErrors } from '../../../foundation/server/Middlewares';
class Route extends CommonRoute {
  controller: LogController = new LogController();
  constructor() {
    super(logger);
    this.prefix = 'logs';
    this.setRoutes();
  }

  protected setRoutes() {
    this.router.get(
      `/${this.prefix}/read`,
      makeHandlerAwareOfAsyncErrors(this.controller.read, logger),
    );
  }
}

export default Route;
