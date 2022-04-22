import path from 'path';
import winston from 'winston';
import 'winston-daily-rotate-file';
export const folderPath = {
  dist: path.join(__dirname, '../../../logs'),
  proj: 'logs',
};

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
    s
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

const loggerHelper = {
  createRotateLogger(
    filename: string,
    timeInterval?: string,
    logInConsole?: boolean
  ) {
    timeInterval = timeInterval ? timeInterval : '14d';
    const transports: winston.transport | winston.transport[] | undefined = [
      new winston.transports.DailyRotateFile({
        filename: `${folderPath.dist}/error/${filename}-%DATE%.log`,
        datePattern: 'YYYY-MM-DD',
        zippedArchive: true,
        maxSize: '128m',
        maxFiles: '14d',
        json: true,
        level: 'error',
        createSymlink: true,
        symlinkName: `${filename}.log`,
        format: winston.format.combine(winston.format.timestamp(), errorFormat),
      }),
      new winston.transports.DailyRotateFile({
        filename: `${folderPath.dist}/info/${filename}-%DATE%.log`,
        datePattern: 'YYYY-MM-DD',
        zippedArchive: true,
        maxSize: '128m',
        maxFiles: '14d',
        json: true,
        level: 'debug',
        createSymlink: true,
        symlinkName: `${filename}.log`,
        format: winston.format.combine(
          infoDebugWarnFilter(),
          winston.format.timestamp(),
          customFormat
        ),
      }),
    ];
    if (logInConsole) {
      transports.push(
        new winston.transports.Console({
          format: winston.format.simple(),
          level: 'debug',
        })
      );
    }
    return winston.createLogger({
      levels: custom.levels,
      exitOnError: false,
      format: winston.format.combine(
        winston.format.timestamp(),
        customFormat,
        winston.format.colorize(),
        winston.format.errors({ stack: true })
      ),
      transports: transports,
    });
  },
};

export default loggerHelper;
