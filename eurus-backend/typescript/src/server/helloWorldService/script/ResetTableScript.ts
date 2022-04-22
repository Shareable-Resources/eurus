#!/usr/bin/env node

import seq from '../sequelize';
import * as SeqModel from '../model/seqModel/0_index';
import * as DBModel from '../model/dbModel/0_index';
import { HelloWorldUserStatus } from '../model/DBModel/HelloWorldUser';

export async function resetTable() {
  const dateNow = new Date();
  await seq.assertDatabaseConnectionOk();
  const modelModule = seq.sequelize.models;
  await seq.sequelize.drop();
  await seq.sequelize.sync({ force: true });

  const dummyUsers: DBModel.HelloWorldUser[] = [
    {
      id: null,
      username: 'user1',
      balance: BigInt(0),
      walletAddress: '',
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: HelloWorldUserStatus.StatusActive,
    },
    {
      id: null,
      username: 'user2',
      balance: BigInt(1500),
      walletAddress: '',
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: HelloWorldUserStatus.StatusActive,
    },
    {
      id: null,
      username: 'user3',
      balance: BigInt(0),
      walletAddress: '',
      createdDate: dateNow,
      lastModifiedDate: dateNow,
      status: HelloWorldUserStatus.StatusActive,
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
      status: HelloWorldUserStatus.StatusActive,
    },
  ];

  await modelModule[SeqModel.name.HelloWorldUser].bulkCreate(dummyUsers);
  //console.log by css
  console.log(
    '%c ----- Done creating dummy table and data for DB! ----- ',
    'color:green',
  );
}
