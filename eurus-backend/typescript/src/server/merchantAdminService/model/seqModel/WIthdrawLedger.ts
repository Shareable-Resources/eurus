/* eslint-disable @typescript-eslint/consistent-type-assertions */
import { BuildOptions, DataTypes, Model, Sequelize } from 'sequelize';
import WithdrawLedger from '../dbModel/WithdrawLedger';
import * as SeqModel from './0_index';
export default WithdrawLedger;
export interface WithdrawLedgerModel
  extends Model<WithdrawLedger>,
    WithdrawLedger {}
export type WithdrawLedgerStatic = typeof Model & {
  new (values?: object, options?: BuildOptions): WithdrawLedgerModel;
};
export function WithdrawLedgerFactory(sequelize: Sequelize) {
  return <WithdrawLedgerStatic>sequelize.define(
    SeqModel.name.WithdrawLedgers,
    {
      txHash: {
        field: 'tx_hash',
        allowNull: false,
        primaryKey: true,
        type: DataTypes.STRING(255),
      },
      fromWalletAddr: {
        field: 'from_asset_id',
        allowNull: false,
        type: DataTypes.STRING(255),
      },
      fromTokenAmt: {
        field: 'from_token_amt',
        allowNull: false,
        type: DataTypes.DECIMAL(78),
      },
      toWalletAddr: {
        field: 'to_wallet_addr',
        allowNull: false,
        type: DataTypes.STRING(255),
      },
      toAssetId: {
        field: 'to_asset_id',
        allowNull: false,
        type: DataTypes.STRING(20),
      },
      toAssetAddr: {
        field: 'toAssetAddr',
        allowNull: false,
        type: DataTypes.STRING(255),
      },
      toUserId: {
        type: DataTypes.BIGINT,
        allowNull: true,
        field: 'to_user_id',
      },
      rate: {
        field: 'rate',
        allowNull: false,
        type: DataTypes.DECIMAL(78),
      },
      createdDate: {
        type: DataTypes.DATE,
        defaultValue: DataTypes.NOW,
        field: 'created_date',
      },
      lastModifiedDate: {
        type: DataTypes.DATE,
        defaultValue: DataTypes.NOW,
        field: 'last_modified_date',
      },
      status: {
        type: DataTypes.SMALLINT,
        allowNull: false,
        field: 'status',
      },
      remarks: {
        type: DataTypes.STRING(1000),
        allowNull: true,
        field: 'remarks',
      },
      reqId: {
        type: DataTypes.BIGINT,
        allowNull: true,
        field: 'req_id',
      },
    },
    {
      freezeTableName: true,
      updatedAt: false,
      createdAt: false,
    },
  );
}
