export enum WithdrawReqStatus {
  StatusPendingApproval = 10,
  StatusApproved = 20,
  StatusRejected = 30,
}

export default class WithdrawReq {
  id: bigint; //Primary key
  fromWalletAddr: string;
  fromTokenAmt: bigint;
  toWalletAddr: string;
  toAssetId: string;
  toAssetAddr: string;
  rate: bigint;
  status: WithdrawReqStatus;
  rejectReason: string;
  createdDate: Date;
  lastModifiedDate: Date;
  approveBy: bigint;
  approveDate: Date;
  userId: bigint;
  txHash: string | null;
  constructor() {
    this.id = BigInt(0);
    this.fromWalletAddr = '';
    this.fromTokenAmt = BigInt(0);
    this.toWalletAddr = '';
    this.toAssetId = '';
    this.toAssetAddr = '';
    this.rate = BigInt(0);
    this.status = WithdrawReqStatus.StatusPendingApproval;
    this.rejectReason = '';
    this.createdDate = new Date();
    this.lastModifiedDate = new Date();
    this.approveBy = BigInt(0);
    this.approveDate = new Date();
    this.userId = BigInt(0);
    this.txHash = '';
  }
}
