export default class ImportWalletReq {
  username: string;
  signature: string;
  walletAddress: string;
  message: string;
  constructor() {
    this.username = '';
    this.signature = '';
    this.walletAddress = '';
    this.message = '';
  }
}
