class ServerConfig {
  int returnCode;
  String message;
  String nonce;
  ServerConfigData? data;

  ServerConfig(this.returnCode, this.message, this.nonce, this.data);

  ServerConfig.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? new ServerConfigData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class ServerConfigData {
  String eurusRPCDomain;
  int eurusRPCPort;
  String eurusPRCProtocol;
  int eurusChainId;
  String mainnetRPCDomain;
  int mainnetRPCPort;
  String mainnetRPCProtocol;
  int mainnetChainId;
  String externalSmartContractConfigAddress;
  String eurusInternalConfigAddress;

  ServerConfigData(
      this.eurusRPCDomain,
      this.eurusRPCPort,
      this.eurusPRCProtocol,
      this.eurusChainId,
      this.mainnetRPCDomain,
      this.mainnetRPCPort,
      this.mainnetRPCProtocol,
      this.mainnetChainId,
      this.externalSmartContractConfigAddress,
      this.eurusInternalConfigAddress);

  ServerConfigData.fromJson(Map<String, dynamic> json)
      : eurusRPCDomain = json['eurusRPCDomain'],
        eurusRPCPort = json['eurusRPCPort'],
        eurusPRCProtocol = json['eurusPRCProtocol'],
        eurusChainId = json['eurusChainId'],
        mainnetRPCDomain = json['mainnetRPCDomain'],
        mainnetRPCPort = json['mainnetRPCPort'],
        mainnetRPCProtocol = json['mainnetRPCProtocol'],
        mainnetChainId = json['mainnetChainId'],
        externalSmartContractConfigAddress =
            json['externalSmartContractConfigAddress'],
        eurusInternalConfigAddress = json['eurusInternalConfigAddress'];

  Map<String, dynamic> toJson() => {
        'eurusRPCDomain': this.eurusRPCDomain,
        'eurusRPCPort': this.eurusRPCPort,
        'eurusPRCProtocol': this.eurusPRCProtocol,
        'eurusChainId': this.eurusChainId,
        'mainnetRPCDomain': this.mainnetRPCDomain,
        'mainnetRPCPort': this.mainnetRPCPort,
        'mainnetRPCProtocol': this.mainnetRPCProtocol,
        'mainnetChainId': this.mainnetChainId,
        'externalSmartContractConfigAddress':
            this.externalSmartContractConfigAddress,
        'eurusInternalConfigAddress': this.eurusInternalConfigAddress,
      };
}
