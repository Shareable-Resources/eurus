class FindEmailWalletAddress {
  int returnCode;
  String message;
  String nonce;
  FindEmailWalletAddressData? data;

  FindEmailWalletAddress(this.returnCode, this.message, this.nonce, this.data);

  FindEmailWalletAddress.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? new FindEmailWalletAddressData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class FindEmailWalletAddressData {
  String? email;
  String? walletAddress;
  int? userType;

  FindEmailWalletAddressData({
    this.email,
    this.walletAddress,
    this.userType,
  });

  FindEmailWalletAddressData.fromJson(Map<String, dynamic> json) {
    email = json['email'];
    walletAddress = json['walletAddress'];
    userType = json['userType'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['email'] = this.email;
    data['walletAddress'] = this.walletAddress;
    data['userType'] = this.userType;
    return data;
  }
}
