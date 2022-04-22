export enum MerchantAdminStatus {
  StatusInactive = 10,
  StatusActive = 20,
}
export default class MerchantAdmin {
  operatorId: bigint; //Auto increment
  merchantId: bigint;
  username: string;
  email: string;
  passwordHash: string;
  status: MerchantAdminStatus;
  createdDate: Date;
  lastModifiedDate: Date;
  constructor() {
    this.operatorId = BigInt(0);
    this.merchantId = BigInt(0);
    this.username = '';
    this.email = '';
    this.passwordHash = '';
    this.status = MerchantAdminStatus.StatusActive;
    this.createdDate = new Date();
    this.lastModifiedDate = new Date();
  }
}
