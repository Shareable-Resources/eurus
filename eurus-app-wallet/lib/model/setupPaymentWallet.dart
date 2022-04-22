class SetupPaymentWallet {
  int returnCode;
  String message;
  String nonce;
  SetupPaymentWalletData? data;

  SetupPaymentWallet(this.returnCode, this.message, this.nonce, this.data);

  SetupPaymentWallet.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? new SetupPaymentWalletData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class SetupPaymentWalletData {
  String token;
  String walletAddress;
  String mainnetWalletAddress;
  bool isMetamaskAddr;

  SetupPaymentWalletData(this.token, this.walletAddress,
      this.mainnetWalletAddress, this.isMetamaskAddr);

  SetupPaymentWalletData.fromJson(Map<String, dynamic> json)
      : token = json['token'],
        walletAddress = json['walletAddress'],
        mainnetWalletAddress = json['mainnetWalletAddress'],
        isMetamaskAddr = json['isMetamaskAddr'];

  Map<String, dynamic> toJson() => {
        'token': this.token,
        'walletAddress': this.walletAddress,
        'mainnetWalletAddress': this.mainnetWalletAddress,
        'isMetamaskAddr': this.isMetamaskAddr,
      };
}
