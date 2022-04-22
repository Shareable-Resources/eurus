import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/dapp_browser_website_item.dart';
import 'package:euruswallet/pages/dappBrowser/dapp_browser_helper.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_neumorphic/flutter_neumorphic.dart';

class DappBrowserAddFavouritePage extends StatefulWidget {
  const DappBrowserAddFavouritePage({
    Key? key,
    this.item,
  }) : super(key: key);

  final DappBrowserWebsiteItem? item;

  @override
  _DappBrowserAddFavouritePageState createState() =>
      _DappBrowserAddFavouritePageState();
}

class _DappBrowserAddFavouritePageState
    extends State<DappBrowserAddFavouritePage> {
  final TextEditingController _nameController = TextEditingController();
  final TextEditingController _urlController = TextEditingController();

  @override
  void initState() {
    _nameController.text = widget.item?.title ?? '';
    _urlController.text = widget.item?.url ?? '';
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: FXColor.lightWhiteColor,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        foregroundColor: Colors.black,
        centerTitle: true,
        elevation: 0.0,
        title: Text('DAPP_BROWSER.ADD_FAVOURITES'.tr()),
      ),
      body: Neumorphic(
        margin: EdgeInsets.symmetric(vertical: 32, horizontal: 24),
        padding: EdgeInsets.fromLTRB(14, 40, 14, 24),
        style: NeumorphicStyle(
          color: Colors.white,
          boxShape: NeumorphicBoxShape.roundRect(FXUI.cricleRadius),
          lightSource: LightSource(0, 0),
          depth: 10,
        ),
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(
                'DAPP_BROWSER.ADD_FAVOURITES_NAME'.tr(),
                style: FXUI.titleTextStyle.copyWith(
                  color: FXColor.placeholderGreyColor,
                  fontSize: 16,
                  fontWeight: FontWeight.normal,
                ),
              ),
              SizedBox(height: 8),
              FXUI.neumorphicTextField(
                context,
                textInputAction: TextInputAction.next,
                keyboardType: TextInputType.name,
                padding: EdgeInsets.symmetric(vertical: 18, horizontal: 18),
                hintText: 'DAPP_BROWSER.ADD_FAVOURITES_NAME_PLACEHOLDER'.tr(),
                prefixIcon: Padding(
                  padding: EdgeInsets.only(left: 4, right: 12),
                  child: Image.asset(
                    'images/icn_browser_bookmark.png',
                    package: 'euruswallet',
                    width: 15,
                  ),
                ),
                controller: _nameController,
              ),
              SizedBox(height: 14),
              Text(
                'DAPP_BROWSER.ADD_FAVOURITES_URL'.tr(),
                style: FXUI.titleTextStyle.copyWith(
                  color: FXColor.placeholderGreyColor,
                  fontSize: 16,
                  fontWeight: FontWeight.normal,
                ),
              ),
              SizedBox(height: 8),
              FXUI.neumorphicTextField(
                context,
                keyboardType: TextInputType.url,
                padding: EdgeInsets.symmetric(vertical: 18, horizontal: 18),
                hintText: 'DAPP_BROWSER.ADD_FAVOURITES_URL_PLACEHOLDER'.tr(),
                prefixIcon: Padding(
                  padding: EdgeInsets.only(left: 4, right: 12),
                  child: Image.asset(
                    'images/icn_browser.png',
                    package: 'euruswallet',
                    width: 15,
                  ),
                ),
                controller: _urlController,
                onSubmitted: (value) async {
                  await _submit();
                },
              ),
              SizedBox(height: 27),
              TextButton(
                onPressed: () async {
                  await _submit();
                },
                style: ButtonStyle(
                  textStyle: MaterialStateProperty.all<TextStyle>(
                      FXUI.titleTextStyle.copyWith(
                    fontSize: 16,
                    fontWeight: FontWeight.w500,
                  )),
                  foregroundColor:
                      MaterialStateProperty.all<Color>(Colors.white),
                  backgroundColor: MaterialStateProperty.all<Color>(
                      common.getBackGroundColor()),
                  shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                      RoundedRectangleBorder(
                    borderRadius: FXUI.cricleRadius,
                  )),
                  fixedSize:
                      MaterialStateProperty.all<Size>(Size.fromHeight(44)),
                ),
                child: Text("COMMON.CONFIRM".tr()),
              ),
              SizedBox(height: 13),
              TextButton(
                onPressed: () => Navigator.pop(context),
                style: ButtonStyle(
                  textStyle: MaterialStateProperty.all<TextStyle>(
                      FXUI.titleTextStyle.copyWith(
                    fontSize: 16,
                    fontWeight: FontWeight.w500,
                  )),
                  foregroundColor: MaterialStateProperty.all<Color>(
                      common.getBackGroundColor()),
                  shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                      RoundedRectangleBorder(
                    side: BorderSide(
                      color: common.getBackGroundColor(),
                    ),
                    borderRadius: FXUI.cricleRadius,
                  )),
                  fixedSize:
                      MaterialStateProperty.all<Size>(Size.fromHeight(44)),
                ),
                child: Text('COMMON.CANCEL'.tr()),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future _submit() async {
    final item = DappBrowserWebsiteItem(_urlController.text,
        title: _nameController.text);
    await DappBrowserHelper().addWebsiteItemToStorage(
      DappBrowserWebsiteItemType.favorite,
      item,
    );
    Navigator.pop(context);
  }
}
