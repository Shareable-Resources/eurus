import 'dart:convert';

import 'package:collection/src/iterable_extensions.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/user_profile.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:flutter_neumorphic/flutter_neumorphic.dart';
import 'package:rxdart/rxdart.dart';

import '../common/biometric_authentication_helper.dart';
import '../common/commonMethod.dart';
import '../commonUI/crypto_asset_list_widget.dart';
import '../commonUI/replacing_qrcode_widget.dart';
import '../commonUI/topAppBar.dart';
import '../commonUI/topSelectBlockChainBar.dart';
import '../commonUI/walletLockerPWDialog.dart';
import '../model/getFaucetList.dart';
import 'assets_detail_page.dart';
import 'centralized/top_up_payment_wallet_page.dart';
import 'coinPriceRanking.dart';
import 'dappBrowser/dapp_browser_home_page.dart';
import 'decentralized/backup_seed_phrases_home_page.dart';
import 'decentralized/decentralized_wallet_base_page.dart';
import 'decentralized/edit_assets_list.dart';
import 'reward_scheme_dialog.dart';
import 'settingPage.dart';

class EurusWalletHomePage extends StatefulWidget {
  EurusWalletHomePage({
    Key? key,
  }) : super(key: key);

  @override
  _EurusWalletHomePageState createState() => _EurusWalletHomePageState();
}

class _EurusWalletHomePageState extends State<EurusWalletHomePage>
    with WidgetsBindingObserver {
  late Future<String?> _setupFuture;

  Map<String, dynamic> _cacheAssetBalance = {};

  int navigationBarTabBuilderSelectedIndex = 0;

  @override
  void initState() {
    _setupFuture = Future.delayed(Duration(milliseconds: 150), () async {
      // TODO - here can await a result pop from a navigation push into wallet account selection dialog
      // TODO - iOS Keychain service keep retaining keychain items even after app uninstall and reinstall due to platform behaviour
      final encryptedAddress =
          await common.findDefaultAccountEncryptedAddress();
      final address = await SecureStorageKit().readValue('currentUserAddress');
      List<UserProfile> userProfiles = await getLocalAcs();
      UserProfile? userProfile = userProfiles.firstWhereOrNull(
          (element) => element.address == address.toString().toLowerCase());
      if (!isEmptyString(string: encryptedAddress) && userProfile != null) {
        common.encryptedAddress = encryptedAddress;
        common.refreshTopBar.add(true);
        var _defaultShowBalance = await NormalStorageKit()
            .readValue('defaultShowAssetBalance_$encryptedAddress');
        common.showBalance = _defaultShowBalance == '0' ? false : true;
        await setAcToLocal(userProfile, delete: true);
        String newAuthValidTime = DateTime.now()
            .add(Duration(days: 7))
            .millisecondsSinceEpoch
            .toString();
        if (userProfile.enableBiometric == true)
          userProfile.bioAuthValidTime = newAuthValidTime;
        if (userProfile.enableTxBiometric == true)
          userProfile.bioTxAuthValidTime = newAuthValidTime;
        common.currentUserProfile = userProfile;
        await setAcToLocal(userProfile);
      }
      return encryptedAddress;
    });

    super.initState();
    print('_EurusWalletHomePageState initState()');

    // set Auth UI of private key Handler to Web3 ETH Client
    if (mounted) {
      web3dart.canGetPrivateKeyHandler = ({bool? isTx}) async {
        final _isTx = isTx ?? true;

        final _submit = (String? _accountPrivateKey,
            TextEditingController _textEditingController) async {
          String? _decryptedValue = _accountPrivateKey != null
              ? CommonMethod.passwordDecrypt(
                  _textEditingController.text, _accountPrivateKey)
              : null;
          return _decryptedValue;
        };
        var _bioAuthHelper = BiometricAuthenticationHelper();
        final _tryBioAuth = () => Future.delayed(
              Duration.zero,
              () async => _isTx == true
                  ? _bioAuthHelper.canPersistTxWithBiometricSecurely(
                      await common.prefix + (common.encryptedAddress ?? ''))
                  : _bioAuthHelper.canPersistWithBiometricSecurely(
                      await common.prefix + (common.encryptedAddress ?? '')),
            ).then((value) => value);

        var _isBioLoginEnabled =
            common.currentUserProfile?.enableBiometric ?? false;
        final _isBioTxEnabled =
            common.currentUserProfile?.enableTxBiometric ?? false;

        return await Navigator.of(context)
            .push(PageRouteBuilder(
              fullscreenDialog: true,
              opaque: false,
              pageBuilder:
                  (pageBuilderContext, animation, secondaryAnimation) =>
                      WalletLockerPWDialog(
                          themeColor: common.getBackGroundColor(),
                          decenUserkey:
                              common.currentUserProfile!.encryptedPrivateKey,
                          submitFnc: _submit,
                          tryBioAuthFnc: _isTx && _isBioTxEnabled ||
                                  !_isTx && _isBioLoginEnabled
                              ? _tryBioAuth
                              : null),
            ))
            .then((value) => value is String ? value : '');
      };
    }
    common.refreshAssetList = BehaviorSubject<bool>();
    common.refreshAssetList.listen((bool value) async {
      common.assetsList = await NormalStorageKit()
          .readValue('assetsList_${common.encryptedAddress}');
      setState(() {
        print("assetsList:${common.assetsList}");
      });
    });
    common.updateNavbarOnChangeLang = () {
      setState(() {});
    };
    if (envType != EnvType.Production) {
      api.getFaucetList();
    }
    WidgetsBinding.instance?.addObserver(this);
    common.logoutFunction = BehaviorSubject<bool>();
    common.logoutFunction?.listen((value) {
      logoutFnc();
    });

    showRewardIfNeeded();
  }

  @override
  void dispose() {
    WidgetsBinding.instance?.removeObserver(this);
    _appPausedTimer?.cancel();
    super.dispose();
  }

  Timer? _appPausedTimer;
  bool _shouldShowAuthPopUp = false;
  bool _isAuthPopUpShowed = false;

  @override
  Future<void> didChangeAppLifecycleState(AppLifecycleState state) async {
    switch (state) {
      case AppLifecycleState.resumed:
        if (_appPausedTimer != null) {
          _appPausedTimer?.cancel();
          _appPausedTimer = null;
        }
        if (mounted && _shouldShowAuthPopUp && !_isAuthPopUpShowed) {
          _shouldShowAuthPopUp = false;
          _isAuthPopUpShowed = true;
          await common.showAuthPopUp(
            context,
            currentUserType: common.currentUserType,
          );
          _isAuthPopUpShowed = false;
        }
        if (_appPausedTimer != null) {
          _appPausedTimer?.cancel();
          _appPausedTimer = null;
        }
        break;
      case AppLifecycleState.inactive:
      case AppLifecycleState.detached:
        break;
      case AppLifecycleState.paused:
        _appPausedTimer = Timer(Duration(minutes: 1), () {
          _shouldShowAuthPopUp = true;
        });
        break;
    }
  }

  Future<void> logoutFnc() async {
    Navigator.popUntil(context, (route) => route.isFirst);
    await SecureStorageKit().setValue(null, await common.prefix + 'accounts');
    await NormalStorageKit().setValue('', 'displayUserName');
    await NormalStorageKit().setValue('', 'apiAccessTokenExpiryTime_');
    common.currentAddress = null;
    common.timer.cancel();
  }

  final pageController = PageController();

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    print('size.screenHeight${size.screenHeight}');
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: FutureBuilder(
        future: _setupFuture,
        builder: (ctx, AsyncSnapshot<String?> snapshot) => snapshot
                        .connectionState ==
                    ConnectionState.done &&
                (snapshot.hasError ||
                    !snapshot.hasData ||
                    snapshot.hasData && isEmptyString(string: snapshot.data))
            ? Container()
            : Scaffold(
                appBar: navigationBarTabBuilderSelectedIndex == 1 ||
                        navigationBarTabBuilderSelectedIndex == 3
                    ? null
                    : TopAppBar(
                        rewardButtonOnTap: () =>
                            common.routeToRewardDetailPage(context: context),
                        refuelButtonOnTap: () {
                          Navigator.push(
                            context,
                            MaterialPageRoute(
                              builder: (context) => TopUpPaymentWalletPage(),
                            ),
                          ).whenComplete(() => setState(() {}));
                        },
                      ),
                backgroundColor: FXColor.lightWhiteColor,
                body: PageView(
                  controller: pageController,
                  onPageChanged: (index) {
                    setState(() {
                      navigationBarTabBuilderSelectedIndex = index;
                    });
                  },
                  physics: NeverScrollableScrollPhysics(),
                  children: createTabBarItemBody(),
                ),
                bottomNavigationBar: BottomNavigationBar(
                  type: BottomNavigationBarType.fixed,
                  showUnselectedLabels: true,
                  selectedItemColor: common.getBackGroundColor(),
                  selectedFontSize: 11,
                  unselectedFontSize: 11,
                  unselectedItemColor: FXColor.DoveGreyColor,
                  items: [
                    createTabBarItem(
                      icon: 'images/icn_wallet.png',
                      title: 'MAIN_NAVBAR.WALLET'.tr(),
                    ),
                    createTabBarItem(
                      icon: 'images/icn_browser.png',
                      title: 'MAIN_NAVBAR.BROWSER'.tr(),
                    ),
                    createTabBarItem(
                      icon: 'images/icn_market.png',
                      title: 'MAIN_NAVBAR.MARKETS'.tr(),
                    ),
                    createTabBarItem(
                      icon: 'images/icn_setting.png',
                      title: 'MAIN_NAVBAR.SETTING'.tr(),
                    ),
                  ],
                  currentIndex: navigationBarTabBuilderSelectedIndex,
                  onTap: (index) {
                    pageController.jumpToPage(index);
                  },
                ),
              ),
      ),
    );

    // method 1: press to select account to push to password unlock
    // method 2: Stack password unlock widget always on top of / above wallet home page
    // method 3: StreamBuilder or FutureBuilder challenge
    // method 4: wallet home page widget internally build() a password unlock widget before actual wallet home page
    // method 5: uses of Widget class lifecycle like onMount, didMount, initState
    // method 6: for App Lifecycle resumes, WidgetBinding
    // method 7: wrapper around wallet home page
  }

  BottomNavigationBarItem createTabBarItem({
    required String icon,
    required String title,
  }) {
    final double iconWidth = 25;
    return BottomNavigationBarItem(
      icon: Image.asset(
        icon,
        package: 'euruswallet',
        width: iconWidth,
        height: iconWidth,
        color: FXColor.DoveGreyColor,
      ),
      activeIcon: Image.asset(
        icon,
        package: 'euruswallet',
        width: iconWidth,
        height: iconWidth,
        color: common.getBackGroundColor(),
      ),
      label: title,
    );
  }

  List<Widget> createTabBarItemBody() {
    int _viewSelectedChainCode =
        EURUS_CHAIN_CODE; // default is EURUS_CHAIN_CODE (1) otherwise Testnet (-1) or Ethereum (0)
    Map<String, dynamic> _cArgs = {
      ...{
        if (common.currentUserProfile?.seedPraseBackuped != true &&
            common.currentUserType == CurrentUserType.decentralized)
          'replacingQRCodeWidget': ReplacingQRCodeWidget(
            backupSeedPhrasesCompletion: () => setState(() {}),
          ),
        'walletAccountEncryptedAddress': common.encryptedAddress,
        'canGetWalletAccountAssetsList': () async => await NormalStorageKit()
            .readValue('assetsList_${common.encryptedAddress}'),
      }
    };

    final _navigateToAssetDetails = (cryptoCurrencyModelMap) async {
      if (common.getFaucetList == null) {
        await api.getFaucetList();
      }
      common.cryptoCurrencyModelMap = cryptoCurrencyModelMap;
      common.showFaucetBtn = false;
      for (GetFaucetListData getFaucetListData
          in common.getFaucetList?.data ?? []) {
        if (getFaucetListData.key == common.cryptoCurrencyModelMap?['symbol']) {
          if (cryptoCurrencyModelMap['address'] ==
              cryptoCurrencyModelMap['address$TEST_NET']) {
            common.showFaucetBtn = true;
          }
        }
      }
      await Navigator.of(context).push(
        MaterialPageRoute(
          builder: (_) => AssetsDetailPage(
            cryptoCurrencyModelMap: cryptoCurrencyModelMap,
            cArgs: _cArgs,
          ),
        ),
      );
    };

    return [
      NestedScrollView(
          headerSliverBuilder: (BuildContext context, bool innerBoxIsScrolled) {
            return <Widget>[
              SliverAppBar(
                foregroundColor: Colors.transparent,
                backgroundColor: Colors.transparent,
                expandedHeight: 205.0,
                flexibleSpace: FlexibleSpaceBar(
                  background: common.getBannerPath(context: context) != null
                      ? Padding(
                          padding:
                              const EdgeInsets.fromLTRB(20.0, 0.0, 20.0, 10.0),
                          child: GestureDetector(
                            onTap: () => common.routeToRewardDetailPage(
                                context: context),
                            child: Neumorphic(
                              child: Image.asset(
                                common.getBannerPath(context: context)!,
                                package: 'euruswallet',
                                fit: BoxFit.fill,
                              ),
                              style: FXUI.neumorphicBannerImage,
                            ),
                          ),
                        )
                      : Container(),
                ),
              ),
            ];
          },
          body: Column(children: [
            IntrinsicHeight(
                child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 4.0),
              child: Row(
                children: common
                    .homeFunctionsBarItemBuilder(
                      context,
                      _cArgs,
                    )
                    .map((e) => Expanded(child: e))
                    .toList(),
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                crossAxisAlignment: CrossAxisAlignment.stretch,
              ),
            )),
            Expanded(
              child: Padding(
                padding: EdgeInsets.only(top: 10),
                child: Container(
                  clipBehavior: Clip.hardEdge,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadiusDirectional.only(
                        topStart: Radius.circular(FXUI.circular),
                        topEnd: Radius.circular(FXUI.circular)),
                    color: Colors.white,
                    boxShadow: kElevationToShadow[3],
                  ),
                  child: Padding(
                    padding: EdgeInsets.only(top: 0, left: 0, right: 0),
                    child: Column(
                      children: [
                        Expanded(
                          child: Stack(
                            children: [
                              // asset list section content body
                              // SafeArea(child:
                              Padding(
                                padding: EdgeInsets.only(top: 57),
                                child: FutureBuilder(
                                  future: common.updateAssetsListIfEmpty(),
                                  builder: (_, AsyncSnapshot<String> s) {
                                    if (s.connectionState !=
                                            ConnectionState.done &&
                                        isEmptyString(
                                            string: common.assetsList)) {
                                      return Center(
                                          child: CircularProgressIndicator());
                                    } else {
                                      if ((isEmptyString(string: s.data) &&
                                          isEmptyString(
                                              string: common.assetsList))) {
                                        return Builder(
                                          builder: (_) {
                                            print(
                                                '${common.encryptedAddress} s.data ${s.data}');
                                            return ListView(
                                              padding: EdgeInsets.only(
                                                  top: 0,
                                                  left: 15,
                                                  right: 15,
                                                  bottom:
                                                      SCREEN_WITH_BOTTOM_NAV_TAB_BAR_SAFE_AREA_BOTTOM_PADDING_CONTENT_INSET),
                                              children: [
                                                Container(
                                                  child: ListTile(
                                                      title: CupertinoButton(
                                                          onPressed: () async {
                                                            await Navigator.of(
                                                                    _)
                                                                .push(
                                                              MaterialPageRoute(
                                                                builder: (_) =>
                                                                    EditAssetsListPage(
                                                                        userSuffix:
                                                                            common.encryptedAddress ??
                                                                                ''),
                                                              ),
                                                            );
                                                            setState(() {});
                                                          },
                                                          child: Text(
                                                              "COMMON.ADD_TOKEN"
                                                                  .tr(),
                                                              style: Theme.of(_)
                                                                  .textTheme
                                                                  .button
                                                                  ?.apply(
                                                                      color: Theme.of(
                                                                              _)
                                                                          .colorScheme
                                                                          .primary)))),
                                                  margin: EdgeInsets.all(4),
                                                  decoration: BoxDecoration(
                                                    shape: BoxShape.rectangle,
                                                    boxShadow:
                                                        kElevationToShadow[3],
                                                    color: Colors.white,
                                                    borderRadius:
                                                        BorderRadius.circular(
                                                            15),
                                                  ),
                                                ),
                                              ],
                                            );
                                          },
                                        );
                                      } else {
                                        common.assetsList =
                                            isEmptyString(string: s.data)
                                                ? common.assetsList
                                                : s.data;
                                        return CryptoAssetListWidget(
                                          data: jsonDecode(s.data ??
                                              common.assetsList ??
                                              '') as List,
                                          addressSuffix: getAddressSuffix(
                                              common.topSelectedBlockchainType),
                                          listViewTopPaddingContentInset: 0.0,
                                          additionalListViewBottomPaddingContentInset:
                                              SCREEN_WITH_BOTTOM_NAV_TAB_BAR_SAFE_AREA_BOTTOM_PADDING_CONTENT_INSET,
                                          onTapHandler: (data) async {
                                            if (data['supported'] as bool? ??
                                                false)
                                              await _navigateToAssetDetails(
                                                      data)
                                                  .whenComplete(
                                                      () => setState(() {}));
                                          },
                                          assetBalanceCache: _cacheAssetBalance,
                                        );
                                      }
                                    }
                                  },
                                ),
                              ),
                              // asset list section header
                              Align(
                                alignment: Alignment.topCenter,
                                child: Container(
                                  decoration: BoxDecoration(),
                                  child: Stack(
                                    children: [
                                      Padding(
                                        padding:
                                            EdgeInsets.symmetric(horizontal: 17)
                                                .add(EdgeInsets.only(top: 5)),
                                        child: IntrinsicHeight(
                                          child: Row(
                                            children: [
                                              common.currentUserType ==
                                                      CurrentUserType
                                                          .decentralized
                                                  ? SizedBox(
                                                      width: 190,
                                                      child:
                                                          TopSelectBlockChainBar(
                                                              onSegmentChosen:
                                                                  (BlockChainType
                                                                      type) {
                                                                setState(() {
                                                                  common.topSelectedBlockchainType =
                                                                      type;
                                                                });
                                                              },
                                                              currentSelection:
                                                                  common
                                                                      .topSelectedBlockchainType,
                                                              topBarType:
                                                                  TopSelectBarType
                                                                      .enable))
                                                  : Container(),
                                              Spacer(flex: 1),
                                              IconButton(
                                                icon: SizedBox(
                                                    child: common.showBalance
                                                        ? Image(
                                                            image: AssetImage(
                                                                'images/eyeOpen.png',
                                                                package:
                                                                    'euruswallet'),
                                                            color: common
                                                                .getBackGroundColor())
                                                        : Image(
                                                            image: AssetImage(
                                                                'images/eyeClose.png',
                                                                package:
                                                                    'euruswallet'),
                                                            color: common
                                                                .getBackGroundColor()),
                                                    width: 24),
                                                color:
                                                    common.getBackGroundColor(),
                                                onPressed: () {
                                                  setState(() {
                                                    common.showBalance =
                                                        !common.showBalance;
                                                    NormalStorageKit().setValue(
                                                        common.showBalance ==
                                                                true
                                                            ? '1'
                                                            : '0',
                                                        'defaultShowAssetBalance_${common.encryptedAddress ?? ''}');
                                                  });
                                                },
                                              ),
                                              IconButton(
                                                onPressed: () async {
                                                  await common.pushPage(
                                                      page: EditAssetsListPage(
                                                          userSuffix: common
                                                                  .encryptedAddress ??
                                                              ''),
                                                      context: context);
                                                  common.refreshAssetList.sink
                                                      .add(true);
                                                },
                                                icon: Image.asset(
                                                  'images/icn_edit_asset.png',
                                                  package: 'euruswallet',
                                                  width: 25,
                                                  height: 25,
                                                  color: common
                                                      .getBackGroundColor(),
                                                ),
                                              ),
                                            ],
                                          ),
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              )
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            )
          ])),
      DappBrowserHomePage(),
      CoinPriceRankingPage(),
      SettingPage(
        userSuffix: common.encryptedAddress ?? '',
        backupSeedPhraseFnc: () async {
          await Navigator.of(context).push(
            MaterialPageRoute(
              builder: (_materialPageRouteBuilderContext) =>
                  DecentralizedWalletBasePage(
                appBarTitle: Text(
                  "BACKUP_SEEDPHRASE_PAGE.MAIN_TITLE".tr(),
                ),
                body: BackupSeedPhrasesHomePage(),
              ),
            ),
          );
        },
      ),
    ];
  }

  void showRewardIfNeeded() async {
    final rewardSchemeResponse = await api.getUserMarketingRewardScheme();
    api.getUserMarketingRewardList().then((value) async {
      setState(() {
        common.rewardSchemes = rewardSchemeResponse.schemes;
        common.rewardedList = value.list;
      });
      if ((common.rewardSchemes ?? []).isNotEmpty && value.list.isEmpty) {
        await showDialog(
          context: context,
          builder: (context) {
            return RewardSchemeDialog();
          },
        );
      }
    }).catchError((e) {});
  }
}
