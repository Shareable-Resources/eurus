import * as ws_message from './WebSocketMessage';

export class AuthenticateRequest implements ws_message.IRequestConent {
  serviceId: number = 0;

  public getMethodName(): string {
    return 'authenticate';
  }
}

export class AuthenticateResponse {
  sessionId: number = 0;
  token: string = '';
  additionalInfo: string = '';
}
