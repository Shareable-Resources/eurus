
import 'dart:async';

import 'package:flutter/services.dart';

class Configure {
  static const MethodChannel _channel =
      const MethodChannel('configure');

  static Future<String> get platformVersion async {
    final String version = await _channel.invokeMethod('getPlatformVersion');
    return version;
  }
}
