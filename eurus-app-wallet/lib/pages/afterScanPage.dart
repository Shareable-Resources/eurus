import 'package:app_authentication_kit/utils/address.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/common/qrcode_scanner.dart';
import 'package:euruswallet/pages/settingSubpages/cardContainer.dart';
import 'package:flutter_neumorphic/flutter_neumorphic.dart';
import 'package:euruswallet/commonUI/wallet_app_bar.dart';

class AfterScanPage extends StatefulWidget {
  const AfterScanPage({
    Key? key,
    required this.qrCodeString,
  });
  final String qrCodeString;

  _AfterScanPageState createState() => _AfterScanPageState();
}

class _AfterScanPageState extends State<AfterScanPage> {
  Color get themeColor => common.getBackGroundColor();
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();
  RoundedLoadingButtonController btnController2 =
      RoundedLoadingButtonController();
  TextEditingController textEditingController = TextEditingController();
  FocusNode textFieldFocusNode = FocusNode();
  Color lightGray = Color(0xFFF4F3F7);
  String? errorText;

  @override
  void initState() {
    super.initState();
    textEditingController.text = widget.qrCodeString
        .replaceAll('eurus:', '')
        .replaceAll('ethereum:', '')
        .replaceAll('binance coin:', '')
        .replaceAll('empty', '');
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: FXColor.veryLightGreyTextColor,
      appBar: WalletAppBarV2(
        title: 'AFTER_SCAN_PAGE.SCAN_RESULT'.tr(),
        backButton: true,
      ),
      body: Container(
        child: SingleChildScrollView(
          child: SafeArea(
            child: CardContainer(
              '',
              Container(
                  child: Column(
                children: [
                  Padding(
                    padding: EdgeInsets.only(left: 13, right: 13, top: 50),
                    child: TextFormField(
                      readOnly: true,
                      style: FXUI.inputStyle,
                      focusNode: textFieldFocusNode,
                      autofocus: false,
                      controller: textEditingController,
                      validator: (value) {
                        return errorText;
                      },
                      decoration: FXUI.inputDecoration.copyWith(
                          hintText:
                              "SEND_PAGE.RECEIVING_ADDRESS.PLACEHOLDER".tr(),
                          prefix: Padding(
                            padding: EdgeInsets.only(right: 9.6),
                            child: Image.asset(
                              'images/icn_address.png',
                              package: 'euruswallet',
                              width: 14,
                            ),
                          )),
                      onChanged: (value) => errorText = null,
                    ),
                  ),
                  Padding(
                    padding: EdgeInsets.only(top: 22, bottom: 20),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Expanded(
                            child: Divider(
                          color: lightGray,
                          thickness: 1,
                        )),
                        TextButton(
                          onPressed: () {},
                          style: TextButton.styleFrom(
                            padding: EdgeInsets.zero,
                            minimumSize: Size(100, 30),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(15),
                              side: BorderSide(color: lightGray),
                            ),
                            textStyle: FXUI.normalTextStyle,
                            tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                          ).copyWith(
                            foregroundColor:
                                MaterialStateProperty.all(lightGray),
                            backgroundColor:
                                MaterialStateProperty.all(Colors.white),
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            crossAxisAlignment: CrossAxisAlignment.center,
                            children: [
                              Text('MAIN_FNC.SEND'.tr(),
                                  style: FXUI.normalTextStyle
                                      .copyWith(color: FXColor.mainBlueColor)),
                            ],
                          ),
                        ),
                        Expanded(
                            child: Divider(
                          color: lightGray,
                          thickness: 1,
                        )),
                      ],
                    ),
                  ),
                  Padding(
                    padding: EdgeInsets.only(left: 13, right: 13),
                    child: Column(
                      children: [
                        Container(
                            child: SubmitButton(
                                btnController: btnController,
                                label: 'AFTER_SCAN_PAGE.TO_EURUS'.tr().tr(),
                                onPressed: () {
                                  Navigator.of(context).pop();
                                  btnController.reset();
                                  common.transferToMySelf = false;
                                  common.isCenWithdraw = false;
                                },
                                buttonBGColor: FXColor.mainBlueColor)),
                        Container(
                            padding: EdgeInsets.only(top: 18),
                            child: SubmitButton(
                                btnController: btnController2,
                                label: 'AFTER_SCAN_PAGE.TO_ETHEREUM'.tr().tr(),
                                onPressed: () {
                                  Navigator.of(context).pop();
                                  btnController.reset();
                                  common.isCenWithdraw = true;
                                  common.transferToMySelf = false;
                                },
                                buttonBGColor: FXColor.mainBlueColor)),
                        Padding(
                          padding: EdgeInsets.only(top: 15, bottom: 15),
                          child: Divider(
                            color: lightGray,
                            thickness: 1,
                          ),
                        ),
                        Container(
                            padding: EdgeInsets.only(bottom: 18),
                            child: TextButton(
                              onPressed: () async {
                                btnController.reset();
                                String receiveText =
                                    await scanQRCode(context: context);
                                setState(() {
                                  if (EthAddress()
                                      .isValidEthereumAddress(receiveText)) {
                                    errorText = null;
                                  } else {
                                    errorText =
                                        "SEND_PAGE.ERROR.ADDRESS_NOT_VALID"
                                            .tr();
                                  }
                                });
                              },
                              style: TextButton.styleFrom(
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(15),
                                  side:
                                      BorderSide(color: FXColor.mainBlueColor),
                                ),
                                textStyle: FXUI.normalTextStyle,
                                padding: EdgeInsets.symmetric(
                                    vertical: 12.0, horizontal: 24.0),
                                tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                              ).copyWith(
                                foregroundColor: MaterialStateProperty.all(
                                    FXColor.mainBlueColor),
                                backgroundColor:
                                    MaterialStateProperty.all(Colors.white),
                              ),
                              child: Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                crossAxisAlignment: CrossAxisAlignment.center,
                                children: [
                                  Image.asset('images/icon_scan.png',
                                      package: 'euruswallet',
                                      width: 24,
                                      height: 24,
                                      fit: BoxFit.contain,
                                      color: FXColor.mainBlueColor),
                                  Container(width: 6),
                                  Text('AFTER_SCAN_PAGE.RESCAN'.tr(),
                                      style: FXUI.normalTextStyle.copyWith(
                                          color: FXColor.mainBlueColor)),
                                ],
                              ),
                            )),
                      ],
                    ),
                  ),
                ],
              )),
              padding: EdgeInsets.all(25),
              titleWidget: Padding(
                padding: EdgeInsets.only(bottom: 28),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Padding(
                      padding: EdgeInsets.only(right: 6),
                      child: Container(
                        width: 8,
                        height: 8,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: Color(0xFF66C381),
                        ),
                      ),
                    ),
                    Text('AFTER_SCAN_PAGE.EURUS_NETWORK'.tr(),
                        style: FXUI.normalTextStyle.copyWith(
                            fontWeight: FontWeight.w600,
                            color: FXColor.mainBlueColor)),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
