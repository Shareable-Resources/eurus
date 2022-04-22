import {
  apiResponse,
  responseBySuccess,
  statusCode,
} from '../../../foundation/server/ApiMessage';
import CommonController from '../../../foundation/server/CommonController';
import logger from '../util/ServiceLogger';
import Service from '../service/HelloWorldService';
import * as DBModel from '../model/dbModel/0_index';
export class Controller implements CommonController {
  /**
   * @description Creates an instance of hello world users controller.
   *
   * @constructor
   * @param {Service} service
   */
  constructor(private service: Service = new Service()) {
    this.getAll = this.getAll.bind(this);
    this.getAllByRaw = this.getAllByRaw.bind(this);
    this.getById = this.getById.bind(this);
    this.updateUser = this.updateUser.bind(this);
    this.removeUser = this.removeUser.bind(this);
    this.hello = this.hello.bind(this);
    this.create = this.create.bind(this);
  }

  async hello(req: any, res: any) {
    logger.info('hello world');
    return apiResponse(
      res,
      responseBySuccess(
        'Hello World',
        true,
        'query',
        'success',
        'Response returned from controller',
      ),
      statusCode(true, 'query'),
    );
  }

  async getAll(req: any, res: any) {
    const data = await this.service.getAll();
    return apiResponse(
      res,
      responseBySuccess(data, true, 'query', 'success', 'Founded'),
      statusCode(true, 'query'),
    );
  }
  async getAllByRaw(req: any, res: any) {
    const data = await this.service.getAllByRaw();
    return apiResponse(
      res,
      responseBySuccess(data, true, 'query', 'success', 'Founded'),
      statusCode(true, 'query'),
    );
  }
  //Using search params -    localhost:8092/helloWorldUsers/1 (GET)
  async getById(req: any, res: any) {
    const id = req.params.id;
    const resp = await this.service.getById(req.params);
    return apiResponse(res, resp, statusCode(resp.success, 'query'));
  }
  //Using request body -   localhost:8092/helloWorldUsers/updateUser (POST) {"id":1}
  async updateUser(req: any, res: any) {
    const reqBody: DBModel.HelloWorldUser = req.body;
    const resp = await this.service.updateUser(reqBody);
    return apiResponse(res, resp, statusCode(resp.success, 'up'));
  }

  async removeUser(req: any, res: any) {
    const reqBody: DBModel.HelloWorldUser = req.body;
    const resp = await this.service.removeUser(reqBody);
    return apiResponse(res, resp, statusCode(resp.success, 'del'));
  }

  async create(req: any, res: any) {
    const reqBody: DBModel.HelloWorldUser = req.body;
    const resp = await this.service.create(reqBody);
    return apiResponse(res, resp, statusCode(resp.success, 'add'));
  }
}
export default Controller;
