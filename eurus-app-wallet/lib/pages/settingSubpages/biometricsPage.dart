import 'package:app_security_kit/app_security_kit.dart';
import 'package:biometric_storage/biometric_storage.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/common/user_profile.dart';
import 'package:euruswallet/commonUI/walletLockerPWDialog.dart';
import 'package:euruswallet/pages/settingSubpages/cardContainer.dart';
import 'package:euruswallet/pages/settingSubpages/settingAppBar.dart';
import 'package:flutter/cupertino.dart';

class BiometricsPage extends StatefulWidget {
  BiometricsPage({
    required this.userSuffix,
  });

  final String userSuffix;

  _BiometricsPageState createState() => _BiometricsPageState();
}

class _BiometricsPageState extends State<BiometricsPage> {
  bool _bioLogin = common.currentUserProfile?.enableBiometric ?? false;
  bool _bioTransaction = common.currentUserProfile?.enableTxBiometric ?? false;
  bool _bioUnsupported = false;
  BiometricStorageFile? _loginBioStorage;
  BiometricStorageFile? _txBioStorage;

  @override
  void initState() {
    _getBioAuthStatus();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    var _themeColor = common.getBackGroundColor();

    return Scaffold(
      backgroundColor: FXColor.veryLightGreyTextColor,
      appBar: SettingAppBar(true),
      body: Container(
        child: SafeArea(
          child: CardContainer(
            'BIOMETRICS_PAGE.TITLE'.tr(),
            Container(
              padding: EdgeInsets.symmetric(horizontal: 35, vertical: 25),
              alignment: Alignment(-1, 0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'BIOMETRICS_PAGE.DESC'.tr(),
                    textAlign: TextAlign.start,
                    style: FXUI.normalTextStyle
                        .copyWith(fontWeight: FontWeight.w600, fontSize: 16),
                  ),
                  Padding(
                    padding: EdgeInsets.only(top: 55),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          children: [
                            Text('BIOMETRICS_PAGE.LOGIN'.tr()),
                            Padding(padding: EdgeInsets.only(left: 10)),
                            CupertinoSwitch(
                              value: _bioLogin,
                              onChanged: (e) => _toggleBioAuth(e, 'login'),
                              activeColor: _themeColor,
                            )
                          ],
                        ),
                        Row(
                          children: [
                            Text('BIOMETRICS_PAGE.TRANSACTION'.tr()),
                            Padding(padding: EdgeInsets.only(left: 10)),
                            CupertinoSwitch(
                              value: _bioTransaction,
                              onChanged: (e) => _toggleBioAuth(e, 'tx'),
                              activeColor: _themeColor,
                            )
                          ],
                        )
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _getBioAuthStatus() async {
    CanAuthenticateResponse bioAuth =
        await BiometricStorage().canAuthenticate();
    if (bioAuth != CanAuthenticateResponse.success) {
      setState(() {
        _bioUnsupported = false;
      });
      if (bioAuth == CanAuthenticateResponse.errorNoBiometricEnrolled) {
        common.showBiometricNotEnable(context: context);
      }
      return;
    }

    // init all bio storage
    final _prefix = await common.prefix;
    _loginBioStorage = await BiometricStorage().getStorage(
      '${_prefix + widget.userSuffix}_authenticated',
      options: StorageFileInitOptions(
        authenticationValidityDurationSeconds: Platform.isIOS ? 0 : 30,
      ),
    );
    _txBioStorage = await BiometricStorage().getStorage(
      '${_prefix + widget.userSuffix}_tx_authenticated',
      options: StorageFileInitOptions(
        authenticationValidityDurationSeconds: Platform.isIOS ? 0 : 30,
      ),
    );
  }

  Future<Null> _toggleBioAuth(bool status, String type) async {
    if (_bioUnsupported) return null;
    status == true ? await _turnOnBioAuth(type) : await _turnOffBioAuth(type);
  }

  Future<void> _turnOnBioAuth(String type) async {
    bool setBioSuccess = true;
    String? pw = await _requestLockerPWDialog();

    if (isEmptyString(string: pw)) return;

    BiometricStorageFile? activeBioStorage =
        type == 'login' ? _loginBioStorage : _txBioStorage;
    await activeBioStorage?.delete().catchError((e) {
      print("nothing to be delete");
      return null;
    });
    await activeBioStorage?.write(pw ?? '').catchError((e) {
      /// Only Android will catch this error
      print('error in write by bio');
      setBioSuccess = false;
    });

    if (Platform.isIOS) {
      await activeBioStorage?.read().catchError((e) {
        print('error in read by bio');
        setBioSuccess = false;
      });
    }

    if (!setBioSuccess) {
      await activeBioStorage?.delete();
      return;
    }

    final UserProfile? user =
        new UserProfile.fromJson(common.currentUserProfile!.toJson());
    if (user != null) {
      var newBioValidUntil = DateTime.now()
          .add(Duration(days: 7))
          .millisecondsSinceEpoch
          .toString();
      if (type == 'login') {
        user.enableBiometric = true;
        user.bioAuthValidTime = newBioValidUntil;
      } else {
        user.enableTxBiometric = true;
        user.bioTxAuthValidTime = newBioValidUntil;
      }
      _updateUserBioAuthStatus(user);
      setState(() {
        type == 'login' ? _bioLogin = true : _bioTransaction = true;
      });
    }
  }

  Future<void> _turnOffBioAuth(String type) async {
    // Verify user with PW, using Wallet locker dialog
    final lockerInput = await _requestLockerPWDialog();
    if (isEmptyString(string: lockerInput)) return;

    final UserProfile? user =
        new UserProfile.fromJson(common.currentUserProfile!.toJson());
    if (user != null) {
      if (type == 'login') {
        await _loginBioStorage?.delete().catchError((e) {
          print('Empty - deleting in bio login');
        });
        user.enableBiometric = false;
        user.bioAuthValidTime = "-1";
      } else if (type == 'tx') {
        await _txBioStorage?.delete().catchError((e) {
          print('Empty - deleting in bio transaction');
        });
        user.enableTxBiometric = false;
        user.bioTxAuthValidTime = "-1";
      }
      _updateUserBioAuthStatus(user);
      setState(() {
        type == 'login' ? _bioLogin = false : _bioTransaction = false;
      });
    }
  }

  Future<void> _updateUserBioAuthStatus(UserProfile user) async {
    List<UserProfile> userProfiles = await getLocalAcs();
    UserProfile userProfile = userProfiles.firstWhere(
        (element) => element.encryptedAddress == user.encryptedAddress);
    await setAcToLocal(userProfile, delete: true);
    await setAcToLocal(user);
    common.currentUserProfile = user;
  }

  Future<String?> _requestLockerPWDialog() async {
    final _tc = TextEditingController();
    final _submit = (_accountPrivateKey, tec) async {
      final _privateKey = PasswordEncryptHelper(password: _tc.text)
          .decryptWPed(_accountPrivateKey ?? '');

      return _privateKey;
    };

    String? decryptedString = await Navigator.of(context)
        .push(PageRouteBuilder(
          fullscreenDialog: true,
          opaque: false,
          pageBuilder: (pageBuilderContext, animation, secondaryAnimation) =>
              WalletLockerPWDialog(
                  themeColor: common.getBackGroundColor(),
                  textEditingController: _tc,
                  decenUserkey: common.currentUserProfile!.encryptedPrivateKey,
                  submitFnc: _submit,
                  tryBioAuthFnc: null),
        ))
        .then((value) => value);

    return isEmptyString(string: decryptedString) ? null : _tc.text;
  }
}
