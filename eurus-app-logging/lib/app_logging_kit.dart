import 'dart:convert';
import 'package:app_logging_kit/utils/local_storage.dart';
import 'package:app_logging_kit/utils/log_msg.dart';
import 'package:connectivity/connectivity.dart';
import 'log_fncs/other_log_config.dart';

class AppLoggingKit {
  static final AppLoggingKit _singleton = AppLoggingKit._internal();

  /// Save log in local or not
  ///
  /// Default [true] for storing log locally
  bool _localLog = true;
  bool get localLog => _localLog;

  /// Local Log File location
  String _path;
  String get path => _path;

  /// file name of local log
  ///
  /// Default "log" file name
  String _logFileName = "log";
  String get logFileName => _logFileName;

  /// Save log with other method or not
  ///
  /// Deafult [false] for using other logging method
  bool _otherLog = false;
  bool get otherLog => _otherLog;

  /// Other Log Config and functions
  OtherLogConfig _otherLogConfig;
  OtherLogConfig get otherLogConfig => _otherLogConfig;

  factory AppLoggingKit({
    bool localLog,
    String path,
    String fileName,
    bool otherLog,
    OtherLogConfig otherLogConfig,
  }) {
    return _singleton;
  }

  AppLoggingKit._internal();

  /// Set local log file directory
  ///
  /// [path]: Location that stores the log file
  /// [fileName]: Local log file name
  /// [localLog]: Determine whether storing logs locally or not
  /// [otherLog]: Determine whether taking logs with other methods
  /// [otherLogConfg]: Other logging config for taking logs with other methods
  Future<void> setConfig({
    String path,
    String fileName,
    bool localLog,
    bool otherLog,
    OtherLogConfig otherLogConfig,
  }) async {
    /// Set local log folder path
    if (path != null) {
      String finalPath = path;
      bool pathValid = await ckPathExists(path);

      if (!pathValid) {
        finalPath = await createFolder(path);
      }

      _path = finalPath;
    }

    /// Set local log file name
    if (fileName != null) {
      _logFileName = fileName;
    }

    /// Set weather saving local log
    if (localLog != null) {
      _localLog = localLog;
    }

    /// Set weather taking logs with other method
    if (otherLog != null) {
      _otherLog = otherLog;
    }

    /// Set Other log configs
    if (otherLogConfig != null) {
      _otherLogConfig = otherLogConfig;
      if (_otherLogConfig.initRequired && _otherLog) _otherLogConfig.initFnc();
    }

    return;
  }

  /// Write log
  ///
  /// [msg] can be in [String], [Error] and [LogMsg]
  /// All messages will eventually convert to [LogMsg]
  Future<bool> writeLog(
    dynamic msg, {
    StackTrace st,
  }) async {
    LogMsg logMsg;
    ConnectivityResult connectivityResult =
        await (Connectivity().checkConnectivity());

    if (msg is String)
      logMsg = LogMsg(
        'String Message',
        DateTime.now(),
        msg,
        networkState: connectivityResult.toString(),
        stackTrace: StackTrace.current,
      );

    if (msg is Error)
      logMsg = LogMsg(
        'Error',
        DateTime.now(),
        msg.toString(),
        networkState: connectivityResult.toString(),
        stackTrace: msg.stackTrace,
      );

    if (msg is LogMsg) {
      logMsg = LogMsg(
        msg.tag,
        msg.datetime,
        msg.message,
        networkState: connectivityResult.toString(),
        stackTrace: msg.stackTrace,
        otherInfo: msg.otherInfo,
      );
    }

    if (_localLog) {
      String formattedLog = _formatLog(logMsg);
      await writeFile(_path, "$_logFileName.txt", formattedLog);
    }

    if (_otherLog && _otherLogConfig != null) {
      _otherLogConfig.logFnc(logMsg);
    }

    return true;
  }

  /// Convert [LogMsg] to log string
  ///
  /// Log contains datetime tag and json string
  /// e.g. 2021-01-01 : {jsoncontent}
  String _formatLog(LogMsg msg) {
    String jsonString = jsonEncode(msg.toJson());
    String datetimeTag = _genReadableDateTime(msg.datetime);

    return "$datetimeTag:$jsonString\n";
  }

  /// Generate date time in specific format
  ///
  /// return datetime in [String] with format yyyy-MM-dd hh:mm:ss
  String _genReadableDateTime(DateTime dt) {
    String year = "${dt.year}";
    String month = "${dt.month > 9 ? '' : '0'}${dt.month}";
    String day = "${dt.day > 9 ? '' : '0'}${dt.day}";
    String hour = "${dt.hour > 9 ? '' : '0'}${dt.hour}";
    String minute = "${dt.minute > 9 ? '' : '0'}${dt.minute}";
    String second = "${dt.second > 9 ? '' : '0'}${dt.second}";

    return "$year-$month-$day $hour:$minute:$second";
  }
}
