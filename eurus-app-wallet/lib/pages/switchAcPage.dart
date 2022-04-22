import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/common/user_profile.dart';
import 'package:euruswallet/pages/centralized/cenLoginPage.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter_slidable/flutter_slidable.dart';

import 'centralized/cenForgetLoginPwPage.dart';

class SwitchAcPage extends StatefulWidget {
  SwitchAcPage({Key? key}) : super(key: key);

  @override
  _SwitchAcPageState createState() => _SwitchAcPageState();
}

class _SwitchAcPageState extends State<SwitchAcPage> {
  List<UserProfile>? acs;

  bool get acReadys => acs != null;

  @override
  void initState() {
    super.initState();
    _getAccounts();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          image: DecorationImage(
            image: AssetImage(
              'images/backgroundImage.png',
              package: 'euruswallet',
            ),
            fit: BoxFit.cover,
            alignment: Alignment.topCenter,
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            AppBar(
              title: Text('SWITCH_AC.TITLE'.tr()),
              backgroundColor: Colors.transparent,
              elevation: 0,
            ),
            Expanded(
              flex: 1,
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.only(
                    topLeft: Radius.circular(15),
                    topRight: Radius.circular(15),
                  ),
                ),
                margin: EdgeInsets.only(top: 12),
                padding: EdgeInsets.only(bottom: 31, top: 15),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Expanded(
                      child: ListView(
                        padding: EdgeInsets.zero,
                        children: genAcListItems(),
                      ),
                    ),
                    GestureDetector(
                      onTap: () {
                        print('Forgot pw btn onpress');
                        common.pushPage(
                            page: CenForgetLoginPwPage(), context: context);
                      },
                      child: Center(
                        child: Text(
                          'CEN_LOGIN.FORGOT_PW'
                              .tr() /*+
                              "/ ${'REPORT_AC.REPORT_ACCOUNT'.tr()}"*/
                          ,
                          style: FXUI.normalTextStyle.copyWith(
                              color: FXColor.mainBlueColor, fontSize: 14),
                        ),
                      ),
                    ),
                    Container(
                      padding: EdgeInsets.all(15),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                        children: [
                          Container(
                            child: Row(
                              children: [
                                SizedBox(
                                  width: 31,
                                  child: Image.asset('images/Eurus_Blue.png',
                                      package: 'euruswallet'),
                                ),
                                SizedBox(width: 8),
                                Text(
                                  'CREATE_WALLET_PAGE.CEN.TITLE'.tr(),
                                  style: FXUI.normalTextStyle
                                      .copyWith(color: FXColor.mainBlueColor),
                                ),
                              ],
                            ),
                          ),
                          Container(
                            child: Row(
                              children: [
                                SizedBox(
                                  width: 31,
                                  child: Image.asset('images/Eurus_Violet.png',
                                      package: 'euruswallet'),
                                ),
                                SizedBox(width: 8),
                                Text(
                                  'CREATE_WALLET_PAGE.DECEN.TITLE'.tr(),
                                  style: FXUI.normalTextStyle.copyWith(
                                      color: FXColor.mainDeepBlueColor),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _getAccounts() async {
    List<UserProfile> _acs = await getLocalAcs();
    setState(() {
      acs = _acs;
    });
  }

  List<Widget> genAcListItems() {
    if (!this.acReadys) return [];

    List<Widget> _items = [];

    acs?.forEach((e) {
      _items.add(
        genAcListItem(
          Image.asset(
            'images/${e.userType == CurrentUserType.centralized ? 'Eurus_Blue' : 'Eurus_Violet'}.png',
            package: 'euruswallet',
          ),
          e.alias ??
              e.email ??
              'SWITCH_AC.${e.userType == CurrentUserType.centralized ? 'CEN' : 'DECEN'}_AC'
                  .tr(),
          subContent: getUserAddress(e),
          onTab: () => openPWDialog(e),
          secondaryActions: [
            IconSlideAction(
              color: FXColor.greyTextColor,
              iconWidget: Image.asset(
                "images/icn_delete.png",
                package: 'euruswallet',
                width: 20,
                height: 25,
              ),
              onTap: () async {
                final shouldDeleteAccount = await showDialog(
                    context: context,
                    builder: (dialogContext) =>
                        getDeleteAccountDialog(dialogContext));
                if (shouldDeleteAccount) {
                  await setAcToLocal(
                    e,
                    delete: true,
                  );
                  _getAccounts();
                }
              },
            ),
          ],
        ),
      );
    });

    _items.add(
      genAcListItem(Icon(Icons.add), 'SWITCH_AC.ADD_OTHER_AC'.tr(),
          onTab: toCenLoginPage),
    );

    return _items;
  }

  Widget genAcListItem(
    Widget icon,
    String mainContent, {
    String? subContent,
    Function()? onTab,
    List<Widget>? actions,
    List<Widget>? secondaryActions,
  }) {
    return Slidable(
      actionPane: SlidableDrawerActionPane(),
      actionExtentRatio: 0.18,
      child: Container(
        color: Colors.white,
        child: ListTile(
          leading: SizedBox(
            width: 35,
            child: icon,
          ),
          title: Text(
            mainContent,
            style: FXUI.normalTextStyle
                .copyWith(fontSize: 16, color: FXColor.lightBlack),
          ),
          subtitle: subContent != null
              ? Text(
                  subContent,
                  style: FXUI.normalTextStyle
                      .copyWith(fontSize: 12, color: FXColor.lightGray),
                  overflow: TextOverflow.ellipsis,
                )
              : Container(),
          onTap: onTab,
        ),
      ),
      actions: actions,
      secondaryActions: secondaryActions,
    );
  }

  Widget getDeleteAccountDialog(BuildContext context) {
    return Dialog(
      insetPadding: EdgeInsets.all(8),
      shape: RoundedRectangleBorder(borderRadius: FXUI.cricleRadius),
      child: Container(
        padding: EdgeInsets.fromLTRB(28, 28, 28, 20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              'COMMON.WARNING'.tr(),
              textAlign: TextAlign.center,
              style: FXUI.titleTextStyle.copyWith(
                fontSize: 18,
                fontWeight: FontWeight.w600,
              ),
            ),
            SizedBox(height: 34),
            Text(
              'SWITCH_AC.REMOVE_DESCRIPTION'.tr(),
              textAlign: TextAlign.center,
              style: FXUI.subtitleTextStyle.copyWith(
                fontSize: 14,
                fontWeight: FontWeight.normal,
              ),
            ),
            SizedBox(height: 41),
            TextButton(
              onPressed: () => Navigator.of(context).pop(true),
              style: ButtonStyle(
                textStyle: MaterialStateProperty.all<TextStyle>(
                  FXUI.titleTextStyle.copyWith(
                    fontSize: 16,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                foregroundColor: MaterialStateProperty.all<Color>(
                  Colors.white,
                ),
                backgroundColor: MaterialStateProperty.all<Color>(
                  FXColor.mainBlueColor,
                ),
                shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                  RoundedRectangleBorder(
                    borderRadius: FXUI.cricleRadius,
                  ),
                ),
                fixedSize: MaterialStateProperty.all<Size>(
                  Size.fromHeight(44),
                ),
              ),
              child: Text('COMMON.CONFIRM'.tr()),
            ),
            OutlinedButton(
              onPressed: () => Navigator.of(context).pop(false),
              style: ButtonStyle(
                textStyle: MaterialStateProperty.all<TextStyle>(
                  FXUI.titleTextStyle.copyWith(
                    fontSize: 16,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                foregroundColor: MaterialStateProperty.all<Color>(
                  FXColor.mainBlueColor,
                ),
                shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                  RoundedRectangleBorder(
                    side: BorderSide(
                      color: FXColor.mainBlueColor,
                    ),
                    borderRadius: FXUI.cricleRadius,
                  ),
                ),
                fixedSize: MaterialStateProperty.all<Size>(
                  Size.fromHeight(44),
                ),
              ),
              child: Text('COMMON.CANCEL'.tr()),
            ),
          ],
        ),
      ),
    );
  }

  void openPWDialog(UserProfile up) async {
    common.showAuthPopUp(
      context,
      currentUserType: up.userType,
      isNeededRoutingToHomePage: true,
      displayUserName: up.userType == CurrentUserType.centralized
          ? up.email
          : getUserAddress(up),
      encryptedAddress: up.encryptedAddress,
    );
  }

  void toCenLoginPage() =>
      common.pushPage(context: context, page: CenLoginPage());
}
