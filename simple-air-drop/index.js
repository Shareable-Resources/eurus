const fs = require('fs');
const EthCrypto = require('eth-crypto');
const logFile = 'output-' + Date.now() + '.log';
const stream = fs.createWriteStream(logFile, { flags: 'a' });

const captains = {
  log: function () {
    console.log(...arguments);
    stream.write([...arguments].join(' ') + '\n');
  },
  error: function () {
    console.log("Error:", ...arguments);
    stream.write("Error: " + [...arguments].join(' ') + '\n');
  },
};

const argv = require('yargs').argv;

const argExit = error => {
  captains.error(error);
  captains.error('Required syntax is: yarn start --config <configpath> --csv <csvpath>');
  process.exit(1);
};

if (!argv.config) {
  argExit('Missing config file path argument... exiting...');
} else if (!argv.csv) {
  argExit('Missing csv file path argument... exiting...');
}

let config;

try {
  config = require(argv.config);
} catch (e) {
  captains.log('Failed to import config file...');
  captains.log(e.message, e.stack);
  process.exit(1);
}

let input;

try {
  input = fs.readFileSync(argv.csv);
} catch (e) {
  captains.log('Failed to import csv file...');
  captains.log(e.message, e.stack);
  process.exit(1);
}

// Dependencies
const BigNumber = require('bignumber.js');
const parse = require('csv-parse/lib/sync');

const records = parse(input, {
  columns: true,
  skip_empty_lines: true,
});

const sleep = time => {
  return new Promise(resolve => {
    setTimeout(resolve, time);
  });
};

const processAccount = async (count, record, web3) => {
  captains.log("--");
  captains.log('Processing (', count, ')', record.address, 'for amount of', record.amount);

  // check balances before
  const balance = new BigNumber(web3.utils.fromWei(await web3.eth.getBalance(record.address), 'ether'));
  const srcBalance = new BigNumber(web3.utils.fromWei(await web3.eth.getBalance(config.srcAccount), 'ether'));
  captains.log('Pre airdrop balance (target, source): (', balance.toString(), ',', srcBalance.toString(), ')');

  // do airdrop
  let tx = {
    value: web3.utils.toHex(new BigNumber(web3.utils.toWei(record.amount, 'ether'))),
    from: config.srcAccount,
    to: record.address,
    gas: config.gasLimit
  };

  const response = await web3.eth.sendTransaction(tx);

  captains.log("TxHash: ", response ? response.transactionHash : "null");
  captains.log(">>", JSON.stringify(response));

  const balance2 = new BigNumber(web3.utils.fromWei(await web3.eth.getBalance(record.address), 'ether'));
  const srcBalance2 = new BigNumber(web3.utils.fromWei(await web3.eth.getBalance(config.srcAccount), 'ether'));
  captains.log('Post airdrop balance (target, source): (', balance2.toString(), ',', srcBalance2.toString(), ')');

  if (!balance.plus(record.amount).eq(balance2)) {
    captains.log('##################################');
    captains.log('### ERROR TARGET BALANCE INCORRECT');
    captains.log('##################################');
    await sleep(5000);
  }

  const gasInWei = new BigNumber(response.gasUsed).multipliedBy(new BigNumber(response.effectiveGasPrice));
  const gas = new BigNumber(web3.utils.fromWei(gasInWei.toFixed(), 'ether')); //2.1 * 10^4 * 2.4 * 10^9

  if (!srcBalance.minus(record.amount).minus(gas).eq(srcBalance2)) {
    captains.log('##################################');
    captains.log('### ERROR SOURCE BALANCE INCORRECT');
    captains.log('##################################');
    await sleep(5000);
  }

  captains.log("Paid>>",
    JSON.stringify(
      {
        "amount": record.amount,
        "gas": gas,
        "address": record.address,
        "tx": response ? response.transactionHash : null
      }
    ));

};

const run = async () => {
  // captains.log(records);
  // captains.log(records.length);
  const Web3 = require('web3');
  const web3 = new Web3(config.rpcUrl);

  const privateKey = config.srcPrivateKey;
  const publicKey = EthCrypto.publicKeyByPrivateKey(privateKey);
  const address = EthCrypto.publicKey.toAddress(publicKey);

  if (address != config.srcAccount) {
    captains.error('Private key not match address.');
    return;
  } else {
    captains.log(`Air drop from this address ${address}`);
  }
  web3.eth.accounts.wallet.add(privateKey);

  const theSrcBalance = new BigNumber(web3.utils.fromWei(await web3.eth.getBalance(config.srcAccount), 'ether'));

  let accounts = {};
  let goodAccounts = [];

  let gasEstimatedInEth = new BigNumber(config.gasEstimatedInEth);
  let sum = new BigNumber(0);
  let duplicateSum = new BigNumber(0);
  let valid = true;
  for (let i = 0; i < records.length; i++) {
    let amount = new BigNumber(records[i].amount);
    sum = sum.plus(amount);
    if (!accounts[records[i].address]) {
      if (!web3.utils.isAddress(records[i].address)) {
        valid = false;
        captains.log('Bad ADDRESS:', records[i].address, 'skipping address...');
      } else if (Number.isNaN(records[i].amount) || Number(records[i].amount) <= 0) {
        valid = false;
        captains.log('Bad AMOUNT:', records[i].address, records[i].amount, 'skipping address...');
      } else {
        goodAccounts.push(records[i]);
        accounts[records[i].address] = true;
        duplicateSum = duplicateSum.plus(amount);
      }
    }
  }

  if (valid === false) {
    captains.log('##########################################################');
    captains.log(
      'There were bad addresses, stop program if you wish to fix, otherwise continuing and skipping those payouts...'
    );
    captains.log('##########################################################');
    await sleep(5000);
  }

  captains.log('Source Account Pre Airdrop Balance:', theSrcBalance.toString());
  captains.log('Total Addresses:', records.length);
  captains.log('Total Amount:', sum);
  captains.log('Total Addresses After Duplciate Removal:', goodAccounts.length);
  captains.log('Total Amount After Duplicate Removal:', duplicateSum);

  const totalGas = gasEstimatedInEth.multipliedBy(new BigNumber(goodAccounts.length));
  const totalAmount = duplicateSum.plus(totalGas);
  captains.log('Total Gas After Duplicate Removal:', totalGas);
  captains.log('Total amount estimated:', totalAmount);

  if (theSrcBalance.lt(totalAmount)) {
    captains.error('Insufficient balance. Total amount needed:', totalAmount, "Account balance:", theSrcBalance);
    return;
  }


  captains.log('Starting payout in 5 seconds');

  await sleep(5000);

  for (let i = 0; i < goodAccounts.length; i++) {
    await processAccount(i + 1 + '/' + goodAccounts.length, goodAccounts[i], web3);
  }

  const theSrcBalance2 = new BigNumber(web3.utils.fromWei(await web3.eth.getBalance(config.srcAccount), 'ether'));

  captains.log('Total Paid To', goodAccounts.length, 'Addresses:', theSrcBalance.minus(theSrcBalance2).toString());
};

run()
  .then(() => {
    captains.log('Program finished...');
    captains.log('Log saved to', logFile);
    stream.end();
  })
  .catch(e => {
    captains.log('Uncaught error in program run: ', e.message, e.stack);
    captains.log('Log saved to', logFile);
    stream.end();
  });
