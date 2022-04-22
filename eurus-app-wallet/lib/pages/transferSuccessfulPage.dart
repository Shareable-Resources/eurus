import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/customDialogBox.dart';
import 'package:euruswallet/commonUI/topBlockchainBar.dart';
import 'package:euruswallet/model/getTransactionReceipt.dart';
import 'package:euruswallet/pages/webViewPage.dart';
import 'package:intl/intl.dart';
import 'package:path_provider/path_provider.dart';
import 'package:screenshot/screenshot.dart';
import 'package:share/share.dart';

import 'confirmationPage.dart';

class TransferSuccessfulPage extends StatefulWidget {
  TransactionType transactionType;
  final String centerMessage;
  final String fromAddress;
  final String toAddress;
  final String date;
  final String txId;
  final String? gasFeeString;
  final String? adminFeeString;
  final String transferAmount;
  final BlockChainType? fromBlockChainType;
  final BlockChainType? toBlockChainType;
  final String? eurusTxType;
  TransferSuccessfulPage({
    Key? key,
    this.transactionType: TransactionType.pending,
    this.centerMessage: "",
    this.transferAmount: "",
    this.gasFeeString,
    this.adminFeeString,
    this.fromAddress: "",
    this.toAddress: "",
    this.date: "",
    this.txId: "",
    this.fromBlockChainType,
    this.toBlockChainType,
    this.eurusTxType,
  }) : super(key: key);
  @override
  _TransferSuccessfulPageState createState() => _TransferSuccessfulPageState();
}

class _TransferSuccessfulPageState extends State<TransferSuccessfulPage> {
  Timer? getBlockNumberTimer;
  int requestSuccessBlock = 3;
  String statusString = "COMMON.PENDING".tr();
  int confirmBlockNumber = 0;
  double? gasFee;
  String gasFeeString = "";
  String? formattedDate;
  ScreenshotController screenshotController = ScreenshotController();
  String fromAddress = "";
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  bool get isAllocate =>
      widget.transactionType == TransactionType.allocationProcessing ||
      widget.transactionType ==
          TransactionType.allocationTransactionHistoryPendingStatus ||
      widget.transactionType ==
          TransactionType.allocationTransactionHistoryProcessingStatus ||
      widget.transactionType ==
          TransactionType.allocationTransactionHistorySuccessfulStatus;

  @override
  void initState() {
    super.initState();
    common.fromBlockChainType =
        widget.fromBlockChainType ?? common.fromBlockChainType;
    web3dart.lastTxId =
        common.isTransactionHistory(transactionType: widget.transactionType)
            ? widget.txId
            : web3dart.lastTxId;

    if (widget.transactionType == TransactionType.successful ||
        widget.transactionType ==
            TransactionType.sendTransactionHistorySuccessfulStatus ||
        widget.transactionType ==
            TransactionType.allocationTransactionHistorySuccessfulStatus) {
      changeToSuccessfulStatus();
    } else {
      Future.delayed(const Duration(seconds: 2), () {
        checkBlockConfirmation();
      });
      new Timer.periodic(const Duration(seconds: 1), changeStatusString);
    }
  }

  @override
  void dispose() {
    getBlockNumberTimer?.cancel();
    super.dispose();
  }

  void changeStatusString(Timer t) async {
    common.targetDepositOrWidthDrawAddresss = isCentralized()
        ? widget.toAddress
        : common.targetDepositOrWidthDrawAddress();

    fromAddress = web3dart.myEthereumAddress.toString();
    if (mounted) {
      setState(() {
        if (widget.transactionType == TransactionType.successful) {
          statusString = "COMMON.SUCCESS".tr();
        } else if (widget.transactionType ==
            TransactionType.allocationProcessing) {
          /// TODO: Display different status description on asset allocation
          if (statusString == '${"COMMON.PENDING".tr()}' ||
              statusString == "${"COMMON.PENDING".tr()}." ||
              statusString == "${"COMMON.PENDING".tr()}..") {
            statusString += ".";
          } else if (statusString == "${"COMMON.PENDING".tr()}...") {
            statusString = "${"COMMON.PENDING".tr()}";
          }
        } else if (widget.transactionType == TransactionType.pending) {
          if (statusString == '${"COMMON.PENDING".tr()}' ||
              statusString == "${"COMMON.PENDING".tr()}." ||
              statusString == "${"COMMON.PENDING".tr()}..") {
            statusString += ".";
          } else if (statusString == "${"COMMON.PENDING".tr()}...") {
            statusString = "${"COMMON.PENDING".tr()}";
          }
        }
      });
    }
  }

  checkBlockConfirmation() async {
    getBlockNumberTimer = Timer.periodic(Duration(seconds: 3), (timer) async {
      var client =
          web3dart.getCurrentClient(blockChainType: common.fromBlockChainType);
      var transactionInformation =
          await client.getTransactionReceipt(web3dart.lastTxId ?? "");
      if (transactionInformation != null) {
        int currentBlockNumber = await client.getBlockNumber();
        if (transactionInformation.blockNumber.useAbsolute) {
          confirmBlockNumber =
              currentBlockNumber - transactionInformation.blockNumber.blockNum;
          if (confirmBlockNumber >= requestSuccessBlock) {
            getTransactionReceipt();
            common.refreshAssetList.add(true);
          }
          gasFee = (transactionInformation.gasUsed?.toDouble() ?? 0.0) *
              web3dart
                  .getGasPrice(blockChainType: common.fromBlockChainType)
                  .getInWei
                  .toDouble() /
              1000000000000000000;
          gasFeeString = common.numberFormat(
              maxDecimal: 12, number: gasFee?.toStringAsFixed(12));
          print("transactionInformation.gas:+$gasFeeString");
        }
      }
    });
  }

  void getTransactionReceipt() async {
    GetTransactionReceipt? transactionReceipt = web3dart.lastTxId != null
        ? await api.getTransactionReceipt(txId: web3dart.lastTxId!)
        : null;
    bool haveErrorLog = false;
    if (transactionReceipt == null ||
        transactionReceipt.result == null ||
        transactionReceipt.result!.logs == null) {
      haveErrorLog = false;
    } else {
      for (Logs logs in transactionReceipt.result?.logs ?? []) {
        String topicsString = logs.topics.toString();
        if (topicsString.contains(
                "35d65561df43bc40f42f800bb9d081adb2c7e4a7d087244849eebde7e071e140") ||
            topicsString.contains(
                "c4a5a5ce167cdfe781dc65f24d03a9bfc41022b6450965c12cd0b0bdd9119236")) {
          haveErrorLog = true;
        }
      }
    }
    if ((transactionReceipt?.result?.status == "0x1" && !haveErrorLog)) {
      changeToSuccessfulStatus();
    } else {
      changeToFailureStatus(receipt: transactionReceipt);
    }
  }

  bool isAllocationType() {
    return (widget.transactionType == TransactionType.allocationProcessing ||
        widget.transactionType ==
            TransactionType.allocationTransactionHistoryProcessingStatus ||
        widget.transactionType ==
            TransactionType.allocationTransactionHistoryPendingStatus);
  }

  bool isSuccessfulType() {
    return (widget.transactionType == TransactionType.successful ||
        widget.transactionType ==
            TransactionType.allocationTransactionHistorySuccessfulStatus ||
        widget.transactionType ==
            TransactionType.sendTransactionHistorySuccessfulStatus);
  }

  changeToSuccessfulStatus() {
    if (mounted) {
      setState(() {
        if (widget.transactionType == TransactionType.pending ||
            widget.transactionType == TransactionType.allocationProcessing) {
          widget.transactionType = TransactionType.successful;
        }
        statusString =
            isAllocationType() ? "Processing" : "COMMON.SUCCESS".tr();
        getBlockNumberTimer?.cancel();
      });
    }
  }

  changeToFailureStatus({GetTransactionReceipt? receipt}) {
    if (mounted) {
      setState(() {
        widget.transactionType = TransactionType.failure;
        statusString =
            isAllocationType() ? "Processing" : "COMMON.FAILURE".tr();
        getBlockNumberTimer?.cancel();
        String? ascii2;
        if (receipt?.result?.revertReason != null) {
          ascii2 = hexToASCII2(receipt: receipt!.result!.revertReason!);
        }
        if (ascii2 == null) {
          for (Logs logs in receipt!.result?.logs ?? []) {
            String dataLogs = logs.data!;
            if (dataLogs != "0x") {
              ascii2 = hexToASCII2(receipt: dataLogs);
            }
          }
        }

        if (ascii2 != null) {
          showDialog(
              context: context,
              builder: (BuildContext context) {
                return CustomDialogBox(
                  descriptions: ascii2!,
                  buttonText: "COMMON.OK".tr(),
                );
              });
        }
      });
    }
  }

  String hexToASCII2({required String receipt}) {
    String? receiptString = receipt.replaceAll('0x', '00');
    receiptString = receiptString.replaceAll('08c379a', '0000000');
    // receiptString =
    //     receiptString.replaceAll('20000000000', '00000000000');
    // receiptString =
    //     receiptString.replaceAll(new RegExp(r'^0+(?=.)'), '');
    print("receiptString:$receiptString");
    List<String> splitted = [];
    for (int i = 0; i < receiptString.length; i = i + 2) {
      splitted.add(receiptString.substring(
          i, i + 2 > receiptString.length ? receiptString.length : i + 2));
    }
    String ascii2 = List.generate(splitted.length,
        (i) => String.fromCharCode(int.parse(splitted[i], radix: 16))).join();
    ascii2 = ascii2.replaceAll(new RegExp(r'[^A-Za-z0-9  *]'), '');
    print('ascii2:$ascii2}');
    return ascii2;
  }

  String getStatusMessage() {
    return widget.transactionType == TransactionType.pending
        ? "lt is waiting for confirmation"
        : widget.transactionType == TransactionType.allocationProcessing
            ? "lt is waiting for processing"
            : "TX_PAGE.TRANSFER_SUCCESS_DESC".tr();
  }

  String getExplorerUrl() {
    String explorerUrl;
    if (common.eurusTXType == "3") {
      explorerUrl = api.eurusExplorerUrl;
    } else if (common.eurusTXType == "2") {
      explorerUrl = api.mainNetExplorerUrl;
    } else if (widget.fromBlockChainType == BlockChainType.BinanceCoin ||
        common.fromBlockChainType == BlockChainType.BinanceCoin) {
      explorerUrl = api.bscExplorerUrl;
    } else if (common.topSelectedBlockchainType == BlockChainType.Ethereum) {
      explorerUrl = api.mainNetExplorerUrl;
    } else {
      explorerUrl = api.eurusExplorerUrl;
    }

    return explorerUrl;
  }

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    DateTime now = DateTime.now();
    final DateFormat formatter = DateFormat('dd-MM-yyyy HH:mm:ss');

    if (isEmptyString(string: formattedDate)) {
      formattedDate = formatter.format(now);
    }

    String gasFeeSymbol = getSymbolByBlockChainType(common.fromBlockChainType);

    return WillPopScope(
        onWillPop: () async => common.isTransactionHistory(
            transactionType: widget.transactionType),
        child: (isSuccessfulType() ||
                widget.transactionType == TransactionType.failure)
            ? BackGroundImage(
                child: Scaffold(
                  backgroundColor: Colors.transparent,
                  appBar: WalletAppBar(
                      backButton: common.isTransactionHistory(
                          transactionType: widget.transactionType),
                      title: "",
                      rightWidget: Icon(Icons.ios_share),
                      function: () {
                        screenshotController
                            .capture()
                            .then((Uint8List? image) async {
                          if (image != null) {
                            final directory =
                                await getApplicationDocumentsDirectory();
                            final imagePath =
                                await File('${directory.path}/image.png')
                                    .create();
                            await imagePath.writeAsBytes(image);

                            /// Share Plugin
                            String shareText =
                                "${'ASSET_TX_HISTORY.TX_ID'.tr()}${'COMMON.COLON'.tr()}\n";
                            if (!isEmptyString(string: web3dart.lastTxId))
                              shareText += web3dart.lastTxId!;
                            Share.shareFiles([imagePath.path], text: shareText);
                          }
                        }).catchError((onError) {
                          print(onError);
                        });
                      }),
                  body: SingleChildScrollView(
                    child: SafeArea(
                      child: Padding(
                        padding: EdgeInsets.only(
                            top: 8,
                            left: size.leftPadding,
                            right: size.leftPadding,
                            bottom: 30),
                        child: Screenshot(
                          controller: screenshotController,
                          child: Container(
                            decoration: new BoxDecoration(
                              color: Colors.white,
                              borderRadius:
                                  BorderRadius.all(Radius.circular(20)),
                            ),
                            child: Padding(
                              padding: EdgeInsets.only(left: 38, right: 38),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.center,
                                children: [
                                  TopBlockchainBar(
                                    transactionType: widget.transactionType,
                                    fromBlockChainType:
                                        widget.fromBlockChainType,
                                    toBlockChainType: widget.toBlockChainType,
                                  ),
                                  SizedBox(
                                    width: double.infinity,
                                    child: Image.asset(
                                      "images/line.png",
                                      package: 'euruswallet',
                                    ),
                                  ),
                                  SizedBox(height: 35),
                                  Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Image.asset(
                                        widget.transactionType ==
                                                TransactionType.failure
                                            ? "images/failure.png"
                                            : "images/${!isCentralized() ? 'decenTickIcon' : 'tickIcon'}.png",
                                        package: 'euruswallet',
                                        width: 40,
                                        height: 40,
                                        fit: BoxFit.cover,
                                      ),
                                      SizedBox(width: 12),
                                      Text(
                                        statusString,
                                        style: FXUI.normalTextStyle.copyWith(
                                          fontSize: 24,
                                          fontWeight: FontWeight.w700,
                                        ),
                                      ),
                                    ],
                                  ),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Padding(
                                          padding: EdgeInsets.only(top: 20),
                                          child: Text(
                                              common.isTransactionHistory(
                                                      transactionType: widget
                                                          .transactionType)
                                                  ? "TX_PAGE.TRANSFER_SUCCESS_DESC"
                                                      .tr()
                                                  : widget.transactionType ==
                                                          TransactionType
                                                              .pending
                                                      ? "TX_PAGE.WAITING_CONFIRM"
                                                          .tr()
                                                      : widget.transactionType ==
                                                              TransactionType
                                                                  .allocationProcessing
                                                          ? "TX_PAGE.WAITING_CONFIRM"
                                                              .tr()
                                                          : "TX_PAGE.TRANSFER_SUCCESS_DESC"
                                                              .tr(),
                                              style: FXUI.normalTextStyle
                                                  .copyWith(
                                                      fontSize: 14,
                                                      color: FXColor.textGray)),
                                        ),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Padding(
                                          padding: EdgeInsets.only(top: 20),
                                          child: Container(
                                            width: 500,
                                            decoration: new BoxDecoration(
                                              color: FXColor
                                                  .transferSuccessGreyColor,
                                              borderRadius: FXUI.cricleRadius,
                                            ),
                                            child: Column(
                                              children: [
                                                Padding(
                                                  padding: EdgeInsets.only(
                                                      top: 10, bottom: 10),
                                                  child: FittedBox(
                                                      child: Text(
                                                          common.isTransactionHistory(
                                                                  transactionType:
                                                                      widget
                                                                          .transactionType)
                                                              ? widget
                                                                  .transferAmount
                                                              : "${common.numberFormat(number: common.transferAmount)} " +
                                                                  (common.selectTokenSymbol ??
                                                                      ''),
                                                          style: FXUI
                                                              .titleTextStyle
                                                              .copyWith(
                                                                  fontSize:
                                                                      28))),
                                                ),
                                                (common.isWithdraw() &&
                                                            !common.isTransactionHistory(
                                                                transactionType:
                                                                    widget
                                                                        .transactionType)) ||
                                                        (common.isTransactionHistory(
                                                                transactionType:
                                                                    widget
                                                                        .transactionType) &&
                                                            !isEmptyString(
                                                                string: widget
                                                                    .adminFeeString))
                                                    ? Padding(
                                                        padding:
                                                            EdgeInsets.only(
                                                                bottom: 10),
                                                        child: Text(
                                                            ((common
                                                                        .adminFeeModel
                                                                        ?.data
                                                                        ?.actualFee
                                                                        .toString() ??
                                                                    widget
                                                                        .adminFeeString ??
                                                                    '') +
                                                                " " +
                                                                (common.selectTokenSymbol ??
                                                                    '') +
                                                                " ${"ADMIN_FEE.NAME".tr()}"),
                                                            style: FXUI
                                                                .titleTextStyle
                                                                .copyWith(
                                                                    fontSize:
                                                                        13)),
                                                      )
                                                    : Container(),
                                                (!common.isTransactionHistory(
                                                                transactionType:
                                                                    widget
                                                                        .transactionType) &&
                                                            isCentralized()) ||
                                                        (common.isTransactionHistory(
                                                                transactionType:
                                                                    widget
                                                                        .transactionType) &&
                                                            widget.gasFeeString ==
                                                                "")
                                                    ? Container()
                                                    : Padding(
                                                        padding:
                                                            EdgeInsets.only(
                                                                bottom: 10),
                                                        child: Text(
                                                            (common.isTransactionHistory(
                                                                        transactionType:
                                                                            widget
                                                                                .transactionType)
                                                                    ? widget
                                                                        .gasFeeString
                                                                    : (gasFeeString +
                                                                        gasFeeSymbol +
                                                                        " ${'TX_PAGE.GAS_FEE'.tr()}")) ??
                                                                '',
                                                            style: FXUI
                                                                .titleTextStyle
                                                                .copyWith(
                                                                    fontSize:
                                                                        13)),
                                                      ),
                                              ],
                                            ),
                                          )),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Padding(
                                          padding: EdgeInsets.only(top: 21),
                                          child: Row(
                                            mainAxisAlignment:
                                                MainAxisAlignment.spaceBetween,
                                            children: [
                                              Text(
                                                  '${'TX_PAGE.FROM'.tr()}${'COMMON.COLON'.tr()}',
                                                  style: FXUI.titleTextStyle
                                                      .copyWith(fontSize: 13)),
                                              // Text("KK456",
                                              //     style: FXUI.titleTextStyle
                                              //         .copyWith(fontSize: 13)),
                                            ],
                                          ),
                                        ),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Container(
                                          padding: EdgeInsets.only(
                                              top: 6, bottom: 8),
                                          alignment: Alignment(-1, 0),
                                          child: Text(
                                            common.isTransactionHistory(
                                                    transactionType:
                                                        widget.transactionType)
                                                ? widget.eurusTxType == '6'
                                                    ? 'TX_PAGE.OFFICIAL_WALLET'
                                                        .tr()
                                                    : isAllocate
                                                        ? 'TX_PAGE.MY_WALLET'
                                                            .tr(args: [
                                                            getBlockChainName(common
                                                                .fromBlockChainType)
                                                          ])
                                                        : widget.fromAddress
                                                : common.targetDepositOrWidthDrawAddresss ==
                                                        ''
                                                    ? (common.transferToMySelf &&
                                                            isCentralized())
                                                        ? 'TX_PAGE.MY_WALLET'
                                                            .tr(args: [
                                                            getBlockChainName(common
                                                                .fromBlockChainType)
                                                          ])
                                                        : fromAddress
                                                    : 'TX_PAGE.MY_WALLET'.tr(
                                                        args: [
                                                            getBlockChainName(common
                                                                .fromBlockChainType)
                                                          ]),
                                            style: FXUI.normalTextStyle
                                                .copyWith(
                                                    fontSize: 12,
                                                    color: FXColor.textGray,
                                                    fontWeight:
                                                        FontWeight.normal),
                                          ),
                                        ),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Divider(),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Padding(
                                          padding: EdgeInsets.only(top: 21),
                                          child: Row(
                                            mainAxisAlignment:
                                                MainAxisAlignment.spaceBetween,
                                            children: [
                                              Text(
                                                  '${'TX_PAGE.TRANSFER_TO'.tr()}${'COMMON.COLON'.tr()}',
                                                  style: FXUI.titleTextStyle
                                                      .copyWith(fontSize: 13)),
                                            ],
                                          ),
                                        ),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Container(
                                          padding: EdgeInsets.only(
                                              top: 6, bottom: 8),
                                          alignment: Alignment(-1, 0),
                                          child: Text(
                                            common.isTransactionHistory(
                                                    transactionType:
                                                        widget.transactionType)
                                                ? widget.eurusTxType == '6' ||
                                                        widget.eurusTxType ==
                                                            '7'
                                                    ? 'TX_PAGE.MY_WALLET'.tr(
                                                        args: [
                                                            getBlockChainName(
                                                                BlockChainType
                                                                    .Eurus)
                                                          ])
                                                    : isAllocate
                                                        ? isCentralized()
                                                            ? widget.toAddress
                                                            : 'TX_PAGE.MY_WALLET'
                                                                .tr(args: [
                                                                getBlockChainName(widget
                                                                        .toBlockChainType ??
                                                                    BlockChainType
                                                                        .Eurus)
                                                              ])
                                                        : widget.toAddress
                                                : common.targetDepositOrWidthDrawAddresss ==
                                                        ''
                                                    ? common.targetAddress ?? ""
                                                    : common
                                                        .targetDepositOrWidthDrawAddresss,
                                            style: FXUI.normalTextStyle
                                                .copyWith(
                                                    fontSize: 12,
                                                    color: FXColor.textGray),
                                          )),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Divider(),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Padding(
                                          padding: EdgeInsets.only(
                                              top: 14, bottom: 8),
                                          child: Row(
                                            mainAxisAlignment:
                                                MainAxisAlignment.spaceBetween,
                                            children: [
                                              Text(
                                                  "${'COMMON.DATE'.tr()}${'COMMON.COLON'.tr()}",
                                                  style: FXUI.titleTextStyle
                                                      .copyWith(fontSize: 13)),
                                              Text(
                                                  (common.isTransactionHistory(
                                                              transactionType:
                                                                  widget
                                                                      .transactionType)
                                                          ? widget.date
                                                          : formattedDate) ??
                                                      '',
                                                  style: FXUI.normalTextStyle
                                                      .copyWith(
                                                          fontSize: 13,
                                                          color: FXColor
                                                              .textGray)),
                                            ],
                                          ),
                                        ),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Divider(),
                                  widget.transactionType ==
                                          TransactionType.failure
                                      ? Container()
                                      : Padding(
                                          padding: EdgeInsets.only(top: 14),
                                          child: Row(
                                            crossAxisAlignment:
                                                CrossAxisAlignment.center,
                                            mainAxisAlignment:
                                                MainAxisAlignment.spaceBetween,
                                            children: [
                                              Padding(
                                                padding:
                                                    EdgeInsets.only(right: 10),
                                                child: Text(
                                                    "${'ASSET_TX_HISTORY.TX_ID'.tr()}${'COMMON.COLON'.tr()}",
                                                    style: FXUI.titleTextStyle
                                                        .copyWith(
                                                            fontSize: 13)),
                                              ),
                                              Expanded(
                                                child: Text(
                                                    web3dart.lastTxId ?? "",
                                                    style: FXUI.normalTextStyle
                                                        .copyWith(
                                                            fontSize: 13,
                                                            color: FXColor
                                                                .textGray),
                                                    maxLines: 4,
                                                    overflow:
                                                        TextOverflow.ellipsis),
                                              ),
                                              Column(
                                                crossAxisAlignment:
                                                    CrossAxisAlignment.center,
                                                children: [
                                                  InkWell(
                                                      child: Icon(Icons.link,
                                                          size: 30,
                                                          color: common
                                                              .getBackGroundColor()),
                                                      onTap: () {
                                                        common.pushPage(
                                                            page: WebViewPage(
                                                                link: getExplorerUrl() +
                                                                    (web3dart
                                                                            .lastTxId ??
                                                                        "empty lastTxId"),
                                                                appTitle:
                                                                    "Transcation Record Page"),
                                                            context: context);
                                                      }),
                                                  Container(
                                                      child: new Material(
                                                          child: InkWell(
                                                            onTap: () {
                                                              Clipboard.setData(
                                                                  ClipboardData(
                                                                      text: web3dart
                                                                              .lastTxId ??
                                                                          "empty lastTxId"));
                                                              common
                                                                  .showCopiedToClipboardSnackBar(
                                                                context,
                                                              );
                                                            },
                                                            child: Padding(
                                                              padding:
                                                                  const EdgeInsets
                                                                          .only(
                                                                      left: 5,
                                                                      top: 5,
                                                                      bottom: 5,
                                                                      right: 5),
                                                              child: SizedBox(
                                                                width: 25,
                                                                height: 25,
                                                                child:
                                                                    Image.asset(
                                                                  "images/paste.png",
                                                                  package:
                                                                      'euruswallet',
                                                                  color: common
                                                                      .getBackGroundColor(),
                                                                  fit: BoxFit
                                                                      .fill,
                                                                ),
                                                              ),
                                                            ),
                                                          ),
                                                          color: Colors.white)),
                                                ],
                                              )
                                            ],
                                          ),
                                        ),
                                  Padding(
                                    padding:
                                        EdgeInsets.only(top: 30, bottom: 27),
                                    child: SubmitButton(
                                      btnController: btnController,
                                      label: "COMMON.CLOSE".tr(),
                                      onPressed: () {
                                        widget.transactionType ==
                                                    TransactionType
                                                        .successful ||
                                                widget.transactionType ==
                                                    TransactionType.failure
                                            ? Navigator.popUntil(context,
                                                ModalRoute.withName("HomePage"))
                                            : Navigator.pop(context);
                                        btnController.reset();
                                      },
                                    ),
                                  )
                                ],
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                ),
              )
            : ConfirmationPage(
                transactionType: widget.transactionType,
                confirmBlockNumber: confirmBlockNumber,
                requestSuccessBlock: requestSuccessBlock,
                statusString: statusString));
  }
}
