import Web3 from 'web3';
import { EthClient } from '../utils/ethereum/EthClient';
import { AuthClient } from '../auth/AuthClient';

import { ServerConfigBase } from './ServerConfigBase';
import { Logger } from 'winston';
import Db from '../sequlelize';
import * as http from 'http';
import express = require('express');
export class ServerBase {
  configBase: ServerConfigBase;
  authClient?: AuthClient;
  logger?: Logger;
  server?: express.Express;
  httpServer?: http.Server;
  db?: Db;
  ethClient?: EthClient;

  public constructor(configBase: ServerConfigBase) {
    this.configBase = configBase;
  }

  public setupRouter() {}

  public listen() {
    if (!this.configBase.httpServerIp || !this.configBase.httpServerPort){
      throw new Error("httpServerIp or httpServerPort is undefined");
    }
    this.server!.listen(this.configBase.httpServerPort, this.configBase.httpServerIp , () => {
      this.printRoutes();
    });
  }

  public initHttpEthClient() {
    if (!this.configBase.ethClientProtocol){
      throw new Error("ethClientProtocol is undefined")
    }
    if (!this.configBase.ethClientPort){
      this.configBase.ethClientPort = 8545;
    }
    if (!this.configBase.ethClientProtocol){
      this.configBase.ethClientProtocol = "http";
    }
    if (!this.configBase.ethClientChainId){
      throw new Error("ethClientChainID is undefined");
    }
    let provider = new Web3.providers.HttpProvider(`${this.configBase.ethClientProtocol}://${this.configBase.ethClientIp}:${this.configBase.ethClientPort}`);
    this.ethClient = new EthClient(provider, this.configBase.ethClientChainId)
  }

  public printRoutes () {
    const routes: any[] = [];
    this.server!._router.stack.forEach((middleware) => {
      if (middleware.route) {
        // routes registered directly on the app
        routes.push(middleware.route);
      } else if (middleware.name === 'router') {
        // router middleware
        middleware.handle.stack.forEach((handler) => {
          const route = handler.route;
          if (route) {
            routes.push(route);
          }
        });
      }
    });

    routes.forEach((r) => {
      const methods: string[] = [];
      for (const key in r.methods) {
        methods.push(key.toUpperCase());
      }
      this.logger!.info(
        `${this.configBase.httpServerIp}:${this.configBase.httpServerPort}${r.path}  (${methods.join(',')})`,
      );
    });
  };

}
