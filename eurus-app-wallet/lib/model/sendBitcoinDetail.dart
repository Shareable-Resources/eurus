class SendBitcoinDetail {
  Tx? tx;
  List<String>? tosign;

  SendBitcoinDetail({this.tx, this.tosign});

  SendBitcoinDetail.fromJson(Map<String, dynamic> json) {
    tx = json['tx'] != null ? new Tx.fromJson(json['tx']) : null;
    tosign = json['tosign'].cast<String>();
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    if (this.tx != null) {
      data['tx'] = this.tx!.toJson();
    }
    data['tosign'] = this.tosign;
    return data;
  }
}

class Tx {
  int? blockHeight;
  int? blockIndex;
  String? hash;
  List<String>? addresses;
  int? total;
  int? fees;
  int? size;
  int? vsize;
  String? preference;
  String? relayedBy;
  String? received;
  int? ver;
  bool? doubleSpend;
  int? vinSz;
  int? voutSz;
  int? confirmations;
  List<InputData>? inputs;
  List<Outputs>? outputs;

  Tx(
      {this.blockHeight,
        this.blockIndex,
        this.hash,
        this.addresses,
        this.total,
        this.fees,
        this.size,
        this.vsize,
        this.preference,
        this.relayedBy,
        this.received,
        this.ver,
        this.doubleSpend,
        this.vinSz,
        this.voutSz,
        this.confirmations,
        this.inputs,
        this.outputs});

  Tx.fromJson(Map<String, dynamic> json) {
    blockHeight = json['block_height'];
    blockIndex = json['block_index'];
    hash = json['hash'];
    addresses = json['addresses'].cast<String>();
    total = json['total'];
    fees = json['fees'];
    size = json['size'];
    vsize = json['vsize'];
    preference = json['preference'];
    relayedBy = json['relayed_by'];
    received = json['received'];
    ver = json['ver'];
    doubleSpend = json['double_spend'];
    vinSz = json['vin_sz'];
    voutSz = json['vout_sz'];
    confirmations = json['confirmations'];
    if (json['inputs'] != null) {
      inputs = [];
      json['inputs'].forEach((v) {
        inputs!.add(new InputData.fromJson(v));
      });
    }
    if (json['outputs'] != null) {
      outputs = [];
      json['outputs'].forEach((v) {
        outputs!.add(new Outputs.fromJson(v));
      });
    }
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['block_height'] = this.blockHeight;
    data['block_index'] = this.blockIndex;
    data['hash'] = this.hash;
    data['addresses'] = this.addresses;
    data['total'] = this.total;
    data['fees'] = this.fees;
    data['size'] = this.size;
    data['vsize'] = this.vsize;
    data['preference'] = this.preference;
    data['relayed_by'] = this.relayedBy;
    data['received'] = this.received;
    data['ver'] = this.ver;
    data['double_spend'] = this.doubleSpend;
    data['vin_sz'] = this.vinSz;
    data['vout_sz'] = this.voutSz;
    data['confirmations'] = this.confirmations;
    if (this.inputs != null) {
      data['inputs'] = this.inputs!.map((v) => v.toJson()).toList();
    }
    if (this.outputs != null) {
      data['outputs'] = this.outputs!.map((v) => v.toJson()).toList();
    }
    return data;
  }
}

class InputData {
  String? prevHash;
  int? outputIndex;
  int? outputValue;
  int? sequence;
  List<String>? addresses;
  String? scriptType;
  int? age;

  InputData(
      {this.prevHash,
        this.outputIndex,
        this.outputValue,
        this.sequence,
        this.addresses,
        this.scriptType,
        this.age});

  InputData.fromJson(Map<String, dynamic> json) {
    prevHash = json['prev_hash'];
    outputIndex = json['output_index'];
    outputValue = json['output_value'];
    sequence = json['sequence'];
    addresses = json['addresses'].cast<String>();
    scriptType = json['script_type'];
    age = json['age'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['prev_hash'] = this.prevHash;
    data['output_index'] = this.outputIndex;
    data['output_value'] = this.outputValue;
    data['sequence'] = this.sequence;
    data['addresses'] = this.addresses;
    data['script_type'] = this.scriptType;
    data['age'] = this.age;
    return data;
  }
}

class Outputs {
  int? value;
  String? script;
  List<String>? addresses;
  String? scriptType;

  Outputs({this.value, this.script, this.addresses, this.scriptType});

  Outputs.fromJson(Map<String, dynamic> json) {
    value = json['value'];
    script = json['script'];
    addresses = json['addresses'].cast<String>();
    scriptType = json['script_type'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['value'] = this.value;
    data['script'] = this.script;
    data['addresses'] = this.addresses;
    data['script_type'] = this.scriptType;
    return data;
  }
}