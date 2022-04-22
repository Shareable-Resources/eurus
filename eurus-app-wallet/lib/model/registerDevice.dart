class RegisterDevice {
  int returnCode;
  String message;
  String nonce;
  Data? data;

  RegisterDevice(this.returnCode, this.message, this.nonce, this.data);

  RegisterDevice.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null ? new Data.fromJson(json['data']!) : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class Data {
  String code;

  Data(this.code);

  Data.fromJson(Map<String, dynamic> json) : code = json['code'];

  Map<String, dynamic> toJson() => {
        'code': this.code,
      };
}
