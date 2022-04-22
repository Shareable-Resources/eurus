export enum MerchantClientStatus {
  StatusInactive = 10,
  StatusActive = 20,
}

export default class MerchantClient {
  id: bigint | null;
  username: string;
  balance: bigint;
  walletAddress: string;
  createdDate: Date;
  lastModifiedDate: Date;
  status: MerchantClientStatus;
  constructor() {
    this.id = null;
    this.username = '';
    this.balance = BigInt(0);
    this.walletAddress = '';
    this.createdDate = new Date();
    this.lastModifiedDate = new Date();
    this.status = MerchantClientStatus.StatusActive;
  }
}
