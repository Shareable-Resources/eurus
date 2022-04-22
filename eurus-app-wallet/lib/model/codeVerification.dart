class CodeVerification {
  int returnCode;
  String message;
  String nonce;
  CodeVerificationData? data;

  CodeVerification(this.returnCode, this.message, this.nonce, this.data);

  CodeVerification.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? new CodeVerificationData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class CodeVerificationData {
  int? userId;
  String? email;
  String mnemonic;
  String? token;
  int? expiredTime;

  CodeVerificationData(
      this.userId, this.email, this.mnemonic, this.token, this.expiredTime);

  CodeVerificationData.fromJson(Map<String, dynamic> json)
      : userId = json['userId'],
        email = json['email'],
        mnemonic = json['mnemonic'],
        token = json['token'],
        expiredTime = json['expiredTime'];

  Map<String, dynamic> toJson() => {
        'userId': this.userId,
        'email': this.email,
        'mnemonic': this.mnemonic,
        'token': this.token,
        'expiredTime': this.expiredTime,
      };
}
