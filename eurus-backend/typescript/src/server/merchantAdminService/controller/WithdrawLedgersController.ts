import {
  apiResponse,
  responseBySuccess,
  statusCode,
} from '../../../foundation/server/ApiMessage';
import CommonController from '../../../foundation/server/CommonController';
import Service from '../service/WithdrawLedgersService';
import * as DBModel from '../model/dbModel/0_index';
import { WithdrawLedgerStatus } from '../model/dbModel/WithdrawLedger';
import { WithdrawReqStatus } from '../model/dbModel/WithdrawReq';
export class Controller implements CommonController {
  /**
   * @description Creates an instance of withdraw ledgers controller.
   *
   * @constructor
   * @param {Service} service
   */
  constructor(private service: Service = new Service()) {
    this.getAll = this.getAll.bind(this);
    this.insertWithdrawLedger = this.insertWithdrawLedger.bind(this);
  }

  async getAll(req: any, res: any) {
    const list = await this.service.getAll();
    return apiResponse(
      res,
      responseBySuccess(list, true, 'query', 'success'),
      statusCode(true, 'query'),
    );
  }

  async insertWithdrawLedger(req: any, res: any) {
    const dateNow = new Date();
    const obj = new DBModel.WithdrawLedger();
    obj.txHash = 'yooooooo';
    obj.fromWalletAddr = 'yooooooo';
    obj.fromTokenAmt = BigInt(100000);
    obj.toWalletAddr = '0xabaf0503d06ac5d222653f3294a134db3ca98e29';
    obj.toAssetId = 'yooo';
    obj.toAssetAmt = BigInt(100000);
    obj.rate = BigInt(1);
    obj.reqId == BigInt(1);
    obj.createdDate = dateNow;
    obj.lastModifiedDate = dateNow;
    obj.status = WithdrawLedgerStatus.StatusCompleted;
    obj.remarks = `Block number : ${12314}`;
    const withdrawReq: DBModel.WithdrawReq = new DBModel.WithdrawReq();
    withdrawReq.lastModifiedDate = dateNow;
    withdrawReq.approveDate = dateNow;
    withdrawReq.txHash = 'yooooooo';
    withdrawReq.approveBy = BigInt(1);
    withdrawReq.id = BigInt(1);
    withdrawReq.status = WithdrawReqStatus.StatusApproved;
    const list = await this.service.insertWithdrawLedger(obj, withdrawReq);
    return apiResponse(
      res,
      responseBySuccess(list, true, 'query', 'success'),
      statusCode(true, 'query'),
    );
  }
}
export default Controller;
