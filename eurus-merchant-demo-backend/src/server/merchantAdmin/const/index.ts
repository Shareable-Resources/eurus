import configJSON from '../config/MerchantAdminServerConfig.json';
const config =
  configJSON[process.env.NODE_ENV ? process.env.NODE_ENV : 'local'];
const serverName = config.express.name;
export { serverName };
