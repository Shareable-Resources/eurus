import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/pages/welcome_page.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:loader_overlay/loader_overlay.dart';

import '../common/commonMethod.dart';
import 'eurus_wallet_home_page.dart';

class MyApp extends StatefulWidget {
  const MyApp({
    Key? key,
  }) : super(key: key);

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    common.currentContext = context;
    return GestureDetector(
      onTap: () {
        FocusScopeNode currentFocus = FocusScope.of(context);

        if (currentFocus.hasFocus) {
          currentFocus.unfocus();
        }
      },
      child: MaterialApp(
        localizationsDelegates: context.localizationDelegates,
        supportedLocales: context.supportedLocales,
        locale: context.locale,
        title: 'Eurus Wallet',
        theme: ThemeData(
          primarySwatch: Colors.blue,
          visualDensity: VisualDensity.adaptivePlatformDensity,
          appBarTheme: AppBarTheme(
            systemOverlayStyle: SystemUiOverlayStyle.dark,
          ),
        ),
        initialRoute: 'WelcomePage',
        onGenerateRoute: (route) {
          if (route.name == 'WelcomePage') {
            return MaterialPageRoute(
              builder: (_) => WelcomePage(),
              settings: RouteSettings(name: 'WelcomePage', arguments: Map()),
            );
          } else if (route.name == 'HomePage') {
            return PageRouteBuilder(
              settings: RouteSettings(name: 'HomePage', arguments: Map()),
              fullscreenDialog: true,
              opaque: false,
              pageBuilder:
                  (pageBuilderContext, animation, secondaryAnimation) =>
                      LoaderOverlay(
                useDefaultLoading: false,
                overlayWidget: Center(
                  child: CircularProgressIndicator(
                    color: common.getBackGroundColor(),
                  ),
                ),
                child: FlutterEasyLoading(
                  child: EurusWalletHomePage(),
                ),
              ),
            );
          }
        },
      ),
    );
  }
}
