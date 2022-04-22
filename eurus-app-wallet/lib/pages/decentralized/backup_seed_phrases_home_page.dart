import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/walletLockerPWDialog.dart';
import 'package:euruswallet/pages/decentralized/backup_seed_phrases_confirmation_page.dart';
import 'package:flutter/cupertino.dart';
import 'package:screenshot_callback/screenshot_callback.dart';

class BackupSeedPhrasesHomePage extends StatefulWidget {
  final Function? backupSeedPhrasesCompletion;

  const BackupSeedPhrasesHomePage({
    Key? key,
    this.backupSeedPhrasesCompletion,
  }) : super(key: key);

  @override
  _BackupSeedPhrasesHomePageState createState() =>
      _BackupSeedPhrasesHomePageState();
}

class _BackupSeedPhrasesHomePageState extends State<BackupSeedPhrasesHomePage> {
  ScreenshotCallback screenshotCallback = ScreenshotCallback();
  OverlayEntry? overlayEntry;
  late OverlayEntry _promptEnvSecurityBeforeRevealOverlayEntry;
  late final future;

  @override
  void initState() {
    super.initState();
    future = _authUI();

    _promptEnvSecurityBeforeRevealOverlayEntry = OverlayEntry(
        builder: (_overlayEntryBuilderContext) => Scaffold(
            backgroundColor: Colors.black45,
            body: Padding(
                padding: getEdgeInsetsSymmetric(horizontal: 8),
                child: Center(
                    child: Container(
                        decoration: BoxDecoration(
                          borderRadius: BorderRadius.circular(30),
                          color: Colors.white,
                        ),
                        padding:
                            EdgeInsets.symmetric(horizontal: 8, vertical: 26),
                        child: Column(
                            mainAxisSize: MainAxisSize.min,
                            crossAxisAlignment: CrossAxisAlignment.stretch,
                            children: [
                              Stack(
                                  clipBehavior: Clip.none,
                                  alignment: Alignment.center,
                                  children: [
                                    ListTile(
                                        title: Text("COMMON.CAUTION".tr(),
                                            textAlign: TextAlign.center,
                                            style: Theme.of(context)
                                                .textTheme
                                                .headline6
                                                ?.copyWith(
                                                    letterSpacing: 2,
                                                    fontWeight:
                                                        FontWeight.w800)))
                                  ]),
                              Padding(
                                  padding: EdgeInsets.symmetric(horizontal: 19),
                                  child: ListTile(
                                      title: Icon(Icons.warning_rounded,
                                          size: 120, color: Colors.red),
                                      subtitle: Text(
                                        "\n\n${'BACKUP_SEEDPHRASE_PAGE.CAUTION_DIALOG.CONTENT'.tr()}",
                                        textAlign: TextAlign.center,
                                      ))),
                              SizedBox(height: 30),
                              Padding(
                                  padding: EdgeInsets.symmetric(horizontal: 19),
                                  child: CupertinoButton(
                                    child: Text("COMMON.BACKUP_NOW".tr()),
                                    color: FXColor.deepBlue,
                                    borderRadius: FXUI.cricleRadius,
                                    padding: EdgeInsets.symmetric(
                                        horizontal: 0, vertical: 12),
                                    onPressed: () {
                                      __dismissEnvSecurityBeforeRevealOverlayEntry();
                                    },
                                  )),
                              SizedBox(height: 12),
                              Padding(
                                  padding: EdgeInsets.symmetric(horizontal: 19),
                                  child: IntrinsicHeight(
                                      child: CupertinoButton(
                                          child: SizedBox.expand(
                                              child: Container(
                                                  decoration: BoxDecoration(
                                                      border: Border.all(
                                                          color: FXColor
                                                              .mainDeepBlueColor),
                                                      borderRadius:
                                                          FXUI.cricleRadius),
                                                  child: Center(
                                                      child: Text(
                                                          "BACKUP_SEEDPHRASE_PAGE.CAUTION_DIALOG.BACKUP_LATER"
                                                              .tr(),
                                                          style: FXUI
                                                              .normalTextStyle
                                                              .copyWith(
                                                                  color: FXColor
                                                                      .mainDeepBlueColor))))),
                                          borderRadius: FXUI.cricleRadius,
                                          padding: EdgeInsets.zero,
                                          onPressed: () {
                                            __dismissEnvSecurityBeforeRevealOverlayEntry();
                                            Navigator.pop(context);
                                          }))),
                              SizedBox(height: 14),
                              // Padding(padding: EdgeInsets.symmetric(horizontal: 19) , child: CupertinoButton(child: Text("Backup Later"), color: FXColor.mainDeepBlueColor, borderRadius: BorderRadius.circular(15), padding: EdgeInsets.symmetric(horizontal: 0, vertical: 12), onPressed: null)),
                            ]))))));

    screenshotCallback.addListener(() {
      if (overlayEntry?.mounted ?? false) overlayEntry?.remove();
      overlayEntry = OverlayEntry(
        builder: (_overlayEntryBuilderContext) => Scaffold(
          backgroundColor: Colors.black45,
          body: Padding(
            padding: getEdgeInsetsSymmetric(horizontal: 26),
            child: Center(
              child: Container(
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(20),
                  color: Colors.white,
                ),
                padding: EdgeInsets.fromLTRB(0, 27, 0, 48),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Text(
                      "BACKUP_SEEDPHRASE_PAGE.SCREENSHOT_DIALOG.TITLE".tr(),
                      textAlign: TextAlign.center,
                      style: FXUI.titleTextStyle.copyWith(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    SizedBox(height: 52),
                    Icon(
                      Icons.warning_rounded,
                      size: 120,
                      color: Colors.red,
                    ),
                    SizedBox(height: 25),
                    Padding(
                      padding: EdgeInsets.symmetric(horizontal: 24),
                      child: Text(
                        "BACKUP_SEEDPHRASE_PAGE.SCREENSHOT_DIALOG.CONTENT".tr(),
                        textAlign: TextAlign.center,
                        style: FXUI.normalTextStyle.copyWith(
                          color: FXColor.textGray,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ),
                    SizedBox(height: 30),
                    Padding(
                      padding: EdgeInsets.symmetric(horizontal: 19),
                      child: CupertinoButton(
                        child: Text(
                          "BACKUP_SEEDPHRASE_PAGE.SCREENSHOT_DIALOG.BUTTON_TEXT"
                              .tr(),
                        ),
                        color: FXColor.deepBlue,
                        borderRadius: FXUI.cricleRadius,
                        padding:
                            EdgeInsets.symmetric(horizontal: 0, vertical: 12),
                        onPressed: () {
                          overlayEntry?.remove();
                        },
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      );
      Overlay.of(context)!.insert(overlayEntry!);
    });
  }

  @override
  void dispose() {
    screenshotCallback.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Center(
          child: ListTile(
            dense: true,
            title: Text(
              'BACKUP_SEEDPHRASE_PAGE.DISPLAY_SEEDPHRASE.TITLE'.tr(),
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.headline5?.copyWith(
                  color: FXColor.blackColor,
                  fontSize: 24,
                  fontFamily: 'packages/euruswallet/SFProDisplay',
                  fontWeight: FontWeight.w600),
            ),
          ),
        ),
        FutureBuilder(
            future: future,
            builder: (context, s) => s.hasData && s.data is String
                ? _mSeedPhraseRevealerWidget(s.data as String)
                : CircularProgressIndicator())
      ],
    );
  }

  Widget _mSeedPhraseRevealerWidget(String _mSeedPhrase) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        ListTile(
            title: Text(
                "BACKUP_SEEDPHRASE_PAGE.DISPLAY_SEEDPHRASE.CONTENT".tr(),
                style: Theme.of(context).textTheme.bodyText2?.apply(
                    color: FXColor.lightBlackColor.withOpacity(.6),
                    fontWeightDelta: -1))),
        SizedBox(height: 35),
        ...List.generate(
            4,
            (y) => Padding(
                padding: EdgeInsets.only(bottom: 12),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: List.generate(
                      3,
                      (x) => Container(
                            margin: x == 1
                                ? EdgeInsets.symmetric(horizontal: 8)
                                : null,
                            padding: EdgeInsets.symmetric(
                                vertical: 12, horizontal: 14),
                            decoration: BoxDecoration(
                                border: Border.all(color: FXColor.deepBlue),
                                borderRadius: FXUI.cricleRadius),
                            child: Text(
                                "${y * 3 + x + 1}. ${_mSeedPhrase.split(' ')[y * 3 + x]}",
                                style: Theme.of(context)
                                    .textTheme
                                    .caption
                                    ?.apply(color: FXColor.deepBlue)),
                          )),
                ))),
        SizedBox(
          height: 149,
        ),
        CupertinoButton(
          color: FXColor.deepBlue,
          borderRadius: FXUI.cricleRadius,
          child: Text('COMMON.NEXT_STEP'.tr()),
          onPressed: () {
            Navigator.of(context).push(
              MaterialPageRoute(
                  builder: (context) => BackupSeedPhrasesConfirmationPage(
                        mSeedPhrase: _mSeedPhrase,
                        parentNotifyOnValueChange:
                            widget.backupSeedPhrasesCompletion,
                      )),
            );
          },
        ),
      ],
    );
  }

  void __dismissEnvSecurityBeforeRevealOverlayEntry() {
    _promptEnvSecurityBeforeRevealOverlayEntry.remove();
  }

  Future<String?> _authUI() async {
    return await Future.delayed(const Duration(milliseconds: 0), () async {
      return await Navigator.of(context).push(PageRouteBuilder(
        fullscreenDialog: true,
        opaque: false,
        pageBuilder: (pageBuilderContext, animation, secondaryAnimation) =>
            WalletLockerPWDialog(
                themeColor: common.getBackGroundColor(),
                decenUserkey: common.currentUserProfile!.mnemonicSeedPhrases,
                submitFnc: _submit),
      ));
    }).then((value) {
      if (isEmptyString(string: value as String?)) {
        Navigator.pop(context);
        return null;
      } else {
        Overlay.of(context)?.insert(_promptEnvSecurityBeforeRevealOverlayEntry);
        return value;
      }
    });
  }

  Future<String?> _submit(
      String? _uriQuery, TextEditingController _textEditingController) async {
    print(
        '_authUI _submit _accountEncryptedAddress ${common.encryptedAddress} $_uriQuery');
    final _encryptedValue = common.getUriVal(
        _uriQuery, common.currentUserProfile?.encryptedAddress ?? '');
    final _decryptedValue = CommonMethod.passwordDecrypt(
        _textEditingController.text, _encryptedValue ?? '');
    return _decryptedValue;
  }
}
