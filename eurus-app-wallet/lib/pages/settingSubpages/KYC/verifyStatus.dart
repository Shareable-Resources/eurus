import 'package:collection/collection.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/kycCamera.dart';
import 'package:euruswallet/model/userKYCStatus.dart';
import 'package:euruswallet/pages/settingSubpages/KYC/kycStatus.dart';
import 'package:euruswallet/pages/settingSubpages/KYC/uploadPhotoPage.dart';
import 'package:livechat_inc/livechat_inc.dart';

class VerifyStatus extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final KYCStatusData? data = common.userKYCStatus?.data?.kycdata.lastOrNull;
    final IdentityType type = IdentityType.values[data?.kycDoc ?? 0];
    final isKYCStatusPending =
        data == null ? false : data.kycStatus == 0; //KYCStatusPending
    final isKYCStatusRejected =
        data == null ? false : data.kycStatus == 4; //KYCStatusRejected
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        SelectDocmentType(
            'images/cameraIcon.png',
            "KYC.SELFIE".tr(),
            data?.images
                ?.lastWhereOrNull((e) => e.docType == DocType.Selfie.index),
            () {
          onUploadTap(IdentityType.Passport, DocType.Selfie, context);
        }, isKYCStatusPending, isKYCStatusRejected),
        (type != IdentityType.IdentityCard
            ? SelectDocmentType(
                'images/passportIcon.png',
                "KYC.PASSPORT".tr(),
                data?.images?.lastWhereOrNull(
                    (e) => e.docType == DocType.Passport.index), () {
                onUploadTap(IdentityType.Passport, DocType.Passport, context);
              }, isKYCStatusPending, isKYCStatusRejected)
            : Container()),
        (type != IdentityType.Passport
            ? Column(children: [
                SelectDocmentType(
                    'images/idCardIcon.png',
                    "KYC.ID_CARD_FRONT".tr(),
                    data?.images?.lastWhereOrNull(
                        (e) => e.docType == DocType.IdFront.index), () {
                  onUploadTap(
                      IdentityType.IdentityCard, DocType.IdFront, context);
                }, isKYCStatusPending, isKYCStatusRejected),
                SelectDocmentType(
                    'images/idCardIcon.png',
                    "KYC.ID_CARD_BACK".tr(),
                    data?.images?.lastWhereOrNull(
                        (e) => e.docType == DocType.IdBack.index), () {
                  onUploadTap(
                      IdentityType.IdentityCard, DocType.IdBack, context);
                }, isKYCStatusPending, isKYCStatusRejected)
              ])
            : Container()),
        (isKYCStatusRejected == true
            ? Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Padding(
                    padding: EdgeInsets.only(top: 24.0),
                    child: Text(
                      "KYC.CONTACT_CS".tr(),
                      style: FXUI.normalTextStyle
                          .copyWith(color: FXColor.alertRedColor),
                    )),
                GestureDetector(
                  onTap: () {
                    common.openLiveChat(
                      visitorName:
                          isCentralized() ? common.email : common.currentAddress,
                      visitorEmail: isCentralized() ? common.email : null,
                    );
                  },
                  child: Padding(
                      padding: EdgeInsets.symmetric(vertical: 8.0),
                      child: Text(
                        'SETTING_PAGE.HELP_CENTER_CARD.CS'.tr(),
                        style: FXUI.normalTextStyle.copyWith(
                            color: FXColor.mainBlueColor,
                            decoration: TextDecoration.underline),
                      )),
                )
              ])
            : Container()),
      ],
    );
  }

  void onUploadTap(IdentityType type, DocType docType, BuildContext context) {
    final countryCode =
        common.userKYCStatus?.data?.kycdata?.last.kycCountryCode;
    openKYCCamera(
        context,
        (String path) => CommonMethod().pushReplacementPage(
              page: UploadPhotoPage(
                identityType: type,
                countryCode: countryCode,
                docType: docType,
                images: [path],
              ),
              context: context,
            ),
        cameraSide:
            docType == DocType.Selfie ? CameraSide.front : CameraSide.back,
        docType: docType);
  }
}

class SelectDocmentType extends StatelessWidget {
  final String imgPath;
  final String title;
  final KYCImage? kycImage;
  final VoidCallback onTap;
  final bool isPendingStatus;
  final bool isKYCStatusRejected;
  SelectDocmentType(this.imgPath, this.title, this.kycImage, this.onTap,
      this.isPendingStatus, this.isKYCStatusRejected);

  @override
  Widget build(BuildContext context) {
    final RejectReasonData? reason = kycImage?.rejectReason;
    final bool enableUpload = isKYCStatusRejected != true &&
        (isPendingStatus == true ||
            kycImage == null ||
            kycImage?.status == 2 ||
            kycImage?.status ==
                4); //KYCImageStatusWaitingForResubmit, KYCImageStatusVoided

    return Padding(
        padding: EdgeInsets.symmetric(vertical: 10.0),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
            Expanded(
                flex: 1,
                child: Row(
                  children: [
                    Image(
                        width: 28,
                        height: 26,
                        fit: BoxFit.contain,
                        color: FXColor.blueGreenColor,
                        image: AssetImage(imgPath, package: 'euruswallet')),
                    Padding(
                        padding: EdgeInsets.only(left: 17.0, right: 10.0),
                        child: Text(title, style: FXUI.normalTextStyle)),
                    (kycImage != null
                        ? Image.asset(getStatusIconPath(kycImage),
                            package: 'euruswallet',
                            width: 24,
                            height: 24,
                            fit: BoxFit.contain)
                        : Container()),
                  ],
                )),
            IconButton(
                onPressed: enableUpload == true ? onTap : null,
                icon: Image.asset('images/uploadDocIcon.png',
                    package: 'euruswallet',
                    width: 24,
                    height: 22,
                    color: enableUpload == true ? null : Colors.grey,
                    fit: BoxFit.contain)),
          ]),
          (kycImage != null && kycImage!.status <= 1
              ? Text('KYC.WAITING_APPROVAL'.tr(),
                  style: FXUI.normalTextStyle
                      .copyWith(color: FXColor.grayTextColor))
              : Container()),
          (reason == null
              ? Container()
              : Padding(
                  padding: EdgeInsets.symmetric(vertical: 0.0),
                  child: Text(
                    getReason(reason, context.locale.toString()),
                    style: FXUI.normalTextStyle
                        .copyWith(color: FXColor.alertRedColor),
                  ))),
        ]));
  }

  String getReason(RejectReasonData reason, String locale) {
    if (locale == 'en_US')
      return reason.en;
    else if (locale == "zh_Hans")
      return reason.zhCN;
    else if (locale == "zh_Hant") return reason.zhHK;
    return "";
  }

  String getStatusIconPath(KYCImage? image) {
    switch (image?.status) {
      case 0:
      case 1:
        return 'images/submitWaitingIcon.png'; //received, uploaded
      case 2:
        return 'images/submitFailIcon.png'; //waiting for resubmit
      case 3:
        return 'images/submitSuccessIcon.png'; //approved
      case 4:
        return 'images/submitDisabled.png'; //voided
      default:
        return 'images/submitWaitingIcon.png';
    }
  }
}
