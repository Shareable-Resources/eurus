/**
 * @type import('hardhat/config').HardhatUserConfig
 */

require('@nomiclabs/hardhat-ethers');

// Change private keys accordingly - ONLY FOR DEMOSTRATION PURPOSES - PLEASE STORE PRIVATE KEYS IN A SAFE PLACE
// Export your private key as
//       export PRIVKEY=0x.....
const privateKey = process.env.PRIVKEY;
const privateKeyDev =
   '0x5555fe01770a5c6dc621f00465e6a0f76bbd4bc1edbde5f2c380fcb00f354b99';

module.exports = {
   defaultNetwork: 'hardhat',

   networks: {
      hardhat: {},

      besudev: {
         url: 'http://13.228.169.25:8545',
         accounts: [privateKeyDev],
         chainId: 2021,
         gasPrice: 2400000000,
      },
      dev: {
         url: 'http://127.0.0.1:9933',
         accounts: [privateKeyDev],
         network_id: '1281',
         chainId: 1281,
      },
   },
   solidity: {
      compilers: [
         {
            version: '0.5.16',
            settings: {
               optimizer: {
                  enabled: true,
                  runs: 200,
               },
            },
         },
         {
            version: '0.6.6',
            settings: {
               optimizer: {
                  enabled: true,
                  runs: 200,
               },
            },
         },
      ],
   },
   paths: {
      sources: './contracts',
      cache: './cache',
      artifacts: './artifacts',
   },
   mocha: {
      timeout: 20000,
   },
};
