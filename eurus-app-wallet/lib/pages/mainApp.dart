import 'package:apihandler/apiHandler.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/customDialogBox.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_crashlytics/firebase_crashlytics.dart';
import 'package:package_info/package_info.dart';
import 'package:uni_links/uni_links.dart';

import 'decentralized/receivePage.dart';
import 'main.dart';

class MainApp extends StatefulWidget {
  MainApp({Key? key}) : super(key: key);

  @override
  _MainAppState createState() => _MainAppState();
}

class _MainAppState extends State<MainApp> with WidgetsBindingObserver {
  String? _initialLink;
  late String version;
  late var subscription;

  @override
  void initState() {
    super.initState();
    // testing code
    //api.pushRawTransaction(alice: common.bitcoinKeyPair, targetAddress: "mhARbrfUzFoepzrAwdzwEVWAi3Wmq22FmR", amount: 200, selfAddress: "mjaufUtcBUdZpx3AFxwjNX2fW3vQqhy5wD");
    //api.getBlockchainAddressInformation(address: 'mjaufUtcBUdZpx3AFxwjNX2fW3vQqhy5wD');
    common.currentContext = context;

    checkNetworkConnection();
  }

  _setupApp() {
    _checkCountry().then((value) {
      if (value) {
        WidgetsBinding.instance?.addObserver(this);
        _getServerConfig(envType: EnvType.Dev);
        initPlatformState();
        Future.delayed(Duration(milliseconds: 1000), () async {
          common.forceAppUpdateIfNeeded(common.currentContext);
          common.alreadyRoot(common.currentContext);
        });
      }
    });
  }

  _getServerConfig({required EnvType envType}) async {
    if (envType == EnvType.Production) {
      await Firebase.initializeApp();
      FirebaseCrashlytics.instance.setCrashlyticsCollectionEnabled(true);
      FlutterError.onError = FirebaseCrashlytics.instance.recordFlutterError;
    }

    api.setUpServerConfig(newEnvType: envType).then((value) {
      if (!value) {
        _getServerConfig(envType: envType);
        return;
      }
      initEThClient();
      api.getFaucetList();
    });
  }

  Future<bool> _checkCountry() async {
    final response = await api.get(
      url: 'http://www.geoplugin.net/json.gp',
      haveApiAccessToken: false,
    );
    final country = response['geoplugin_countryName'].toString();
    if (country.toLowerCase() == 'china') {
      showDialog(
        context: common.currentContext,
        barrierDismissible: false,
        builder: (_) {
          return WillPopScope(
            onWillPop: () async => false,
            child: Dialog(
              backgroundColor: Colors.transparent,
              insetPadding: EdgeInsets.all(10),
              child: Container(
                width: double.infinity,
                padding: EdgeInsets.symmetric(vertical: 25, horizontal: 17),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: FXUI.cricleRadius,
                ),
                child: SingleChildScrollView(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        'COMMON.WARNING'.tr(),
                        style: FXUI.inputStyle.copyWith(
                          fontWeight: FontWeight.bold,
                          fontSize: 20,
                        ),
                      ),
                      SizedBox(height: 42),
                      Image.asset(
                        'images/icn_warning_rounded.png',
                        package: 'euruswallet',
                        width: 80,
                      ),
                      SizedBox(height: 34),
                      Text(
                        'COMMON.COUNTRY_NOT_AVAILABLE'.tr(),
                        style: FXUI.subtitleTextStyle.copyWith(
                          fontSize: 13,
                        ),
                      ),
                      SizedBox(height: 34),
                      TextButton(
                        onPressed: () => exit(0),
                        child: Container(
                          decoration: BoxDecoration(
                            color: common.getBackGroundColor(),
                            borderRadius: FXUI.cricleRadius,
                          ),
                          padding: EdgeInsets.all(15),
                          width: double.infinity,
                          child: Center(
                            child: Text(
                              // 'Back',
                              'COMMON.CLOSE'.tr(),
                              style: FXUI.normalTextStyle.copyWith(
                                  color: Colors.white,
                                  fontWeight: FontWeight.bold,
                                  fontSize: 16),
                            ),
                          ),
                        ),
                      ),
                      SizedBox(height: 9),
                    ],
                  ),
                ),
              ),
            ),
          );
        },
      );
      return false;
    }
    return true;
  }

  Future<void> checkNetworkConnection() async {
    await Future.delayed(Duration(seconds: 1));
    bool hasConnection = false;
    try {
      final result = await InternetAddress.lookup('baidu.com');
      if (result.isNotEmpty && result[0].rawAddress.isNotEmpty) {
        hasConnection = true;
      } else {
        hasConnection = false;
      }
    } on SocketException catch (_) {
      hasConnection = false;
    }

    if (!hasConnection) {
      showNotNetworkConnectionError().then((value) {
        checkNetworkConnection();
      });
    }

    subscription = Connectivity()
        .onConnectivityChanged
        .listen((ConnectivityResult result) async {
      if (result == ConnectivityResult.none) {
        showNotNetworkConnectionError();
      }
      print("result:$result");
    });

    _setupApp();
  }

  Future<void> showNotNetworkConnectionError() async {
    await showDialog(
        context: common.currentContext,
        builder: (BuildContext context) {
          return CustomDialogBox(
            descriptions: "COMMON.NOT_NETWORK_CONNECTION".tr(),
            buttonText: "COMMON.OK".tr(),
          );
        });
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    switch (state) {
      case AppLifecycleState.resumed:
        api.reFreshToken();
        if (mounted) {
          common.forceAppUpdateIfNeeded(common.currentContext);
        }
        break;
      case AppLifecycleState.inactive:
      case AppLifecycleState.paused:
      case AppLifecycleState.detached:
        break;
    }
  }

  Future<void> initEThClient() async {
    web3dart.setErc20Contract(
        blockChainType: BlockChainType.Ethereum,
        contractAddress: '0x022E292b44B5a146F2e8ee36Ff44D3dd863C915c');
    web3dart.setErc20Contract(
        blockChainType: BlockChainType.Eurus,
        contractAddress: '0xa54Dee79c3bB34251DEbf86C1BA7D21898FFb7AC');
    web3dart.setErc20Contract(
        blockChainType: BlockChainType.BinanceCoin,
        contractAddress: '0x0F603B818a93E1e5E007CaD19065cfe8645F5d12');
  }

  Future<void> initPlatformState() async {
    await initPlatformStateForStringUniLinks();
    await initPlatformStateForUriUniLinks();
    version = await getVersionNumber();
  }

  Future<String> getVersionNumber() async {
    PackageInfo packageInfo = await PackageInfo.fromPlatform();
    String version = 'Eurus wallet v';
    version += packageInfo.version;
    return version;
  }

  /// An implementation using a [String] link
  Future<void> initPlatformStateForStringUniLinks() async {
    // Get the latest link
    // Platform messages may fail, so we use a try/catch PlatformException.
    try {
      _initialLink = await getInitialLink();
      print('initial link: $_initialLink');
      if (_initialLink != null) {
        // common.pushPage(page: SelectTargetPage(), context: context);
      }
    } on PlatformException {
      _initialLink = 'Failed to get initial link.';
    } on FormatException {
      _initialLink = 'Failed to parse the initial link as Uri.';
    }
  }

  /// An implementation using the [Uri] convenience helpers
  Future<void> initPlatformStateForUriUniLinks() async {
    // Attach a second listener to the stream
    uriLinkStream.listen((Uri? uri) {
      print('got uri: ${uri?.path} ${uri?.queryParametersAll}');
      if (uri?.normalizePath() != null) {
        common.pushPage(
            page: ReceivePage(
              ethereumAddress: '',
              errorPromptBuilder: (_) {
                return Container();
              },
            ),
            context: context);
      }
    }, onError: (Object err) {
      print('got err: $err');
    });
  }

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);

    Future.microtask(() async {
      await common.pushReplacementPage(
          settings: RouteSettings(name: "/"), page: MyApp(), context: context);
    });

    return MaterialApp(
      localizationsDelegates: context.localizationDelegates,
      supportedLocales: context.supportedLocales,
      locale: context.locale,
      builder: (BuildContext context, Widget? widget) {
        return Scaffold(
          backgroundColor: Colors.transparent,
        );
      },
    );
  }
}
