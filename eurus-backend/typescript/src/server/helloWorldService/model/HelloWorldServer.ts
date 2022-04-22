import winston from 'winston';
import { ServerBase } from '../../../foundation/server/ServerBase';
import { ServerConfigBase } from '../../../foundation/server/ServerConfigBase';
import logger from '../util/ServiceLogger';

export class HelloWorldServer extends ServerBase {
  logger?: winston.Logger; // you can import this logger by import or just use HelloWorldUserServer.logger by instance
  constructor(configBase: ServerConfigBase) {
    super(configBase);
    this.logger = logger;
  }
}
