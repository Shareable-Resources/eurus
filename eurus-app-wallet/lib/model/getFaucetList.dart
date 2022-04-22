class GetFaucetList {
  int returnCode;
  String message;
  String nonce;
  List<GetFaucetListData> data = [];

  GetFaucetList(this.returnCode, this.message, this.nonce, this.data);

  GetFaucetList.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? json['data']!
                .map<GetFaucetListData>((e) => GetFaucetListData.fromJson(e))
                .toList()
            : [];

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        'data': this.data.map((v) => v.toJson()).toList(),
      };
}

class GetFaucetListData {
  String key;
  int amount;
  int decimal;
  int dateInterval;

  GetFaucetListData(this.key, this.amount, this.decimal, this.dateInterval);

  GetFaucetListData.fromJson(Map<String, dynamic> json)
      : key = json['key'],
        amount = json['amount'],
        decimal = json['decimal'],
        dateInterval = json['dateInterval'];

  Map<String, dynamic> toJson() => {
        'key': this.key,
        'amount': this.amount,
        'decimal': this.decimal,
        'dateInterval': this.dateInterval,
      };
}
