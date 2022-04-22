import { Request, Response } from 'express';

import {
  apiResponse,
  ResponseBase,
  responseBySuccess,
  statusCode,
} from '../../../foundation/server/ApiMessage';
import CommonController from '../../../foundation/server/CommonController';
import Service from '../service/MerchantClientsService';
import * as ReqModel from '../model/reqModel/0_index';
import { EthAccount } from '../../../foundation/utils/ethereum/0_index';
export class Controller implements CommonController {
  /**
   * @description Creates an instance of users controller.
   *
   * @constructor
   * @param {Service} service
   */
  ethAccount: EthAccount;
  constructor(private service: Service = new Service()) {
    this.ethAccount = new EthAccount();
    this.importWallet = this.importWallet.bind(this);
    this.getAll = this.getAll.bind(this);
    this.getById = this.getById.bind(this);
  }

  async getById(req: any, res: any) {
    const obj = await this.service.getById(req.params.id);
    const respType = obj ? 'success' : 'info';
    const msg = obj ? 'Found' : 'Not Found';
    return apiResponse(
      res,
      responseBySuccess(obj, true, 'query', respType, msg),
      statusCode(true, 'query'),
    );
  }
  async getAll(req: any, res: any) {
    const list = await this.service.getAll();
    return apiResponse(
      res,
      responseBySuccess(list, true, 'query', 'success'),
      statusCode(true, 'query'),
    );
  }

  async importWallet(req: any, res: any) {
    const reqBody: ReqModel.ImportWalletReq = req.body;
    // Verify the signed data is from actual user
    const addressFromSignature = this.ethAccount.recover(
      reqBody.message,
      reqBody.signature,
    );
    const isVerified =
      addressFromSignature.toUpperCase() == reqBody.walletAddress.toUpperCase();
    let responseBase: ResponseBase = new ResponseBase();
    if (isVerified) {
      responseBase = await this.service.importWallet(reqBody);
    } else {
      responseBase.msg = 'Signature mismatched';
      responseBase.success = false;
      responseBase.respType = 'warning';
    }

    return apiResponse(
      res,
      responseBySuccess(
        responseBase.data,
        responseBase.success,
        'query',
        responseBase.respType,
        responseBase.msg,
      ),
      statusCode(responseBase.success, 'query'),
    );
  }
}
export default Controller;
