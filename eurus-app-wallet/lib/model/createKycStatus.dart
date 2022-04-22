class CreateKycStatus {
  String nonce;
  int returnCode;
  String message;
  CreateKycStatusData? data;

  CreateKycStatus(this.nonce, this.returnCode, this.message, this.data);

  CreateKycStatus.fromJson(Map<String, dynamic> json)
      : nonce = json['nonce'],
        returnCode = json['returnCode'],
        message = json['message'],
        data = json['data'] != null
            ? CreateKycStatusData.fromJson(json['data'])
            : null;

  Map<String, dynamic> toJson() => {
        'nonce': this.nonce,
        'returnCode': this.returnCode,
        'message': this.message,
        if (this.data != null) 'data': this.data?.toJson(),
      };
}

class CreateKycStatusData {
  int id;

  CreateKycStatusData(this.id);

  CreateKycStatusData.fromJson(Map<String, dynamic> json) : id = json['id'];

  Map<String, dynamic> toJson() => {
        'id': this.id,
      };
}
