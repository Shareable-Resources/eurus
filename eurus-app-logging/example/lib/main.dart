import 'package:app_logging_kit/utils/log_msg.dart';
import 'package:flutter/material.dart';

import 'package:app_logging_kit/app_logging_kit.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();

  AppLoggingKit();

  runApp(MyApp());
}

class MyApp extends StatefulWidget {
  @override
  _MyAppState createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  final AppLoggingKit logKit = AppLoggingKit();

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: Scaffold(
        appBar: AppBar(
          title: const Text('Plugin example app'),
        ),
        body: Center(
          child: Column(
            children: [
              FlatButton(onPressed: errorLog, child: Text("String Log")),
              FlatButton(onPressed: errorLog, child: Text("Error Log")),
              FlatButton(
                  onPressed: errorLog, child: Text("Custom (LogMsg) Log")),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> stringLog() async {
    await logKit.writeLog("String Log");
  }

  Future<void> errorLog() async {
    try {
      List<String> a = ['a', 'b'];
      print(a[3]);
    } catch (e) {
      await logKit.writeLog(e);
    }
  }

  Future<void> logMsg() async {
    LogMsg msg = LogMsg("custom", DateTime.now(), "msg");
    await logKit.writeLog(msg);
  }
}
