class BlockchainAddressInformation {
  String? address;
  int? totalReceived;
  int?totalSent;
  int? balance;
  int? unconfirmedBalance;
  int? finalBalance;
  int? nTx;
  int? unconfirmedNTx;
  int? finalNTx;
  List<Txrefs>? txrefs;
  String? txUrl;

  BlockchainAddressInformation(
      {this.address,
        this.totalReceived,
        this.totalSent,
        this.balance,
        this.unconfirmedBalance,
        this.finalBalance,
        this.nTx,
        this.unconfirmedNTx,
        this.finalNTx,
        this.txrefs,
        this.txUrl});

  BlockchainAddressInformation.fromJson(Map<String, dynamic> json) {
    address = json['address'];
    totalReceived = json['total_received'];
    totalSent = json['total_sent'];
    balance = json['balance'];
    unconfirmedBalance = json['unconfirmed_balance'];
    finalBalance = json['final_balance'];
    nTx = json['n_tx'];
    unconfirmedNTx = json['unconfirmed_n_tx'];
    finalNTx = json['final_n_tx'];
    if (json['txrefs'] != null) {
      txrefs = [];
      json['txrefs'].forEach((v) {
        txrefs!.add(new Txrefs.fromJson(v));
      });
    }
    txUrl = json['tx_url'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['address'] = this.address;
    data['total_received'] = this.totalReceived;
    data['total_sent'] = this.totalSent;
    data['balance'] = this.balance;
    data['unconfirmed_balance'] = this.unconfirmedBalance;
    data['final_balance'] = this.finalBalance;
    data['n_tx'] = this.nTx;
    data['unconfirmed_n_tx'] = this.unconfirmedNTx;
    data['final_n_tx'] = this.finalNTx;
    if (this.txrefs != null) {
      data['txrefs'] = this.txrefs!.map((v) => v.toJson()).toList();
    }
    data['tx_url'] = this.txUrl;
    return data;
  }
}

class Txrefs {
  String? txHash;
  int? blockHeight;
  int? txInputN;
  int? txOutputN;
  int? value;
  int? refBalance;
  int? confirmations;
  String? confirmed;
  bool? doubleSpend;
  bool? spent;
  String? spentBy;

  Txrefs(
      {this.txHash,
        this.blockHeight,
        this.txInputN,
        this.txOutputN,
        this.value,
        this.refBalance,
        this.confirmations,
        this.confirmed,
        this.doubleSpend,
        this.spent,
        this.spentBy});

  Txrefs.fromJson(Map<String, dynamic> json) {
    txHash = json['tx_hash'];
    blockHeight = json['block_height'];
    txInputN = json['tx_input_n'];
    txOutputN = json['tx_output_n'];
    value = json['value'];
    refBalance = json['ref_balance'];
    confirmations = json['confirmations'];
    confirmed = json['confirmed'];
    doubleSpend = json['double_spend'];
    spent = json['spent'];
    spentBy = json['spent_by'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['tx_hash'] = this.txHash;
    data['block_height'] = this.blockHeight;
    data['tx_input_n'] = this.txInputN;
    data['tx_output_n'] = this.txOutputN;
    data['value'] = this.value;
    data['ref_balance'] = this.refBalance;
    data['confirmations'] = this.confirmations;
    data['confirmed'] = this.confirmed;
    data['double_spend'] = this.doubleSpend;
    data['spent'] = this.spent;
    data['spent_by'] = this.spentBy;
    return data;
  }
}