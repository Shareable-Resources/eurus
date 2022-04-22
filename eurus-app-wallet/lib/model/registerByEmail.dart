class RegisterByEmail {
  int returnCode;
  String message;
  String nonce;
  RegisterByEmailData? data;

  RegisterByEmail(this.returnCode, this.message, this.nonce, this.data);

  RegisterByEmail.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? new RegisterByEmailData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class RegisterByEmailData {
  int userId;
  String code;

  RegisterByEmailData(this.userId, this.code);

  RegisterByEmailData.fromJson(Map<String, dynamic> json)
      : userId = json['userId'],
        code = json['code'];

  Map<String, dynamic> toJson() => {
        'userId': this.userId,
        'code': this.code,
      };
}
