import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/topBlockchainBar.dart';
import 'package:euruswallet/commonUI/walletLockerPWDialog.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:percent_indicator/circular_percent_indicator.dart';

import 'transferSuccessfulPage.dart';

class ConfirmationPage extends StatefulWidget {
  final String titleName;
  final TransactionType transactionType;
  final int confirmBlockNumber;
  final int requestSuccessBlock;
  final String statusString;
  ConfirmationPage({
    Key? key,
    this.titleName: '',
    this.transactionType: TransactionType.confirmation,
    this.confirmBlockNumber: 0,
    this.statusString: '',
    this.requestSuccessBlock: 3,
  }) : super(key: key);

  @override
  _ConfirmationPageState createState() => _ConfirmationPageState();
}

class _ConfirmationPageState extends State<ConfirmationPage> {
  String? lastTxId;

  bool get isAllocation => common.isDeposit() || common.isWithdraw();
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    super.initState();
    setState(() {
      common.targetDepositOrWidthDrawAddress();
    });
  }

  @override
  void dispose() {
    super.dispose();
  }

  String getFromAddress() {
    String isDepositOrWidthDraw = '';
    if (common.isDeposit()) {
      isDepositOrWidthDraw = 'TX_PAGE.MY_WALLET'.tr(args: ["Ethereum"]);
    } else if (common.isWithdraw()) {
      isDepositOrWidthDraw = 'TX_PAGE.MY_WALLET'.tr(args: ["Eurus"]);
    }
    return isDepositOrWidthDraw;
  }

  void showError({String? errorMessage}) {
    showDialog(
        context: context,
        builder: (BuildContext context) => CupertinoAlertDialog(
              title: Text("Transfer Fail"),
              content: Text(errorMessage ?? ""),
              actions: <Widget>[
                CupertinoDialogAction(
                  isDefaultAction: true,
                  child: Text("OK"),
                  onPressed: () {
                    Navigator.pop(context);
                  },
                ),
              ],
            ));
  }

  Future<String?> _requestLockerPWDialog() async {
    final _tc = TextEditingController();

    final _submit = (tec) async {
      TextEditingController txt = tec;
      common.serverAddressPair = await common.getAddressPair(
          email: common.email,
          password: txt.text,
          mnemonic: common.serverMnemonic,
          addressPairType: AddressPairType.paymentPw);
      if (common.serverAddressPair?.address.toUpperCase() ==
          common.ownerWalletAddress?.toUpperCase()) {
        common.cenSignKey = common.serverAddressPair?.privateKey;
        return common.cenSignKey;
      } else {
        common.cenSignKey = null;
        return null;
      }
    };

    String? decryptedString = await Navigator.of(context)
        .push(PageRouteBuilder(
          fullscreenDialog: true,
          opaque: false,
          pageBuilder: (pageBuilderContext, animation, secondaryAnimation) =>
              WalletLockerPWDialog(
                  themeColor: common.getBackGroundColor(),
                  textEditingController: _tc,
                  cenSubmitFnc: _submit,
                  tryBioAuthFnc: null),
        ))
        .then((value) => value);

    return isEmptyString(string: decryptedString) ? null : _tc.text;
  }

  double getPercent() {
    return widget.confirmBlockNumber / widget.requestSuccessBlock > 1.0
        ? 1.0
        : widget.confirmBlockNumber / widget.requestSuccessBlock;
  }

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);

    return WillPopScope(
        onWillPop: () async =>
            widget.transactionType == TransactionType.confirmation,
        child: BackGroundImage(
            child: Scaffold(
                backgroundColor: Colors.transparent,
                appBar: WalletAppBar(
                    backButton:
                        widget.transactionType == TransactionType.confirmation,
                    title:
                        widget.transactionType == TransactionType.confirmation
                            ? common.transferToMySelf
                                ? "TX_PAGE.ALLOCATE_CONFIRMATION_TITLE".tr()
                                : "TX_PAGE.TRANSFER_CONFIRMATION_TITLE".tr()
                            : widget.transactionType == TransactionType.pending
                                ? "COMMON.PENDING".tr()
                                : "COMMON.PENDING".tr()),
                body: SingleChildScrollView(
                  child: Padding(
                      padding: EdgeInsets.only(
                          top: 8,
                          left: size.leftPadding,
                          right: size.leftPadding),
                      child: Container(
                          decoration: new BoxDecoration(
                            borderRadius: BorderRadius.all(Radius.circular(20)),
                            color: Colors.white,
                          ),
                          child: Column(
                            children: [
                              TopBlockchainBar(),
                              SizedBox(
                                width: double.infinity,
                                child: Image.asset(
                                  "images/line.png",
                                  package: 'euruswallet',
                                ),
                              ),
                              widget.transactionType ==
                                          TransactionType.pending ||
                                      widget.transactionType ==
                                          TransactionType.allocationProcessing
                                  ? Padding(
                                      padding: EdgeInsets.only(
                                          top: 30.0, bottom: 8.0),
                                      child: Container(
                                        alignment: Alignment.center,
                                        child: CircularPercentIndicator(
                                          animation: true,
                                          animationDuration:
                                              common.fromBlockChainType ==
                                                      BlockChainType.Ethereum
                                                  ? 6000
                                                  : 2000,
                                          animateFromLastPercent: true,
                                          radius: 60.0,
                                          lineWidth: 5.0,
                                          percent: getPercent(),
                                          center: new Text(
                                              "${(getPercent() * 100).round()}%"),
                                          progressColor: common
                                              .getBackGroundColor()
                                              .withOpacity(0.5),
                                        ),
                                      ),
                                    )
                                  : Container(),
                              widget.transactionType ==
                                          TransactionType.pending ||
                                      widget.transactionType ==
                                          TransactionType.allocationProcessing
                                  ? Text(widget.statusString,
                                      style: FXUI.titleTextStyle.copyWith(
                                          fontSize: 14,
                                          color: common
                                              .getBackGroundColor()
                                              .withOpacity(0.5)))
                                  : Container(),
                              Padding(
                                  padding: EdgeInsets.only(
                                      left: size.leftPadding,
                                      right: size.leftPadding),
                                  child: Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      Padding(
                                        padding: EdgeInsets.only(top: 20),
                                        child: Text("TX_PAGE.FROM".tr(),
                                            style: FXUI.normalTextStyle
                                                .copyWith(
                                                    color: FXColor.lightGray)),
                                      ),
                                      Padding(
                                        padding: EdgeInsets.only(top: 10),
                                        child: Row(
                                          crossAxisAlignment:
                                              CrossAxisAlignment.start,
                                          children: [
                                            Padding(
                                              padding: EdgeInsets.only(
                                                  top: 5, right: 5),
                                              child: SizedBox(
                                                width: 15,
                                                height: 7.5,
                                                child: Image.asset(
                                                  "images/addressIcon.png",
                                                  package: 'euruswallet',
                                                  fit: BoxFit.contain,
                                                ),
                                              ),
                                            ),
                                            Expanded(
                                                child: Text(
                                                    common.transferToMySelf
                                                        ? getFromAddress()
                                                        : web3dart
                                                            .myEthereumAddress
                                                            .toString(),
                                                    maxLines: 2,
                                                    overflow:
                                                        TextOverflow.ellipsis,
                                                    style: FXUI.normalTextStyle
                                                        .copyWith(
                                                            fontSize: 14,
                                                            color: FXColor
                                                                .centralizedGrayTextColor))),
                                          ],
                                        ),
                                      ),
                                    ],
                                  )),
                              Padding(
                                padding: EdgeInsets.only(top: 20),
                                child: SizedBox(
                                  width: 4000,
                                  height: 3,
                                  child: Image.asset(
                                    "images/line.png",
                                    package: 'euruswallet',
                                    fit: BoxFit.contain,
                                  ),
                                ),
                              ),
                              Padding(
                                padding: EdgeInsets.only(
                                    left: size.leftPadding,
                                    right: size.leftPadding),
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Padding(
                                      padding: EdgeInsets.only(top: 20),
                                      child: Text(
                                          common.transferToMySelf &&
                                                  !isCentralized()
                                              ? "TX_PAGE.ALLOCATE_TO".tr()
                                              : "TX_PAGE.TRANSFER_TO".tr(),
                                          style: FXUI.normalTextStyle.copyWith(
                                              color: FXColor.lightGray)),
                                    ),
                                    Padding(
                                      padding: EdgeInsets.only(top: 10),
                                      child: Row(
                                        crossAxisAlignment:
                                            CrossAxisAlignment.start,
                                        children: [
                                          Padding(
                                            padding: EdgeInsets.only(
                                                top: 5, right: 5),
                                            child: SizedBox(
                                              width: 15,
                                              height: 7.5,
                                              child: Image.asset(
                                                "images/addressIcon.png",
                                                package: 'euruswallet',
                                                fit: BoxFit.contain,
                                              ),
                                            ),
                                          ),
                                          Expanded(
                                              child: Text(
                                                  common.targetDepositOrWidthDrawAddresss ==
                                                              '' ||
                                                          isCentralized()
                                                      ? common.targetAddress ??
                                                          ""
                                                      : common
                                                          .targetDepositOrWidthDrawAddresss,
                                                  maxLines: 2,
                                                  overflow:
                                                      TextOverflow.ellipsis,
                                                  style: FXUI.normalTextStyle
                                                      .copyWith(
                                                          fontSize: 14,
                                                          color: FXColor
                                                              .centralizedGrayTextColor))),
                                        ],
                                      ),
                                    ),
                                    Column(
                                      crossAxisAlignment:
                                          CrossAxisAlignment.center,
                                      children: [
                                        common.isWithdraw()
                                            ? Padding(
                                                padding:
                                                    EdgeInsets.only(top: 34),
                                                child: Text(
                                                    "ADMIN_FEE.NAME".tr(),
                                                    style: FXUI.normalTextStyle
                                                        .copyWith(
                                                            color: FXColor
                                                                .lightGray)),
                                              )
                                            : Container(),
                                        common.isWithdraw()
                                            ? Padding(
                                                padding:
                                                    EdgeInsets.only(top: 10),
                                                child: Text(
                                                  common.numberFormat(
                                                          maxDecimal: 12,
                                                          number: common
                                                              .adminFeeModel
                                                              ?.data
                                                              ?.actualFee
                                                              ?.toStringAsFixed(
                                                                  12)) +
                                                      " " +
                                                      (common.selectTokenSymbol ??
                                                          ''),
                                                  style: FXUI.normalTextStyle
                                                      .copyWith(
                                                    fontSize: 14,
                                                    color: FXColor.textGray,
                                                  ),
                                                ),
                                              )
                                            : Container(),
                                        Padding(
                                          padding: EdgeInsets.only(
                                              top: common.isWithdraw()
                                                  ? 20
                                                  : 34),
                                          child: Text(
                                              "TX_PAGE.E_GAS_PRICE".tr(),
                                              style: FXUI.normalTextStyle
                                                  .copyWith(
                                                      color:
                                                          FXColor.lightGray)),
                                        ),
                                        Padding(
                                          padding: EdgeInsets.only(top: 10),
                                          child: Text(
                                              common.numberFormat(
                                                      maxDecimal: 12,
                                                      number: common.currentGas
                                                          ?.toStringAsFixed(
                                                              12)) +
                                                  (getSymbolByBlockChainType(
                                                      common
                                                          .fromBlockChainType)),
                                              style:
                                                  FXUI.normalTextStyle.copyWith(
                                                fontSize: 14,
                                                color: FXColor.textGray,
                                              )),
                                        ),
                                        Padding(
                                            padding: EdgeInsets.only(top: 20),
                                            child: Row(
                                              mainAxisAlignment:
                                                  MainAxisAlignment.center,
                                              children: [
                                                Text(
                                                    common.transferToMySelf
                                                        ? 'TX_PAGE.ALLOCATE_AMOUNT'
                                                            .tr()
                                                        : 'TX_PAGE.TRANSFER_AMOUNT'
                                                            .tr(),
                                                    style: FXUI.normalTextStyle
                                                        .copyWith(
                                                            color: FXColor
                                                                .lightGray,
                                                            fontSize: 14)),
                                              ],
                                            )),
                                        Padding(
                                          padding: EdgeInsets.only(
                                              top: 10, bottom: 35),
                                          child: Text(
                                            "${common.numberFormat(maxDecimal: 12, number: common.transferAmount)} " +
                                                (common.selectTokenSymbol ??
                                                    ''),
                                            style: FXUI.titleTextStyle.copyWith(
                                              fontSize: 40,
                                              color: FXColor.textGray,
                                            ),
                                            textAlign: TextAlign.center,
                                          ),
                                        ),
                                      ],
                                    ),
                                    widget.transactionType ==
                                            TransactionType.confirmation
                                        ? Padding(
                                            padding: EdgeInsets.only(
                                                top: 30, bottom: 27),
                                            child: SubmitButton(
                                              btnController: btnController,
                                              label: common.transferToMySelf
                                                  ? 'TX_PAGE.ALLOCATE_NOW'.tr()
                                                  : 'TX_PAGE.SEND_NOW'.tr(),
                                              onPressed: () async {
                                                // EasyLoading.show(
                                                //     status:
                                                //         'COMMON.LOADING_W_DOT'
                                                //             .tr());
                                                // Future.delayed(
                                                //     Duration(
                                                //         milliseconds: 6000),
                                                //     () async {
                                                //   EasyLoading.dismiss();
                                                // });
                                                try {
                                                  double amount = double.parse(
                                                      common.transferAmount ??
                                                          '');
                                                  String toAddress = '';
                                                  if (isCentralized()) {
                                                    toAddress =
                                                        common.targetAddress ??
                                                            '';
                                                    await _requestLockerPWDialog();
                                                    if (common.cenSignKey !=
                                                        null) {
                                                      if (common.isWithdraw() &&
                                                          common
                                                              .isCenWithdraw) {
                                                        Transaction
                                                            transaction =
                                                            await web3dart
                                                                .cenSubmitWithdraw(
                                                          userWalletAddress:
                                                              common.cenUserWalletAddress ??
                                                                  '',
                                                          deployedContract: web3dart
                                                              .erc20ContractFromEurus,
                                                          enterAmount: amount,
                                                          toAddress: toAddress,
                                                          enterAmountWithFee: amount +
                                                              (common
                                                                      .adminFeeModel
                                                                      ?.data
                                                                      ?.actualFee ??
                                                                  0.0),
                                                          selectTokenSymbol:
                                                              common.selectTokenSymbol ??
                                                                  '',
                                                          maxGas: web3dart
                                                              .estimateMaxGas,
                                                        );

                                                        lastTxId = await web3dart
                                                            .broadCastTranscation(
                                                          credentials: EthPrivateKey
                                                              .fromHex(common
                                                                  .cenSignKey!),
                                                          transaction:
                                                              transaction,
                                                          blockChainType:
                                                              BlockChainType
                                                                  .Eurus,
                                                        );
                                                      } else {
                                                        Transaction
                                                            transaction =
                                                            await web3dart
                                                                .requestTransfer(
                                                          selectTokenSymbol:
                                                              common.selectTokenSymbol ??
                                                                  '',
                                                          userWalletAddress:
                                                              common.cenUserWalletAddress ??
                                                                  '',
                                                          deployedContract: web3dart
                                                              .erc20ContractFromEurus,
                                                          enterAmount: amount,
                                                          toAddress: toAddress,
                                                          blockChainType:
                                                              BlockChainType
                                                                  .Eurus,
                                                          maxGas: web3dart
                                                              .estimateMaxGas,
                                                        );

                                                        lastTxId = await web3dart
                                                            .broadCastTranscation(
                                                          credentials: EthPrivateKey
                                                              .fromHex(common
                                                                  .cenSignKey!),
                                                          transaction:
                                                              transaction,
                                                          blockChainType:
                                                              BlockChainType
                                                                  .Eurus,
                                                        );
                                                      }
                                                    }
                                                  } else {
                                                    if (common.isWithdraw()) {
                                                      toAddress = web3dart
                                                              .myEthereumAddress
                                                              ?.hex ??
                                                          '';
                                                      Transaction? transaction =
                                                          await web3dart
                                                              .submitWithdrawERC20(
                                                        deployedContract: web3dart
                                                            .erc20ContractFromEurus,
                                                        enterAmount: amount,
                                                        toAddress: toAddress,
                                                        enterAmountWithFee: amount +
                                                            (common
                                                                    .adminFeeModel
                                                                    ?.data
                                                                    ?.actualFee ??
                                                                0.0),
                                                        maxGas: web3dart
                                                            .estimateMaxGas,
                                                      );
                                                      final credentials = web3dart
                                                              .credentials ??
                                                          await web3dart
                                                              .canGetCredentialsHandler();
                                                      if (credentials != null &&
                                                          transaction != null) {
                                                        lastTxId = await web3dart
                                                            .broadCastTranscation(
                                                          credentials:
                                                              credentials,
                                                          transaction:
                                                              transaction,
                                                          blockChainType:
                                                              BlockChainType
                                                                  .Eurus,
                                                        );
                                                        print(
                                                            "sendERC20 result:$lastTxId");
                                                      }
                                                    } else {
                                                      if (common.isDeposit()) {
                                                        toAddress = await web3dart
                                                                .getEurusUserDepositAddress() ??
                                                            '';
                                                      } else if (!common
                                                          .transferToMySelf) {
                                                        toAddress = common
                                                                .targetAddress ??
                                                            '';
                                                      }

                                                      if (common.isEthOrEun() ||
                                                          common.isBSC()) {
                                                        Transaction?
                                                            transaction =
                                                            await web3dart
                                                                .sendETH(
                                                          enterAmount: amount,
                                                          toAddress: toAddress,
                                                          type: common
                                                              .fromBlockChainType,
                                                          maxGas: web3dart
                                                              .estimateMaxGas,
                                                        );

                                                        final credentials = web3dart
                                                                .credentials ??
                                                            await web3dart
                                                                .canGetCredentialsHandler();
                                                        if (credentials !=
                                                                null &&
                                                            transaction !=
                                                                null) {
                                                          lastTxId = await web3dart
                                                              .broadCastTranscation(
                                                            credentials:
                                                                credentials,
                                                            transaction:
                                                                transaction,
                                                            blockChainType: common
                                                                .fromBlockChainType,
                                                          );
                                                          print(
                                                              "sendETHTransaction resultString:$lastTxId");
                                                        }
                                                      } else {
                                                        Transaction?
                                                            transaction =
                                                            await web3dart
                                                                .sendERC20(
                                                          deployedContract: web3dart
                                                              .getERC20Contract(
                                                                  common
                                                                      .fromBlockChainType),
                                                          enterAmount: amount,
                                                          toAddress: toAddress,
                                                          blockChainType: common
                                                              .fromBlockChainType,
                                                          maxGas: web3dart
                                                              .estimateMaxGas,
                                                        );
                                                        final credentials = web3dart
                                                                .credentials ??
                                                            await web3dart
                                                                .canGetCredentialsHandler();
                                                        if (credentials !=
                                                                null &&
                                                            transaction !=
                                                                null) {
                                                          lastTxId = await web3dart
                                                              .broadCastTranscation(
                                                            credentials:
                                                                credentials,
                                                            transaction:
                                                                transaction,
                                                            blockChainType: common
                                                                .fromBlockChainType,
                                                          );
                                                          print(
                                                              "sendERC20 result:$lastTxId");
                                                        }
                                                      }
                                                    }
                                                  }
                                                } on Exception catch (exception) {
                                                  print("exception:$exception");
                                                  showError(
                                                      errorMessage:
                                                          exception.toString());
                                                }
                                                if (lastTxId != null &&
                                                    lastTxId!.contains("0x") &&
                                                    lastTxId!.length > 20) {
                                                  EasyLoading.dismiss();
                                                  web3dart.lastTxId = lastTxId;

                                                  // TxRecordModel txRd = TxRecordModel(
                                                  //     transactionHash:
                                                  //         web3dart.lastTxId ??
                                                  //             '',
                                                  //     decodedInputAmount:
                                                  //         amount)
                                                  //   ..txTo = common
                                                  //               .isEthOrEun() ||
                                                  //           common.isBSC()
                                                  //       ? toAddress
                                                  //       : web3dart
                                                  //           .getERC20Contract(common
                                                  //               .fromBlockChainType)
                                                  //           .address
                                                  //           .hex
                                                  //   ..decodedInputRecipientAddress =
                                                  //       toAddress
                                                  //   ..txFrom = web3dart
                                                  //       .myEthereumAddress?.hex
                                                  //   ..eurusTxType = isAllocation
                                                  //       ? common.isDeposit()
                                                  //           ? 2
                                                  //           : 3
                                                  //       : null
                                                  //   ..adminFee = isAllocation &&
                                                  //           common.isWithdraw()
                                                  //       ? common.adminFeeModel
                                                  //           ?.data?.actualFee
                                                  //           ?.toStringAsFixed(
                                                  //               12)
                                                  //       : null
                                                  //   ..chain =
                                                  //       getSymbolByBlockChainType(
                                                  //               common
                                                  //                   .fromBlockChainType)
                                                  //           .toLowerCase();
                                                  // await TxRecordHandler()
                                                  //     .addSentTx(txRd);

                                                  print(getFromAddress());

                                                  print("lastTxId:$lastTxId");
                                                  common.targetType =
                                                      TargetType.ethAddress;
                                                  common.pushPage(
                                                      page:
                                                          TransferSuccessfulPage(),
                                                      context: context);
                                                }
                                                btnController.reset();
                                              },
                                            ),
                                          )
                                        : Padding(
                                            padding: EdgeInsets.only(
                                                top: 30, bottom: 27),
                                            child: SubmitButton(
                                              btnController: btnController,
                                              label: "COMMON.CLOSE".tr(),
                                              onPressed: () {
                                                btnController.reset();
                                                common.isTransactionHistory(
                                                        transactionType: widget
                                                            .transactionType)
                                                    ? Navigator.pop(context)
                                                    : Navigator.popUntil(
                                                        context,
                                                        ModalRoute.withName(
                                                            "HomePage"));
                                              },
                                            ),
                                          )
                                  ],
                                ),
                              ),
                            ],
                          ))),
                ))));
  }
}
