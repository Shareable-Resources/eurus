export class RequestBase {
  method: string = '';
  requestPath: string = '';
  nonce: string = '';
}
/*
protected onConnected(): void {
    const content: auth_message.AuthenticateRequest =
      new auth_message.AuthenticateRequest();
    content.serviceId = this.config!.getServiceId();

    const req: ws_message.MasterRequestMessage =
      new ws_message.MasterRequestMessage(content);

    let reqStr: string = JSON.stringify(req);
    this.sendMessage(req);
  }


protected sendMessage(
    masterMessage: ws_message.MasterRequestMessage,
  ): boolean {
    if (!this.conn) {
      return false;
    }
    let sign: string = crypto.Crypto.generateRSASignFromBase64(
      this.config!.getPrivateKey(),
      Buffer.from(masterMessage.message),
    );
    masterMessage.sign = sign;

    let reqStr: string = JSON.stringify(masterMessage);
    this.conn.sendUTF(reqStr, (err?: Error) => {
      if (err) {
        console.log(
          'Unable to send message to auth server: ',
          masterMessage.nonce,
        );
      }
    });
    return true;
  }*/
