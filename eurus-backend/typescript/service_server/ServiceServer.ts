import { AuthClient } from './../foundation/src/auth/AuthClient';
import { ServerReturnCode } from '../foundation/src/ServerReturnCode';
import { ServerBase } from './../foundation/src/server/ServerBase';
import { ServerConfigBase } from '../foundation/src/server/ServerConfigBase';
import { Server } from 'http';

export class ServiceServer extends ServerBase{
    
    loginHandler? : (ServerReturnCode) => void;
    dummy : number = 0;

    constructor(config: ServerConfigBase){
        super(config);
    }

    public startup(configPath: string, loginHandler?: (ServerReturnCode) => void) : boolean{
        let isSuccess : boolean = false;
        isSuccess = this.loadConfig(configPath);
        if (!isSuccess){
            return isSuccess;
        }
        return this.initAuth(loginHandler);
    }

    public initAuth(loginHandler?: (ServerReturnCode) => void) : boolean{
      
        this.loginHandler = loginHandler;

        let authClientConcrete : AuthClient = new AuthClient();
        let _this: ServiceServer = this;
        
        return super.loginAuth(authClientConcrete, (returnCode :ServerReturnCode)=> {
            if (_this.loginHandler) {
                _this.loginHandler(returnCode);
            }else{
                _this.defaultPostLoginProcessing(returnCode);
            }
        });
    }

    public async defaultPostLoginProcessing(returnCode: ServerReturnCode){
        if (returnCode == ServerReturnCode.Success){
            let isSuccess : boolean = await this.initDatabase();
            if (!isSuccess){
                return;
            }
            this.initHttpServer();
        }
    }

    protected processPostLogin (returnCode: ServerReturnCode){
        
    }
}
