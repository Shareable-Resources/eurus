export enum DepositLedgerStatus {
  StatusError = -1,
  DepositReceiptCollected = 10,
  // StatusPendingApproval = 10,
  // StatusApproved = 20,
  // StatusRejected = 30,
  DepositAssetCollected = 20,
  DepositMintRequesting = 25,
  DepositMintConfirming = 30,
  DepositCompleted = 40,
  SweepTrans = 0x7fff,
}

export default class DepositLedger {
  txHash: string;
  fromWalletAddr: string;
  fromAssetId: string;
  fromAssetAddr: string;
  fromAssetAmt: bigint;
  toWalletAddr: string;
  toTokenAmt: bigint;
  rate: bigint;
  createdDate: Date;
  lastModifiedDate: Date;
  status: DepositLedgerStatus;
  userId: bigint;
  remarks: string;
  constructor() {
    this.txHash = '';
    this.fromWalletAddr = '';
    this.fromAssetId = '';
    this.fromAssetAddr = '';
    this.fromAssetAmt = BigInt(0);
    this.toWalletAddr = '';
    this.toTokenAmt = BigInt(0);
    this.rate = BigInt(0);
    this.createdDate = new Date();
    this.lastModifiedDate = new Date();
    this.status = DepositLedgerStatus.DepositReceiptCollected;
    this.userId = BigInt(0);
    this.remarks = '';
  }
}
