class GetTransactionReceipt {
  String? jsonrpc;
  int? id;
  Result? result;

  GetTransactionReceipt({this.jsonrpc, this.id, this.result});

  GetTransactionReceipt.fromJson(Map<String, dynamic> json) {
    jsonrpc = json['jsonrpc'];
    id = json['id'];
    result =
    json['result'] != null ? new Result.fromJson(json['result']) : null;
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['jsonrpc'] = this.jsonrpc;
    data['id'] = this.id;
    if (this.result != null) {
      data['result'] = this.result?.toJson();
    }
    return data;
  }
}

class Result {
  String? blockHash;
  String? blockNumber;
  String? contractAddress;
  String? cumulativeGasUsed;
  String? from;
  String? gasUsed;
  String? effectiveGasPrice;
  List<Logs>? logs;
  String? logsBloom;
  String? status;
  String? to;
  String? transactionHash;
  String? transactionIndex;
  String? revertReason;

  Result(
      {this.blockHash,
        this.blockNumber,
        this.contractAddress,
        this.cumulativeGasUsed,
        this.from,
        this.gasUsed,
        this.effectiveGasPrice,
        this.logs,
        this.logsBloom,
        this.status,
        this.to,
        this.transactionHash,
        this.transactionIndex,
        this.revertReason});

  Result.fromJson(Map<String, dynamic> json) {
    blockHash = json['blockHash'];
    blockNumber = json['blockNumber'];
    contractAddress = json['contractAddress'];
    cumulativeGasUsed = json['cumulativeGasUsed'];
    from = json['from'];
    gasUsed = json['gasUsed'];
    effectiveGasPrice = json['effectiveGasPrice'];
    if (json['logs'] != null) {
      logs = [];
      json['logs'].forEach((v) {
        logs?.add(new Logs.fromJson(v));
      });
    }
    logsBloom = json['logsBloom'];
    status = json['status'];
    to = json['to'];
    transactionHash = json['transactionHash'];
    transactionIndex = json['transactionIndex'];
    revertReason = json['revertReason'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['blockHash'] = this.blockHash;
    data['blockNumber'] = this.blockNumber;
    data['contractAddress'] = this.contractAddress;
    data['cumulativeGasUsed'] = this.cumulativeGasUsed;
    data['from'] = this.from;
    data['gasUsed'] = this.gasUsed;
    data['effectiveGasPrice'] = this.effectiveGasPrice;
    if (this.logs != null) {
      data['logs'] = this.logs?.map((v) => v.toJson()).toList();
    }
    data['logsBloom'] = this.logsBloom;
    data['status'] = this.status;
    data['to'] = this.to;
    data['transactionHash'] = this.transactionHash;
    data['transactionIndex'] = this.transactionIndex;
    data['revertReason'] = this.revertReason;
    return data;
  }
}

class Logs {
  String? address;
  List<String>? topics;
  String? data;
  String? blockNumber;
  String? transactionHash;
  String? transactionIndex;
  String? blockHash;
  String? logIndex;
  bool? removed;

  Logs(
      {this.address,
        this.topics,
        this.data,
        this.blockNumber,
        this.transactionHash,
        this.transactionIndex,
        this.blockHash,
        this.logIndex,
        this.removed});

  Logs.fromJson(Map<String, dynamic> json) {
    address = json['address'];
    topics = json['topics'].cast<String>();
    data = json['data'];
    blockNumber = json['blockNumber'];
    transactionHash = json['transactionHash'];
    transactionIndex = json['transactionIndex'];
    blockHash = json['blockHash'];
    logIndex = json['logIndex'];
    removed = json['removed'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['address'] = this.address;
    data['topics'] = this.topics;
    data['data'] = this.data;
    data['blockNumber'] = this.blockNumber;
    data['transactionHash'] = this.transactionHash;
    data['transactionIndex'] = this.transactionIndex;
    data['blockHash'] = this.blockHash;
    data['logIndex'] = this.logIndex;
    data['removed'] = this.removed;
    return data;
  }
}