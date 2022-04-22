import 'package:app_logging_kit/app_logging_kit.dart';
import 'package:app_logging_kit/log_fncs/other_log_config.dart';
import 'package:app_logging_kit/utils/log_msg.dart';
import 'package:flutter/material.dart';

class OtherlogTestConfig extends OtherLogConfig {
  final bool initRequired;

  OtherlogTestConfig({@required this.initRequired})
      : super(initRequired: initRequired);

  @override
  Future<String> initFnc() async {
    await AppLoggingKit().writeLog("Init Success");

    return "Init Success";
  }

  @override
  Future<String> logFnc(LogMsg msg) async {
    await AppLoggingKit().writeLog("${msg.tag} - ${msg.message} :: to local");

    return "${msg.tag} - ${msg.message}";
  }
}
