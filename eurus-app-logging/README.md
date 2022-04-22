# app_logging_kit

app_logging_kit is a plugin for application to store different kinds of log locally or into third party logging systems e.g. Sentry

## Usage
### Initialize Plugin
```dart
import 'package:app_logging_kit/app_logging_kit.dart';

void main() {
    WidgetsFlutterBinding.ensureInitialized();
    AppLoggingKit();

    runApp(MyApp());
}
```

### Update plugin settings
```dart
import 'package:app_logging_kit/app_logging_kit.dart';

final AppLoggingKit logKit = AppLoggingKit();

logKit.setConfig(
    // Save log locally or not
    localLog: true,
    // Log folder location
    path: "YOUR/CUSTOME/LOCATION",
    // Log file name
    fileName: "CUSTOM_FILENAME",
    // Save log with other service
    otherLog: false,
    // Other log service config with [OtherLogConfig]
    otherLogConfig: OtherLogConfig(),
);
```

### Write log locally
```dart
import 'package:app_logging_kit/app_logging_kit.dart';

final AppLoggingKit logKit = AppLoggingKit();

// In String
await logKit.writeLog("Custom Log");
/// In Error
try {
    List<String> a = ['a', 'b'];
    // Simulates error call non-existing  index 3
    print(a[3]);
} catch (e) {
    await logKit.writeLog(e);
}
/// In LogMsg
LogMsg msg = LogMsg("custom", DateTime.now(), "msg");
await logKit.writeLog(msg);
```

### LogMsg Class
#### Log class we use to handle every log
```dart
LogMsg msg = LogMsg(
    // Required: Tag
    "other",
    // Required: DateTime
    DateTime.parse("2021-01-01"),
    // Required: Message
    "Other testing",
    // Optional: Network state
    networkState: "wifi",
    // Optional: StackTrace
    stackTrace: StackTrace.empty,
    // Optional: Other infos wanted to be logged
    otherInfo: {"dm_2": 1, "dm_1": "two"},
);
```

### OtherLogConfig
#### For developer to connect their own logging system into this plugin
#### Example
```dart
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
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
