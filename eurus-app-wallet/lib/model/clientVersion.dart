class ClientVersion {
  int returnCode;
  String message;
  String nonce;
  ClientVersionData? data;

  ClientVersion(this.returnCode, this.message, this.nonce, this.data);

  ClientVersion.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? new ClientVersionData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class ClientVersionData {
  String iPhoneMinimumVersion;
  String androidMinimumVersion;

  ClientVersionData(this.iPhoneMinimumVersion, this.androidMinimumVersion);

  ClientVersionData.fromJson(Map<String, dynamic> json)
      : iPhoneMinimumVersion = json['iPhoneMinimumVersion'],
        androidMinimumVersion = json['androidMinimumVersion'];

  Map<String, dynamic> toJson() => {
        'iPhoneMinimumVersion': this.iPhoneMinimumVersion,
        'androidMinimumVersion': this.androidMinimumVersion,
      };
}
