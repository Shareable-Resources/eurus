import 'dart:convert';

import 'package:app_logging_kit/log_fncs/other_log_config.dart';
import 'package:app_logging_kit/utils/log_msg.dart';
import 'package:flutter/material.dart';
import 'package:sentry/sentry.dart';

/// Sentry Log Functions
class SentryLog extends OtherLogConfig {
  /// Sentry endpoint that going to be connect
  final String dsn;

  SentryLog({@required this.dsn}) : super(initRequired: true);

  /// Initialize Sentry Log function
  @override
  Future<void> initFnc() async {
    await Sentry.init((options) => options..dsn = this.dsn);
    return;
  }

  /// Send logs to Sentry
  @override
  Future<void> logFnc(LogMsg msg) async {
    String strMsg = _formatLogMsg(msg);
    await Sentry.captureException(strMsg, stackTrace: msg.stackTrace);
  }

  /// Format [LogMsg] into [String]
  String _formatLogMsg(LogMsg msg) {
    LogMsg newMsg = LogMsg(msg.tag, msg.datetime, msg.message,
        networkState: msg.networkState,
        stackTrace: null,
        otherInfo: msg.otherInfo);

    Object jsonObj = newMsg.toJson();
    return jsonEncode(jsonObj);
  }
}
