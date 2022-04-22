import { EthClient } from '../../../foundation/utils/ethereum/0_index';

abstract class CommonObserver {
  protected abstract startMonitoringSmartContract(): void;
  protected abstract proccessBlock(event?: any): void;
  ethClient: EthClient;
  blocksRead: number[];
  constructor(ethClient: EthClient) {
    this.ethClient = ethClient;
    this.blocksRead = [];
  }

  getBlockNumber() {
    return new Promise((resolve) => {
      this.ethClient.web3Client.eth.getBlockNumber((blockNum) =>
        resolve(blockNum),
      );
    });
  }
}

export default CommonObserver;
