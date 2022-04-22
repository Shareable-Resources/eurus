import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/pages/create_wallet_page.dart';
import 'package:euruswallet/pages/decentralized/decentralized_import_wallet_page.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:livechat_inc/livechat_inc.dart';

import '../common/commonMethod.dart';
import '../common/user_profile.dart';
import '../pages/switchAcPage.dart';
import 'centralized/cenLoginPage.dart';
import 'centralized/centralized_wallet_base_page.dart';

class WelcomePage extends StatefulWidget {
  const WelcomePage({Key? key}) : super(key: key);

  @override
  _WelcomePageState createState() => _WelcomePageState();
}

class _WelcomePageState extends State<WelcomePage> {
  @override
  void initState() {
    super.initState();
    //common.showTestNetPopUp(context: context);

    Future.delayed(Duration(milliseconds: 1000), () {
      _showAuthPopUp();
    });

    EasyLoading.instance
      ..displayDuration = const Duration(milliseconds: 2000)
      ..userInteractions = false
      ..dismissOnTap = false;
    // EasyLoading.show(status: 'COMMON.LOADING_W_DOT'.tr());
  }

  void _showAuthPopUp() async {
    final currentUserTypeString =
        await NormalStorageKit().readValue('currentUserType');
    final currentUserType =
        currentUserTypeString == 'CurrentUserType.centralized'
            ? CurrentUserType.centralized
            : currentUserTypeString == 'CurrentUserType.decentralized'
                ? CurrentUserType.decentralized
                : null;
    if (!mounted || currentUserType == null) return;
    await common.showAuthPopUp(context,
        currentUserType: currentUserType, isNeededRoutingToHomePage: true);
    _showAuthPopUpIfNeeded();
  }

  void _showAuthPopUpIfNeeded() {
    final arguments = (ModalRoute.of(context)?.settings.arguments as Map);
    if (arguments['shouldShowAuthPopUp'] ?? false) {
      arguments['shouldShowAuthPopUp'] = false;
      _showAuthPopUp();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: DecoratedBox(
      decoration: BoxDecoration(
        color: Colors.white,
        image: DecorationImage(
            image: AssetImage('images/bgLogin.png', package: 'euruswallet'),
            fit: BoxFit.cover),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: <Widget>[
          Expanded(
            child: Container(),
          ),
          Image.asset(
            'images/imgLogo.png',
            package: 'euruswallet',
            fit: BoxFit.fitHeight,
            height: 137,
          ),
          SizedBox(
            height: 109,
          ),
          // Bottom Sheet style bottom aligned half-shown vertical Card
          DecoratedBox(
            decoration: BoxDecoration(
              borderRadius: BorderRadiusDirectional.only(
                topStart: Radius.circular(30),
                topEnd: Radius.circular(30),
              ),
              color: FXColor.mainBlueColor,
            ),
            child: Padding(
              padding: EdgeInsets.symmetric(horizontal: 35).add(
                EdgeInsets.only(
                    top: 35,
                    bottom: 14 + MediaQuery.of(context).padding.bottom),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.start,
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: <Widget>[
                  CupertinoButton(
                    color: Colors.white,
                    padding: EdgeInsets.zero,
                    child: Text(
                      'LAUNCH_PAGE.CREATE_ACCOUNT_BTN'.tr(),
                      style: FXUI.normalTextStyle.copyWith(
                        color: FXColor.mainBlueColor,
                        fontSize: 16,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                    onPressed: () async {
                      await Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => CentralizedWalletBasePage(
                            appBarTitle:
                                Text('CREATE_WALLET_PAGE.CREATE_WALLET'.tr()),
                            body: CreateWalletPage(),
                          ),
                        ),
                      );
                      _showAuthPopUpIfNeeded();
                    },
                    borderRadius: FXUI.cricleRadius,
                  ),
                  SizedBox(height: 12),
                  CupertinoButton(
                    color: Colors.white,
                    padding: EdgeInsets.zero,
                    child: Text(
                      'LAUNCH_PAGE.LOGIN_BTN'.tr(),
                      style: FXUI.normalTextStyle.copyWith(
                          color: FXColor.mainBlueColor,
                          fontSize: 16,
                          fontWeight: FontWeight.w500),
                    ),
                    onPressed: () async {
                      List<UserProfile> localAcs = await getLocalAcs();
                      if (localAcs.length > 0) {
                        await common.pushPage(
                            page: SwitchAcPage(), context: context);
                      } else {
                        await common.pushPage(
                            page: CenLoginPage(), context: context);
                      }
                      _showAuthPopUpIfNeeded();
                    },
                    borderRadius: FXUI.cricleRadius,
                  ),
                  SizedBox(
                    height: 20 - 3.0,
                  ),
                  Row(
                    children: [
                      Expanded(
                        child: Divider(
                          color: Colors.white,
                          thickness: 2,
                        ),
                      ),
                      Padding(
                        padding: EdgeInsets.symmetric(horizontal: 10),
                        child: Text(
                          'COMMON.OR'.tr(),
                          style: FXUI.normalTextStyle.copyWith(
                              color: Colors.white,
                              fontSize: 16,
                              fontWeight: FontWeight.w500),
                        ),
                      ),
                      Expanded(
                        child: Divider(
                          color: Colors.white,
                          thickness: 2,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: 17,
                  ),
                  ListTile(
                    contentPadding: EdgeInsets.zero,
                    title: IntrinsicHeight(
                      child: CupertinoButton(
                        /// OutlineButton
                        padding: EdgeInsets.zero,
                        borderRadius: FXUI.cricleRadius,
                        onPressed: () async {
                          await Navigator.push(
                            context,
                            MaterialPageRoute(
                              builder: (_) => DecentralizedImportWalletPage(),
                            ),
                          );
                          _showAuthPopUpIfNeeded();
                        },

                        /// SizedBox.expand and Expanded both require IntrinsicHeight wrapping CupertinoButton
                        child: SizedBox.expand(
                          child: Container(
                            decoration: BoxDecoration(
                              border: Border.all(color: Colors.white),
                              borderRadius: FXUI.cricleRadius,
                            ),
                            child: Center(
                              child: Text(
                                'LAUNCH_PAGE.IMPORT_WALLET_BTN'.tr(),
                                style: FXUI.normalTextStyle
                                    .copyWith(color: Colors.white),
                                textAlign: TextAlign.center,
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                  InkWell(
                    child: Text('LAUNCH_PAGE.CUSTOMER_SERVICE'.tr(),
                        style: Theme.of(context)
                            .textTheme
                            .overline
                            ?.apply(color: Colors.white),
                        textAlign: TextAlign.center),
                    onTap: () {
                      common.openLiveChat();
                    },
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    ));
  }
}
