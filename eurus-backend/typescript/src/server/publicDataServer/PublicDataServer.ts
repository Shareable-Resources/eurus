import { ServerBase } from '../../foundation/server/ServerBase';
import express from 'express';

import {
  timeMiddleware,
  accessControlAllowMiddleware,
  clientErrorHandler,
} from '../../foundation/server/Middlewares';

import PublicDataRoute from './route/PublicDataRoute';
import { ServerConfigBase } from '../../foundation/server/ServerConfigBase';
import winston from 'winston';

export interface PublicDataServerConfig extends ServerConfigBase {
  initializeEunAmount : number;
}

export class PublicDataServer extends ServerBase {

  config : PublicDataServerConfig;

  constructor (config : PublicDataServerConfig, logger: winston.Logger){
    super(config);
    this.config = config;
    this.server = express();
    let app = this.server;
    app.use(accessControlAllowMiddleware);
    //app.use(jwtVertificaitonMiddleWare);
    app.use(express.urlencoded({ extended: true }));
    // app.use(express.json());
    app.use(timeMiddleware);
    app.use(clientErrorHandler);

    // We provide a root route just as an example
    app.get('/', (req: any, res: any) => {
      res.send(
        `<p>Eurus Public Data Server Starts</p>`,
      );
    });


    this.logger = logger;

  }

  public override setupRouter(){
    let route = new PublicDataRoute(this);
    this.server!.use(route.getRouter());
  }


}

export default PublicDataServer;
