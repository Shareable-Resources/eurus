import seq from '../sequelize';
import CommonService from '../../../foundation/server/CommonService';
import * as DBModel from '../model/dbModel/0_index';
import * as SeqModel from '../model/seqModel/0_index';
import { Transaction } from 'sequelize/types';
import {
  ErrorResponseBase,
  ResponseBase,
  WarningResponseBase,
} from '../../../foundation/server/ApiMessage';
import logger from '../util/ServiceLogger';
import { ServerReturnCode } from '../../../foundation/server/ServerReturnCode';
import { HelloWorldUserStatus } from '../model/DBModel/HelloWorldUser';
const modelModule = seq.sequelize.models;
export default class Service implements CommonService {
  async getAll(searchParams?: any): Promise<DBModel.HelloWorldUser[]> {
    const results: any = await modelModule[
      SeqModel.name.HelloWorldUser
    ].findAll({});
    return results as DBModel.HelloWorldUser[];
  }
  async getAllByRaw(): Promise<any> {
    const [results, metadata]: any = await seq.sequelize.query(
      `select id as id, last_modified_date as "lastModifiedDate" from "hello_world"."hello_world_users" hwu`,
    );
    return results;
  }
  async getById(searchParams?: any): Promise<ResponseBase> {
    const funcMsg = `[helloWorldService][getById](obj.id : ${searchParams.id})`;
    let resp = new ResponseBase();
    const result = await modelModule[SeqModel.name.HelloWorldUser].findOne({
      where: {
        id: searchParams.id,
      },
    });
    if (!result) {
      resp = new WarningResponseBase(
        ServerReturnCode.RecordNotFound,
        'Record not found',
      );
      logger.info(funcMsg + ' - fail ', { message: resp.msg });
    } else {
      resp.data = result;
      resp.success = true;
      resp.msg = 'Record found';
    }
    return resp;
  }
  async updateUser(newObj: DBModel.HelloWorldUser): Promise<ResponseBase> {
    const funcMsg = `[helloWorldService][update](obj.id : ${newObj.id})`;
    const t = await seq.sequelize.transaction();
    let resp = new ResponseBase();
    const dateNow = new Date();
    newObj.createdDate = dateNow;
    newObj.lastModifiedDate = dateNow;
    try {
      //Validate the record is in db
      const recordInDb = await modelModule[
        SeqModel.name.HelloWorldUser
      ].findOne({
        where: {
          id: newObj.id,
        },
        transaction: t,
      });
      if (!recordInDb) {
        resp = new WarningResponseBase(
          ServerReturnCode.RecordNotFound,
          'Record not found',
        );
        logger.info(funcMsg + ' - fail ', { message: resp.msg });
      } else {
        const updateResult: any = await modelModule[
          SeqModel.name.HelloWorldUser
        ].update(newObj, {
          transaction: t,
          where: {
            id: newObj.id,
          },
          fields: ['balance', 'lastModifiedDate'], //fields will limit the columns that need to be updated, use the DBModel attributes name instead of DB column name
        });
        const affectedRowsMsg = `Updated affected rows (${updateResult[0]})`;
        logger.info(funcMsg, {
          message: affectedRowsMsg,
        });
        resp.success = updateResult[0] > 0 ? true : false;
        resp.msg = affectedRowsMsg;
        resp.returnCode =
          updateResult[0] > 0
            ? ServerReturnCode.Success
            : ServerReturnCode.InternalServerError;
        resp.respType = updateResult[0] > 0 ? 'success' : 'warning';
      }
      await t.commit(); //Transaction must be commited or rollback, otherwise sequlize connection to DB will be locked
    } catch (e) {
      logger.error(funcMsg + ' - fail ');
      logger.error(e);
      resp = new ErrorResponseBase(
        ServerReturnCode.InternalServerError,
        'Update hello world user fail :' + e,
      );
      await t.rollback();
    }
    return resp;
  }
  async create(obj: DBModel.HelloWorldUser): Promise<ResponseBase> {
    const newObj = obj;
    newObj.id = null;
    const dateNow = new Date();
    newObj.createdDate = dateNow;
    newObj.lastModifiedDate = dateNow;
    const funcMsg = `[helloWorldService][create](obj.username : ${newObj.username})`;
    const t = await seq.sequelize.transaction();
    let resp = new ResponseBase();

    try {
      const recordInDb = await modelModule[
        SeqModel.name.HelloWorldUser
      ].findOne({
        where: {
          username: newObj.username,
        },
        transaction: t,
      });
      if (recordInDb) {
        resp = new WarningResponseBase(
          ServerReturnCode.UniqueViolationError,
          'Duplicated Record',
        );
        logger.info(funcMsg + ' - fail ', { message: resp.msg });
      } else {
        const insertResult: any = await modelModule[
          SeqModel.name.HelloWorldUser
        ].create(newObj, {
          transaction: t,
        });
        resp.msg = `New hello world user is created. New record id :(${insertResult.id})`;
        logger.info(funcMsg + ' - success ', { message: resp.msg });
      }
      await t.commit(); //Transaction must be commited or rollback, otherwise sequlize connection to DB will be locked
    } catch (e) {
      logger.error(funcMsg + ' - fail ');
      logger.error(e);
      resp = new ErrorResponseBase(
        ServerReturnCode.InternalServerError,
        'Insert hello world user fail :' + e,
      );
      await t.rollback();
    }
    return resp;
  }
  async removeUser(obj?: any): Promise<ResponseBase> {
    const funcMsg = `[helloWorldService][remove](obj.id : ${obj.id})`;
    const t = await seq.sequelize.transaction();
    let resp = new ResponseBase();
    try {
      //Validate the record is in db
      const recordInDb = await modelModule[
        SeqModel.name.HelloWorldUser
      ].findOne({
        where: {
          id: obj.id,
        },
        transaction: t,
      });
      if (!recordInDb) {
        resp = new WarningResponseBase(
          ServerReturnCode.RecordNotFound,
          'Record has already been deleted',
        );
        logger.info(funcMsg + ' - fail ', { message: resp.msg });
      } else {
        const deleteRows: any = await modelModule[
          SeqModel.name.HelloWorldUser
        ].destroy({
          transaction: t,
          where: {
            id: obj.id,
          },
        });
        const affectedRowsMsg = `Deleted affected rows (${deleteRows})`;
        console.log(deleteRows);
        resp.success = deleteRows ? true : false;
        resp.msg = affectedRowsMsg;
        resp.returnCode =
          deleteRows > 0
            ? ServerReturnCode.Success
            : ServerReturnCode.InternalServerError;
        resp.respType = deleteRows > 0 ? 'success' : 'warning';
      }
      await t.commit(); //Transaction must be commited or rollback, otherwise sequlize connection to DB will be locked
    } catch (e) {
      logger.error(funcMsg + ' - fail ');
      logger.error(e);
      resp = new ErrorResponseBase(
        ServerReturnCode.InternalServerError,
        'Delete hello world user fail :' + e,
      );
      await t.rollback();
    }
    return resp;
  }
}
