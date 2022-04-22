/*
import Web3 from 'web3';
import elliptic from 'elliptic';
const SignatureHelper = {
  //recover wallet address from pubKey
  //pubKey
  getSignatureVerifyResult(pubKey: string) {
    let signatureSignedByPrivateKey = getSignatureByInput(input);
    let pem = fs.readFileSync('PUBLIC_KEY_FILE_PATH_GOES_HERE');
    let publicKey = pem.toString('ascii');
    const verifier = crypto.createVerify('RSA-SHA256');
    verifier.update(input, 'ascii');
    const publicKeyBuf = new Buffer(publicKey, 'ascii');
    const signatureBuf = new Buffer(signatureSignedByPrivateKey, 'hex');
    const result = verifier.verify(publicKeyBuf, signatureBuf);

    return result;
  },
  const getSignatureByInput = (input) => {
    let privatePem = fs.readFileSync('PRIVATE_KEY_FILE_PATH_GOES_HERE')
    let key = privatePem.toString('ascii')
    let sign = crypto.createSign('RSA-SHA256')
    sign.update(input)
    let signature = sign.sign(key, 'hex')
  
    return signature
  }
};

export default SignatureHelper;
*/
