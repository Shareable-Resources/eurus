import 'package:apihandler/apiHandler.dart';
import 'package:decimal/decimal.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/acknowledgementDialog.dart';
import 'package:euruswallet/commonUI/customDialogBox.dart';
import 'package:euruswallet/commonUI/topBlockchainBar.dart';
import 'package:euruswallet/model/ethgasstation.dart';
import 'package:euruswallet/pages/centralized/top_up_payment_wallet_page.dart';
import 'package:euruswallet/pages/confirmationPage.dart';
import 'package:material_segmented_control/material_segmented_control.dart';
import 'package:pin_code_fields/pin_code_fields.dart';

class TransferPage extends StatefulWidget {
  final String titleName;
  final bool transferToMySelf;
  final BlockChainType fromBlockChainType;
  final String? ethereumErc20ContractAddress;
  final String? eurusErc20ContractAddress;

  TransferPage({
    Key? key,
    required this.titleName,
    this.transferToMySelf = false,
    this.ethereumErc20ContractAddress,
    this.eurusErc20ContractAddress,
    required this.fromBlockChainType,
  }) : super(key: key);

  @override
  _TransferPageState createState() => _TransferPageState();
}

class _TransferPageState extends State<TransferPage> {
  //String amount = "0";
  //bool showTransactionCode = false;
  // String paymentPassword;
  bool alreadyClickNextBtn = false;
  bool showGasPriceRow = false;
  StreamController<ErrorAnimationType> errorController =
      StreamController<ErrorAnimationType>();
  // TextEditingController paymentPasswordController = TextEditingController();
  // FocusNode myFocusNode = FocusNode();
  bool hasError = false;
  String? balance;
  double? balanceAmount;
  String? amountErrorMessage;
  String? gasPriceErrorMessage;
  TextEditingController amountController = TextEditingController();
  TextEditingController gasPriceController = TextEditingController();
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    super.initState();
    common.fromBlockChainType = widget.fromBlockChainType;
    common.currentSelectionSpeed = 0;
    common.currentGas = 0;
    initGasPrice();
    web3dart.transactionSpeed = common.currentSelectionSpeed.toDouble() + 1;

    common.setUpErc20TokenContract(
        ethereumErc20ContractAddress: widget.ethereumErc20ContractAddress,
        eurusErc20ContractAddress: widget.eurusErc20ContractAddress,
        fromBlockChainType: widget.fromBlockChainType);
    refreshBalance();
  }

  @override
  void dispose() {
    super.dispose();
  }

  Future<Ethgasstation?> getGasStationData() async {
    if (common.fromBlockChainType == BlockChainType.Ethereum) {
      var gasStationDataResponse = await apiHandler
          .get("https://ethgasstation.info/api/ethgasAPI.json?");
      final gasStationData = Ethgasstation.fromJson(gasStationDataResponse);
      common.gasStationData = gasStationData;
      web3dart.ethereumGasPrice = gasStationData.fast.toDouble() * 100000000;
      return gasStationData;
    }
    return common.gasStationData;
  }

  Future<bool> initGasPrice() async {
    gasPriceErrorMessage = null;
    gasPriceController.text = "0";
    if (common.fromBlockChainType == BlockChainType.Ethereum) {
      if (common.gasStationData != null) {
        web3dart.ethereumGasPrice =
            common.gasStationData!.fast.toDouble() * 100000000;
        gasPriceController.text =
            (common.gasStationData!.fast.toDouble() / 10).toString();
      }
    }
    setState(() {});
    return true;
  }

  Future<void> refreshBalance() async {
    await getGasStationData();
    final selectTokenSymbol = common.selectTokenSymbol ?? '';
    if (widget.transferToMySelf) {
      await api.getAdminFee(selectTokenSymbol: selectTokenSymbol);
    }
    common.transferToMySelf = widget.transferToMySelf;
    if (common.transferToMySelf) {
      common.targetType = TargetType.ethAddress;
      if (!isCentralized()) {
        common.targetAddress =
            common.fromBlockChainType == BlockChainType.Ethereum
                ? await web3dart.getEurusUserDepositAddress()
                : web3dart.myEthereumAddress.toString();
      }
    }
    await web3dart.getErc20Balance(
        type: common.fromBlockChainType, isEthOrEun: common.isEthOrEun());

    if (common.isEthOrEun()) {
      balanceAmount = double.parse(
          common.fromBlockChainType == BlockChainType.Ethereum
              ? web3dart.ethBalanceFromEthereum ?? "0"
              : web3dart.ethBalanceFromEurus ?? "0");
      balance = "TX_PAGE.ASSET_BALANCE".tr() +
          common.numberFormat(
              maxDecimal: 8,
              number: common.fromBlockChainType == BlockChainType.Ethereum
                  ? web3dart.ethBalanceFromEthereum
                  : web3dart.ethBalanceFromEurus) +
          " " +
          selectTokenSymbol;
    } else {
      balanceAmount = double.parse(
          common.fromBlockChainType == BlockChainType.Ethereum
              ? web3dart.erc20TokenBalanceFromEthereum ?? "0"
              : web3dart.erc20TokenBalanceFromEurus ?? "0");
      balance = "TX_PAGE.ASSET_BALANCE".tr() +
          common.numberFormat(
              maxDecimal: 8,
              number: common.fromBlockChainType == BlockChainType.Ethereum
                  ? web3dart.erc20TokenBalanceFromEthereum
                  : web3dart.erc20TokenBalanceFromEurus) +
          " " +
          selectTokenSymbol;
    }

    await _updateEstimateGas(amountInDecimal: Decimal.zero);
  }

  double? getCurrentGasFee({bool enterGasPrice: false}) {
    final estimateGasFee = web3dart.estimateMaxGas;
    if (estimateGasFee == null) return null;

    if (enterGasPrice) {
      return estimateGasFee * web3dart.ethereumGasPrice / pow(10, 18);
    } else {
      return estimateGasFee *
          web3dart
              .getGasPrice(blockChainType: common.fromBlockChainType)
              .getInWei
              .toDouble() /
          pow(10, 18);
    }
  }

  Widget speedButton({
    required String imagePath,
    required String buttonText,
    required int currentSelection,
  }) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Image.asset(imagePath,
            package: 'euruswallet',
            width: 16,
            height: 16,
            fit: BoxFit.contain,
            color: currentSelection == common.currentSelectionSpeed
                ? Colors.white
                : isEmptyString(string: buttonText)
                    ? FXColor.lightGray
                    : FXColor.textGray),
        SizedBox(width: 4),
        Text(buttonText,
            style: FXUI.normalTextStyle.copyWith(
                color: currentSelection == common.currentSelectionSpeed
                    ? Colors.white
                    : isEmptyString(string: buttonText)
                        ? FXColor.lightGray
                        : FXColor.textGray,
                fontSize: 13)),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    return BackGroundImage(
      child: Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: Colors.transparent,
        appBar: WalletAppBar(title: widget.titleName),
        body: Column(
          children: [
            SingleChildScrollView(
              child: Padding(
                padding: EdgeInsets.only(
                    top: 8, left: size.leftPadding, right: size.leftPadding),
                child: Container(
                  decoration: new BoxDecoration(
                    borderRadius: BorderRadius.all(Radius.circular(20)),
                    color: Colors.white,
                  ),
                  child: Padding(
                    padding: EdgeInsets.only(
                        left: size.leftPadding, right: size.leftPadding),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        TopBlockchainBar(
                          transactionType: TransactionType.transferInput,
                          fromBlockChainType: widget.fromBlockChainType,
                        ),
                        SizedBox(
                          width: double.infinity,
                          child: Image.asset(
                            "images/line.png",
                            package: 'euruswallet',
                          ),
                        ),
                        Row(
                            crossAxisAlignment: CrossAxisAlignment.center,
                            children: [
                              Expanded(
                                  child: TextFormField(
                                      maxLines: 1,
                                      keyboardType:
                                          TextInputType.numberWithOptions(
                                              signed: true, decimal: true),
                                      style: FXUI.inputStyle.copyWith(
                                          color: FXColor.textGray,
                                          fontSize: 48),
                                      inputFormatters: [
                                        FilteringTextInputFormatter.allow(
                                            RegExp('[0-9.]+')),
                                      ],
                                      autofocus: true,
                                      onChanged: (string) async {
                                        final amount = Decimal.tryParse(string);
                                        if (amount == null) return;
                                        _updateEstimateGas(
                                            amountInDecimal: amount);
                                      },
                                      controller: amountController,
                                      decoration: InputDecoration(
                                        fillColor: Colors.transparent,
                                        filled: true,
                                        enabledBorder: OutlineInputBorder(
                                            borderRadius: FXUI.cricleRadius,
                                            borderSide: BorderSide(
                                                width: 0,
                                                color: Colors.transparent)),
                                        focusedBorder: OutlineInputBorder(
                                            borderRadius: FXUI.cricleRadius,
                                            borderSide: BorderSide(
                                                width: 0,
                                                color: Colors.transparent)),
                                        border: OutlineInputBorder(
                                            borderRadius: FXUI.cricleRadius,
                                            borderSide: BorderSide(
                                                width: 0,
                                                color: Colors.transparent)),
                                      ))),
                              Text(common.selectTokenSymbol ?? "USDT",
                                  style: FXUI.titleTextStyle.copyWith(
                                      color: FXColor.textGray.withOpacity(0.5),
                                      fontSize: 18)),
                            ]),
                        Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            crossAxisAlignment: CrossAxisAlignment.center,
                            children: [
                              Expanded(
                                  child: Text(balance ?? "",
                                      style: FXUI.titleTextStyle.copyWith(
                                          color:
                                              FXColor.textGray.withOpacity(0.5),
                                          fontSize: 15),
                                      maxLines: 4,
                                      overflow: TextOverflow.ellipsis)),
                              Container(
                                width: 47,
                                height: 24,
                                child: ElevatedButton(
                                    style: ElevatedButton.styleFrom(
                                        shape: RoundedRectangleBorder(
                                          borderRadius: FXUI.cricleRadius,
                                        ),
                                        padding: EdgeInsets.zero,
                                        onSurface: common.getBackGroundColor(),
                                        shadowColor: Colors.black,
                                        primary: common.getBackGroundColor(),
                                        elevation: isCentralized() ? 1 : 10),
                                    child: Text('TX_PAGE.MAX_AMOUNT'.tr()),
                                    onPressed: () async {
                                      Decimal totalAmount = Decimal.parse(
                                          (balanceAmount ?? 0.0).toString());

                                      if (common.isWithdraw()) {
                                        totalAmount -= Decimal.parse((common
                                                    .adminFeeModel
                                                    ?.data
                                                    ?.actualFee ??
                                                0.0)
                                            .toString());
                                      }

                                      await _updateEstimateGas(
                                          amountInDecimal: totalAmount);

                                      if ((common.isEthOrEun()) &&
                                          common.currentGas != null) {
                                        totalAmount -= Decimal.parse(
                                            common.currentGas!.toString());
                                      }

                                      if (totalAmount < Decimal.zero) {
                                        totalAmount = Decimal.zero;
                                        common.showSnackBar(
                                          errorMessage: isEmptyString(
                                                  string: balance)
                                              ? "TX_PAGE.ERROR.NOT_GET_BALANCE"
                                                  .tr()
                                              : "TX_PAGE.ERROR.INSUFFICIENT_ASSET"
                                                  .tr(),
                                          context: context,
                                        );
                                      }
                                      amountController.text =
                                          totalAmount.toDouble().toString();
                                      btnController.reset();
                                    }),
                              ),
                            ]),
                        Container(height: 24),
                        widget.transferToMySelf &&
                                common.fromBlockChainType ==
                                    BlockChainType.Eurus
                            ? Padding(
                                padding: EdgeInsets.only(bottom: 20),
                                child: Row(
                                  mainAxisAlignment:
                                      MainAxisAlignment.spaceBetween,
                                  children: [
                                    Row(
                                      children: [
                                        Text("ADMIN_FEE.NAME".tr(),
                                            style: FXUI.titleTextStyle.copyWith(
                                                color: FXColor.lightGray
                                                    .withOpacity(0.5),
                                                fontSize: 14)),
                                        Padding(
                                            padding: EdgeInsets.only(left: 5)),
                                        InkWell(
                                            onTap: () async {
                                              await showDialog(
                                                  context: context,
                                                  builder:
                                                      (BuildContext context) {
                                                    return CustomDialogBox(
                                                      title:
                                                          "ADMIN_FEE.NAME".tr(),
                                                      titleIcon: Icon(
                                                          Icons.info,
                                                          color: common
                                                              .getBackGroundColor()),
                                                      descriptions:
                                                          "ADMIN_FEE.DESC".tr(),
                                                      buttonText:
                                                          "COMMON.CLOSE".tr(),
                                                    );
                                                  });
                                            },
                                            child: Padding(
                                                padding: EdgeInsets.all(3),
                                                child: Image.asset(
                                                  "images/gasAlert.png",
                                                  package: 'euruswallet',
                                                  width: 13,
                                                  height: 13,
                                                  fit: BoxFit.cover,
                                                  color: FXColor.lightGray,
                                                ))),
                                      ],
                                    ),
                                    Row(
                                      children: [
                                        Text(
                                            common.numberFormat(
                                                maxDecimal: 12,
                                                number: common.adminFeeModel
                                                        ?.data?.actualFee
                                                        .toString() ??
                                                    ""),
                                            style: FXUI.titleTextStyle.copyWith(
                                                color: FXColor.lightGray,
                                                fontSize: 14)),
                                        Padding(
                                          padding: EdgeInsets.only(left: 5),
                                          child: Text(
                                              common.selectTokenSymbol ?? '',
                                              style: FXUI.titleTextStyle
                                                  .copyWith(
                                                      color: FXColor.lightGray
                                                          .withOpacity(0.5),
                                                      fontSize: 14)),
                                        ),
                                      ],
                                    ),
                                  ],
                                ),
                              )
                            : Container(),
                        Padding(
                          padding: EdgeInsets.only(bottom: 20),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Row(children: [
                                Text("TX_PAGE.GAS_FEE".tr(),
                                    style: FXUI.titleTextStyle.copyWith(
                                        color:
                                            FXColor.lightGray.withOpacity(0.5),
                                        fontSize: 14)),
                                Padding(padding: EdgeInsets.only(left: 5)),
                                InkWell(
                                    onTap: () async {
                                      await showDialog(
                                          context: context,
                                          builder: (BuildContext context) {
                                            return CustomDialogBox(
                                              title: "GAS_FEE.NAME".tr(),
                                              titleIcon: Icon(Icons.info,
                                                  color: common
                                                      .getBackGroundColor()),
                                              descriptions: "GAS_FEE.DESC"
                                                  .tr(args: [
                                                getSymbolByBlockChainType(
                                                    common.fromBlockChainType)
                                              ]),
                                              buttonText: "COMMON.CLOSE".tr(),
                                            );
                                          });
                                    },
                                    child: Padding(
                                        padding: EdgeInsets.all(3),
                                        child: Image.asset(
                                          "images/gasAlert.png",
                                          package: 'euruswallet',
                                          width: 13,
                                          height: 13,
                                          fit: BoxFit.cover,
                                          color: FXColor.lightGray,
                                        )))
                              ]),
                              Row(
                                children: [
                                  Text(
                                      common.numberFormat(
                                          maxDecimal: 12,
                                          number: common.currentGas
                                                  ?.toStringAsFixed(12) ??
                                              ""),
                                      style: FXUI.titleTextStyle.copyWith(
                                          color: FXColor.lightGray,
                                          fontSize: 14)),
                                  Padding(
                                    padding: EdgeInsets.only(left: 5),
                                    child: Text(
                                        getSymbolByBlockChainType(
                                            common.fromBlockChainType),
                                        style: FXUI.titleTextStyle.copyWith(
                                            color: FXColor.lightGray
                                                .withOpacity(0.5),
                                            fontSize: 14)),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                        common.fromBlockChainType == BlockChainType.Ethereum
                            ? Padding(
                                padding: EdgeInsets.only(bottom: 20),
                                child: Container(
                                  width: size.blockSizeVertical * 100,
                                  child: Column(
                                    children: [
                                      MaterialSegmentedControl(
                                        horizontalPadding: EdgeInsets.zero,
                                        children: {
                                          0: Container(
                                              child: speedButton(
                                                  buttonText:
                                                      'TX_PAGE.GAS_FEE_SPEED.SLOW'
                                                          .tr(),
                                                  imagePath:
                                                      "images/slowSpeed.png",
                                                  currentSelection: 0),
                                              width: 100),
                                          1: Container(
                                              child: speedButton(
                                                  buttonText:
                                                      'TX_PAGE.GAS_FEE_SPEED.MEDIUM'
                                                          .tr(),
                                                  imagePath:
                                                      "images/mediumSpeed.png",
                                                  currentSelection: 1),
                                              width: 200),
                                          2: Container(
                                              child: speedButton(
                                                  buttonText:
                                                      'TX_PAGE.GAS_FEE_SPEED.RAPID'
                                                          .tr(),
                                                  imagePath:
                                                      "images/rapidSpeed.png",
                                                  currentSelection: 2),
                                              width: 400),
                                          3: Container(
                                              child: speedButton(
                                                  buttonText: '',
                                                  imagePath:
                                                      "images/moreSpeed.png",
                                                  currentSelection: 3),
                                              width: 20),
                                        },
                                        selectionIndex:
                                            common.currentSelectionSpeed,
                                        borderColor: FXColor.lightGray,
                                        selectedColor:
                                            common.getBackGroundColor(),
                                        unselectedColor: Colors.white,
                                        borderRadius: 10.0,
                                        onSegmentChosen: (int index) async {
                                          common.currentSelectionSpeed = index;
                                          if (common.currentSelectionSpeed
                                                  .toDouble() ==
                                              3) {
                                            showGasPriceRow = true;
                                            await initGasPrice();
                                          } else {
                                            showGasPriceRow = false;
                                            await initGasPrice();
                                            web3dart.transactionSpeed = (common
                                                    .currentSelectionSpeed
                                                    .toDouble() +
                                                1);
                                          }
                                          setState(() {
                                            common.currentGas =
                                                getCurrentGasFee();
                                          });
                                        },
                                      ),
                                      showGasPriceRow
                                          ? Padding(
                                              padding: EdgeInsets.only(top: 19),
                                              child: Row(
                                                crossAxisAlignment:
                                                    CrossAxisAlignment.center,
                                                children: [
                                                  Text("TX_PAGE.GAS_PRICE".tr(),
                                                      style: FXUI.titleTextStyle
                                                          .copyWith(
                                                              color: FXColor
                                                                  .lightGray
                                                                  .withOpacity(
                                                                      0.5),
                                                              fontSize: 14)),
                                                  InkWell(
                                                    onTap: () async {
                                                      await showDialog(
                                                          context: context,
                                                          builder: (BuildContext
                                                              context) {
                                                            return CustomDialogBox(
                                                              title:
                                                                  "TX_PAGE.GAS_PRICE"
                                                                      .tr(),
                                                              titleIcon: Icon(
                                                                  Icons.info,
                                                                  color: common
                                                                      .getBackGroundColor()),
                                                              descriptions:
                                                                  "TX_PAGE.GAS_PRICE_DESC"
                                                                      .tr(),
                                                              buttonText:
                                                                  "COMMON.CLOSE"
                                                                      .tr(),
                                                            );
                                                          });
                                                    },
                                                    child: Padding(
                                                        padding:
                                                            EdgeInsets.only(
                                                                left: 5),
                                                        child: Image.asset(
                                                          "images/gasAlert.png",
                                                          package:
                                                              'euruswallet',
                                                          width: 13,
                                                          height: 13,
                                                          fit: BoxFit.cover,
                                                          color:
                                                              FXColor.lightGray,
                                                        )),
                                                  ),
                                                  Padding(
                                                      padding: EdgeInsets.only(
                                                          left: 18),
                                                      child: Container(
                                                          height: 30,
                                                          width: 100,
                                                          child: TextFormField(
                                                            onChanged:
                                                                (String value) {
                                                              final gasPrice =
                                                                  double
                                                                      .tryParse(
                                                                          value);
                                                              if (gasPrice !=
                                                                      null &&
                                                                  gasPrice >
                                                                      0) {
                                                                web3dart.ethereumGasPrice =
                                                                    gasPrice *
                                                                        1000000000;

                                                                setState(() {
                                                                  common.currentGas =
                                                                      getCurrentGasFee(
                                                                          enterGasPrice:
                                                                              true);
                                                                });
                                                              }
                                                            },
                                                            keyboardType:
                                                                TextInputType
                                                                    .numberWithOptions(
                                                                        decimal:
                                                                            true),
                                                            style: FXUI
                                                                .inputStyle
                                                                .copyWith(
                                                                    color: FXColor
                                                                        .textGray),
                                                            autofocus: false,
                                                            controller:
                                                                gasPriceController,
                                                            decoration:
                                                                InputDecoration(
                                                              hintText: "",
                                                              contentPadding:
                                                                  EdgeInsets.symmetric(
                                                                      vertical:
                                                                          0,
                                                                      horizontal:
                                                                          15.5),
                                                              // fillColor: FXColor
                                                              //     .mainDeepBlueColor,
                                                              filled: true,
                                                              enabledBorder: OutlineInputBorder(
                                                                  borderRadius:
                                                                      BorderRadius
                                                                          .circular(
                                                                              8.0),
                                                                  borderSide:
                                                                      BorderSide(
                                                                          width:
                                                                              0,
                                                                          color:
                                                                              Colors.transparent)),
                                                              focusedBorder: OutlineInputBorder(
                                                                  borderRadius:
                                                                      BorderRadius
                                                                          .circular(
                                                                              8.0),
                                                                  borderSide:
                                                                      BorderSide(
                                                                          width:
                                                                              0,
                                                                          color:
                                                                              Colors.transparent)),
                                                              border: OutlineInputBorder(
                                                                  borderRadius:
                                                                      BorderRadius
                                                                          .circular(
                                                                              8.0),
                                                                  borderSide:
                                                                      BorderSide(
                                                                          width:
                                                                              0,
                                                                          color:
                                                                              Colors.transparent)),
                                                            ),
                                                          ))),
                                                  Padding(
                                                    padding: EdgeInsets.only(
                                                        left: 8),
                                                    child: Text("GWEI",
                                                        style: FXUI
                                                            .titleTextStyle
                                                            .copyWith(
                                                                color: FXColor
                                                                    .lightGray
                                                                    .withOpacity(
                                                                        0.5),
                                                                fontSize: 14)),
                                                  ),
                                                ],
                                              ),
                                            )
                                          : Container(),
                                      alreadyClickNextBtn &&
                                              gasPriceErrorMessage != null
                                          ? Padding(
                                              padding: EdgeInsets.only(top: 10),
                                              child: Row(
                                                children: [
                                                  Text(
                                                    gasPriceErrorMessage!,
                                                    style: FXUI.normalTextStyle
                                                        .copyWith(
                                                            color: FXColor
                                                                .alertRedColor,
                                                            fontSize: 12),
                                                  ),
                                                ],
                                              ),
                                            )
                                          : Container(height: 0)
                                    ],
                                  ),
                                ))
                            : Container(),
                        Padding(
                          padding: EdgeInsets.only(
                              top: 30, left: 25, right: 25, bottom: 30),
                          child: SubmitButton(
                            loadingSecond: 4,
                            btnController: btnController,
                            label: 'COMMON.NEXT_STEP'.tr(),
                            onPressed: () async {
                              double gasPrice = double.parse(
                                  isEmptyString(string: gasPriceController.text)
                                      ? "0"
                                      : gasPriceController.text);
                              if (showGasPriceRow) {
                                if (gasPrice > 0 &&
                                    gasPriceController.text.length > 0) {
                                  // if(gasPrice < 100){
                                  //   gasPriceErrorMessage = "TX_PAGE.ERROR.GAS_MUST_GREATER_100".tr();
                                  // } else {
                                  gasPriceErrorMessage = null;
                                  // }
                                } else {
                                  gasPriceErrorMessage =
                                      "TX_PAGE.ERROR.EMPTY_GAS_PRICE".tr();
                                }
                              } else {
                                gasPriceErrorMessage = null;
                              }

                              if (isCentralized() &&
                                  common.currentGas != null) {
                                final currentGasBalance =
                                    await web3dart.eurusEthClient.getBalance(
                                  EthereumAddress.fromHex(
                                      common.ownerWalletAddress ?? ''),
                                );
                                if (currentGasBalance.getInWei.toDouble() /
                                        pow(10, 18) <
                                    (common.currentGas ?? 0.0)) {
                                  await common.showGasLeakDialog(
                                      context: context);
                                  btnController.reset();
                                  return;
                                }
                              }

                              String amountString =
                                  amountController.text.replaceAll(',', '');
                              amountString =
                                  (isEmptyString(string: amountString))
                                      ? "0"
                                      : amountString;
                              double totalAmount = double.parse(amountString);
                              // if(common.isEthOrEun()){
                              //   totalAmount += common.currentGas * 10;
                              // }

                              if (common.isWithdraw()) {
                                totalAmount +=
                                    common.adminFeeModel?.data?.actualFee ??
                                        0.0;
                              }
                              String gasBalance = await web3dart.getETHBalance(
                                  blockChainType: common.fromBlockChainType);
                              if (double.parse(amountString) > 0 &&
                                  amountString.length > 0) {
                                if (totalAmount > (balanceAmount ?? 0.0)) {
                                  amountErrorMessage =
                                      "TX_PAGE.ERROR.INSUFFICIENT_ASSET".tr();
                                } else if (double.parse(gasBalance) <= 0) {
                                  amountErrorMessage =
                                      "TX_PAGE.ERROR.INSUFFICIENT_GAS_FEE".tr();
                                } else {
                                  amountErrorMessage = null;
                                  common.transferAmount = amountString;
                                }
                              } else {
                                amountErrorMessage = widget.transferToMySelf
                                    ? "TX_PAGE.ERROR.EMPTY_ALLOCATE_AMOUNT".tr()
                                    : "TX_PAGE.ERROR.EMPTY_TRANSFER_AMOUNT"
                                        .tr();
                              }

                              if (isEmptyString(string: balance)) {
                                amountErrorMessage =
                                    "TX_PAGE.ERROR.NOT_GET_BALANCE".tr();
                              }
                              if (amountErrorMessage != null) {
                                common.showSnackBar(
                                    errorMessage: amountErrorMessage!,
                                    context: context);
                              }
                              if (isEmptyString(string: gasPriceErrorMessage) &&
                                  isEmptyString(string: amountErrorMessage)) {
                                if (widget.transferToMySelf) {
                                  String statmentKey =
                                      common.fromBlockChainType ==
                                              BlockChainType.Eurus
                                          ? 'ackWithdrawStatement'
                                          : 'ackDepositStatement';

                                  String? doNotShowStatement =
                                      await NormalStorageKit()
                                          .readValue(statmentKey);

                                  if (common.getUriVal(
                                          doNotShowStatement,
                                          web3dart.myEthereumAddress?.hex ??
                                              '') !=
                                      '1') {
                                    bool? ackResponse = await showDialog(
                                      barrierDismissible: false,
                                      context: context,
                                      builder: (_) => AcknowledgementDialog(
                                        statement:
                                            '${common.fromBlockChainType == BlockChainType.Eurus ? 'ACKNOWLEDGEMENT_DIALOG.WITHDRAW_STATEMENT'.tr() : 'ACKNOWLEDGEMENT_DIALOG.DEPOSIT_STATEMENT'.tr()}',
                                        mainIcon: Image.asset(
                                          'images/${isCentralized() ? 'icn_cross_chain_warning_centralized' : 'icn_cross_chain_warning_decentralized'}.png',
                                          package: 'euruswallet',
                                          width: MediaQuery.of(context)
                                                  .size
                                                  .width /
                                              4,
                                        ),
                                        buttonText:
                                            'ACKNOWLEDGEMENT_DIALOG.I_AGREE'
                                                .tr(),
                                        buttonHandler: () {
                                          common.pushPage(
                                              page: ConfirmationPage(),
                                              context: context);
                                        },
                                        dontAskAgainText:
                                            'ACKNOWLEDGEMENT_DIALOG.DONT_ASK_AGAIN'
                                                .tr(),
                                      ),
                                    );

                                    if (ackResponse != null &&
                                        ackResponse == true) {
                                      String updatedAckState =
                                          common.updateUriVal(
                                              doNotShowStatement,
                                              web3dart.myEthereumAddress?.hex ??
                                                  '',
                                              '1');
                                      await NormalStorageKit().setValue(
                                          updatedAckState, statmentKey);
                                    }
                                  } else {
                                    common.pushPage(
                                        page: ConfirmationPage(),
                                        context: context);
                                  }
                                } else {
                                  common.pushPage(
                                      page: ConfirmationPage(),
                                      context: context);
                                }
                              }

                              setState(() {
                                alreadyClickNextBtn = true;
                              });
                              btnController.reset();
                            },
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future _updateEstimateGas({required Decimal amountInDecimal}) async {
    if (amountInDecimal < Decimal.zero) return;

    String? to = common.targetAddress ?? '';
    final amount = amountInDecimal.toDouble();
    final enterAmountWithFee = (amountInDecimal +
            Decimal.parse(
                (common.adminFeeModel?.data?.actualFee ?? 0.0).toString()))
        .toDouble();
    Transaction? transaction;
    BlockChainType blockChainType;
    if (isCentralized()) {
      if (common.isCenWithdraw) {
        blockChainType = BlockChainType.Eurus;
        transaction = await web3dart.cenSubmitWithdraw(
          userWalletAddress: common.cenUserWalletAddress ?? '',
          deployedContract: web3dart.erc20ContractFromEurus,
          enterAmount: amount,
          toAddress: to,
          enterAmountWithFee: enterAmountWithFee,
          selectTokenSymbol: common.selectTokenSymbol ?? '',
        );
      } else {
        blockChainType = BlockChainType.Eurus;
        transaction = await web3dart.requestTransfer(
          selectTokenSymbol: common.selectTokenSymbol ?? '',
          userWalletAddress: common.cenUserWalletAddress ?? '',
          deployedContract: web3dart.erc20ContractFromEurus,
          enterAmount: amount,
          toAddress: to,
          blockChainType: BlockChainType.Eurus,
        );
      }
    } else {
      if (!common.transferToMySelf || common.isDeposit()) {
        if (common.isDeposit()) {
          to = await web3dart.getEurusUserDepositAddress() ?? '';
        }
        blockChainType = common.fromBlockChainType;
        if (common.isEthOrEun() || common.isBSC()) {
          transaction = await web3dart.sendETH(
              enterAmount: amount,
              toAddress: to,
              type: common.fromBlockChainType);
        } else {
          transaction = await web3dart.sendERC20(
              deployedContract:
                  web3dart.getERC20Contract(common.fromBlockChainType),
              enterAmount: amount,
              toAddress: to,
              blockChainType: common.fromBlockChainType);
        }
      } else {
        blockChainType = BlockChainType.Eurus;
        transaction = await web3dart.submitWithdrawERC20(
          deployedContract: web3dart.erc20ContractFromEurus,
          enterAmount: amount,
          toAddress: to,
          enterAmountWithFee: enterAmountWithFee,
        );
      }
    }

    if (transaction == null) return;

    final maxGas = await web3dart.estimateGas(
      blockChainType: blockChainType,
      transaction: transaction,
    );

    setState(() {
      web3dart.estimateMaxGas = maxGas.toInt();
      common.currentGas = getCurrentGasFee();
    });
  }
}
