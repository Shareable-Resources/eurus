import { ResponseBase } from '../server/ApiMessage';
import * as uuid from 'uuid';

export interface IRequestConent {
  getMethodName(): string;
}

export class RequestBase {
  method: string = '';
  requestPath: string = '';
  nonce: string = '';
}

export class RequestMessage extends RequestBase {
  methodName: string = '';
  timestamp: number = 0;
  data?: any;
}

export class MasterRequestMessage extends RequestMessage {
  message: string = '';
  sign: string = '';

  public constructor(content?: IRequestConent, methodName?: string) {
    super();
    this.timestamp = Date.now();
    this.nonce = uuid.v4();

    let tempReq: RequestMessage = new RequestMessage();
    tempReq.data = content;
    if (methodName) {
      tempReq.methodName = methodName;
    } else {
      tempReq.methodName = content ? content.getMethodName() : '';
    }
    tempReq.timestamp = this.timestamp;

    let messageStr: string = JSON.stringify(tempReq);
    this.message = messageStr;

    this.methodName = tempReq.methodName;
    this.data = content;
  }
}

export class ResponseMesssage extends ResponseBase {
  methodName: string = '';
  timestamp: number = 0;
}
