import 'package:apihandler/apiHandler.dart';
import 'package:app_transaction_manager/app_transaction_manager.dart';
import 'package:collection/collection.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/common/tx_record_model.dart';
import 'package:web3dart/crypto.dart';

import '../model/crypto_currency_model.dart';

class TxRecordHandler extends AppTransactionManager {
  static String _etherscanAPIKey = 'EC6FFYX2JI8IM8HAFQ9DS1HZZ3YNEB8YEY';

  Map<String, int> _ethTokenDecimal = {};
  Map<String, int> _eunTokenDecimal = {};

  @override
  Future<bool> addSentTx(txRd) async {
    await initDB();

    txRd..sendTimestamp = getCurrentTimeStamp().toString();

    print('setRecord ${txRd.toJson()}');

    int response = await super.db.setRecord(txRd);
    print('setRecord int $response');
    return true;
  }

  /// Read Trasaction Record from Database
  ///
  /// If local record is complete (Have )
  /// return local record
  /// If local record is not complete (missing info / receipt)
  /// Get info using web3 -> update local data -> return data
  @override
  Future<List<TxRecordModel>> readTxs({
    String? where,
    List<String>? whereArgs,
    int? limit,
    int? offset,
    String? order,
  }) async {
    await initDB();

    List<Map<String, dynamic>> records = await super.db.getRecords(
          where: where,
          whereArgs: whereArgs,
          offset: offset,
          order: order ?? 'sendTimestamp DESC',
        );

    List<TxRecordModel> finalRecords = [];

    for (int i = 0; i < records.length; i++) {
      Map<String, dynamic> val = records[i];
      val["transactionHash"] ??= '';
      final rd = TxRecordModel.fromJson(val);
      print('txrecordhandler : ${rd.transactionHash}');

      TxRecordModel? updatedTxRd = await fetchTxRecordContent(rd);

      if (updatedTxRd != null) await updateTx(updatedTxRd);

      finalRecords.add(updatedTxRd ?? rd);
      print(
          'txRdHandler recordList: ${(updatedTxRd ?? rd).transactionHash} :: ${(updatedTxRd ?? rd).decodedInputFncIdentifierHex} :: ${updatedTxRd == null ? 'rd' : 'udrd'}');
    }

    return finalRecords;
  }

  Future<List<TxRecordModel>> readMainnetTxs(
    String address, {
    String? contractAddress,
    bool? isPlatformToken,
    int? limit,
    int? offset,
    String? order,
  }) async {
    return await _fetchNReadTxs(
      BlockChainType.Ethereum,
      address,
      contractAddress: contractAddress,
      isPlatformToken: isPlatformToken,
      limit: limit,
      offset: offset,
      order: order,
    );
  }

  Future<List<TxRecordModel>> readSidechainTxs(
    String address,
    String authToken, {
    String? symbol,
    String? contractAddress,
    bool? isPlatformToken,
    int? limit,
    int? offset,
    String? order,
  }) async {
    String? result = await NormalStorageKit().readValue('apiAccessToken_');

    return await _fetchNReadTxs(
      BlockChainType.Eurus,
      address,
      contractAddress: contractAddress,
      isPlatformToken: isPlatformToken,
      symbol: symbol,
      authToken: result,
      limit: limit,
      offset: offset,
      order: order,
    );
  }

  Future<List<TxRecordModel>> _fetchNReadTxs(
    BlockChainType chain,
    String address, {
    String? symbol,
    String? authToken,
    String? contractAddress,
    bool? isPlatformToken,
    int? limit,
    int? offset,
    String? order,
  }) async {
    await initDB();

    String _prefix = await common.prefix;
    String? stBlockNum = await NormalStorageKit().readValue(
        'last_fetched_${chain}_blocknum_$_prefix${web3dart.myEthereumAddress?.hex}');

    if (chain == BlockChainType.Ethereum) {
      //if erc20 token get decimal first
      int decimal = !isEmptyString(string: contractAddress)
          ? (await web3dart.getContractDecimal(
                      deployedContract: web3dart.getEthereumERC20Contract(
                          contractAddress: contractAddress ?? ''),
                      blockChainType: BlockChainType.Ethereum))
                  ?.toInt() ??
              18
          : 18;
      // return getTxFromMainnet(
      //     useTestnet: true,
      //     stBlock: stBlockNum != null && int.tryParse(stBlockNum) != null
      //         ? int.tryParse(stBlockNum)! + 1
      //         : 0,
      //     contractAddress: contractAddress,
      //     address: web3dart.myEthereumAddress?.hex,
      //     chain: chain,
      //     decimal: decimal,
      //     maxRecord: 50);
      return getTxRecord(
        chain,
        address,
        '',
        common.selectTokenSymbol!,
        contractAddress,
      );
    } else if (chain == BlockChainType.Eurus && authToken != null) {
      String? _symbol = isEmptyString(string: contractAddress)
          ? 'EUN'
          : symbol != null
              ? symbol
              : await getContractSymbol(contractAddress ?? '');

      if (_symbol != null) {
        return getTxRecord(
          chain,
          address,
          '',
          _symbol,
          contractAddress,
        );
      }
    }

    String _chain;
    switch (chain) {
      case BlockChainType.Eurus:
        _chain = 'eun';
        break;
      case BlockChainType.Ethereum:
        _chain = 'eth';
        break;
      case BlockChainType.BinanceCoin:
        _chain = 'bnb';
        break;
    }

    String where =
        '(txFrom LIKE ? OR decodedInputRecipientAddress LIKE ?) AND chain = ? AND txTo LIKE ?';
    List<String> whereArgs = [
      '%$address%',
      '%$address%',
      _chain,
      '%$contractAddress%'
    ];

    if (isPlatformToken == true) {
      where =
          '(txFrom LIKE ? OR decodedInputRecipientAddress LIKE ? OR txTo LIKE ?) AND chain = ? AND txTo LIKE decodedInputRecipientAddress';
      whereArgs = ['%$address%', '%$address%', '%$address%', _chain];
    }

    return await readTxs(
      where: where,
      whereArgs: whereArgs,
      order: order,
      limit: limit,
      offset: offset,
    );
  }

  @override
  Future<bool> updateTx(r) async {
    await super.db.setRecord(r);

    return true;
  }

  /// Get TransactionInformation using web3dart
  Future<TransactionInformation> getTxInfo(String hash, {String? chain}) async {
    Future<TransactionInformation> Function(String) fnc;

    if (chain == 'eun') {
      fnc = web3dart.eurusEthClient.getTransactionByHash;
    } else {
      fnc = web3dart.mainNetEthClient.getTransactionByHash;
    }

    return await fnc(hash).catchError((e, t) {
      /// Log error
      print('txrecordhandler fail: getTxInfo $e - $t');
      return e;
    });
  }

  /// Get TransactionReceipt using web3dart
  Future<TransactionReceipt?> getTxReceipt(String hash, {String? chain}) async {
    Future<TransactionReceipt?> Function(String) fnc;

    if (chain == 'eun') {
      fnc = web3dart.eurusEthClient.getTransactionReceipt;
    } else {
      fnc = web3dart.mainNetEthClient.getTransactionReceipt;
    }

    return await fnc(hash).catchError((e, t) {
      print('txrecordhandler fail: getTxReceipt $e - $t');
      return e;
    });
  }

  /// Get Token decimal by contract address using web3dart
  Future<int?> getTokenDecimal(String contractAddress, {String? chain}) async {
    int? loadedDecimal = chain == 'eun'
        ? _eunTokenDecimal[contractAddress]
        : _ethTokenDecimal[contractAddress];

    if (loadedDecimal != null) return loadedDecimal;

    DeployedContract contract = chain == 'eun'
        ? web3dart.getEurusERC20Contract(contractAddress: contractAddress)
        : web3dart.getEthereumERC20Contract(contractAddress: contractAddress);

    print('chain:$chain');
    BigInt? decimal = await web3dart
        .getContractDecimal(
            deployedContract: contract,
            blockChainType:
                chain == 'eun' ? BlockChainType.Eurus : BlockChainType.Ethereum)
        .catchError(
      (e) {
        print('txrecordhandler fail:  fail in get decimal : $e $chain');
      },
    );

    int? finalDecimal = decimal?.toInt();

    if (finalDecimal != null) {
      chain == 'eun'
          ? _eunTokenDecimal.addAll({contractAddress: finalDecimal})
          : _ethTokenDecimal.addAll({contractAddress: finalDecimal});
    }

    return finalDecimal;
  }

  /// Transform amount to readable value
  double getDecodedInputAmount(BigInt amount, int decimals) {
    return amount / BigInt.from(10).pow(decimals);
  }

  @override
  Map<String, dynamic> decodedInput(Uint8List raw, {int decimals = 0}) {
    if (raw.length < 3) return Map<String, dynamic>();

    String input = bytesToHex(raw);

    String fncIdentifier = '0x' + input.substring(0, 8);

    /// Only if function identifier is transfer 0xa9059cbb and submitwithdraw 0x0269dfb7
    if (fncIdentifier == '0xa9059cbb' || fncIdentifier == '0x0269dfb7') {
      String addressRaw = input.substring(10, 72);
      RegExp exp = new RegExp(r"^0+(.+)$");
      var matches = exp.allMatches(addressRaw);
      if (matches.elementAt(0).group(1) != null) {
        String address = '0x' + matches.elementAt(0).group(1)!;

        String amountRaw = input.substring(72, 136);
        double amountBigInt =
            hexToInt(amountRaw) / (BigInt.from(10).pow(decimals));

        return {
          'fncIdentifier': fncIdentifier,
          'address': address,
          'amount': amountBigInt,
        };
      }
    }

    return Map<String, dynamic>();
  }

  /// Fetch transaction history online
  ///
  /// Transactions in Mainnet (Rinkeby in testing)
  ///   Uses etherscan api
  /// Transactions in Besu
  ///   Uses our api
  Future<List<TxRecordModel>> getTxFromMainnet(
      {String? address,
      required bool useTestnet,
      int? stBlock,
      String? contractAddress,
      required String chain,
      int? decimal,
      required int maxRecord}) async {
    Map<String, TxRecordModel> fetchedRds = {};
    int _stBlock = stBlock != null ? stBlock : 0;

    /// Latest block number that just fetched
    int latestBlockIdx = 0;

    /// During this fetch, select the smallest block that identified as incomplete data as latestBlockIdx
    int smallestIncompleteBlock = -1;

    /// Get transaction from using txlist action
    String txListUri =
        'https://api${useTestnet == true ? '-rinkeby' : ''}.etherscan.io/api?module=account&action=txlist&address=$address&startblock=$_stBlock&sort=desc&apikey=$_etherscanAPIKey';

    Map<String, dynamic>? txListResponse =
        await apiHandler.get(txListUri).then((v) => v).catchError((e, t) {
      print('tx_record_handler-debugg : txList etherscan api : $e - $t');
      return null;
    });

    if (txListResponse != null &&
        txListResponse['status'] == '1' &&
        txListResponse['result'] != null &&
        txListResponse['result'].length > 0) {
      List<dynamic> txListRecords = txListResponse['result'];
      int txListRecordsLenghth =
          txListRecords.length > maxRecord ? maxRecord : txListRecords.length;
      for (var i = 0; i < txListRecordsLenghth; i++) {
        Map<String, dynamic> val = txListRecords[i];
        val["transactionHash"] ??= '';
        // decimal is new function
        // this line is new add
        final value =
            val['value'] is String ? BigInt.tryParse(val['value']) : null;
        double? decodedAmount = value != null && decimal != null
            ? getDecodedInputAmount(value, decimal)
            : null;

        TxRecordModel tempRd = TxRecordModel.fromJson(val);
        tempRd = tempRd
          ..sendTimestamp = '${val['timeStamp']}000'
          ..confirmTimestamp = '${val['timeStamp']}000'
          ..txFrom = val['from']
          ..txTo = val['to']
          ..decodedInputRecipientAddress =
              val['value'] != '0' && val['input'] == '0x' ? val['to'] : null
          ..decodedInputAmount = decodedAmount
          ..chain = 'eth'
          ..gasPrice = val['gasPrice'].toString()
          ..gasUsed = val['gasUsed'].toString()
          ..transactionHash = val['hash'];

        fetchedRds.addAll({val['hash']: tempRd});

        int? blockNum = val['blockNumber'] is String
            ? int.tryParse(val['blockNumber'] as String)
            : null;
        if (blockNum != null) {
          if (blockNum > latestBlockIdx) latestBlockIdx = blockNum;
          if (!tempRd.gasfeeCompleted)
            smallestIncompleteBlock = smallestIncompleteBlock == -1 ||
                    smallestIncompleteBlock > blockNum
                ? blockNum
                : smallestIncompleteBlock;
        }
      }
    }

    /// Get transaction from using tokentx action
    String tokentxUri =
        'https://api${useTestnet == true ? '-rinkeby' : ''}.etherscan.io/api?module=account&action=tokentx&address=$address&startblock=$_stBlock&sort=desc&apikey=$_etherscanAPIKey';

    // when contractAddress not null call api
    Map<String, dynamic>? tokentxResponse = contractAddress != null
        ? await apiHandler.get(tokentxUri).then((v) => v).catchError((e, t) {
            print(
                'tx_record_handler-debugg : tokenlist etherscan api : $e - $t');
            return null;
          })
        : null;

    if (tokentxResponse != null &&
        tokentxResponse['status'] == '1' &&
        tokentxResponse['result'] != null &&
        tokentxResponse['result'].length > 0) {
      List<dynamic> tokenListRecords = tokentxResponse['result'];
      Map<String, int> map = {};
      int tokenListLenghth = tokenListRecords.length > maxRecord
          ? maxRecord
          : tokenListRecords.length;
      for (var i = 0; i < tokenListLenghth; i++) {
        Map<String, dynamic> val = tokenListRecords[i];
        val["transactionHash"] ??= '';
        String contractAddressString = val['contractAddress'];
        bool haveContractAddress = contractAddressString.isNotEmpty &&
            contractAddressString.contains('0x');
        int decimal = 0;

        if (haveContractAddress) {
          if (map[contractAddressString] != null) {
            decimal = map[contractAddressString]!;
          } else {
            final contractDecimal = await web3dart.getContractDecimal(
                deployedContract: web3dart.getEthereumERC20Contract(
                    contractAddress: contractAddressString),
                blockChainType: BlockChainType.Ethereum);
            decimal = contractDecimal != null ? contractDecimal.toInt() : 0;
            map[contractAddressString] = decimal;
          }
        }

        final value =
            val['value'] is String ? BigInt.tryParse(val['value']) : null;
        final tokenDecimal = val['tokenDecimal'] is String
            ? int.tryParse(val['tokenDecimal'])
            : null;
        double? decodedAmount = value != null
            ? haveContractAddress
                ? getDecodedInputAmount(value, decimal)
                : tokenDecimal != null
                    ? getDecodedInputAmount(value, tokenDecimal)
                    : null
            : null;

        TxRecordModel tempRd = TxRecordModel.fromJson(val);
        tempRd = tempRd
          ..sendTimestamp = '${val['timeStamp']}000'
          ..confirmTimestamp = '${val['timeStamp']}000'
          ..txFrom = val['from']
          ..txTo = val['contractAddress']
          ..decodedInputRecipientAddress = val['to']
          ..decodedInputAmount = decodedAmount
          ..chain = 'eth'
          ..gasPrice = val['gasPrice'].toString()
          ..gasUsed = val['gasUsed'].toString()
          ..transactionHash = val['hash'];

        fetchedRds.addAll({val['hash']: tempRd});

        int? blockNum = val['blockNumber'] is String
            ? int.tryParse(val['blockNumber'])
            : null;
        if (blockNum != null) {
          if (blockNum > latestBlockIdx) latestBlockIdx = blockNum;
          if (!tempRd.dataCompleted || !tempRd.gasfeeCompleted)
            smallestIncompleteBlock = smallestIncompleteBlock == -1 ||
                    smallestIncompleteBlock > blockNum
                ? blockNum
                : smallestIncompleteBlock;
        }
      }
    }

    /// Get eth withdraw record using txlistinternal
    String txInternalUri =
        'https://api${useTestnet == true ? '-rinkeby' : ''}.etherscan.io/api?module=account&action=txlistinternal&address=$address&startblock=$_stBlock&sort=desc&apikey=$_etherscanAPIKey';

    Map<String, dynamic>? txInternalResponse =
        await apiHandler.get(txInternalUri).then((v) => v).catchError((e, t) {
      print(
          'tx_record_handler-debugg : txlistinternal etherscan api : $e - $t');
      return null;
    });

    print('txrecordhandler kk: $txInternalResponse');

    if (txInternalResponse != null &&
        txInternalResponse['status'] == '1' &&
        txInternalResponse['result'] != null &&
        txInternalResponse['result'].length > 0) {
      List<dynamic> tokenListRecords = txInternalResponse['result'];
      int tokenListRecordsLenghth = tokenListRecords.length > maxRecord
          ? maxRecord
          : tokenListRecords.length;
      for (var i = 0; i < tokenListRecordsLenghth; i++) {
        Map<String, dynamic> val = tokenListRecords[i];
        val["transactionHash"] ??= '';
        final value =
            val['value'] is String ? BigInt.tryParse(val['value']) : null;
        double? decodedAmount =
            value != null ? getDecodedInputAmount(value, 18) : null;

        TxRecordModel tempRd = TxRecordModel.fromJson(val);
        tempRd = tempRd
          ..sendTimestamp = '${val['timeStamp']}000'
          ..confirmTimestamp = '${val['timeStamp']}000'
          ..txFrom = val['from']
          ..txTo = val['to']
          ..decodedInputRecipientAddress = val['to']
          ..decodedInputAmount = decodedAmount
          ..chain = 'eth'
          ..gasPrice = val['gasPrice'].toString()
          ..gasUsed = val['gasUsed'].toString()
          ..transactionHash = val['hash'];

        fetchedRds.addAll({val['hash']: tempRd});

        int? blockNum = val['blockNumber'] is String
            ? int.tryParse(val['blockNumber'])
            : null;
        if (blockNum != null) {
          if (blockNum > latestBlockIdx) latestBlockIdx = blockNum;
          if (!tempRd.dataCompleted || !tempRd.gasfeeCompleted)
            smallestIncompleteBlock = smallestIncompleteBlock == -1 ||
                    smallestIncompleteBlock > blockNum
                ? blockNum
                : smallestIncompleteBlock;
        }
      }
    }

    List<TxRecordModel> rdToSave = fetchedRds.values.toList();
    List<TxRecordModel> rdToSave2 = [];
    rdToSave.forEach((TxRecordModel txRecordModel) {
      print("txRecordModel$txRecordModel");
      if (address != null &&
          (txRecordModel.txFrom != null &&
                  txRecordModel.txFrom!.contains(address) ||
              txRecordModel.decodedInputRecipientAddress != null &&
                  txRecordModel.decodedInputRecipientAddress!
                      .contains(address))) {
        //eth case
        if ((isEmptyString(string: contractAddress)) &&
            (txRecordModel.chain != null &&
                txRecordModel.chain!.contains(chain)) &&
            isEmptyString(string: txRecordModel.tokenSymbol)) {
          //input == "" or input == "0x" case or input == null
          if ((txRecordModel.input != null &&
              txRecordModel.input!.length <= 2)) {
            if (txRecordModel.isError != "1") {
              rdToSave2.add(txRecordModel);
            }
          }
        } else if (contractAddress != null &&
            (txRecordModel.txTo != null &&
                txRecordModel.txTo!.contains(contractAddress)) &&
            (txRecordModel.chain != null &&
                txRecordModel.chain!.contains(chain))) {
          //erc 20 token case
          if (txRecordModel.isError != "1") {
            rdToSave2.add(txRecordModel);
          }
        }
      }
    });

    rdToSave2.sort((TxRecordModel a, TxRecordModel b) {
      var adate = a.confirmTimestamp;
      var bdate = b.confirmTimestamp;

      return adate == null
          ? 1
          : bdate == null
              ? -1
              : bdate.compareTo(
                  adate); //to get the order other way just switch `adate & bdate`
    });
    print("rdToSave:$rdToSave");
    print("rdToSave2:$rdToSave2");
    return rdToSave2;
  }

  Future<List<TxRecordModel>> getTxRecord(
    BlockChainType chain,
    String address,
    String nonce,
    String symbol,
    String? symbolContract,
  ) async {
    Map<String, dynamic> adminfeeDetail =
        await api.withdrawAdminFeeDetail(symbol: symbol);
    Map<String, dynamic> adminfeeDetailData =
        adminfeeDetail['data'] is Map<String, dynamic>
            ? adminfeeDetail['data']
            : Map<String, dynamic>();
    int? adminfeeDecimal = adminfeeDetailData['decimal'] is int
        ? adminfeeDetailData['decimal']
        : null;

    Map<String, dynamic> response = await api.recentTransaction(
      chain: chain,
      nonce: nonce,
      symbol: symbol,
    );
    List<TxRecordModel> fetchedRds = [];
    int? returnCode =
        response['returnCode'] is int ? response['returnCode'] : null;
    Map<String, dynamic> recentTransactionData =
        response['data'] is Map<String, dynamic>
            ? response['data']
            : Map<String, dynamic>();
    int? decimalPlace = recentTransactionData['decimalPlace'] is int
        ? recentTransactionData['decimalPlace']
        : null;
    List txsList = recentTransactionData['transList'] is List
        ? recentTransactionData['transList']
        : [];
    final eurusUserDepositAddress = await web3dart.getEurusUserDepositAddress();

    if (returnCode == 0) {
      for (var i = 0; i < txsList.length; i++) {
        Map<String, dynamic> tempTx = txsList[i];
        String? txHash = tempTx['txHash'] is String ? tempTx['txHash'] : null;
        tempTx["transactionHash"] ??= '';
        int? transType =
            tempTx['transType'] is int ? tempTx['transType'] : null;
        int? status = tempTx['status'] is int ? tempTx['status'] : null;
        BigInt? amount = tempTx['amount'] is String
            ? BigInt.tryParse(tempTx['amount'])
            : tempTx['amount'] is num
                ? BigInt.from(tempTx['amount'])
                : null;
        int? adminFee = tempTx['adminFee'] is int ? tempTx['adminFee'] : null;
        String? transDate =
            tempTx['transDate'] is String ? tempTx['transDate'] : null;
        String? fromAddress =
            tempTx['fromAddress'] is String ? tempTx['fromAddress'] : null;
        String? toAddress =
            tempTx['toAddress'] is String ? tempTx['toAddress'] : null;
        String? destAddress =
            tempTx['destAddress'] is String ? tempTx['destAddress'] : null;
        String? targetAddress =
            tempTx['targetAddress'] is String ? tempTx['targetAddress'] : null;
        int? chainLocation =
            tempTx['chainLocation'] is int ? tempTx['chainLocation'] : null;
        String? depositTxHash =
            tempTx['depositTxHash'] is String ? tempTx['depositTxHash'] : null;
        double? decodedInputAmount = amount != null && decimalPlace != null
            ? getDecodedInputAmount(amount, decimalPlace)
            : null;
        String? withdrawTxHash = tempTx['withdrawTxHash'] is String
            ? tempTx['withdrawTxHash']
            : null;

        print('recentTransactionResponse::: $tempTx');
        print(
            'response-detail::: $txHash :: $transType :: $status :: $amount :: $adminFee');

        final _chain = chain == BlockChainType.Eurus
            ? 'eun'
            : chain == BlockChainType.Ethereum
                ? 'eth'
                : chain == BlockChainType.BinanceCoin
                    ? 'bnb'
                    : null;
        if (transType == 1 && status == 50) {
          TxRecordModel tempRd = TxRecordModel.fromJson(tempTx);
          tempRd = tempRd
            ..sendTimestamp = '${tempRd.transDate}000'
            ..confirmTimestamp = '${tempRd.transDate}000'
            ..chain = _chain
            ..eurusTxType = transType
            ..txTo = symbolContract
            ..txFrom = fromAddress
            ..decodedInputRecipientAddress = toAddress
            ..decodedInputAmount = decodedInputAmount
            ..eurusTxStatus = status
            ..transactionHash = txHash ?? '';

          //TxRecordModel finalRd = await fetchTxRecordContent(tempRd) ?? tempRd;

          if (tempRd.toAddress != eurusUserDepositAddress)
            fetchedRds.add(tempRd);
        } else if (transType == 2 && status == 40) {
          TxRecordModel tempRd = TxRecordModel.fromJson(tempTx);
          tempRd = tempRd
            ..sendTimestamp = '${tempRd.transDate}000'
            ..confirmTimestamp = '${tempRd.transDate}000'
            ..txTo = symbolContract
            ..txFrom = toAddress
            ..decodedInputRecipientAddress = destAddress
            ..decodedInputAmount = decodedInputAmount
            ..chain = _chain
            ..eurusTxType = transType
            ..eurusTxStatus = status
            ..transactionHash = depositTxHash ?? txHash ?? '';

          fetchedRds.add(tempRd);
        } else if (transType == 3 && status == 80) {
          TxRecordModel tempRd = TxRecordModel.fromJson(tempTx);
          tempRd = tempRd
            ..sendTimestamp = '${tempRd.transDate}000'
            ..confirmTimestamp = '${tempRd.transDate}000'
            ..txTo = symbolContract
            ..txFrom = fromAddress
            ..decodedInputRecipientAddress = targetAddress
            ..decodedInputAmount = decodedInputAmount
            ..chain = _chain
            ..eurusTxType = transType
            ..eurusTxStatus = status
            ..adminFee = adminfeeDecimal != null && adminFee != null
                ? (adminFee / pow(10, adminfeeDecimal)).toString()
                : null
            ..transactionHash = withdrawTxHash ?? txHash ?? '';

          fetchedRds.add(tempRd);
        } else if (transType == 6) {
          TxRecordModel tempRd = TxRecordModel.fromJson(tempTx);
          tempRd = tempRd
            ..sendTimestamp = '${tempRd.transDate}000'
            ..confirmTimestamp = '${tempRd.transDate}000'
            ..txTo = symbolContract
            ..txFrom = fromAddress
            ..decodedInputRecipientAddress = targetAddress
            ..decodedInputAmount = decodedInputAmount
            ..chain = _chain
            ..eurusTxType = transType
            ..eurusTxStatus = status
            ..adminFee = adminfeeDecimal != null && adminFee != null
                ? (adminFee / pow(10, adminfeeDecimal)).toString()
                : null
            ..transactionHash = withdrawTxHash ?? txHash ?? '';

          fetchedRds.add(tempRd);
        } else if (transType == 7 && status == 1) {
          TxRecordModel tempRd = TxRecordModel.fromJson(tempTx);
          tempRd = tempRd
            ..sendTimestamp = '${tempRd.transDate}000'
            ..confirmTimestamp = '${tempRd.transDate}000'
            ..txTo = symbolContract
            ..txFrom = fromAddress
            ..decodedInputRecipientAddress = targetAddress
            ..decodedInputAmount = decodedInputAmount
            ..chain = _chain
            ..eurusTxType = transType
            ..eurusTxStatus = status
            ..adminFee = adminfeeDecimal != null && adminFee != null
                ? (adminFee / pow(10, adminfeeDecimal)).toString()
                : null
            ..transactionHash = withdrawTxHash ?? txHash ?? '';
          fetchedRds.add(tempRd);
        }
      }
      fetchedRds.sort((TxRecordModel a, TxRecordModel b) {
        var adate = a.confirmTimestamp;
        var bdate = b.confirmTimestamp;

        return adate == null
            ? 1
            : bdate == null
                ? -1
                : bdate.compareTo(
                    adate); //to get the order other way just switch `adate & bdate`
      });
      print("fetchedRds");
      return fetchedRds;
    }
    return [];
  }

  Future<TxRecordModel?> fetchTxRecordContent(TxRecordModel rd) async {
    bool updated = false;

    if (rd.dataCompleted &&
        (rd.gasfeeCompleted || (rd.eurusTxType != 2 && rd.eurusTxType != 3)))
      return null;

    /// Update txInfo if found null or missing blockhash
    if (rd.txInfo == null || rd.txInfo!.blockHash == null) {
      updated = true;

      await getTxInfo(rd.transactionHash, chain: rd.chain)
          .then((value) => rd.txInfo = value)
          .catchError((e) {});

      if (rd.txInfo != null) {
        final txInfo = rd.txInfo!;
        rd.txFrom = rd.txFrom ?? txInfo.from.hex;
        rd.txTo = rd.txTo ?? txInfo.to?.hex;
        rd.txInput = rd.txInput ?? txInfo.input;
        rd.gasPrice = rd.gasPrice ?? txInfo.gasPrice.getInWei.toRadixString(10);

        if (txInfo.input.length > 2) {
          int? decimal = rd.decodedInputAmount != null || txInfo.to?.hex == null
              ? 0
              : await getTokenDecimal(txInfo.to!.hex, chain: rd.chain);

          Map<String, dynamic>? decodedVal = rd.txInfo != null
              ? decodedInput(txInfo.input, decimals: decimal ?? 0)
              : null;

          if (decodedVal != null) {
            rd.decodedInputFncIdentifierHex =
                rd.decodedInputFncIdentifierHex ?? decodedVal['fncIdentifier'];
            rd.decodedInputRecipientAddress =
                rd.decodedInputRecipientAddress ?? decodedVal['address'];
            rd.decodedInputAmount =
                rd.decodedInputAmount ?? decodedVal['amount'];
          }
        } else {
          rd.decodedInputAmount =
              getDecodedInputAmount(txInfo.value.getInWei, 18);
          rd.decodedInputRecipientAddress =
              rd.decodedInputRecipientAddress ?? txInfo.to?.hex;
        }
      }
    }

    if (rd.txReceipt == null && rd.txInfo?.blockHash != null) {
      updated = true;
      await getTxReceipt(rd.transactionHash, chain: rd.chain)
          .then((value) => rd.txReceipt = value)
          .catchError((e) {});

      if (rd.txReceipt != null) {
        rd.confirmTimestamp =
            rd.confirmTimestamp ?? getCurrentTimeStamp().toString();

        rd.gasUsed = rd.txReceipt?.gasUsed?.toRadixString(10);
      }
    }

    return updated ? rd : null;
  }

  Future<String?> getContractSymbol(String address) async {
    Map<String, CryptoCurrencyModel> supportedList =
        await CommonMethod().getSupportedTokens().catchError((e) {
      print('txrecordhandler : getContractSymbolFailed $e');
      return Map<String, CryptoCurrencyModel>();
    });

    final symbol = supportedList.keys.toList().firstWhereOrNull(
        (element) => supportedList[element]?.addressEurus == address);

    return symbol;
  }
}
