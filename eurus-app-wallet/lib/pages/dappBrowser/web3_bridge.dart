import 'package:collection/collection.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:eth_sig_util/eth_sig_util.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/common/web3dartTransaction.dart';
import 'package:euruswallet/model/crypto_currency_model.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:web3dart/crypto.dart';
import 'package:web3dart/json_rpc.dart';
import 'package:webview_flutter/webview_flutter.dart';

import 'dapp_browser_confirmation_dialog.dart';

enum JsonRpcApiMethod {
  net_version,
  eth_accounts,
  eth_blockNumber,
  eth_call,
  eth_chainId,
  eth_estimateGas,
  eth_getTransactionByHash,
  eth_getTransactionReceipt,
  eth_requestAccounts,
  eth_sendTransaction,
  eth_sign,
  personal_sign,
  personal_ecRecover,
  eth_signTypedData,
  eth_signTypedData_v3,
  eth_signTypedData_v4,
  wallet_watchAsset,
  wallet_switchEthereumChain,
  wallet_addEthereumChain
}

extension JsonRpcApiMethodExtension on JsonRpcApiMethod {
  String get displayValue => describeEnum(this);
}

enum NetworkType {
  ethMainNet,
  ethRinkebyTestNet,
  eurusMainNet,
  eurusTestNet,
  eurusDev
}

extension NetworkTypeExtension on NetworkType {
  bool get isEnabled {
    if (isCentralized() &&
        (this == NetworkType.ethMainNet ||
            this == NetworkType.ethRinkebyTestNet)) return false;

    switch (this) {
      case NetworkType.ethMainNet:
      case NetworkType.eurusMainNet:
        return envType == EnvType.Production;
      case NetworkType.ethRinkebyTestNet:
        return envType == EnvType.Testnet || envType == EnvType.Dev;
      case NetworkType.eurusTestNet:
        return envType == EnvType.Testnet;
      case NetworkType.eurusDev:
        return envType == EnvType.Dev;
    }
  }

  bool get isEurusNetwork {
    return this == NetworkType.eurusDev ||
        this == NetworkType.eurusTestNet ||
        this == NetworkType.eurusMainNet;
  }

  int get netVersion {
    switch (this) {
      case NetworkType.ethMainNet:
      case NetworkType.ethRinkebyTestNet:
        return web3dart.mainNetChainId;
      case NetworkType.eurusMainNet:
      case NetworkType.eurusTestNet:
      case NetworkType.eurusDev:
        return web3dart.eurusChainId;
    }
  }

  String get chainId {
    return '0x' + this.netVersion.toRadixString(16);
  }

  String get displayValue {
    switch (this) {
      case NetworkType.ethMainNet:
        return 'Ethereum Mainnet';
      case NetworkType.ethRinkebyTestNet:
        return 'Rinkeby Test Network';
      case NetworkType.eurusMainNet:
        return 'Eurus Mainnet';
      case NetworkType.eurusTestNet:
        return 'Eurus Test Network';
      case NetworkType.eurusDev:
        return 'Eurus Dev Network';
    }
  }
}

class Web3Bridge {
  InAppWebViewController? webViewController;
  late BuildContext context;
  NetworkType currentNetworkType = envType == EnvType.Dev
      ? NetworkType.eurusDev
      : envType == EnvType.Testnet
          ? NetworkType.eurusTestNet
          : NetworkType.eurusMainNet;
  bool popped = false;
  bool isAccountConfirmed = false;
  Web3Client get currentClient {
    return web3dart.getCurrentClient(blockChainType: blockChainType);
  }

  BlockChainType get blockChainType {
    return currentNetworkType.isEurusNetwork
        ? BlockChainType.Eurus
        : BlockChainType.Ethereum;
  }

  Web3Bridge._internal();

  static final Web3Bridge instance = Web3Bridge._internal();

  Web3Bridge();

  void setNetwork(NetworkType type) async {
    if (currentNetworkType == type || !type.isEnabled) return;

    currentNetworkType = type;
    await webViewController?.reload();
  }

  Future<dynamic> request(String method, List? params) async {
    final methodType = JsonRpcApiMethod.values
        .firstWhere((element) => element.displayValue == method);
    final paramsList = params?.first is List ? params?.first as List : [];

    switch (methodType) {
      case JsonRpcApiMethod.net_version:
        return currentNetworkType.netVersion;
      case JsonRpcApiMethod.eth_accounts:
        return getAccounts();
      case JsonRpcApiMethod.eth_blockNumber:
        return currentClient.getBlockNumber();
      case JsonRpcApiMethod.eth_call:
        return ethCall(paramsList);
      case JsonRpcApiMethod.eth_chainId:
        return currentNetworkType.chainId;
      case JsonRpcApiMethod.eth_estimateGas:
        return ethEstimateGas(paramsList);
      case JsonRpcApiMethod.eth_getTransactionByHash:
        return ethGetTransactionByHash(paramsList);
      case JsonRpcApiMethod.eth_getTransactionReceipt:
        return ethGetTransactionReceipt(paramsList);
      case JsonRpcApiMethod.eth_requestAccounts:
        return ethRequestAccounts(paramsList);
      case JsonRpcApiMethod.eth_sendTransaction:
        return ethSendTransaction(paramsList);
      case JsonRpcApiMethod.eth_sign:
        return ethSign(paramsList);
      case JsonRpcApiMethod.personal_sign:
        return personalSign(paramsList);
      case JsonRpcApiMethod.eth_signTypedData:
        return ethSignTypedData(
          paramsList,
          version: TypedDataVersion.V1,
        );
      case JsonRpcApiMethod.eth_signTypedData_v3:
        return ethSignTypedData(
          paramsList,
          version: TypedDataVersion.V3,
        );
      case JsonRpcApiMethod.eth_signTypedData_v4:
        return ethSignTypedData(
          paramsList,
          version: TypedDataVersion.V4,
        );
      case JsonRpcApiMethod.personal_ecRecover:
        return personalEcRecover(paramsList);
      case JsonRpcApiMethod.wallet_watchAsset:
        // TODO: Handle this case.
        break;
      case JsonRpcApiMethod.wallet_switchEthereumChain:
        return switchEthereumChain(paramsList);
      case JsonRpcApiMethod.wallet_addEthereumChain:
        return switchEthereumChain(paramsList);
    }
  }

  List<String> getAccounts() {
    if (isAccountConfirmed && common.currentAddress is String)
      return [common.currentAddress as String];

    return [];
  }

  Future<List<String>> ethRequestAccounts(List paramsList) async {
    if (isAccountConfirmed) return getAccounts();

    List<String> result = [];

    if (!popped) {
      popped = true;
      final shouldGetAccounts = await showDialog(
        context: context,
        builder: (BuildContext dialogContext) {
          return AlertDialog(
            title: Text('DAPP_BROWSER.CONNECT_SITE_DIALOG_TITLE'.tr()),
            actions: <Widget>[
              TextButton(
                child: Text("COMMON.OK".tr()),
                onPressed: () {
                  Navigator.pop(dialogContext, true);
                },
              ),
              TextButton(
                child: Text("COMMON.CANCEL".tr()),
                onPressed: () {
                  Navigator.pop(dialogContext, false);
                },
              )
            ],
          );
        },
      );
      popped = false;

      if (shouldGetAccounts) {
        isAccountConfirmed = true;
        result = getAccounts();
      } else {
        throw RPCError(
          4001,
          'User denied account authorization.',
          null,
        );
      }
    }
    return result;
  }

  Future<String?> ethCall(List paramsList) async {
    final params =
        paramsList.firstWhere((element) => element is Map<String, dynamic>);
    final to = params["to"] is String ? params["to"] as String : null;
    final data = params["data"] is String
        ? hexToBytes((params["data"] as String).replaceFirst("0x", ''))
        : null;

    if (to == null || data == null) {
      return null;
    }

    return currentClient.callRaw(
      contract: EthereumAddress.fromHex(to),
      data: data,
    );
  }

  Future<String?> ethEstimateGas(List paramsList) async {
    final params =
        paramsList.firstWhere((element) => element is Map<String, dynamic>);
    final from = params["from"] is String
        ? EthereumAddress.fromHex(params["from"])
        : null;
    final to =
        params["to"] is String ? EthereumAddress.fromHex(params["to"]) : null;
    final value = params["value"] is String
        ? EtherAmount.inWei(BigInt.parse(
            (params["value"] as String).replaceFirst("0x", ''),
            radix: 16))
        : null;
    final data = params["data"] is String
        ? hexToBytes((params["data"] as String).replaceFirst("0x", ''))
        : null;

    final result = await currentClient.estimateGas(
      sender: from,
      to: to,
      value: value,
      data: data,
    );

    return result.toString();
  }

  Future<String?> ethSendTransaction(List paramsList) async {
    final params = paramsList.firstWhere((e) => e is Map<String, dynamic>);
    final from = params["from"] is String
        ? EthereumAddress.fromHex(params["from"])
        : null;
    final to =
        params["to"] is String ? EthereumAddress.fromHex(params["to"]) : null;
    final value = params["value"] is String
        ? EtherAmount.inWei(BigInt.parse(
            (params["value"] as String).replaceFirst("0x", ''),
            radix: 16))
        : null;
    final gas = params['gas'] is String
        ? int.parse((params['gas'] as String).replaceFirst('0x', ''), radix: 16)
        : null;
    final data = params["data"] is String
        ? hexToBytes((params["data"] as String).replaceFirst("0x", ''))
        : null;

    final credentials = await getCredentials(
      from: from,
      to: to,
      value: value,
      gas: gas,
      data: data,
    );

    if (credentials == null)
      throw RPCError(
        4001,
        'User denied transaction signature.',
        null,
      );

    if (isCentralized()) {
      return await web3dart.invokeSmartContract(
        credentials: credentials,
        scAddr: to,
        eun: value?.getInWei ?? BigInt.zero,
        data: data,
        blockChainType: currentNetworkType.isEurusNetwork
            ? BlockChainType.Eurus
            : BlockChainType.Ethereum,
      );
    } else {
      return currentClient.sendTransaction(
        credentials,
        Transaction(
          from: from,
          to: to,
          value: value,
          gasPrice: web3dart.getGasPrice(blockChainType: blockChainType),
          maxGas: gas,
          data: data,
        ),
        chainId: web3dart.getChainId(blockChainType),
      );
    }
  }

  Future<String> ethSign(List paramsList) async {
    String from = paramsList[0];
    String message = paramsList[1];

    if (from != common.currentAddress)
      throw RPCError(-32602,
          'Invalid parameters: must provide an Ethereum address.', null);

    EthPrivateKey? credentials = await getCredentials(
      from: EthereumAddress.fromHex(from),
      message: message,
      signType: DappBrowserConfirmationDialogType.ethSign,
    );
    if (credentials == null)
      throw RPCError(4001, 'User denied message signature.', null);

    return EthSigUtil.signMessage(
      privateKey: bytesToHex(credentials.privateKey),
      message: hexToBytes(message),
    );
  }

  Future<String> personalSign(List paramsList) async {
    String firstParam = paramsList[0];
    String secondParam = paramsList[1];

    String message = firstParam;
    String from = secondParam;

    if ((firstParam.length == EthereumAddress.addressByteLength * 2 + 2) &&
        secondParam.length != EthereumAddress.addressByteLength * 2 + 2) {
      message = secondParam;
      from = firstParam;
    }

    String displayMessage = message;
    try {
      displayMessage = utf8.decode(hexToBytes(message));
    } catch (e) {
      message = bytesToHex(Uint8List.fromList(utf8.encode(message)));
    }

    if (from != common.currentAddress)
      throw RPCError(-32602,
          'Invalid parameters: must provide an Ethereum address.', null);

    EthPrivateKey? credentials = await getCredentials(
      from: EthereumAddress.fromHex(from),
      message: displayMessage,
      signType: DappBrowserConfirmationDialogType.personalSign,
    );
    if (credentials == null)
      throw RPCError(4001, 'User denied message signature.', null);

    final result = EthSigUtil.signPersonalMessage(
      privateKey: bytesToHex(credentials.privateKey),
      message: hexToBytes(message),
    );
    return result;
  }

  Future<String> personalEcRecover(List paramsList) async {
    if (paramsList.length < 2)
      throw RPCError(-32602, 'Invalid parameters.', null);

    final message = paramsList[0] as String;
    final signature = paramsList[1] as String;

    if (isEmptyString(string: message) || isEmptyString(string: signature))
      throw RPCError(-32602, 'Invalid parameters.', null);

    return web3dart
        .personalEcRecover(hexToBytes(message), hexToBytes(signature))
        .hex;
  }

  Future<String> ethSignTypedData(
    List paramsList, {
    required TypedDataVersion version,
  }) async {
    final data;
    String from;
    switch (version) {
      case TypedDataVersion.V1:
        data = (paramsList[0] as List)
            .map((e) => e as Map<String, dynamic>)
            .toList();
        from = (paramsList[1] as String);
        break;
      case TypedDataVersion.V3:
      case TypedDataVersion.V4:
        from = (paramsList[0] as String);
        data = jsonDecode(paramsList[1] as String);
        final chainId = (data['domain'] as Map<String, dynamic>?)?['chainId'];
        int? _chainId =
            chainId is String ? int.tryParse(chainId) : chainId as int?;
        int activeChainId = currentNetworkType.netVersion;

        // eslint-disable-next-line
        if (_chainId != null && _chainId != activeChainId) {
          throw RPCError(
            -32600,
            'Provided chainId ($_chainId) must match the active chainId ($activeChainId})',
            null,
          );
        }
        break;
    }

    if (from != common.currentAddress)
      throw RPCError(-32602,
          'Invalid parameters: must provide an Ethereum address.', null);

    DappBrowserConfirmationDialogType signType;
    switch (version) {
      case TypedDataVersion.V1:
        signType = DappBrowserConfirmationDialogType.signTypeDataV1;
        break;
      case TypedDataVersion.V3:
        signType = DappBrowserConfirmationDialogType.signTypeDataV3;
        break;
      case TypedDataVersion.V4:
        signType = DappBrowserConfirmationDialogType.signTypeDataV4;
        break;
    }

    EthPrivateKey? credentials = await getCredentials(
      from: EthereumAddress.fromHex(from),
      message: new JsonEncoder.withIndent("   ").convert(data),
      signType: signType,
    );
    if (credentials == null)
      throw RPCError(4001, 'User denied message signature.', null);

    return EthSigUtil.signTypedData(
      privateKey: bytesToHex(credentials.privateKey),
      jsonData: jsonEncode(data),
      version: version,
    );
  }

  Future<dynamic> ethGetTransactionReceipt(List paramsList) async {
    final result = await api.post(
      url: currentNetworkType.isEurusNetwork
          ? web3dart.eurusRPCUrl
          : web3dart.blockchainRpcUrl,
      emptyHeader: true,
      body: jsonEncode(<String, Object>{
        "jsonrpc": "2.0",
        "method": "eth_getTransactionReceipt",
        "params": [paramsList.first as String],
        "id": 1,
      }),
    );
    return result["result"];
  }

  Future<dynamic> ethGetTransactionByHash(List paramsList) async {
    final result = await api.post(
      url: currentNetworkType.isEurusNetwork
          ? web3dart.eurusRPCUrl
          : web3dart.blockchainRpcUrl,
      emptyHeader: true,
      body: jsonEncode(<String, Object>{
        "jsonrpc": "2.0",
        "method": "eth_getTransactionByHash",
        "params": [paramsList.first as String],
        "id": 1,
      }),
    );
    return result["result"];
  }

  Future<dynamic> switchEthereumChain(List paramList) async {
    final chainId = paramList.firstWhereOrNull(
        (element) => element is Map<String, dynamic>)['chainId'];
    if (chainId is! String)
      throw RPCError(
          -32602,
          'Expected 0x-prefixed, unpadded, non-zero hexadecimal string \'chainId\'. Received:\n$chainId',
          null);
    final type = NetworkType.values.firstWhereOrNull(
        (element) => (element.isEnabled && element.chainId == chainId));
    if (type == null)
      throw RPCError(4902, 'Unrecognized chain ID $chainId', null);

    dynamic result = await getCredentials(
      signType: DappBrowserConfirmationDialogType.switchNetwork,
      chainId: chainId,
    );
    if (result is! bool) result = false;
    if (result) {
      setNetwork(type);
    } else {
      throw RPCError(
        4001,
        'User rejected the request.',
        null,
      );
    }
  }

  Future<dynamic> getCredentials({
    EthereumAddress? from,
    EthereumAddress? to,
    EtherAmount? value,
    int? gas,
    Uint8List? data,
    String? message,
    String? chainId,
    DappBrowserConfirmationDialogType? signType,
  }) async {
    DappBrowserConfirmationDialogType _transactionType;
    String? title;
    EthereumAddress? _from = from;
    EthereumAddress? _to = to;
    EtherAmount? _value = value;
    CryptoCurrencyModel? _token;
    int? _gas = gas;
    Uint8List? _data = data;

    final tokens = (await common.getSupportedTokens()).values.toList();

    if (data != null) {
      final functionId = data.sublist(0, 4);
      final _data = data.sublist(4);

      DeployedContract deployedContract =
          web3dart.getERC20Contract(this.blockChainType);
      ContractFunction? contractFunction = deployedContract.functions
          .firstWhereOrNull(
              (element) => listEquals(functionId, element.selector));
      _transactionType = contractFunction?.name.toLowerCase() == 'approve'
          ? DappBrowserConfirmationDialogType.approve
          : contractFunction?.name.toLowerCase() == 'transfer'
              ? DappBrowserConfirmationDialogType.transfer
              : DappBrowserConfirmationDialogType.unknown;
      if (_transactionType == DappBrowserConfirmationDialogType.unknown)
        title = contractFunction?.name;

      final decodedValue =
          contractFunction?.decodeReturnValuesFromParameters(_data);
      if (decodedValue != null) {
        _to = decodedValue
            .firstWhereOrNull((element) => element is EthereumAddress);
        _value = EtherAmount.inWei(
            decodedValue.firstWhereOrNull((element) => element is BigInt) ??
                value?.getInWei ??
                BigInt.zero);
      }

      if (value != null && value != EtherAmount.zero()) {
        _token = tokens.firstWhereOrNull((element) =>
            element.symbol.toLowerCase() ==
            (this.blockChainType == BlockChainType.Eurus ? 'eun' : 'eth'));
      } else {
        _token = tokens.firstWhereOrNull((element) =>
            (this.blockChainType == BlockChainType.Eurus
                    ? element.addressEurus
                    : element.addressEthereum)
                ?.toLowerCase() ==
            to?.hex.toLowerCase());
      }
    } else {
      _transactionType = DappBrowserConfirmationDialogType.transfer;
      _token = tokens.firstWhereOrNull((element) =>
          element.symbol.toLowerCase() ==
          (this.blockChainType == BlockChainType.Eurus ? 'eun' : 'eth'));
    }

    final result = await showDialog(
      context: context,
      barrierDismissible: false,
      builder: (_) {
        return DappBrowserConfirmationDialog(
          signType != null && (message != null || chainId != null)
              ? signType
              : _transactionType,
          title: title,
          from: _from,
          to: _to,
          value: _value,
          token: _token,
          gas: _gas,
          data: _data,
          message: message,
          chainId: chainId,
        );
      },
    );
    return result;
  }
}
