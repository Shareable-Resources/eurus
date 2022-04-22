import 'package:app_authentication_kit/mnemonic_kit.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/pages/decentralized/decentralized_import_keystore_locker_password_setup_page.dart';
import 'package:flutter/cupertino.dart';
import 'package:web3dart/crypto.dart' as _web3dartCrypto;

import '../../common/commonMethod.dart';
import '../decentralized/decentralized_wallet_base_page.dart';

enum ImportWalletType { mnemonic, privateKey }

extension ImportWalletTypeExtension on ImportWalletType {
  String get name {
    switch (this) {
      case ImportWalletType.mnemonic:
        return 'IMPORT_WALLET_PAGE.MPHRASE.NAME'.tr();
      case ImportWalletType.privateKey:
        return 'IMPORT_WALLET_PAGE.PRIVATE_KEY.NAME'.tr();
    }
  }

  String get title {
    switch (this) {
      case ImportWalletType.mnemonic:
        return 'IMPORT_WALLET_PAGE.MPHRASE.TITLE'.tr();
      case ImportWalletType.privateKey:
        return 'IMPORT_WALLET_PAGE.PRIVATE_KEY.TITLE'.tr();
    }
  }

  String get inputHint {
    switch (this) {
      case ImportWalletType.mnemonic:
        return 'IMPORT_WALLET_PAGE.MPHRASE.INPUT_HINT'.tr();
      case ImportWalletType.privateKey:
        return 'IMPORT_WALLET_PAGE.PRIVATE_KEY.INPUT_HINT'.tr();
    }
  }

  String get invalidInputError {
    switch (this) {
      case ImportWalletType.mnemonic:
        return 'IMPORT_WALLET_PAGE.MPHRASE.INVALID_INPUT'.tr();
      case ImportWalletType.privateKey:
        return 'IMPORT_WALLET_PAGE.PRIVATE_KEY.INVALID_INPUT'.tr();
    }
  }
}

class DecentralizedImportWalletPage extends StatefulWidget {
  const DecentralizedImportWalletPage({Key? key}) : super(key: key);

  @override
  _DecentralizedWalletPageState createState() =>
      _DecentralizedWalletPageState();
}

class _DecentralizedWalletPageState
    extends State<DecentralizedImportWalletPage> {
  final _textEditingController = TextEditingController();
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  bool _isFormValid = false;
  ImportWalletType _importWalletType = ImportWalletType.mnemonic;
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    super.initState();
    common.currentUserType = CurrentUserType.decentralized;
    autoFillText(importWalletType: _importWalletType);
  }

  @override
  void dispose() {
    common.currentUserType = CurrentUserType.centralized;
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      top: false,
      child: DecentralizedWalletBasePage(
        appBarTitle: Text('LAUNCH_PAGE.IMPORT_WALLET_BTN'.tr()),
        body: _getImportWalletPageBody(),
      ),
    );
  }

  void autoFillText({required ImportWalletType importWalletType}) {
    if (AutoFillAccount) {
      Future.delayed(const Duration(milliseconds: 500), () {
        if (importWalletType == ImportWalletType.privateKey) {
          _textEditingController.text = envType == EnvType.Production
              ? "4c0575dc329e621460f73d4c4fa731ade88c61804fa9a89672832eeb48e2620f"
              : 'd1bdc683fbeb9fa0b4ceb26adb39eaffb21b16891ea28e4cf1bc3118fdd39295';
        } else {
          _textEditingController.text = envType == EnvType.Production
              ? 'dove clock chalk front spike prefer people spike capable word gasp congress'
              : 'carbon shuffle shoot knock alter bottom polar maple husband poet match spring';
        }
      });
    }
  }

  Widget _getImportWalletPageBody() {
    return Form(
      key: _formKey,
      onChanged: () {
        if (_formKey.currentState != null)
          setState(() {
            _isFormValid = _formKey.currentState!.validate();
          });
      },
      child: Column(
        children: [
          getSelectionSegmentButton(),
          Center(
            child: ListTile(
              dense: true,
              title: Text(
                _importWalletType.title,
                textAlign: TextAlign.center,
                style: Theme.of(context).textTheme.headline5?.copyWith(
                    color: FXColor.blackColor,
                    fontSize: 24,
                    fontFamily: 'packages/euruswallet/SFProDisplay',
                    fontWeight: FontWeight.w600),
              ),
            ),
          ),
          Expanded(
            flex: 1,
            child: Padding(
              padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              child: TextFormField(
                autovalidateMode: AutovalidateMode.always,
                validator: _importWalletType == ImportWalletType.mnemonic
                    ? _mnemonicValidator
                    : _privateKeyValidator,
                decoration: FXUI.defaultTextFieldInputDecoration.copyWith(
                  hintText: _importWalletType.inputHint,
                  hintMaxLines: 4,
                  hintStyle: Theme.of(context)
                      .textTheme
                      .subtitle1
                      ?.apply(color: Theme.of(context).hintColor),
                ),
                maxLines: null,
                expands: true,
                textAlignVertical: TextAlignVertical.center,
                style: Theme.of(context).textTheme.headline5,
                autofocus: true,
                controller: _textEditingController,
              ),
            ),
          ),
          SubmitButton(
            btnController: btnController,
            buttonBGColor: FXColor.mainDeepBlueColor,
            label: 'IMPORT_WALLET_PAGE.IMPORT_BTN'.tr(),
            onPressed: () async {
              setState(() {
                _isFormValid = _formKey.currentState!.validate();
              });
              if (_isFormValid)
                Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (context) => DecentralizedWalletBasePage(
                      appBarTitle: Text("CREATE_LOCKER_PAGE.MAIN_TITLE".tr()),
                      body: DecentralizedImportKeyStoreLockerPasswordSetupPage(
                        importWalletType: _importWalletType,
                        privateKey:
                            _importWalletType == ImportWalletType.privateKey
                                ? _textEditingController.text
                                : null,
                        mnemonicSeedPhrase:
                            _importWalletType == ImportWalletType.mnemonic
                                ? _textEditingController.text
                                : null,
                        isSeedPhraseImported: true,
                      ),
                    ),
                  ),
                );
              btnController.reset();
            },
          ),
          SizedBox(height: 8),
        ],
      ),
    );
  }

  Widget getSelectionSegmentButton() {
    final buttons = ImportWalletType.values.map((e) {
      final _isSelected = e == _importWalletType;
      return CupertinoButton(
        padding: EdgeInsets.all(4),
        child: Container(
          padding: EdgeInsets.all(4),
          decoration: ShapeDecoration(
            shape: StadiumBorder(),
            color: _isSelected ? Colors.white : Colors.transparent,
          ),
          child: Text(
            e.name,
            style: FXUI.normalTextStyle.copyWith(
                color: _isSelected
                    ? common.getBackGroundColor()
                    : FXColor.lightPurpleColor,
                fontSize: 14,
                fontWeight: FontWeight.w600),
          ),
        ),
        onPressed: () {
          setState(() {
            _importWalletType = e;
            autoFillText(importWalletType: _importWalletType);
          });
        },
      );
    }).toList();

    return Container(
      decoration: ShapeDecoration(
        color: common.getBackGroundColor(),
        shape: StadiumBorder(),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        mainAxisSize: MainAxisSize.min,
        children: buttons,
      ),
    );
  }

  // form texteditingcontroller validator
  String? _privateKeyValidator(String? privateKey) {
    final aPrivateKey = privateKey ?? '';
    if (!isEmptyString(string: aPrivateKey) &&
        _web3dartCrypto.strip0x(aPrivateKey).length == 64) return null;

    return "IMPORT_WALLET_PAGE.PRIVATE_KEY.INVALID_INPUT".tr();
  }

  // form texteditingcontroller validator
  String? _mnemonicValidator(String? mPhrase) {
    // TODO - this mnemonic phrase validator currently keep complaining all inputs even typing is not yet completed nor ready to be dictated, i.e. return red warning messages since typing 1st ~ 11th letter, words, phrases
    if (isEmptyString(string: mPhrase)) return null;

    return _validateMnemonic(mPhrase ?? '')
        ? null
        : "IMPORT_WALLET_PAGE.MPHRASE.INVALID_INPUT".tr();
  }

  /// Check if mnemonic phrase is valid
  bool _validateMnemonic(String mPhrase) {
    return MnemonicKit().validateMnemonic(mPhrase);
  }
}
