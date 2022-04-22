import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';

import '../common/commonMethod.dart';
import '../pages/decentralized/backup_seed_phrases_home_page.dart';
import '../pages/decentralized/decentralized_wallet_base_page.dart';

class ReplacingQRCodeWidget extends StatelessWidget {
  final Function? backupSeedPhrasesCompletion;

  const ReplacingQRCodeWidget({
    Key? key,
    this.backupSeedPhrasesCompletion,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        backgroundColor: Colors.black45,
        body: Padding(
            padding: getEdgeInsetsSymmetric(horizontal: 8),
            child: Center(
                child: Container(
              decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(30), color: Colors.white),
              padding: EdgeInsets.symmetric(horizontal: 8, vertical: 26),
              child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    // Stack(alignment: Alignment.center, overflow: Overflow.visible, children: [ListTile(title: Text("WARNING", textAlign: TextAlign.center, style: Theme.of(context).textTheme.headline6.copyWith(letterSpacing: 2, fontWeight: FontWeight.w800))), Positioned(right: -24, child: IconButton(icon: Icon(Icons.close), onPressed: () {print(1); Navigator.of(_backupWalletSecretMnemonicSeedPhrasesWidgetBuilderContext).popUntil(ModalRoute.withName('HomePage'));},))]),
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
                                          fontWeight: FontWeight.w800))),
                          Positioned(
                              right: -24.0 + 19,
                              child: IconButton(
                                icon: Icon(Icons.close),
                                onPressed: () {
                                  Navigator.of(context).pop();
                                },
                              ))
                        ]),
                    // Icon(Icons.info_outline_rounded, size: 100,),
                    Padding(
                        padding: EdgeInsets.symmetric(horizontal: 19),
                        child: ListTile(
                            title: Icon(Icons.warning_rounded,
                                size: 120, color: Colors.red),
                            subtitle: Text(
                              '\n\n${"BACKUP_MPHRASE.RECEIVE_ALERT_DIALOG.CONTENT".tr()}',
                              textAlign: TextAlign.center,
                            ))),
                    SizedBox(height: 34),
                    Padding(
                        padding: EdgeInsets.symmetric(horizontal: 19),
                        child: CupertinoButton(
                          child: Text("COMMON.BACKUP_NOW".tr()),
                          color: FXColor.deepBlue,
                          borderRadius: FXUI.cricleRadius,
                          padding:
                              EdgeInsets.symmetric(horizontal: 0, vertical: 12),
                          onPressed: () async {
                            final _popResult = await Navigator.of(context).push(
                                MaterialPageRoute(
                                    fullscreenDialog: true,
                                    builder: (_materialPageRouteBuilderCtx) =>
                                        DecentralizedWalletBasePage(
                                            appBarTitle: Text(
                                                "BACKUP_SEEDPHRASE_PAGE.MAIN_TITLE"
                                                    .tr()),
                                            body: BackupSeedPhrasesHomePage(
                                              backupSeedPhrasesCompletion:
                                                  backupSeedPhrasesCompletion,
                                            ))));
                            if (_popResult is bool && _popResult)
                              Navigator.of(context).popUntil(
                                  (route) => route.settings.name == 'HomePage');
                          },
                        )),
                    SizedBox(height: 14),
                  ]),
            ))));
  }
}
