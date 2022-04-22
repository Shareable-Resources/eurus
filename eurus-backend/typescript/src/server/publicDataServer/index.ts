import { exit } from 'process';
import PublicDataServer, { PublicDataServerConfig } from './PublicDataServer';
import fs from 'fs';
import { ServerConfigBase } from '../../foundation/server/ServerConfigBase';
import loggerHelper from '../../foundation/utils/logger';

function init(){

  let args = process.argv.slice(2);
  if (args.length == 0){
    console.log("Usage: node index.js <config file path>");
    exit(1);
  }
  let configData = fs.readFileSync(args[0]);
  let config = JSON.parse(configData.toString());
  
  const logger = loggerHelper.createRotateLogger(config.logFilePath);

  console.log("Going to start server");
  logger.info("Going to start server");

  let server = new PublicDataServer(config as PublicDataServerConfig, logger);
  server.setupRouter();
  server.initHttpEthClient();
  server.listen();
}

init();
