import { ServerConfigBase } from '../../foundation/server/ServerConfigBase';
import seq from './sequelize';
import expressServer, { printRoutes } from './server';
import { ContentTypeMiddleWare } from '../../foundation/server/Middlewares';
import config from './config/HelloWorldServerConfig.json';
import { EthClient } from '../../foundation/utils/ethereum/EthClient';
import foundationConst from '../../foundation/const';
import { HelloWorldServer } from './model/HelloWorldServer';
import Web3 from 'web3';
import logger from './util/ServiceLogger';

let serverConfig: ServerConfigBase = new ServerConfigBase();
let serverBase: HelloWorldServer = new HelloWorldServer(serverConfig);
const assertDatabaseConnectionOk = async () => {
  serverBase.db = seq;
  await serverBase.db.assertDatabaseConnectionOk();
};

async function init() {
  try {
    await assertDatabaseConnectionOk();
    //Middleware should add to expressServer
    expressServer.use(ContentTypeMiddleWare);
    serverBase.server = expressServer;
    //expressServer listening
    serverBase.httpServer = expressServer.listen(config.express.port, () => {
      printRoutes(expressServer);
    });

    //2. Observer (Optional) - observer block in a eth chain
    const webSocketProvider = new Web3.providers.WebsocketProvider(
      foundationConst.EthClientProviderUrl.webSocket,
    );
    const ethClient = new EthClient(
      webSocketProvider,
      foundationConst.ChainId.SideChain,
    );
    serverBase.ethClient = ethClient;
  } catch (e) {
    logger.error(e);
  }
}

init();
