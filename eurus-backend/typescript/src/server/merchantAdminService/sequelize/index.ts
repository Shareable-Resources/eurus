import config from '../config/MerchantAdminServerConfig.json';
import * as SeqModel from '../model/seqModel/0_index';
import Db from '../../../foundation/sequlelize';
import logger from '../util/ServiceLogger';

export class ServerDb extends Db {
  constructor() {
    super(config, 'postgres', logger);
  }

  bindModelsToSeq() {
    SeqModel.factory.MerchantClientFactory(this.sequelize);
    SeqModel.factory.MerchantAdminFactory(this.sequelize);
    SeqModel.factory.DepositLedgerFactory(this.sequelize);
    SeqModel.factory.WithdrawLedgerFactory(this.sequelize);
    SeqModel.factory.WithdrawReqFactory(this.sequelize);
  }
}

const db = new ServerDb();
export default db;
