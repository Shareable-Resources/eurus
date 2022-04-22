class SignTransaction {
  String nonce;
  int returnCode;
  String message;
  SignTransactionData? data;

  SignTransaction(this.nonce, this.returnCode, this.message, this.data);

  SignTransaction.fromJson(Map<String, dynamic> json)
      : nonce = json['nonce'],
        returnCode = json['returnCode'],
        message = json['message'],
        data = json['data'] != null
            ? new SignTransactionData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'nonce': this.nonce,
        'returnCode': this.returnCode,
        'message': this.message,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class SignTransactionData {
  String signedTx;

  SignTransactionData(this.signedTx);

  SignTransactionData.fromJson(Map<String, dynamic> json)
      : signedTx = json['signedTx'];

  Map<String, dynamic> toJson() => {
        'signedTx': this.signedTx,
      };
}
