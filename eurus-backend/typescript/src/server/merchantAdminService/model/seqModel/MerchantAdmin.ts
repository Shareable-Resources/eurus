/* eslint-disable @typescript-eslint/consistent-type-assertions */
import { BuildOptions, DataTypes, Model, Sequelize } from 'sequelize';
import MerchantAdmin from '../dbModel/MerchantAdmin';
import * as SeqModel from './0_index';
export default MerchantAdmin;
export interface MerchantAdminModel
  extends Model<MerchantAdmin>,
    MerchantAdmin {}
export type MerchantAdminStatic = typeof Model & {
  new (values?: object, options?: BuildOptions): MerchantAdminModel;
};
export function MerchantAdminFactory(sequelize: Sequelize) {
  return <MerchantAdminStatic>sequelize.define(
    SeqModel.name.MerchantAdmins,
    {
      operatorId: {
        field: 'operator_id',
        allowNull: false,
        primaryKey: true,
        type: DataTypes.INTEGER,
        autoIncrement: true,
      },
      merchantId: {
        field: 'merchant_id',
        allowNull: false,
        type: DataTypes.BIGINT,
      },
      username: {
        field: 'username',
        allowNull: false,
        type: DataTypes.STRING(255),
      },
      email: {
        field: 'email',
        allowNull: false,
        type: DataTypes.STRING(255),
      },
      passwordHash: {
        field: 'password_hash',
        allowNull: false,
        type: DataTypes.STRING(300),
      },
      status: {
        field: 'status',
        allowNull: false,
        type: DataTypes.SMALLINT,
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
    },
    {
      freezeTableName: true,
      createdAt: false,
      updatedAt: false,
      indexes: [
        {
          unique: true,
          fields: ['merchant_id', 'username'],
        },
      ],
    },
  );
}
