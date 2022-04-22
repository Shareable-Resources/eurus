/* eslint-disable @typescript-eslint/consistent-type-assertions */
import { BuildOptions, DataTypes, Model, Sequelize } from 'sequelize';
import MerchantClient from '../dbModel/MerchantClient';
import * as SeqModel from './0_index';
export default MerchantClient;
export interface MerchantClientModel
  extends Model<MerchantClient>,
    MerchantClient {}
export type MerchantClientStatic = typeof Model & {
  new (values?: object, options?: BuildOptions): MerchantClientModel;
};
export function MerchantClientFactory(sequelize: Sequelize) {
  return <MerchantClientStatic>sequelize.define(
    SeqModel.name.MerchantClient,
    {
      id: {
        field: 'id',
        allowNull: false,
        primaryKey: true,
        type: DataTypes.BIGINT,
        autoIncrement: true,
      },
      username: {
        field: 'username',
        allowNull: false,
        type: DataTypes.STRING(255),
      },
      balance: {
        field: 'balance',
        allowNull: false,
        defaultValue: 0,
        type: DataTypes.DECIMAL(75),
      },
      walletAddress: {
        field: 'wallet_address',
        allowNull: false,
        type: DataTypes.STRING(50),
      },
      createdDate: {
        type: DataTypes.DATE,
        defaultValue: DataTypes.NOW,
        field: 'created_date',
      },
      lastModifiedDate: {
        field: 'last_modified_date',
        defaultValue: DataTypes.NOW,
        type: DataTypes.DATE,
      },
      status: {
        field: 'status',
        allowNull: false,
        type: DataTypes.SMALLINT,
      },
    },
    {
      freezeTableName: true,
      createdAt: false,
      updatedAt: false,
      indexes: [
        {
          unique: true,
          fields: ['username'],
        },
      ],
    },
  );
}
