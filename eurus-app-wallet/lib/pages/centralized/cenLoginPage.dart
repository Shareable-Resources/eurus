import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/pages/centralized/centralized_wallet_base_page.dart';
import 'package:euruswallet/pages/create_wallet_page.dart';

import 'cenForgetLoginPwPage.dart';

class CenLoginPage extends StatefulWidget {
  CenLoginPage({Key? key}) : super(key: key);

  @override
  _CenLoginPageState createState() => _CenLoginPageState();
}

class _CenLoginPageState extends State<CenLoginPage> {
  TextEditingController emailTc = TextEditingController();
  TextEditingController pwTc = TextEditingController();
  bool isPasswordMasked = false;
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    super.initState();
    //testing code
    if (AutoFillAccount) {
      emailTc.text =
          envType == EnvType.Testnet || envType == EnvType.Production ? "eu31@18m.dev" :  "ken575@goldhub.hk";
      pwTc.text = 'aaaaaaa1';
    }
  }

  @override
  Widget build(BuildContext context) {
    return BackGroundImage(
      currentUserType: CurrentUserType.centralized,
      child: Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: Colors.transparent,
        body: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            AppBar(
              centerTitle: true,
              title: Text('CEN_LOGIN.TITLE'.tr()),
              backgroundColor: Colors.transparent,
              elevation: 0,
            ),
            Expanded(
              flex: 1,
              child: Container(
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(15),
                      topRight: Radius.circular(15),
                    ),
                  ),
                  margin: EdgeInsets.only(top: 12),
                  padding: EdgeInsets.symmetric(vertical: 31, horizontal: 35),
                  child: SingleChildScrollView(
                      child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      _inputRow('CEN_LOGIN.INPUT.EMAIL_TITLE'.tr(), emailTc,
                          hintText: 'CEN_LOGIN.INPUT.EMAIL_DESC'.tr()),
                      _inputRow(
                        'CEN_LOGIN.INPUT.PW_TITLE'.tr(),
                        pwTc,
                        hintText: 'CEN_LOGIN.INPUT.PW_DESC'.tr(),
                        obscureText: !isPasswordMasked,
                        suffixIcon: IconButton(
                          onPressed: () {
                            setState(() {
                              isPasswordMasked = !isPasswordMasked;
                            });
                          },
                          icon: Image.asset(
                            isPasswordMasked
                                ? 'images/eyeClose.png'
                                : 'images/eyeOpen.png',
                            package: 'euruswallet',
                            width: 16,
                            height: 16,
                            color: FXColor.mainBlueColor,
                          ),
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          print('Forgot pw btn onpress');
                          common.pushPage(
                              page: CenForgetLoginPwPage(), context: context);
                        },
                        child: Text(
                          'CEN_LOGIN.FORGOT_PW'.tr(),
                          style: FXUI.normalTextStyle.copyWith(
                              color: FXColor.mainBlueColor, fontSize: 14),
                        ),
                      ),
                      SizedBox(
                        height: 35,
                        width: double.infinity,
                      ),
                      SizedBox(
                        width: double.infinity,
                        child: SubmitButton(
                          btnController: btnController,
                          loadingSecond: 3,
                          buttonBGColor: FXColor.mainBlueColor,
                          label: 'CEN_LOGIN.LOGIN_BTN'.tr(),
                          onPressed: onLoginBtnPress,
                        ),
                      ),
                      SizedBox(
                        height: 51,
                        width: double.infinity,
                      ),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text('CEN_LOGIN.REG_TEXT'.tr(),
                              style: FXUI.normalTextStyle
                                  .copyWith(color: FXColor.textGray)),
                          TextButton(
                            onPressed: () {
                              Navigator.push(
                                context,
                                MaterialPageRoute(
                                  builder: (context) =>
                                      CentralizedWalletBasePage(
                                    appBarTitle: Text(
                                        'CREATE_WALLET_PAGE.CREATE_WALLET'
                                            .tr()),
                                    body: CreateWalletPage(),
                                  ),
                                ),
                              );
                            },
                            child: Text('CEN_LOGIN.REG_BTN'.tr()),
                          )
                        ],
                      ),
                    ],
                  ))),
            )
          ],
        ),
      ),
    );
  }

  void onLoginBtnPress() async {
    common.email = emailTc.text.toLowerCase();
    common.loginPassword = pwTc.text;

    await common.successMoveToHomePage(
        userType: CurrentUserType.centralized,
        email: common.email ?? '',
        context: context,
        loginPassword: common.loginPassword ?? '');
    btnController.reset();
  }

  Widget _inputRow(
    String title,
    TextEditingController tc, {
    String? hintText,
    String? errorMsg,
    String? Function(String?)? vFnc,
    bool? obscureText,
    Widget? suffixIcon,
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
            obscureText: obscureText ?? false,
            autovalidateMode: AutovalidateMode.onUserInteraction,
            // onChanged: (v) => _clearErrorMsg(),
            validator: vFnc,
          )
        ],
      ),
    );
  }
}
