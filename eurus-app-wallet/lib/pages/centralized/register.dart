import 'package:easy_localization/easy_localization.dart';
import 'package:email_validator/email_validator.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/pages/centralized/verifyCode.dart';
import 'package:url_launcher/url_launcher.dart';

class RegisterPage extends StatefulWidget {
  RegisterPage();
  @override
  _RegisterPageState createState() => _RegisterPageState();
}

enum TextFieldType {
  loginPassword,
  confirmLoginPassword,
  paymentPassword,
  confirmPaymentPassword
}

class _RegisterPageState extends State<RegisterPage> {
  final TextEditingController _emailTc = TextEditingController();
  final TextEditingController _curLoginPwTc = TextEditingController();
  final TextEditingController _confirmLoginPwTc = TextEditingController();
  final TextEditingController _curPaymentPwTc = TextEditingController();
  final TextEditingController _confirmPaymentPwTc = TextEditingController();
  bool isLoginPasswordMasked = false;
  bool isConfirmLoginPasswordMasked = false;
  bool isPaymentPasswordMasked = false;
  bool isConfirmPaymentPasswordMasked = false;
  final _pwForm = GlobalKey<FormState>();
  bool _checkbox = false;
  bool loginPwNumberAndAlphabeticalError = true;
  bool loginPwLengthError = true;
  bool confirmLoginPwNumberAndAlphabeticalError = true;
  bool confirmLoginPwLengthError = true;
  bool loginPwAndConfirmLoginPwNotSameError = true;

  bool paymentPwNumberAndAlphabeticalError = true;
  bool paymentPwLengthError = true;
  bool confirmPaymentPwNumberAndAlphabeticalError = true;
  bool confirmPaymentPwLengthError = true;
  bool paymentPwAndConfirmPaymentPwNotSameError = true;
  bool loginPwAndPaymentPwSameError = true;

  String hintEnterLoginPassword = 'REGISTER.INPUT.LOGIN_PW_DESC';
  String hintEnterYourTransactionCode = 'REGISTER.PAYMENT_PW_SECTION_TITLE';
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    if (AutoFillAccount) {
      _emailTc.text =
          envType == EnvType.Testnet || envType == EnvType.Production
              ? "eu31@18m.dev"
              : "ken575@goldhub.hk";
      _curLoginPwTc.text = "aaaaaaa1";
      _confirmLoginPwTc.text = "aaaaaaa1";
      _curPaymentPwTc.text = "bbbbbbb1";
      _confirmPaymentPwTc.text = "bbbbbbb1";
    }
    super.initState();
  }

  String? errorMessage({
    required TextFieldType type,
    String? v,
  }) {
    String? errorMessageString;
    String textFieldName = '';
    if (type == TextFieldType.loginPassword) {
      textFieldName = 'REGISTER.INPUT.CONFIRM_LOGIN_PW_DESC'.tr();
    }

    if (type == TextFieldType.confirmLoginPassword) {
      textFieldName = 'REGISTER.INPUT.CONFIRM_LOGIN_PW_TITLE'.tr();
    }
    if (type == TextFieldType.paymentPassword) {
      textFieldName = 'REGISTER.INPUT.PAYMENT_PW_TITLE'.tr();
    }

    if (type == TextFieldType.confirmPaymentPassword) {
      textFieldName = 'REGISTER.INPUT.PAYMENT_PW_DESC'.tr();
    }

    if (type == TextFieldType.loginPassword ||
        type == TextFieldType.confirmLoginPassword) {
      if (isEmptyString(string: v)) return hintEnterLoginPassword.tr();
      if (common.isNotContain8To20Characters(text: _curLoginPwTc.text)) {
        errorMessageString = errorMessageString ??
            'REGISTER.CONTAIN_8_TO_20'.tr(args: ['$textFieldName']);
        loginPwLengthError = true;
      } else {
        loginPwLengthError = false;
      }

      if (common.isNotContain8To20Characters(text: _confirmLoginPwTc.text)) {
        errorMessageString = errorMessageString ??
            'REGISTER.CONTAIN_8_TO_20'.tr(args: ['$textFieldName']);
        confirmLoginPwLengthError = true;
      } else {
        confirmLoginPwLengthError = false;
      }

      loginPwNumberAndAlphabeticalError =
          common.isNotContainDigitAndCharacter(text: _curLoginPwTc.text);
      if (loginPwNumberAndAlphabeticalError) {
        errorMessageString = errorMessageString ??
            'REGISTER.DIGIT_AND_CHARACTER'.tr(args: ['$textFieldName']);
      }

      confirmLoginPwNumberAndAlphabeticalError =
          common.isNotContainDigitAndCharacter(text: _confirmLoginPwTc.text);
      if (confirmLoginPwNumberAndAlphabeticalError) {
        errorMessageString = errorMessageString ??
            'REGISTER.DIGIT_AND_CHARACTER'.tr(args: ['$textFieldName']);
      }

      if (_curLoginPwTc.text != _confirmLoginPwTc.text) {
        errorMessageString =
            errorMessageString ?? 'REGISTER.LOGIN_PW_CONFIRM'.tr();
        loginPwAndConfirmLoginPwNotSameError = true;
      } else {
        loginPwAndConfirmLoginPwNotSameError = false;
      }
    }

    if (type == TextFieldType.paymentPassword) {
      if (isEmptyString(string: v)) return hintEnterYourTransactionCode.tr();

      if (common.isNotContain8To20Characters(text: _curPaymentPwTc.text)) {
        errorMessageString = errorMessageString ??
            'REGISTER.CONTAIN_8_TO_20'.tr(args: ['$textFieldName']);
        paymentPwLengthError = true;
      } else {
        paymentPwLengthError = false;
      }

      paymentPwNumberAndAlphabeticalError =
          common.isNotContainDigitAndCharacter(text: _curPaymentPwTc.text);
      if (paymentPwNumberAndAlphabeticalError) {
        errorMessageString = errorMessageString ??
            'REGISTER.DIGIT_AND_CHARACTER'.tr(args: ['$textFieldName']);
      }

      if (_curLoginPwTc.text == _curPaymentPwTc.text) {
        errorMessageString =
            errorMessageString ?? 'REGISTER.LOGIN_PW_NOT_EQUAL_PAYMENT_PW'.tr();
        loginPwAndPaymentPwSameError = true;
      } else {
        loginPwAndPaymentPwSameError = false;
      }
    }

    if (type == TextFieldType.confirmPaymentPassword) {
      if (isEmptyString(string: v)) return hintEnterYourTransactionCode.tr();

      if (common.isNotContain8To20Characters(text: _confirmPaymentPwTc.text)) {
        errorMessageString = errorMessageString ??
            'REGISTER.CONTAIN_8_TO_20'.tr(args: ['$textFieldName']);
        confirmPaymentPwLengthError = true;
      } else {
        confirmPaymentPwLengthError = false;
      }

      confirmPaymentPwNumberAndAlphabeticalError =
          common.isNotContainDigitAndCharacter(text: _confirmPaymentPwTc.text);
      if (confirmPaymentPwNumberAndAlphabeticalError) {
        errorMessageString = errorMessageString ??
            'REGISTER.DIGIT_AND_CHARACTER'.tr(args: ['$textFieldName']);
      }

      if (_curPaymentPwTc.text != _confirmPaymentPwTc.text) {
        errorMessageString =
            errorMessageString ?? 'REGISTER.PAYMENT_PW_CONFIRM'.tr();
        paymentPwAndConfirmPaymentPwNotSameError = true;
      } else {
        paymentPwAndConfirmPaymentPwNotSameError = false;
      }
    }

    return errorMessageString;
  }

  @override
  Widget build(BuildContext context) {
    return BackGroundImage(
        currentUserType: CurrentUserType.centralized,
        child: Scaffold(
          backgroundColor: Colors.transparent,
          appBar: WalletAppBar(
            title: 'REGISTER.MAIN_TITLE'.tr(),
            backButton: true,
          ),
          body: TopCircularContainer(
              height: size.heightWithoutAppBar,
              child: SingleChildScrollView(
                child: Padding(
                    padding: EdgeInsets.only(left: 30, right: 30),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Padding(padding: EdgeInsets.only(bottom: 35)),
                        Form(
                          key: _pwForm,
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              _inputRow(
                                  'REGISTER.INPUT.EMAIL_TITLE'.tr(), _emailTc,
                                  hintText: 'REGISTER.INPUT.EMAIL_DESC'.tr(),
                                  vFnc: (v) {
                                if (isEmptyString(string: v))
                                  return 'REGISTER.INPUT.EMAIL_DESC'.tr();

                                if (!EmailValidator.validate(_emailTc.text))
                                  return 'REGISTER.VALIDATE_EMAIL'.tr();

                                return null;
                              }, obscureText: false),
                              Row(
                                  mainAxisAlignment:
                                      MainAxisAlignment.spaceBetween,
                                  children: [
                                    Container(
                                      width: 63,
                                      height: 1,
                                      child: Image(
                                        image: AssetImage(
                                          'images/registerLine.png',
                                          package: 'euruswallet',
                                        ),
                                      ),
                                    ),
                                    Container(
                                      width: 10,
                                      height: 14,
                                      child: Image(
                                        image: AssetImage(
                                          'images/lockIcon.png',
                                          package: 'euruswallet',
                                        ),
                                      ),
                                    ),
                                    Text('REGISTER.INPUT.LOGIN_PW_TITLE'.tr(),
                                        style: FXUI.normalTextStyle
                                            .copyWith(fontSize: 18)),
                                    Container(
                                        width: 63,
                                        height: 1,
                                        child: Image(
                                          image: AssetImage(
                                            'images/registerLine.png',
                                            package: 'euruswallet',
                                          ),
                                        )),
                                  ]),
                              _inputRow('REGISTER.INPUT.LOGIN_PW_TITLE'.tr(),
                                  _curLoginPwTc,
                                  hintText: hintEnterLoginPassword.tr(),
                                  vFnc: (v) {
                                return errorMessage(
                                    type: TextFieldType.loginPassword, v: v);
                              },
                                  suffixIcon: IconButton(
                                    onPressed: () {
                                      setState(() {
                                        isLoginPasswordMasked =
                                            !isLoginPasswordMasked;
                                      });
                                    },
                                    icon: Image.asset(
                                      isLoginPasswordMasked
                                          ? 'images/eyeClose.png'
                                          : 'images/eyeOpen.png',
                                      package: 'euruswallet',
                                      width: 16,
                                      height: 16,
                                      color: FXColor.mainDeepBlueColor,
                                    ),
                                  ),
                                  obscureText: !isLoginPasswordMasked),
                              _inputRow('REGISTER.INPUT.LOGIN_PW_DESC'.tr(),
                                  _confirmLoginPwTc,
                                  hintText: hintEnterLoginPassword.tr(),
                                  vFnc: (v) {
                                return errorMessage(
                                    type: TextFieldType.confirmLoginPassword,
                                    v: v);
                              },
                                  suffixIcon: IconButton(
                                    onPressed: () {
                                      setState(() {
                                        isConfirmLoginPasswordMasked =
                                            !isConfirmLoginPasswordMasked;
                                      });
                                    },
                                    icon: Image.asset(
                                      isConfirmLoginPasswordMasked
                                          ? 'images/eyeClose.png'
                                          : 'images/eyeOpen.png',
                                      package: 'euruswallet',
                                      width: 16,
                                      height: 16,
                                      color: FXColor.mainDeepBlueColor,
                                    ),
                                  ),
                                  obscureText: !isConfirmLoginPasswordMasked),
                              Padding(
                                  padding: const EdgeInsets.only(bottom: 10),
                                  child: Text(
                                      'REGISTER.DIGIT_AND_CHARACTER'.tr(args: [
                                        '・${'REGISTER.LOGIN_PW_SECTION_TITLE'.tr()}'
                                      ]),
                                      style: FXUI.normalTextStyle.copyWith(
                                          fontSize: 12,
                                          color:
                                              loginPwNumberAndAlphabeticalError
                                                  ? FXColor
                                                      .centralizedGrayTextColor
                                                  : FXColor.mainBlueColor))),
                              Padding(
                                padding: const EdgeInsets.only(bottom: 10),
                                child: Text(
                                    'REGISTER.CONTAIN_8_TO_20'.tr(args: [
                                      '・${'REGISTER.LOGIN_PW_SECTION_TITLE'.tr()}'
                                    ]),
                                    style: FXUI.normalTextStyle.copyWith(
                                        fontSize: 12,
                                        color: loginPwLengthError
                                            ? FXColor.centralizedGrayTextColor
                                            : FXColor.mainBlueColor)),
                              ),
                              Padding(
                                  padding: const EdgeInsets.only(bottom: 10),
                                  child: Text(
                                      'REGISTER.DIGIT_AND_CHARACTER'.tr(args: [
                                        '・${'REGISTER.INPUT.CONFIRM_LOGIN_PW_TITLE'.tr()}'
                                      ]),
                                      style: FXUI.normalTextStyle.copyWith(
                                          fontSize: 12,
                                          color:
                                              confirmLoginPwNumberAndAlphabeticalError
                                                  ? FXColor
                                                      .centralizedGrayTextColor
                                                  : FXColor.mainBlueColor))),
                              Padding(
                                padding: const EdgeInsets.only(bottom: 10),
                                child: Text(
                                    'REGISTER.CONTAIN_8_TO_20'.tr(args: [
                                      '・${'REGISTER.INPUT.CONFIRM_LOGIN_PW_TITLE'.tr()}'
                                    ]),
                                    style: FXUI.normalTextStyle.copyWith(
                                        fontSize: 12,
                                        color: confirmLoginPwLengthError
                                            ? FXColor.centralizedGrayTextColor
                                            : FXColor.mainBlueColor)),
                              ),
                              Padding(
                                  padding: const EdgeInsets.only(bottom: 10),
                                  child: Text(
                                      '・${'REGISTER.LOGIN_PW_CONFIRM'.tr()}',
                                      style: FXUI.normalTextStyle.copyWith(
                                          fontSize: 12,
                                          color:
                                              loginPwAndConfirmLoginPwNotSameError
                                                  ? FXColor
                                                      .centralizedGrayTextColor
                                                  : FXColor.mainBlueColor))),
                              Padding(
                                padding: const EdgeInsets.only(bottom: 10),
                                child: Text(
                                    'REGISTER.DIGIT_AND_CHARACTER'.tr(args: [
                                      '・${'REGISTER.INPUT.PAYMENT_PW_TITLE'.tr()}'
                                    ]),
                                    style: FXUI.normalTextStyle.copyWith(
                                        fontSize: 12,
                                        color:
                                            paymentPwNumberAndAlphabeticalError
                                                ? FXColor
                                                    .centralizedGrayTextColor
                                                : FXColor.mainBlueColor)),
                              ),
                              Padding(
                                padding: const EdgeInsets.only(bottom: 10),
                                child: Text(
                                    'REGISTER.CONTAIN_8_TO_20'.tr(args: [
                                      '・${'REGISTER.INPUT.PAYMENT_PW_TITLE'.tr()}'
                                    ]),
                                    style: FXUI.normalTextStyle.copyWith(
                                        fontSize: 12,
                                        color: paymentPwLengthError
                                            ? FXColor.centralizedGrayTextColor
                                            : FXColor.mainBlueColor)),
                              ),
                              Padding(
                                padding: const EdgeInsets.only(bottom: 10),
                                child: Text(
                                    'REGISTER.DIGIT_AND_CHARACTER'.tr(args: [
                                      '・${'REGISTER.INPUT.PAYMENT_PW_DESC'.tr()}'
                                    ]),
                                    style: FXUI.normalTextStyle.copyWith(
                                        fontSize: 12,
                                        color:
                                            confirmPaymentPwNumberAndAlphabeticalError
                                                ? FXColor
                                                    .centralizedGrayTextColor
                                                : FXColor.mainBlueColor)),
                              ),
                              Padding(
                                padding: const EdgeInsets.only(bottom: 10),
                                child: Text(
                                    'REGISTER.CONTAIN_8_TO_20'.tr(args: [
                                      '・${'REGISTER.INPUT.PAYMENT_PW_DESC'.tr()}'
                                    ]),
                                    style: FXUI.normalTextStyle.copyWith(
                                        fontSize: 12,
                                        color: confirmPaymentPwLengthError
                                            ? FXColor.centralizedGrayTextColor
                                            : FXColor.mainBlueColor)),
                              ),
                              Padding(
                                  padding: const EdgeInsets.only(bottom: 10),
                                  child: Text(
                                      '・${'REGISTER.PAYMENT_PW_CONFIRM'.tr()}',
                                      style: FXUI.normalTextStyle.copyWith(
                                          fontSize: 12,
                                          color:
                                              paymentPwAndConfirmPaymentPwNotSameError
                                                  ? FXColor
                                                      .centralizedGrayTextColor
                                                  : FXColor.mainBlueColor))),
                              Padding(
                                  padding: const EdgeInsets.only(bottom: 35),
                                  child: Text(
                                      '・${'REGISTER.LOGIN_PW_NOT_EQUAL_PAYMENT_PW'.tr()}',
                                      style: FXUI.normalTextStyle.copyWith(
                                          fontSize: 12,
                                          color: loginPwAndPaymentPwSameError
                                              ? FXColor.centralizedGrayTextColor
                                              : FXColor.mainBlueColor))),
                              Row(
                                  mainAxisAlignment:
                                      MainAxisAlignment.spaceBetween,
                                  children: [
                                    Container(
                                      width: 63,
                                      height: 1,
                                      child: Image(
                                        image: AssetImage(
                                          'images/registerLine.png',
                                          package: 'euruswallet',
                                        ),
                                      ),
                                    ),
                                    Container(
                                      width: 12,
                                      height: 13,
                                      child: Image(
                                        image: AssetImage(
                                          'images/registerExchangeIcon.png',
                                          package: 'euruswallet',
                                        ),
                                      ),
                                    ),
                                    Text("REGISTER.INPUT.PAYMENT_PW_TITLE".tr(),
                                        style: FXUI.normalTextStyle
                                            .copyWith(fontSize: 18)),
                                    Container(
                                        width: 63,
                                        height: 1,
                                        child: Image(
                                          image: AssetImage(
                                            'images/registerLine.png',
                                            package: 'euruswallet',
                                          ),
                                        )),
                                  ]),
                              _inputRow("REGISTER.INPUT.PAYMENT_PW_TITLE".tr(),
                                  _curPaymentPwTc,
                                  hintText: hintEnterYourTransactionCode.tr(),
                                  vFnc: (v) {
                                return errorMessage(
                                    type: TextFieldType.paymentPassword, v: v);
                              },
                                  suffixIcon: IconButton(
                                    onPressed: () {
                                      setState(() {
                                        isPaymentPasswordMasked =
                                            !isPaymentPasswordMasked;
                                      });
                                    },
                                    icon: Image.asset(
                                      isPaymentPasswordMasked
                                          ? 'images/eyeClose.png'
                                          : 'images/eyeOpen.png',
                                      package: 'euruswallet',
                                      width: 16,
                                      height: 16,
                                      color: FXColor.mainDeepBlueColor,
                                    ),
                                  ),
                                  obscureText: !isPaymentPasswordMasked),
                              _inputRow("REGISTER.INPUT.PAYMENT_PW_DESC".tr(),
                                  _confirmPaymentPwTc,
                                  hintText: hintEnterYourTransactionCode.tr(),
                                  vFnc: (v) {
                                return errorMessage(
                                    type: TextFieldType.confirmPaymentPassword,
                                    v: v);
                              },
                                  suffixIcon: IconButton(
                                    onPressed: () {
                                      setState(() {
                                        isConfirmPaymentPasswordMasked =
                                            !isConfirmPaymentPasswordMasked;
                                      });
                                    },
                                    icon: Image.asset(
                                      isConfirmPaymentPasswordMasked
                                          ? 'images/eyeClose.png'
                                          : 'images/eyeOpen.png',
                                      package: 'euruswallet',
                                      width: 16,
                                      height: 16,
                                      color: FXColor.mainDeepBlueColor,
                                    ),
                                  ),
                                  obscureText: !isConfirmPaymentPasswordMasked),
                            ],
                          ),
                        ),
                        Padding(
                          padding: EdgeInsets.only(top: 15),
                          child: Row(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Container(
                                width: 18,
                                height: 18,
                                child: Checkbox(
                                  value: _checkbox,
                                  onChanged: (value) {
                                    setState(() {
                                      _checkbox = !_checkbox;
                                    });
                                  },
                                ),
                              ),
                              Padding(
                                padding: EdgeInsets.only(left: 8),
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Padding(
                                      padding: EdgeInsets.only(top: 2),
                                      child: Text('REGISTER.AGREE_T&C'.tr(),
                                          style: FXUI.normalTextStyle.copyWith(
                                              color: FXColor
                                                  .centralizedGrayTextColor,
                                              fontSize: 12)),
                                    ),
                                    InkWell(
                                      child: Padding(
                                        padding:
                                            EdgeInsets.only(top: 5, bottom: 5),
                                        child: Text(
                                          'REGISTER.T&C'.tr(),
                                          style: FXUI.normalTextStyle.copyWith(
                                            fontSize: 12,
                                            fontWeight: FontWeight.w400,
                                            decoration:
                                                TextDecoration.underline,
                                          ),
                                        ),
                                      ),
                                      onTap: () => launch(
                                        'SETTING_PAGE.HELP_CENTER_CARD.TNC.LINK'
                                            .tr(),
                                      ),
                                    ),
                                    if (!_checkbox)
                                      Text('REGISTER.MUST_AGREE_T&C'.tr(),
                                          style: FXUI.normalTextStyle.copyWith(
                                              color: FXColor.alertRedColor,
                                              fontFamily: 'Roboto',
                                              fontSize: 12)),
                                  ],
                                ),
                              ),
                            ],
                          ),
                        ),
                        Container(
                          padding: EdgeInsets.only(top: 25, bottom: 25),
                          child: SizedBox(
                              width: double.infinity,
                              child: SubmitButton(
                                  btnController: btnController,
                                  buttonBGColor: FXColor.mainBlueColor,
                                  label: 'REGISTER.GET_CODE_FROM_EMAIL'.tr(),
                                  onPressed: () async {
                                    if (_pwForm.currentState != null &&
                                        _pwForm.currentState!.validate() &&
                                        _checkbox) {
                                      setState(() {});
                                      common.registerByEmail =
                                          await api.registerByEmail(
                                              email:
                                                  _emailTc.text.toLowerCase(),
                                              password: _curLoginPwTc.text);
                                      if (await common.checkApiError(
                                          context: context,
                                          errorString:
                                              common.registerByEmail?.message,
                                          returnCode: common
                                              .registerByEmail?.returnCode,
                                          btnColor: FXColor.mainBlueColor)) {
                                        common.email =
                                            _emailTc.text.toLowerCase();
                                        common.loginPassword =
                                            _curLoginPwTc.text;
                                        common.paymentPassword =
                                            _curPaymentPwTc.text;
                                        common.verifyCodePageType =
                                            VerifyCodePageType.register;
                                        common.pushPage(
                                            page: VerifyCodePage(),
                                            context: context);
                                      } else {}
                                    }
                                    btnController.reset();
                                  })),
                        ),
                      ],
                    )),
              )),
        ));
  }

  Widget _inputRow(
    String title,
    TextEditingController tc, {
    String? hintText,
    String? errorMsg,
    String? Function(String?)? vFnc,
    Widget? suffixIcon,
    bool? obscureText,
  }) {
    final _defaultTextFieldInputDecoration = InputDecoration(
      filled: true,
      fillColor: FXColor.lightGreyTextColor,
      hintStyle: Theme.of(context)
          .textTheme
          .subtitle1
          ?.apply(color: Theme.of(context).hintColor),
      border: OutlineInputBorder(
        borderSide: BorderSide.none,
        borderRadius: FXUI.cricleRadius,
      ),
      contentPadding: EdgeInsets.all(16),
    );

    // double _hPadding = MediaQuery.of(context).size.width / 13;

    return Container(
      width: double.infinity,
      child: Column(
        children: [
          Container(
            padding: getEdgeInsetsSymmetric(),
            alignment: Alignment(-1, 0),
            child: Text(
              title,
              style: Theme.of(context)
                  .textTheme
                  .bodyText2
                  ?.apply(color: FXColor.lightGray),
            ),
          ),
          TextFormField(
            decoration: _defaultTextFieldInputDecoration.copyWith(
              hintText: hintText ?? '',
              hintStyle: FXUI.normalTextStyle.copyWith(
                color: FXColor.centralizedGrayTextColor,
                fontSize: 14,
              ),
              errorText: errorMsg ?? '',
              errorMaxLines: 2,
              suffixIcon: suffixIcon,
            ),
            controller: tc,
            obscureText: obscureText ?? true,
            autovalidateMode: AutovalidateMode.onUserInteraction,
            onChanged: (v) => _clearErrorMsg(tc),
            validator: vFnc,
          )
        ],
      ),
    );
  }

  void _clearErrorMsg(TextEditingController v) {
    if (_pwForm.currentState != null)
      setState(() {
        _pwForm.currentState!.validate();
      });
  }
}
