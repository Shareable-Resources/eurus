let deployHelper = require('./deploy_helper');

const multiSigWallet = artifacts.require("GeneralMultiSigWallet");

let ownerAddressList = [
    '0x7fe4d0271f2988e7be4Aba6f825C191281f72c0B', 
    '0x1C6E9f40BB3035b59cd720e16B5E2BBB5a478Dea', 
    '0x8bC5E866be51C81ac75dCa282E118F634A36A33C'
    // '0x99B32dAD54F630D9ED36E193Bc582bbed273d666',
    // '0xf865eAbd5C55887e20572383Da8e03fAa140E731',
    // '0xC6e8AB9849aD2529cC0C023790Fa04616E2E1aDe'
];

module.exports = async (deployer)=>{
    deployHelper.readDeployLog();
    await deployHelper.deploy(deployer, multiSigWallet, ownerAddressList, 2);
    deployHelper.writeDeployLog();
}
