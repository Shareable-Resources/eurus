import 'package:app_storage_kit/secure_storage.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/userKYCStatus.dart';
import 'package:euruswallet/pages/centralized/changeLoginPW.dart';
import 'package:euruswallet/pages/centralized/verifyCode.dart';
import 'package:euruswallet/pages/settingSubpages/KYC/kycStatus.dart';
import 'package:euruswallet/pages/settingSubpages/biometricsPage.dart';
import 'package:euruswallet/pages/settingSubpages/cardBtnRow.dart';
import 'package:euruswallet/pages/settingSubpages/cardContainer.dart';
import 'package:euruswallet/pages/settingSubpages/changeLockerPW.dart';
import 'package:euruswallet/pages/settingSubpages/languagePage.dart';
import 'package:euruswallet/pages/settingSubpages/settingAppBar.dart';
import 'package:euruswallet/pages/settingSubpages/versionPage.dart';
import 'package:flutter/cupertino.dart';
import 'package:livechat_inc/livechat_inc.dart';
import 'package:url_launcher/url_launcher.dart';

import 'centralized/changePaymentPW.dart';

class SettingPage extends StatefulWidget {
  SettingPage({
    required this.userSuffix,
    required this.backupSeedPhraseFnc,
  }) : super();

  final String userSuffix;
  final Function backupSeedPhraseFnc;

  @override
  _SettingPageState createState() => _SettingPageState();
}

class _SettingPageState extends State<SettingPage> {
  bool? _defaultShowAssetBalance;
  UserKYCStatus? _kycStatus;

  @override
  void initState() {
    _getDefaultAssetBalancePrivacyVal();
    if (isCentralized()) _getKYCStatus();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: FXColor.veryLightGreyTextColor,
      appBar: SettingAppBar(false),
      body: Container(
        child: SingleChildScrollView(
          child: SafeArea(
            child: Column(
              children: [
                CardContainer('SETTING_PAGE.SECURITY_CARD.MAIN_TITLE'.tr(),
                    _genSecurityItems()),
                CardContainer(
                    'SETTING_PAGE.GENERAL_CARD.TITLE'.tr(), _genGeneralItems()),
                CardContainer('SETTING_PAGE.HELP_CENTER_CARD.TITLE'.tr(),
                    _genHelpCenterItems()),
                Padding(
                    padding:
                        EdgeInsets.only(bottom: BOTTOM_NAV_TAB_BAR_HEIGHT)),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _genSecurityItems() {
    var _themeColor = common.getBackGroundColor();
    var isActive = common.currentUserProfile!.seedPraseBackuped == true;

    List<Widget> items = [
      common.currentUserProfile!.decenUserType != DecenUserType.imported
          ? !isCentralized()
              ? CardBtnRow(
                  'SETTING_PAGE.SECURITY_CARD.BACKUP_MPHRASE.STATUS'.tr(),
                  onPressFnc: () async {
                    if (isActive) return;
                    await widget.backupSeedPhraseFnc();
                    setState(() {});
                  },
                  borderBtm: true,
                  btnContent: Text(
                    isActive
                        ? 'SETTING_PAGE.SECURITY_CARD.BACKUP_MPHRASE.COMPLETED'
                            .tr()
                        : 'SETTING_PAGE.SECURITY_CARD.BACKUP_MPHRASE.INCOMPLETED'
                            .tr(),
                    style: FXUI.normalTextStyle.copyWith(
                        color: isActive
                            ? FXColor.greenColor
                            : FXColor.pinkRedColor),
                  ),
                  customBtmContent: !isActive
                      ? Container(
                          decoration: BoxDecoration(
                            borderRadius: FXUI.cricleRadius,
                            color: FXColor.mainDeepBlueColor,
                          ),
                          margin: EdgeInsets.only(left: 30, right: 32, top: 20),
                          padding: EdgeInsets.all(13),
                          width: double.infinity,
                          child: Text(
                            'COMMON.BACKUP_NOW'.tr(),
                            textAlign: TextAlign.center,
                            style: FXUI.normalTextStyle
                                .copyWith(color: Colors.white, fontSize: 16),
                          ),
                        )
                      : Container(),
                )
              : Container()
          : Container(),

      CardBtnRow(
        !isCentralized()
            ? 'SETTING_PAGE.SECURITY_CARD.WALLET_LOCKER_BTN'.tr()
            : 'SETTING_PAGE.SECURITY_CARD.CHANGE_LOGINPW_BTN'.tr(),
        onPressFnc: () async {
          String result = '';
          if (!isCentralized()) {
            result = await common.pushPage(
              page: ChangeLockerPWPage(userSuffix: widget.userSuffix),
              context: context,
            );
          }

          if (isCentralized()) {
            result = await common.pushPage(
              page: ChangeLoginPWPage(
                  changePasswordType:
                      ChangeLoginPasswordType.changeLoginPassword),
              context: context,
            );
          }

          if (result == 'done') common.logoutFunction?.add(true);
        },
        borderBtm: true,
      ),
      if (isCentralized())
        CardBtnRow(
          'SETTING_PAGE.SECURITY_CARD.CHANGE_PAYMENTPW_BTN'.tr(),
          onPressFnc: () async {
            common.changePaymentPasswordType =
                ChangePaymentPasswordType.changePaymentPassword;
            String? result;
            result = await common.pushPage(
              page: ChangePaymentPWPage(),
              context: context,
            );
            if (result == 'done') common.logoutFunction?.add(true);
          },
          borderBtm: true,
        ),
      if (isCentralized())
        CardBtnRow(
          'SETTING_PAGE.SECURITY_CARD.FORGET_PAYMENTPW_BTN'.tr(),
          onPressFnc: () {
            showDialog(
                context: context,
                builder: (BuildContext context) {
                  return CustomConfirmDialog(
                    title:
                        'SETTING_PAGE.SECURITY_CARD.FORGET_PAYMENTPW_BTN'.tr(),
                    descriptions: 'FORGOT_PAYMENT_PWD_DIALOG.CONTENT'.tr(),
                    showCancelButton: true,
                    useSubmitButton: false,
                    buttonText: "COMMON.CONFIRM".tr(),
                    btnHandler: () async {
                      common.verifyCodePageType =
                          VerifyCodePageType.forgotPaymentPw;
                      await common.pushPage(
                          page: VerifyCodePage(), context: context);
                    },
                  );
                });
          },
          borderBtm: true,
        ),
      isCentralized()
          ? CardBtnRow(
              'SETTING_PAGE.SECURITY_CARD.KYC_STATUS'.tr(),
              onPressFnc: () {
                CommonMethod().pushPage(
                  page: KYCStatusPage(),
                  context: context,
                );
              },
              btnContent: Row(children: [
                _kycStatus != null
                    ? Row(children: [
                        Text('SETTING_PAGE.SECURITY_CARD.LEVEL'.tr(),
                            style: FXUI.normalTextStyle.copyWith(
                              fontSize: 12,
                              color: FXColor.mediumGrayColor,
                            )),
                        Padding(
                            padding: EdgeInsets.only(left: 4),
                            child: Text(
                                _kycStatus?.data?.kycdata?.length != null
                                    ? _kycStatus!.data!.kycLevel.toString()
                                    : "0",
                                style: FXUI.normalTextStyle.copyWith(
                                  fontSize: 12,
                                  color: FXColor.mainBlueColor,
                                )))
                      ])
                    : Container(),
                Icon(Icons.navigate_next_outlined, color: FXColor.untabColor),
              ]),
              borderBtm: true,
            )
          : Container(),
      !isCentralized()
          ? CardBtnRow(
              'SETTING_PAGE.SECURITY_CARD.BIOMETRICS'.tr(),
              onPressFnc: () {
                CommonMethod().pushPage(
                  page: BiometricsPage(userSuffix: widget.userSuffix),
                  context: context,
                );
              },
              borderBtm: true,
            )
          : Container(),
      // CardBtnRow(
      //   'SETTING_PAGE.SECURITY_CARD.ASSET_PRIVACY.TITLE'.tr(),
      //   subTitle: 'SETTING_PAGE.SECURITY_CARD.ASSET_PRIVACY.DESC'.tr(),
      //   onPressFnc: () {},
      //   btnContent: CupertinoSwitch(
      //     value: _defaultShowAssetBalance ?? false,
      //     onChanged: _setDefaultAssetBalancePrivacyVal,
      //     activeColor: _themeColor,
      //   ),
      //   borderBtm: true,
      // ),
      CardBtnRow(
        'SETTING_PAGE.SECURITY_CARD.LOGOUT'.tr(),
        onPressFnc: _confirmLogoutDialog, //widget.logoutFnc,
        btnContent: Icon(
          Icons.logout,
          color: FXColor.mediumGrayColor,
        ),
      ),
    ];

    return Column(children: items);
  }

  Widget _genHelpCenterItems() {
    return Column(
      children: [
        CardBtnRow('SETTING_PAGE.HELP_CENTER_CARD.FAQ'.tr(),
            borderBtm: true,
            onPressFnc: () async =>
                await launch('SETTING_PAGE.HELP_CENTER_CARD.FAQ_LINK'.tr())),
        CardBtnRow(
          'SETTING_PAGE.HELP_CENTER_CARD.CS'.tr(),
          borderBtm: true,
          onPressFnc: () => common.openLiveChat(
            visitorName: isCentralized() ? common.email : common.currentAddress,
            visitorEmail: isCentralized() ? common.email : null,
          ),
        ),
        CardBtnRow(
          'SETTING_PAGE.HELP_CENTER_CARD.TNC.TITLE'.tr(),
          borderBtm: true,
          onPressFnc: () async =>
              await launch('SETTING_PAGE.HELP_CENTER_CARD.TNC.LINK'.tr()),
        ),
        CardBtnRow(
          'SETTING_PAGE.HELP_CENTER_CARD.WALLET_VERSION'.tr(),
          onPressFnc: () {
            CommonMethod().pushPage(page: VersionPage(), context: context);
          },
        ),
      ],
    );
  }

  Widget _genGeneralItems() {
    return Column(
      children: [
        CardBtnRow('SETTING_PAGE.GENERAL_CARD.LANGUAGE'.tr(),
            onPressFnc: () async {
          await CommonMethod().pushPage(page: LanguagePage(), context: context);
          setState(() {});
          if (common.updateNavbarOnChangeLang != null)
            common.updateNavbarOnChangeLang!();
        }),
      ],
    );
  }

  Future<bool> _getDefaultAssetBalancePrivacyVal() async {
    String defaultShowAB = await NormalStorageKit()
            .readValue('defaultShowAssetBalance_${widget.userSuffix}') ??
        '';

    if (isEmptyString(string: defaultShowAB)) {
      NormalStorageKit()
          .setValue('1', 'defaultShowAssetBalance_${widget.userSuffix}');

      return true;
    }

    setState(() {
      _defaultShowAssetBalance = defaultShowAB == '1' ? false : true;
    });

    return defaultShowAB == '1' ? true : false;
  }

  Future<void> _setDefaultAssetBalancePrivacyVal(bool val) async {
    await NormalStorageKit().setValue(val == true ? '0' : '1',
        'defaultShowAssetBalance_${widget.userSuffix}');

    setState(() {
      _defaultShowAssetBalance = val;
      common.showBalance = !val;
    });
  }

  Future<void> _getKYCStatus() async {
    if (_kycStatus == null) {
      UserKYCStatus status = await api.getKYCStatus();
      setState(() {
        _kycStatus = status;
      });
    }
  }

  Future<bool?> _confirmLogoutDialog() async {
    var logoutWarnIconWidth = MediaQuery.of(context).size.width / 3.5;

    return showDialog(
      context: context,
      builder: (_) {
        return Dialog(
          backgroundColor: Colors.transparent,
          insetPadding: EdgeInsets.all(10),
          child: Container(
            width: double.infinity,
            padding: EdgeInsets.symmetric(vertical: 25, horizontal: 27),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: FXUI.cricleRadius,
            ),
            child: Stack(
              children: [
                SingleChildScrollView(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Center(
                        child: Text(
                          'LOGOUT_DIALOG.TITLE'.tr(),
                          style: FXUI.normalTextStyle.copyWith(
                            fontWeight: FontWeight.bold,
                            fontSize: 18,
                          ),
                        ),
                      ),
                      Padding(padding: EdgeInsets.only(top: 50)),
                      SizedBox(
                        width: logoutWarnIconWidth,
                        child: Padding(
                          padding: EdgeInsets.only(
                            left: logoutWarnIconWidth / 5,
                          ),
                          child: Image(
                            image: AssetImage(
                              'images/logoutWarnIcon.png',
                              package: 'euruswallet',
                            ),
                          ),
                        ),
                      ),
                      Padding(padding: EdgeInsets.only(top: 30)),
                      Text(
                        'LOGOUT_DIALOG.CONTENT'.tr(),
                        style: FXUI.normalTextStyle
                            .copyWith(color: FXColor.textGray, fontSize: 14),
                      ),
                      Padding(padding: EdgeInsets.only(top: 60)),
                      SizedBox(
                        width: double.infinity,
                        child: TextButton(
                          style: TextButton.styleFrom(
                            padding: EdgeInsets.all(15),
                            backgroundColor: common.getBackGroundColor(),
                            shape: RoundedRectangleBorder(
                              borderRadius: FXUI.cricleRadius,
                            ),
                          ),
                          onPressed: () {
                            Navigator.pop(_);
                            common.logoutFunction?.add(true);
                          },
                          child: Text(
                            'COMMON.CONFIRM'.tr(),
                            style: FXUI.normalTextStyle
                                .copyWith(color: Colors.white),
                          ),
                        ),
                      ),
                      Padding(padding: EdgeInsets.only(top: 15)),
                      SizedBox(
                        width: double.infinity,
                        child: TextButton(
                          style: TextButton.styleFrom(
                            padding: EdgeInsets.all(15),
                            shape: RoundedRectangleBorder(
                              side: BorderSide(
                                width: 1,
                                color: common.getBackGroundColor(),
                              ),
                              borderRadius: FXUI.cricleRadius,
                            ),
                          ),
                          onPressed: () => Navigator.pop(_),
                          child: Text(
                            'COMMON.CANCEL'.tr(),
                            style: FXUI.normalTextStyle.copyWith(
                              color: common.getBackGroundColor(),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                Positioned(
                  top: -12,
                  right: -17,
                  child: IconButton(
                    icon: Icon(
                      Icons.close,
                      color: Colors.black.withOpacity(0.5),
                    ),
                    onPressed: () => Navigator.pop(_),
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}
