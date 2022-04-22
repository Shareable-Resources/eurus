import 'dart:io';

import 'package:app_logging_kit/utils/local_storage.dart';
import 'package:app_logging_kit/utils/log_msg.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:app_logging_kit/app_logging_kit.dart';

import 'otherlog_test_config.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  setUpAll(() async {
    // Create a temporary directory.
    final directory = await Directory.systemTemp.createTemp();

    // Mock out the MethodChannel for the path_provider plugin.
    const MethodChannel('plugins.flutter.io/path_provider')
        .setMockMethodCallHandler((MethodCall methodCall) async {
      // If you're getting the apps documents directory, return the path to the
      // temp directory on the test environment instead.
      if (methodCall.method == 'getApplicationDocumentsDirectory') {
        return directory.path;
      }
      return null;
    });
  });

  AppLoggingKit logKit = AppLoggingKit();

  String orgPath;

  group("Plugin Basic Functions", () {
    test("Init Plugin", () {
      expect(true, logKit.path != null);
      expect(true, logKit.localLog);
      expect(false, logKit.otherLog);
      expect("log", logKit.logFileName);
      expect(null, logKit.otherLogConfig);
    });

    orgPath = logKit.path;

    test("Update local log flag", () {
      logKit.setConfig(localLog: false);
      expect(false, logKit.localLog);
    });
    test("Update local log folder path", () {
      logKit.setConfig(path: "${logKit.path}/customPath");
      expect("$orgPath/customPath", logKit.path);
    });
    test("Update local log filename", () {
      logKit.setConfig(fileName: "customName");
      expect("customName", logKit.logFileName);
    });
    test("Update other log flag", () {
      logKit.setConfig(otherLog: true);
      expect(false, logKit.localLog);
    });
    test("Update other log config", () {
      logKit.setConfig(otherLogConfig: OtherlogTestConfig(initRequired: true));
      expect(true, logKit.otherLogConfig != null);
    });
  });

  group("Local Storage", () {
    test("Check Path exists", () async {
      bool exists = await ckPathExists(logKit.path);
      expect(true, exists);
    });
    test("Create Folder", () async {
      String folderPath = await createFolder("$orgPath/newfolder");
      bool exists = await ckPathExists(folderPath);
      expect(true, exists);
    });
    test("Check file exists", () async {
      bool exists = await ckFileExists(logKit.path, logKit.logFileName);
      expect(false, exists);
    });
    test("Write new file", () async {
      bool wrote =
          await writeFile(logKit.path, logKit.logFileName, "A testing Log");
      expect(true, wrote);
      bool exists = await ckFileExists(logKit.path, logKit.logFileName);
      expect(true, exists);
    });
    test("Read log file", () async {
      List<String> logs = await readFile(logKit.path, logKit.logFileName);
      expect(1, logs.length);
      expect("A testing Log", logs[0]);
    });
    test("Delete file", () async {
      bool deleted = await deleteFile(logKit.path, logKit.logFileName);
      expect(true, deleted);
      bool exists = await ckFileExists(logKit.path, logKit.logFileName);
      expect(false, exists);
    });
  });

  group("Write local log", () {
    logKit.setConfig(
      localLog: true,
      path: "$orgPath/testing",
      fileName: "testingLog",
      otherLog: false,
    );
    test("Write log with String", () async {
      bool wrote = await logKit.writeLog("Dummy Log");
      expect(true, wrote);
    });
    test("Write log with Error", () async {
      try {
        List<String> a = ['a', 'b'];
        print(a[3]);
      } catch (e) {
        bool wrote = await logKit.writeLog(e);
        expect(true, wrote);
      }
    });
    test("Write log with LogMsg", () async {
      LogMsg msg = LogMsg(
        "testing",
        DateTime.parse("2021-01-01"),
        "This is a testing",
        stackTrace: StackTrace.empty,
        otherInfo: {"dm_2": 1, "dm_1": "two"},
      );
      bool wrote = await logKit.writeLog(msg);
      expect(true, wrote);
    });
  });
  group("Write other log", () {
    logKit.setConfig(
      localLog: false,
      otherLog: true,
      otherLogConfig: OtherlogTestConfig(initRequired: true),
    );

    LogMsg msg = LogMsg(
      "other",
      DateTime.parse("2021-01-01"),
      "Other testing",
      stackTrace: StackTrace.empty,
      otherInfo: {"dm_2": 1, "dm_1": "two"},
    );
    test("Init Function", () async {
      String testMsg = await logKit.otherLogConfig.initFnc();
      expect("Init Success", testMsg);
      List<String> logs = await readFile(logKit.path, logKit.logFileName);
      expect("Init Success", logs[logs.length - 1]);
    });
    test("Log Function", () async {
      String testMsg = await logKit.otherLogConfig.logFnc(msg);
      expect("other - Other testing", testMsg);
    });
    test("Write log", () async {
      bool wrote = await logKit.writeLog(msg);
      expect(true, wrote);
      List<String> logs = await readFile(logKit.path, logKit.logFileName);
      expect("${msg.tag} - ${msg.message} :: to local", logs[logs.length - 1]);
    });
  });
}
