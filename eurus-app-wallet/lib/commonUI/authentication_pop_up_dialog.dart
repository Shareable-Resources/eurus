import 'package:app_authentication_kit/mnemonic_kit.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/biometric_authentication_helper.dart';
import 'package:euruswallet/common/user_profile.dart';
import 'package:euruswallet/model/loginBySignModel.dart';
import '../common/commonMethod.dart';
import '../pages/centralized/cenForgetLoginPwPage.dart';

class AuthenticationPopUpDialog extends StatefulWidget {
  final CurrentUserType currentUserType;
  final String? displayUserName;
  final String encryptedAddress;

  const AuthenticationPopUpDialog({
    Key? key,
    required this.currentUserType,
    this.displayUserName,
    required this.encryptedAddress,
  }) : super(key: key);

  @override
  _AuthenticationPopUpDialogState createState() =>
      _AuthenticationPopUpDialogState();
}

class _AuthenticationPopUpDialogState extends State<AuthenticationPopUpDialog> {
  final _textEditingController = TextEditingController();
  bool isPasswordMasked = false;
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();
  late Color _mainColor;

  @override
  void initState() {
    _mainColor = widget.currentUserType == CurrentUserType.decentralized
        ? FXColor.mainDeepBlueColor
        : FXColor.mainBlueColor;
    super.initState();
    if (widget.currentUserType == CurrentUserType.decentralized) {
      Future.delayed(Duration(milliseconds: 500),
          () => _tryBiometricAuthentication(widget.encryptedAddress));
    }

    if(AutoFillAccount){
      _textEditingController.text = "aaaaaaa1";
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: GestureDetector(
        onTap: () {
          FocusScopeNode currentFocus = FocusScope.of(context);

          if (!currentFocus.hasPrimaryFocus) {
            currentFocus.unfocus();
          }
        },
        child: SafeArea(
          child: Column(
            children: [
              Align(
                alignment: Alignment.topRight,
                child: Padding(
                  padding: EdgeInsets.only(right: 15.0),
                  child: TextButton.icon(
                    onPressed: () {
                      Navigator.of(context).pop();
                    },
                    icon: Image.asset(
                      'images/icn_switch_account.png',
                      package: 'euruswallet',
                      color: _mainColor,
                      width: 15,
                    ),
                    label: Text(
                      'WELCOME_DIALOG.SWITCH_ACCOUNT'.tr(),
                      style: FXUI.hintStyle.copyWith(color: _mainColor),
                    ),
                  ),
                ),
              ),
              Expanded(
                child: SingleChildScrollView(
                  child: Container(
                    padding: EdgeInsets.symmetric(vertical: 24, horizontal: 36),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.start,
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        Padding(
                          padding: EdgeInsets.symmetric(vertical: 64),
                          child: Image.asset(
                            'images/imgLogo.png',
                            package: 'euruswallet',
                            height: 130,
                          ),
                        ),
                        Text("WELCOME_DIALOG.MAIN_TITLE".tr(),
                            style: FXUI.normalTextStyle.copyWith(
                              fontWeight: FontWeight.w500,
                              fontSize: 16,
                              color: Colors.black,
                            ),
                            textAlign: TextAlign.center),
                        SizedBox(height: 8),
                        if (!isEmptyString(string: widget.displayUserName))
                          Text(
                            widget.displayUserName!,
                            textAlign: TextAlign.center,
                            style: Theme.of(context)
                                .textTheme
                                .caption
                                ?.apply(color: FXColor.lightBlack),
                          ),
                        SizedBox(height: 16),
                        SingleChildScrollView(
                          child: TextField(
                            textInputAction: TextInputAction.go,
                            onSubmitted: (_) async => await _submit(),
                            autofocus: false,
                            obscureText: !isPasswordMasked,
                            controller: _textEditingController,
                            decoration: InputDecoration(
                              focusedBorder: OutlineInputBorder(
                                borderSide: BorderSide(
                                  color: _mainColor,
                                  width: 2,
                                ),
                                borderRadius: FXUI.cricleRadius,
                              ),
                              enabledBorder: OutlineInputBorder(
                                borderSide: BorderSide(
                                  color: _mainColor,
                                  width: 2,
                                ),
                                borderRadius: FXUI.cricleRadius,
                              ),
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
                                  color: _mainColor,
                                ),
                              ),
                              hintText: 'COMMON.PASSWORD_HINT'.tr(),
                              hintStyle: FXUI.normalTextStyle.copyWith(
                                fontWeight: FontWeight.w600,
                                fontSize: 14,
                                color: Colors.grey.shade400,
                              ),
                            ),
                          ),
                        ),
                        SizedBox(height: 16),
                        SubmitButton(
                          btnController: btnController,
                          loadingSecond: 3,
                          label: "WELCOME_DIALOG.LOGIN_BTN".tr(),
                          buttonBGColor: _mainColor,
                          onPressed: () async {
                            await _submit();
                          },
                        ),
                        if (widget.currentUserType ==
                            CurrentUserType.centralized)
                          Padding(
                            padding: EdgeInsets.only(top: 16.0),
                            child: TextButton(
                              child: Text('WELCOME_DIALOG.FORGOT_PW_BTN'.tr(),
                                  style: Theme.of(context)
                                      .textTheme
                                      .caption
                                      ?.apply(color: Colors.blue)),
                              onPressed: () {
                                common.pushPage(
                                    page: CenForgetLoginPwPage(),
                                    context: context);
                              },
                            ),
                          ),
                      ],
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future _tryBiometricAuthentication(String encryptedAddress) async {
    List<UserProfile> userProfiles = await getLocalAcs();
    UserProfile userProfile = userProfiles
        .firstWhere((element) => element.encryptedAddress == encryptedAddress);
    if (getCurrentTimeStamp() <
            (int.tryParse(userProfile.bioAuthValidTime ?? "") ?? -1) &&
        userProfile.enableBiometric == true) {
      final result = await BiometricAuthenticationHelper()
          .canPersistWithBiometricSecurely(
              await common.prefix + encryptedAddress);
      if (result != null) {
        _textEditingController.text = result;
        await _submit();
      }
    }
  }

  Future<void> _submit() async {
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    String textEditingControllerText = _textEditingController.text;

    if (widget.currentUserType == CurrentUserType.centralized) {
      LoginBySignModel loginBySignModel = await api.loginBySignature(
        email: widget.displayUserName ?? '',
        password: textEditingControllerText,
      );

      if (await common.checkApiError(
          context: context,
          errorString: loginBySignModel.message,
          returnCode: loginBySignModel.returnCode)) {
        if (loginBySignModel.returnCode == 0) {
          Navigator.of(context)
              .pop([loginBySignModel, textEditingControllerText]);
        } else {
          showErrorSnackbar();
        }
      }
    } else {
      List<UserProfile> userProfiles = await getLocalAcs();
      UserProfile userProfile = userProfiles.firstWhere(
          (element) => element.encryptedAddress == widget.encryptedAddress);
      final privateKey = CommonMethod.passwordDecrypt(
          textEditingControllerText, userProfile.encryptedPrivateKey);
      final decryptedAddress = CommonMethod.passwordDecrypt(
          textEditingControllerText, widget.encryptedAddress);

      if (decryptedAddress != null && privateKey != null) {
        String? currentUserType =
            await NormalStorageKit().readValue('currentUserType');
        if (currentUserType == 'CurrentUserType.decentralized') {
          common.currentUserType = CurrentUserType.decentralized;
        } else if (currentUserType == 'CurrentUserType.centralized') {
          common.currentUserType = CurrentUserType.centralized;
        }

        // MARK: - use decryptedAddress instead as address extracted from EthPrivateKey is lowered case
        final addressPair = common.getAddressPairFromPrivateKey(privateKey);
        final anAddressPair = AddressPair(decryptedAddress, privateKey,
            publicKey: addressPair.publicKey);
        ScaffoldMessenger.of(context).removeCurrentSnackBar();
        Navigator.of(context).pop([anAddressPair, textEditingControllerText]);
      } else {
        showErrorSnackbar();
      }
    }
    btnController.reset();
  }

  void showErrorSnackbar() {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(
          "COMMON_ERROR.AUTH_FAIL".tr(),
          textAlign: TextAlign.center,
        ),
        backgroundColor: Colors.redAccent.shade100.withOpacity(.9),
      ),
    );
  }
}
