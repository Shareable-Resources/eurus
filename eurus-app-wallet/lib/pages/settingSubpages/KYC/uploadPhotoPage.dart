import 'dart:math' as math;
import 'package:dio/dio.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/customDialogBox.dart';
import 'package:euruswallet/commonUI/kycCamera.dart';
import 'package:euruswallet/model/createKycStatus.dart';
import 'package:euruswallet/model/kYCCountryList.dart';
import 'package:euruswallet/model/submitKycDocument.dart';
import 'kycStatus.dart';

class UploadPhotoPage extends StatefulWidget {
  const UploadPhotoPage({
    Key? key,
    required this.identityType,
    required this.countryCode,
    this.images,
    this.docType: DocType.Unknown,
  }) : super(key: key);

  final IdentityType identityType;
  final DocType docType;
  final List? images;
  final String? countryCode;

  @override
  _UploadPhotoPageState createState() => _UploadPhotoPageState();
}

class _UploadPhotoPageState extends State<UploadPhotoPage> {
  late List _images;
  late bool _isSelfie;
  final List<DocType> idDocmentOrder = [
    DocType.Selfie,
    DocType.IdFront,
    DocType.IdBack
  ];
  final List<DocType> passportDocmentOrder = [DocType.Selfie, DocType.Passport];
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    super.initState();
    _images = List.from(widget.images ?? []);
    _isSelfie =
        widget.docType == DocType.Unknown || widget.docType == DocType.Selfie;
  }

  void onTakedPhoto() async {
    openKYCCamera(
        context,
        (String path) => setState(() {
              _images.add(path);
              Navigator.pop(context);
            }),
        cameraSide: _isSelfie == true ? CameraSide.front : CameraSide.back,
        docType: getDocType(widget.docType, math.max(0, _images.length - 1)));
  }

  void onNext() {
    openKYCCamera(
        context,
        (String path) => setState(() {
              _isSelfie = false;
              _images.add(path);
              Navigator.pop(context);
            }),
        docType: getDocType(widget.docType, _images.length));
  }

  void onRetake() {
    setState(() {
      _images.removeLast();
    });
    onTakedPhoto();
  }

  @override
  Widget build(BuildContext context) {
    return TopBar(
        childWidget: _isSelfie == true
            ? genUploadSelfiePhoto()
            : genUploadPassportPhoto(),
        identityType: widget.identityType);
  }

  Widget genUploadSelfiePhoto() {
    if (_images.isEmpty) {
      return Column(children: [
        Padding(
          padding: const EdgeInsets.only(top: 83.0),
        ),
        GestureDetector(
            onTap: () {
              onTakedPhoto();
            }, // handle your image tap here
            child: CircleAvatar(
              backgroundColor: FXColor.lightGreyColor2,
              child: Image.asset('images/cameraIcon.png',
                  package: 'euruswallet',
                  width: 48,
                  height: 39,
                  fit: BoxFit.contain),
              minRadius: 64,
            )),
        Padding(
          padding: const EdgeInsets.only(bottom: 71.0),
        ),
        CustomTextButton(
          text: 'KYC.UPLOAD_PHOTO'.tr(),
          onPressed: () {
            onTakedPhoto();
          },
        ),
        Padding(
          padding: const EdgeInsets.only(top: 14.0),
          child: Text(
            'KYC.NO_SCAN_OR_COPIES'.tr(),
            style: FXUI.normalTextStyle.copyWith(
              fontSize: 14,
              color: FXColor.textGray,
            ),
          ),
        ),
      ]);
    } else {
      return Column(children: [
        Padding(
          padding: const EdgeInsets.only(top: 27.0),
        ),
        CircleAvatar(
          backgroundColor: FXColor.lightGreyColor2,
          child: ClipOval(
              child: Image.file(File(_images.last),
                  width: 200, height: 200, fit: BoxFit.cover)),
          minRadius: 100,
        ),
        Padding(
          padding: const EdgeInsets.only(top: 19.0),
        ),
        (widget.docType == DocType.Selfie
            ? SubmitButton(
                btnController: btnController,
                label: 'KYC.CONTINUE_ON_PHONE'.tr(),
                onPressed: () {
                  submitDocument(DocType.Selfie);
                  btnController.reset();
                })
            : CustomTextButton(
                text: 'COMMON.NEXT_STEP'.tr(),
                onPressed: () {
                  onNext();
                },
              )),
        RetakeButton(onPressed: () {
          onRetake();
        }),
      ]);
    }
  }

  Widget genUploadPassportPhoto() {
    bool isReadyToSubmit = widget.docType != DocType.Unknown ||
        (widget.identityType == IdentityType.IdentityCard &&
            _images.length == idDocmentOrder.length) ||
        (widget.identityType == IdentityType.Passport &&
            _images.length == passportDocmentOrder.length);
    return Column(children: [
      Padding(
        padding: const EdgeInsets.only(top: 27.0),
      ),
      (widget.identityType == IdentityType.IdentityCard
          ? Transform.rotate(
              angle: -90 * math.pi / 180,
              child: Image.file(File(_images.last),
                  width: MediaQuery.of(context).size.width / 2,
                  fit: BoxFit.contain),
            )
          : Image.file(File(_images.last),
              width: 165, height: 224, fit: BoxFit.contain)),
      Padding(
        padding: const EdgeInsets.only(top: 19.0),
      ),
      (isReadyToSubmit == true
          ? SubmitButton(
              btnController: btnController,
              label: 'KYC.CONTINUE_ON_PHONE'.tr(),
              onPressed: () {
                submitDocument(widget.docType);
                btnController.reset();
              })
          : CustomTextButton(
              text: 'COMMON.NEXT_STEP'.tr(),
              onPressed: () {
                onNext();
              })),
      RetakeButton(onPressed: () {
        onRetake();
      }),
    ]);
  }

  DocType getDocType(DocType docType, int index) {
    DocType type = docType;
    if (docType == DocType.Unknown) {
      if (widget.identityType == IdentityType.IdentityCard) {
        type = idDocmentOrder[index];
      } else if (widget.identityType == IdentityType.Passport) {
        type = passportDocmentOrder[index];
      }
    }
    return type;
  }

  void submitDocument(DocType imageType) async {
    var id = common.userKYCStatus != null &&
            common.userKYCStatus!.data!.kycdata.isNotEmpty
        ? common.userKYCStatus!.data!.kycdata.last.id
        : null;
    bool uploadDocumentSuccess = true;
    String errorMsg = "COMMON.FAILURE".tr();
    if (id == null) {
      CreateKycStatus createKycStatus = await api.createKYCStatus(
          kycCountryCode: widget.countryCode,
          kycDoc: widget.identityType.index);
      id = createKycStatus.data?.id;
    }
    for (var i = 0; i < _images.length; i++) {
      DocType type = getDocType(imageType, i);
      SubmitKycDocument submitKycDocument = await api.submitKYCDocument(
          id: id, imageType: type.index, imgPath: _images[i]);
      if (submitKycDocument.returnCode != 0) {
        uploadDocumentSuccess = false;
        errorMsg = submitKycDocument.message;
        break;
      }
    }
    if (uploadDocumentSuccess == true) {
      SubmitKycDocument submitKycDocument = await api.submitKYCApproval(id: id);
      // if (submitKycDocument.returnCode == 0) {
      _uploadSuccessDialog();
      // }
    } else {
      await showDialog(
          context: context,
          builder: (BuildContext context) {
            return CustomDialogBox(
              btnColor: FXColor.mainBlueColor,
              descriptions: errorMsg,
              buttonText: "COMMON.OK".tr(),
            );
          });
    }
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
                  child: Stack(children: [
                    SingleChildScrollView(
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
                            Padding(padding: EdgeInsets.only(top: 87)),
                            Image(
                              width: 86,
                              height: 86,
                              image: AssetImage(
                                'images/uploadSuccessIcon.png',
                                package: 'euruswallet',
                              ),
                            ),
                            Padding(padding: EdgeInsets.only(top: 15)),
                            Text(
                              'COMMON.SUCCESS'.tr(),
                              style: FXUI.normalTextStyle.copyWith(
                                  fontWeight: FontWeight.bold, fontSize: 20),
                            ),
                            Text(
                              'KYC.UPLOAD_SUCCESS_DIALOG_CONTENT'.tr(),
                              style: FXUI.normalTextStyle.copyWith(
                                  color: FXColor.textGray, fontSize: 14),
                            ),
                            Padding(padding: EdgeInsets.only(top: 69)),
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
                                  Navigator.pushAndRemoveUntil(
                                      context,
                                      MaterialPageRoute(
                                        builder: (_) => KYCStatusPage(),
                                      ),
                                      ModalRoute.withName("HomePage"));
                                },
                                child: Text(
                                  'COMMON.CLOSE'.tr(),
                                  style: FXUI.normalTextStyle
                                      .copyWith(color: Colors.white),
                                ),
                              ),
                            )
                          ]),
                    )
                  ])));
        });
  }
}

class RetakeButton extends StatelessWidget {
  const RetakeButton({
    Key? key,
    this.onPressed,
  }) : super(key: key);

  final VoidCallback? onPressed;

  @override
  Widget build(BuildContext context) {
    return Container(
        width: double.infinity,
        height: 50,
        margin: EdgeInsets.only(top: 16.0),
        child: OutlinedButton(
          onPressed: onPressed,
          style: OutlinedButton.styleFrom(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(10.0),
            ),
            side: BorderSide(width: 1, color: FXColor.mainBlueColor),
          ),
          child: Text('KYC.RETAKE_PHOTO'.tr(),
              style:
                  FXUI.normalTextStyle.copyWith(color: FXColor.mainBlueColor)),
        ));
  }
}

class NextButton extends StatelessWidget {
  NextButton({
    Key? key,
    this.onPressed,
    this.text: 'COMMON.NEXT_STEP',
  }) : super(key: key);

  final VoidCallback? onPressed;
  final String text;

  @override
  Widget build(BuildContext context) {
    return CustomTextButton(
      text: 'COMMON.NEXT_STEP'.tr(),
      onPressed: onPressed,
    );
  }
}

class TopBar extends StatelessWidget {
  const TopBar({
    Key? key,
    this.childWidget,
    required this.identityType,
  }) : super(key: key);

  final Widget? childWidget;
  final IdentityType identityType;

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    return BackGroundImage(
      currentUserType: CurrentUserType.centralized,
      child: Scaffold(
        backgroundColor: Colors.transparent,
        appBar: WalletAppBar(
          title: 'KYC.SUBMIT_DOCUMENT'.tr(),
          backButton: true,
        ),
        body: SingleChildScrollView(
            child: TopCircularContainer(
          height: size.heightWithoutAppBar,
          child: Padding(
            padding: getEdgeInsetsSymmetric(horizontal: 34),
            child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 60.0),
                    child: Text(
                      'KYC.TAKE_PHOTO_TITLE'.tr(),
                      style: FXUI.normalTextStyle.copyWith(fontSize: 14),
                      textAlign: TextAlign.center,
                    ),
                  ),
                  if (childWidget != null) childWidget!,
                ]),
          ),
        )),
      ),
    );
  }
}
