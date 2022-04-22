import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/userKYCStatus.dart';

class SecurityQuestionForm extends StatefulWidget {
  const SecurityQuestionForm({
    Key? key,
  }) : super(key: key);

  @override
  _SecurityQuestionFormState createState() => _SecurityQuestionFormState();
}

class _SecurityQuestionFormState extends State<SecurityQuestionForm> {
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();
  final TextEditingController _question1Controller = TextEditingController();
  final TextEditingController _question2Controller = TextEditingController();
  final TextEditingController _question3Controller = TextEditingController();
  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    return BackGroundImage(
        currentUserType: CurrentUserType.centralized,
        child: Scaffold(
            backgroundColor: Colors.transparent,
            appBar: WalletAppBar(
              title: 'KYC.SECURITY_QUESTION'.tr(),
              backButton: true,
            ),
            body: SingleChildScrollView(
              child: TopCircularContainer(
                  height: size.heightWithoutAppBar,
                  child: Padding(
                    padding: getEdgeInsetsSymmetric(horizontal: 34),
                    child: Column(
                      children: [
                        Padding(padding: EdgeInsets.only(top: 40.0)),
                        genQuestion(
                            'What is the name of the high school you attended?',
                            _question1Controller),
                        genQuestion('What is your favorite food?',
                            _question2Controller),
                        genQuestion('What is the name of your first pet?',
                            _question3Controller),
                        SubmitButton(
                            btnController: btnController,
                            label: 'COMMON.SUBMIT'.tr(),
                            onPressed: () {
                              print(
                                  'answer1: ${_question1Controller.text}, answer2: ${_question2Controller.text}, answer3: ${_question3Controller.text}');
                              btnController.reset();
                            })
                      ],
                    ),
                  )),
            )));
  }

  Widget genQuestion(String question, TextEditingController controller) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(question, style: FXUI.normalTextStyle.copyWith()),
        Padding(
          padding: const EdgeInsets.only(top: 14.0, bottom: 50.0),
          child: TextField(
            controller: controller,
            decoration: FXUI.defaultTextFieldInputDecoration.copyWith(),
            maxLines: 1,
            onChanged: (String value) => {},
          ),
        ),
      ],
    );
  }
}
