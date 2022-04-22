import 'package:app_authentication_kit/utils/address.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:email_validator/email_validator.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/common/qrcode_scanner.dart';
import 'package:euruswallet/common/web3dart.dart';
import 'package:euruswallet/commonUI/acknowledgementDialog.dart';
import 'package:euruswallet/commonUI/topSelectBlockChainBar.dart';
import 'package:euruswallet/pages/afterScanPage.dart';
import 'package:euruswallet/pages/centralized/selectTargetDetailPage.dart';
import 'package:euruswallet/pages/centralized/top_up_payment_wallet_page.dart';
import 'package:euruswallet/pages/decentralized/selectCurrencyPage.dart';
import 'package:euruswallet/pages/decentralized/transferPage.dart';

class SelectTargetPage extends StatefulWidget {
  final String titleName;
  final String ethereumErc20ContractAddress;
  final String eurusErc20ContractAddress;
  final BlockChainType fromBlockChainType;
  final bool disableSelectBlockchain;
  final String reciverText;

  final Future<String> Function()? canGetPrivateKeyHandler;
  final String? userWalletAccountAvailableAssetsList;

  SelectTargetPage({
    Key? key,
    required this.titleName,
    required this.ethereumErc20ContractAddress,
    required this.eurusErc20ContractAddress,
    required this.fromBlockChainType,
    this.disableSelectBlockchain: false,
    this.reciverText = "",
    this.canGetPrivateKeyHandler,
    this.userWalletAccountAvailableAssetsList,
  }) : super(key: key);

  @override
  _SelectTargetPageState createState() => _SelectTargetPageState();
}

class _SelectTargetPageState extends State<SelectTargetPage> {
  TextEditingController targetController = TextEditingController();
  final _formKey = GlobalKey<FormState>();
  bool alreadyClickNextBtn = false;
  FocusNode textFieldFocusNode = FocusNode();
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  BlockChainType? currentNetworkSelection;
  String? errorText;

  @override
  void initState() {
    super.initState();
    targetController = TextEditingController();
    if (isCentralized()) {
      common.fromBlockChainType = BlockChainType.Eurus;
      _checkTopUpGasBalance();
    }

    if (!isEmptyString(string: widget.reciverText))
      currentNetworkSelection = widget.fromBlockChainType;

    common
        .setUpErc20TokenContract(
            ethereumErc20ContractAddress: widget.ethereumErc20ContractAddress,
            eurusErc20ContractAddress: widget.eurusErc20ContractAddress,
            fromBlockChainType: widget.fromBlockChainType)
        .then((value) {
      setState(() {
        reFreshUI(receiveText: widget.reciverText);
      });
    });
  }

  _checkTopUpGasBalance() async {
    String statmentKey = 'ackTopUpGasStatement';

    String? doNotShowStatement =
        await NormalStorageKit().readValue(statmentKey);
    if (common.getUriVal(
            doNotShowStatement, web3dart.myEthereumAddress?.hex ?? '') !=
        '1') {
      Future.wait([
        web3dart.getMaxTopUpGasAmount(),
        web3dart.eurusEthClient.getBalance(
          EthereumAddress.fromHex(common.ownerWalletAddress ?? ''),
        ),
      ]).then((value) async {
        if (value.isNotEmpty) {
          final maxTopUpGasAmount = value[0] as double?;
          final balance = value[1] as EtherAmount;
          if (maxTopUpGasAmount == null ||
              balance.getInWei.toDouble() /
                      web3dart
                          .getGasPrice(blockChainType: BlockChainType.Eurus)
                          .getInWei
                          .toDouble() <
                  maxTopUpGasAmount.toDouble() * 0.1) {
            bool? result = await showDialog(
              barrierDismissible: false,
              context: context,
              builder: (_) => AcknowledgementDialog(
                statement: 'REFUEL.DIALOG.WARNING_CONTENT'.tr(),
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
                  _checkTopUpGasBalance();
                },
                dontAskAgainText: 'ACKNOWLEDGEMENT_DIALOG.DONT_ASK_AGAIN'.tr(),
              ),
            );
            if (result != null && result == true) {
              String updatedAckState = common.updateUriVal(doNotShowStatement,
                  web3dart.myEthereumAddress?.hex ?? '', '1');
              await NormalStorageKit().setValue(updatedAckState, statmentKey);
            }
          }
        }
      });
    }
  }

  String selectAssetsString() {
    return isEmptyString(string: common.selectTokenSymbol) ||
            (common.isCenWithdraw &&
                common.selectTokenSymbol?.toUpperCase() == 'EUN')
        ? "SEND_PAGE.SELECT_ASSET_PLACEHOLDER".tr()
        : common.selectTokenSymbol ?? "SEND_PAGE.SELECT_ASSET_PLACEHOLDER".tr();
  }

  String? checkReceiverString({String? receiver}) {
    String? errorString;
    if (isEmptyString(string: receiver)) {
      errorString = "SEND_PAGE.ERROR.EMPTY_RECEIVE_ADDRESS".tr();
    } else {
      if (common.isCenWithdraw || !isCentralized()) {
        errorString = EthAddress().isValidEthereumAddress(receiver ?? '')
            ? null
            : "SEND_PAGE.ERROR.ADDRESS_NOT_VALID".tr();
      } else {
        errorString = (EthAddress().isValidEthereumAddress(receiver ?? '') ||
                EmailValidator.validate(receiver ?? ''))
            ? null
            : "SEND_PAGE.ERROR.INPUT_NOT_VALID".tr();
      }
      if (isEmptyString(string: errorString)) {
        errorString = receiver?.toLowerCase() ==
                    web3dart.myEthereumAddress.toString().toLowerCase() ||
                receiver?.toLowerCase() == common.email
            ? "SEND_PAGE.ERROR.CANNOT_SEND_TO_SELF".tr()
            : null;
      }
    }
    return errorString;
  }

  void reFreshUI({String? receiveText}) {
    if (isEmptyString(string: receiveText) && common.isWithdraw()) {
      if (isEmptyString(string: widget.eurusErc20ContractAddress) &&
          isEmptyString(string: widget.ethereumErc20ContractAddress))
        common.selectTokenSymbol = null;
    }

    if (!widget.disableSelectBlockchain &&
        !isEmptyString(string: receiveText)) {
      bool isValidEthereumAddress = true;

      receiveText = receiveText!;
      final isEurusAddress = receiveText.contains('eurus:');
      if (!isCentralized()) {
        if (currentNetworkSelection == BlockChainType.Ethereum &&
            isEurusAddress) {
          isValidEthereumAddress = false;
        }

        // if (receiveText.contains('eurus:')) {
        //   common.fromBlockChainType = BlockChainType.Eurus;
        //   common.selectTokenSymbol = "EUN";
        // } else {
        //   common.fromBlockChainType = BlockChainType.Ethereum;
        //   common.selectTokenSymbol = "ETH";
        // }
        // currently not eurus prefix will withdraw
        //  else if (receiveText.contains('binance coin:')) {
        //   common.fromBlockChainType = BlockChainType.BinanceCoin;
        //   common.selectTokenSymbol = "BNB";
        // }
      } else {
        if (common.isCenWithdraw && isEurusAddress) {
          isValidEthereumAddress = false;
        }
        // common.fromBlockChainType = BlockChainType.Eurus;
        // if (receiveText.contains('eurus:')) {
        //   common.isCenWithdraw = false;
        //   common.selectTokenSymbol = "EUN";
        // } else {
        //   common.isCenWithdraw = true;
        //   common.selectTokenSymbol = "ETH";
        // }
      }

      receiveText = receiveText
          .replaceAll('eurus:', '')
          .replaceAll('ethereum:', '')
          .replaceAll('binance coin:', '')
          .replaceAll('empty', '');
      targetController.text = receiveText;

      if (!EthAddress().isValidEthereumAddress(receiveText))
        isValidEthereumAddress = false;

      if (isValidEthereumAddress) {
        reFreshSelectAssets(type: common.fromBlockChainType);
        print(
            "Send Page assetsList ${widget.userWalletAccountAvailableAssetsList}");

        errorText = null;
      } else {
        errorText = "SEND_PAGE.ERROR.ADDRESS_NOT_VALID".tr();
      }
    }
  }

  void reFreshSelectAssets({BlockChainType? type}) {
    if (!widget.disableSelectBlockchain) {
      type = type == null ? BlockChainType.Eurus : type;
      if (type == BlockChainType.BinanceCoin) {}
    }

    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    return BackGroundImage(
        child: Scaffold(
      backgroundColor: Colors.transparent,
      appBar: WalletAppBar(
          title: widget.titleName,
          rightWidget: isCentralized() ||
                  !isCentralized() && currentNetworkSelection != null
              ? Container(
                  width: 20,
                  height: 20,
                  child: Image.asset(
                    "images/scanQrCode.png",
                    package: 'euruswallet',
                    fit: BoxFit.contain,
                  ))
              : null,
          function: () async {
            FocusScope.of(context).requestFocus(FocusNode());
            String result = await scanQRCode(context: context);
            reFreshUI(receiveText: result);
            print("clickButton");
          }),
      body: SingleChildScrollView(
          child: TopCircularContainer(
              height: size.heightWithoutAppBar,
              width: size.blockSizeHorizontal * 100,
              child: Padding(
                padding: EdgeInsets.only(
                    left: size.leftPadding, right: size.leftPadding, top: 20),
                child: Column(
                  children: [
                    Form(
                        key: _formKey,
                        child: Column(
                          children: [
                            !isCentralized()
                                ? Padding(
                                    padding: EdgeInsets.symmetric(
                                        vertical:
                                            currentNetworkSelection == null
                                                ? 20.0
                                                : 0.0),
                                    child: DropdownButtonHideUnderline(
                                      child: ButtonTheme(
                                        alignedDropdown: true,
                                        child: DropdownButton<BlockChainType>(
                                          alignment: Alignment.center,
                                          dropdownColor: Colors.white,
                                          borderRadius:
                                              BorderRadius.circular(12),
                                          hint: Text(
                                            'SEND_PAGE.SELECT_NETWORK'.tr(),
                                            style: FXUI.titleTextStyle.copyWith(
                                              fontSize: 16,
                                              color: FXColor.deepGreyColor,
                                            ),
                                          ),
                                          value: currentNetworkSelection,
                                          icon: currentNetworkSelection == null
                                              ? Padding(
                                                  padding:
                                                      const EdgeInsets.all(8.0),
                                                  child: Image.asset(
                                                    'images/icn_arrow_down.png',
                                                    package: 'euruswallet',
                                                    width: 16,
                                                    color:
                                                        FXColor.deepGreyColor,
                                                  ),
                                                )
                                              : Container(),
                                          items: [
                                            BlockChainType.Eurus,
                                            BlockChainType.Ethereum,
                                          ]
                                              .map(
                                                (e) => DropdownMenuItem<
                                                    BlockChainType>(
                                                  value: e,
                                                  child: Container(
                                                    decoration: BoxDecoration(
                                                      borderRadius:
                                                          FXUI.cricleRadius,
                                                      color: common
                                                          .getBackGroundColor(),
                                                    ),
                                                    padding:
                                                        EdgeInsets.fromLTRB(
                                                            12, 5, 12, 5),
                                                    child: Row(
                                                      children: [
                                                        Image.asset(
                                                          'images/${e == BlockChainType.Eurus ? 'icn_eun.png' : 'icn_eth.png'}',
                                                          package:
                                                              'euruswallet',
                                                          width: 16,
                                                          height: 16,
                                                          color: Colors.white,
                                                        ),
                                                        SizedBox(
                                                          width: 3,
                                                        ),
                                                        Text(
                                                          getBlockChainName(e),
                                                          style:
                                                              Theme.of(context)
                                                                  .textTheme
                                                                  .bodyText1
                                                                  ?.apply(
                                                                    color: Colors
                                                                        .white,
                                                                    fontWeightDelta:
                                                                        2,
                                                                  )
                                                                  .copyWith(
                                                                      fontSize:
                                                                          13),
                                                        ),
                                                      ],
                                                    ),
                                                  ),
                                                ),
                                              )
                                              .toList(),
                                          onChanged: currentNetworkSelection ==
                                                  null
                                              ? (value) {
                                                  if (value != null) {
                                                    if (value ==
                                                        currentNetworkSelection)
                                                      return;

                                                    setState(
                                                      () {
                                                        currentNetworkSelection =
                                                            value;
                                                        common.currentBlockchainSelection =
                                                            value;
                                                        common.selectTokenSymbol =
                                                            null;

                                                        common.fromBlockChainType =
                                                            value;
                                                        FocusScope.of(context)
                                                            .requestFocus(
                                                                FocusNode());
                                                        reFreshSelectAssets(
                                                            type: value);
                                                      },
                                                    );
                                                  }
                                                }
                                              : null,
                                        ),
                                      ),
                                    ),
                                  )
                                : Container(height: 32),
                            if (!isCentralized() &&
                                currentNetworkSelection == null)
                              Text(
                                'SEND_PAGE.SELECT_NETWORK_STATEMENT'.tr(),
                                style: FXUI.hintStyle.copyWith(
                                  color: FXColor.grey44,
                                  fontSize: 16,
                                  fontWeight: FontWeight.normal,
                                ),
                                textAlign: TextAlign.center,
                              ),
                            if (isCentralized() ||
                                (!isCentralized() &&
                                    currentNetworkSelection != null))
                              Padding(
                                padding: EdgeInsets.only(top: 20, bottom: 10),
                                child: Text("SEND_PAGE.ENTER_INFO".tr(),
                                    style: FXUI.titleTextStyle.copyWith(
                                        color: Colors.black, fontSize: 24)),
                              ),
                            if (isCentralized() ||
                                (!isCentralized() &&
                                    currentNetworkSelection != null))
                              Row(children: [
                                Text(
                                    !isCentralized()
                                        ? "SEND_PAGE.RECEIVING_ADDRESS.LABEL"
                                            .tr()
                                        : common.isCenWithdraw
                                            ? "SEND_PAGE.RECEIVING_ADDRESS.LABEL"
                                                .tr()
                                            : "SEND_PAGE.RECEIVING_ADDRESS.CENLABEL"
                                                .tr(),
                                    style: FXUI.titleTextStyle.copyWith(
                                        color: FXColor.lightGray, fontSize: 14))
                              ]),
                            if (isCentralized() ||
                                (!isCentralized() &&
                                    currentNetworkSelection != null))
                              Padding(
                                padding: EdgeInsets.only(top: 16, bottom: 16),
                                child: TextFormField(
                                  style: FXUI.inputStyle,
                                  focusNode: textFieldFocusNode,
                                  autofocus: true,
                                  controller: targetController,
                                  validator: (value) {
                                    return errorText ??
                                        checkReceiverString(receiver: value);
                                  },
                                  decoration: FXUI.inputDecoration.copyWith(
                                    hintText:
                                        "SEND_PAGE.RECEIVING_ADDRESS.PLACEHOLDER"
                                            .tr(),
                                    // errorText: errorText,
                                  ),
                                  onChanged: (value) => errorText = null,
                                ),
                              ),
                            if (isCentralized() ||
                                (!isCentralized() &&
                                    currentNetworkSelection != null))
                              Padding(
                                padding: EdgeInsets.only(bottom: 10),
                                child: Row(
                                  crossAxisAlignment: CrossAxisAlignment.center,
                                  children: [
                                    Text("SEND_PAGE.ASSET_PREFIX".tr(),
                                        style: FXUI.titleTextStyle.copyWith(
                                            color: FXColor.lightGray,
                                            fontSize: 14)),
                                    InkWell(
                                        onTap: widget.disableSelectBlockchain ==
                                                true
                                            ? null
                                            : () async {
                                                await common.pushPage(
                                                    page: SelectCurrencyPage(
                                                      titleName:
                                                          "SEND_PAGE.MAIN_TITLE"
                                                              .tr(),
                                                      userWalletAccountAvailableAssetsList:
                                                          widget
                                                              .userWalletAccountAvailableAssetsList,
                                                    ),
                                                    context: context);
                                                setState(() {});
                                              },
                                        child: Padding(
                                            padding: EdgeInsets.only(left: 10),
                                            child: Row(
                                              children: [
                                                Text(selectAssetsString(),
                                                    style: FXUI.normalTextStyle
                                                        .copyWith(
                                                            color:
                                                                Colors.black)),
                                                widget.disableSelectBlockchain ==
                                                        true
                                                    ? Container()
                                                    : Icon(
                                                        Icons
                                                            .arrow_drop_down_sharp,
                                                        color: Colors.black,
                                                        size: 20.0,
                                                      ),
                                              ],
                                            ))),
                                  ],
                                ),
                              ),
                            if (isCentralized() ||
                                (!isCentralized() &&
                                    currentNetworkSelection != null))
                              Padding(
                                  padding: EdgeInsets.only(bottom: 0),
                                  child: isEmptyString(
                                              string:
                                                  common.selectTokenSymbol ??
                                                      '') &&
                                          alreadyClickNextBtn
                                      ? Row(
                                          children: [
                                            Padding(
                                              padding:
                                                  EdgeInsets.only(left: 15),
                                              child: Text(
                                                  "SEND_PAGE.ERROR.EMPTY.ASSET_TYPE"
                                                      .tr(),
                                                  style: FXUI.normalTextStyle
                                                      .copyWith(
                                                          color: FXColor
                                                              .alertRedColor,
                                                          fontSize: 12)),
                                            ),
                                          ],
                                        )
                                      : Container(
                                          padding: EdgeInsets.only(bottom: 0))),
                            if (isCentralized() ||
                                (!isCentralized() &&
                                    currentNetworkSelection != null))
                              widget.disableSelectBlockchain
                                  ? Padding(
                                      padding: EdgeInsets.only(bottom: 10),
                                      child: Row(
                                        crossAxisAlignment:
                                            CrossAxisAlignment.center,
                                        children: [
                                          Text(
                                              '(${getBlockChainName(common.fromBlockChainType)})',
                                              style: FXUI.titleTextStyle
                                                  .copyWith(
                                                      color: FXColor.lightGray,
                                                      fontSize: 14)),
                                        ],
                                      ),
                                    )
                                  : Container(),
                            if (isCentralized() ||
                                (!isCentralized() &&
                                    currentNetworkSelection != null))
                              Padding(
                                padding: EdgeInsets.only(
                                    top: size.screenHeight > 800 ? 120 : 80),
                                child: SubmitButton(
                                  loadingSecond: 4,
                                  btnController: btnController,
                                  label: "COMMON.NEXT_STEP".tr(),
                                  onPressed: () async {
                                    if ((common.selectTokenSymbol == "EUN") &&
                                        common.isCenWithdraw) {
                                      common.showPopUpError(
                                          context: context,
                                          descriptions:
                                              "COMMON_ERROR.EUN_WITHDRAW_ERROR"
                                                  .tr());
                                      return;
                                    }
                                    setState(() {
                                      alreadyClickNextBtn = true;
                                    });

                                    if (_formKey.currentState != null &&
                                        _formKey.currentState!.validate() &&
                                        common.selectTokenSymbol != null) {
                                      common.findEmailWalletAddress = null;
                                      if (common.currentUserType ==
                                          CurrentUserType.decentralized) {
                                        if (EmailValidator.validate(
                                            targetController.text)) {
                                          common.findEmailWalletAddress =
                                              await api.findEmailWalletAddress(
                                                  targetAddress:
                                                      targetController.text);
                                          if (await common.checkApiError(
                                              context: context,
                                              errorString: common
                                                  .findEmailWalletAddress
                                                  ?.message,
                                              returnCode: common
                                                  .findEmailWalletAddress
                                                  ?.returnCode)) {
                                            common.targetAddress = common
                                                .findEmailWalletAddress
                                                ?.data
                                                ?.walletAddress;
                                            common.pushPage(
                                                page: TransferPage(
                                                    titleName:
                                                        "SEND_PAGE.MAIN_TITLE"
                                                            .tr(),
                                                    fromBlockChainType: common
                                                        .fromBlockChainType),
                                                context: context);
                                          }
                                        } else {
                                          common.targetAddress =
                                              targetController.text;
                                          common.pushPage(
                                              page: TransferPage(
                                                  titleName:
                                                      "SEND_PAGE.MAIN_TITLE"
                                                          .tr(),
                                                  fromBlockChainType: common
                                                      .fromBlockChainType),
                                              context: context);
                                        }
                                      } else {
                                        if (common.isCenWithdraw) {
                                          common.targetAddress =
                                              targetController.text;
                                          common.pushPage(
                                              page: TransferPage(
                                                  titleName: widget.titleName,
                                                  fromBlockChainType:
                                                      common.fromBlockChainType,
                                                  transferToMySelf: true),
                                              context: context);
                                        } else {
                                          if (EmailValidator.validate(
                                              targetController.text)) {
                                            common.findEmailWalletAddress =
                                                await api
                                                    .findEmailWalletAddress(
                                                        targetAddress:
                                                            targetController
                                                                .text);
                                            if (await common.checkApiError(
                                                apiName:
                                                    "findEmailWalletAddress",
                                                context: context,
                                                errorString: common
                                                    .findEmailWalletAddress
                                                    ?.message,
                                                returnCode: common
                                                    .findEmailWalletAddress
                                                    ?.returnCode)) {
                                              common.targetAddress = common
                                                  .findEmailWalletAddress
                                                  ?.data
                                                  ?.walletAddress;
                                              common.pushPage(
                                                  page: SelectTargetDetailPage(
                                                    titleName:
                                                        "SEND_PAGE.MAIN_TITLE"
                                                            .tr(),
                                                  ),
                                                  context: context);
                                            }
                                          } else {
                                            common.targetAddress =
                                                targetController.text;
                                            common.pushPage(
                                                page: TransferPage(
                                                  titleName:
                                                      "SEND_PAGE.MAIN_TITLE"
                                                          .tr(),
                                                  fromBlockChainType:
                                                      common.fromBlockChainType,
                                                  transferToMySelf: false,
                                                ),
                                                context: context);
                                          }
                                        }
                                      }
                                    }
                                    btnController.reset();
                                  },
                                ),
                              )
                          ],
                        ))
                  ],
                ),
              ))),
    ));
  }
}
