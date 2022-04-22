import { makeHandlerAwareOfAsyncErrors } from '../../../foundation/server/Middlewares';
import { HelloWorldUsersController } from '../controller/0_index';
import CommonRoute from '../../../foundation/server/CommonRoute';
import logger from '../util/ServiceLogger';
class Route extends CommonRoute {
  controller: HelloWorldUsersController = new HelloWorldUsersController();
  constructor() {
    super(logger);
    this.prefix = 'helloWorldUsers';
    this.setRoutes();
  }

  protected setRoutes() {
    this.setDefaultRoutes();
    this.router.get(
      `/${this.prefix}/hello`,
      makeHandlerAwareOfAsyncErrors(this.controller.hello, logger),
    );
    this.router.post(
      `/${this.prefix}/removeUser`,
      makeHandlerAwareOfAsyncErrors(this.controller.removeUser, logger),
    );
    this.router.post(
      `/${this.prefix}/updateUser`,
      makeHandlerAwareOfAsyncErrors(this.controller.updateUser, logger),
    );
    this.router.get(
      `/${this.prefix}/getAllByRaw`,
      makeHandlerAwareOfAsyncErrors(this.controller.getAllByRaw, logger),
    );
  }
}

export default Route;
