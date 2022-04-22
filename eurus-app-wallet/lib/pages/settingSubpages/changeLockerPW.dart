import 'package:app_security_kit/app_security_kit.dart';
import 'package:biometric_storage/biometric_storage.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/common/user_profile.dart';
import 'package:euruswallet/pages/settingSubpages/cardContainer.dart';
import 'package:euruswallet/pages/settingSubpages/settingAppBar.dart';

class ChangeLockerPWPage extends StatefulWidget {
  ChangeLockerPWPage({
    required this.userSuffix,
  });

  final String userSuffix;

  @override
  _ChangeLockerPWPageState createState() => _ChangeLockerPWPageState();
}

class _ChangeLockerPWPageState extends State<ChangeLockerPWPage> {
  final TextEditingController _curPwTc = TextEditingController();
  final TextEditingController _newPwTc = TextEditingController();
  final TextEditingController _confirmPwTc = TextEditingController();

  bool _isCurPwMasked = false;
  bool _isNewPwMasked = false;
  bool _isConfirmPwMasked = false;

  final _pwForm = GlobalKey<FormState>();

  bool _onUpdatePW = false;

  bool? _curPWValid;
  bool? _newPWValid;
  bool? _confirmPWValid;

  String? _curPWErrorMsg;
  String? _newPWErrorMsg;
  String? _confirmPWErrorMsg;

  late NormalStorageKit nStorage;
  late SecureStorageKit sStorage;

  String get _newUserSuffix => _switchValEncryption(widget.userSuffix);
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    nStorage = NormalStorageKit();
    sStorage = SecureStorageKit();

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SettingAppBar(true),
      body: Container(
        child: SingleChildScrollView(
          child: SafeArea(
            child: CardContainer(
              'CHANGE_LOCKER_PW.TITLE'.tr(),
              Container(
                padding: EdgeInsets.symmetric(vertical: 25),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Row(
                      children: [
                        Spacer(flex: 4),
                        Expanded(
                          flex: 3,
                          child: Image(
                            image: AssetImage(
                              'images/WLP.png',
                              package: 'euruswallet',
                            ),
                          ),
                        ),
                        Spacer(flex: 4),
                      ],
                    ),
                    Padding(padding: EdgeInsets.only(bottom: 35)),
                    Form(
                      key: _pwForm,
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          _inputRow(
                            'CHANGE_LOCKER_PW.CURRENT_PW.LABEL'.tr(),
                            _curPwTc,
                            hintText:
                                'CHANGE_LOCKER_PW.CURRENT_PW.PLACEHOLDER'.tr(),
                            vFnc: (v) {
                              if (isEmptyString(string: v))
                                return 'CHANGE_LOCKER_PW.ERROR.EMPTY_CUR_PW'
                                    .tr();

                              if (_curPWValid == false) return _curPWErrorMsg;

                              return null;
                            },
                            obscureText: !_isCurPwMasked,
                            suffixIconOnPressed: () {
                              setState(() {
                                _isCurPwMasked = !_isCurPwMasked;
                              });
                            },
                          ),
                          _inputRow(
                            'CHANGE_LOCKER_PW.NEW_PW.LABEL'.tr(),
                            _newPwTc,
                            hintText:
                                'CHANGE_LOCKER_PW.NEW_PW.PLACEHOLDER'.tr(),
                            vFnc: (v) {
                              if (isEmptyString(string: v))
                                return 'CHANGE_LOCKER_PW.ERROR.EMPTY_NEW_PW'
                                    .tr();

                              if (_newPWValid == false) return _newPWErrorMsg;

                              return null;
                            },
                            obscureText: !_isNewPwMasked,
                            suffixIconOnPressed: () {
                              setState(() {
                                _isNewPwMasked = !_isNewPwMasked;
                              });
                            },
                          ),
                          _inputRow(
                            'CHANGE_LOCKER_PW.CONFIRM_PW.LABEL'.tr(),
                            _confirmPwTc,
                            hintText:
                                'CHANGE_LOCKER_PW.CONFIRM_PW.PLACEHOLDER'.tr(),
                            vFnc: (v) {
                              if (isEmptyString(string: v))
                                return 'CHANGE_LOCKER_PW.ERROR.EMPTY_CONFIRM_PW'
                                    .tr();

                              if (_confirmPWValid == false)
                                return _confirmPWErrorMsg;

                              return null;
                            },
                            obscureText: !_isConfirmPwMasked,
                            suffixIconOnPressed: () {
                              setState(() {
                                _isConfirmPwMasked = !_isConfirmPwMasked;
                              });
                            },
                          ),
                        ],
                      ),
                    ),
                    Container(
                        padding: EdgeInsets.only(left: 35, right: 35, top: 25),
                        child: SubmitButton(
                            btnController: btnController,
                            label: 'COMMON.CONFIRM'.tr(),
                            onPressed: _changePW)),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _inputRow(
    String title,
    TextEditingController tc, {
    String? hintText,
    String? errorMsg,
    String? Function(String?)? vFnc,
    bool obscureText = true,
    void Function()? suffixIconOnPressed,
  }) {
    final _defaultTextFieldInputDecoration = InputDecoration(
      filled: true,
      fillColor: FXColor.lightGreyTextColor,
      hintStyle: Theme.of(context)
          .textTheme
          .subtitle1
          ?.apply(color: Theme.of(context).hintColor),
      border: OutlineInputBorder(
        borderSide: BorderSide.none,
        borderRadius: FXUI.cricleRadius,
      ),
      contentPadding: EdgeInsets.all(16),
    );

    double _hPadding = MediaQuery.of(context).size.width / 13;

    return Container(
      width: double.infinity,
      padding:
          EdgeInsets.symmetric(horizontal: _hPadding > 35 ? 35 : _hPadding),
      child: Column(
        children: [
          Container(
            padding: getEdgeInsetsSymmetric(),
            alignment: Alignment(-1, 0),
            child: Text(
              title,
              style: Theme.of(context)
                  .textTheme
                  .bodyText2
                  ?.apply(color: FXColor.lightGray),
            ),
          ),
          TextFormField(
            decoration: _defaultTextFieldInputDecoration.copyWith(
              hintText: hintText ?? '',
              hintStyle: FXUI.normalTextStyle.copyWith(
                color: FXColor.centralizedGrayTextColor,
                fontSize: 12,
              ),
              errorText: errorMsg ?? '',
              suffixIcon: IconButton(
                onPressed: suffixIconOnPressed,
                icon: Image.asset(
                  !obscureText ? 'images/eyeClose.png' : 'images/eyeOpen.png',
                  package: 'euruswallet',
                  width: 16,
                  height: 16,
                  color: common.getBackGroundColor(),
                ),
              ),
            ),
            controller: tc,
            obscureText: obscureText,
            autovalidateMode: AutovalidateMode.onUserInteraction,
            onChanged: (v) => _clearErrorMsg(),
            validator: vFnc,
          )
        ],
      ),
    );
  }

  Future<Null> _changePW() async {
    if (_onUpdatePW) {
      btnController.reset();
      return;
    }

    if (_pwForm.currentState == null || !_pwForm.currentState!.validate()) {
      btnController.reset();
      return;
    }

    // Check if current password is correct
    PasswordEncryptHelper pwHelper = PasswordEncryptHelper(
      password: _curPwTc.text,
    );
    String acAddress = pwHelper.decryptWPed(widget.userSuffix) ?? '';

    if (isEmptyString(string: acAddress)) {
      setState(() {
        _curPWValid = false;
        _curPWErrorMsg = 'COMMON_ERROR.INCORRECT_PW'.tr();
      });
      btnController.reset();
      return;
    }

    int newPwMatch = _newPwTc.text.compareTo(_confirmPwTc.text);
    if (newPwMatch != 0) {
      setState(() {
        _confirmPWValid = false;
        _confirmPWErrorMsg = 'COMMON_ERROR.PW_INCONSISTENT'.tr();
      });
      btnController.reset();
      return;
    }

    int diffNewPw = _newPwTc.text.compareTo(_curPwTc.text);
    if (diffNewPw == 0) {
      setState(() {
        _newPWValid = false;
        _newPWErrorMsg = 'CHANGE_LOCKER_PW.ERROR.DIFF_NEW_PW'.tr();
      });
      btnController.reset();
      return;
    }

    setState(() {
      _onUpdatePW = true;
    });

    /// Content migeation needed
    ///
    /// 1. User Profile
    /// 2. Account
    /// 3. Asset List
    /// 4. Biometric storage
    ///    1. Unlock wallet (login) using biometric
    ///    2. Confirm transaction using biometric
    /// 5. Asset Privacy (Show/Hide asset balance on home by default)
    await _profileMigration();
    await _acountMigration();
    await _tokenListMigration();
    await _bioStorageMigration();
    await _assetPrivacyMigration();

    setState(() {
      _onUpdatePW = false;
    });

    // await _successDialog().then((value) => Navigator.pop(context, 'done'));
    await _successDialog();
    btnController.reset();
    // Navigator.pop(context, 'done');
  }

  Future<void> _profileMigration() async {
    List<UserProfile> userProfiles = await getLocalAcs();
    UserProfile userProfile = userProfiles
        .firstWhere((element) => element.encryptedAddress == widget.userSuffix);

    await setAcToLocal(userProfile, delete: true);

    userProfile.encryptedAddress = _newUserSuffix;
    userProfile.encryptedPrivateKey =
        _switchValEncryption(userProfile.encryptedPrivateKey);
    userProfile.bioAuthValidTime = "-1";
    userProfile.bioTxAuthValidTime = "-1";
    String newSeedPhrase = _updateUriKeyVal(
      userProfile.mnemonicSeedPhrases ?? "",
      widget.userSuffix,
      _switchValEncryption(widget.userSuffix),
      switchEncryption: true,
    );
    userProfile.mnemonicSeedPhrases = newSeedPhrase;

    await setAcToLocal(userProfile);
  }

  Future<bool> _acountMigration() async {
    String valKey = await common.prefix + 'accounts';

    String ac = await sStorage.readValue(valKey) ?? '';

    await sStorage.setValue(
      _updateUriKeyVal(ac, widget.userSuffix, _newUserSuffix),
      valKey,
    );

    return true;
  }

  Future<bool> _bioStorageMigration() async {
    bool migrationSuccess = true;

    // Check if biometric authentication is available
    final bioAuth = await BiometricStorage().canAuthenticate();

    // Skip biostorage migration if biometric authentication is not available
    if (bioAuth != CanAuthenticateResponse.success) return migrationSuccess;

    var __prefix = await common.prefix;

    /// Login Bio Key: [__prefix][widget.userSuffix]_authenticated
    final _authStorageFile = await BiometricStorage().getStorage(
      '${__prefix + widget.userSuffix}_authenticated',
      options: StorageFileInitOptions(
        authenticationValidityDurationSeconds: Platform.isIOS ? 0 : 30,
      ),
    );

    await _authStorageFile.delete().then((v) async {
      if (v == true) {
        final _newAuthStorageFile = await BiometricStorage().getStorage(
          '${__prefix + _newUserSuffix}_authenticated',
          options: StorageFileInitOptions(
            authenticationValidityDurationSeconds: Platform.isIOS ? 0 : 30,
          ),
        );

        await _newAuthStorageFile.write(_newPwTc.text).catchError((e, t) {
          migrationSuccess = false;
          print('bioauth write new val $e - $t');
        });
      }
    }).catchError((e, t) {
      print('Empyt - Login Bio migration abort - Tread as migration success');
    });

    /// Tx Bio Key: [__prefix][widget.userSuffix]_tx_authenticated
    final _txAuthStorageFile = await BiometricStorage().getStorage(
      '${__prefix + widget.userSuffix}_tx_authenticated',
      options: StorageFileInitOptions(
        authenticationValidityDurationSeconds: Platform.isIOS ? 0 : 30,
      ),
    );

    await _txAuthStorageFile.delete().then((v) async {
      if (v == true) {
        final _newTxAuthStorageFile = await BiometricStorage().getStorage(
          '${__prefix + _newUserSuffix}_tx_authenticated',
          options: StorageFileInitOptions(
            authenticationValidityDurationSeconds: Platform.isIOS ? 0 : 30,
          ),
        );

        await _newTxAuthStorageFile.write(_newPwTc.text).catchError((e, t) {
          migrationSuccess = false;
          print('bioauth write new val $e - $t');
        });
      }
    }).catchError((e, t) {
      print('Empyt - Tx Bio migration abort - Tread as migration success');
    });
    return migrationSuccess;
  }

  Future<bool> _tokenListMigration() async {
    String? tkenList =
        await nStorage.readValue('assetsList_${widget.userSuffix}');

    if (tkenList != null) {
      await nStorage.deleteValue('assetsList_${widget.userSuffix}');
      await nStorage.setValue(tkenList, 'assetsList_$_newUserSuffix');
    }

    return true;
  }

  Future<bool> _assetPrivacyMigration() async {
    var defaultShowAB = await NormalStorageKit()
        .readValue('defaultShowAssetBalance_${widget.userSuffix}');

    if (defaultShowAB != null) {
      await nStorage
          .deleteValue('defaultShowAssetBalance_${widget.userSuffix}');
      await nStorage.setValue(
          defaultShowAB, 'defaultShowAssetBalance_$_newUserSuffix');
    }

    return true;
  }

  String _updateUriKeyVal(
    String uri,
    String orgKey,
    String newKey, {
    bool? switchEncryption,
    String? newVal,
  }) {
    var orgVal = Uri(query: uri).queryParameters[orgKey];

    if (orgVal == null) return uri;

    var uriParam = {
      ...Uri(query: uri).queryParameters,
      newKey: switchEncryption == true
          ? _switchValEncryption(newVal ?? orgVal)
          : newVal ?? orgVal
    }..remove(orgKey);

    return Uri(queryParameters: uriParam).query;
  }

  String _switchValEncryption(String v) {
    var decodedVal =
        PasswordEncryptHelper(password: _curPwTc.text).decryptWPed(v) ?? '';
    return PasswordEncryptHelper(password: _newPwTc.text)
        .encryptWPwd(decodedVal);
  }

  void _clearErrorMsg() {
    setState(() {
      _curPWValid = null;
      _newPWValid = null;
      _confirmPWValid = null;
      _curPWErrorMsg = null;
      _newPWErrorMsg = null;
      _confirmPWErrorMsg = null;
    });
  }

  Future<void> _successDialog() async {
    double _hPadding = MediaQuery.of(context).size.width / 13;
    double _wPadding = MediaQuery.of(context).size.height / 30;

    return showDialog(
      context: context,
      builder: (_) {
        return Dialog(
          backgroundColor: Colors.transparent,
          insetPadding: EdgeInsets.all(10),
          child: Container(
            width: double.infinity,
            padding: EdgeInsets.symmetric(
                vertical: _hPadding > 25 ? 25 : _hPadding,
                horizontal: _wPadding > 17 ? 17 : _wPadding),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: FXUI.cricleRadius,
            ),
            child: Stack(
              children: [
                SingleChildScrollView(
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Center(
                        child: Text(
                          'CHANGE_LOCKER_PW.CHANGE_SUCCESS'.tr(),
                          style: FXUI.normalTextStyle.copyWith(
                            fontWeight: FontWeight.bold,
                            fontSize: 18,
                          ),
                        ),
                      ),
                      Padding(padding: EdgeInsets.only(top: 70)),
                      SizedBox(
                        width: MediaQuery.of(context).size.width / 4.5,
                        child: Image(
                          image: AssetImage(
                              "images/${!isCentralized() ? 'decenTickIcon' : 'tickIcon'}.png",
                              package: 'euruswallet'),
                        ),
                      ),
                      Padding(
                        padding: EdgeInsets.symmetric(vertical: 15),
                        child: Text(
                          'COMMON.SUCCESS'.tr(),
                          style: FXUI.normalTextStyle.copyWith(
                            fontWeight: FontWeight.bold,
                            fontSize: 23,
                          ),
                        ),
                      ),
                      Padding(
                        padding: EdgeInsets.symmetric(horizontal: 10),
                        child: Text(
                          'CHANGE_LOCKER_PW.CHANGE_SUCCESS_DESC'.tr(),
                          textAlign: TextAlign.center,
                          style: FXUI.normalTextStyle
                              .copyWith(color: FXColor.textGray),
                        ),
                      ),
                      Padding(
                        padding: EdgeInsets.only(top: 50),
                        child: TextButton(
                          style: TextButton.styleFrom(
                            padding: EdgeInsets.zero,
                          ),
                          onPressed: () {
                            Navigator.pop(_);
                            Navigator.of(context).popUntil((route) {
                              if (route.isFirst) {
                                (route.settings.arguments
                                    as Map)['shouldShowAuthPopUp'] = true;
                                return true;
                              } else {
                                return false;
                              }
                            });
                          },
                          child: Container(
                            decoration: BoxDecoration(
                              color: FXColor.mainDeepBlueColor,
                              borderRadius: FXUI.cricleRadius,
                            ),
                            padding: EdgeInsets.all(15),
                            width: double.infinity,
                            child: Center(
                              child: Text(
                                // 'Back',
                                'CHANGE_LOCKER_PW.RELOGIN'.tr(),
                                style: FXUI.normalTextStyle.copyWith(
                                    color: Colors.white, fontSize: 16),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                Positioned(
                  right: -13,
                  top: -12,
                  child: IconButton(
                    icon: Icon(
                      Icons.close,
                      color: Colors.black.withOpacity(0.5),
                    ),
                    onPressed: () => Navigator.pop(_),
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}
