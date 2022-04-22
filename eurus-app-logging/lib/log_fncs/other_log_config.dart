import 'package:app_logging_kit/utils/log_msg.dart';
import 'package:flutter/material.dart';

/// Basic format for creating other logging methods
abstract class OtherLogConfig {
  /// Determine if a init function is needed
  final bool initRequired;

  OtherLogConfig({
    @required this.initRequired,
  });

  dynamic initFnc() async {}

  dynamic logFnc(LogMsg msg) async {}
}
