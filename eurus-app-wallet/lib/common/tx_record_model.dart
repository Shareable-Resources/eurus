import 'package:app_transaction_manager/data_models/tran_record.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:web3dart/crypto.dart';

class TxRecordModel extends TransactionRecord {
  TxRecordModel({
    required this.transactionHash,
    this.decodedInputAmount,
    this.sendTimestamp,
    this.confirmTimestamp,
  }) : super(transactionHash: transactionHash);

  String transactionHash;
  TransactionInformation? txInfo;
  TransactionReceipt? txReceipt;
  Uint8List? txInput;
  String? decodedInputFncIdentifierHex;
  String? decodedInputRecipientAddress;
  double? decodedInputAmount;
  String? txTo;
  String? txFrom;
  String? sendTimestamp;
  String? confirmTimestamp;
  String? chain;
  String? gasPrice;
  String? gasUsed;
  int? transferGas;
  int? targetGas;
  int? transGasUsed;
  int? eurusTxType;
  int? eurusTxStatus;
  String? adminFee;
  int? transType;
  int? transDate;
  String? txHash;
  int? chainLocation;
  String? fromAddress;
  String? toAddress;
  bool? isSend;
  String? amount;
  int? status;
  String? remarks;
  String? tokenSymbol;
  String? input;
  String? isError;

  bool get dataCompleted =>
      decodedInputRecipientAddress != null &&
      decodedInputAmount != null &&
      txTo != null &&
      txFrom != null &&
      sendTimestamp != null &&
      confirmTimestamp != null &&
      chain != null;

  bool get gasfeeCompleted => gasPrice != null && gasUsed != null;

  bool get isAllocate => eurusTxType != null && eurusTxType != 1;

  @override
  TxRecordModel.fromJson(Map<String, dynamic> r)
      : transactionHash = r['transactionHash'],
        txInfo = jsonToInfo(r['txInfo']),
        txReceipt = jsonToReceipt(r['txReceipt']),
        txInput = r['txInput'],
        decodedInputFncIdentifierHex = r['decodedInputFncIdentifierHex'],
        decodedInputRecipientAddress = r['decodedInputRecipientAddress'],
        decodedInputAmount = r['decodedInputAmount'] != null
            ? double.tryParse(r['decodedInputAmount'])
            : null,
        sendTimestamp = r['sendTimestamp'],
        confirmTimestamp = r['confirmTimestamp'],
        txTo = r['txTo'],
        txFrom = r['txFrom'],
        chain = r['chain'],
        gasPrice = r['gasPrice'].toString(),
        gasUsed = r['gasUsed'].toString(),
        transferGas = r['transferGas'],
        targetGas = r['targetGas'],
        transGasUsed = r['transGasUsed'],
        eurusTxType =
            r['eurusTxType'] != null ? int.tryParse(r['eurusTxType']) : null,
        eurusTxStatus = r['eurusTxStatus'] != null
            ? int.tryParse(r['eurusTxStatus'])
            : null,
        adminFee = r['adminFee'].toString(),
        this.transType = r['transType'],
        this.transDate = r['transDate'],
        this.txHash = r['txHash'],
        this.chainLocation = r['chainLocation'],
        this.fromAddress = r['fromAddress'],
        this.toAddress = r['toAddress'],
        this.isSend = r['isSend'],
        this.amount = r['amount'].toString(),
        this.status = r['status'],
        this.remarks = r['remarks'],
        this.tokenSymbol = r['tokenSymbol'],
        this.input = r['input'],
        this.isError = r['isError'],
        super.fromJson(r);

  @override
  Map<String, dynamic> toJson() => {
        'transactionHash': transactionHash,
        'txInfo': txInfo != null ? jsonEncode(txInfoToMap()) : null,
        'txReceipt': txReceipt != null ? jsonEncode(txReceiptToMap()) : null,
        'txFrom': txFrom ?? txInfo?.from.hex,
        'txTo': txTo ?? txInfo?.to?.hex,
        'txInput': txInput ?? txInfo?.input,
        'decodedInputFncIdentifierHex': decodedInputFncIdentifierHex,
        'decodedInputRecipientAddress': decodedInputRecipientAddress,
        'decodedInputAmount': decodedInputAmount,
        'sendTimestamp': sendTimestamp,
        'confirmTimestamp': confirmTimestamp,
        'chain': chain,
        'gasPrice': gasPrice,
        'gasUsed': gasUsed,
        'transferGas': transferGas,
        'targetGas': targetGas,
        'transGasUsed': transGasUsed,
        'eurusTxType': eurusTxType?.toString(),
        'eurusTxStatus': eurusTxStatus?.toString(),
        'adminFee': adminFee,
        'isSend': isSend,
      };

  @override
  Map<String, dynamic> txInfoToMap() {
    return {
      'blockHash': txInfo!.blockHash,
      'blockNumber': txInfo!.blockNumber.toString(),
      'from': txInfo!.from.hex,
      'gas': txInfo!.gas.toString(),
      'gasPrice': txInfo!.gasPrice.getInWei.toString(),
      'hash': txInfo!.hash,
      'input': bytesToHex(txInfo!.input),
      'nonce': txInfo!.nonce.toString(),
      'to': txInfo!.to?.hex,
      'transactionIndex': txInfo!.transactionIndex.toString(),
      'value': txInfo!.value.getInWei.toString(),
      'v': '0x' + txInfo!.v.toRadixString(16),
      'r': '0x' + txInfo!.r.toRadixString(16),
      's': '0x' + txInfo!.s.toRadixString(16),
    };
  }

  @override
  Map<String, dynamic> txReceiptToMap() {
    return {
      'transactionHash': bytesToHex(txReceipt!.transactionHash),
      'transactionIndex': txReceipt!.transactionIndex.toString(),
      'blockHash': bytesToHex(txReceipt!.blockHash),
      'blockNumber': txReceipt!.blockNumber.toString(),
      'from': txReceipt!.from?.hex,
      'to': txReceipt!.to?.hex,
      'cumulativeGasUsed':
          '0x' + txReceipt!.cumulativeGasUsed.toRadixString(16),
      'gasUsed': txReceipt!.gasUsed != null
          ? '0x' + txReceipt!.gasUsed!.toRadixString(16)
          : null,
      'contractAddress': txReceipt!.contractAddress?.hex,
      'status': txReceipt!.status == true ? '0x1' : '0x0',
      'logs': filterEventToMap(txReceipt!.logs),
      // 'logs': filterEventToMap(txReceipt!.logs).map((e) => jsonEncode(e)).toList(),
      // 'logs': filterEventToMapToJsonString(txReceipt!.logs),
      // 'logs': jsonEncode(filterEventToMapToJsonString(txReceipt!.logs)),
    };
  }

  List<Map<String, dynamic>> filterEventToMap(List<FilterEvent> l) {
    List<Map<String, dynamic>> result = [];

    l.forEach((e) {
      final Map<String, dynamic> log = {
        'removed': e.removed,
        // 'logIndex': '0x' + e.logIndex.toRadixString(16), // ok
        if (e.logIndex != null)
          'logIndex': '0x' + e.logIndex!.toRadixString(16)
        else
          'logIndex': null, // ok
        // 'transactionIndex': '0x' + e.transactionIndex.toRadixString(16), // ok
        ...(e.transactionIndex != null
            ? {'transactionIndex': '0x' + e.transactionIndex!.toRadixString(16)}
            : {'transactionIndex': null}), // ok
        'transactionHash': e.transactionHash,
        'blockHash': e.blockHash,
        // 'blockNum': '0x' + e.blockNum.toRadixString(16), // error
        // 'blockNum': '0x' + e.blockNum?.toRadixString(16), // error
        'blockNum': e.blockNum != null
            ? '0x' + e.blockNum!.toRadixString(16)
            : null, // testing
        'address': e.address?.hex,
        'data': e.data,
        'topics': e.topics,
      };

      result.add(log);
    });

    return result;
  }

  List<String> filterEventToMapToJsonString(List<FilterEvent> l) {
    List<String> result = [];

    l.forEach((e) {
      final Map<String, dynamic> log = {
        'removed': e.removed,
        // 'logIndex': '0x' + e.logIndex.toRadixString(16), // ok
        if (e.logIndex != null)
          'logIndex': '0x' + e.logIndex!.toRadixString(16)
        else
          'logIndex': null, // ok
        // 'transactionIndex': '0x' + e.transactionIndex.toRadixString(16), // ok
        ...(e.transactionIndex != null
            ? {'transactionIndex': '0x' + e.transactionIndex!.toRadixString(16)}
            : {'transactionIndex': null}), // ok
        'transactionHash': e.transactionHash,
        'blockHash': e.blockHash,
        // 'blockNum': '0x' + e.blockNum.toRadixString(16), // error
        // 'blockNum': '0x' + e.blockNum?.toRadixString(16), // error
        'blockNum': e.blockNum != null
            ? '0x' + e.blockNum!.toRadixString(16)
            : null, // testing
        'address': e.address?.hex,
        'data': e.data,
        'topics': e.topics,
      };

      result.add(jsonEncode(log));
    });

    return result;
  }

  static TransactionInformation? jsonToInfo(String? info) {
    if (isEmptyString(string: info)) return null;

    try {
      return TransactionInformation.fromMap(jsonDecode(info!));
    } catch (e, t) {
      print("txrecordmodel: format json to info error $e - $t");
    }

    return null;
  }

  static TransactionReceipt? jsonToReceipt(String? receipt) {
    if (isEmptyString(string: receipt)) return null;

    try {
      return TransactionReceipt.fromMap(jsonDecode(receipt!));
    } catch (e, t) {
      print("txrecordmodel: format json to receipt error $e - $t");
    }

    return null;
  }
}
