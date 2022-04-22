import loggerHelper from '../../../foundation/utils/logger';
import { serverName } from '../const/index';
import config from '../config/HelloWorldServerConfig.json';
const logger = loggerHelper.createRotateLogger(
  serverName,
  undefined,
  config.winston.console,
);

export function writeToWinstonLogDebug(msg: any) {
  logger.debug(msg);
}

export default logger;
