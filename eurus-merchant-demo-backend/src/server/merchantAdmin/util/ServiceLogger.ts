import loggerHelper from '../../../foundation/utils/logger';
import { serverName } from '../const/index';
import configJSON from '../config/MerchantAdminServerConfig.json';
const config =
  configJSON[process.env.NODE_ENV ? process.env.NODE_ENV : 'dev'];
const logger = loggerHelper.createRotateLogger(
  serverName,
  undefined,
  config.winston.console,
);

export default logger;
