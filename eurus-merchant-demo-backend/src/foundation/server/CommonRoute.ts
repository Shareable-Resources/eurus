/* eslint-disable @typescript-eslint/no-inferrable-types */

import { Router } from 'express';
import CommonController from './CommonController';
import { makeHandlerAwareOfAsyncErrors } from './Middlewares';
import winston from 'winston';

abstract class Route {
  protected controller: CommonController = new CommonController();
  protected router = Router();
  protected abstract setRoutes(controller: CommonController): void;
  protected prefix: string = '/';
  protected logger: winston.Logger;

  public getRouter() {
    return this.router;
  }
  public getPrefix() {
    return this.prefix;
  }

  protected setDefaultRoutes() {
    if (this.controller.getAll) {
      this.router.get(
        `/${this.prefix}`,
        makeHandlerAwareOfAsyncErrors(this.controller.getAll)
      );
    }
    if (this.controller.getById) {
      this.router.get(
        `/${this.prefix}/:id`,
        makeHandlerAwareOfAsyncErrors(this.controller.getById)
      );
    }
    if (this.controller.create) {
      this.router.post(
        `/${this.prefix}`,
        makeHandlerAwareOfAsyncErrors(this.controller.create, this.logger)
      );
    }
  }

  constructor(logger: winston.Logger) {
    this.logger = logger;
  }
}

export default Route;
