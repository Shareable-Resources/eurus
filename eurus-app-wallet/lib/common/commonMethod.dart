import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'dart:io' show Platform;
import 'dart:typed_data';

import 'package:apihandler/apiHandler.dart';
import 'package:app_authentication_kit/mnemonic_kit.dart';
import 'package:app_crypto_icons/app_crypto_icons.dart';
import 'package:app_security_kit/app_security_kit.dart';
import 'package:app_storage_kit/app_storage_kit.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:collection/collection.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/bitcoinLibrary/bitcoin_flutter.dart';
import 'package:euruswallet/common/qrcode_scanner.dart';
import 'package:euruswallet/common/track_manager.dart';
import 'package:euruswallet/common/tx_record_model.dart';
import 'package:euruswallet/common/user_profile.dart';
import 'package:euruswallet/common/web3dart.dart';
import 'package:euruswallet/commonUI/acknowledgementDialog.dart';
import 'package:euruswallet/commonUI/authentication_pop_up_dialog.dart';
import 'package:euruswallet/commonUI/cameraFocus.dart';
import 'package:euruswallet/commonUI/customDialogBox.dart';
import 'package:euruswallet/commonUI/kycCamera.dart';
import 'package:euruswallet/model/adminFeeModel.dart';
import 'package:euruswallet/model/blockchainAddressInformation.dart';
import 'package:euruswallet/model/codeVerification.dart';
import 'package:euruswallet/model/coinPrice.dart';
import 'package:euruswallet/model/crypto_currency_model.dart';
import 'package:euruswallet/model/ethgasstation.dart';
import 'package:euruswallet/model/faucetRequest.dart';
import 'package:euruswallet/model/findEmailWalletAddress.dart';
import 'package:euruswallet/model/forgetLoginPw.dart';
import 'package:euruswallet/model/getFaucetList.dart';
import 'package:euruswallet/model/importWallet.dart';
import 'package:euruswallet/model/kYCCountryList.dart';
import 'package:euruswallet/model/loginBySignModel.dart';
import 'package:euruswallet/model/registerByEmail.dart';
import 'package:euruswallet/model/registerDevice.dart';
import 'package:euruswallet/model/userKYCStatus.dart';
import 'package:euruswallet/model/user_marketing_reward_list_response_model.dart';
import 'package:euruswallet/model/user_marketing_reward_scheme_response_model.dart';
import 'package:euruswallet/model/user_storage.dart';
import 'package:euruswallet/pages/afterScanPage.dart';
import 'package:euruswallet/pages/asset_allocation_token_list_page.dart';
import 'package:euruswallet/pages/centralized/top_up_payment_wallet_page.dart';
import 'package:euruswallet/pages/centralized/verifyCode.dart';
import 'package:euruswallet/pages/decentralized/add_token.dart';
import 'package:euruswallet/pages/decentralized/receivePage.dart';
import 'package:euruswallet/pages/decentralized/selectTargetPage.dart';
import 'package:euruswallet/pages/decentralized/transferPage.dart';
import 'package:euruswallet/pages/reward_detail_page.dart';
import 'package:euruswallet/pages/settingSubpages/KYC/kycStatus.dart';
import 'package:firebase_crashlytics/firebase_crashlytics.dart';
import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:flutter_jailbreak_detection/flutter_jailbreak_detection.dart';
import 'package:intl/intl.dart';
import 'package:livechat_inc/livechat_inc.dart';
import 'package:loader_overlay/src/overlay_controller_widget_extension.dart';
import 'package:package_info/package_info.dart';
import 'package:rxdart/rxdart.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:web3dart/crypto.dart' as _web3dartCrypto;
import 'package:web3dart/src/utils/typed_data.dart';

import '../commonUI/constant.dart';
import 'advance_gen_address_pair_args.dart';
import 'callApiHandler.dart';
import 'tx_record_handler.dart';

export 'dart:async';
export 'dart:convert';
export 'dart:io';
export 'dart:math';
export 'dart:typed_data';

export 'package:app_storage_kit/app_storage_kit.dart';
export 'package:app_storage_kit/normal_storage.dart';
export 'package:euruswallet/common/commonMethod.dart';
export 'package:euruswallet/common/web3dart.dart';
export 'package:euruswallet/commonUI/CustomConfirmDialog.dart';
export 'package:euruswallet/commonUI/SizeConfig.dart';
export 'package:euruswallet/commonUI/backgroundImage.dart';
export 'package:euruswallet/commonUI/customTextButton.dart';
export 'package:euruswallet/commonUI/submitButton.dart';
export 'package:euruswallet/commonUI/topCircularContainer.dart';
export 'package:euruswallet/commonUI/walletAppBar.dart';
export 'package:flutter/foundation.dart';
export 'package:flutter/material.dart';
export 'package:flutter/services.dart';
export 'package:flutter/widgets.dart';
export 'package:rounded_loading_button/rounded_loading_button.dart';

export '../commonUI/constant.dart';
export 'callApiHandler.dart';

const AutoFillAccount = true;
const TEST_NET = 'Rinkeby';
const haveFreeTokenButton = false;
const EURUS_CHAIN_CODE = 1;
const BOTTOM_NAV_TAB_BAR_HEIGHT = 50.0;
const SCREEN_WITH_BOTTOM_NAV_TAB_BAR_SAFE_AREA_BOTTOM_PADDING_CONTENT_INSET =
    BOTTOM_NAV_TAB_BAR_HEIGHT + 4.0;

enum TransactionType {
  transferInput,
  sendTransactionHistorySuccessfulStatus,
  allocationTransactionHistorySuccessfulStatus,
  sendTransactionHistoryPendingStatus,
  allocationTransactionHistoryPendingStatus,
  allocationTransactionHistoryProcessingStatus,
  confirmation,
  pending,
  allocationProcessing,
  successful,
  failure
}
enum TargetType { username, ethAddress, email }
enum CurrentUserType { centralized, decentralized }
enum DecenUserType { created, imported }
enum VerifyCodePageType {
  register,
  forgotLoginPw,
  forgotPaymentPw,
  newDeviceResetPublicKey
}
enum AddressPairType { loginPw, paymentPw }
enum ChangeLoginPasswordType { changeLoginPassword, resetLoginPassword }
enum ChangePaymentPasswordType { changePaymentPassword, forgetPaymentPassword }

class CommonMethod {
  static final CommonMethod _instance = CommonMethod._internal();
  TargetType? targetType;
  String? targetAddress;
  String? targetEmail;
  String? selectTokenSymbol;
  String? transferAmount;
  String? forgotPasswordEmailOrUserName;
  CurrentUserType currentUserType = CurrentUserType.centralized;
  BlockChainType fromBlockChainType = BlockChainType.Eurus;
  double? currentGas;
  bool transferToMySelf = false;
  String targetDepositOrWidthDrawAddresss = "";
  BlockChainType currentBlockchainSelection = BlockChainType.Eurus;
  BlockChainType backupCurrentBlockchainSelection = BlockChainType.Eurus;
  BlockChainType topSelectedBlockchainType = BlockChainType.Eurus;
  int currentSelectionSpeed = 0;
  Ethgasstation? gasStationData;
  BehaviorSubject<bool> refreshTopBar = BehaviorSubject<bool>();
  BehaviorSubject<bool>? logoutFunction;
  BehaviorSubject<bool> refreshAssetList = BehaviorSubject<bool>();
  Function? updateNavbarOnChangeLang;
  AdminFeeModel? adminFeeModel;
  late Timer timer;
  String? email;
  String? loginPassword;
  String? paymentPassword;
  GetFaucetList? getFaucetList;
  FaucetRequest? faucetRequest;
  RegisterByEmail? registerByEmail;
  ForgetLoginPw? forgetLoginPw;
  String? forgetPwEmail;
  CodeVerification? codeVerification;
  Map? cryptoCurrencyModelMap;
  bool? showFaucetBtn;
  AddressPair? loginAddressPair;
  AddressPair? serverAddressPair;
  String? serverMnemonic;
  String? rsaPublicKey;
  String? rsaPrivateKey =
      "MIICXQIBAAKBgQDU5YufylzObiWijgdmBfeZAyKrSOxq6Nrh+5Oh4crA/QQPkp8ZAcOoRvApgLLqQuUbfd7egWMOwczLID/yrA0Wi3k+Tk9Z7z4SfrPxaWu1+elHBktzLhQLkf4bj2PVAJH1zePJIy4PJCIO8gFelstEHFvWcso2ePA1nnZ6knunJQIDAQABAoGAYtUNNGjlHI/VuNjmZl5uywHBnnKEDj17H12C86u2TFEpCXGvmhRPmFcWNq4gYNAdO937EKBQNBGT2Nhn12g3yl4+4edPkf39URg3ZwCq4uxAWQ9z5rwPa1eOOAok2VNrJd1NX4WEUerXRgzT1MU449jaZkGz4m8LEiKYKEIW+CECQQDZ2hQy80HXM/gJVIaYN2khzW7Q5LigrVzGf1SJbXGb10XfSbm1r0rN6DSa2haSfChb/lj8eN1xjXMaxiZTsRfbAkEA+i1VRWtdHNMdTnXW/gU3cfbGxo7/nj7WUU4LaZ2AvZ12K/dRkwZYPEyxlId2HnJwnw8TRX8FIR8H3mkV5HXs/wJAcmOcL5Sjgch7+Qo1EkAmJ+WixnUSrOvaxy+cx/x7pwTGX5RquwesE6pV1Omm6Ivg9Uz8lLUyManAQtLA1Tkr+QJBAI5BzOUmgdHsMhP1agUTzk1dd/ZcRfoj3RZqfI7X4ubvbMzfW2FxECdprOi6hm4VwPiRR/ISokYNMRpFQw+gBt0CQQCzxWkD/jeGCI9F4qYI9qwgT8lxce702Y2D9huaQ3sxc9NiD/Bhm1Nw8P8W2tgXENh9BV9xcq/gEsHu7OaidSfF";
  VerifyCodePageType? verifyCodePageType;
  String? cenUserWalletAddress;
  String? cenMainNetWalletAddress;
  FindEmailWalletAddress? findEmailWalletAddress;
  ImportWallet? importWallet;
  bool isCenWithdraw = false;
  String? ownerWalletAddress;
  String? cenSignKey;
  String? sendAmount;
  String? _address;
  String? get currentAddress => _address;
  set currentAddress(String? newValue) {
    _address = newValue;
    // hacky trick to enforce override web3dart previously initialized member variable of the shared instance
    // by invoking getbalance via the web3dart shared instance per each wallet public address once newly set through this address setter
    getERC20BalanceAndInit(newValue, '0x0');
  }

  String? encryptedAddress;
  bool showBalance = true;
  String? assetsList;
  List coingeckoCoinsList = [];
  CoinPriceList? coinPriceList;
  RegisterDevice? registerDevice;
  ChangePaymentPasswordType? changePaymentPasswordType;
  ForgetLoginPw? forgotPaymentPw;
  UserKYCStatus? userKYCStatus;
  KYCCountryList? kYCCountryList;
  late BuildContext currentContext;
  String eurusTXType = "0";
  BlockchainAddressInformation blockchainAddressInformation =
      BlockchainAddressInformation();
  ECPair? bitcoinKeyPair = ECPair.fromWIF(
      'cNR1r2WgKiTdTd1VwhvBbYbJcsERDYyYXzA2FxWtGNxcNWjjZV2k',
      network: testnet);
  HDWallet? bitcoinWallet;
  UserStorage userStorage = UserStorage();
  List<UserMarketingRewardScheme>? rewardSchemes;
  List<UserMarkingRewardList>? rewardedList;

  UserProfile? currentUserProfile;

  /// init method
  CommonMethod._internal();

  factory CommonMethod() {
    return _instance;
  }

  static String? passwordEncrypt(String password, String value) =>
      PasswordEncryptHelper(password: password).encryptWPwd(value);
  static String? passwordDecrypt(String password, String value) =>
      PasswordEncryptHelper(password: password).decryptWPed(value);

  Future<String> get prefix async {
    final _prefixKey = 'APP_INSTALL_FIRST_LAUNCH_ENTER_WALLET';
    String? _prefix = await NormalStorageKit().readValue(_prefixKey);
    if (isEmptyString(string: _prefix)) {
      _prefix = DateTime.now().microsecondsSinceEpoch.toString();
      await NormalStorageKit().setValue(_prefix, _prefixKey);
    }
    return _prefix ?? '';
  }

  String? getBannerPath({required BuildContext context}) {
    final locale = context.locale.toString();
    var suffix = 'en';
    if (locale == 'zh_Hans') {
      suffix = 'cn';
    } else if (locale == 'zh_Hant') {
      suffix = 'tw';
    }

    if (shouldShowReward == null) return null;

    if (shouldShowReward!) {
      return 'images/reward_banner_type_0_$suffix.png';
    } else {
      String imageName = 'banner1';
      if (isCentralized()) imageName += '_centralized';
      return 'images/${imageName}_$suffix.png';
    }
  }

  void startTimer() {
    timer = Timer.periodic(Duration(seconds: 20), (timer) {
      api.reFreshToken();
    });
  }

  bool isTransactionHistory({required TransactionType transactionType}) {
    return transactionType ==
            TransactionType.sendTransactionHistorySuccessfulStatus ||
        transactionType ==
            TransactionType.allocationTransactionHistorySuccessfulStatus ||
        transactionType ==
            TransactionType.sendTransactionHistoryPendingStatus ||
        transactionType ==
            TransactionType.allocationTransactionHistoryProcessingStatus;
  }

  Future<String?> findDefaultAccountEncryptedAddress() async {
    final account =
        await SecureStorageKit().readValue(await common.prefix + 'accounts');

    if (isEmptyString(string: account)) {
      return null;
    }

    final encodedComponent =
        RegExp(r'([^=&?]*)=0').firstMatch(account ?? '')?.group(1);
    return !isEmptyString(string: encodedComponent)
        ? Uri.decodeComponent(encodedComponent ?? '')
        : null;
  }

  // TODO - should push an Auth Page onto upfront when AppLifecycleState or LifecycleReactor or LifecycleEventHandler or WidgetsBindingObserver onResumed
  Future showAuthPopUp(
    BuildContext context, {
    required CurrentUserType currentUserType,
    bool isNeededRoutingToHomePage = false,
    String? displayUserName,
    String? encryptedAddress,
  }) async {
    EasyLoading.show(status: 'COMMON.LOADING_W_DOT'.tr());

    final _encryptedAddress =
        encryptedAddress ?? await findDefaultAccountEncryptedAddress();

    // final userProfiles = await getLocalAcs();
    List<UserProfile> userProfiles = await getLocalAcs();
    UserProfile? userProfile = userProfiles.firstWhereOrNull(
        (element) => element.encryptedAddress == _encryptedAddress);
    if (userProfile == null) {
      EasyLoading.dismiss();
      return;
    }

    String? _displayUserName = displayUserName ??
        await NormalStorageKit().readValue('displayUserName');
    if (userProfile.seedPraseBackuped == false) {
      _displayUserName = getUserAddress(userProfile);
    }

    if (isEmptyString(string: _encryptedAddress)) {
      EasyLoading.dismiss();
      return;
    }

    final List<dynamic>? result = await showGeneralDialog(
      context: context,
      pageBuilder: (BuildContext context, Animation<double> animation,
              Animation<double> secondaryAnimation) =>
          AuthenticationPopUpDialog(
        currentUserType: currentUserType,
        displayUserName: _displayUserName,
        encryptedAddress: _encryptedAddress ?? '',
      ),
    );

    if (result == null) {
      EasyLoading.dismiss();
      if (!isNeededRoutingToHomePage) logoutFunction?.add(true);
      return;
    }

    AddressPair? addressPair;
    int passwordIndex = result.indexWhere((element) => element is String);
    String? password =
        passwordIndex != -1 ? result.elementAt(passwordIndex) : null;
    String? privateKey;
    LoginBySignModel? loginBySignModel;

    if (currentUserType == CurrentUserType.centralized) {
      int loginBySignModelIndex =
          result.indexWhere((element) => element is LoginBySignModel);
      loginBySignModel = loginBySignModelIndex != -1
          ? result.elementAt(loginBySignModelIndex)
          : null;
      privateKey = null;
    } else {
      int addressPairIndex =
          result.indexWhere((element) => element is AddressPair);
      addressPair =
          addressPairIndex != -1 ? result.elementAt(addressPairIndex) : null;
      if (addressPair == null || isEmptyString(string: password)) {
        Navigator.popUntil(context, (route) => route.isFirst);
        EasyLoading.dismiss();
        return;
      }
      currentAddress = addressPair.address.toLowerCase();
      privateKey = addressPair.privateKey;

      await updateApiAccessTokenByImportWallet(
        context: context,
        addressPair: addressPair,
      );
    }

    if (isNeededRoutingToHomePage)
      await common.successMoveToHomePage(
        userType: currentUserType,
        email: await NormalStorageKit().readValue('cenUserEmail') ?? '',
        context: context,
        loginPassword: password ?? '',
        address: currentAddress,
        privateKey: privateKey,
        loginBySignModel: loginBySignModel,
      );

    EasyLoading.dismiss();
  }

  Future<bool> updateApiAccessTokenByImportWallet({
    required BuildContext context,
    required AddressPair addressPair,
  }) async {
    ImportWallet importWallet =
        await api.importWallet(anAddressPair: addressPair);

    if (await checkApiError(
        context: context,
        errorString: importWallet.message,
        returnCode: importWallet.returnCode)) {
      String? token = importWallet.data == null
          ? importWallet.token
          : importWallet.data?.token;
      String? expiryTime = importWallet.data == null
          ? importWallet.expiryTime.toString()
          : importWallet.data?.expiryTime.toString();
      NormalStorageKit().setValue(token ?? '', 'apiAccessToken_');
      NormalStorageKit()
          .setValue(expiryTime ?? '', 'apiAccessTokenExpiryTime_');
      NormalStorageKit()
          .setValue("${getCurrentTimeStamp().toString()}", 'currentTimeStamp');
      return true;
    } else {
      return false;
    }
  }

  Future<List<Map<String, dynamic>>> fetchTransactionHistory(
      String from, int chainCode,
      [String? toContractAddress]) async {
    // 0x0 chain root asset, i.e. ETH in Ethereum main-net or test-net, and EUN in Eurus side-chain, transaction history recordset
    // if (toContractAddress == null || toContractAddress.isEmpty || toContractAddress == '0x0') return await TxRecordHandler().readTxs(where: 'txFrom = ? AND chain = ? AND txTo = ?', whereArgs: []);
    if (isEmptyString(string: toContractAddress)) {
      // final recordset =  (await TxRecordHandler().readTxs(where: '(txFrom LIKE ? OR decodedInputRecipientAddress LIKE ? OR txTo LIKE ?) AND chain = ? AND txTo LIKE decodedInputRecipientAddress', whereArgs: ["%${from}%", "%${from}%", "%${from}%", chainCode == 0 ? 'eth' : 'eun']).catchError((e) {print("TxRecordHandler().readTxs().catchError() ${e}"); return null;}))?.map((e) => e.toJson())?.toList();
      final recordset = chainCode == 0
          ? (await TxRecordHandler()
                  .readMainnetTxs(from, isPlatformToken: true)
                  .catchError((e) {
              print("TxRecordHandler().readTxs().catchError() $e");
              return <TxRecordModel>[];
            }))
              .map((e) => e.toJson())
              .toList()
          : (await TxRecordHandler()
                  .readSidechainTxs(from, '', isPlatformToken: true)
                  .catchError((e) {
              print("TxRecordHandler().readTxs().catchError() $e");
              return <TxRecordModel>[];
            }))
              .map((e) => e.toJson())
              .toList();
      print(recordset.isNotEmpty
          ? 'recordset[0] ${recordset[0]}'
          : 'recordset is list but empty list');
      print("recordset for ${[
        from,
        chainCode == 0 ? 'eth' : 'eun',
        toContractAddress
      ]} is ${recordset.length} $recordset");
      return recordset;
    }

    // ERC20 transaction history recordset
    // final recordset = await TxRecordHandler().readTxs(where: 'txFrom = ? AND chain = ? AND txTo = ?', whereArgs: [from, chainCode == 0 ? 'eth' : 'eun', toContractAddress]).then((value) => value.map((e) => e.toJson()).toList()??jsonEncode(value));
    // final recordset = await TxRecordHandler().readTxs().then((value) => value.map((e) => e.toJson()).toList());
    // final recordset = await TxRecordHandler().readTxs().then((value) => value.map((e) => e.toJson())).then((value) => value.toList());
    // final recordset = (await TxRecordHandler().readTxs()).map((e) => e.toJson()).toList();
    // final recordset = (await TxRecordHandler().readTxs(where: '(txFrom LIKE ? OR decodedInputRecipientAddress LIKE ?) AND chain = ? AND txTo LIKE ?', whereArgs: ["%${from}%", "%${from}%", chainCode == 0 ? 'eth' : 'eun', "%$toContractAddress%"]).catchError((e) {print('TxRecordHandler().readTxs().catchError() ${e}'); return null;}))?.map((e) => e.toJson())?.toList(); // option 1: db uses COLLATE NOCASE // option 2: uses rawQuery for COLLATE NOCASE // uses where clause LIKE %% // uses toLowerCase in both where clause and where args
    final recordset = chainCode == 0
        ? (await TxRecordHandler()
                .readMainnetTxs(from, contractAddress: toContractAddress)
                .catchError((e) {
            print("TxRecordHandler().readTxs().catchError() $e");
            return <TxRecordModel>[];
          }))
            .map((e) => e.toJson())
            .toList()
        : (await TxRecordHandler()
                .readSidechainTxs(from, '', contractAddress: toContractAddress)
                .catchError((e) {
            print("TxRecordHandler().readTxs().catchError() $e");
            return <TxRecordModel>[];
          }))
            .map((e) => e.toJson())
            .toList();
    // final recordset = (await TxRecordHandler().readTxs()).map((e) => e.toJson()).where((e) => e['txFrom'] == from && (e['chain'] == (chainCode == 0 ? 'eth' : 'eun')) && e['txTo'] == toContractAddress).toList();
    // final recordset = (await TxRecordHandler().readTxs(where: '1')).where((e) => e.txFrom == from && e.chain == (chainCode == 0 ? 'eth' : 'eun') && e.txTo == toContractAddress).map((e) => e.toJson()).toList();
    // final recordset = (await TxRecordHandler().readTxs(where: '1')).where((e) => e.chain == (chainCode == 0 ? 'eth' : 'eun') && e.txTo == toContractAddress).map((e) => e.toJson()).toList(); // issues root cause identified: from address (with checksum encoded) upper-lower case (case sensitive) v.s. db recordset txFrom address (without checksum encoding) lower case only (case insensitive)
    print(recordset.isNotEmpty
        ? 'recordset[0] ${recordset[0]}'
        : 'recordset is list but empty list');

    //// testing without .toJson()
    ////
    // final recordset = (await TxRecordHandler().readTxs(where: '1', whereArgs: ['true'])).map((e) => e).toList(); // work if without .toJson()
    // print('recordset[0]?.toJson() ${recordset[0]?.toJson()}');
    ////

    print("recordset for ${[
      from,
      chainCode == 0 ? 'eth' : 'eun',
      toContractAddress
    ]} is ${recordset.length} $recordset");
    return recordset;
  }

  double divisionDecimal(int decimal, int number) {
    double divisionNumber = 1.0;
    for (var i = 0; i < decimal; i++) {
      divisionNumber = divisionNumber * 10;
    }
    return number.toDouble() / divisionNumber;
  }

  String targetDepositOrWidthDrawAddress() {
    targetDepositOrWidthDrawAddresss = "";
    if (common.isDeposit()) {
      targetDepositOrWidthDrawAddresss =
          'TX_PAGE.MY_WALLET'.tr(args: ["Eurus"]);
    } else if (common.isWithdraw()) {
      targetDepositOrWidthDrawAddresss =
          'TX_PAGE.MY_WALLET'.tr(args: ["Ethereum"]);
    }
    return targetDepositOrWidthDrawAddresss;
  }

  String numberFormat({
    String? number,
    int? maxDecimal,
    int? minDecimal,
  }) {
    if (isEmptyString(string: number)) {
      return 0.toStringAsFixed(minDecimal ?? 8);
    }
    if (double.tryParse(number ?? '') == null ||
        double.parse(number ?? '') == 0)
      return 0.toStringAsFixed(minDecimal ?? 8);

    minDecimal = minDecimal == null ? 0 : minDecimal;
    var numberFormat = new NumberFormat("#,##0.00", "en_US");
    numberFormat.minimumFractionDigits = minDecimal;
    numberFormat.maximumFractionDigits = maxDecimal == null ? 8 : maxDecimal;

    return numberFormat.format(double.parse(number ?? ''));
  }

  /// use this method can use loading on pushing page
  Future pushPage({
    required Widget page,
    required BuildContext context,
  }) async {
    return Navigator.push(
      context,
      new MaterialPageRoute(
          builder: (context) => FlutterEasyLoading(child: page)),
    );
  }

  /// use this method can use loading on pushing page
  Future pushReplacementPage({
    required Widget page,
    required BuildContext context,
    RouteSettings? settings,
  }) async {
    final completer = Completer();
    final result = Navigator.pushReplacement(
        context,
        PageRouteBuilder(
            settings: settings,
            fullscreenDialog: true,
            opaque: false,
            pageBuilder: (pageBuilderContext, animation, secondaryAnimation) =>
                FlutterEasyLoading(child: page)),
        result: completer.future);
    return completer.complete(result);
  }

  bool isEthOrEun() {
    return common.fromBlockChainType == BlockChainType.Ethereum &&
            common.selectTokenSymbol == "ETH" ||
        common.fromBlockChainType == BlockChainType.Eurus &&
            common.selectTokenSymbol == "EUN";
  }

  bool isBSC() {
    return common.fromBlockChainType == BlockChainType.BinanceCoin &&
        common.selectTokenSymbol == "BNB";
  }

  Future<bool> getBalance() async {
    web3dart.erc20TokenBalanceFromEurus = await web3dart.getERC20Balance(
        blockChainType: BlockChainType.Eurus,
        deployedContract: web3dart.erc20ContractFromEurus);
    web3dart.erc20TokenBalanceFromEthereum = await web3dart.getERC20Balance(
        blockChainType: BlockChainType.Ethereum,
        deployedContract: web3dart.erc20ContractFromEthereum);
    web3dart.ethBalanceFromEurus =
        await web3dart.getETHBalance(blockChainType: BlockChainType.Eurus);
    web3dart.ethBalanceFromEthereum =
        await web3dart.getETHBalance(blockChainType: BlockChainType.Ethereum);
    return true;
  }

  Future<String> getERC20BalanceAndInit(publicAddress, contractAddress,
      [blockChainType = BlockChainType.Ethereum]) async {
    print('will web3dart.initEthClient()');
    await web3dart.initEthClient(publicAddress: publicAddress);
    print('contractAddress: $contractAddress');
    if (contractAddress == '0x0') {
      final _canGet0x0Balance =
          ([blockChainType = BlockChainType.Ethereum]) async {
        return await web3dart.getETHBalance(blockChainType: blockChainType);
      };
      return await _canGet0x0Balance(blockChainType);
    }
    DeployedContract deployedContract;
    if (blockChainType == BlockChainType.BinanceCoin) {
      deployedContract =
          await web3dart.getBSCTokenContract(contractAddress: contractAddress);
    } else {
      deployedContract =
          web3dart.getEurusERC20Contract(contractAddress: contractAddress);
    }
    final _ = await web3dart.getERC20Balance(
        blockChainType: blockChainType, deployedContract: deployedContract);
    print('contractAddress: $contractAddress balance: $_');
    return common.numberFormat(number: _);
  }

  Future<bool> setUpErc20TokenContract({
    String? ethereumErc20ContractAddress,
    String? eurusErc20ContractAddress,
    BlockChainType? fromBlockChainType,
  }) async {
    if (ethereumErc20ContractAddress != null &&
        fromBlockChainType == BlockChainType.Ethereum) {
      if (ethereumErc20ContractAddress == "0x0") {
        common.selectTokenSymbol = "ETH";
      } else {
        web3dart.setErc20Contract(
            blockChainType: BlockChainType.Ethereum,
            contractAddress: ethereumErc20ContractAddress);
      }
    }
    if (eurusErc20ContractAddress != null &&
        fromBlockChainType == BlockChainType.Eurus) {
      if (eurusErc20ContractAddress == "0x0") {
        common.selectTokenSymbol = "EUN";
      } else {
        web3dart.setErc20Contract(
            blockChainType: BlockChainType.Eurus,
            contractAddress: eurusErc20ContractAddress);
      }
    }
    if (fromBlockChainType != null) {
      common.fromBlockChainType = fromBlockChainType;
      bool emptyContract = false;
      if (fromBlockChainType == BlockChainType.Ethereum) {
        if (isEmptyString(string: ethereumErc20ContractAddress)) {
          emptyContract = true;
        }
      } else if (fromBlockChainType == BlockChainType.Eurus) {
        if (isEmptyString(string: ethereumErc20ContractAddress)) {
          emptyContract = true;
        }
      }
      common.selectTokenSymbol = emptyContract
          ? common.selectTokenSymbol
          : await web3dart.getTokenSymbol(
              blockChainType: common.fromBlockChainType,
              deployedContract:
                  common.fromBlockChainType == BlockChainType.Ethereum
                      ? web3dart.erc20ContractFromEthereum
                      : web3dart.erc20ContractFromEurus);
      print("selectTokenSymbol:$selectTokenSymbol");
    }

    return true;
  }

  Future<Map<String, CryptoCurrencyModel>> getSupportedTokens() async {
    var web3dart = Web3dart();

    await web3dart.getERC20TokenList(blockChainType: BlockChainType.Ethereum);
    var tempEthMap = web3dart.tokenListMap;

    Map<String, CryptoCurrencyModel> tokensInMap = {
      'ETH': CryptoCurrencyModel(
        currency: 'Ethereum',
        symbol: 'ETH',
        showAssets: true,
        supported: true,
        addressRinkeby: '0x0',
        addressEthereum: '0x0',
      ),
      'EUN': CryptoCurrencyModel(
        currency: 'Eurus',
        symbol: 'EUN',
        showAssets: true,
        supported: true,
        addressRinkeby: null,
        addressEurus: '0x0',
        addressEthereum: null,
        addressBSC: null,
      ),
      // 'BNB': CryptoCurrencyModel(
      //   currency: 'Binance Coin',
      //   symbol: 'BNB',
      //   showAssets: true,
      //   supported: true,
      //   addressBSC: '0x0',
      // ),
    };

    tempEthMap.forEach((key, value) {
      if (tokensInMap[key] != null) {
        tokensInMap[key]!.addressEthereum = value;
        tokensInMap[key]!.addressEurus = value;
      } else {
        CryptoCurrencyModel c = CryptoCurrencyModel(
          currency: getTokenNameBySymbol(key),
          symbol: key,
          showAssets: true,
          supported: true,
          addressRinkeby: value,
          addressEthereum: value,
        );
        tokensInMap.addAll({key: c});
      }
    });

    await web3dart.getERC20TokenList(blockChainType: BlockChainType.Eurus);
    var tempEurMap = web3dart.tokenListMap;

    tempEurMap.forEach((key, value) {
      if (tokensInMap[key] != null) {
        tokensInMap[key]!.addressEurus = value;
      } else {
        CryptoCurrencyModel c = CryptoCurrencyModel(
          currency: getTokenNameBySymbol(key),
          symbol: key,
          showAssets: true,
          supported: true,
          addressEurus: value,
          addressEthereum: null,
          addressRinkeby: null,
        );
        tokensInMap.addAll({key: c});
      }
    });

    // var tempBscMap = web3dart.bscTokenListMap;
    // tempBscMap.forEach((key, value) {
    //   if (tokensInMap[key] != null) {
    //     tokensInMap[key]!.addressBSC = value;
    //   } else {
    //     CryptoCurrencyModel c = CryptoCurrencyModel(
    //       currency: getTokenNameBySymbol(key),
    //       symbol: key,
    //       showAssets: true,
    //       supported: true,
    //       addressBSC: value,
    //     );
    //     tokensInMap.addAll({key: c});
    //   }
    // });

    return tokensInMap;
  }

  String getTokenNameBySymbol(String s) {
    Map<String, String> mappedCurrency = {
      'ETH': 'Ethereum',
      'EUN': 'Eurus',
      'USDT': 'Tether USD',
      'USDC': 'USD Coin',
      'LINK': 'ChainLink Token',
      'UNI': 'Uniswap',
      'BNB': 'BNB',
      'BUSD': 'Binance USD',
      'YFI': 'yearn.finance',
      'DAI': 'Dai Stablecoin',
      'OMG': 'OMG Network',
      'VEN': 'VeChain',
      'AAVE': 'Aave Token',
      'HT': 'HuobiToken',
      'SUSHI': 'SushiToken',
      'TUSD': 'TrueUSD',
      'cDAI': 'Compound Dai',
      'SXP': 'Swipe',
      'BAT': 'BAT',
      'USDK': 'USDK',
      'WBTC': 'Wrapped BTC',
      'ZIL': 'Zilliqa',
      'SNX': 'Synthetix Network Token',
      'OKB': 'OKB',
      'BAND': 'BandToken',
      'MKR': 'Maker',
      'HUSD': 'HUSD',
      'ZRX': 'ZRX',
      'PAX': 'Paxos Standard',
      'COMP': 'Compound',
      'RSR': 'Reserve Rights',
      'BAL': 'Balancer',
    };
    return mappedCurrency[s] ?? s;
  }

  void showSnackBar({
    required String errorMessage,
    required BuildContext context,
  }) {
    final snackBar = SnackBar(
        padding: EdgeInsets.zero,
        elevation: 0,
        content: Padding(
          padding: EdgeInsets.only(bottom: 66),
          child: Container(
              child: Center(
                  child: Text(errorMessage, textAlign: TextAlign.center)),
              color: Colors.redAccent.shade100.withOpacity(.9),
              height: 60),
        ),
        backgroundColor: Colors.transparent);
    ScaffoldMessenger.of(context).showSnackBar(snackBar);
  }

  void showCopiedToClipboardSnackBar(
    BuildContext context, {
    EdgeInsetsGeometry? margin,
    SnackBarBehavior behavior = SnackBarBehavior.fixed,
  }) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(
          'COMMON.COPIED'.tr(),
          textAlign: TextAlign.center,
          style: FXUI.titleTextStyle.copyWith(
            fontSize: 16,
            color: Colors.white,
          ),
        ),
        backgroundColor: FXColor.textGray,
        behavior: behavior,
        margin: margin,
        duration: Duration(milliseconds: 700),
        shape: RoundedRectangleBorder(),
        padding: EdgeInsets.symmetric(vertical: 10),
      ),
    );
  }

  bool isWithdraw() {
    return (common.transferToMySelf &&
            common.fromBlockChainType == BlockChainType.Eurus &&
            !isCentralized()) ||
        (common.isCenWithdraw && isCentralized());
  }

  bool isDeposit() {
    return common.transferToMySelf &&
        common.fromBlockChainType == BlockChainType.Ethereum;
  }

  String? getUriVal(String? uri, String valKey) {
    return Uri(query: uri).queryParameters[valKey];
  }

  String updateUriVal(
    String? uri,
    String valKey,
    String? val, {
    bool? addNewIfNotFound,
    String? orgVal,
  }) {
    if ((addNewIfNotFound ?? true) == false) {
      var orgVal = getUriVal(uri, valKey);
      if (isEmptyString(string: orgVal) && !isEmptyString(string: uri))
        return uri ?? '';
    }

    /// 1.add new key and value , if key exist update val(may be 0)
    var uriParam = {
      ...Uri(query: uri).queryParameters,
      valKey: val ?? '',
    };

    if (isEmptyString(string: val)) uriParam.remove(valKey);

    if (orgVal != null) {
      uriParam.forEach((key, value) {
        /// key != valKey not change 1. add new key and value, change other value to 1
        if (key != valKey && value == (val ?? '')) uriParam[key] = orgVal;
      });
    }

    return Uri(queryParameters: uriParam).query;
  }

  Future<AddressPair> getAddressPair({
    String? email,
    String? password,
    String? mnemonic,
    AddressPairType addressPairType: AddressPairType.loginPw,
  }) async {
    String pw = "${email?.toLowerCase() ?? ''}${password ?? ''}";
    mnemonic = mnemonic;
    if (isEmptyString(string: mnemonic)) {
      switch (envType) {
        case EnvType.Dev:
          mnemonic =
              'carbon shuffle shoot knock alter bottom polar maple husband poet match spring';
          break;
        case EnvType.Staging:
        case EnvType.Testnet:
          mnemonic =
              'buddy build text drill aisle stone robot fringe duty mother assault please';
          break;
        case EnvType.Production:
          mnemonic =
              'dove clock chalk front spike prefer people spike capable word gasp congress';
          break;
      }
    }
    AdvanceGenAddressPairArgs args =
        AdvanceGenAddressPairArgs(mnemonic: mnemonic ?? '', pw: pw);
    AddressPair adPair;
    if (addressPairType == AddressPairType.loginPw) {
      adPair = advanceGenAddressPair(args);
    } else {
      adPair = advanceGenAddressPairPaymentPw(args);
    }

    return adPair;
  }

  AddressPair getAddressPairFromPrivateKey(String privateKey) {
    final ethPrivateKey = EthPrivateKey.fromHex(privateKey);
    final encodedPrivateKey = _web3dartCrypto.hexToBytes(privateKey);
    final encodedPublicKey =
        _web3dartCrypto.privateKeyBytesToPublic(encodedPrivateKey);
    var uncompressedPublicKey = encodedPublicKey.toList();
    // MARK: - compressPublicKey need prefix 04
    uncompressedPublicKey.insert(0, 4);
    final compressedPublicKey = _web3dartCrypto
        .compressPublicKey(Uint8List.fromList(uncompressedPublicKey));
    final publicKey = _web3dartCrypto.bytesToHex(compressedPublicKey);
    return AddressPair(
      ethPrivateKey.address.toString(),
      privateKey,
      publicKey: publicKey,
    );
  }

  Future<bool> checkApiError({
    String? apiName,
    required BuildContext context,
    String? errorString,
    int? returnCode,
    Color? btnColor,
  }) async {
    EasyLoading.dismiss();
    if (returnCode == 0) {
      return true;
    } else {
      Function? btnHandler;
      if (errorString?.toLowerCase() == "user not found") {
        if (apiName == "findEmailWalletAddress") {
          errorString = "COMMON_ERROR.USER_NOT_FOUND2".tr();
        } else {
          errorString = "COMMON_ERROR.USER_NOT_FOUND".tr();
        }
      }

      if (errorString?.toLowerCase() == "invalid verification code") {
        errorString = "COMMON_ERROR.INVALID_CODE".tr();
      }

      if (errorString?.toLowerCase() == "code expired") {
        errorString = "COMMON_ERROR.CODE_EXPIRED".tr();
      }

      if (errorString?.toLowerCase() == "server maintenance") {
        errorString = "COMMON_ERROR.SERVER_MAINTENANCE".tr();
        btnHandler = () {
          if (common.logoutFunction != null) {
            common.logoutFunction?.add(true);
          } else {
            Navigator.popUntil(context, (route) => route.isFirst);
          }
        };
      }

      if (errorString?.toLowerCase() ==
          "request too frequent, wait for 5 minute") {
        errorString = "COMMON_ERROR.REQUEST_TOO_FREQUENT".tr();
      }

      if (errorString?.toLowerCase() ==
          "request too frequent, wait for 1 minute") {
        errorString = "COMMON_ERROR.REQUEST_TOO_FREQUENT_1_MIN".tr();
      }

      if (errorString?.toLowerCase() == "old login address not match" ||
          errorString?.toLowerCase() ==
              "old owner wallet address not does match with user wallet owner wallet address") {
        errorString = "COMMON_ERROR.PW_IS_INCORRECT".tr();
      }

      if (errorString?.toLowerCase() == "cannot get user by email") {
        errorString = "COMMON_ERROR.USER_EMAIL_NOT_FOUND".tr();
      }

      await showDialog(
        context: context,
        builder: (BuildContext context) {
          return CustomDialogBox(
            btnColor: btnColor,
            descriptions: errorString ?? "Error",
            buttonText: "COMMON.OK".tr(),
            btnHandler: btnHandler,
          );
        },
      );

      return false;
    }
  }

  Color getBackGroundColor() {
    return !isCentralized() ? FXColor.mainDeepBlueColor : FXColor.mainBlueColor;
  }

  void showPopUpError(
      {required BuildContext context,
      Color? btnColor,
      required String descriptions}) {
    Future.delayed(Duration.zero, () {
      showDialog(
          context: context,
          builder: (BuildContext context) {
            return CustomDialogBox(
              btnColor: btnColor ?? common.getBackGroundColor(),
              //  title: 'DISCLAIMER.DISCLAIMER_TITLE'.tr(),
              descriptions: descriptions,
              buttonText: "COMMON.OK".tr(),
              btnHandler: () {
                Navigator.of(context).popUntil((route) => route.isFirst);
              },
            );
          });
    });
  }

  void showBiometricNotEnable({
    required BuildContext context,
    Color? btnColor,
    CurrentUserType? currentUserType,
  }) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return CustomDialogBox(
          btnColor: (currentUserType ?? common.currentUserType) ==
                  CurrentUserType.centralized
              ? FXColor.mainBlueColor
              : FXColor.mainDeepBlueColor,
          title: 'BIOMETRIC_ERROR.TIPS'.tr(),
          descriptions: 'BIOMETRIC_ERROR.BIOMETRIC_ENABLE_TIPS'.tr(),
          buttonText: "COMMON.OK".tr(),
        );
      },
    );
  }

  Future<bool> forceAppUpdateIfNeeded(BuildContext context) async {
    final clientVersion = await api.clientVersion();
    return await PackageInfo.fromPlatform().then((package) {
      final currentAppVersion = package.version;
      final minimumAppVersion = Platform.isIOS
          ? clientVersion.data?.iPhoneMinimumVersion
          : Platform.isAndroid
              ? clientVersion.data?.androidMinimumVersion
              : null;
      final updateLink = Platform.isIOS
          ? envType == EnvType.Production
              ? 'https://apps.apple.com/us/app/eurus-wallet/id1585783815'
              : 'https://apps.apple.com/us/app/eurus-testnet-wallet/id1566514594'
          : Platform.isAndroid
              ? envType == EnvType.Production
                  ? 'https://play.google.com/store/apps/details?id=me.mobi.eurusWallet'
                  : envType == EnvType.Testnet
                      ? 'https://play.google.com/store/apps/details?id=wallet.euruswallet'
                      : 'https://drive.google.com/file/d/1Og3kOIAqzQXQcWFf7myhybd7iK8Y2Ipy/view?usp=sharing'
              : null;
      if (minimumAppVersion != null &&
          (getExtendedVersionNumber(currentAppVersion) <
              getExtendedVersionNumber(minimumAppVersion))) {
        showDialog(
          context: context,
          barrierDismissible: false,
          builder: (_) {
            final desc = Platform.isIOS
                ? 'FORCE_UPDATE.DESC_IOS'.tr()
                : Platform.isAndroid
                    ? 'FORCE_UPDATE.DESC_ANDROID'.tr()
                    : '';
            return WillPopScope(
              onWillPop: () async => false,
              child: CustomDialogBox(
                dismissIble: false,
                title: 'FORCE_UPDATE.TITLE'.tr(),
                descriptions: desc,
                buttonText: 'FORCE_UPDATE.BTN_TITLE'.tr(),
                btnHandler: () async {
                  updateLink != null && await canLaunch(updateLink)
                      ? await launch(updateLink)
                      : throw 'Could not launch $updateLink';
                },
              ),
            );
          },
        );
        return true;
      }
      return false;
    });
  }

  int getExtendedVersionNumber(String version) {
    // Note that if you want to support bigger version cells than 99,
    // just increase the returned versionCells multipliers
    List versionCells = version.split('.');
    versionCells = versionCells.map((i) => int.parse(i)).toList();
    return versionCells[0] * 10000 + versionCells[1] * 100 + versionCells[2];
  }

  void alreadyRoot(BuildContext context) async {
    if (await FlutterJailbreakDetection.jailbroken) {
      showDialog(
        context: context,
        barrierDismissible: false,
        builder: (_) {
          return WillPopScope(
            onWillPop: () async => false,
            child: CustomDialogBox(
              dismissIble: false,
              descriptions: 'COMMON.ALREADY_ROOT_DEVICE'.tr(),
              buttonText: 'COMMON.OK'.tr(),
              btnHandler: () async {},
            ),
          );
        },
      );
    }
  }

  bool isNotContainDigitAndCharacter({required String text}) {
    final alphanumeric = RegExp(r'^(?=.*?[a-zA-Z])(?=.*?[0-9])');
    bool loginPwNumberAndAlphabeticalError;
    if (!alphanumeric.hasMatch(text)) {
      loginPwNumberAndAlphabeticalError = true;
    } else {
      loginPwNumberAndAlphabeticalError = false;
    }
    return loginPwNumberAndAlphabeticalError;
  }

  bool isNotContain8To20Characters({
    required String text,
  }) {
    return text.length > 20 || text.length < 8;
  }

  String? getCustomCryptoSymbol(symbol) {
    List customToken = ['ETH', 'BTCM', 'ETHM', 'USDM', 'MST'];
    if (symbol == "EUN" || symbol == "WEUN") {
      return 'images/Eurus_Blue.png';
    } else if (customToken.contains(symbol)) {
      return 'images/crypto_${symbol.toLowerCase()}.png';
    }
  }

  Future<CoinPriceList?> getCoingekoCryptoImage(
      List<CryptoCurrencyModel> list) async {
    if (common.coingeckoCoinsList.isEmpty) {
      common.coingeckoCoinsList = await api.get(
        url: 'https://api.coingecko.com/api/v3/coins/list',
        shouldPrintWrapped: false,
      );
    }

    final ids = common.coingeckoCoinsList
        .where((element) => list
            .map((e) => e.symbol)
            .contains((element['symbol'] as String).toUpperCase()))
        .map((e) => e['id'])
        .toList()
        .join(',');
    return await api
        .get(
          url:
              "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=$ids",
          shouldPrintWrapped: false,
        )
        .then((value) =>
            value is List<dynamic> ? CoinPriceList.fromJson(value) : null);
  }

  Future<Widget> getCryptoIcon(symbol, double size,
      {imgUrl, placeholder}) async {
    String? path = getCustomCryptoSymbol(symbol);
    if (path != null) {
      return SizedBox(
          width: size,
          height: size,
          child: Image.asset(path, package: 'euruswallet'));
    }
    return AppCryptoIcons.getIcon(symbol, size,
        source: await AppCryptoIcons.ckIconSource(symbol, imgUrl: imgUrl),
        imgUrl: imgUrl,
        placeholder: placeholder);
  }

  Widget getIcon(symbol, double size,
      {imgUrl, placeholder, required IconSourceType source}) {
    String? path = getCustomCryptoSymbol(symbol);
    if (path != null) {
      return SizedBox(
          width: size,
          height: size,
          child: Image.asset(path, package: 'euruswallet'));
    }
    return AppCryptoIcons.getIcon(symbol, size,
        source: source, imgUrl: imgUrl, placeholder: placeholder);
  }

  Future<void> navigateToAssetAllocationTransfer(
      BuildContext _navigatorContext, Map<String, dynamic> _tArgs) async {
    selectTokenSymbol = isEmptyString(string: _tArgs['symbol'])
        ? selectTokenSymbol
        : _tArgs['symbol'];
    common.selectTokenSymbol = selectTokenSymbol;
    BlockChainType type = BlockChainType
        .values[_tArgs['address'] == _tArgs['address$TEST_NET'] ? 0 : 1];
    String ethereumErc20ContractAddress = _tArgs['address$TEST_NET'];
    String eurusErc20ContractAddress = _tArgs['addressEurus'];
    await pushPage(
        page: TransferPage(
          titleName: "ASSET_ALLOCATION_PAGE.TITLE".tr(),
          fromBlockChainType: type,
          transferToMySelf: true,
          ethereumErc20ContractAddress: ethereumErc20ContractAddress,
          eurusErc20ContractAddress: eurusErc20ContractAddress,
        ),
        context: _navigatorContext);
    FocusScope.of(_navigatorContext).requestFocus(FocusNode());
  }

  List<Widget> homeFunctionsBarItemBuilder(context, [cArgs = const {}]) {
    final _scanner = (context, handler) async {
      /// await QRCode content as [String]
      String result = await scanQRCode(context: context);
      print("result:$result");
      return handler != null ? await handler(result) : result;
    };

    return [
// cArgs keys for send: 'ethereumErc20ContractAddress': String, 'eurusErc20ContractAddress': String, 'fromBlockChainType': int, 'disableSelectBlockchain': bool, canGetPrivateKeyHandler: Future<String> Function()
      GestureDetector(
          behavior: HitTestBehavior.translucent,
          onTap: () async {
            common.backupCurrentBlockchainSelection =
                common.currentBlockchainSelection;
            print("walletAccountEncryptedAddress" +
                (cArgs['walletAccountEncryptedAddress'] ?? ''));
            isCenWithdraw = false;
            await pushPage(
                page: SelectTargetPage(
                  titleName: "SEND_PAGE.MAIN_TITLE".tr(),
                  ethereumErc20ContractAddress:
                      cArgs['ethereumErc20ContractAddress'] ?? '0x0',
                  eurusErc20ContractAddress:
                      cArgs['eurusErc20ContractAddress'] ?? '0x0',
                  fromBlockChainType:
                      BlockChainType.values[cArgs['fromBlockChainType'] ?? 1],
                  disableSelectBlockchain:
                      cArgs['disableSelectBlockchain'] ?? false,
                  canGetPrivateKeyHandler: cArgs[
                      'canGetPrivateKeyHandler'] ?? /*() => Future.value("")*/ null,
                  userWalletAccountAvailableAssetsList:
                      cArgs['canGetWalletAccountAssetsList'] != null
                          ? await cArgs['canGetWalletAccountAssetsList']()
                          : null,
                ),
                context: context);
            common.currentBlockchainSelection =
                common.backupCurrentBlockchainSelection;
          },
          child: Column(children: [
            Padding(
                padding:
                    EdgeInsets.only(left: 17, right: 17, top: 17, bottom: 7),
                child: Image.asset('images/icon_send.png',
                    package: 'euruswallet',
                    width: 24,
                    height: 24,
                    fit: BoxFit.contain,
                    color: cArgs['btnColor'] ?? FXColor.blackColor)),
            Text('MAIN_FNC.SEND'.tr(),
                style: Theme.of(context).textTheme.caption?.merge(FXUI
                    .normalTextStyle
                    .copyWith(color: cArgs['btnColor'] ?? FXColor.blackColor))),
          ], mainAxisSize: MainAxisSize.max)),
// cArgs keys for receive: 'ethereumAddress': String, 'ethereumErrorPopUp': bool, 'eurusErrorPopUp': bool, 'errorText': String, 'errorPromptBuilder': Widget, 'replacingQRCodeWidget': Widget
      GestureDetector(
          behavior: HitTestBehavior.translucent,
          onTap: () async {
            common.backupCurrentBlockchainSelection =
                common.currentBlockchainSelection;

            print("ethereumAddress" + (cArgs['ethereumAddress'] ?? ''));
            await pushPage(
                page: ReceivePage(
                  ethereumAddress: cArgs['ethereumAddress'],
                  ethereumErrorPopUp: cArgs['ethereumErrorPopUp'],
                  eurusErrorPopUp: cArgs['eurusErrorPopUp'],
                  errorPromptBuilder: cArgs['errorPromptBuilder'],
                  errorText: cArgs['errorText'],
                  replacingQRCodeWidget: cArgs['replacingQRCodeWidget'],
                  blockChainType: cArgs['fromBlockChainType'] != null
                      ? BlockChainType.values[cArgs['fromBlockChainType']!]
                      : null,
                  disableSelectBlockchain:
                      cArgs['disableSelectBlockchain'] ?? false,
                ),
                context: context);
            common.currentBlockchainSelection =
                common.backupCurrentBlockchainSelection;

// emitting a data event in the Stream, only if there was a `replacingQRCodeWidget` being shown in the ReceivePage,
// which triggers subscribed listeners to react to data event,
// where one of the subscribers is the TopAppBar stateful widget instance
// TopAppBar instance is subscribing to this Stream at its State<TopAppBar>.initState()
            if (cArgs['replacingQRCodeWidget'] != null) refreshTopBar.add(true);
          },
          child: Column(children: [
            Padding(
                padding:
                    EdgeInsets.only(left: 17, right: 17, top: 17, bottom: 7),
                child: Image.asset('images/icon_receive.png',
                    package: 'euruswallet',
                    width: 24,
                    height: 24,
                    fit: BoxFit.contain,
                    color: cArgs['btnColor'] ?? FXColor.blackColor)),
            Text('MAIN_FNC.RECEIVE'.tr(),
                style: Theme.of(context).textTheme.caption?.merge(FXUI
                    .normalTextStyle
                    .copyWith(color: cArgs['btnColor'] ?? FXColor.blackColor))),
          ])),
// cArgs keys for scanner: 'scannerHandler': FutureOr<dynamic> Function(String)
      GestureDetector(
          behavior: HitTestBehavior.translucent,
          onTap: () async {
            var qrcode = await _scanner(context, cArgs['scannerHandler']);
            print("qrcode:$qrcode");
            if (qrcode == null || qrcode is String && qrcode.isEmpty) return;
            if (qrcode != null &&
                qrcode is String &&
                qrcode.isNotEmpty &&
                qrcode != 'empty') {
// contemporarily guard away qrcode scanner handler from cross chain transfer by qrcode scan result of uri scheme
// TODO: prompt modal dialog for "cross-chain transfer support coming soon"
              if ((cArgs['disableSelectBlockchain'] ?? false) &&
                  cArgs['fromBlockChainType'] != null &&
                  BlockChainType.values[cArgs['fromBlockChainType'] ?? 0] ==
                      BlockChainType.Eurus &&
                  !qrcode.startsWith('eurus:')) {
                ScaffoldMessenger.of(context).showSnackBar(SnackBar(
                    content: Text(
                        "Support Coming Soon for Transfer Sending to Ethereum from Eurus")));
                return;
              }
              if ((cArgs['disableSelectBlockchain'] ?? false) &&
                  cArgs['fromBlockChainType'] != null &&
                  BlockChainType.values[cArgs['fromBlockChainType'] ?? 1] ==
                      BlockChainType.Ethereum &&
                  qrcode.startsWith('eurus:')) {
                ScaffoldMessenger.of(context).showSnackBar(SnackBar(
                    content: Text(
                        "Support Coming Soon for Transfer Sending to Ethereum from Eurus")));
              }
              //isCenWithdraw = !qrcode.startsWith('eurus:');
              common.transferToMySelf = false;
              await pushPage(
                  page: SelectTargetPage(
                    titleName: "SEND_PAGE.MAIN_TITLE".tr(),
                    ethereumErc20ContractAddress:
                        cArgs['ethereumErc20ContractAddress'] ?? '0x0',
                    eurusErc20ContractAddress:
                        cArgs['eurusErc20ContractAddress'] ?? '0x0',
                    fromBlockChainType: BlockChainType.values[
                        cArgs['fromBlockChainType'] ??
                            (qrcode.startsWith('eurus:') || isCentralized()
                                ? 1
                                : 0)],
                    disableSelectBlockchain:
                        cArgs['disableSelectBlockchain'] ?? false,
                    canGetPrivateKeyHandler: cArgs[
                            'canGetPrivateKeyHandler'] ?? /*() => Future.value("")*/
                        null,
                    reciverText: qrcode,
                    userWalletAccountAvailableAssetsList:
                        cArgs['canGetWalletAccountAssetsList'] != null
                            ? await cArgs['canGetWalletAccountAssetsList']()
                            : null,
                  ),
                  context: context);
            }
          },
          child: Column(children: [
            Padding(
                padding:
                    EdgeInsets.only(left: 17, right: 17, top: 17, bottom: 7),
                child: Image.asset('images/icon_scan.png',
                    package: 'euruswallet',
                    width: 24,
                    height: 24,
                    fit: BoxFit.contain,
                    color: cArgs['btnColor'] ?? FXColor.blackColor)),
            Text('MAIN_FNC.SCAN'.tr(),
                style: Theme.of(context).textTheme.caption?.merge(FXUI
                    .normalTextStyle
                    .copyWith(color: cArgs['btnColor'] ?? FXColor.blackColor)))
          ])),
// cArgs keys for asset allocation: 'navigateToAssetAllocationPage' FutureOr<dynamic> Function()
// for contextual navigated-scene-aware asset allocation button in function bar item to be shown or not
      if ((!(cArgs['disableSelectBlockchain'] ?? false) ||
              (cArgs['${BlockChainType.values[cArgs['fromBlockChainType'] ?? 0] == BlockChainType.Ethereum ? 'eurus' : 'ethereum'}Erc20ContractAddress'] ??
                      '') !=
                  '') &&
          currentUserType == CurrentUserType.decentralized &&
          (cArgs['enableAllocation'] ?? true))
        GestureDetector(
            behavior: HitTestBehavior.translucent,
            onTap: () async {
              selectTokenSymbol = cArgs['symbol'];
              common.transferToMySelf = true;
              isCenWithdraw = false;
              common.pushPage(
                  page: AssetAllocationTokenListPage(), context: context);
            },
            child: Column(children: [
              Padding(
                  padding:
                      EdgeInsets.only(left: 17, right: 17, top: 17, bottom: 7),
                  child: Image.asset('images/icon_asset_allocation.png',
                      package: 'euruswallet',
                      width: 24,
                      height: 24,
                      fit: BoxFit.contain,
                      color: cArgs['btnColor'] ?? FXColor.blackColor)),
              Text.rich(
// TextSpan(
//     text: 'Asset', children: [TextSpan(text: "\nallocation")]),
                TextSpan(text: 'ASSET_ALLOCATION_PAGE.TITLE'.tr()),
                style: Theme.of(context).textTheme.caption?.merge(FXUI
                    .normalTextStyle
                    .copyWith(color: cArgs['btnColor'] ?? FXColor.blackColor)),
                textAlign: TextAlign.center,
              )
            ])),
      if (currentUserType == CurrentUserType.centralized &&
          cArgs['symbol'] != 'EUN')
        GestureDetector(
          behavior: HitTestBehavior.translucent,
          onTap: () async {
            print("walletAccountEncryptedAddress" +
                (cArgs['walletAccountEncryptedAddress'] ?? ''));
            isCenWithdraw = true;
            common.transferToMySelf = false;
            await pushPage(
                page: SelectTargetPage(
                  titleName: "ASSET_ALLOCATION_PAGE.WITHDRAWAL".tr(),
                  ethereumErc20ContractAddress:
                      cArgs['ethereumErc20ContractAddress'] ?? '0x0',
                  eurusErc20ContractAddress:
                      cArgs['eurusErc20ContractAddress'] ?? '0x0',
                  fromBlockChainType:
                      BlockChainType.values[cArgs['fromBlockChainType'] ?? 1],
                  disableSelectBlockchain:
                      cArgs['disableSelectBlockchain'] ?? false,
                  canGetPrivateKeyHandler: cArgs[
                      'canGetPrivateKeyHandler'] ?? /*() => Future.value("")*/ null,
                  userWalletAccountAvailableAssetsList:
                      cArgs['canGetWalletAccountAssetsList'] != null
                          ? await cArgs['canGetWalletAccountAssetsList']()
                          : null,
                ),
                context: context);
          },
          child: Column(
            children: [
              Padding(
                  padding:
                      EdgeInsets.only(left: 17, right: 17, top: 17, bottom: 7),
                  child: Image.asset('images/icon_withdraw.png',
                      package: 'euruswallet',
                      width: 24,
                      height: 24,
                      fit: BoxFit.contain,
                      color: cArgs['btnColor'] ?? FXColor.blackColor)),
              Text.rich(
// TextSpan(
//     text: 'Asset', children: [TextSpan(text: "\nallocation")]),
                TextSpan(text: 'ASSET_ALLOCATION_PAGE.WITHDRAWAL'.tr()),
                style: Theme.of(context).textTheme.caption?.merge(FXUI
                    .normalTextStyle
                    .copyWith(color: cArgs['btnColor'] ?? FXColor.blackColor)),
                textAlign: TextAlign.center,
              )
            ],
          ),
        ),
    ];
  }

  Future<LoginBySignModel?> successMoveToHomePage({
    required CurrentUserType userType,
    required BuildContext context,
    required String loginPassword,
    required String email,
    String? address,
    String? privateKey,
    LoginBySignModel? loginBySignModel,
    UserProfile? userProfile,
    bool? isRegister = false,
  }) async {
    LoginBySignModel? result;
    if (userType == CurrentUserType.centralized) {
      result = loginBySignModel == null
          ? await api.loginBySignature(
              email: email.toLowerCase(), password: loginPassword)
          : loginBySignModel;
      if (await common.checkApiError(
          context: context,
          errorString: result.message,
          returnCode: result.returnCode)) {
        common.loginPassword = loginPassword;
        if (result.decryptedMnemonic != null &&
            result.isMetaMaskUser != null &&
            !(result.isMetaMaskUser)! &&
            common.rsaPrivateKey != null) {
          await moveToHomePage(
            userType: userType,
            context: context,
            privateKey: common.loginAddressPair?.privateKey ?? '',
            address: common.cenUserWalletAddress ?? '0x0',
            pw: common.loginPassword ?? '',
            isRegister: isRegister,
          );
        } else {
          common.registerDevice = await api.registerDevice();
          if (await common.checkApiError(
            context: context,
            errorString: common.registerDevice?.message,
            returnCode: common.registerDevice?.returnCode,
          )) {
            common.verifyCodePageType =
                VerifyCodePageType.newDeviceResetPublicKey;
            common.pushPage(page: VerifyCodePage(), context: context);
          }
        }
      }
    } else {
      if (userProfile != null) common.currentUserProfile = userProfile;
      await moveToHomePage(
        userType: userType,
        context: context,
        privateKey: privateKey,
        address: address,
        pw: loginPassword,
        isRegister: isRegister,
      );
    }

    return result;
  }

  Future<dynamic> moveToHomePage({
    required CurrentUserType userType,
    String? address,
    String? privateKey,
    required String pw,
    required BuildContext context,
    bool? isRegister = false,
  }) async {
    final _encrypt = (String? c) => CommonMethod.passwordEncrypt(pw, c ?? '');
    final encryptedAddress = _encrypt(address) ?? '';

    await setActiveAccount(encryptedAddress,
        displayUserName:
            userType == CurrentUserType.centralized ? email : address,
        enPrivateKey:
            isEmptyString(string: privateKey) ? null : _encrypt(privateKey));

    final decryptedString = CommonMethod.passwordDecrypt(pw, encryptedAddress);

    this.currentAddress = address;
    if (decryptedString != null) this.encryptedAddress = encryptedAddress;
    currentUserType = userType;
    await NormalStorageKit().setValue(
        userType == CurrentUserType.centralized
            ? 'CurrentUserType.centralized'
            : 'CurrentUserType.decentralized',
        'currentUserType');
    await SecureStorageKit()
        .setValue(address!.toLowerCase(), 'currentUserAddress');

    common.startTimer();

    final userStorageResponse = await api.getUserStorage();
    userStorage = userStorageResponse.storage ?? UserStorage();
    final currentLanguage = userStorage.currentLanguage;
    if (currentLanguage != null) {
      Locale langToSet = currentLanguage == 'cn'
          ? Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hans')
          : currentLanguage == 'tw'
              ? Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hant')
              : Locale('en', 'US');
      await NormalStorageKit().setValue(langToSet.toString(), 'APP_LANGUAGE');
      context.setLocale(langToSet);
    }

    currentBlockchainSelection = BlockChainType.Eurus;
    topSelectedBlockchainType = BlockChainType.Eurus;

    if (isRegister ?? false) {
      TrackManager.instance.trackSignUp(userType: userType);
    } else {
      TrackManager.instance.trackLogin(userType: userType);
    }
    if (envType == EnvType.Production)
      FirebaseCrashlytics.instance.setUserIdentifier(
          (isCentralized() ? this.email : this.currentAddress) ?? '');

    EasyLoading.dismiss();

    if (ModalRoute.of(context)?.canPop ?? false) {
      final completer = Completer();
      await Navigator.pushReplacementNamed(
        context,
        'HomePage',
        result: completer.future,
      );
      completer.complete(true);
      return true;
    } else {
      return await Navigator.pushNamed(
        context,
        'HomePage',
      );
    }
  }

  Future<String> updateAssetsListIfEmpty() async {
    final assetsListKey = 'assetsList_${encryptedAddress ?? ''}';
    String? alistString = await NormalStorageKit().readValue(assetsListKey);
    List<CryptoCurrencyModel> storedList = isEmptyString(string: alistString)
        ? []
        : (jsonDecode(alistString ?? '') as List)
            .map((e) => CryptoCurrencyModel.fromJson(e as Map<String, dynamic>))
            .toList();

    Map<String, CryptoCurrencyModel> supportedTokens =
        await common.getSupportedTokens();

    final tokenSymbols = [
      'EUN',
      'USDM',
      'BTCM',
      'ETHM',
      'MST',
      'ETH',
      'USDT',
      if (envType == EnvType.Dev) 'BNB',
      if (envType == EnvType.Dev) 'BTC'
    ];

    if (storedList.isEmpty)
      storedList =
          tokenSymbols.map((e) => supportedTokens[e]).whereNotNull().toList();

    for (var entry in supportedTokens.entries) {
      final key = entry.key;
      final value = entry.value;
      if (!tokenSymbols.contains(key)) continue;

      final index = storedList
          .indexWhere((element) => element.symbol.toUpperCase() == key);
      final newToken = value
        ..iconSource = await AppCryptoIcons.ckIconSource(key);
      if (index != -1) {
        storedList[index] = newToken..showAssets = storedList[index].showAssets;
      } else {
        storedList.add(newToken);
      }
    }

    String rdString = jsonEncode(storedList);

    await NormalStorageKit().setValue(rdString, assetsListKey);

    return rdString;
  }

  bool? get shouldShowReward =>
      rewardSchemes != null ? rewardSchemes!.isNotEmpty : null;
  Future<void> routeToRewardDetailPage({required BuildContext context}) async {
    if (shouldShowReward ?? false) {
      context.loaderOverlay.show();
      final result = await api.getUserMarketingRewardList();
      context.loaderOverlay.hide();
      rewardedList = result.list;

      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (context) => RewardDetailPage(
            scheme: (rewardSchemes ?? []).first,
          ),
        ),
      );
    }
  }

  void openLiveChat({
    String? visitorName,
    String? visitorEmail,
  }) {
    LivechatInc.start_chat(
      "12610959",
      isCentralized() ? "15" : "16",
      visitorName ?? "guest",
      visitorEmail ?? "guest@gmail.com",
    );
  }

  Future<T?> showGasLeakDialog<T>({
    required BuildContext context,
  }) {
    return showDialog<T>(
      barrierDismissible: false,
      context: context,
      builder: (_) => AcknowledgementDialog(
        statement: 'REFUEL.DIALOG.ERROR_CONTENT'.tr(),
        mainIcon: Image.asset(
          'images/icn_fuel.png',
          package: 'euruswallet',
          width: MediaQuery.of(context).size.width / 4,
        ),
        buttonText: 'REFUEL.TITLE'.tr(),
        buttonHandler: () async {
          await Navigator.of(context).push(MaterialPageRoute(
            builder: (_) {
              return TopUpPaymentWalletPage();
            },
          ));
        },
      ),
    );
  }
}

bool isAssetEditable(CryptoCurrencyModel cryptoCurrencyModel) {
  final nonEditableSymbolList = ['EUN', 'ETH', 'USDT'];
  if (nonEditableSymbolList.contains(cryptoCurrencyModel.symbol)) return false;
  return true;
}

bool isCentralized() {
  return common.currentUserType == CurrentUserType.centralized;
}

AddressPair advanceGenAddressPairPaymentPw(AdvanceGenAddressPairArgs args) {
  print('advanceGenAddressPairPaymentPw args:$args');
  return MnemonicKit()
      .mnemonicPhraseToAddressPair(args.mnemonic, pw: args.pw, account: 1);
}

AddressPair advanceGenAddressPair(AdvanceGenAddressPairArgs args) {
  print('advanceGenAddressPair args:$args');
  return MnemonicKit().mnemonicPhraseToAddressPair(args.mnemonic, pw: args.pw);
}

/// Set which account is currently active
/// 0 == active , 1 == inActive
/// 1. Update [{prefix}accounts] in secure storage
///   a. Add / Update active account encrypted address key value to 0
///   c. Update previous active account (if exists) to 1
/// TODO: 2. Disable all previous account biometric (if turned on)
/// TODO: 3. Turn on biometric (if checked)
Future<void> setActiveAccount(
  String encryptedAddress, {
  String? displayUserName,
  String? enPrivateKey,
  String? enMPhrase,
  bool? bioOn,
}) async {
  final String prefix = await common.prefix;

  /// 1. Update [{prefix}accounts] in secure storage
  final String? accountsUri =
      await SecureStorageKit().readValue('${prefix}accounts');

  /// encryptedAddress to 0 , other to 1 , addNewIfNotFound == false or null -> can't found not add new
  String udpatedAcUri = common.updateUriVal(
    accountsUri,
    encryptedAddress,
    "0",
    addNewIfNotFound: true,
    orgVal: "1",
  );
  if (accountsUri == null) {
    await SecureStorageKit().setValue(udpatedAcUri, '${prefix}accounts');
  }

  await NormalStorageKit().setValue(displayUserName ?? '', 'displayUserName');
}

/// Set account info to local
///
/// Accounts are localed in secure storage under the key [${prefix}_local_accounts_w_info]
Future<void> setAcToLocal(
  UserProfile up, {
  bool delete = false,
}) async {
  final String valKey = up.encryptedAddress;

  final String _prefix = await common.prefix;
  final String? _localAcs =
      await SecureStorageKit().readValue('${_prefix}_local_accounts_w_info');

  final String? _orgVal = common.getUriVal(_localAcs, valKey);

  if (_orgVal != null) {
    up.alias = UserProfile.fromJson(jsonDecode(_orgVal)).alias;
  }

  final String updatedLocalAcs = common.updateUriVal(
    _localAcs,
    valKey,
    delete == true ? null : jsonEncode(up.toJson()),
    addNewIfNotFound: true,
  );

  await SecureStorageKit()
      .setValue(updatedLocalAcs, '${_prefix}_local_accounts_w_info');
}

/// Get accounts from local
Future<List<UserProfile>> getLocalAcs() async {
  final String _prefix = await common.prefix;

  final String? _localAcs =
      await SecureStorageKit().readValue('${_prefix}_local_accounts_w_info');
  final Map<String, String> _localAcsUri =
      Uri(query: _localAcs).queryParameters;

  final String? _accounts =
      await SecureStorageKit().readValue('${_prefix}accounts');
  final Map<String, String> _accountsUri =
      Uri(query: _accounts).queryParameters;

  List<UserProfile> acList = [];

  _localAcsUri.forEach((key, value) {
    UserProfile dummyUp = UserProfile.fromJson(jsonDecode(value));
    if (dummyUp.userType == CurrentUserType.centralized) {
      dummyUp.address = dummyUp.address.toLowerCase();
    }
    acList = _accountsUri[dummyUp.encryptedAddress] == "0"
        ? [dummyUp, ...acList]
        : [...acList, dummyUp];
  });
  print("acList1:$acList");
  final ids = acList.map((UserProfile e) => e.address).toSet();
  acList.retainWhere((x) => ids.remove(x.address));
  print("acList:$acList");

  return acList;
}

String encryptStringFromPrivateKey({
  required String message,
  required String privateKey,
}) {
  // /web3dart-2.0.0-dev.9/lib/crypto.dart
  // /web3dart-2.0.0-dev.9/lib/src/crypto/secp256k1.dart
  // /pointycastle-1.0.2/lib/signers/ecdsa_signer.dart
  // /web3dart-2.0.0-dev.9/lib/src/utils/typed_data.dart
  // /web3dart-2.0.0-dev.9/lib/src/core/expensive_operations.dart
  // /web3dart-2.0.0-dev.9/lib/src/core/client.dart
  // /web3dart-2.0.0-dev.9/lib/src/core/transaction_signer.dart
  // /web3dart-2.0.0-dev.9/lib/src/credentials/credentials.dart
  // /web3dart-2.0.0-dev.9/lib/src/crypto/formatting.dart
  final Uint8List _messageHashBytes = _web3dartCrypto.keccakUtf8(message);
  final Uint8List _privateKeyBytes = _web3dartCrypto.hexToBytes(privateKey);

  final _web3dartCrypto.MsgSignature _msgSignature =
      _web3dartCrypto.sign(_messageHashBytes, _privateKeyBytes);
  print("_msgSignature.r:${_msgSignature.r}");
  print("_msgSignature.s:${_msgSignature.s}");
  final Uint8List _rBytes =
      padUint8ListTo32(_web3dartCrypto.unsignedIntToBytes(_msgSignature.r));
  final Uint8List _sBytes =
      padUint8ListTo32(_web3dartCrypto.unsignedIntToBytes(_msgSignature.s));
  final Uint8List _vBytes =
      _web3dartCrypto.unsignedIntToBytes(BigInt.from(_msgSignature.v));
  final Uint8List _concatenatedInto65Bytes =
      uint8ListFromList(_rBytes + _sBytes + _vBytes);
  print(
      "_canGetSignature ${_web3dartCrypto.bytesToHex(_concatenatedInto65Bytes)}");
  return _web3dartCrypto.bytesToHex(_concatenatedInto65Bytes);
}

/// you can use Transaction
CommonMethod common = CommonMethod();

int getCurrentTimeStamp() {
  int timeStamp = DateTime.now().millisecondsSinceEpoch;
  print("timeStamp:$timeStamp");
  return timeStamp;
}

bool isEmptyString({String? string}) {
  return (string?.isEmpty ?? true) ||
      string == "null" ||
      string == '0x0' ||
      string == 'empty';
}

void openKYCCamera(BuildContext context, Function(String) onCapture,
    {CameraSide cameraSide = CameraSide.back,
    DocType docType = DocType.Unknown}) {
  common.pushPage(
      page: KYCCamera(
        onCapture: (String path) => {onCapture(path)},
        imageMask: cameraSide == CameraSide.back
            ? CameraFocus.rectangle(
                color: Colors.black.withOpacity(0.5),
              )
            : CameraFocus.circle(
                color: Colors.black.withOpacity(0.5),
              ),
        cameraSide: cameraSide,
        docType: docType,
      ),
      context: context);
}

/// Generate 12 / 24 words mnemonic phrase base on strength
///
/// [strength] = 128 for 12 words
/// [strength] = 256 for 24 words
String genMnemonic([int strength = 128]) {
  return MnemonicKit().genMnemonicPhrase(strength: strength);
}

/// Generate Base58 from mnemonic phrase
String? genBase58(String? mPhrase) {
  if (mPhrase == null) return null;
  return MnemonicKit().mnemonicToBase58(mPhrase);
}

/// Generate address and private key from Base58
AddressPair genAddressPair(String b58) {
  return MnemonicKit().genAddressPairFromBase58(b58);
}

void printWrapped(String text, {shouldWrapped = true}) {
  if (shouldWrapped) {
    final pattern = RegExp('.{1,800}'); // 800 is the size of each chunk
    pattern.allMatches(text).forEach((match) => print(match.group(0)));
  } else {
    print(text);
  }
}

String getBlockChainName(BlockChainType blockChainType) {
  Map<BlockChainType, String> list = {
    BlockChainType.Eurus: "Eurus",
    BlockChainType.Ethereum: "Ethereum",
    BlockChainType.BinanceCoin: "Binance Coin",
  };
  return list[blockChainType] ?? "";
}

BlockChainType getBlockChainTypeBySymbol(String symbol) {
  Map<String, BlockChainType> list = {
    'EUN': BlockChainType.Eurus,
    'ETH': BlockChainType.Ethereum,
    'BNB': BlockChainType.BinanceCoin,
  };
  return list[symbol] ?? BlockChainType.Ethereum;
}

String getSymbolByBlockChainType(BlockChainType blockChainType) {
  Map<BlockChainType, String> list = {
    BlockChainType.Eurus: "EUN",
    BlockChainType.Ethereum: "ETH",
    BlockChainType.BinanceCoin: "BNB",
  };
  return list[blockChainType] ?? "";
}

String getAddressSuffix(BlockChainType blockChainType) {
  switch (blockChainType) {
    case BlockChainType.Eurus:
      return "Eurus";
    case BlockChainType.Ethereum:
      if (!isEmptyString(string: TEST_NET))
        return TEST_NET;
      else
        return 'Ethereum';
    case BlockChainType.BinanceCoin:
      return "BSC";
    default:
      return "";
  }
}

String getUserAddress(UserProfile e) {
  String address = e.address;
  if (e.userType == CurrentUserType.decentralized &&
      e.seedPraseBackuped != true) {
    address = 'HOME.HEADER.SEED_PHRASE_NOT_BACKUP'.tr();
  }
  return address;
}

Future<String> scanQRCode({required BuildContext context}) async {
  String result = await QRCodeScanner().tryOpenScanner(context);
  await common.pushPage(
    page: AfterScanPage(qrCodeString: result),
    context: context,
  );
  return result;
}
