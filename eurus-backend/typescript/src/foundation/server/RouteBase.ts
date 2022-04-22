/* eslint-disable @typescript-eslint/no-inferrable-types */
// src/routes/route.ts

import { Router } from 'express';
import CommonController from './CommonController';
import { makeHandlerAwareOfAsyncErrors } from './Middlewares';
import ExcludeRouteOption from './ExcludeRoutes';

abstract class RouteBase {
  protected controller: CommonController = new CommonController();
  protected router = Router();
  protected abstract setRoutes(controller: CommonController): void;
  protected prefix: string = '/';

  public getRouter() {
    return this.router;
  }
  public getPrefix() {
    return this.prefix;
  }

  protected setDefaultRoutes(exclude?: ExcludeRouteOption) {
    exclude = exclude ? exclude : new ExcludeRouteOption();
    if (this.controller.getAll && !exclude.getAll) {
      this.router.get(
        `/api/${this.prefix}`,
        makeHandlerAwareOfAsyncErrors(this.controller.getAll),
      );
    }
    if (this.controller.getById && !exclude.getById) {
      this.router.get(
        `/api/${this.prefix}/:id`,
        makeHandlerAwareOfAsyncErrors(this.controller.getById),
      );
    }
    if (this.controller.create && !exclude.create) {
      this.router.post(
        `/api/${this.prefix}`,
        makeHandlerAwareOfAsyncErrors(this.controller.create),
      );
    }
    if (this.controller.update && !exclude.update) {
      this.router.put(
        `/api/${this.prefix}/:id`,
        makeHandlerAwareOfAsyncErrors(this.controller.update),
      );
    }
    if (this.controller.remove && !exclude.remove) {
      this.router.delete(
        `/api/${this.prefix}/:id`,
        makeHandlerAwareOfAsyncErrors(this.controller.remove),
      );
    }
    if (this.controller.patch && !exclude.patch) {
      this.router.patch(
        `/api/${this.prefix}/:id`,
        makeHandlerAwareOfAsyncErrors(this.controller.patch),
      );
    }
  }

  constructor() {}
}

export default RouteBase;
