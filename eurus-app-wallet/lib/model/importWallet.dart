class ImportWallet {
  int returnCode;
  String message;
  String nonce;
  ImportWalletData? data;
  String? token;
  int? expiryTime;

  ImportWallet(this.returnCode, this.message, this.nonce, this.data, this.token,
      this.expiryTime);

  ImportWallet.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        token = json['token'],
        expiryTime = json['expiryTime'],
        data = json['data'] != null
            ? new ImportWalletData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class ImportWalletData {
  String token;
  int expiryTime;
  String lastLoginTime;
  int status;
  String txHash;

  ImportWalletData(this.token, this.expiryTime, this.lastLoginTime, this.status,
      this.txHash);

  ImportWalletData.fromJson(Map<String, dynamic> json)
      : token = json['token'],
        expiryTime = json['expiryTime'],
        lastLoginTime = json['lastLoginTime'],
        status = json['status'],
        txHash = json['txHash'];

  Map<String, dynamic> toJson() => {
        'token': this.token,
        'expiryTime': this.expiryTime,
        'lastLoginTime': this.lastLoginTime,
        'status': this.status,
        'txHash': this.txHash,
      };
}
