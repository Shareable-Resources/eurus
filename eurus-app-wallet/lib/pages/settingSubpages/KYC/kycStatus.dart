import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/pages/settingSubpages/KYC/selectDocumentType.dart';
import 'package:euruswallet/pages/settingSubpages/KYC/verifyStatus.dart';
import 'package:euruswallet/model/userKYCStatus.dart';

enum IdentityType { Unknown, Passport, IdentityCard }
enum DocType { Unknown, Passport, IdFront, IdBack, Selfie }

class KYCStatusPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    return BackGroundImage(
        currentUserType: CurrentUserType.centralized,
        child: Scaffold(
            backgroundColor: Colors.transparent,
            appBar: WalletAppBar(
              title: 'KYC.VERIFY_IDENTITY_TITLE'.tr(),
              backButton: true,
            ),
            body: FutureBuilder(
                future: api.getKYCStatus(),
                builder: (_, AsyncSnapshot<UserKYCStatus> kycStatus) {
                  return SingleChildScrollView(
                    child: TopCircularContainer(
                        height: size.heightWithoutAppBar,
                        child: Padding(
                          padding: getEdgeInsetsSymmetric(horizontal: 34),
                          child: kycStatus.data != null
                              ? (kycStatus.data?.data?.kycdata != null &&
                                      kycStatus.data!.data!.kycdata.isNotEmpty
                                  ? VerifyStatus()
                                  : SelectDocumentType())
                              : Container(
                                  height: size.blockSizeVertical * 70,
                                  child: Center(
                                      child: CircularProgressIndicator())),
                        )),
                  );
                })));
  }
}
