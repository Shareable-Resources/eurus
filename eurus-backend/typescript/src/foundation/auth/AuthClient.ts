import { ServerReturnCode } from '../server/ServerReturnCode';
import { ResponseMesssage } from './WebSocketMessage';
import { ResponseBase } from '../server/ApiMessage';
import * as websocket from 'websocket';
import * as auth_message from './AuthMessage';
import * as ws_message from './WebSocketMessage';
import * as crypto from '../utils/crypto/Crypto';
import { time } from 'console';

export interface IAuthConfig {
  getAuthIp(): string;
  getAuthPort(): number;
  getServiceId(): number;
  getPrivateKey(): string;
  getAuthPath(): string;
}

export class AuthClient {
  config?: IAuthConfig;
  isStop: boolean = false;
  sessionId: number = 0;
  loggedIn: boolean = false;
  loginToken: string = '';
  additionalInfo: string = '';
  wsClient?: websocket.client;
  conn?: websocket.connection;

  protected loginHandler?: (returnCode: ServerReturnCode) => void;

  constructor() {
    this.wsClient = new websocket.client();

    this.wsClient.on('connectFailed', (err: Error) => {
      console.log('Connect error: ' + err.message);
      this.onConnectionError(err);
    });

    this.wsClient.on('connect', (conn: websocket.connection) => {
      console.log('Connected');

      conn.on('error', (err: Error) => {
        this.onConnectionError(err);
      });
      conn.on('close', () => {
        this.onConnectionClosed();
      });
      conn.on('message', (message: websocket.Message) => {
        this.onMessageReceived(message);
      });
      this.conn = conn;
      this.onConnected();
    });
  }

  public isLoggedIn(): boolean {
    return this.loggedIn;
  }

  public loginAuthServer(
    config: IAuthConfig,
    loginHandler: (code: ServerReturnCode) => void,
  ) {
    this.config = config;
    this.loginHandler = loginHandler;
    this.login(config);
  }

  protected login(config: IAuthConfig) {
    let serverUrl: string =
      'ws://' +
      config.getAuthIp() +
      ':' +
      config.getAuthPort() +
      config.getAuthPath();

    this.wsClient?.connect(serverUrl, 'json', undefined);
  }

  protected onConnectionClosed() {
    console.log('Connection Closed');

    if (!this.isStop) {
      this.login(this.config!);
    }
  }

  protected onConnectionError(err: Error) {
    console.log('onConnection Error: ' + err.message);
    if (!this.isStop) {
      setTimeout(() => {
        this.login(this.config!);
      }, 1000);
    }
  }

  protected onConnected(): void {
    const content: auth_message.AuthenticateRequest =
      new auth_message.AuthenticateRequest();
    content.serviceId = this.config!.getServiceId();

    const req: ws_message.MasterRequestMessage =
      new ws_message.MasterRequestMessage(content);

    let reqStr: string = JSON.stringify(req);
    this.sendMessage(req);
  }

  protected sendMessage(
    masterMessage: ws_message.MasterRequestMessage,
  ): boolean {
    if (!this.conn) {
      return false;
    }
    let sign: string = crypto.Crypto.generateRSASignFromBase64(
      this.config!.getPrivateKey(),
      Buffer.from(masterMessage.message),
    );
    masterMessage.sign = sign;

    let reqStr: string = JSON.stringify(masterMessage);
    this.conn.sendUTF(reqStr, (err?: Error) => {
      if (err) {
        /*
        console.log(
          'Unable to send message to auth server: ',
          masterMessage.nonce,
        );*/
      }
    });
    return true;
  }

  protected onMessageReceived(message: websocket.Message) {
    if (message.type === 'utf8') {
      let content: string = message.utf8Data;
      let j = JSON.parse(content);
      let methodName: string = j['methodName'];
      switch (methodName) {
        case 'authenticate':
          let resContent = new auth_message.AuthenticateResponse();
          let res = new ResponseMesssage();
          res.data = resContent;

          Object.assign(res, j);

          this.loginToken = resContent.token;
          /*
          if (res.returnCode == 0) {
            this.loggedIn = true;
          }

          this.loginHandler!(res.returnCode);*/

          break;
        default:
          console.log('AuthClient receive unknown message: ', methodName);
      }
    }
  }
}
