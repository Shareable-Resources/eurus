class ForgetLoginPw {
  int returnCode;
  String message;
  String nonce;
  ForgetLoginPwData? data;

  ForgetLoginPw(this.returnCode, this.message, this.nonce, this.data);

  ForgetLoginPw.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? new ForgetLoginPwData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class ForgetLoginPwData {
  int userId;
  String code;
  int type;
  String createdDate;
  String lastModifiedDate;
  String expiredTime;
  int count;

  ForgetLoginPwData(this.userId, this.code, this.type, this.createdDate,
      this.lastModifiedDate, this.expiredTime, this.count);

  ForgetLoginPwData.fromJson(Map<String, dynamic> json)
      : userId = json['UserId'],
        code = json['Code'],
        type = json['Type'],
        createdDate = json['CreatedDate'],
        lastModifiedDate = json['LastModifiedDate'],
        expiredTime = json['ExpiredTime'],
        count = json['Count'];

  Map<String, dynamic> toJson() => {
        'UserId': this.userId,
        'Code': this.code,
        'Type': this.type,
        'CreatedDate': this.createdDate,
        'LastModifiedDate': this.lastModifiedDate,
        'ExpiredTime': this.expiredTime,
        'Count': this.count,
      };
}
