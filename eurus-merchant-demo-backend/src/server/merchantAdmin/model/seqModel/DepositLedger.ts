/* eslint-disable @typescript-eslint/consistent-type-assertions */
import { BuildOptions, DataTypes, Model, Sequelize } from 'sequelize';
import DepositLedger from '../dbModel/DepositLedger';
import * as SeqModel from './0_index';
export default DepositLedger;
export interface DepositLedgerModel
  extends Model<DepositLedger>,
    DepositLedger {}
export type DepositLedgerStatic = typeof Model & {
  new (values?: object, options?: BuildOptions): DepositLedgerModel;
};
export function DepositLedgerFactory(sequelize: Sequelize) {
  return <DepositLedgerStatic>sequelize.define(
    SeqModel.name.DepositLedgers,
    {
      txHash: {
        field: 'tx_hash',
        allowNull: false,
        primaryKey: true,
        type: DataTypes.STRING(255),
      },
      fromWalletAddr: {
        field: 'from_wallet_addr',
        allowNull: false,
        type: DataTypes.STRING(255),
      },
      fromAssetId: {
        field: 'from_asset_id',
        allowNull: false,
        type: DataTypes.STRING(20),
      },
      fromAssetAddr: {
        field: 'from_asset_addr',
        allowNull: false,
        type: DataTypes.STRING(255),
      },
      fromAssetAmt: {
        field: 'from_asset_amt',
        allowNull: false,
        type: DataTypes.DECIMAL(78),
      },
      toWalletAddr: {
        field: 'to_wallet_addr',
        allowNull: false,
        primaryKey: true,
        type: DataTypes.STRING(255),
      },
      toTokenAmt: {
        field: 'to_token_amt',
        allowNull: false,
        type: DataTypes.DECIMAL(78),
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
      userId: {
        type: DataTypes.BIGINT,
        allowNull: false,
        field: 'user_id',
      },
      remarks: {
        type: DataTypes.STRING(1000),
        allowNull: true,
        field: 'remarks',
      },
    },
    {
      freezeTableName: true,
      updatedAt: false,
      createdAt: false,
    },
  );
}
