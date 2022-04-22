import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/userKYCStatus.dart';
import 'package:euruswallet/commonUI/kycCamera.dart';
import 'package:livechat_inc/livechat_inc.dart';

class ReportAcSubmitPage extends StatefulWidget {
  const ReportAcSubmitPage({
    Key? key,
  }) : super(key: key);

  @override
  _ReportAcSubmitPagState createState() => _ReportAcSubmitPagState();
}

class _ReportAcSubmitPagState extends State<ReportAcSubmitPage> {
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();
  static const double padding = 10.0;
  String? imagePath;
  String caseId = "A00000001";

  void onTakedPhoto() async {
    openKYCCamera(
        context,
        (String path) => setState(() {
              imagePath = path;
              Navigator.pop(context);
            }),
        cameraSide: CameraSide.front);
  }

  Future<Widget?> _uploadSuccessDialog() async {
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
                child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Center(
                        child: Text(
                          'KYC.UPLOAD_SUCCESS_DIALOG'.tr(),
                          style: FXUI.normalTextStyle.copyWith(
                            fontWeight: FontWeight.bold,
                            fontSize: 18,
                          ),
                        ),
                      ),
                      Padding(padding: EdgeInsets.only(top: 25)),
                      Text('REPORT_AC.CASE_ID'.tr(),
                          style: FXUI.normalTextStyle
                              .copyWith(color: FXColor.greyTextColor)),
                      Padding(padding: EdgeInsets.only(top: 6)),
                      Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(caseId,
                                style: FXUI.normalTextStyle.copyWith()),
                            InkWell(
                                onTap: () {
                                  Clipboard.setData(
                                      ClipboardData(text: caseId));
                                  common.showCopiedToClipboardSnackBar(
                                    context,
                                    margin: EdgeInsets.only(bottom: 56),
                                    behavior: SnackBarBehavior.floating,
                                  );
                                },
                                child: Padding(
                                    padding: const EdgeInsets.only(left: 8),
                                    child: Image.asset('images/copyIcon.png',
                                        package: 'euruswallet',
                                        width: 19,
                                        height: 19,
                                        fit: BoxFit.contain)))
                          ]),
                      Padding(padding: EdgeInsets.only(top: 15)),
                      Image(
                        width: 86,
                        height: 86,
                        image: AssetImage(
                          'images/uploadSuccessIcon.png',
                          package: 'euruswallet',
                        ),
                      ),
                      Padding(padding: EdgeInsets.only(top: 15)),
                      Text('COMMON.PENDING'.tr(),
                          style: FXUI.normalTextStyle.copyWith(fontSize: 20.0)),
                      Padding(padding: EdgeInsets.only(top: 15)),
                      Text(
                        'REPORT_AC.UPLOAD_IMAGE_SUCCESS'.tr(),
                        style: FXUI.normalTextStyle
                            .copyWith(color: FXColor.textGray, fontSize: 14),
                      ),
                      Padding(padding: EdgeInsets.only(top: 25)),
                      CustomTextButton(
                          text: "SETTING_PAGE.HELP_CENTER_CARD.CS".tr(),
                          onPressed: () {
                            Navigator.of(context)
                                .popUntil((route) => route.isFirst);
                            // LivechatInc.start_chat(
                            //   "12610959",
                            //   "",
                            //   web3dart.myEthereumAddress.toString(),
                            //   "guest@gmail.com",
                            // );
                          }),
                    ]),
              ));
        });
  }

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
                      children: [
                        Padding(
                          padding:
                              const EdgeInsets.symmetric(vertical: padding),
                          child: Text('REPORT_AC.CASE_ID'.tr(),
                              style: FXUI.normalTextStyle
                                  .copyWith(color: FXColor.greyTextColor)),
                        ),
                        Text('A0000', style: FXUI.normalTextStyle.copyWith()),
                        Padding(padding: EdgeInsets.all(padding)),
                        Container(
                          padding: getEdgeInsetsSymmetric(horizontal: 12),
                          decoration: FXUI.boxDecorationWithShadow,
                          child: Text('181818',
                              style: FXUI.normalTextStyle
                                  .copyWith(letterSpacing: 18.0)),
                        ),
                        Padding(padding: EdgeInsets.all(padding)),
                        Text('REPORT_AC.SELFIE_INSTRUCTION'.tr(),
                            style:
                                FXUI.normalTextStyle.copyWith(fontSize: 13.0)),
                        imagePath == null
                            ? Container(
                                margin: EdgeInsets.only(top: 20.0),
                                child: GestureDetector(
                                  onTap: () {
                                    onTakedPhoto();
                                  },
                                  child: CircleAvatar(
                                    backgroundColor: FXColor.lightGreyColor2,
                                    child: Image.asset('images/cameraIcon.png',
                                        package: 'euruswallet',
                                        width: 48,
                                        height: 39,
                                        fit: BoxFit.contain),
                                    minRadius: 64,
                                  ),
                                ),
                              )
                            : Column(children: [
                                Padding(padding: EdgeInsets.all(padding)),
                                Image.file(File(imagePath!),
                                    width: 200, height: 200, fit: BoxFit.cover),
                                Padding(padding: EdgeInsets.all(padding)),
                                SubmitButton(
                                    btnController: btnController,
                                    label: 'COMMON.SUBMIT'.tr(),
                                    onPressed: () {
                                      btnController.reset();
                                      _uploadSuccessDialog();
                                    }),
                                Padding(padding: EdgeInsets.all(padding)),
                                CustomTextButton(
                                    text: "KYC.RETAKE_PHOTO".tr(),
                                    buttonBGColor: FXColor.cancelGrayButton,
                                    onPressed: () {
                                      setState(() {
                                        imagePath = null;
                                      });
                                      onTakedPhoto();
                                    })
                              ])
                      ],
                    ),
                  )),
            )));
  }
}
