const HDWalletProvider = require("@truffle/hdwallet-provider");
const PrivateKeyDevelopment_PrivateChain = "906b9fa7fab30a174fae300b9580906738b68b28afed1a47aaf0a6d7e420cb9d"//"5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99"; //"8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63";



const PrivateKeyProviderDevelopment_Rinkeby = new HDWalletProvider(
  {
    privateKeys: [PrivateKeyDevelopment_PrivateChain],
    providerOrUrl: "http://18.142.45.37:8545",
    pollingInterval: 80000
  }
);


const PrivateKeyTestnet_Rinkeby = "71d870d4cfe9d39c664bb99a4ee3afb272cc8af0844dabe40a02961fca326e10";
const PrivateKeyProviderTestnet_Rinkeby = new HDWalletProvider({
privateKeys: [PrivateKeyTestnet_Rinkeby],
providerOrUrl:"http://18.142.45.37:8545",
pollingInterval: 70000000
});
/**
 * Use this file to configure your truffle project. It's seeded with some
 * common settings for different networks and features like migrations,
 * compilation and testing. Uncomment the ones you need or modify
 * them to suit your project as necessary.
 *
 * More information about configuration can be found at:
 *
 * trufflesuite.com/docs/advanced/configuration
 *
 * To deploy via Infura you'll need a wallet provider (like @truffle/hdwallet-provider)
 * to sign your transactions before they're sent to a remote public node. Infura accounts
 * are available for free at: infura.io/register.
 *
 * You'll also need a mnemonic - the twelve word phrase the wallet uses to generate
 * public/private key pairs. If you're publishing your code to GitHub make sure you load this
 * phrase from a file you've .gitignored so it doesn't accidentally become public.
 *
 */

// const HDWalletProvider = require('@truffle/hdwallet-provider');
// const infuraKey = "fj4jll3k.....";
//
// const fs = require('fs');
// const mnemonic = fs.readFileSync(".secret").toString().trim();

module.exports = {
  /**
   * Networks define how you connect to your ethereum client and let you set the
   * defaults web3 uses to send transactions. If you don't specify one truffle
   * will spin up a development blockchain for you on port 9545 when you
   * run `develop` or `test`. You can ask a truffle command to use a specific
   * network from the command line, e.g
   *
   * $ truffle test --network <network-name>
   */

  networks: {
    dev_rinkeby:{
      provider: PrivateKeyProviderDevelopment_Rinkeby,
      network_id: "4",       // Ropsten's id
      gas: 10000000,        // Ropsten has a lower block limit than mainnet
      confirmations: 1,    // # of confs to wait between deployments. (default: 0)
      timeoutBlocks: 2000,  // # of blocks before a deployment times out  (minimum/default: 50)
      skipDryRun: true,    // Skip dry run before migrations? (default: false for public nets )
      gasPrice: 10000000000,
    },
    testnet_rinkeby:{
      provider: PrivateKeyProviderTestnet_Rinkeby,
      network_id: "4",       // Ropsten's id
      gas: 10000000,        // Ropsten has a lower block limit than mainnet
      confirmations: 1,    // # of confs to wait between deployments. (default: 0)
      timeoutBlocks: 2000,  // # of blocks before a deployment times out  (minimum/default: 50)
      skipDryRun: true,    // Skip dry run before migrations? (default: false for public nets )
      gasPrice: 10000000000,
    },
    // Useful for testing. The `development` name is special - truffle uses it by default
    // if it's defined here and no other network is specified at the command line.
    // You should run a client (like ganache-cli, geth or parity) in a separate terminal
    // tab if you use this network and you must also set the `host`, `port` and `network_id`
    // options below to some value.
    //
    // development: {
    //  host: "127.0.0.1",     // Localhost (default: none)
    //  port: 8545,            // Standard Ethereum port (default: none)
    //  network_id: "*",       // Any network (default: none)
    // },
    // Another network with more advanced options...
    // advanced: {
    // port: 8777,             // Custom port
    // network_id: 1342,       // Custom network
    // gas: 8500000,           // Gas sent with each transaction (default: ~6700000)
    // gasPrice: 20000000000,  // 20 gwei (in wei) (default: 100 gwei)
    // from: <address>,        // Account to send txs from (default: accounts[0])
    // websockets: true        // Enable EventEmitter interface for web3 (default: false)
    // },
    // Useful for deploying to a public network.
    // NB: It's important to wrap the provider as a function.
    // ropsten: {
    // provider: () => new HDWalletProvider(mnemonic, `https://ropsten.infura.io/v3/YOUR-PROJECT-ID`),
    // network_id: 3,       // Ropsten's id
    // gas: 5500000,        // Ropsten has a lower block limit than mainnet
    // confirmations: 2,    // # of confs to wait between deployments. (default: 0)
    // timeoutBlocks: 200,  // # of blocks before a deployment times out  (minimum/default: 50)
    // skipDryRun: true     // Skip dry run before migrations? (default: false for public nets )
    // },
    // Useful for private networks
    // private: {
    // provider: () => new HDWalletProvider(mnemonic, `https://network.io`),
    // network_id: 2111,   // This network is yours, in the cloud.
    // production: true    // Treats this network as if it was a public net. (default: false)
    // }
  },

  // Set default mocha options here, use special reporters etc.
  mocha: {
    // timeout: 100000
  },
  // Configure your compilers
  compilers: {
    solc: {
      version: "0.6.2",    // Fetch exact version from solc-bin (default: truffle's version)
      // docker: true,        // Use "0.5.1" you've installed locally with docker (default: false)
      settings: {          // See the solidity docs for advice about optimization and evmVersion
        optimizer: {
          enabled: true,
          runs: 200
        },
        //  evmVersion: "byzantium"
      }
    }
  }
};