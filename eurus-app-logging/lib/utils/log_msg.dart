import 'dart:convert';

/// Log Message class
class LogMsg {
  final String tag;
  final DateTime datetime;
  final String message;
  final String networkState;
  final StackTrace stackTrace;
  final Object otherInfo;

  LogMsg(
    this.tag,
    this.datetime,
    this.message, {
    this.networkState,
    this.stackTrace,
    this.otherInfo,
  });

  LogMsg.fromJson(Map<String, dynamic> json)
      : tag = json['tag'],
        datetime = DateTime.fromMillisecondsSinceEpoch(json['datetime']),
        message = json['message'],
        networkState = json['networkState'],
        stackTrace = json['stackTrace'] == ''
            ? ''
            : StackTrace.fromString(json['stackTrace']),
        otherInfo =
            json['otherInfo'] == '' ? '' : jsonDecode(json['otherInfo']);

  Map<String, dynamic> toJson() => {
        'tag': tag,
        'datetime': DateTime.now().millisecondsSinceEpoch,
        'message': message,
        'networkState': networkState,
        'stackTrace': stackTrace != null ? stackTrace.toString() : '',
        'otherInfo': otherInfo != null ? jsonEncode(otherInfo) : ''
      };
}
