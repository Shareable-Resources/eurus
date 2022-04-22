import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/pages/centralized/verifyCode.dart';
import './reportAcForm.dart';

class CenForgetLoginPwPage extends StatefulWidget {
  CenForgetLoginPwPage({Key? key}) : super(key: key);

  @override
  _CenForgetLoginPwPageState createState() => _CenForgetLoginPwPageState();
}

class _CenForgetLoginPwPageState extends State<CenForgetLoginPwPage> {
  TextEditingController emailTc = TextEditingController();
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    super.initState();

    //testing code
    if (AutoFillAccount) {
      emailTc.text =
      envType == EnvType.Testnet || envType == EnvType.Production ? "eu31@18m.dev" : "ken575@goldhub.hk";
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          image: DecorationImage(
            image: AssetImage(
              'images/backgroundImage.png',
              package: 'euruswallet',
            ),
            fit: BoxFit.cover,
            alignment: Alignment.topCenter,
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            AppBar(
              centerTitle: true,
              title: Text('FORGOT_PW.FORGOT_PW_BTN'.tr()),
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
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _inputRow('CEN_LOGIN.INPUT.EMAIL_TITLE'.tr(), emailTc,
                        hintText: 'CEN_LOGIN.INPUT.EMAIL_DESC'.tr()),
                    SizedBox(
                        width: double.infinity,
                        child: SubmitButton(
                            btnController: btnController,
                            buttonBGColor: FXColor.mainBlueColor,
                            label: 'FORGOT_PW.FORGOT_PW_BTN'.tr(),
                            onPressed: onForgetLoginPwBtnPress)),
                    Padding(
                      padding: const EdgeInsets.symmetric(vertical: 12.0),
                      child: const Divider(
                        height: 20,
                        thickness: 1,
                      ),
                    ),
                    /* Hide for mainnet */
                    // CustomTextButton(
                    //   buttonBGColor: FXColor.warningRedColor,
                    //   text: 'REPORT_AC.REPORT_ACCOUNT'.tr(),
                    //   onPressed: () {
                    //     showDialog(
                    //         context: context,
                    //         builder: (BuildContext context) {
                    //           return CustomConfirmDialog(
                    //             title: "COMMON.WARNING".tr(),
                    //             icon: Image(
                    //               width: 88,
                    //               height: 88,
                    //               image: AssetImage(
                    //                 'images/warningIcon.png',
                    //                 package: 'euruswallet',
                    //               ),
                    //             ),
                    //             descriptions:
                    //                 "REPORT_AC.WARNING_REPORT_AC".tr(),
                    //             buttonText: "COMMON.OK".tr(),
                    //             btnColor: FXColor.warningRedColor,
                    //             showCancelButton: true,
                    //             btnHandler: () {
                    //               common.pushPage(
                    //                   page: ReportAcForm(), context: context);
                    //             },
                    //           );
                    //         });
                    //   },
                    // ),
                  ],
                ),
              ),
            )
          ],
        ),
      ),
    );
  }

  void onForgetLoginPwBtnPress() async {
    common.forgetLoginPw =
        await api.forgetLoginPw(email: emailTc.text.toLowerCase());

    if (await common.checkApiError(
        context: context,
        errorString: common.forgetLoginPw?.message,
        returnCode: common.forgetLoginPw?.returnCode)) {
      common.verifyCodePageType = VerifyCodePageType.forgotLoginPw;
      common.forgetPwEmail = emailTc.text.toLowerCase();
      common.pushPage(page: VerifyCodePage(), context: context);
    }
    btnController.reset();
  }

  Widget _inputRow(
    String title,
    TextEditingController tc, {
    String? hintText,
    String? errorMsg,
    String? Function(String?)? vFnc,
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
                errorMaxLines: 2),
            controller: tc,
            obscureText: obscureText ?? false,
            autovalidateMode: AutovalidateMode.onUserInteraction,
            // onChanged: (v) => _clearErrorMsg(),
            validator: vFnc ?? (v) => null,
          )
        ],
      ),
    );
  }
}
