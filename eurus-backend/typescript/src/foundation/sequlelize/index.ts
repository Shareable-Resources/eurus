import { Dialect, Options, Sequelize } from 'sequelize';
import winston from 'winston';

export default abstract class Db {
  sequelize: Sequelize;
  logger: winston.Logger;
  constructor(
    config: any,
    dialect: Dialect | undefined,
    logger: winston.Logger,
  ) {
    this.logger = logger;
    const parsedConfig = config.sequelize as Options;
    parsedConfig.dialect = dialect;
    parsedConfig.dialectOptions = {
      supportBigNumbers: true,
      bigNumberStrings: true,
    };
    parsedConfig.hooks = {};
    parsedConfig.logging = (msg) => {
      this.logger.debug(msg);
    };
    this.sequelize = new Sequelize(
      config.sequelize.database,
      config.sequelize.username,
      config.sequelize.password,
      parsedConfig,
    );
    this.bindModelsToSeq();
  }

  public async assertDatabaseConnectionOk() {
    try {
      await this.sequelize?.authenticate();
    } catch (error) {
      this.logger.error(error);
      process.exit(1);
    }
  }

  abstract bindModelsToSeq();
  /*
  Implement those model u need in micro service
  bindModelsToSeq() {
    SeqModel.factory.MerchantClientFactory(this.sequelize);
    SeqModel.factory.MerchantAdminFactory(this.sequelize);
    SeqModel.factory.DepositLedgerFactory(this.sequelize);
    SeqModel.factory.WithdrawLedgerFactory(this.sequelize);
    SeqModel.factory.WithdrawReqFactory(this.sequelize);
  }
  */
}
