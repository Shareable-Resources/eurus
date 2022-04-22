import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/userKYCStatus.dart';
import './ReportAcSubmitPage.dart';

class ReportAcForm extends StatefulWidget {
  const ReportAcForm({
    Key? key,
  }) : super(key: key);

  @override
  _ReportAcFormState createState() => _ReportAcFormState();
}

class _ReportAcFormState extends State<ReportAcForm> {
  @override
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();
  final TextEditingController _emailTextController = TextEditingController();
  final TextEditingController _securityTextController = TextEditingController();
  static const double padding = 10.0;
  final List<String> questionsList = [
    'What is the name of the high school you attended?',
    'What is your favorite food?',
    'What is the name of your first pet?'
  ];
  String dropdownValue = 'What is the name of the high school you attended?';

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    return BackGroundImage(
        currentUserType: CurrentUserType.centralized,
        child: Scaffold(
            backgroundColor: Colors.transparent,
            appBar: WalletAppBar(
              title: 'REPORT_AC.REPORT_ACCOUNT'.tr(),
              backButton: true,
            ),
            body: SingleChildScrollView(
              child: TopCircularContainer(
                  height: size.heightWithoutAppBar,
                  child: Padding(
                    padding: getEdgeInsetsSymmetric(horizontal: 34),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding:
                              const EdgeInsets.symmetric(vertical: padding),
                          child: Text('REPORT_AC.CASE_ID'.tr(),
                              style: FXUI.normalTextStyle
                                  .copyWith(color: FXColor.greyTextColor)),
                        ),
                        Text('A0000', style: FXUI.normalTextStyle.copyWith()),
                        Padding(
                          padding:
                              const EdgeInsets.symmetric(vertical: padding),
                          child: Text('CEN_LOGIN.INPUT.EMAIL_TITLE'.tr(),
                              style: FXUI.normalTextStyle
                                  .copyWith(color: FXColor.greyTextColor)),
                        ),
                        TextField(
                          controller: _emailTextController,
                          decoration:
                              FXUI.defaultTextFieldInputDecoration.copyWith(),
                          maxLines: 1,
                          onChanged: (String value) => {},
                        ),
                        Padding(
                          padding:
                              const EdgeInsets.symmetric(vertical: padding),
                          child: Text('KYC.SECURITY_QUESTION'.tr(),
                              style: FXUI.normalTextStyle
                                  .copyWith(color: FXColor.greyTextColor)),
                        ),
                        Container(
                            padding: EdgeInsets.symmetric(
                                vertical: 2.0, horizontal: 14.0),
                            decoration: FXUI.boxDecorationWithShadow,
                            child: DropdownButton<String>(
                              value: dropdownValue,
                              icon: const Icon(Icons.arrow_drop_down),
                              iconSize: 24,
                              isExpanded: true,
                              underline: Container(),
                              style: FXUI.normalTextStyle,
                              onChanged: (String? newValue) {
                                setState(() {
                                  dropdownValue = newValue!;
                                });
                              },
                              items: questionsList
                                  .map<DropdownMenuItem<String>>(
                                      (String value) {
                                return DropdownMenuItem<String>(
                                  value: value,
                                  child: Text(value),
                                );
                              }).toList(),
                            )),
                        Padding(padding: EdgeInsets.all(padding)),
                        TextField(
                          controller: _securityTextController,
                          decoration:
                              FXUI.defaultTextFieldInputDecoration.copyWith(),
                          maxLines: 1,
                          onChanged: (String value) => {},
                        ),
                        Padding(padding: EdgeInsets.all(padding)),
                        SubmitButton(
                            btnController: btnController,
                            label: 'COMMON.NEXT_STEP'.tr(),
                            onPressed: () {
                              showDialog(
                                  context: context,
                                  builder: (BuildContext context) {
                                    return CustomConfirmDialog(
                                      title: "COMMON.SUCCESS".tr(),
                                      icon: Image(
                                        width: 86,
                                        height: 86,
                                        image: AssetImage(
                                          'images/uploadSuccessIcon.png',
                                          package: 'euruswallet',
                                        ),
                                      ),
                                      descriptions:
                                          "REPORT_AC.SUCCESS_DEACTIVATED_AC"
                                              .tr(),
                                      useSubmitButton: false,
                                      buttonText: "COMMON.NEXT_STEP".tr(),
                                      btnHandler: () {
                                        common.pushReplacementPage(
                                            page: ReportAcSubmitPage(),
                                            context: context);
                                        btnController.reset();
                                      },
                                    );
                                  });
                            }),
                        Padding(padding: EdgeInsets.all(padding)),
                        SubmitButton(
                            btnController: btnController,
                            buttonBGColor: FXColor.cancelGrayButton,
                            label: 'REPORT_AC.FORGOT_SECURITY_QUESTIONS'.tr(),
                            onPressed: () {
                              showDialog(
                                  context: context,
                                  builder: (BuildContext context) {
                                    return CustomConfirmDialog(
                                      title: "COMMON.REMIND".tr(),
                                      icon: Image(
                                        width: 88,
                                        height: 88,
                                        image: AssetImage(
                                          'images/warningIcon.png',
                                          package: 'euruswallet',
                                        ),
                                      ),
                                      descriptions:
                                          "REPORT_AC.FORGOT_SECURITY_QUESTIONS_REMIND"
                                              .tr(),
                                      useSubmitButton: false,
                                      buttonText: "COMMON.NEXT_STEP".tr(),
                                      btnHandler: () {
                                        common.pushPage(
                                            page: ReportAcSubmitPage(),
                                            context: context);
                                        btnController.reset();
                                      },
                                    );
                                  });
                            })
                      ],
                    ),
                  )),
            )));
  }
}
