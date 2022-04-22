class AdminFeeModel {
  final int returnCode;
  final String message;
  final String nonce;
  final AdminFeeModelData? data;

  AdminFeeModel(this.returnCode, this.message, this.nonce, this.data);

  AdminFeeModel.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? new AdminFeeModelData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class AdminFeeModelData {
  final String currency;
  final int fee;
  final int decimal;
  double? actualFee;

  AdminFeeModelData(this.currency, this.fee, this.decimal, this.actualFee);

  AdminFeeModelData.fromJson(Map<String, dynamic> json)
      : currency = json['currency'],
        fee = json['fee'],
        decimal = json['decimal'],
        actualFee = json['actualFee'];

  Map<String, dynamic> toJson() => {
        'currency': this.currency,
        'fee': this.fee,
        'decimal': this.decimal,
        'actualFee': this.actualFee,
      };
}
