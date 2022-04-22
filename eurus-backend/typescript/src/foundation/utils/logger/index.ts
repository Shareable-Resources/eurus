import path from 'path';
import winston, { createLogger, format } from 'winston';
import 'winston-daily-rotate-file';
import DailyRotateFile from 'winston-daily-rotate-file';
const { colorize, combine, timestamp, label, printf, json } = winston.format;


const custom = {
  levels: {
    error: 0,
    warn: 1,
    info: 2,
    verbose: 3,
    debug: 4,
    http: 5,
  },
  colors: {
    error: 'red',
    warn: 'orange',
    info: 'white bold yellow',
    verbose: 'blue',
    debug: 'green',
    http: 'pink',
  },
};
winston.addColors(custom.colors);

function retriveSymbolLevel(info: winston.Logform.TransformableInfo) {
  const symbolLevel: any = Object.getOwnPropertySymbols(info).find(function (
    s,
  ) {
    return String(s) === 'Symbol(level)';
  });
  const level = info[symbolLevel];
  return level;
}

const infoDebugWarnFilter = winston.format((info, opts) => {
  const level = retriveSymbolLevel(info);

  const levelToBeLogged = ['info', 'warn', 'debug'];
  return levelToBeLogged.includes(level) ? info : false;
});

const customFormat = winston.format.printf((info) => {
  const level = retriveSymbolLevel(info);
  return `[${level.toUpperCase()}]: [${info.timestamp}] [${info.message}]`;
});

const errorFormat = winston.format.printf((info) => {
  const level = retriveSymbolLevel(info);
  return `[${level.toUpperCase()}]: [${info.timestamp}] [${info.message}] [${
    info.stack
  }]`;
});

function parseFileName(filePath : string): any {
  let returnInfo  = { path: "" , fileName : "", ext: ""};
  let pos = filePath.lastIndexOf(".")
  if (pos >= 0) {
    let pathPos = filePath.lastIndexOf("/")
    let start = 0;
    if (pathPos >= 0){
      start = pathPos + 1;
      returnInfo["path"] = filePath.substring(0, pathPos + 1);
    }else {
      returnInfo["path"] = "./";
    }
    returnInfo["fileName"] = filePath.substring(start, pos);
    returnInfo["ext"] = filePath.substring(pos)
  }else {
    let pathPos = filePath.lastIndexOf("/")
    let start = 0;
    if (pathPos >= 0){
      start = pathPos + 1;
      returnInfo["path"] = filePath.substring(0, pathPos + 1);
    }else {
      returnInfo["path"] = "./";
    } 
    returnInfo["fileName"] = filePath.substring(start);
  }
  return returnInfo;
}

const loggerHelper = {
  createRotateLogger(
    filename: string,
    timeInterval?: string,
    logInConsole?: boolean,
  ) {
    if (!filename){
      throw new Error("log file name is undefined");
    }
    let fileInfo = parseFileName(filename);
    timeInterval = timeInterval ? timeInterval : '14d';
    const transports: winston.transport | winston.transport[] | undefined = [
      new winston.transports.DailyRotateFile({
        filename: `${fileInfo["path"]}/error/${fileInfo["fileName"]}-%DATE%${fileInfo["ext"]}`,
        datePattern: 'YYYY-MM-DD',
        zippedArchive: true,
        maxSize: '128m',
        maxFiles: '14d',
        json: true,
        level: 'error',
        createSymlink: true,
        symlinkName: `${fileInfo["fileName"]}${fileInfo["ext"]}`,
        format: winston.format.combine(winston.format.timestamp(), errorFormat),
      }),
      new winston.transports.DailyRotateFile({
        filename: `${fileInfo["path"]}/info/${fileInfo["fileName"]}-%DATE%${fileInfo["ext"]}`,
        datePattern: 'YYYY-MM-DD',
        zippedArchive: true,
        maxSize: '128m',
        maxFiles: '14d',
        json: true,
        level: 'debug',
        createSymlink: true,
        symlinkName: `${fileInfo["fileName"]}${fileInfo["ext"]}`,
        format: winston.format.combine(
          infoDebugWarnFilter(),
          winston.format.timestamp(),
          customFormat,
        ),
      }),
    ];
    if (logInConsole) {
      transports.push(
        new winston.transports.Console({
          format: winston.format.simple(),
          level: 'debug',
        }),
      );
    }
    return winston.createLogger({
      levels: custom.levels,
      exitOnError: false,
      format: winston.format.combine(
        winston.format.timestamp(),
        customFormat,
        winston.format.colorize(),
        winston.format.errors({ stack: true }),
      ),
      transports: transports,
    });
  },
};

export default loggerHelper;
