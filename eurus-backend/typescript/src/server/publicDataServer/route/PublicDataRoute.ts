import { PublicDataController } from '../controller/PublicDataController';
import * as controller from '../controller/0_index';
import { Router } from 'express';

import CommonRoute from '../../../foundation/server/CommonRoute';
import { makeHandlerAwareOfAsyncErrors } from '../../../foundation/server/Middlewares';
import winston from 'winston';
import { ServerBase } from '../../../foundation/server/ServerBase';
import PublicDataServer from '../PublicDataServer';


export class PublicDataRoute extends CommonRoute {
    
    publicDataController : controller.PublicDataController
    server : ServerBase
    constructor(server : PublicDataServer) {
        super(server.logger!);
        this.prefix = '';
        this.server = server;
        this.publicDataController = new PublicDataController(server);
        this.setRoutes();

      }

    protected setRoutes() {   
        this.router.get(
            `/eunAvail`,
            makeHandlerAwareOfAsyncErrors( this.publicDataController.getEUNAvailability, this.logger),
        );
    }

}


  export default PublicDataRoute;
  
