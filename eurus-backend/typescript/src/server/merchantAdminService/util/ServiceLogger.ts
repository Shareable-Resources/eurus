import loggerHelper from '../../../foundation/utils/logger';
import { serverName } from '../const/index';
import config from '../config/MerchantAdminServerConfig.json';
const logger = loggerHelper.createRotateLogger(
  serverName,
  undefined,
  config.winston.console,
);

export default logger;
