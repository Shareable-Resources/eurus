import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:flutter/cupertino.dart';
import 'decentralized_wallet_base_page.dart';
import 'package:euruswallet/common/user_profile.dart';

class BackupSeedPhrasesConfirmationPage extends StatefulWidget {
  final String mSeedPhrase;
  final Function? parentNotifyOnValueChange;

  const BackupSeedPhrasesConfirmationPage({
    Key? key,
    required this.mSeedPhrase,
    this.parentNotifyOnValueChange,
  }) : super(key: key);

  @override
  _BackupSeedPhrasesConfirmationPageState createState() =>
      _BackupSeedPhrasesConfirmationPageState();
}

class _BackupSeedPhrasesConfirmationPageState
    extends State<BackupSeedPhrasesConfirmationPage> {
  late List<String> _randomlyOrderedPhrases;
  List<int> inputPhrasesIndices = List.filled(12, -1);
  int explicitWordCursor = -1;
  FocusNode focusNode = FocusNode();

  @override
  void initState() {
    super.initState();

    _randomlyOrderedPhrases = widget.mSeedPhrase.split(' ')..shuffle();
  }

  Future<void> _updateProfilePhraseBackup() async {
    List<UserProfile> userProfiles = await getLocalAcs();
    UserProfile userProfile = userProfiles.firstWhere((element) =>
        element.encryptedAddress ==
        common.currentUserProfile!.encryptedAddress);
    await setAcToLocal(userProfile, delete: true);
    userProfile.seedPraseBackuped = true;
    common.currentUserProfile = userProfile;
    await setAcToLocal(userProfile);
  }

  @override
  Widget build(BuildContext context) {
    // Future.delayed(Duration.zero, () {focusNode.requestFocus();});
    return DecentralizedWalletBasePage(
      appBarTitle: Text("BACKUP_SEEDPHRASE_PAGE.MAIN_TITLE".tr()),
      body: StatefulBuilder(
        builder: (_statefulBuilderCtx, __setState) {
          Future.delayed(Duration.zero, () => focusNode.requestFocus());
          return Padding(
              padding: EdgeInsets.only(
                  bottom: MediaQuery.of(_statefulBuilderCtx).padding.bottom),
              child: Column(children: [
                Center(
                    child: ListTile(
                        dense: true,
                        title: Text(
                          'BACKUP_SEEDPHRASE_PAGE.CONFIRM_SEEDPHRASE.TITLE'
                              .tr(),
                          textAlign: TextAlign.center,
                          style: Theme.of(_statefulBuilderCtx)
                              .textTheme
                              .headline5
                              ?.copyWith(
                                  color: FXColor.blackColor,
                                  fontSize: 24,
                                  fontFamily:
                                      'packages/euruswallet/SFProDisplay',
                                  fontWeight: FontWeight.w600),
                        ))),
                ListTile(
                    dense: true,
                    title: Text(
                        'BACKUP_SEEDPHRASE_PAGE.CONFIRM_SEEDPHRASE.CONTENT'
                            .tr(),
                        textAlign: TextAlign.left,
                        style: Theme.of(_statefulBuilderCtx)
                            .textTheme
                            .bodyText2
                            ?.apply(
                                color: FXColor.lightBlackColor.withOpacity(.6),
                                fontWeightDelta: -1))),
                SizedBox(
                  height: 35,
                ),
                // ════════ Exception caught by rendering library ═════════════════════════════════
                // An InputDecorator, which is typically created by a TextField, cannot have an unbounded width.
                // This happens when the parent widget does not provide a finite width constraint. For example, if the InputDecorator is contained by a Row, then its width must be constrained. An Expanded widget or a SizedBox can be used to constrain the width of the InputDecorator or the TextField that contains it.
                // 'package:flutter/src/material/input_decorator.dart':
                // Failed assertion: line 948 pos 7: 'layoutConstraints.maxWidth < double.infinity'
                // ...List.generate(4, (y) => Row(mainAxisAlignment: MainAxisAlignment.spaceAround, children: List.generate(3, (x) => Expanded(child: Padding(padding: EdgeInsets.all(2), child: TextField(onTap: () {inputPhrases[y*3+x] = ''; __setState(() => explicitWordCursor = inputPhrases.indexWhere((e) => e == null || e.isEmpty) == y*3+x ? -1 : y*3+x);}, focusNode: explicitWordCursor > -1 && explicitWordCursor == y*3+x || explicitWordCursor == -1 && inputPhrases.indexWhere((e) => e == null || e.isEmpty) == y*3+x ? focusNode : null, readOnly: true, enableInteractiveSelection: false, decoration: InputDecoration(border: OutlineInputBorder(borderRadius: BorderRadius.circular(15),), focusedBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(15), borderSide: BorderSide(color: Color(0xFF4406dd)))/*, fillColor: Color(0xFF4406dd), filled: false,*/, labelText: "${y*3+x+1}"), controller: TextEditingController(text: "${inputPhrases[y*3+x]??''}"),)))))),
                // ...List.generate(4, (y) => Row(mainAxisAlignment: MainAxisAlignment.spaceAround, children: List.generate(3, (x) => OutlineButton(onPressed: () {inputPhrases[explicitWordCursor > -1 ? explicitWordCursor : inputPhrases.indexWhere((e) => e == null || e.isEmpty)] = randomlyOrderedPhrases[y*3+x]; __setState(() => explicitWordCursor = -1);}, color: Color(0xFF4406dd), child: Text("${randomlyOrderedPhrases[y*3+x]}"),)))),
                // if (_mSeedPhrase == (_mSeedPhrase.split(' ')??inputPhrases).join(' ')) CupertinoButton(onPressed: () => Navigator.popUntil(_statefulBuilderCtx, ModalRoute.withName('HomePage')), child: Text('Finished')),
                ...List.generate(
                    4,
                    (y) => Row(
                        mainAxisAlignment: MainAxisAlignment.spaceAround,
                        children: List.generate(
                            3,
                            (x) => Expanded(
                                child: Padding(
                                    padding: EdgeInsets.only(bottom: 12) +
                                        (x == 1
                                            ? EdgeInsets.symmetric(
                                                horizontal: 8)
                                            : EdgeInsets.zero),
                                    child: TextField(
                                      onTap: () {
                                        inputPhrasesIndices[y * 3 + x] = -1;
                                        __setState(() => explicitWordCursor =
                                            inputPhrasesIndices.indexWhere(
                                                        (e) => e == -1) ==
                                                    y * 3 + x
                                                ? -1
                                                : y * 3 + x);
                                      },
                                      focusNode: explicitWordCursor > -1 &&
                                                  explicitWordCursor ==
                                                      y * 3 + x ||
                                              explicitWordCursor == -1 &&
                                                  inputPhrasesIndices
                                                          .indexWhere(
                                                              (e) => e == -1) ==
                                                      y * 3 + x
                                          ? focusNode
                                          : null,
                                      readOnly: true,
                                      enableInteractiveSelection: false,
                                      style: Theme.of(_statefulBuilderCtx)
                                          .textTheme
                                          .caption,
                                      textAlign: TextAlign.center,
                                      decoration: InputDecoration(
                                          isCollapsed: true,
                                          contentPadding: EdgeInsets.symmetric(
                                              vertical: 12, horizontal: 14),
                                          border: OutlineInputBorder(
                                            borderRadius: FXUI.cricleRadius,
                                          ),
                                          focusedBorder: OutlineInputBorder(
                                              borderRadius: FXUI.cricleRadius,
                                              borderSide: BorderSide(
                                                  color: FXColor.deepBlue,
                                                  width:
                                                      1.5)) /*, fillColor: Color(0xFF4406dd), filled: false,*/,
                                          labelText: "${y * 3 + x + 1}."),
                                      controller: TextEditingController(
                                          text:
                                              "${(inputPhrasesIndices[y * 3 + x]) > -1 ? _randomlyOrderedPhrases[inputPhrasesIndices[y * 3 + x]] : ''}"),
                                    )))))),
                Expanded(child: Container()),
                ...List.generate(
                    4,
                    (y) => Row(
                        mainAxisAlignment: MainAxisAlignment.spaceAround,
                        children: List.generate(
                            3,
                            (x) => Expanded(
                                child: Padding(
                                    padding: EdgeInsets.only(bottom: 12) +
                                        (x == 1
                                            ? EdgeInsets.symmetric(
                                                horizontal: 8)
                                            : EdgeInsets.zero),
                                    child: CupertinoButton(
                                      onPressed: inputPhrasesIndices
                                                  .indexOf(y * 3 + x) >
                                              -1
                                          ? null
                                          : () {
                                              inputPhrasesIndices[
                                                      explicitWordCursor > -1
                                                          ? explicitWordCursor
                                                          : inputPhrasesIndices
                                                              .indexWhere((e) =>
                                                                  e == -1)] =
                                                  y * 3 + x;
                                              __setState(() =>
                                                  explicitWordCursor = -1);
                                            },
                                      color: FXColor.deepBlue,
                                      padding:
                                          EdgeInsets.symmetric(vertical: 8),
                                      minSize: null,
                                      borderRadius: FXUI.cricleRadius,
                                      child: Text(
                                          "${_randomlyOrderedPhrases[y * 3 + x]}",
                                          style: Theme.of(_statefulBuilderCtx)
                                              .textTheme
                                              .caption
                                              ?.apply(color: Colors.white)),
                                    )))))),
                CupertinoButton(
                    onPressed:
                        (inputPhrasesIndices.indexWhere((e) => e == -1) == -1 &&
                                widget.mSeedPhrase ==
                                    (inputPhrasesIndices
                                            .map((e) => e == -1
                                                ? ''
                                                : _randomlyOrderedPhrases[e])
                                            .toList())
                                        .join(' '))
                            ? () async {
                                if (widget.parentNotifyOnValueChange != null)
                                  widget.parentNotifyOnValueChange!();
                                await _updateProfilePhraseBackup();
                                Navigator.popUntil(_statefulBuilderCtx,
                                    ModalRoute.withName('HomePage'));
                              }
                            : null,
                    child: Text('COMMON.FINISH'.tr())),
              ]));
        },
      ),
    );
  }
}
