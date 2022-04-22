import * as cr from "crypto";

export class Crypto {

  private static aesBlockSize : number = 16;


  public static generateRSASignFromBase64(
    privateKey: string,
    message: Buffer
  ): string {
    let pkcs1FormateKey: string =
      "-----BEGIN RSA PRIVATE KEY-----\n" +
      privateKey +
      "\n-----END RSA PRIVATE KEY-----";

    let privKeyObj: cr.KeyObject = cr.createPrivateKey({
      key: pkcs1FormateKey,
      type: "pkcs1",
    });

    let signature : Buffer = cr.sign("sha256", message, {
      key: privKeyObj,
      padding: cr.constants.RSA_PKCS1_PSS_PADDING,
    });

    let signBase64: string = signature.toString("base64");

    return signBase64;
  }



  public static decryptAESFromBase64(base64Data: string, base64Key: string): Buffer {
      let aesKeyBuff : Buffer = Buffer.from(base64Key, 'base64');
      // let aesKey : cr.KeyObject = cr.createSecretKey(Uint8Array.from(aesKeyBase64));
      let cipher : Buffer = Buffer.from(base64Data, 'base64')

      let iv : Buffer = cipher.slice(0, Crypto.aesBlockSize)
      let data : Buffer = cipher.slice(Crypto.aesBlockSize)

      let decipher = cr.createDecipheriv("aes-256-cbc", Uint8Array.from(aesKeyBuff), iv)
      decipher.update(data);
      return decipher.final();
  }
}
