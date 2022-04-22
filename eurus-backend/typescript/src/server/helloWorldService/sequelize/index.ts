import config from '../config/HelloWorldServerConfig.json';
import * as SeqModel from '../model/seqModel/0_index';
import Db from '../../../foundation/sequlelize';
import logger from '../util/serviceLogger';
export class ServerDb extends Db {
  constructor() {
    super(config, 'postgres', logger);
  }

  bindModelsToSeq() {
    SeqModel.factory.HelloWorldUserFactory(this.sequelize);
  }
}

const db = new ServerDb();
export default db;
