import { Response } from 'express';
import { StatusCodes } from 'http-status-codes';
import { ServerReturnCode } from './ServerReturnCode';

export class ResponseBase {
  success: boolean;
  data: any;
  respType: ResponseType;
  returnCode: ServerReturnCode;
  msg: string;
  constructor() {
    this.success = false;
    this.data = undefined;
    this.respType = 'info';
    this.returnCode = ServerReturnCode.Success;
    this.msg = '';
  }
}

export class WarningResponseBase extends ResponseBase {
  constructor(returnCode?: ServerReturnCode, msg?: string) {
    super();
    this.respType = 'warning';
    this.returnCode = returnCode
      ? returnCode
      : ServerReturnCode.InternalServerError;
    this.msg = msg ? msg : 'Warning';
  }
}

export class ErrorResponseBase extends ResponseBase {
  constructor(returnCode?: ServerReturnCode, msg?: string) {
    super();
    this.respType = 'error';
    this.returnCode = returnCode
      ? returnCode
      : ServerReturnCode.InternalServerError;
    this.msg = msg ? msg : 'Error';
  }
}

export interface ApiResponse {
  <T>(res: Response, data: T, statusCode: number): Response;
}

export const apiResponse: ApiResponse = (res, data, statusCode): Response => {
  return res.format({
    json: () => {
      res.type('application/json');
      res.status(statusCode).send(data);
    },
  });
};

export function successResponse(data: any, msg?: string): ResponseBase {
  const resData = new ResponseBase();
  resData.success = true;
  resData.data = data;
  resData.msg = msg ? msg : 'Success';
  return resData;
}

export function failedResponse(data: any, msg?: string): ResponseBase {
  const resData = new ResponseBase();
  resData.success = false;
  resData.data = data;
  resData.msg = msg ? msg : 'Fail';
  return resData;
}

export function statusCode(success: boolean, type: APIType) {
  switch (type) {
    case 'up':
      return StatusCodes.OK;
    case 'del':
      return StatusCodes.OK;
    case 'add':
      return StatusCodes.OK;
    case 'query':
      return StatusCodes.OK;
    case 'error':
      return StatusCodes.INTERNAL_SERVER_ERROR;
  }
}

export function responseBySuccess(
  data: any,
  success: boolean,
  type: APIType,
  respType?: ResponseType,
  msg?: string,
  returnCode?: ServerReturnCode,
): ResponseBase {
  const resData = new ResponseBase();
  resData.success = success ? true : false;
  resData.data = data;
  resData.respType = respType ? respType : 'info';
  if (success) {
    resData.returnCode = ServerReturnCode.Success;
    switch (type) {
      case 'up':
        resData.msg = msg ? msg : 'Update Success';
        break;
      case 'del':
        resData.msg = msg ? msg : 'Delete Success';
        break;
      case 'add':
        resData.msg = msg ? msg : 'Create Success';
        break;
      case 'query':
        resData.msg = msg ? msg : 'Found';
        break;
      case 'error':
        resData.msg = msg ? msg : 'Error';
        break;
    }
  } else {
    switch (type) {
      case 'up':
        resData.returnCode = ServerReturnCode.DatabaseError;
        resData.msg = msg ? msg : 'Update Fail';
        break;
      case 'del':
        resData.returnCode = ServerReturnCode.DatabaseError;
        resData.msg = msg ? msg : 'Delete Fail';
        break;
      case 'add':
        resData.returnCode = ServerReturnCode.DatabaseError;
        resData.msg = msg ? msg : 'Create Fail';
        break;
      case 'query':
        resData.returnCode = ServerReturnCode.RecordNotFound;
        resData.msg = msg ? msg : 'Found';
        break;
      case 'error':
        resData.returnCode = ServerReturnCode.InternalServerError;
        resData.msg = msg ? msg : 'Error';
        break;
    }
  }
  if (returnCode) {
    resData.returnCode = returnCode;
  }
  return resData;
}

export type APIType = 'up' | 'del' | 'add' | 'query' | 'error';
export type ResponseType = 'warning' | 'info' | 'success' | 'error';
