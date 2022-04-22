/* eslint-disable @typescript-eslint/consistent-type-assertions */
import { BuildOptions, DataTypes, Model, Sequelize } from 'sequelize';
import WithdrawReq from '../dbModel/WithdrawReq';
import * as SeqModel from './0_index';
export default WithdrawReq;
export interface WithdrawReqModel extends Model<WithdrawReq>, WithdrawReq {}
export type WithdrawReqStatic = typeof Model & {
  new (values?: object, options?: BuildOptions): WithdrawReqModel;
};
export function WithdrawReqFactory(sequelize: Sequelize) {
  return <WithdrawReqStatic>sequelize.define(
    SeqModel.name.WithdrawReqs,
    {
      id: {
        field: 'id',
        allowNull: false,
        primaryKey: true,
        type: DataTypes.BIGINT,
        autoIncrement: true,
      },
      fromWalletAddr: {
        field: 'from_wallet_addr',
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
        field: 'to_asset_addr',
        allowNull: false,
        type: DataTypes.STRING(255),
      },
      rate: {
        field: 'rate',
        allowNull: false,
        type: DataTypes.DECIMAL(78),
      },
      status: {
        field: 'status',
        allowNull: false,
        type: DataTypes.SMALLINT,
      },
      rejectReason: {
        field: 'reject_reason',
        allowNull: true,
        type: DataTypes.STRING(1000),
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
      approveBy: {
        type: DataTypes.BIGINT,
        allowNull: true,
        field: 'approve_by',
      },
      approveDate: {
        type: DataTypes.DATE,
        allowNull: true,
        field: 'approve_date',
      },
      userId: {
        type: DataTypes.BIGINT,
        allowNull: false,
        field: 'user_id',
      },
      txHash: {
        type: DataTypes.STRING,
        allowNull: true,
        field: 'tx_hash',
      },
    },
    {
      freezeTableName: true,
      updatedAt: false,
      createdAt: false,
    },
  );
}
