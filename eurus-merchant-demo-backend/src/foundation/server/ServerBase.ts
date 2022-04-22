import { EthClient } from '../utils/ethereum/EthClient';
import express from 'express';
import { Logger } from 'winston';
import Db from '../sequlelize';
import * as http from 'http';
export class ServerBase {
  logger?: Logger;
  server?: express.Express;
  httpServer?: http.Server;
  db?: Db;
  ethClient?: EthClient;

  public constructor() {}

  protected setupRouter() {}
}
