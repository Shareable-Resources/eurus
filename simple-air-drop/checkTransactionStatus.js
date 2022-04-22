const fs = require('fs');
const lineReader = require('line-reader');
const argv = require('yargs').argv;
const Web3 = require('web3');
const BigNumber = require('bignumber.js');

const logFile = 'check-' + Date.now() + '.log';
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

const LINE_PREFIX = "Paid>>";
const argExit = error => {
    captains.error(error);
    captains.error('Required syntax is: yarn start --log <logFilePath> --config <configFilePath>');
    process.exit(1);
};

if (!argv.config) {
    argExit('Missing config file path argument... exiting...');
} else if (!argv.log) {
    argExit('Missing log file path argument... exiting...');
}

let config;

try {
    config = require(argv.config);
} catch (e) {
    captains.log('Failed to import config file...');
    captains.log(e.message, e.stack);
    process.exit(1);
}

const web3 = new Web3(config.rpcUrl);

async function readline(logPath) {
    let records = [];
    return new Promise((resolve, reject) => {
        lineReader.eachLine(logPath, function (line, last) {
            if (line.startsWith(LINE_PREFIX)) {
                const json = line.substring(LINE_PREFIX.length);
                //captains.log(json);
                const record = JSON.parse(json);
                records.push(record);
            }
        }, function (err) {
            if (err) {
                reject(err);
            } else {
                resolve(records);
            }
        });
    })
}

async function processRecords(records) {
    const currentBlockNumber = await web3.eth.getBlockNumber();

    let total = records.length;
    let count = 0;
    let successCount = 0;
    let failCount = 0;

    for (let i = 0; i < total; i++) {
        let record = records[i];
        let transaction = await web3.eth.getTransaction(record.tx)
        //captains.log(transaction);
        if (transaction) {
            const confirmedBlock = currentBlockNumber - transaction.blockNumber;
            const valueInRecord = new BigNumber(record.amount);
            const valueInTxn = new BigNumber(transaction.value).div(new BigNumber(1000000000000000000));
            if (valueInRecord.eq(valueInTxn)) {
                captains.log(`${confirmedBlock} blocks confirmed ${record.tx}, ${valueInRecord}`);
                successCount++
            } else {
                captains.log(`Txn amount not match ${record.tx}, valueInRecord=${valueInRecord}, valueInTxn=${valueInTxn}`);
                failCount++
            }
        } else {
            captains.error("***Fail***", record)
            failCount++
        }
        count++
    }
    return { count, successCount, failCount };
}


async function run() {

    captains.log(`Start to read log file ${argv.log}`);

    let records = await readline(argv.log);

    captains.log(`Total to process ${records.length}`);

    let processResult = await processRecords(records);

    captains.log(`Processed ${processResult.count}, Sucess:${processResult.successCount}, Fail:${processResult.failCount}`);

}

run()
    .then(() => {
        captains.log('Program finished...');
        captains.log('Log saved to', logFile);
        stream.end();
    })
    .catch(e => {
        captains.error('Uncaught error in program run: ', e.message, e.stack);
        captains.log('Log saved to', logFile);
        stream.end();
    });