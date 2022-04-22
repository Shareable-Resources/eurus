import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/changeLoginPW.dart';
import 'package:euruswallet/model/codeVerification.dart';
import 'package:euruswallet/pages/settingSubpages/cardContainer.dart';
import 'package:euruswallet/pages/settingSubpages/settingAppBar.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';

class ChangePaymentPWPage extends StatefulWidget {
  ChangePaymentPWPage();

  @override
  _ChangePaymentPWPageState createState() => _ChangePaymentPWPageState();
}

class _ChangePaymentPWPageState extends State<ChangePaymentPWPage> {
  final TextEditingController _curPwTc = TextEditingController();
  final TextEditingController _newPwTc = TextEditingController();
  final TextEditingController _confirmPwTc = TextEditingController();

  bool _isCurPwMasked = false;
  bool _isNewPwMasked = false;
  bool _isConfirmPwMasked = false;

  final _pwForm = GlobalKey<FormState>();

  bool? _curPWValid;
  bool? _newPWValid;
  bool? _confirmPWValid;

  String? _curPWErrorMsg;
  String? _newPWErrorMsg;
  String? _confirmPWErrorMsg;
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    super.initState();
    //testing code
    // if (AutoFillAccount) {
    //   _curPwTc.text = "bbbbbbb1";
    //   _newPwTc.text = "bbbbbbb1";
    //   _confirmPwTc.text = "bbbbbbb1";
    // }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SettingAppBar(true),
      body: Container(
        child: SingleChildScrollView(
          child: SafeArea(
            child: CardContainer(
              '',
              Container(
                padding: EdgeInsets.symmetric(vertical: 9),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Form(
                      key: _pwForm,
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          if (common.changePaymentPasswordType ==
                              ChangePaymentPasswordType.changePaymentPassword)
                            _inputRow(
                              'CHANGE_LOCKER_PW.CURRENT_PW.LABEL'.tr(),
                              _curPwTc,
                              hintText:
                                  'CHANGE_LOCKER_PW.CURRENT_PW.PLACEHOLDER'
                                      .tr(),
                              vFnc: (v) {
                                if (isEmptyString(string: v))
                                  return 'CHANGE_LOCKER_PW.ERROR.EMPTY_CUR_PW'
                                      .tr();

                                if (common.isNotContainDigitAndCharacter(
                                    text: v ?? '')) {
                                  return 'CHANGE_LOCKER_PW.ERROR.PW_NOT_DIGIT_OR_CHARACTER'
                                      .tr();
                                }

                                if (common.isNotContain8To20Characters(
                                    text: v ?? '')) {
                                  return 'CHANGE_LOCKER_PW.ERROR.PW_NOT_CONTAIN_8_TO_20_CHARACTER'
                                      .tr();
                                }

                                if (_curPWValid == false) return _curPWErrorMsg;

                                return null;
                              },
                              obscureText: !_isCurPwMasked,
                              suffixIconOnPressed: () {
                                setState(() {
                                  _isCurPwMasked = !_isCurPwMasked;
                                });
                              },
                            ),
                          _inputRow(
                            'CHANGE_LOCKER_PW.NEW_PW.LABEL'.tr(),
                            _newPwTc,
                            hintText:
                                'CHANGE_LOCKER_PW.NEW_PW.PLACEHOLDER'.tr(),
                            vFnc: (v) {
                              if (isEmptyString(string: v))
                                return 'CHANGE_LOCKER_PW.ERROR.EMPTY_NEW_PW'
                                    .tr();

                              if (common.isNotContainDigitAndCharacter(
                                  text: v ?? '')) {
                                return 'CHANGE_LOCKER_PW.ERROR.PW_NOT_DIGIT_OR_CHARACTER'
                                    .tr();
                              }

                              if (common.isNotContain8To20Characters(
                                  text: v ?? '')) {
                                return 'CHANGE_LOCKER_PW.ERROR.PW_NOT_CONTAIN_8_TO_20_CHARACTER'
                                    .tr();
                              }

                              int newPwMatch =
                                  _newPwTc.text.compareTo(_confirmPwTc.text);
                              if (newPwMatch != 0) {
                                _confirmPWErrorMsg =
                                    'COMMON_ERROR.PW_INCONSISTENT'.tr();
                                return _confirmPWErrorMsg;
                              }

                              if (_newPWValid == false) return _newPWErrorMsg;

                              return null;
                            },
                            obscureText: !_isNewPwMasked,
                            suffixIconOnPressed: () {
                              setState(() {
                                _isNewPwMasked = !_isNewPwMasked;
                              });
                            },
                          ),
                          _inputRow(
                            'CHANGE_LOCKER_PW.CONFIRM_PW.LABEL'.tr(),
                            _confirmPwTc,
                            hintText:
                                'CHANGE_LOCKER_PW.CONFIRM_PW.PLACEHOLDER'.tr(),
                            vFnc: (v) {
                              if (isEmptyString(string: v))
                                return 'CHANGE_LOCKER_PW.ERROR.EMPTY_CONFIRM_PW'
                                    .tr();

                              if (common.isNotContainDigitAndCharacter(
                                  text: v ?? '')) {
                                return 'CHANGE_LOCKER_PW.ERROR.PW_NOT_DIGIT_OR_CHARACTER'
                                    .tr();
                              }

                              if (common.isNotContain8To20Characters(
                                  text: v ?? '')) {
                                return 'CHANGE_LOCKER_PW.ERROR.PW_NOT_CONTAIN_8_TO_20_CHARACTER'
                                    .tr();
                              }

                              int newPwMatch =
                                  _newPwTc.text.compareTo(_confirmPwTc.text);
                              if (newPwMatch != 0) {
                                _confirmPWErrorMsg =
                                    'COMMON_ERROR.PW_INCONSISTENT'.tr();
                                return _confirmPWErrorMsg;
                              }

                              if (_confirmPWValid == false)
                                return _confirmPWErrorMsg;

                              return null;
                            },
                            obscureText: !_isConfirmPwMasked,
                            suffixIconOnPressed: () {
                              setState(() {
                                _isConfirmPwMasked = !_isConfirmPwMasked;
                              });
                            },
                          ),
                        ],
                      ),
                    ),
                    Container(
                        padding: EdgeInsets.only(
                            left: 35, right: 35, top: 25, bottom: 25),
                        child: SubmitButton(
                            btnController: btnController,
                            label: 'COMMON.CONFIRM'.tr(),
                            onPressed: _changePW)),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _inputRow(
    String title,
    TextEditingController tc, {
    String? hintText,
    String? errorMsg,
    String? Function(String?)? vFnc,
    bool obscureText = true,
    void Function()? suffixIconOnPressed,
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

    double _hPadding = MediaQuery.of(context).size.width / 13;

    return Container(
      width: double.infinity,
      padding:
          EdgeInsets.symmetric(horizontal: _hPadding > 35 ? 35 : _hPadding),
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
                fontSize: 12,
              ),
              errorText: errorMsg ?? '',
              suffixIcon: IconButton(
                onPressed: suffixIconOnPressed,
                icon: Image.asset(
                  !obscureText ? 'images/eyeClose.png' : 'images/eyeOpen.png',
                  package: 'euruswallet',
                  width: 16,
                  height: 16,
                  color: common.getBackGroundColor(),
                ),
              ),
            ),
            controller: tc,
            obscureText: obscureText,
            autovalidateMode: AutovalidateMode.onUserInteraction,
            onChanged: (v) => _clearErrorMsg(),
            validator: vFnc,
          )
        ],
      ),
    );
  }

  Future<void> _changePW() async {
    if (_pwForm.currentState?.validate() ?? false) {
      // EasyLoading.show(status: 'COMMON.LOADING_W_DOT'.tr());
      if (common.changePaymentPasswordType ==
          ChangePaymentPasswordType.changePaymentPassword) {
        final result = await api.changePaymentPassword(
            context: context,
            oldPaymentPassword: _curPwTc.text,
            newPaymentPassword: _newPwTc.text);
        if (result) {
          _successDialog();
        }
      }

      if (common.changePaymentPasswordType ==
          ChangePaymentPasswordType.forgetPaymentPassword) {
        CodeVerification codeVerification = await api.resetPaymentPassword(
            email: common.email,
            password: _newPwTc.text,
            mnemonic: common.serverMnemonic);
        if (await common.checkApiError(
            context: context,
            errorString: codeVerification.message,
            returnCode: codeVerification.returnCode)) {
          _successDialog();
        }
      }

      EasyLoading.dismiss();
    }
    btnController.reset();
  }

  void _clearErrorMsg() {
    setState(() {
      _curPWValid = null;
      _newPWValid = null;
      _confirmPWValid = null;
      _curPWErrorMsg = null;
      _newPWErrorMsg = null;
      _confirmPWErrorMsg = null;
    });
  }

  Future<void> _successDialog() async {
    double _hPadding = MediaQuery.of(context).size.width / 13;
    double _wPadding = MediaQuery.of(context).size.height / 30;

    return showDialog(
      context: context,
      builder: (_) {
        return Dialog(
          backgroundColor: Colors.transparent,
          insetPadding: EdgeInsets.all(10),
          child: Container(
            width: double.infinity,
            padding: EdgeInsets.symmetric(
                vertical: _hPadding > 25 ? 25 : _hPadding,
                horizontal: _wPadding > 17 ? 17 : _wPadding),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: FXUI.cricleRadius,
            ),
            child: Stack(
              children: [
                SingleChildScrollView(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Center(
                        child: Text(
                          'CHANGE_LOCKER_PW.CHANGE_SUCCESS'.tr(),
                          style: FXUI.normalTextStyle.copyWith(
                            fontWeight: FontWeight.bold,
                            fontSize: 18,
                          ),
                        ),
                      ),
                      Padding(padding: EdgeInsets.only(top: 70)),
                      SizedBox(
                        width: MediaQuery.of(context).size.width / 4.5,
                        child: Image(
                          image: AssetImage(
                              "images/${!isCentralized() ? 'decenTickIcon' : 'tickIcon'}.png",
                              package: 'euruswallet'),
                        ),
                      ),
                      Padding(
                        padding: EdgeInsets.symmetric(vertical: 15),
                        child: Text(
                          'COMMON.SUCCESS'.tr(),
                          style: FXUI.normalTextStyle.copyWith(
                            fontWeight: FontWeight.bold,
                            fontSize: 23,
                          ),
                        ),
                      ),
                      Padding(
                        padding: EdgeInsets.symmetric(horizontal: 10),
                        child: Text(
                          'CHANGE_LOCKER_PW.CHANGE_SUCCESS_PW'.tr(),
                          textAlign: TextAlign.center,
                          style: FXUI.normalTextStyle
                              .copyWith(color: FXColor.textGray),
                        ),
                      ),
                      Padding(
                        padding: EdgeInsets.only(top: 50),
                        child: TextButton(
                          style: TextButton.styleFrom(padding: EdgeInsets.zero),
                          onPressed: () {
                            Navigator.pop(_);
                            common.logoutFunction?.add(true);
                          },
                          child: Container(
                            decoration: BoxDecoration(
                              color: common.getBackGroundColor(),
                              borderRadius: FXUI.cricleRadius,
                            ),
                            padding: EdgeInsets.all(15),
                            width: double.infinity,
                            child: Center(
                              child: Text(
                                // 'Back',
                                'COMMON.OK'.tr(),
                                style: FXUI.normalTextStyle.copyWith(
                                    color: Colors.white, fontSize: 16),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                Positioned(
                  right: -13,
                  top: -12,
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
