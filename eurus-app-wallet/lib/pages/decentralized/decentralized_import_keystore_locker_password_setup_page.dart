import 'package:app_authentication_kit/mnemonic_kit.dart';
import 'package:app_security_kit/password_encrypt_helper.dart';
import 'package:app_storage_kit/secure_storage.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:rounded_loading_button/rounded_loading_button.dart';
import '../../common/biometric_authentication_helper.dart';
import '../../common/commonMethod.dart';
import '../../common/user_profile.dart';
import 'decentralized_import_wallet_page.dart';

class DecentralizedImportKeyStoreLockerPasswordSetupPage
    extends StatefulWidget {
  final ImportWalletType importWalletType;
  final String? mnemonicSeedPhrase;
  final String? privateKey;
  final bool isSeedPhraseImported;

  const DecentralizedImportKeyStoreLockerPasswordSetupPage({
    Key? key,
    required this.importWalletType,
    this.mnemonicSeedPhrase,
    this.privateKey,
    this.isSeedPhraseImported = false,
  }) : super(key: key);

  @override
  _DecentralizedImportKeyStoreLockerPasswordSetupState createState() =>
      _DecentralizedImportKeyStoreLockerPasswordSetupState();
}

class _DecentralizedImportKeyStoreLockerPasswordSetupState
    extends State<DecentralizedImportKeyStoreLockerPasswordSetupPage>
    with WidgetsBindingObserver {
  bool isPasswordMasked = false;
  bool isConfirmPasswordMasked = false;
  final _textEditingController0 = TextEditingController();
  final _textEditingController1 = TextEditingController();
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  bool _isFormValid = false;

  late bool _canSupportBiometric;
  bool _willEnableBiometricsAuth = false;
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    super.initState();
    checkCanSupportBiometric(() {
      if (!_canSupportBiometric) {
        common.showBiometricNotEnable(
          context: context,
          currentUserType: CurrentUserType.decentralized,
        );
      }
    });
    WidgetsBinding.instance?.addObserver(this);
  }

  @override
  void dispose() {
    WidgetsBinding.instance?.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);

    if (state == AppLifecycleState.resumed) checkCanSupportBiometric(() {});
  }

  void checkCanSupportBiometric(Function failureHandler) async {
    _canSupportBiometric = await BiometricAuthenticationHelper()
        .canSupportsBiometricAuthenticated();
    failureHandler();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Padding(
          padding: EdgeInsets.symmetric(vertical: 30),
          child: FractionallySizedBox(
            widthFactor: .25,
            child: Image(
              image: AssetImage('images/WLP.png', package: 'euruswallet'),
            ),
          ),
        ),
        Expanded(
          flex: 1,
          child: SingleChildScrollView(
            child: Form(
              key: _formKey,
              onChanged: () {
                if (_formKey.currentState != null)
                  setState(() {
                    _isFormValid = _formKey.currentState!.validate();
                  });
              },
              child: Column(
                children: [
                  ListTile(
                      title: Padding(
                        child: Text(
                          "CREATE_LOCKER_PAGE.NEW_PW.LABEL".tr(),
                          style: Theme.of(context)
                              .textTheme
                              .bodyText2
                              ?.apply(color: FXColor.lightGray),
                        ),
                        padding: getEdgeInsetsSymmetric(),
                      ),
                      subtitle: TextFormField(
                        validator: (value) => value ==
                                    _textEditingController1.text &&
                                !isEmptyString(string: value)
                            ? null
                            : !isEmptyString(string: value)
                                ? 'COMMON_ERROR.PW_INCONSISTENT'.tr()
                                : 'CREATE_LOCKER_PAGE.ERROR.EMPTY_CONFIRM_PW'
                                    .tr(),
                        decoration:
                            FXUI.defaultTextFieldInputDecoration.copyWith(
                          hintText:
                              "CREATE_LOCKER_PAGE.NEW_PW.PLACEHOLDER".tr(),
                          hintStyle: FXUI.normalTextStyle.copyWith(
                              color: FXColor.centralizedGrayTextColor),
                          suffixIcon: IconButton(
                            onPressed: () {
                              setState(() {
                                isPasswordMasked = !isPasswordMasked;
                              });
                            },
                            icon: Image.asset(
                              isPasswordMasked
                                  ? 'images/eyeClose.png'
                                  : 'images/eyeOpen.png',
                              package: 'euruswallet',
                              width: 16,
                              height: 16,
                              color: FXColor.mainDeepBlueColor,
                            ),
                          ),
                        ),
                        obscureText: !isPasswordMasked,
                        controller: _textEditingController0,
                      )),
                  ListTile(
                      title: Padding(
                          child: Text(
                              "CREATE_LOCKER_PAGE.CONFIRM_PW.LABEL".tr(),
                              style: Theme.of(context)
                                  .textTheme
                                  .bodyText2
                                  ?.apply(color: FXColor.lightGray)),
                          padding: getEdgeInsetsSymmetric()),
                      subtitle: Padding(
                          padding: EdgeInsets.only(
                              bottom: MediaQuery.of(context).viewInsets.bottom),
                          child: TextFormField(
                            decoration:
                                FXUI.defaultTextFieldInputDecoration.copyWith(
                              hintText:
                                  "CREATE_LOCKER_PAGE.CONFIRM_PW.PLACEHOLDER"
                                      .tr(),
                              hintStyle: FXUI.normalTextStyle.copyWith(
                                  color: FXColor.centralizedGrayTextColor),
                              suffixIcon: IconButton(
                                onPressed: () {
                                  setState(() {
                                    isConfirmPasswordMasked =
                                        !isConfirmPasswordMasked;
                                  });
                                },
                                icon: Image.asset(
                                  isConfirmPasswordMasked
                                      ? 'images/eyeClose.png'
                                      : 'images/eyeOpen.png',
                                  package: 'euruswallet',
                                  width: 16,
                                  height: 16,
                                  color: FXColor.mainDeepBlueColor,
                                ),
                              ),
                            ),
                            obscureText: !isConfirmPasswordMasked,
                            controller: _textEditingController1,
                            validator: (value) => value ==
                                        _textEditingController0.text &&
                                    !isEmptyString(string: value)
                                ? null
                                : !isEmptyString(string: value)
                                    ? 'COMMON_ERROR.PW_INCONSISTENT'.tr()
                                    : 'CREATE_LOCKER_PAGE.ERROR.EMPTY_CONFIRM_PW'
                                        .tr(),
                          ))),
                  SwitchListTile.adaptive(
                    value: _willEnableBiometricsAuth,
                    onChanged: (newValue) async {
                      if (_canSupportBiometric) {
                        setState(() {
                          _willEnableBiometricsAuth = newValue;
                        });
                      }
                    },
                    title: Text("CREATE_LOCKER_PAGE.ENABLE_BIO_AUTH".tr(),
                        style: Theme.of(context)
                            .textTheme
                            .subtitle1
                            ?.apply(color: FXColor.lightGray)),
                    contentPadding: EdgeInsets.all(16),
                  ),
                  SubmitButton(
                      loadingSecond: 30,
                      btnController: btnController,
                      label: 'COMMON.CONFIRM'.tr(),
                      buttonBGColor: FXColor.mainDeepBlueColor,
                      onPressed: () async {
                        if (_formKey.currentState != null)
                          setState(() {
                            _isFormValid = _formKey.currentState!.validate();
                          });
                        if (!_formKey.currentState!.validate()) {
                          btnController.reset();
                          return;
                        }
// Generate root key for application to store in local
// Uses this key to derivative Address and Private key
                        String _mnemonicSeedPhrase =
                            widget.mnemonicSeedPhrase ??
                                await compute(genMnemonic, 128);
                        AddressPair? anAddressPair;
                        switch (widget.importWalletType) {
                          case ImportWalletType.mnemonic:
                            final base58 =
                                await compute(genBase58, _mnemonicSeedPhrase) ??
                                    '';
                            anAddressPair =
                                await compute(genAddressPair, base58);
                            common.bitcoinKeyPair = api.getBitcoinKeyPair(
                                mnemonicSeedPhrase: _mnemonicSeedPhrase);
                            // SendBitcoinDetail sendBitcoinDetail = await api.getSendBitcoinDetail(targetAddress: "mhARbrfUzFoepzrAwdzwEVWAi3Wmq22FmR", myAddress:common.bitcoinWallet!.address!, amount: 132000);
                            // api.pushRawTransaction(alice: common.bitcoinKeyPair, sendBitcoinDetail: sendBitcoinDetail);
                            break;
                          case ImportWalletType.privateKey:
                            final privateKey = widget.privateKey ?? '';
                            if (!isEmptyString(string: privateKey))
                              anAddressPair = common
                                  .getAddressPairFromPrivateKey(privateKey);
                            common.bitcoinKeyPair = null;
                            break;
                        }

                        if (anAddressPair == null) return;

                        final password = _textEditingController0.text;

                        final _encrypt = (c) =>
                            PasswordEncryptHelper(password: password)
                                .encryptWPwd(c);
                        final _pwdEncryptedAddress =
                            _encrypt(anAddressPair.address);
                        final __prefix = await common.prefix;
                        final encryptPasswordAddress =
                            __prefix + _pwdEncryptedAddress;

                        await setActiveAccount(_pwdEncryptedAddress,
                            enPrivateKey: _encrypt(anAddressPair.privateKey),
                            enMPhrase: widget.importWalletType ==
                                    ImportWalletType.mnemonic
                                ? widget.mnemonicSeedPhrase == null
                                    ? _encrypt(_mnemonicSeedPhrase)
                                    : null
                                : null);
                        // secure persistence of private key
                        // await _securePersist({k_accountPrivateKeys: Uri(queryParameters: {_pwdEncryptedAddress: _encrypt(anAddressPair.privateKey)}).query});
                        // secure persistence of newly generated unique mnemonic seed phrase if HD wallet is newly created
                        final encryptedAddress =
                            await common.findDefaultAccountEncryptedAddress() ??
                                '';
                        await _handleBiometricAuthentication(
                          encryptPasswordAddress,
                          password,
                          encryptedAddress,
                        );

                        final isUpdateSuccess =
                            await common.updateApiAccessTokenByImportWallet(
                          context: context,
                          addressPair: anAddressPair,
                        );

                        if (isUpdateSuccess) {
                          UserProfile userProfile = UserProfile.fromJson({
                            "userType": CurrentUserType.decentralized,
                            "address": anAddressPair.address.toLowerCase(),
                            "encryptedAddress": _pwdEncryptedAddress,
                            "encryptedPrivateKey":
                                _encrypt(anAddressPair.privateKey),
                            "seedPraseBackuped": widget.isSeedPhraseImported,
                            "decenUserType": widget.isSeedPhraseImported == true
                                ? DecenUserType.imported
                                : DecenUserType.created,
                            "enableBiometric": _willEnableBiometricsAuth,
                            "bioAuthValidTime":
                                _willEnableBiometricsAuth == true
                                    ? DateTime.now()
                                        .add(Duration(days: 7))
                                        .millisecondsSinceEpoch
                                        .toString()
                                    : "-1",
                            "mnemonicSeedPhrases": widget.importWalletType ==
                                        ImportWalletType.mnemonic &&
                                    widget.mnemonicSeedPhrase == null
                                ? Uri(queryParameters: {
                                    encryptedAddress:
                                        _encrypt(_mnemonicSeedPhrase)
                                  }).query
                                : null,
                          });

                          await setAcToLocal(userProfile);

                          common.successMoveToHomePage(
                            userType: CurrentUserType.decentralized,
                            email: common.email ?? '',
                            context: context,
                            loginPassword: password,
                            address: anAddressPair.address,
                            privateKey: anAddressPair.privateKey,
                            userProfile: userProfile,
                            isRegister: widget.mnemonicSeedPhrase == null,
                          );
                          btnController.reset();
                        }
                      })
                ],
                mainAxisSize: MainAxisSize.max,
                mainAxisAlignment: MainAxisAlignment.start,
              ),
            ),
          ),
        ),
      ],
      mainAxisSize: MainAxisSize.max,
    );
  }

  Future<void> _handleBiometricAuthentication(String encryptPasswordAddress,
      String password, String pwdEncryptedAddress) async {
    final _biometricAuthenticationHelper = BiometricAuthenticationHelper();

    if (_willEnableBiometricsAuth) {
      // in order to ensure and enforce that iOS LocalAuthentication Prompt for user auth permission
      if (Platform.isIOS) {
        final _deleteResult = await _biometricAuthenticationHelper
            .canPersistWithBiometricSecurely(encryptPasswordAddress,
                delete: true);
        print(
            "Platform.isIOS SecItemDelete Flutter MethodChannel.invokeMethod StorageCallback result in Future $_deleteResult");
        if (isEmptyString(string: _deleteResult))
          print("iOS KeyChain SecItemDelete errSecItemNotFound");
        await _biometricAuthenticationHelper.canPersistWithBiometricSecurely(
            {encryptPasswordAddress: password});
      }
      final _bioAuth = await _biometricAuthenticationHelper
          .canPersistWithBiometricSecurely(
              {encryptPasswordAddress: password}).catchError((e) {
        return null;
      });
      if (isEmptyString(string: _bioAuth)) {
        return;
      }
      // Bit Mask enum flags definitions: {0: OFF, 1: ON_login, 2: ON_transaction, 4: other_definitions, 8: other_definitions}
      // e.g. both login and transaction all flags OFF == 0, only login flags ON == 1, only transaction flags ON == 2, both login and transaction flags ON == 3 (i.e. 1+2)
      // p.s. in querystring, use `pwdEncryptedAddress` for account unique key might be enough

      // for any newly Created Wallet Locker, should always ensure a clean state of Biometric UI for Transaction Auth, by deleting / removing storage, regardless of if there was any prior biometric setup for the same wallet account
      await _biometricAuthenticationHelper.canPersistTxWithBiometricSecurely(
          encryptPasswordAddress,
          delete: true);

      // final _bioAuth = await _biometricAuthenticationHelper.canPersistWithBiometricSecurely({k_encryptPasswordAddress: password}, deleteBeforeWriteTwiceEnforceUpdateForiOSPrompt: Platform.isIOS);
      // TODO - fork Flutter Biometric_Storage, CRUD call be purely CRUD, not implicitly update / overwrite existing conflicting key item while commanding only to create / insert an item; complex operation command can accept optional args or instruction or handler or callback to do the update / overwrite operation if key already exist and create / insert error
      // final _bioAuthRead = await _biometricAuthenticationHelper.canPersistWithBiometricSecurely(k_encryptPasswordAddress);
      // TODO - fork Flutter Biometric_Storage, custom Platform System Biometric / Local Authentication prompts text string
      // TODO - fork Flutter Biometric_Storage, custom enum LAPolicy.deviceOwnerAuthenticationWithBiometrics
      print("_willEnableBiometricsAuth $_bioAuth");
    } else if (await _biometricAuthenticationHelper
        .canSupportsBiometricAuthenticated()) {
      final _deleteResult = await _biometricAuthenticationHelper
          .canPersistWithBiometricSecurely(encryptPasswordAddress,
              delete: true);
      print(
          "shall not enable biometrics auth Platform BioStorage Flutter MethodChannel.invokeMethod StorageCallback result in Future $_deleteResult");
      if (isEmptyString(string: _deleteResult)) print("prior item not found");

      // Bit Mask enum flags definitions: {0: OFF, 1: ON_login, 2: ON_transaction, 4: other_definitions, 8: other_definitions}
      // e.g. both login and transaction all flags OFF == 0, only login flags ON == 1, only transaction flags ON == 2, both login and transaction flags ON == 3 (i.e. 1+2)
      // p.s. in querystring, use `_pwdEncryptedAddress` for account unique key might be enough
      // for any newly Created Wallet Locker, should always ensure a clean state of Biometric UI for Transaction Auth, by deleting / removing storage, regardless of if there was any prior biometric setup for the same wallet account
      await _biometricAuthenticationHelper.canPersistTxWithBiometricSecurely(
          encryptPasswordAddress,
          delete: true);
    }
  }
}
