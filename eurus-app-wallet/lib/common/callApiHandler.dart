import 'dart:convert';

import 'package:app_authentication_kit/app_authentication_kit.dart';
import 'package:app_authentication_kit/utils/address.dart';
import 'package:app_security_kit/password_encrypt_helper.dart';
import 'package:app_security_kit/rsa_pem.dart';
import 'package:app_storage_kit/normal_storage.dart';
import 'package:bip39/bip39.dart' as bip39;
import 'package:collection/collection.dart';
import 'package:device_info/device_info.dart';
import 'package:dio/dio.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/bitcoinLibrary/bitcoin_flutter.dart';
import 'package:euruswallet/common/user_profile.dart';
import 'package:euruswallet/common/web3dart.dart';
import 'package:euruswallet/model/adminFeeModel.dart';
import 'package:euruswallet/model/blockchainAddressInformation.dart';
import 'package:euruswallet/model/changeLoginPW.dart';
import 'package:euruswallet/model/clientVersion.dart';
import 'package:euruswallet/model/codeVerification.dart';
import 'package:euruswallet/model/createKycStatus.dart';
import 'package:euruswallet/model/faucetRequest.dart';
import 'package:euruswallet/model/findEmailWalletAddress.dart';
import 'package:euruswallet/model/forgetLoginPw.dart';
import 'package:euruswallet/model/getFaucetList.dart';
import 'package:euruswallet/model/getTransactionReceipt.dart';
import 'package:euruswallet/model/importWallet.dart';
import 'package:euruswallet/model/kYCCountryList.dart';
import 'package:euruswallet/model/loginBySignModel.dart';
import 'package:euruswallet/model/pushRawTransaction.dart';
import 'package:euruswallet/model/refreshToken.dart';
import 'package:euruswallet/model/registerByEmail.dart';
import 'package:euruswallet/model/registerDevice.dart';
import 'package:euruswallet/model/sendBitcoinDetail.dart';
import 'package:euruswallet/model/serverConfig.dart';
import 'package:euruswallet/model/setupPaymentWallet.dart';
import 'package:euruswallet/model/signTransaction.dart';
import 'package:euruswallet/model/submitKycDocument.dart';
import 'package:euruswallet/model/userKYCStatus.dart';
import 'package:euruswallet/model/user_marketing_reward_list_response_model.dart';
import 'package:euruswallet/model/user_marketing_reward_scheme_response_model.dart';
import 'package:euruswallet/model/user_storage_request_model.dart';
import 'package:euruswallet/model/user_storage_response_model.dart';
import 'package:http/http.dart' as http;
import 'package:path/path.dart';
import 'package:pointycastle/api.dart';
import 'package:pointycastle/asymmetric/api.dart';
import 'package:pointycastle/key_generators/api.dart';
import 'package:pointycastle/key_generators/rsa_key_generator.dart';
import 'package:web3dart/crypto.dart';
import 'package:web3dart/web3dart.dart' as web3dart_lib;

import 'commonMethod.dart';

class CallApiHandler {
  final client = http.Client();
  static final CallApiHandler _instance = CallApiHandler._internal();

  /// init method
  CallApiHandler._internal();

  factory CallApiHandler() {
    return _instance;
  }

  final stagingRpcUrl = 'http://18.141.43.75:8082';
  final devRpcUrl = 'http://besudevapi.eurus.network:80';
  final testingRpcUrl = 'https://testnetobs.eurus.network:443';
  final productionRpcUrl = 'https://walletapi.eurus.network';
  late String serverUrl;
  late String eurusExplorerUrl;
  late String mainNetExplorerUrl;
  String bscExplorerUrl = "https://bscscan.com/tx/"; //transaction details link
  String bscTxAPIUrl = "https://api.bscscan.com/";

  Future<bool> setUpServerConfig({required EnvType newEnvType}) async {
    envType = newEnvType;
    serverUrl = envType == EnvType.Staging
        ? stagingRpcUrl
        : envType == EnvType.Dev
            ? devRpcUrl
            : envType == EnvType.Testnet
                ? testingRpcUrl
                : productionRpcUrl;

    ServerConfig serverConfig = await setupServerConfig();
    if (await common.checkApiError(
        context: common.currentContext,
        errorString: serverConfig.message,
        returnCode: serverConfig.returnCode)) {
      mainNetExplorerUrl = envType == EnvType.Production
          ? "https://etherscan.io/tx/"
          : "https://rinkeby." "etherscan.io/tx/";
      bscExplorerUrl = "https://testnet.bscscan.com/tx/";
      bscTxAPIUrl = "https://api-testnet.bscscan.com/";
      if (envType == EnvType.Dev) {
        eurusExplorerUrl = "http://13.228.169.25/transactions/";
      } else if (envType == EnvType.Testnet) {
        eurusExplorerUrl =
            "https://testnetexplorer.eurus.network/transactions/";
      } else if (envType == EnvType.Production) {
        eurusExplorerUrl = "https://explorer.eurus.network/transactions/";
      }
      // mainnetExplorerUrl = serverConfig.data.mainnetRPCDomain + "/transactions/";
      web3dart.blockchainRpcUrl = envType == EnvType.Production
          ? (serverConfig.data?.mainnetRPCProtocol ?? '') +
              "://" +
              (serverConfig.data?.mainnetRPCDomain ?? '') +
              ":" +
              (serverConfig.data?.mainnetRPCPort ?? '').toString()
          : "https://rinkeby.infura.io/v3/fa89761e51884ca48dce5c0b6cfef565";
      web3dart.bscRpcUrl = "https://data-seed-prebsc-1-s1.binance.org:8545";
      //'https://bsc-dataseed1.binance.org:443'
      web3dart.eurusRPCUrl = (serverConfig.data?.eurusPRCProtocol ?? '') +
          "://" +
          (serverConfig.data?.eurusRPCDomain ?? '') +
          ":" +
          (serverConfig.data?.eurusRPCPort ?? '').toString();
      web3dart.mainNetEthClient =
          new Web3Client(web3dart.blockchainRpcUrl, web3dart.httpClient);
      web3dart.mainNetBscClient =
          new Web3Client(web3dart.bscRpcUrl, web3dart.httpClient);
      web3dart.eurusEthClient =
          new Web3Client(web3dart.eurusRPCUrl, web3dart.httpClient);
      //web3dart.rinkebyRpcUrl = serverConfig.data.mainnetRPCProtocol+"://"+serverConfig.data.mainnetRPCDomain+":"+serverConfig.data.mainnetRPCPort.toString();
      web3dart.externalSmartContractConfigAddress =
          (serverConfig.data?.externalSmartContractConfigAddress ?? '');
      web3dart.eurusInternalConfigAddress =
          (serverConfig.data?.eurusInternalConfigAddress ?? '');
      if (serverConfig.data != null) {
        web3dart.eurusChainId = serverConfig.data!.eurusChainId;
        web3dart.mainNetChainId = serverConfig.data!.mainnetChainId;
      }
      print("web3dart.eurusRPCUrl:${web3dart.eurusRPCUrl}");
      print("web3dart.rinkebyRpcUrl:${web3dart.blockchainRpcUrl}");
      print("web3dart.bscRpcUrl:${web3dart.bscRpcUrl}");
      print("mainNetExplorerUrl:$mainNetExplorerUrl");
      print("eurusExplorerUrl:$eurusExplorerUrl");
      print("setUpServerConfig");
      return true;
    } else {
      return false;
    }
  }

  dynamic getHeaders({String? apiAccessToken, bool haveApiAccessToken: true}) {
    //print('apiAccessToken:${apiAccessToken??"empty"}');
    return haveApiAccessToken
        ? {
            'Authorization': 'Bearer $apiAccessToken',
            'Content-Type': 'application/json; charset=UTF-8',
          }
        : {
            'Content-Type': 'application/json; charset=UTF-8',
          };
  }

  dynamic post({
    required String url,
    dynamic body,
    bool haveApiAccessToken: true,
    bool emptyHeader: false,
    bool shouldPrintWrapped: true,
    Duration? timeout,
  }) async {
    final _timeout = timeout ?? Duration(seconds: 30);
    if (haveApiAccessToken) {
      await reFreshToken();
    }
    String? apiAccessToken =
        await NormalStorageKit().readValue('apiAccessToken_');
    http.Response rawResponse;
    printWrapped(
      'request: $url header: ${getHeaders(apiAccessToken: apiAccessToken, haveApiAccessToken: haveApiAccessToken)} body: ${body ?? "empty"}',
      shouldWrapped: shouldPrintWrapped,
    );
    try {
      rawResponse = await client
          .post(
            Uri.parse(url),
            headers: emptyHeader
                ? null
                : getHeaders(
                    apiAccessToken: apiAccessToken,
                    haveApiAccessToken: haveApiAccessToken),
            body: body,
          )
          .timeout(_timeout);
      printWrapped(
        'response: $url body: ${rawResponse.body}',
        shouldWrapped: shouldPrintWrapped,
      );
    } catch (e) {
      rawResponse = _getServerMaintanceHttpResponse();
    }
    return jsonDecode(rawResponse.body);
  }

  dynamic get({
    required String url,
    bool haveApiAccessToken: true,
    bool shouldPrintWrapped: true,
    Duration? timeout,
  }) async {
    final _timeout = timeout ?? Duration(seconds: 30);
    if (haveApiAccessToken) {
      await reFreshToken();
    }
    String? apiAccessToken =
        await NormalStorageKit().readValue('apiAccessToken_');
    // if(url.contains('/user/withdrawAdminFee/')){
    //   apiAccessToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTU3ODQxNTcsImlzcyI6IkV1cnVzQXV0aCIsIm5iZiI6MTYxNTc3Njk1NywiY2xpZW50SW5mbyI6IntcIkNsaWVudEluZm9cIjpcIntcXFwibG9naW5BZGRyZXNzXFxcIjpcXFwiMHhjOWRiNWYzZGYzMzRhYmFjYTk4Njc4NWEzOTk0YzAzMWFmZDhlNzBiXFxcIixcXFwidXNlcklkXFxcIjo1MH1cIixcIlNlcnZpY2VJZFwiOjB9In0.MciLijoXEZJY8NG4ZZzh0ke53dxcldD8B_b8_We-kno';
    // }
    http.Response rawResponse;
    printWrapped(
      'request: $url header: ${getHeaders(apiAccessToken: apiAccessToken, haveApiAccessToken: haveApiAccessToken)}',
      shouldWrapped: shouldPrintWrapped,
    );
    try {
      rawResponse = await client
          .get(
            Uri.parse(url),
            headers: getHeaders(
                apiAccessToken: apiAccessToken,
                haveApiAccessToken: haveApiAccessToken),
          )
          .timeout(_timeout);
    } catch (e) {
      rawResponse = _getServerMaintanceHttpResponse();
    }
    printWrapped(
      'response: $url body: ${rawResponse.body}}',
      shouldWrapped: shouldPrintWrapped,
    );
    String rawResponseString = rawResponse.body.toString().toLowerCase();
    if (rawResponseString.contains("502 bad gateway") ||
        rawResponseString.contains("server maintenance") ||
        rawResponseString.contains("<html") ||
        rawResponseString.contains("internal server error") ||
        rawResponseString.contains("404 not found") ||
        rawResponseString.contains("http 404") ||
        rawResponseString.contains('connection refused') ||
        (rawResponse.statusCode >= 400 && rawResponse.statusCode < 600)) {
      rawResponse = _getServerMaintanceHttpResponse();
    }
    return jsonDecode(rawResponse.body);
  }

  dynamic dioPost({
    required String url,
    dynamic data,
    bool haveApiAccessToken: true,
    bool shouldPrintWrapped: true,
  }) async {
    Dio dio = new Dio();
    String? apiAccessToken =
        await NormalStorageKit().readValue('apiAccessToken_');
    Map<String, dynamic> header = haveApiAccessToken == true
        ? {'Authorization': 'Bearer $apiAccessToken'}
        : {};
    printWrapped('request: $url  body: ${data ?? "empty"}',
        shouldWrapped: shouldPrintWrapped);
    Response rawResponse =
        await dio.post(url, data: data, options: Options(headers: header));
    printWrapped('response: $url  body: ${rawResponse.data}',
        shouldWrapped: shouldPrintWrapped);
    return rawResponse.data;
  }

  http.Response _getServerMaintanceHttpResponse() {
    return http.Response(
      jsonEncode({
        "returnCode": 503,
        "message": "Server maintenance",
        "nonce": "",
        "data": null
      }),
      503,
    );
  }

  dynamic recentTransaction({
    required BlockChainType chain,
    String? nonce,
    String? symbol,
  }) async {
    var response = await post(
        url: serverUrl + '/user/recentTransaction',
        body: jsonEncode(
          {
            "nonce": nonce,
            "currencySymbol": symbol?.toUpperCase(),
            "chainId": chain == BlockChainType.Eurus
                ? web3dart.eurusChainId
                : chain == BlockChainType.Ethereum
                    ? web3dart.mainNetChainId
                    : chain == BlockChainType.BinanceCoin
                        ? web3dart.bscChainId
                        : null,
          },
        ));
    return response;
  }

  dynamic withdrawAdminFeeDetail({
    String? nonce,
    required String symbol,
  }) async {
    var response = await get(
        url: serverUrl + '/user/withdrawAdminFee/${symbol.toUpperCase()}');
    return response;
  }

  Future<bool> reFreshToken() async {
    String? apiAccessTokenExpiryTime =
        await NormalStorageKit().readValue('apiAccessTokenExpiryTime_');
    print('apiAccessTokenExpiryTime_:$apiAccessTokenExpiryTime');
    int? expiryTime = int.tryParse(apiAccessTokenExpiryTime ?? '');
    if (expiryTime == null) return false;
    DateTime dateTimeCreatedAt = DateTime.fromMillisecondsSinceEpoch(
      (expiryTime * 1000),
      isUtc: true,
    );
    DateTime dateTimeNow = DateTime.now().toUtc();
    final differenceInSeconds =
        dateTimeNow.difference(dateTimeCreatedAt).inSeconds;

    // after 1 hours then logout or remaining 20 min (only use 40 min) call refreshToken // you must call refreshToken at least remaining 30 min or token update too frequently
    int compareResult = differenceInSeconds.compareTo(-1200);
    print('differenceInSeconds:$differenceInSeconds');
    if (differenceInSeconds > 0) {
      common.logoutFunction?.add(true);
    } else if (compareResult >= 0) {
      String? apiAccessToken =
          await NormalStorageKit().readValue('apiAccessToken_');
      var rawResponse;
      String url = serverUrl + '/user/refreshToken';
      rawResponse = await client.post(
        Uri.parse(url),
        headers: getHeaders(apiAccessToken: apiAccessToken),
      );
      Map<String, dynamic> response = jsonDecode(rawResponse.body);
      print('reFreshToken.body${rawResponse.body}');
      RefreshToken refreshToken = RefreshToken.fromJson(response);
      if (refreshToken.data != null) {
        await NormalStorageKit()
            .setValue(refreshToken.data!.token, 'apiAccessToken_');
        await NormalStorageKit().setValue(
            refreshToken.data!.expiryTime.toString(),
            'apiAccessTokenExpiryTime_');
      }
      print('response:$response');
    }

    return true;
  }

  Future<AdminFeeModel> getAdminFee({required String selectTokenSymbol}) async {
    Map<String, dynamic> response =
        await get(url: serverUrl + '/user/withdrawAdminFee/$selectTokenSymbol');
    final adminFeeModel = AdminFeeModel.fromJson(response);
    common.adminFeeModel = adminFeeModel;
    if (adminFeeModel.data != null)
      adminFeeModel.data!.actualFee = common.divisionDecimal(
          adminFeeModel.data!.decimal, adminFeeModel.data!.fee);
    return adminFeeModel;
  }

  Future<GetFaucetList?> getFaucetList() async {
    if (envType == EnvType.Production) {
      return null;
    }
    Map<String, dynamic> response = await get(
      url: serverUrl + '/user/testnet/asset/faucet',
      haveApiAccessToken: true,
    );
    final getFaucetList = GetFaucetList.fromJson(response);
    common.getFaucetList = getFaucetList;
    if (getFaucetList.returnCode != 0) {
      common.getFaucetList = null;
    }
    return common.getFaucetList;
  }

  Future<FaucetRequest> faucetRequest({String? symbol}) async {
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/testnet/asset/faucet' + '/$symbol',
      haveApiAccessToken: true,
    );
    final faucetRequest = FaucetRequest.fromJson(response);
    common.faucetRequest = faucetRequest;
    return faucetRequest;
  }

  Future<RegisterByEmail> registerByEmail({
    required String email,
    required String password,
  }) async {
    String deviceId = await getDeviceID(prefix: email);

    int currentTimeStamp = getCurrentTimeStamp();
    final loginAddressPair =
        await common.getAddressPair(email: email, password: password);
    common.loginAddressPair = loginAddressPair;
    var json = jsonEncode({
      "nonce": currentTimeStamp.toString(),
      "timestamp": currentTimeStamp,
      "deviceId": deviceId,
      "email": email,
      "loginAddress": loginAddressPair.address,
      "signature": getSignature(
          deviceId: deviceId,
          currentTimeStamp: currentTimeStamp,
          walletAddress: loginAddressPair.address,
          privateKey: loginAddressPair.privateKey),
      "publicKey": loginAddressPair.publicKey
    });

    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/registerByEmail',
      body: json,
      haveApiAccessToken: false,
    );
    final registerByEmail = RegisterByEmail.fromJson(response);
    common.registerByEmail = registerByEmail;
    return registerByEmail;
  }

  Future<LoginBySignModel> loginBySignature({
    required String email,
    required String password,
  }) async {
    String deviceId = await getDeviceID(prefix: email);

    int currentTimeStamp = getCurrentTimeStamp();
    final loginAddressPair =
        await common.getAddressPair(email: email, password: password);
    common.loginAddressPair = loginAddressPair;
    final json = jsonEncode({
      "nonce": currentTimeStamp.toString(),
      "timestamp": currentTimeStamp,
      "deviceId": deviceId,
      "sign": getSignature(
          deviceId: deviceId,
          currentTimeStamp: currentTimeStamp,
          walletAddress: loginAddressPair.address,
          privateKey: loginAddressPair.privateKey),
      "publicKey": loginAddressPair.publicKey,
      "walletAddress": loginAddressPair.address,
    });

    Map<String, dynamic> response =
        await post(url: serverUrl + '/user/loginBySignature', body: json);
    LoginBySignModel result = LoginBySignModel.fromJson(response);
    common.cenUserWalletAddress = result.walletAddress;
    common.cenMainNetWalletAddress = result.mainnetWalletAddress;
    common.ownerWalletAddress = result.ownerWalletAddress;
    await NormalStorageKit().setValue(email, 'cenUserEmail');
    await NormalStorageKit().setValue(password, 'cenUserPassword');
    String rsaPrivateKeyName = '${email}rsaPrivateKey';
    common.rsaPrivateKey =
        await NormalStorageKit().readValue(rsaPrivateKeyName);
    if (result.returnCode == 0 && !isEmptyString(string: password)) {
      await NormalStorageKit().setValue(result.token ?? '', 'apiAccessToken_');
      await NormalStorageKit()
          .setValue(result.expiryTime.toString(), 'apiAccessTokenExpiryTime_');
      common.email = email.toLowerCase();
      common.serverMnemonic = result.decryptedMnemonic;
      common.serverAddressPair = null;
      print("common.serverMnemonic:${common.serverMnemonic}");

      final PasswordEncryptHelper pwHelper =
          PasswordEncryptHelper(password: password);

      final walletAddress = result.walletAddress ?? '0x0';
      final encryptedAddress = pwHelper.encryptWPwd(walletAddress);

      final userProfiles = await getLocalAcs();
      UserProfile? userProfile = userProfiles.firstWhereOrNull(
          (element) => element.address == walletAddress.toLowerCase());
      userProfile?.email = email;
      userProfile?.encryptedAddress = encryptedAddress;
      userProfile?.encryptedPrivateKey =
          pwHelper.encryptWPwd(loginAddressPair.privateKey);
      userProfile?.lastLoginTime = getCurrentTimeStamp().toString();
      List<Uint8List>? successTransactions = [];
      for (var transaction in userProfile?.pendingTransaction ?? []) {
        try {
          final txHash =
              await web3dart.eurusEthClient.sendRawTransaction(transaction);
          if (!isEmptyString(string: txHash)) {
            successTransactions.add(transaction);
          }
        } on Exception {
          continue;
        }
      }
      successTransactions.forEach((element) {
        userProfile?.pendingTransaction?.remove(element);
      });

      if (walletAddress.toLowerCase() != '0x0')
        await setAcToLocal(userProfile ??
            UserProfile.fromJson({
              "userType": CurrentUserType.centralized,
              "address": walletAddress.toLowerCase(),
              "email": email,
              "encryptedAddress": encryptedAddress,
              "encryptedPrivateKey":
                  pwHelper.encryptWPwd(loginAddressPair.privateKey),
              "lastLoginTime": getCurrentTimeStamp().toString(),
            }));
    }
    return result;
  }

  String getSignature({
    required String deviceId,
    required int currentTimeStamp,
    required String walletAddress,
    required String privateKey,
  }) {
    String signatureRawString = "deviceId=" +
        deviceId +
        "&timestamp=" +
        currentTimeStamp.toString() +
        "&walletAddress=" +
        walletAddress;
    String signature = encryptStringFromPrivateKey(
            message: signatureRawString, privateKey: privateKey)
        .substring(0, 128);
    return signature;
  }

  Future<ChangeLoginPW> changeLoginPassword({
    required String oldLoginPassword,
    required String newLoginPassword,
  }) async {
    String deviceId = await getDeviceID(prefix: common.email ?? '');

    int currentTimeStamp = getCurrentTimeStamp();
    AddressPair oldAddressPair = await common.getAddressPair(
        email: common.email, password: oldLoginPassword);
    AddressPair newAddressPair = await common.getAddressPair(
        email: common.email, password: newLoginPassword);
    final serverAddressPair = await common.getAddressPair(
        email: common.email,
        password: newLoginPassword,
        mnemonic: common.serverMnemonic,
        addressPairType: AddressPairType.paymentPw);
    ChangeLoginPW changeLoginPW = ChangeLoginPW(1, "", "", null);
    if (common.ownerWalletAddress?.toUpperCase() ==
        serverAddressPair.address.toUpperCase()) {
      changeLoginPW = ChangeLoginPW(
          1222, 'REGISTER.LOGIN_PW_NOT_EQUAL_PAYMENT_PW'.tr(), "", null);
    } else {
      final json = jsonEncode({
        "nonce": currentTimeStamp.toString(),
        "oldLoginAddress": oldAddressPair.address,
        "loginAddress": newAddressPair.address,
        "timestamp": currentTimeStamp,
        "deviceId": deviceId,
        "sign": getSignature(
            deviceId: deviceId,
            currentTimeStamp: currentTimeStamp,
            walletAddress: newAddressPair.address,
            privateKey: newAddressPair.privateKey),
        "oldSign": getSignature(
            deviceId: deviceId,
            currentTimeStamp: currentTimeStamp,
            walletAddress: oldAddressPair.address,
            privateKey: oldAddressPair.privateKey),
        "publicKey": newAddressPair.publicKey,
        "oldPublicKey": oldAddressPair.publicKey
      });

      Map<String, dynamic> response = await post(
        url: serverUrl + '/user/changeLoginPassword',
        body: json,
        haveApiAccessToken: true,
      );

      changeLoginPW = ChangeLoginPW.fromJson(response);
    }
    return changeLoginPW;
  }

  Future<bool> changePaymentPassword({
    required BuildContext context,
    required String oldPaymentPassword,
    required String newPaymentPassword,
  }) async {
    String deviceId = await getDeviceID(prefix: common.email ?? '');

    int currentTimeStamp = getCurrentTimeStamp();
    AddressPair oldAddressPair = await common.getAddressPair(
        email: common.email,
        password: oldPaymentPassword,
        mnemonic: common.serverMnemonic,
        addressPairType: AddressPairType.paymentPw);
    AddressPair newAddressPair = await common.getAddressPair(
        email: common.email,
        password: newPaymentPassword,
        mnemonic: common.serverMnemonic,
        addressPairType: AddressPairType.paymentPw);
    print("common.email:${common.email}");
    print("oldPaymentPassword:$oldPaymentPassword");
    print("newPaymentPassword:$newPaymentPassword");
    ChangeLoginPW changeLoginPW = ChangeLoginPW(1, "", "", null);
    final loginAddressPair = await common.getAddressPair(
        email: common.email, password: newPaymentPassword);
    if (common.loginAddressPair?.address == loginAddressPair.address) {
      changeLoginPW = ChangeLoginPW(
          1222, 'REGISTER.LOGIN_PW_NOT_EQUAL_PAYMENT_PW'.tr(), "", null);
    } else {
      final json = jsonEncode({
        "nonce": currentTimeStamp.toString(),
        "oldOwnerWalletAddress": oldAddressPair.address,
        "ownerWalletAddress": newAddressPair.address,
        "timestamp": currentTimeStamp,
        "deviceId": deviceId,
        "sign": getSignature(
            deviceId: deviceId,
            currentTimeStamp: currentTimeStamp,
            walletAddress: newAddressPair.address,
            privateKey: newAddressPair.privateKey),
        "oldSign": getSignature(
            deviceId: deviceId,
            currentTimeStamp: currentTimeStamp,
            walletAddress: oldAddressPair.address,
            privateKey: oldAddressPair.privateKey),
        "publicKey": newAddressPair.publicKey,
        "oldPublicKey": oldAddressPair.publicKey
      });

      Map<String, dynamic> response = await post(
        url: serverUrl + '/user/changePaymentPassword',
        body: json,
        haveApiAccessToken: true,
      );

      changeLoginPW = ChangeLoginPW.fromJson(response);

      if (await common.checkApiError(
        context: context,
        errorString: changeLoginPW.message,
        returnCode: changeLoginPW.returnCode,
      )) {
        final client = web3dart.eurusEthClient;
        final gasPrice = EtherAmount.inWei(
            BigInt.from(web3dart.getBlockChainGasPrice(BlockChainType.Eurus)));
        EtherAmount value = (await client
            .getBalance(EthereumAddress.fromHex(oldAddressPair.address)));
        final estimateGas = await client.estimateGas(
          sender: EthereumAddress.fromHex(oldAddressPair.address),
          to: EthereumAddress.fromHex(newAddressPair.address),
          value: value,
          gasPrice: gasPrice,
        );

        value = EtherAmount.inWei(BigInt.from(value.getInWei.toDouble() -
            estimateGas.toDouble() * gasPrice.getInWei.toDouble()));

        if (value.getInWei.toDouble() >= 0) {
          final signedTransaction =
              await web3dart.eurusEthClient.signTransaction(
            EthPrivateKey.fromHex(oldAddressPair.privateKey),
            web3dart_lib.Transaction(
              gasPrice: gasPrice,
              to: EthereumAddress.fromHex(newAddressPair.address),
              value: value,
            ),
            chainId: web3dart.eurusChainId,
          );

          try {
            final txHash = await web3dart.eurusEthClient
                .sendRawTransaction(signedTransaction);
            print('cen changePaymentPassword wallet transfer: $txHash');
          } on Exception catch (exception) {
            final userProfile = common.currentUserProfile;
            if (userProfile != null) {
              var pendingTransaction = userProfile.pendingTransaction ?? [];
              pendingTransaction.add(signedTransaction);
              userProfile.pendingTransaction = pendingTransaction;
              await setAcToLocal(userProfile);
              print(
                  'cen changePaymentPassword wallet transfer: ${exception.toString()}');
            }
          }
        }

        return true;
      } else {
        return false;
      }
    }
    return false;
  }

  Future<ForgetLoginPw> forgetLoginPw({
    String? email,
  }) async {
    int currentTimeStamp = getCurrentTimeStamp();
    var json = jsonEncode({
      "nonce": currentTimeStamp.toString(),
      "email": email,
    });

    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/forgetLoginPassword',
      body: json,
      haveApiAccessToken: false,
    );
    final forgetLoginPw = ForgetLoginPw.fromJson(response);
    common.forgetLoginPw = forgetLoginPw;
    return forgetLoginPw;
  }

  Future<CodeVerification> forGetLoginPwCodeVerification({
    required String code,
    required String email,
  }) async {
    int currentTimeStamp = getCurrentTimeStamp();
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/verifyForgetLoginPassword',
      body: jsonEncode(
        {"nonce": currentTimeStamp.toString(), "code": code, "email": email},
      ),
      haveApiAccessToken: false,
    );
    final codeVerification = CodeVerification.fromJson(response);
    common.codeVerification = codeVerification;
    return codeVerification;
  }

  Future<CodeVerification> resetLoginPassword({
    required String email,
    required String password,
  }) async {
    String deviceId = await getDeviceID(prefix: email);
    int currentTimeStamp = getCurrentTimeStamp();
    final loginAddressPair =
        await common.getAddressPair(email: email, password: password);
    common.loginAddressPair = loginAddressPair;

    final json = jsonEncode({
      "nonce": currentTimeStamp.toString(),
      "timestamp": currentTimeStamp,
      "deviceId": deviceId,
      "sign": getSignature(
          deviceId: deviceId,
          currentTimeStamp: currentTimeStamp,
          walletAddress: loginAddressPair.address,
          privateKey: loginAddressPair.privateKey),
      "publicKey": loginAddressPair.publicKey,
      "loginAddress": loginAddressPair.address,
    });

    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/resetLoginPassword',
      body: json,
      haveApiAccessToken: true,
    );
    final codeVerification = CodeVerification.fromJson(response);
    common.codeVerification = codeVerification;
    return codeVerification;
  }

  Future<CodeVerification> codeVerification({
    required String code,
    required String email,
  }) async {
    String deviceId = await getDeviceID(prefix: email);
    int currentTimeStamp = getCurrentTimeStamp();
    await genPublicKeyAndPrivateKey();
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/verification',
      body: jsonEncode(
        {
          "nonce": currentTimeStamp.toString(),
          "code": code,
          "email": email,
          "deviceId": deviceId,
          "publicKey": common.rsaPublicKey
        },
      ),
      haveApiAccessToken: false,
    );
    final codeVerification = CodeVerification.fromJson(response);
    common.codeVerification = codeVerification;
    return codeVerification;
  }

  Future<RegisterByEmail> reSendCodeVerification({
    int? userID,
  }) async {
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/resendVerificationEmail',
      body: jsonEncode(
        {
          "userId": userID,
        },
      ),
      haveApiAccessToken: false,
    );
    RegisterByEmail registerByEmail = RegisterByEmail.fromJson(response);
    return registerByEmail;
  }

  Future<SetupPaymentWallet> setupPaymentWallet({
    int? userID,
    String? walletAddress,
  }) async {
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/setupPaymentWallet',
      body: jsonEncode(
        {
          "userId": userID,
          "address": walletAddress,
        },
      ),
      haveApiAccessToken: true,
      timeout: Duration(seconds: 65),
    );
    SetupPaymentWallet setupPaymentWallet =
        SetupPaymentWallet.fromJson(response);
    return setupPaymentWallet;
  }

  Future<ServerConfig> setupServerConfig() async {
    Map<String, dynamic> response = await get(
      url: serverUrl + '/user/serverConfig',
      haveApiAccessToken: false,
    );
    ServerConfig serverConfig = ServerConfig.fromJson(response);
    return serverConfig;
  }

  Future<FindEmailWalletAddress> findEmailWalletAddress({
    required String targetAddress,
  }) async {
    String walletAddress = "";
    String email = "";
    if (EthAddress().isValidEthereumAddress(targetAddress)) {
      walletAddress = targetAddress;
    } else {
      if (targetAddress.contains('@')) {
        email = targetAddress;
      }
    }

    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/findEmailWalletAddress',
      body: jsonEncode(
        {
          "walletAddress": walletAddress,
          "email": email,
        },
      ),
      haveApiAccessToken: true,
    );
    FindEmailWalletAddress findEmailWalletAddress =
        FindEmailWalletAddress.fromJson(response);
    return findEmailWalletAddress;
  }

  Future<ImportWallet> importWallet({
    required AddressPair anAddressPair,
  }) async {
    final deviceId = await getDeviceID(prefix: anAddressPair.address),
        ts = getCurrentTimeStamp(),
        walletAddress = anAddressPair.address.replaceFirst(r'0x', ''),
        publicKey = anAddressPair.publicKey.replaceFirst(r'0x', ''),
        privateKey = anAddressPair.privateKey;
    Map<String, dynamic> response = await post(
      url: api.serverUrl + '/user/importWallet',
      body: jsonEncode(<String, Object>{
        'nonce': ts.toString(),
        'deviceId': deviceId,
        'timestamp': ts,
        'walletAddress': walletAddress,
        'publicKey': publicKey,
        'sign': getSignature(
          deviceId: deviceId,
          currentTimeStamp: ts,
          privateKey: privateKey,
          walletAddress: walletAddress,
        ),
      }),
      haveApiAccessToken: false,
    );

    final importWallet = ImportWallet.fromJson(response);
    common.importWallet = importWallet;
    return importWallet;
  }

  Future<GetTransactionReceipt> getTransactionReceipt({
    required String txId,
  }) async {
    var result = await post(
      url: web3dart.getRpcUrl(common.fromBlockChainType),
      haveApiAccessToken: false,
      body: jsonEncode(<String, Object>{
        "jsonrpc": "2.0",
        "method": "eth_getTransactionReceipt",
        "params": [txId],
        "id": 1,
      }),
    );
    GetTransactionReceipt getTransactionReceipt =
        GetTransactionReceipt.fromJson(result);
    return getTransactionReceipt;
  }

  Future<SignTransaction> signTransactionFromServer({
    required String value,
    required int gasPrice,
    required String inputFunction,
  }) async {
    int currentTimeStamp = getCurrentTimeStamp();
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/signTransaction',
      body: jsonEncode(
        {
          "nonce": currentTimeStamp.toString(),
          "value": value,
          "gasPrice": gasPrice,
          "inputFunction": inputFunction
        },
      ),
      haveApiAccessToken: true,
    );
    SignTransaction signTransaction = SignTransaction.fromJson(response);
    return signTransaction;
  }

  Future<String?> topUpPaymentWallet({
    required int targetGasAmount,
  }) async {
    final userWalletContract = web3dart.getUserWallet(
        contractAddrString: common.cenUserWalletAddress ?? '');

    final signature = web3dart.parameterSignature(
      functionName: 'topUpPaymentWallet',
      userWalletContract: userWalletContract,
      params: [
        EthereumAddress.fromHex(common.ownerWalletAddress ?? ''),
      ],
      paramsName: ['paymentWalletAddr'],
    );

    int currentTimeStamp = getCurrentTimeStamp();
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/topUpPaymentWallet',
      body: jsonEncode(
        {
          "nonce": currentTimeStamp.toString(),
          "targetGasAmount": targetGasAmount,
          "signature": bytesToHex(signature),
        },
      ),
      haveApiAccessToken: true,
    );

    if (response['returnCode'] as int == 0) {
      final txHash = await web3dart.sendRawTransaction(
        signedTransactionString: response['data']['tx'] as String,
      );
      return getReceipt(txHash);
    } else {
      throw ErrorDescription(response['returnCode'] as int == -29
          ? 'TX_PAGE.ERROR.INSUFFICIENT_ASSET'.tr()
          : response['message'] as String);
    }
  }

  Future<String> getReceipt(
    String txHash, {
    int numberOfAttempts = 0,
  }) async {
    final result = (await api.getTransactionReceipt(txId: txHash)).result;
    if (result != null) {
      final revertReason = result.revertReason;
      if (revertReason != null) {
        String? receiptString = revertReason
            .replaceAll('0x', '00')
            .replaceAll('08c379a', '0000000');
        List<String> splitted = [];
        for (int i = 0; i < receiptString.length; i = i + 2) {
          splitted.add(receiptString.substring(
              i, i + 2 > receiptString.length ? receiptString.length : i + 2));
        }
        String ascii2 = List.generate(splitted.length,
                (i) => String.fromCharCode(int.parse(splitted[i], radix: 16)))
            .join()
            .replaceAll(new RegExp(r'[^A-Za-z0-9  *]'), '');
        throw ErrorDescription(ascii2);
      } else {
        return txHash;
      }
    } else {
      if (numberOfAttempts >= 10) throw ErrorDescription('Request Timed Out');
      return Future.delayed(
        Duration(seconds: 3),
        () => getReceipt(
          txHash,
          numberOfAttempts: numberOfAttempts + 1,
        ),
      );
    }
  }

  Future<ClientVersion> clientVersion() async {
    Map<String, dynamic> response = await get(
      url: serverUrl + '/user/clientVersion',
      haveApiAccessToken: false,
    );
    ClientVersion clientVersion = ClientVersion.fromJson(response);
    return clientVersion;
  }

  Future<String> getDeviceID({required String prefix}) async {
    DeviceInfoPlugin deviceInfo = DeviceInfoPlugin();
    String deviceIdKey = '${prefix.toLowerCase()}deviceId';
    String? deviceId = await NormalStorageKit().readValue(deviceIdKey);
    if (!isEmptyString(string: deviceId)) return deviceId ?? '';

    String _deviceId = '';
    if (Platform.isAndroid) {
      AndroidDeviceInfo androidInfo = await deviceInfo.androidInfo;
      print('Running on ${androidInfo.model}'); // e
      _deviceId = androidInfo.androidId;
    } else if (Platform.isIOS) {
      IosDeviceInfo iosInfo = await deviceInfo.iosInfo;
      print('Running on ${iosInfo.model}'); //
      _deviceId = iosInfo.identifierForVendor;
    }

    await NormalStorageKit().setValue(_deviceId, deviceIdKey);
    return _deviceId;
  }

  Future<RegisterDevice> registerDevice() async {
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/registerDevice',
      haveApiAccessToken: true,
    );
    RegisterDevice registerDevice = RegisterDevice.fromJson(response);
    return registerDevice;
  }

  Future<CodeVerification> verifyDevice({
    required String code,
    String? publicKey,
  }) async {
    String deviceId = await getDeviceID(prefix: common.email ?? '');
    int currentTimeStamp = getCurrentTimeStamp();
    await genPublicKeyAndPrivateKey();
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/verifyDevice',
      body: jsonEncode(
        {
          "nonce": currentTimeStamp.toString(),
          "code": code,
          "deviceId": deviceId,
          "publicKey": common.rsaPublicKey
        },
      ),
      haveApiAccessToken: true,
    );
    final codeVerification = CodeVerification.fromJson(response);
    common.codeVerification = codeVerification;
    return codeVerification;
  }

  Future<bool> genPublicKeyAndPrivateKey() async {
    RsaKeyHelper rsaKeyHelper = RsaKeyHelper();
    SecureRandom random = rsaKeyHelper.getSecureRandom();
    AsymmetricKeyPair keyPair = await compute(getRsaKeyPair, random);
    RSAPrivateKey privateKey = keyPair.privateKey as RSAPrivateKey;
    final rsaPrivateKey = rsaKeyHelper.encodePrivateKeyToPemPKCS1(privateKey);
    common.rsaPrivateKey = rsaPrivateKey;
    RSAPublicKey publicKey = keyPair.publicKey as RSAPublicKey;
    final rsaPublicKey = rsaKeyHelper.removePemHeaderAndFooter(
        rsaKeyHelper.encodePublicKeyToPemPKCS1(publicKey));
    common.rsaPublicKey = rsaPublicKey;
    await NormalStorageKit()
        .setValue(rsaPrivateKey, (common.email ?? '') + "rsaPrivateKey");
    await NormalStorageKit()
        .setValue(rsaPublicKey, (common.email ?? '') + "rsaPublicKey");
    return true;
  }

  Future<ForgetLoginPw> forgetPaymentPw() async {
    int currentTimeStamp = getCurrentTimeStamp();
    var json = jsonEncode({
      "nonce": currentTimeStamp.toString(),
    });

    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/forgetPaymentPassword',
      body: json,
      haveApiAccessToken: true,
    );
    final forgetLoginPw = ForgetLoginPw.fromJson(response);
    common.forgetLoginPw = forgetLoginPw;
    return forgetLoginPw;
  }

  Future<CodeVerification> verifyForgetPaymentPassword({
    required String code,
  }) async {
    int currentTimeStamp = getCurrentTimeStamp();
    String deviceId = await getDeviceID(prefix: common.email ?? '');
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/verifyForgetPaymentPassword',
      body: jsonEncode(
        {
          "nonce": currentTimeStamp.toString(),
          "code": code,
          "deviceId": deviceId
        },
      ),
      haveApiAccessToken: true,
    );
    final codeVerification = CodeVerification.fromJson(response);
    common.codeVerification = codeVerification;
    return codeVerification;
  }

  Future<CodeVerification> resetPaymentPassword({
    String? email,
    String? password,
    String? mnemonic,
  }) async {
    String deviceId = await getDeviceID(prefix: email ?? '');
    int currentTimeStamp = getCurrentTimeStamp();

    CodeVerification codeVerification = CodeVerification(1, "", "", null);
    final loginAddressPair =
        await common.getAddressPair(email: common.email, password: password);
    if (common.loginAddressPair?.address == loginAddressPair.address) {
      codeVerification = CodeVerification(
          1222, 'REGISTER.LOGIN_PW_NOT_EQUAL_PAYMENT_PW'.tr(), "", null);
    } else {
      AddressPair serverAddressPair = await common.getAddressPair(
          email: email ?? '',
          password: password ?? '',
          mnemonic: mnemonic,
          addressPairType: AddressPairType.paymentPw);

      final json = jsonEncode({
        "nonce": currentTimeStamp.toString(),
        "timestamp": currentTimeStamp,
        "deviceId": deviceId,
        "sign": getSignature(
            deviceId: deviceId,
            currentTimeStamp: currentTimeStamp,
            walletAddress: serverAddressPair.address,
            privateKey: serverAddressPair.privateKey),
        "publicKey": serverAddressPair.publicKey,
        "ownerWalletAddress": serverAddressPair.address,
      });

      Map<String, dynamic> response = await post(
        url: serverUrl + '/user/resetPaymentPassword',
        body: json,
        haveApiAccessToken: true,
      );
      codeVerification = CodeVerification.fromJson(response);
    }
    return codeVerification;
  }

  Future<UserKYCStatus> getKYCStatus() async {
    Map<String, dynamic> response = await get(
        url: serverUrl + '/user/kyc/getKYCStatusByToken',
        haveApiAccessToken: true);
    final userKYCStatus = UserKYCStatus.fromJson(response);
    common.userKYCStatus = userKYCStatus.returnCode == 0 ? userKYCStatus : null;
    return userKYCStatus;
  }

  Future<KYCCountryList> getKYCCountryList() async {
    Map<String, dynamic> response = await get(
        url: serverUrl + '/user/kyc/getKYCCountryList',
        haveApiAccessToken: true);
    final kYCCountryList = KYCCountryList.fromJson(response);
    common.kYCCountryList =
        kYCCountryList.returnCode == 0 ? kYCCountryList : null;
    return kYCCountryList;
  }

  Future<CreateKycStatus> createKYCStatus({
    String? kycCountryCode,
    int? kycDoc,
  }) async {
    int currentTimeStamp = getCurrentTimeStamp();
    Map<String, dynamic> response = await post(
        url: serverUrl + '/user/kyc/createKYCStatus',
        body: jsonEncode(
          {
            "nonce": currentTimeStamp.toString(),
            "kycCountryCode": kycCountryCode,
            "kycDoc": kycDoc
          },
        ),
        haveApiAccessToken: true);
    CreateKycStatus createKycStatus = CreateKycStatus.fromJson(response);
    return createKycStatus;
  }

  Future<SubmitKycDocument> submitKYCApproval({int? id}) async {
    int currentTimeStamp = getCurrentTimeStamp();
    Map<String, dynamic> response = await post(
        url: serverUrl + '/user/kyc/submitKYCApproval',
        body: jsonEncode(
          {"nonce": currentTimeStamp.toString(), "id": id},
        ),
        haveApiAccessToken: true);
    SubmitKycDocument submitKycDocument = SubmitKycDocument.fromJson(response);
    return submitKycDocument;
  }

  Future<dynamic> submitKYCDocument({
    int? id,
    int? imageType,
    String? imgPath,
  }) async {
    dynamic json = jsonEncode({
      "userKYCStatusId": id,
      "imageType": imageType,
      "fileExtension": 'jpg',
    });
    print('data: $json');
    FormData formData = FormData.fromMap({
      "requestJson": json,
      if (imgPath != null)
        "submitImage":
            await MultipartFile.fromFile(imgPath, filename: basename(imgPath))
    });
    Map<String, dynamic> response = await dioPost(
        url: serverUrl + '/user/kyc/submitKYCDocument',
        data: formData,
        haveApiAccessToken: true);
    SubmitKycDocument submitKycDocument = SubmitKycDocument.fromJson(response);
    return submitKycDocument;
  }

  Future<UserStorageResponseModel> getUserStorage() async {
    Map<String, dynamic> response = await get(url: serverUrl + '/user/storage');
    return UserStorageResponseModel.fromJson(response);
  }

  Future<UserStorageResponseModel> postUserStorage({
    String? currentLanguage,
  }) async {
    final nonce = getCurrentTimeStamp().toString();
    var userStorage = common.userStorage;
    if (currentLanguage != null) userStorage.currentLanguage = currentLanguage;
    final requestModel =
        UserStorageRequestModel(nonce: nonce, storage: userStorage);
    Map<String, dynamic> response = await post(
      url: serverUrl + '/user/storage',
      body: jsonEncode(requestModel.toJson()),
    );
    return UserStorageResponseModel.fromJson(response);
  }

  Future<UserMarketingRewardSchemeResponseModel>
      getUserMarketingRewardScheme() async {
    Map<String, dynamic> response = await get(
      url: serverUrl + '/user/marketing/rewardScheme',
    );
    return UserMarketingRewardSchemeResponseModel.fromJson(response);
  }

  Future<UserMarketingRewardListResponseModel>
      getUserMarketingRewardList() async {
    Map<String, dynamic> response = await get(
      url: serverUrl + '/user/marketing/rewardList',
    );
    return UserMarketingRewardListResponseModel.fromJson(response);
  }

  ECPair getBitcoinKeyPair({required String mnemonicSeedPhrase}) {
    var seed = bip39.mnemonicToSeed(mnemonicSeedPhrase);
    HDWallet hdWallet = new HDWallet.fromSeed(seed, network: testnet);
    print('hdWallet.wif:${hdWallet.wif}');
    common.bitcoinWallet = hdWallet;
    common.bitcoinKeyPair = ECPair.fromWIF(hdWallet.wif!, network: testnet);
    return common.bitcoinKeyPair!;
  }

  Future<PushRawTransaction> pushRawTransaction(
      {required var alice,
      required SendBitcoinDetail sendBitcoinDetail}) async {
    final txb = new TransactionBuilder(network: testnet);
    txb.setVersion(1);
    if (sendBitcoinDetail.tx?.inputs != null) {
      sendBitcoinDetail.tx?.inputs!.forEach((InputData inputData) {
        txb.addInput(inputData.prevHash, inputData.outputIndex!);
      });
    }

    if (sendBitcoinDetail.tx?.outputs != null) {
      sendBitcoinDetail.tx?.outputs!.forEach((outPut) {
        txb.addOutput(outPut.addresses!.first, outPut.value);
      });
    }

    if (sendBitcoinDetail.tx?.inputs != null) {
      int i = 0;
      sendBitcoinDetail.tx?.inputs!.forEach((InputData inputData) {
        txb.sign(vin: i, keyPair: alice);
        i++;
      });
    }

    String transcationData = txb.build().toHex();
    print("txb.build().toHex():$transcationData}");

    Map<String, dynamic> response = await api.post(
        url:
            "https://api.blockcypher.com/v1/btc/test3/txs/push?token=e16ba72d5c85477fa99e9066936b98ea",
        emptyHeader: true,
        body: jsonEncode(
          {
            "tx": transcationData,
          },
        ));
    PushRawTransaction pushRawTransaction =
        PushRawTransaction.fromJson(response);
    return pushRawTransaction;
  }

  Future<BlockchainAddressInformation> getBlockchainAddressInformation(
      {required String address}) async {
    Map<String, dynamic> response = await api.get(
        url: 'https://api.blockcypher.com/v1/btc/test3/addrs/$address',
        haveApiAccessToken: false);
    BlockchainAddressInformation blockchainAddressInformation =
        BlockchainAddressInformation.fromJson(response);
    return blockchainAddressInformation;
  }

  Future<SendBitcoinDetail> getSendBitcoinDetail(
      {required String targetAddress,
      required myAddress,
      required int amount}) async {
    final json = jsonEncode({
      "inputs": [
        {
          "addresses": [myAddress]
        }
      ],
      "outputs": [
        {
          "addresses": [targetAddress],
          "value": amount
        }
      ]
    });
    Map<String, dynamic> response = await api.post(
        url:
            'https://api.blockcypher.com/v1/btc/test3/txs/new?token=e16ba72d5c85477fa99e9066936b98ea',
        haveApiAccessToken: false,
        body: json);
    SendBitcoinDetail sendBitcoinDetail = SendBitcoinDetail.fromJson(response);
    return sendBitcoinDetail;
  }

  //https://api.blockcypher.com/v1/btc/test3/addrs/mjaufUtcBUdZpx3AFxwjNX2fW3vQqhy5wD balance api and transcation api
}

AsymmetricKeyPair<PublicKey, PrivateKey> getRsaKeyPair(
    SecureRandom secureRandom) {
  /// Set BitStrength to [1024, 2048 or 4096]
  var rsapars = new RSAKeyGeneratorParameters(BigInt.from(65537), 1024, 5);
  var params = new ParametersWithRandom(rsapars, secureRandom);
  var keyGenerator = new RSAKeyGenerator();
  keyGenerator.init(params);
  return keyGenerator.generateKeyPair();
}

/// you can use Transaction
CallApiHandler api = CallApiHandler();
