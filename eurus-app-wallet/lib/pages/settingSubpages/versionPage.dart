import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/pages/settingSubpages/cardContainer.dart';
import 'package:euruswallet/pages/settingSubpages/settingAppBar.dart';
import 'package:package_info/package_info.dart';
import 'package:easy_localization/easy_localization.dart';

class VersionPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: FXColor.veryLightGreyTextColor,
      appBar: SettingAppBar(true),
      body: Container(
        child: SingleChildScrollView(
          child: SafeArea(
            child: CardContainer(
              'VERSION_PAGE.VERSION'.tr(),
              Container(
                padding: EdgeInsets.symmetric(horizontal: 13, vertical: 25),
                child: Column(
                  children: [
                    Row(
                      children: [
                        Spacer(flex: 1),
                        Expanded(
                          flex: 1,
                          child: Image(
                            image: AssetImage(
                              'images/Eurus_Vertical_logo.png',
                              package: 'euruswallet',
                            ),
                          ),
                        ),
                        Spacer(flex: 1),
                      ],
                    ),
                    Padding(
                      padding: EdgeInsets.only(top: 15),
                      child: FutureBuilder(
                        future: PackageInfo.fromPlatform(),
                        builder: (_, AsyncSnapshot<PackageInfo> snapshot) {
                          return Text(
                            'v${snapshot.data?.version ?? ''}',
                            style: FXUI.normalTextStyle.copyWith(
                              color: FXColor.mediumGrayColor,
                            ),
                          );
                        },
                      ),
                    ),
                    Padding(
                      padding: EdgeInsets.only(left: 27, right: 27, top: 35),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            'VERSION_PAGE.VERSION'.tr(),
                            style: FXUI.normalTextStyle.copyWith(
                              fontWeight: FontWeight.w600,
                              fontSize: 16,
                            ),
                          ),
                          FutureBuilder(
                            future: PackageInfo.fromPlatform(),
                            builder: (_, AsyncSnapshot<PackageInfo> snapshot) {
                              return Text(
                                '${snapshot.data?.version ?? ''}',
                                style: FXUI.normalTextStyle.copyWith(
                                  color: FXColor.mediumGrayColor,
                                ),
                              );
                            },
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
