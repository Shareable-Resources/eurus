export enum ChainId {
  DevChain = 2021, //Development side chain
  Mainnet = 1001,
  Testnet = 1984,
  Rinkeby = 4,
}

export enum EthClientProviderUrl {
  webSocket = 'ws://54.254.124.206:8546',
  http = 'http://54.254.124.206:8545',
}

export default { ChainId, EthClientProviderUrl };
