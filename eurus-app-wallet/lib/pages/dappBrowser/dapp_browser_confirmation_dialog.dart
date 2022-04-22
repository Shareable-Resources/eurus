import 'package:collection/collection.dart';
import 'package:decimal/decimal.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:eth_sig_util/constant/typed_data_version.dart';
import 'package:eth_sig_util/model/typed_data.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/cta_button.dart';
import 'package:euruswallet/extension/ethereum_address_extension.dart';
import 'package:euruswallet/model/crypto_currency_model.dart';
import 'package:euruswallet/pages/dappBrowser/web3_bridge.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:flutter_neumorphic/flutter_neumorphic.dart';
import 'package:web3dart/crypto.dart';
import 'package:collection/collection.dart';

enum DappBrowserConfirmationDialogType {
  transfer,
  ethSign,
  personalSign,
  signTypeDataV1,
  signTypeDataV3,
  signTypeDataV4,
  approve,
  switchNetwork,
  unknown
}

extension DappBrowserConfirmationDialogTypeExtension
    on DappBrowserConfirmationDialogType {
  bool get isTransactionType {
    return this == DappBrowserConfirmationDialogType.transfer ||
        this == DappBrowserConfirmationDialogType.approve ||
        this == DappBrowserConfirmationDialogType.unknown;
  }

  bool get isSignatureType {
    return this == DappBrowserConfirmationDialogType.ethSign ||
        this == DappBrowserConfirmationDialogType.personalSign ||
        this == DappBrowserConfirmationDialogType.signTypeDataV1 ||
        this == DappBrowserConfirmationDialogType.signTypeDataV3 ||
        this == DappBrowserConfirmationDialogType.signTypeDataV4;
  }

  String get title {
    switch (this) {
      case DappBrowserConfirmationDialogType.transfer:
        return 'DAPP_BROWSER.TRANSACTION_TYPE.TRANSFER'.tr();
      case DappBrowserConfirmationDialogType.ethSign:
      case DappBrowserConfirmationDialogType.personalSign:
      case DappBrowserConfirmationDialogType.signTypeDataV1:
      case DappBrowserConfirmationDialogType.signTypeDataV3:
      case DappBrowserConfirmationDialogType.signTypeDataV4:
        return 'DAPP_BROWSER.TRANSACTION_TYPE.REQUEST_SIGNATURE'.tr();
      case DappBrowserConfirmationDialogType.approve:
        return 'DAPP_BROWSER.TRANSACTION_TYPE.APPROVE'.tr();
      case DappBrowserConfirmationDialogType.switchNetwork:
        return 'DAPP_BROWSER.TRANSACTION_TYPE.SWITCH_NETWORK.TITLE'.tr();
      case DappBrowserConfirmationDialogType.unknown:
        return 'DAPP_BROWSER.TRANSACTION_TYPE.CONTRACT_INTERACTION'.tr();
    }
  }

  String get toTitle {
    switch (this) {
      case DappBrowserConfirmationDialogType.transfer:
      case DappBrowserConfirmationDialogType.unknown:
        return 'TX_PAGE.TRANSFER_TO'.tr();
      case DappBrowserConfirmationDialogType.ethSign:
      case DappBrowserConfirmationDialogType.personalSign:
      case DappBrowserConfirmationDialogType.signTypeDataV1:
      case DappBrowserConfirmationDialogType.signTypeDataV3:
      case DappBrowserConfirmationDialogType.signTypeDataV4:
      case DappBrowserConfirmationDialogType.switchNetwork:
        return '';
      case DappBrowserConfirmationDialogType.approve:
        return 'DAPP_BROWSER.GRANTED_TO'.tr();
    }
  }
}

class DappBrowserConfirmationDialog extends StatefulWidget {
  const DappBrowserConfirmationDialog(
    this.type, {
    Key? key,
    this.title,
    this.from,
    this.to,
    this.value,
    this.token,
    this.gas,
    this.data,
    this.message,
    this.chainId,
  }) : super(key: key);

  final DappBrowserConfirmationDialogType type;
  final String? title;
  final EthereumAddress? from;
  final EthereumAddress? to;
  final EtherAmount? value;
  final CryptoCurrencyModel? token;
  final int? gas;
  final Uint8List? data;
  final String? message;
  final String? chainId;

  @override
  _DappBrowserConfirmationDialogState createState() =>
      _DappBrowserConfirmationDialogState();
}

class _DappBrowserConfirmationDialogState
    extends State<DappBrowserConfirmationDialog> {
  Uri? uri;
  Favicon? favicon;

  TextEditingController passwordEditingController = TextEditingController();

  bool isPasswordMasked = false;
  String? errorMessage;
  bool isCheckingData = false;

  Decimal value = Decimal.zero;
  Decimal formattedValue = Decimal.zero;
  Decimal? usdExchangeRate;
  String? balance;
  Decimal estimatedGasFee = Decimal.zero;

  List<EIP712TypedData>? eip712TypedDatas;
  TypedMessage? typedMessage;

  @override
  void initState() {
    if (widget.type.isTransactionType)
      estimatedGasFee = (Decimal.fromInt((widget.gas ?? 0)) *
          Decimal.fromBigInt(web3dart
              .getGasPrice(blockChainType: Web3Bridge.instance.blockChainType)
              .getInWei) /
          Decimal.fromInt(pow(10, 18).toInt()));

    if (widget.type == DappBrowserConfirmationDialogType.signTypeDataV1 ||
        widget.type == DappBrowserConfirmationDialogType.signTypeDataV3 ||
        widget.type == DappBrowserConfirmationDialogType.signTypeDataV4) {
      initSignTypeData();
    }

    final currentPrice = common.coinPriceList?.coinPriceList
        .firstWhereOrNull((element) =>
            element.symbol.toLowerCase() == widget.token?.symbol.toLowerCase())
        ?.currentPrice;
    if (currentPrice != null)
      usdExchangeRate = Decimal.parse(currentPrice.toString());

    value = Decimal.parse((widget.value?.getInWei ?? 0.0).toString());

    super.initState();

    if (Web3Bridge.instance.webViewController != null)
      try {
        Future.wait([
          Web3Bridge.instance.webViewController!.getUrl(),
          Web3Bridge.instance.webViewController!.getFavicons(),
        ]).then((value) {
          uri = value[0] as Uri?;
          favicon = (value[1] as List<Favicon>?)?.firstOrNull;
          if (widget.type.isTransactionType) {
            String? contractAddress =
                Web3Bridge.instance.blockChainType == BlockChainType.Eurus
                    ? widget.token?.addressEurus
                    : widget.token?.addressEthereum;
            Future.wait([
              getFormattedValue(),
              if (contractAddress != null)
                common.getERC20BalanceAndInit(
                  widget.from.toString(),
                  contractAddress,
                  Web3Bridge.instance.blockChainType,
                )
            ]).then((value) {
              setState(() {
                formattedValue = value[0] as Decimal;
                if (contractAddress != null) balance = value[1] as String?;
              });
            });
          } else {
            setState(() {});
          }
        });
      } catch (e) {}
  }

  initSignTypeData() {
    late final rawTypedData;
    try {
      rawTypedData = jsonDecode(widget.message ?? '');

      if (widget.type == DappBrowserConfirmationDialogType.signTypeDataV3 ||
          widget.type == DappBrowserConfirmationDialogType.signTypeDataV4) {
        typedMessage = TypedMessage.fromJson(rawTypedData);
        return;
      }

      try {
        if (rawTypedData is List) {
          eip712TypedDatas =
              rawTypedData.map((e) => EIP712TypedData.fromJson(e)).toList();
        } else {
          eip712TypedDatas = [EIP712TypedData.fromJson(rawTypedData)];
        }
      } catch (_) {
        throw ArgumentError(
            'jsonData format is not corresponding to EIP712TypedData');
      }
    } catch (_) {
      throw ArgumentError('jsonData format must be correct');
    }
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {
        FocusScopeNode currentFocus = FocusScope.of(context);

        if (currentFocus.hasFocus) {
          currentFocus.unfocus();
        }
      },
      child: Scaffold(
        backgroundColor: Colors.transparent,
        body: Center(
          child: Container(
            margin: EdgeInsets.symmetric(horizontal: 26),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Padding(
                  padding: EdgeInsets.only(
                    left: 16.0,
                    top: 20.0,
                    right: 16.0,
                    bottom: 8,
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      favicon != null
                          ? Image.network(
                              favicon!.url.toString(),
                              width: 16,
                              errorBuilder: (BuildContext context,
                                  Object exception, StackTrace? stackTrace) {
                                return Container();
                              },
                            )
                          : Container(),
                      uri != null
                          ? Flexible(
                              child: Padding(
                                padding:
                                    const EdgeInsets.symmetric(horizontal: 8.0),
                                child: Text(
                                  uri!.origin.toString(),
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                  textAlign: TextAlign.center,
                                  style: FXUI.subtitleTextStyle.copyWith(
                                    color: FXColor.middleBlack,
                                    fontWeight: FontWeight.normal,
                                  ),
                                ),
                              ),
                            )
                          : Container(),
                      uri != null && uri!.scheme == "https"
                          ? Image.asset(
                              'images/icn_lock.png',
                              package: 'euruswallet',
                              width: 12,
                            )
                          : Container(),
                    ],
                  ),
                ),
                TextButton.icon(
                  icon: Image.asset(
                    'images/icn_browser_network_status.png',
                    package: 'euruswallet',
                    width: 8,
                  ),
                  label: Text(
                    Web3Bridge.instance.currentNetworkType.displayValue,
                    style: FXUI.subtitleTextStyle.copyWith(
                      fontSize: 9,
                      fontWeight: FontWeight.normal,
                      color: common.getBackGroundColor(),
                    ),
                  ),
                  onPressed: null,
                  style: TextButton.styleFrom(
                    padding: EdgeInsets.zero,
                    minimumSize: Size.zero,
                    tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                  ),
                ),
                SizedBox(height: 12),
                Stack(
                  alignment: AlignmentDirectional.center,
                  children: [
                    Divider(height: 1),
                    Positioned(
                      child: TextButton(
                        style: TextButton.styleFrom(
                          shape: RoundedRectangleBorder(
                            side: BorderSide(color: FXColor.dividerColor),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          backgroundColor: Colors.white,
                          textStyle: FXUI.titleTextStyle.copyWith(fontSize: 12),
                          padding: EdgeInsets.symmetric(
                            vertical: 4,
                            horizontal: 28,
                          ),
                          minimumSize: Size.zero,
                          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                        ),
                        onPressed: null,
                        child: Text(
                          widget.title ?? widget.type.title,
                          style: TextStyle(
                            color: common.getBackGroundColor(),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
                if (widget.type == DappBrowserConfirmationDialogType.ethSign)
                  Padding(
                    padding: EdgeInsets.only(
                      top: 16.0,
                      left: 16.0,
                      right: 16.0,
                    ),
                    child: Text(
                      'DAPP_BROWSER.ETH_SIGN_WARNING'.tr(),
                      style: FXUI.titleTextStyle.copyWith(
                        fontSize: 12,
                        color: FXColor.alertRedColor,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
                Flexible(
                  child: SingleChildScrollView(
                    child: Padding(
                      padding: EdgeInsets.all(16.0),
                      child: widget.type.isTransactionType
                          ? getTransactionWidget()
                          : widget.type.isSignatureType
                              ? getSignatureMessageWidget()
                              : getSwitchNetworkWidget(),
                    ),
                  ),
                ),
                Divider(height: 1),
                Padding(
                  padding: const EdgeInsets.all(16.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      if (widget.type !=
                          DappBrowserConfirmationDialogType.switchNetwork)
                        Padding(
                          padding: const EdgeInsets.only(bottom: 25.0),
                          child: FXUI.neumorphicTextField(
                            context,
                            padding: const EdgeInsets.symmetric(
                              vertical: 12,
                              horizontal: 16,
                            ),
                            shape: NeumorphicBoxShape.roundRect(
                                BorderRadius.circular(14)),
                            keyboardType: TextInputType.visiblePassword,
                            hintText: isCentralized()
                                ? 'WALLET_LOCKER_DIALOG.CEN_TITLE'.tr()
                                : 'WALLET_LOCKER_DIALOG.TITLE'.tr(),
                            errorText: errorMessage,
                            prefixIcon: Padding(
                              padding: EdgeInsets.only(right: 12),
                              child: Image.asset(
                                'images/icn_lock.png',
                                package: 'euruswallet',
                                color: FXColor.placeholderGreyColor,
                                width: 15,
                              ),
                            ),
                            suffixIcon: GestureDetector(
                              onTap: () {
                                setState(() {
                                  isPasswordMasked = !isPasswordMasked;
                                });
                              },
                              child: Image.asset(
                                isPasswordMasked
                                    ? 'images/eyeClose.png'
                                    : 'images/eyeOpen.png',
                                package: 'euruswallet',
                                width: 16,
                                height: 16,
                                color: common.getBackGroundColor(),
                              ),
                            ),
                            controller: passwordEditingController,
                            obscureText: !isPasswordMasked,
                            autocorrect: false,
                          ),
                        ),
                      Row(
                        children: [
                          Expanded(
                            child: CtaButton(
                              type: CtaButtonType.primary,
                              onPressed: () async {
                                if (widget.type ==
                                    DappBrowserConfirmationDialogType
                                        .switchNetwork) {
                                  Navigator.pop(context, true);
                                } else {
                                  setState(() {
                                    errorMessage = null;
                                  });

                                  await Future.delayed(
                                      Duration(milliseconds: 250));

                                  final credentials = await validatePassword();
                                  if (credentials != null) {
                                    final isGasSufficient =
                                        await checkGasBalance();
                                    if (!isGasSufficient) {
                                      if (isCentralized()) {
                                        common.showGasLeakDialog(
                                            context: context);
                                      } else {
                                        setState(() {
                                          errorMessage =
                                              "TX_PAGE.ERROR.INSUFFICIENT_GAS_FEE"
                                                  .tr();
                                        });
                                      }
                                      return;
                                    }
                                    Navigator.pop(context, credentials);
                                  } else {
                                    setState(() {
                                      errorMessage =
                                          "COMMON_ERROR.AUTH_FAIL".tr();
                                    });
                                  }
                                }
                              },
                              text: widget.type ==
                                      DappBrowserConfirmationDialogType
                                          .switchNetwork
                                  ? widget.type.title
                                  : 'COMMON.CONFIRM'.tr(),
                            ),
                          ),
                          SizedBox(width: 16),
                          Expanded(
                            child: CtaButton(
                              type: CtaButtonType.secondary,
                              onPressed: () => Navigator.pop(context),
                              text: widget.type ==
                                      DappBrowserConfirmationDialogType
                                          .switchNetwork
                                  ? 'COMMON.CANCEL'.tr()
                                  : 'COMMON.REJECT'.tr(),
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget getSwitchNetworkWidget() {
    NetworkType? networkType = NetworkType.values.firstWhereOrNull(
        (element) => element.isEnabled && element.chainId == widget.chainId);
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 24.0, horizontal: 6.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(
            'DAPP_BROWSER.TRANSACTION_TYPE.SWITCH_NETWORK.SUBTITLE'.tr(),
            style: FXUI.inputStyle.copyWith(
              fontSize: 18.0,
            ),
          ),
          SizedBox(height: 11.0),
          Text(
            'DAPP_BROWSER.TRANSACTION_TYPE.SWITCH_NETWORK.DESCRIPTION'.tr(),
            style: FXUI.inputStyle.copyWith(
              fontSize: 14.0,
              color: FXColor.lightGrey,
            ),
          ),
          SizedBox(height: 40),
          if (networkType != null)
            Center(
              child: Container(
                padding: EdgeInsets.symmetric(vertical: 5, horizontal: 14),
                decoration: BoxDecoration(
                  border: Border.all(
                    color: FXColor.dividerColor,
                  ),
                  borderRadius: BorderRadius.circular(14),
                ),
                child: TextButton.icon(
                  icon: Image.asset(
                    'images/icn_browser_network_status.png',
                    package: 'euruswallet',
                    width: 8,
                  ),
                  label: Text(
                    networkType.displayValue,
                    style: FXUI.subtitleTextStyle.copyWith(
                      fontSize: 9,
                      fontWeight: FontWeight.normal,
                      color: common.getBackGroundColor(),
                    ),
                  ),
                  onPressed: null,
                  style: TextButton.styleFrom(
                    padding: EdgeInsets.zero,
                    minimumSize: Size.zero,
                    tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget getSignatureMessageWidget() {
    return Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: <Widget>[
              if (widget.from != null)
                Padding(
                  padding: const EdgeInsets.only(bottom: 12.0),
                  child: getAddressInfoWidget(
                    title: 'TX_PAGE.FROM'.tr(),
                    subtitle: widget.from?.eip55TruncatedString ?? '',
                  ),
                ),
            ] +
            (widget.type == DappBrowserConfirmationDialogType.signTypeDataV1 &&
                    eip712TypedDatas != null
                ? getSignTypedDataV1Widget()
                : (widget.type ==
                                DappBrowserConfirmationDialogType
                                    .signTypeDataV3 ||
                            widget.type ==
                                DappBrowserConfirmationDialogType
                                    .signTypeDataV4) &&
                        typedMessage != null
                    ? [
                        getAddressInfoWidget(
                          title: 'COMMON.MESSAGE'.tr(),
                          content: Column(
                            children: getSignTypedDataV3V4WidgetContent(
                                typedMessage!.message),
                          ),
                        )
                      ]
                    : [
                        if (!isEmptyString(string: widget.message))
                          getAddressInfoWidget(
                            title: 'COMMON.MESSAGE'.tr(),
                            content: Text(
                              widget.message ?? '',
                              style: FXUI.subtitleTextStyle.copyWith(
                                color: FXColor.middleBlack,
                              ),
                            ),
                          ),
                      ]));
  }

  List<Widget> getSignTypedDataV1Widget() {
    if (eip712TypedDatas == null) return [];
    return eip712TypedDatas!.map((e) {
      return Padding(
        padding: EdgeInsets.only(
            bottom: eip712TypedDatas!.indexOf(e) != eip712TypedDatas!.length - 1
                ? 8.0
                : 0.0),
        child: getAddressInfoWidget(
          title: e.name,
          content: Text(
            e.value.toString(),
            style: FXUI.subtitleTextStyle.copyWith(
              color: FXColor.middleBlack,
            ),
          ),
        ),
      );
    }).toList();
  }

  List<Widget> getSignTypedDataV3V4WidgetContent(Map<String, dynamic> message) {
    return message.entries
        .map((e) => e.value is Map
            ? Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    '${e.key}:',
                    style: FXUI.subtitleTextStyle.copyWith(
                      color: FXColor.middleBlack,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.only(left: 8.0),
                    child: Column(
                      children: getSignTypedDataV3V4WidgetContent(e.value),
                    ),
                  ),
                ],
              )
            : Column(
                children: [
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '${e.key}:',
                        style: FXUI.subtitleTextStyle.copyWith(
                          color: FXColor.middleBlack,
                        ),
                      ),
                      SizedBox(width: 4),
                      if (e.value is! List)
                        Expanded(
                          child: Text(
                            e.value,
                            maxLines: 2,
                            overflow: TextOverflow.ellipsis,
                            style: FXUI.subtitleTextStyle.copyWith(
                              color: FXColor.middleBlack,
                            ),
                          ),
                        ),
                    ],
                  ),
                  if (e.value is List)
                    Padding(
                      padding: const EdgeInsets.only(left: 8.0),
                      child: Column(
                        children: getSignTypedDataV3V4WidgetContent(
                          (e.value as List).asMap().map(
                              (key, value) => MapEntry(key.toString(), value)),
                        ),
                      ),
                    ),
                ],
              ))
        .toList();
  }

  Widget getTransactionWidget() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        isCheckingData && widget.data != null
            ? getDataColumn()
            : getDetailColumn(),
        if (widget.data != null)
          Padding(
            padding: EdgeInsets.only(top: 4),
            child: TextButton(
              onPressed: () {
                setState(() {
                  isCheckingData = !isCheckingData;
                });
              },
              child: Text(
                isCheckingData
                    ? 'COMMON.BACK'.tr()
                    : 'DAPP_BROWSER.VIEW_DATA'.tr(),
                style: FXUI.subtitleTextStyle.copyWith(
                  fontSize: 12,
                  color: common.getBackGroundColor(),
                ),
              ),
              style: TextButton.styleFrom(
                padding: const EdgeInsets.all(4.0),
                minimumSize: Size.zero,
                tapTargetSize: MaterialTapTargetSize.shrinkWrap,
              ),
            ),
          ),
        if (!isEmptyString(string: getValueString()))
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 25.0),
            child: Column(
              children: [
                Container(
                  decoration: BoxDecoration(
                    border: Border(
                      bottom: BorderSide(
                        color: common.getBackGroundColor(),
                      ),
                    ),
                  ),
                  child: Text(
                    getValueString(),
                    style: FXUI.titleTextStyle.copyWith(
                      fontSize: 24,
                      color: common.getBackGroundColor(),
                    ),
                    textAlign: TextAlign.center,
                  ),
                ),
                SizedBox(height: 4),
                if (usdExchangeRate != null)
                  Text(
                    'USD ${(formattedValue * usdExchangeRate!).toDouble()}',
                    style: FXUI.inputStyle.copyWith(
                      fontWeight: FontWeight.normal,
                    ),
                    textAlign: TextAlign.center,
                  ),
              ],
            ),
          ),
      ],
    );
  }

  Widget getDetailColumn() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (widget.from != null)
          getAddressInfoWidget(
            title: 'TX_PAGE.FROM'.tr(),
            subtitle: widget.from?.eip55TruncatedString ?? '',
            content: isEmptyString(string: balance)
                ? null
                : Text(
                    '${'TX_PAGE.ASSET_BALANCE'.tr()}$balance ${widget.token?.symbol ?? ''}',
                    style: FXUI.subtitleTextStyle.copyWith(
                      color: FXColor.middleBlack,
                    ),
                  ),
          ),
        if (widget.to != null)
          Padding(
            padding: EdgeInsets.only(top: 12.0),
            child: getAddressInfoWidget(
              title: widget.type.toTitle,
              subtitle: widget.to?.eip55TruncatedString ?? '',
            ),
          ),
        Padding(
          padding: EdgeInsets.only(top: 12.0),
          child: getAmountInfoWidget(
            estimatedGasFeeString: getEstimatedGasFeeString(),
            totalString: getTotalValueString(),
          ),
        ),
      ],
    );
  }

  Widget getDataColumn() {
    return Container(
      decoration: BoxDecoration(
        color: FXColor.dividerColor,
        borderRadius: BorderRadius.circular(14),
      ),
      padding: const EdgeInsets.all(16),
      child: Text(
        bytesToHex(
          widget.data ?? Uint8List.fromList([]),
          include0x: true,
        ),
        style: FXUI.subtitleTextStyle.copyWith(
          fontSize: 10,
          color: FXColor.grey44,
        ),
      ),
    );
  }

  Widget getAddressInfoWidget({
    required String title,
    String? subtitle,
    Widget? content,
  }) {
    return Container(
      decoration: BoxDecoration(
        color: FXColor.dividerColor,
        borderRadius: BorderRadius.circular(14),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Padding(
            padding: EdgeInsets.symmetric(
              vertical: 8.0,
              horizontal: 16.0,
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                isEmptyString(string: subtitle)
                    ? Container()
                    : TextButton.icon(
                        style: TextButton.styleFrom(
                          padding: EdgeInsets.zero,
                          minimumSize: Size.zero,
                          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                        ),
                        onPressed: null,
                        icon: Image.asset(
                          'images/icn_address.png',
                          package: 'euruswallet',
                          width: 14,
                        ),
                        label: Text(
                          subtitle ?? '',
                          style: FXUI.subtitleTextStyle.copyWith(
                            fontWeight: FontWeight.normal,
                            color: FXColor.dimGrey,
                          ),
                        ),
                      ),
                Container(
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(10),
                  ),
                  padding:
                      EdgeInsets.symmetric(vertical: 2.0, horizontal: 10.0),
                  child: Text(
                    title,
                    style: FXUI.subtitleTextStyle.copyWith(
                      fontSize: 9,
                      color: FXColor.placeholderGreyColor,
                    ),
                  ),
                ),
              ],
            ),
          ),
          if (content != null)
            Divider(
              color: Colors.white,
              height: 1,
            ),
          if (content != null)
            Padding(
              padding: const EdgeInsets.symmetric(
                vertical: 8.0,
                horizontal: 16.0,
              ),
              child: content,
            ),
        ],
      ),
    );
  }

  Widget getAmountInfoWidget({
    required String estimatedGasFeeString,
    required String totalString,
  }) {
    return Container(
      decoration: BoxDecoration(
        color: FXColor.dividerColor,
        borderRadius: BorderRadius.circular(14),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Padding(
            padding: EdgeInsets.symmetric(
              vertical: 8.0,
              horizontal: 16.0,
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'COMMON.ESTIMATED_GAS_FEE'.tr(),
                  style: FXUI.subtitleTextStyle.copyWith(
                    color: FXColor.middleBlack,
                  ),
                ),
                Text(
                  estimatedGasFeeString,
                  style: FXUI.subtitleTextStyle.copyWith(
                    color: FXColor.middleBlack,
                  ),
                  textAlign: TextAlign.right,
                ),
              ],
            ),
          ),
          Divider(
            color: Colors.white,
            height: 1,
          ),
          Padding(
            padding: EdgeInsets.symmetric(
              vertical: 8.0,
              horizontal: 16.0,
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  widget.type == DappBrowserConfirmationDialogType.approve
                      ? 'DAPP_BROWSER.APPROVED_AMOUNT'.tr()
                      : 'COMMON.TOTAL'.tr(),
                  style: FXUI.subtitleTextStyle.copyWith(
                    color: FXColor.middleBlack,
                  ),
                ),
                SizedBox(width: 8),
                Expanded(
                  child: Text(
                    totalString,
                    style: FXUI.subtitleTextStyle.copyWith(
                      color: FXColor.middleBlack,
                    ),
                    textAlign: TextAlign.right,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Future<EthPrivateKey?> validatePassword() async {
    final password = passwordEditingController.text;
    if (isEmptyString(string: password)) return null;

    if (isCentralized()) {
      final serverAddressPair = await common.getAddressPair(
          email: common.email,
          password: password,
          mnemonic: common.serverMnemonic,
          addressPairType: AddressPairType.paymentPw);

      if (serverAddressPair.address.toLowerCase() !=
          common.ownerWalletAddress?.toLowerCase()) return null;

      return EthPrivateKey.fromHex(serverAddressPair.privateKey);
    } else {
      final privateKey = CommonMethod.passwordDecrypt(
          password, common.currentUserProfile!.encryptedPrivateKey);

      if (privateKey == null) return null;
      return EthPrivateKey.fromHex(privateKey);
    }
  }

  Future<Decimal> getFormattedValue() async {
    final contractAddress =
        Web3Bridge.instance.blockChainType == BlockChainType.Eurus
            ? widget.token?.addressEurus
            : widget.token?.addressEthereum;

    if (contractAddress == null) return value;
    if (contractAddress == '0x0')
      return value / Decimal.fromInt(pow(10, 18).toInt());

    DeployedContract contract = Web3Bridge.instance.blockChainType ==
            BlockChainType.Eurus
        ? web3dart.getEurusERC20Contract(contractAddress: contractAddress)
        : web3dart.getEthereumERC20Contract(contractAddress: contractAddress);
    BigInt? decimal = await web3dart.getContractDecimal(
        deployedContract: contract,
        blockChainType: Web3Bridge.instance.blockChainType);
    return decimal != null
        ? value / Decimal.fromInt(pow(10, decimal.toInt()).toInt())
        : value;
  }

  String getTotalValueString() {
    String? valueSymbol = widget.token?.symbol.toUpperCase();
    String estimatedGasFeeSymbol =
        Web3Bridge.instance.blockChainType == BlockChainType.Eurus
            ? 'EUN'
            : 'ETH';
    bool isValueAndGasFeeSameSymbol = estimatedGasFeeSymbol == valueSymbol;

    return isValueAndGasFeeSameSymbol
        ? (formattedValue + estimatedGasFee).toDouble().toString() +
            (widget.type != DappBrowserConfirmationDialogType.approve
                ? ' ' + estimatedGasFeeSymbol
                : '')
        : formattedValue != Decimal.zero
            ? getValueString() +
                (widget.type != DappBrowserConfirmationDialogType.approve
                    ? ' + ' + getEstimatedGasFeeString()
                    : '')
            : (widget.type != DappBrowserConfirmationDialogType.approve
                ? getEstimatedGasFeeString()
                : '');
  }

  String getValueString() {
    String? valueSymbol = widget.token?.symbol.toUpperCase();
    return (formattedValue.toDouble().toString() +
        (!isEmptyString(string: valueSymbol) ? ' $valueSymbol' : ''));
  }

  String getEstimatedGasFeeString() {
    return estimatedGasFee.toDouble().toString() +
        ' ' +
        (Web3Bridge.instance.blockChainType == BlockChainType.Eurus
            ? 'EUN'
            : 'ETH');
  }

  Future<bool> checkGasBalance() async {
    EtherAmount? currentGasBalance;
    if (isCentralized()) {
      currentGasBalance = await web3dart.eurusEthClient.getBalance(
        EthereumAddress.fromHex(common.ownerWalletAddress ?? ''),
      );
    } else {
      if (web3dart.myEthereumAddress != null)
        currentGasBalance = await web3dart
            .getCurrentClient(
                blockChainType: Web3Bridge.instance.blockChainType)
            .getBalance(web3dart.myEthereumAddress!);
    }

    return currentGasBalance != null &&
        (Decimal.fromBigInt(currentGasBalance.getInWei) /
                Decimal.fromBigInt(web3dart
                    .getGasPrice(
                        blockChainType: Web3Bridge.instance.blockChainType)
                    .getInWei) >=
            Decimal.fromInt(widget.gas ?? 0));
  }
}
