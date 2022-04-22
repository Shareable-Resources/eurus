import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/customDialogBox.dart';
import 'package:euruswallet/pages/decentralized/decentralized_import_keystore_locker_password_setup_page.dart';
import 'package:euruswallet/pages/decentralized/decentralized_import_wallet_page.dart';
import 'package:euruswallet/pages/decentralized/decentralized_wallet_base_page.dart';

import 'centralized/register.dart';

class CreateWalletPage extends StatefulWidget {
  const CreateWalletPage({Key? key}) : super(key: key);

  @override
  _CreateWalletPageState createState() => _CreateWalletPageState();
}

class _CreateWalletPageState extends State<CreateWalletPage> {
  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.max,
      children: <Widget>[
        ListTile(
          dense: true,
          title: Text(
            "CREATE_WALLET_PAGE.SELECT_TYPE".tr(),
            textAlign: TextAlign.center,
            style: Theme.of(context).textTheme.headline5?.apply(
                color: FXColor.blackColor,
                fontFamily: 'packages/euruswallet/SFProDisplay'),
          ),
        ),
        SizedBox(height: 54),
        GestureDetector(
          child: SizedBox(
            width: MediaQuery.of(context).size.width - 70,
            height: (MediaQuery.of(context).size.width - 70) / 2.2,
            child: Container(
              width: double.infinity,
              height: double.infinity,
              padding: EdgeInsets.zero,
              decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: FXUI.cricleRadius,
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.15),
                      spreadRadius: 0,
                      blurRadius: 11,
                      offset: Offset(1, 4), // changes position of shadow
                    ),
                  ]),
              child: Stack(children: [
                Positioned(
                    top: 0,
                    left: 0,
                    child: Container(
                      decoration:
                          BoxDecoration(borderRadius: FXUI.cricleRadius),
                      height: (MediaQuery.of(context).size.width - 56) / 5,
                      child: Image.asset('images/cenBtnBg.png',
                          package: 'euruswallet', fit: BoxFit.fitHeight),
                    )),
                SizedBox.expand(
                    child: Row(
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                      Spacer(flex: 1),
                      Expanded(
                          flex: 7,
                          child:
                              Column(mainAxisSize: MainAxisSize.min, children: [
                            Image.asset('images/icnCentralized.png',
                                package: 'euruswallet'),
                            SizedBox(height: 10),
                            Text("CREATE_WALLET_PAGE.CEN.TITLE".tr(),
                                style: FXUI.normalTextStyle.copyWith(
                                    fontWeight: FontWeight.bold,
                                    color: FXColor.lightBlack)),
                          ])),
                      Spacer(flex: 1),
                      Expanded(
                          flex: 10,
                          child: Text('CREATE_WALLET_PAGE.CEN_BTN.DESC'.tr(),
                              style: FXUI.normalTextStyle.copyWith(
                                  color: FXColor.lightBlackColor
                                      .withOpacity(0.6)))),
                    ])),
                Positioned(
                    right: 0,
                    top: 0,
                    child: GestureDetector(
                        child: Padding(
                            padding: EdgeInsets.all(5),
                            child:
                                Icon(Icons.info, color: FXColor.mainBlueColor)),
                        onTap: () async {
                          await showDialog(
                              context: context,
                              builder: (BuildContext context) {
                                return CustomDialogBox(
                                    title:
                                        "CREATE_WALLET_PAGE.CEN.INFO_DIALOG.TITLE"
                                            .tr(),
                                    descriptions:
                                        "CREATE_WALLET_PAGE.CEN.INFO_DIALOG.DESC"
                                            .tr(),
                                    buttonText: "COMMON.CLOSE".tr());
                              });
                        })),
              ]),
            ),
          ),
          onTap: () {
            // TODO - create a new decentralized wallet with new random seed
            common.pushPage(page: RegisterPage(), context: context);
          },
        ),
        SizedBox(height: 50),
        GestureDetector(
          child: SizedBox(
            width: MediaQuery.of(context).size.width - 70,
            height: (MediaQuery.of(context).size.width - 70) / 2.2,
            child: Container(
              width: double.infinity,
              height: double.infinity,
              padding: EdgeInsets.zero,
              decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: FXUI.cricleRadius,
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withOpacity(0.15),
                      spreadRadius: 0,
                      blurRadius: 11,
                      offset: Offset(1, 4), // changes position of shadow
                    ),
                  ]),
              child: Stack(children: [
                Positioned(
                    top: 0,
                    left: 0,
                    child: Container(
                      decoration:
                          BoxDecoration(borderRadius: FXUI.cricleRadius),
                      height: (MediaQuery.of(context).size.width - 56) / 5,
                      child: Image.asset('images/decenBtnBg.png',
                          package: 'euruswallet', fit: BoxFit.fitHeight),
                    )),
                SizedBox.expand(
                    child: Row(
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                      Spacer(flex: 1),
                      Expanded(
                          flex: 7,
                          child:
                              Column(mainAxisSize: MainAxisSize.min, children: [
                            Image.asset('images/icnDecentralized.png',
                                package: 'euruswallet'),
                            SizedBox(height: 10),
                            Text("CREATE_WALLET_PAGE.DECEN.TITLE".tr(),
                                style: FXUI.normalTextStyle.copyWith(
                                    fontWeight: FontWeight.bold,
                                    color: FXColor.lightBlack)),
                          ])),
                      Spacer(flex: 1),
                      Expanded(
                          flex: 10,
                          child: Text('CREATE_WALLET_PAGE.DECEN_BTN.DESC'.tr(),
                              style: FXUI.normalTextStyle.copyWith(
                                  color: FXColor.lightBlackColor
                                      .withOpacity(0.6)))),
                    ])),
                Positioned(
                    right: 0,
                    top: 0,
                    child: GestureDetector(
                        child: Padding(
                            padding: EdgeInsets.all(5),
                            child:
                                Icon(Icons.info, color: FXColor.mainBlueColor)),
                        onTap: () async {
                          await showDialog(
                              context: context,
                              builder: (BuildContext context) {
                                return CustomDialogBox(
                                    title:
                                        "CREATE_WALLET_PAGE.DECEN.INFO_DIALOG.TITLE"
                                            .tr(),
                                    descriptions:
                                        "CREATE_WALLET_PAGE.DECEN.INFO_DIALOG.DESC"
                                            .tr(),
                                    buttonText: "COMMON.CLOSE".tr());
                              });
                        })),
              ]),
            ),
          ),
          onTap: () {
            Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) => DecentralizedWalletBasePage(
                  appBarTitle: Text("CREATE_LOCKER_PAGE.MAIN_TITLE".tr()),
                  body: DecentralizedImportKeyStoreLockerPasswordSetupPage(
                    importWalletType: ImportWalletType.mnemonic,
                  ),
                ),
              ),
            );
          },
        ),
      ],
    );
  }
}
