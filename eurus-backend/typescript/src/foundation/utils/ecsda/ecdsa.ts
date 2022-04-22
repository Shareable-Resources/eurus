import Web3 from 'web3';
import elliptic from 'elliptic';
const ECSDAHelper = {
  //recover wallet address from pubKey
  //pubKey
  walletAddrFromPubKey(pubKey: string) {
    const ec = new elliptic.ec('secp256k1');
    const ecKey = ec.keyFromPublic(pubKey, 'hex');
    const pubKeyHex = `0x${ecKey.getPublic(false, 'hex').slice(2)}`;
    const pubHash = Web3.utils.sha3(pubKeyHex);
    const address = Web3.utils.toChecksumAddress(`0x${pubHash?.slice(-40)}`);
    return address;
  },
};

export default ECSDAHelper;
