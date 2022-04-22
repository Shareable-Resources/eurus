#!/usr/bin/env node

import seq from '../sequelize';
import * as SeqModel from '../model/seqModel/0_index';
import * as DBModel from '../model/dbModel/0_index';
import { DepositLedgerStatus } from '../model/dbModel/DepositLedger';
import { WithdrawLedgerStatus } from '../model/dbModel/WithdrawLedger';
import { WithdrawReqStatus } from '../model/dbModel/WithdrawReq';

export async function resetTable() {
  const dateNow = new Date();
  await seq.assertDatabaseConnectionOk();
  const modelModule = seq.sequelize.models;
  await seq.sequelize.drop();
  await seq.sequelize.sync({ force: true });

  const dummyMerchantAdmins: DBModel.MerchantAdmin[] = [
    {
      operatorId: BigInt(1),
      merchantId: BigInt(1),
      username: 'MC1',
      email: 'abcdefg@gmail.com',
      passwordHash:
        '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92',
      status: 1,
      createdDate: dateNow,
      lastModifiedDate: dateNow,
    },
    {
      operatorId: BigInt(2),
      merchantId: BigInt(1),
      username: 'MC2',
      email: 'abcdefg2@gmail.com',
      passwordHash:
        '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92',
      status: 1,
      createdDate: dateNow,
      lastModifiedDate: dateNow,
    },
    {
      operatorId: BigInt(3),
      merchantId: BigInt(2),
      username: 'TC1',
      email: 'abcdefg3@gmail.com',
      passwordHash:
        '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92',
      status: 1,
      createdDate: dateNow,
      lastModifiedDate: dateNow,
    },
    {
      operatorId: BigInt(4),
      merchantId: BigInt(2),
      username: 'TC2',
      email: 'abcdefg4@gmail.com',
      passwordHash:
        '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92',
      status: 1,
      createdDate: dateNow,
      lastModifiedDate: dateNow,
    },
    {
      operatorId: BigInt(5),
      merchantId: BigInt(2),
      username: 'TC3',
      email: 'abcdefg5@gmail.com',
      passwordHash:
        '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92',
      status: 1,
      createdDate: dateNow,
      lastModifiedDate: dateNow,
    },
  ];

  const dummyUsers: DBModel.MerchantClient[] = [
    {
      id: null,
      username: 'user1',
      balance: BigInt(0),
      walletAddress: '',
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: 1,
    },
    {
      id: null,
      username: 'user2',
      balance: BigInt(1500),
      walletAddress: '',
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: 1,
    },
    {
      id: null,
      username: 'user3',
      balance: BigInt(0),
      walletAddress: '',
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: 1,
    },
    {
      id: null,
      username: 'user4',
      balance: BigInt(0),
      walletAddress: '',
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: 1,
    },
    {
      id: null,
      username: 'user5',
      balance: BigInt(0),
      walletAddress: '',
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: 1,
    },
  ];

  const dummyWithdrawLedgers: DBModel.WithdrawLedger[] = [
    {
      txHash: '0xsafasfasfasf412412',
      fromWalletAddr: '0xfakeMerchantWalletAddr',
      fromTokenAmt: BigInt(1000),
      toWalletAddr: '0xfakeUserWalletAddr',
      toAssetId: 'USDT',
      toAssetAddr: '0xfakeAssetAddr',
      toAssetAmt: BigInt(1000),
      toUserId: BigInt(0),
      rate: BigInt(1),
      createdDate: new Date(),
      lastModifiedDate: new Date(),
      status: WithdrawLedgerStatus.StatusCreated,
      reqId: null,
      remarks: '',
    },
    {
      txHash: '0xsafasfasfasf412413',
      fromWalletAddr: '0xfakeMerchantWalletAddr',
      fromTokenAmt: BigInt(1000),
      toWalletAddr: '0xfakeUserWalletAddr',
      toAssetId: 'USDT',
      toAssetAddr: '0xfakeAssetAddr',
      toAssetAmt: BigInt(1000),
      toUserId: BigInt(0),
      rate: BigInt(1),
      createdDate: new Date(),
      lastModifiedDate: new Date(),
      status: WithdrawLedgerStatus.StatusTransferProcessing,
      reqId: null,
      remarks: '',
    },
    {
      txHash: '0xsafasfasfasf412414',
      fromWalletAddr: '0xfakeMerchantWalletAddr',
      fromTokenAmt: BigInt(1000),
      toWalletAddr: '0xfakeUserWalletAddr',
      toAssetId: 'USDT',
      toAssetAddr: '0xfakeAssetAddr',
      toAssetAmt: BigInt(1000),
      toUserId: BigInt(0),
      rate: BigInt(1),
      createdDate: new Date(),
      lastModifiedDate: new Date(),
      status: WithdrawLedgerStatus.StatusConfirmingTransfer,
      reqId: null,
      remarks: '',
    },
    {
      txHash: '0xsafasfasfasf412415',
      fromWalletAddr: '0xfakeMerchantWalletAddr',
      fromTokenAmt: BigInt(1000),
      toWalletAddr: '0xfakeUserWalletAddr',
      toAssetId: 'USDT',
      toAssetAddr: '0xfakeAssetAddr',
      toAssetAmt: BigInt(1000),
      toUserId: BigInt(0),
      rate: BigInt(1),
      createdDate: new Date(),
      lastModifiedDate: new Date(),
      status: WithdrawLedgerStatus.StatusCompleted,
      reqId: null,
      remarks: '',
    },
  ];
  const dummyDeposits: DBModel.DepositLedger[] = [
    {
      txHash: '0x999999999989898989898989999999999999',
      fromWalletAddr: '0xfakeUserWalletAddr',
      fromAssetId: 'USDT',
      fromAssetAddr: '0xfakeEUNSmartContractAddr',
      fromAssetAmt: BigInt(20),
      toWalletAddr: '0xfakeChantWalletAddr',
      toTokenAmt: BigInt(500),
      rate: BigInt(1),
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: DepositLedgerStatus.DepositReceiptCollected,
      userId: BigInt(0),
      remarks: 'FAKE',
    },
    {
      txHash: '0x999999999989898989898989999999999998',
      fromWalletAddr: '0xfakeUserWalletAddr',
      fromAssetId: 'USDT',
      fromAssetAddr: '0xfakeEUNSmartContractAddr',
      fromAssetAmt: BigInt(20),
      toWalletAddr: '0xfakeChantWalletAddr',
      toTokenAmt: BigInt(500),
      rate: BigInt(1),
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: DepositLedgerStatus.DepositReceiptCollected,
      userId: BigInt(0),
      remarks: 'FAKE',
    },
    {
      txHash: '0x999999999989898989898989999999999997',
      fromWalletAddr: '0xfakeUserWalletAddr',
      fromAssetId: 'USDT',
      fromAssetAddr: '0xfakeEUNSmartContractAddr',
      fromAssetAmt: BigInt(20),
      toWalletAddr: '0xfakeChantWalletAddr',
      toTokenAmt: BigInt(250),
      rate: BigInt(1),
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: DepositLedgerStatus.DepositReceiptCollected,
      userId: BigInt(1),
      remarks: 'FAKE',
    },
    {
      txHash: '0x999999999989898989898989999999999996',
      fromWalletAddr: '0xfakeUserWalletAddr',
      fromAssetId: 'USDT',
      fromAssetAddr: '0xfakeEUNSmartContractAddr',
      fromAssetAmt: BigInt(20),
      toWalletAddr: '0xfakeChantWalletAddr',
      toTokenAmt: BigInt(150),
      rate: BigInt(1),
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: DepositLedgerStatus.DepositReceiptCollected,
      userId: BigInt(1),
      remarks: 'FAKE',
    },
  ];

  const dummyWithdrawReqs: DBModel.WithdrawReq[] = [
    {
      id: BigInt(1),
      fromWalletAddr: '0xFromFakeMerchantSmartContractDAppAddres',
      fromTokenAmt: BigInt(20),
      toWalletAddr: '0xFakeUserWalletAddress',
      toAssetId: 'USDT',
      toAssetAddr: '0xFakeUSDTSmartContractAddress',
      rate: BigInt(1),
      status: WithdrawReqStatus.StatusRejected,
      rejectReason: 'Fake reject reason',
      createdDate: new Date(),
      lastModifiedDate: new Date(),
      approveBy: BigInt(1),
      approveDate: new Date(),
      userId: BigInt(2),
      txHash: null,
    },
    {
      id: BigInt(2),
      fromWalletAddr: '0xFromFakeMerchantSmartContractDAppAddres',
      fromTokenAmt: BigInt(20),
      toWalletAddr: '0xFakeUserWalletAddress',
      toAssetId: 'USDT',
      toAssetAddr: '0xFakeUSDTSmartContractAddress',
      rate: BigInt(1),
      status: WithdrawReqStatus.StatusApproved,
      rejectReason: '',
      createdDate: new Date(),
      lastModifiedDate: new Date(),
      approveBy: BigInt(1),
      approveDate: new Date(),
      userId: BigInt(2),
      txHash: null,
    },
  ];

  await modelModule[SeqModel.name.MerchantAdmins].bulkCreate(
    dummyMerchantAdmins,
  );
  //await modelModule[SeqModel.name.WithdrawReqs].bulkCreate(dummyWithdrawReqs);
  //await modelModule[SeqModel.name.MerchantClient].bulkCreate(dummyUsers);
  //await modelModule[SeqModel.name.DepositLedgers].bulkCreate(dummyDeposits);
  //await modelModule[SeqModel.name.WithdrawLedgers].bulkCreate(dummyWithdrawLedgers);
  //console.log by css
  console.log(
    '%c ----- Done creating dummy table and data for DB! ----- ',
    'color:green',
  );
}
