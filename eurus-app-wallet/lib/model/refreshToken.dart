class RefreshToken {
  int returnCode;
  String message;
  String nonce;
  RefreshTokenData? data;

  RefreshToken(this.returnCode, this.message, this.nonce, this.data);

  RefreshToken.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? new RefreshTokenData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class RefreshTokenData {
  String token;
  int expiryTime;

  RefreshTokenData(this.token, this.expiryTime);

  RefreshTokenData.fromJson(Map<String, dynamic> json)
      : token = json['token'],
        expiryTime = json['expiryTime'];

  Map<String, dynamic> toJson() => {
        'token': this.token,
        'expiryTime': this.expiryTime,
      };
}
