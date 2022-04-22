import seq from './sequelize';
import expressServer, { printRoutes } from './server';
import { ContentTypeMiddleWare } from '../../foundation/server/Middlewares';
import configJSON from './config/MerchantAdminServerConfig.json';
import { EthClient } from '../../foundation/utils/ethereum/EthClient';
import foundationConst from '../../foundation/const';
import { MerchantAdminServer } from './model/MerchantAdminServer';
import { DepositObserver } from './observer/DepositObserver';
import { WithdrawObserver } from './observer/WithdrawObserver';
import Web3 from 'web3';
import logger from './util/ServiceLogger';
const config =
  configJSON[process.env.NODE_ENV ? process.env.NODE_ENV : 'local'];

let serverBase: MerchantAdminServer = new MerchantAdminServer();
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
      foundationConst.EthClientProviderUrl.webSocket
    );
    const ethClient = new EthClient(
      webSocketProvider,
      foundationConst.ChainId.DevChain
    );
    serverBase.ethClient = ethClient;
    serverBase.depositObserver = new DepositObserver(serverBase.ethClient);
    serverBase.withdrawObserver = new WithdrawObserver(serverBase.ethClient);
    let maxTries = 10;
    try {
      await serverBase.depositObserver.startMonitoringSmartContract();
      await serverBase.withdrawObserver.startMonitoringSmartContract();
    } catch (e) {
      logger.error('Fail to observe');
      logger.error(e);
      if (maxTries > 10) {
        throw e;
      }
      maxTries++;
    }
  } catch (e) {
    logger.error(e);
  }
}

init();
