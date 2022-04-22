import config from '../config/HelloWorldServerConfig.json';
const serviceConfigName = 'HelloWorldServerConfig';
const serverName = config.express.name;
export { serverName, serviceConfigName };
