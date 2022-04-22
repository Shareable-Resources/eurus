import Web3 from 'web3';
import { provider } from 'web3-core';
import { HttpProviderOptions } from 'web3-core-helpers';

export class EthClient {
  web3Client: Web3;
  chainId: number;
  constructor(provider: provider, chainId: number, timeout?: number) {
    let timeoutVal: number = 30000;
    if (timeout) {
      timeoutVal = timeout;
    }
    let options: HttpProviderOptions = {
      keepAlive: false,
      timeout: timeoutVal,
    };

    this.web3Client = new Web3(provider);
    this.chainId = chainId;
  }
}
