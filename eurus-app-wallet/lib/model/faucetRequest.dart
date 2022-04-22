class FaucetRequest {
  String nonce;
  int returnCode;
  String message;
  FaucetRequestData? data;

  FaucetRequest(this.nonce, this.returnCode, this.message, this.data);

  FaucetRequest.fromJson(Map<String, dynamic> json)
      : nonce = json['nonce'],
        returnCode = json['returnCode'],
        message = json['message'],
        data = json['data'] != null
            ? new FaucetRequestData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'nonce': this.nonce,
        'returnCode': this.returnCode,
        'message': this.message,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class FaucetRequestData {
  String txHash;
  int status;

  FaucetRequestData(this.txHash, this.status);

  FaucetRequestData.fromJson(Map<String, dynamic> json)
      : txHash = json['txHash'],
        status = json['status'];

  Map<String, dynamic> toJson() => {
        'txHash': this.txHash,
        'status': this.status,
      };
}
