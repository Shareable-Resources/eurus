/* eslint-disable @typescript-eslint/consistent-type-assertions */
import { BuildOptions, DataTypes, Model, Sequelize } from 'sequelize';
import HelloWorldUser from '../DBModel/HelloWorldUser';
import * as SeqModel from './0_index';
export default HelloWorldUser;

export interface HelloWorldUserModel
  extends Model<HelloWorldUser>,
    HelloWorldUser {}
export type HelloWorldUserStatic = typeof Model & {
  new (values?: object, options?: BuildOptions): HelloWorldUserModel;
};
export function HelloWorldUserFactory(sequelize: Sequelize) {
  return <HelloWorldUserStatic>sequelize.define(
    SeqModel.name.HelloWorldUser,
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
