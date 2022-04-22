import 'dart:async';
import 'package:app_security_kit/app_security_kit.dart';
import 'package:app_storage_kit/app_storage_kit.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/callApiHandler.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/common/user_profile.dart';
import 'package:euruswallet/commonUI/customDialogBox.dart';
import 'package:euruswallet/commonUI/pinCode.dart';
import 'package:euruswallet/model/loginBySignModel.dart';
import 'package:euruswallet/model/registerByEmail.dart';
import 'package:euruswallet/model/setupPaymentWallet.dart';
import 'package:euruswallet/pages/centralized/changeLoginPW.dart';
import 'package:euruswallet/pages/centralized/changePaymentPW.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:pin_code_fields/pin_code_fields.dart';
import 'loadingTimerDialog.dart';

class VerifyCodePage extends StatefulWidget {
  VerifyCodePage();
  @override
  _VerifyCodePageState createState() => _VerifyCodePageState();
}

class _VerifyCodePageState extends State<VerifyCodePage> {
  StreamController<ErrorAnimationType> errorController =
      StreamController<ErrorAnimationType>();
  TextEditingController paymentPasswordController = TextEditingController();
  FocusNode myFocusNode = FocusNode();
  bool hasError = false;
  final _pwForm = GlobalKey<FormState>();
  bool alreadyClickReSendBtn = false;
  int remainingSecond = 0;
  Timer? _timer;
  bool isWalletSetup = false;
  final GlobalKey<LoadingTimerDialogState> _dialogKey = GlobalKey();

  @override
  void initState() {
    super.initState();

    forgetPaymentPassword();

    //testing Code
    autoFillVerifyCode();
    startTimer();
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  void startTimer() {
    const oneSec = const Duration(seconds: 1);
    _timer?.cancel();
    _timer = new Timer.periodic(
      oneSec,
      (Timer timer) {
        if (remainingSecond == 0) {
          setState(() {
            alreadyClickReSendBtn = false;
            timer.cancel();
          });
        } else {
          setState(() {
            remainingSecond--;
          });
        }
      },
    );
  }

  Future<void> forgetPaymentPassword() async {
    if (common.verifyCodePageType == VerifyCodePageType.forgotPaymentPw) {
      common.forgotPaymentPw = await api.forgetPaymentPw();
    }
  }

  void autoFillVerifyCode() {
    if (AutoFillAccount) {
      Future.delayed(const Duration(milliseconds: 1000), () async {
        setState(() {
          if (common.verifyCodePageType == VerifyCodePageType.register) {
            paymentPasswordController.text =
                common.registerByEmail?.data?.code ??
                    paymentPasswordController.text;
          }

          if (common.verifyCodePageType == VerifyCodePageType.forgotLoginPw) {
            paymentPasswordController.text = common.forgetLoginPw?.data?.code ??
                paymentPasswordController.text;
          }

          if (common.verifyCodePageType ==
              VerifyCodePageType.newDeviceResetPublicKey) {
            paymentPasswordController.text =
                common.registerDevice?.data?.code ??
                    paymentPasswordController.text;
          }

          if (common.verifyCodePageType == VerifyCodePageType.forgotPaymentPw) {
            paymentPasswordController.text =
                common.forgotPaymentPw?.data?.code ??
                    paymentPasswordController.text;
          }
        });
      });
    }
  }

  String? getEmail() {
    String? email;
    if (common.verifyCodePageType == VerifyCodePageType.register ||
        common.verifyCodePageType ==
            VerifyCodePageType.newDeviceResetPublicKey ||
        common.verifyCodePageType == VerifyCodePageType.forgotPaymentPw) {
      email = common.email;
    } else if (common.verifyCodePageType == VerifyCodePageType.forgotLoginPw) {
      email = common.forgetPwEmail;
    }
    return email;
  }

  Future<void> _successDialog() async {
    return showDialog(
        context: context,
        builder: (_) {
          return WillPopScope(
              onWillPop: () async => false,
              child: Dialog(
                  backgroundColor: Colors.transparent,
                  insetPadding: EdgeInsets.all(10),
                  child: Container(
                    width: double.infinity,
                    padding: EdgeInsets.symmetric(vertical: 28),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: FXUI.cricleRadius,
                    ),
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          "REGISTER.CREATE_SUCCESS".tr(),
                          style: FXUI.titleTextStyle.copyWith(
                            fontSize: 18,
                          ),
                        ),
                        Padding(padding: EdgeInsets.only(bottom: 87)),
                        SizedBox(
                          width: MediaQuery.of(context).size.width / 4.5,
                          child: Image(
                            image: AssetImage("images/uploadSuccessIcon.png",
                                package: 'euruswallet'),
                          ),
                        ),
                        Padding(padding: EdgeInsets.only(bottom: 15)),
                        Text(
                          'COMMON.SUCCESS'.tr(),
                          style: FXUI.titleTextStyle.copyWith(
                            fontSize: 24,
                          ),
                        ),
                        Padding(padding: EdgeInsets.only(bottom: 7)),
                        Text(
                          "REGISTER.CREATE_WALLET_SUCCESS".tr(),
                          style: FXUI.normalTextStyle.copyWith(
                            fontSize: 14,
                            color: FXColor.textGray,
                          ),
                        ),
                        Padding(padding: EdgeInsets.only(bottom: 60)),
                        Container(
                          padding: EdgeInsets.symmetric(horizontal: 20),
                          width: double.infinity,
                          child: TextButton(
                            style: TextButton.styleFrom(
                              padding: EdgeInsets.all(15),
                              backgroundColor: FXColor.mainBlueColor,
                              shape: RoundedRectangleBorder(
                                borderRadius: FXUI.cricleRadius,
                              ),
                            ),
                            onPressed: () {
                              Navigator.pop(_);
                              common.successMoveToHomePage(
                                userType: CurrentUserType.centralized,
                                email: common.email ?? '',
                                context: context,
                                loginPassword: common.loginPassword ?? '',
                                isRegister: true,
                              );
                            },
                            child: Text(
                              'REGISTER.WELCOME_TO_EURUS'.tr(),
                              style: FXUI.normalTextStyle
                                  .copyWith(color: Colors.white),
                            ),
                          ),
                        )
                      ],
                    ),
                  )));
        });
  }

  @override
  Widget build(BuildContext context) {
    return BackGroundImage(
      currentUserType: CurrentUserType.centralized,
      child: Scaffold(
        backgroundColor: Colors.transparent,
        appBar: WalletAppBar(
          title: common.verifyCodePageType == VerifyCodePageType.register
              ? 'REGISTER.MAIN_TITLE'.tr()
              : 'VERIFY_CODE_PAGE.VERIFY_CODE'.tr(),
          backButton: true,
        ),
        body: TopCircularContainer(
          height: size.heightWithoutAppBar,
          child: SingleChildScrollView(
            child: Padding(
              padding: EdgeInsets.only(left: 28, right: 28),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Form(
                    key: _pwForm,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Container(
                            height: 320,
                            child: TopCircularContainer(
                              child: Column(children: [
                                TopCircularContainer(
                                    width: size.blockSizeHorizontal * 100,
                                    child: Column(children: [
                                      Padding(
                                        padding: EdgeInsets.only(top: 30),
                                        child: Text(
                                            'VERIFY_CODE_PAGE.VERIFY_CODE'.tr(),
                                            style: FXUI.titleTextStyle
                                                .copyWith(fontSize: 24)),
                                      ),
                                      Padding(
                                        padding: EdgeInsets.only(top: 25),
                                        child: Text(
                                            'VERIFY_CODE_PAGE.VERIFY_CODE_TO_EMAIL'
                                                .tr(args: [getEmail() ?? ""]),
                                            style: FXUI.normalTextStyle.copyWith(
                                                fontSize: 12,
                                                color: FXColor
                                                    .centralizedGrayTextColor)),
                                      ),
                                      Padding(
                                        padding: EdgeInsets.only(top: 15),
                                        child: PinCode(
                                            hasError: hasError,
                                            errorController: errorController,
                                            myFocusNode: myFocusNode,
                                            paymentPasswordController:
                                                paymentPasswordController,
                                            onChanged: (value) {
                                              print(value);
                                              setState(() {
                                                hasError = false;
                                              });
                                            },
                                            onTap: () {},
                                            onCompleted:
                                                (String pinCode) async {
                                              print("Completed");
                                              if (pinCode.isNotEmpty) {
                                                if (pinCode.length == 6) {
                                                  EasyLoading.show(
                                                      status:
                                                          'COMMON.LOADING_W_DOT'
                                                              .tr());
                                                  if (common
                                                          .verifyCodePageType ==
                                                      VerifyCodePageType
                                                          .forgotLoginPw) {
                                                    common.codeVerification = await api
                                                        .forGetLoginPwCodeVerification(
                                                            code: pinCode,
                                                            email: common
                                                                    .forgetPwEmail ??
                                                                '');
                                                    if (await common.checkApiError(
                                                        context: context,
                                                        errorString: common
                                                            .codeVerification
                                                            ?.message,
                                                        returnCode: common
                                                            .codeVerification
                                                            ?.returnCode,
                                                        btnColor: FXColor
                                                            .mainBlueColor)) {
                                                      await NormalStorageKit()
                                                          .setValue(
                                                              common
                                                                      .codeVerification
                                                                      ?.data
                                                                      ?.token ??
                                                                  '',
                                                              'apiAccessToken_');
                                                      common.pushPage(
                                                          page: ChangeLoginPWPage(
                                                              changePasswordType:
                                                                  ChangeLoginPasswordType
                                                                      .resetLoginPassword),
                                                          context: context);
                                                    }
                                                  }
                                                  if (common
                                                          .verifyCodePageType ==
                                                      VerifyCodePageType
                                                          .register) {
                                                    common.codeVerification =
                                                        await api
                                                            .codeVerification(
                                                                code: pinCode,
                                                                email: common
                                                                        .email ??
                                                                    '');
                                                    print(
                                                        'call api successful');
                                                    if (await common.checkApiError(
                                                        context: context,
                                                        errorString: common
                                                            .codeVerification
                                                            ?.message,
                                                        returnCode: common
                                                            .codeVerification
                                                            ?.returnCode,
                                                        btnColor: FXColor
                                                            .mainBlueColor)) {
                                                      EasyLoading.dismiss();
                                                      showDialog(
                                                          context: context,
                                                          builder: (BuildContext
                                                              context) {
                                                            return LoadingTimerDialog(
                                                                key:
                                                                    _dialogKey);
                                                          });
                                                      print(
                                                          'codeVerification successful');
                                                      await NormalStorageKit()
                                                          .setValue(
                                                              common
                                                                      .codeVerification
                                                                      ?.data
                                                                      ?.token ??
                                                                  '',
                                                              'apiAccessToken_');
                                                      await NormalStorageKit().setValue(
                                                          common
                                                                  .codeVerification
                                                                  ?.data
                                                                  ?.expiredTime
                                                                  .toString() ??
                                                              '',
                                                          'apiAccessTokenExpiryTime_');
                                                      // await NormalStorageKit().setValue(refreshToken.data.expiryTime.toString(),'apiAccessTokenExpiryTime_');
                                                      final decryptHepler =
                                                          DecryptionHelper(
                                                              privateKey: common
                                                                      .rsaPrivateKey ??
                                                                  '');
                                                      print(
                                                          "common.addressPair.privateKey2:${common.rsaPrivateKey}");
                                                      print(
                                                          "common.addressPair.privateKey:${common.loginAddressPair?.privateKey}");
                                                      print(
                                                          "common.codeVerification.data.mnemonic:${common.codeVerification?.data?.mnemonic}");
                                                      common.serverMnemonic = decryptHepler
                                                          .decryptRAEncryption(common
                                                                  .codeVerification
                                                                  ?.data
                                                                  ?.mnemonic ??
                                                              '');
                                                      common.serverAddressPair =
                                                          await common.getAddressPair(
                                                              email:
                                                                  common.email,
                                                              password: common
                                                                  .paymentPassword,
                                                              mnemonic: common
                                                                  .serverMnemonic,
                                                              addressPairType:
                                                                  AddressPairType
                                                                      .paymentPw);
                                                      print(
                                                          "common.serverMnemonic:${common.serverMnemonic}");
                                                      SetupPaymentWallet
                                                          setupPaymentWallet =
                                                          await api.setupPaymentWallet(
                                                              userID: common
                                                                  .registerByEmail
                                                                  ?.data
                                                                  ?.userId,
                                                              walletAddress: common
                                                                  .serverAddressPair
                                                                  ?.address);
                                                      if (await common.checkApiError(
                                                          context: context,
                                                          errorString:
                                                              setupPaymentWallet
                                                                  .message,
                                                          returnCode:
                                                              setupPaymentWallet
                                                                  .returnCode,
                                                          btnColor: FXColor
                                                              .mainBlueColor)) {
                                                        // EasyLoading.show(
                                                        //     status:
                                                        //         'COMMON.LOADING_W_DOT'
                                                        //             .tr());
                                                        common.cenMainNetWalletAddress =
                                                            setupPaymentWallet
                                                                .data
                                                                ?.mainnetWalletAddress;
                                                        common.cenUserWalletAddress =
                                                            setupPaymentWallet
                                                                .data
                                                                ?.walletAddress;
                                                        common.ownerWalletAddress =
                                                            common
                                                                .serverAddressPair
                                                                ?.address;
                                                        await NormalStorageKit()
                                                            .setValue(
                                                                common.email ??
                                                                    '',
                                                                'cenUserEmail');
                                                        await NormalStorageKit()
                                                            .setValue(
                                                                common.loginPassword ??
                                                                    '',
                                                                'cenUserPassword');

                                                        final PasswordEncryptHelper
                                                            pwHelper =
                                                            PasswordEncryptHelper(
                                                                password: common
                                                                        .loginPassword ??
                                                                    '');

                                                        UserProfile
                                                            userProfile =
                                                            UserProfile
                                                                .fromJson({
                                                          "userType":
                                                              CurrentUserType
                                                                  .centralized,
                                                          "address":
                                                              setupPaymentWallet
                                                                      .data
                                                                      ?.walletAddress ??
                                                                  '0x0',
                                                          "email": common.email,
                                                          "encryptedAddress":
                                                              pwHelper.encryptWPwd(
                                                                  setupPaymentWallet
                                                                          .data
                                                                          ?.walletAddress ??
                                                                      ''),
                                                          "encryptedPrivateKey":
                                                              pwHelper.encryptWPwd(common
                                                                      .loginAddressPair
                                                                      ?.privateKey ??
                                                                  ''),
                                                          "lastLoginTime": DateTime
                                                                  .now()
                                                              .millisecondsSinceEpoch
                                                              .toString(),
                                                        });

                                                        setActiveAccount(
                                                            pwHelper.encryptWPwd(
                                                                setupPaymentWallet
                                                                        .data
                                                                        ?.walletAddress ??
                                                                    ''));

                                                        if (userProfile
                                                                .address !=
                                                            '0x0')
                                                          await setAcToLocal(
                                                              userProfile);

                                                        _dialogKey.currentState!
                                                            .onCompleted();
                                                        await Future.delayed(
                                                            const Duration(
                                                                milliseconds:
                                                                    600), () {
                                                          _successDialog();
                                                        });
                                                      } else {
                                                        _dialogKey.currentState!
                                                            .onCompleted();
                                                      }
                                                      EasyLoading.dismiss();
                                                    } else {
                                                      EasyLoading.dismiss();
                                                      errorController.add(
                                                          ErrorAnimationType
                                                              .shake);
                                                      setState(() {
                                                        hasError = true;
                                                      });
                                                    }
                                                  }
                                                  if (common
                                                          .verifyCodePageType ==
                                                      VerifyCodePageType
                                                          .newDeviceResetPublicKey) {
                                                    common.codeVerification =
                                                        await api.verifyDevice(
                                                            code: pinCode);
                                                    if (await common.checkApiError(
                                                        context: context,
                                                        errorString: common
                                                            .codeVerification
                                                            ?.message,
                                                        returnCode: common
                                                            .codeVerification
                                                            ?.returnCode,
                                                        btnColor: FXColor
                                                            .mainBlueColor)) {
                                                      LoginBySignModel result =
                                                          await api.loginBySignature(
                                                              email: common
                                                                      .email
                                                                      ?.toLowerCase() ??
                                                                  '',
                                                              password: common
                                                                      .loginPassword ??
                                                                  '');

                                                      if (await common
                                                          .checkApiError(
                                                              context: context,
                                                              errorString:
                                                                  result
                                                                      .message,
                                                              returnCode: result
                                                                  .returnCode)) {
                                                        if (result.decryptedMnemonic != null &&
                                                            result.isMetaMaskUser !=
                                                                null &&
                                                            !(result
                                                                .isMetaMaskUser!)) {
                                                          await common.successMoveToHomePage(
                                                              userType:
                                                                  CurrentUserType
                                                                      .centralized,
                                                              email: common
                                                                      .email ??
                                                                  '',
                                                              context: context,
                                                              loginPassword:
                                                                  common.loginPassword ??
                                                                      '',
                                                              loginBySignModel:
                                                                  result);
                                                        } else {
                                                          common.registerDevice =
                                                              await api
                                                                  .registerDevice();
                                                          if (await common.checkApiError(
                                                              context: context,
                                                              errorString: common
                                                                  .registerDevice
                                                                  ?.message,
                                                              returnCode: common
                                                                  .registerDevice
                                                                  ?.returnCode)) {
                                                            common.verifyCodePageType =
                                                                VerifyCodePageType
                                                                    .newDeviceResetPublicKey;
                                                            common.pushPage(
                                                                page:
                                                                    VerifyCodePage(),
                                                                context:
                                                                    context);
                                                          }
                                                        }
                                                      }
                                                    }
                                                    EasyLoading.dismiss();
                                                  }
                                                  if (common
                                                          .verifyCodePageType ==
                                                      VerifyCodePageType
                                                          .forgotPaymentPw) {
                                                    common.codeVerification =
                                                        await api
                                                            .verifyForgetPaymentPassword(
                                                                code: pinCode);
                                                    if (await common.checkApiError(
                                                        context: context,
                                                        errorString: common
                                                            .codeVerification
                                                            ?.message,
                                                        returnCode: common
                                                            .codeVerification
                                                            ?.returnCode,
                                                        btnColor: FXColor
                                                            .mainBlueColor)) {
                                                      await NormalStorageKit()
                                                          .setValue(
                                                              common
                                                                      .codeVerification
                                                                      ?.data
                                                                      ?.token ??
                                                                  '',
                                                              'apiAccessToken_');
                                                      common
                                                          .serverMnemonic = DecryptionHelper(
                                                              privateKey: common
                                                                      .rsaPrivateKey ??
                                                                  '')
                                                          .decryptRAEncryption(common
                                                                  .codeVerification
                                                                  ?.data
                                                                  ?.mnemonic ??
                                                              '');
                                                      print(
                                                          "common.serverMnemonic:${common.serverMnemonic}");
                                                      String result;
                                                      common.changePaymentPasswordType =
                                                          ChangePaymentPasswordType
                                                              .forgetPaymentPassword;
                                                      result =
                                                          await common.pushPage(
                                                        page:
                                                            ChangePaymentPWPage(),
                                                        context: context,
                                                      );
                                                      if (result == 'done')
                                                        common.logoutFunction
                                                            ?.add(true);
                                                    }
                                                    EasyLoading.dismiss();
                                                  }
                                                }
                                              }
                                            }),
                                      ),
                                      InkWell(
                                          child: Padding(
                                            padding: EdgeInsets.only(top: 13),
                                            child: Text(
                                                alreadyClickReSendBtn
                                                    ? 'VERIFY_CODE_PAGE.RESENT_SECONDS'
                                                        .tr(args: [
                                                        remainingSecond
                                                            .toString()
                                                      ])
                                                    : 'VERIFY_CODE_PAGE.RESENT'
                                                        .tr(),
                                                style: FXUI.normalTextStyle
                                                    .copyWith(
                                                        color: FXColor
                                                            .mainBlueColor)),
                                          ),
                                          onTap: alreadyClickReSendBtn
                                              ? () {}
                                              : () async {
                                                  String? errorString;
                                                  int? returnCode;

                                                  if (common
                                                          .verifyCodePageType ==
                                                      VerifyCodePageType
                                                          .register) {
                                                    RegisterByEmail
                                                        registerByEmail =
                                                        await api
                                                            .reSendCodeVerification(
                                                                userID: common
                                                                    .registerByEmail
                                                                    ?.data
                                                                    ?.userId);
                                                    errorString =
                                                        registerByEmail.message;
                                                    returnCode = registerByEmail
                                                        .returnCode;
                                                  }
                                                  if (common
                                                          .verifyCodePageType ==
                                                      VerifyCodePageType
                                                          .forgotLoginPw) {
                                                    common.forgetLoginPw =
                                                        await api.forgetLoginPw(
                                                            email: common
                                                                .forgetPwEmail
                                                                ?.toLowerCase());
                                                    errorString = common
                                                        .forgetLoginPw?.message;
                                                    returnCode = common
                                                        .forgetLoginPw
                                                        ?.returnCode;
                                                  }

                                                  if (common
                                                          .verifyCodePageType ==
                                                      VerifyCodePageType
                                                          .newDeviceResetPublicKey) {
                                                    common.registerDevice =
                                                        await api
                                                            .registerDevice();
                                                    errorString = common
                                                        .registerDevice
                                                        ?.message;
                                                    returnCode = common
                                                        .registerDevice
                                                        ?.returnCode;
                                                  }

                                                  if (common
                                                          .verifyCodePageType ==
                                                      VerifyCodePageType
                                                          .forgotPaymentPw) {
                                                    common.forgotPaymentPw =
                                                        await api
                                                            .forgetPaymentPw();
                                                    errorString = common
                                                        .forgotPaymentPw
                                                        ?.message;
                                                    returnCode = common
                                                        .forgotPaymentPw
                                                        ?.returnCode;
                                                  }

                                                  if (returnCode == 0) {
                                                    await showDialog(
                                                      context: context,
                                                      builder: (BuildContext
                                                          context) {
                                                        return CustomDialogBox(
                                                          btnColor: FXColor
                                                              .mainBlueColor,
                                                          descriptions:
                                                              'VERIFY_CODE_PAGE.SENT'
                                                                  .tr(),
                                                          buttonText:
                                                              "COMMON.OK".tr(),
                                                        );
                                                      },
                                                    );
                                                  } else {
                                                    setState(() {
                                                      alreadyClickReSendBtn =
                                                          true;
                                                      remainingSecond = (int.tryParse(
                                                                  errorString?.replaceAll(
                                                                          new RegExp(
                                                                              r'[^0-9]'),
                                                                          '') ??
                                                                      '') ??
                                                              0) *
                                                          60;
                                                      startTimer();
                                                    });
                                                  }

                                                  if (await common
                                                      .checkApiError(
                                                          context: context,
                                                          errorString:
                                                              errorString,
                                                          returnCode:
                                                              returnCode,
                                                          btnColor: FXColor
                                                              .mainBlueColor)) {
                                                    autoFillVerifyCode();
                                                  }
                                                })
                                    ])),
                              ]),
                            ))
                      ],
                    ),
                  )
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
