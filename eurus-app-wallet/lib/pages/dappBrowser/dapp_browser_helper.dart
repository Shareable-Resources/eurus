import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/dapp_browser_website_item.dart';

enum DefaultDappType { mappedSwap, raijinSwap }

extension DefaultDappTypeExtension on DefaultDappType {
  String get url {
    switch (this) {
      case DefaultDappType.mappedSwap:
        switch (envType) {
          case EnvType.Dev:
            return 'https://decatsdevapp.eurus.dev/';
          case EnvType.Testnet:
            return 'https://testapp.mappedswap.io/';
          case EnvType.Staging:
          case EnvType.Production:
            return 'https://app.mappedswap.io/';
        }

      case DefaultDappType.raijinSwap:
        switch (envType) {
          case EnvType.Dev:
            return 'https://devswap.raijinswap.org/';
          case EnvType.Testnet:
            return 'https://testnet.raijinswap.org/';
          case EnvType.Staging:
          case EnvType.Production:
            return 'https://app.raijinswap.org/';
        }
    }
  }

  String get logoPath {
    switch (this) {
      case DefaultDappType.mappedSwap:
        return 'images/icn_mapped_swap_logo.png';
      case DefaultDappType.raijinSwap:
        return 'images/icn_raijin_swap_logo.png';
    }
  }

  String get description {
    switch (this) {
      case DefaultDappType.mappedSwap:
        return 'DAPP_BROWSER.MAPPED_SWAP_DESCRIPTION'.tr();
      case DefaultDappType.raijinSwap:
        return 'DAPP_BROWSER.RAIJIN_SWAP_DESCRIPTION'.tr();
    }
  }
}

class DappBrowserHelper {
  static final DappBrowserHelper _dappBrowserHelper =
      DappBrowserHelper._internal();

  factory DappBrowserHelper() {
    return _dappBrowserHelper;
  }

  DappBrowserHelper._internal();

  Future<List<DappBrowserWebsiteItem>> getWebsiteItemsFromStorage(
      DappBrowserWebsiteItemType type) async {
    final value = await NormalStorageKit().readValue(type.storageKey);
    final list =
        !isEmptyString(string: value) ? (jsonDecode(value!) as List) : [];
    return list.map((s) => DappBrowserWebsiteItem.fromJson(s)).toList();
  }

  Future _setWebsitesItemsToStorage(
    DappBrowserWebsiteItemType type,
    List<DappBrowserWebsiteItem> websiteItems,
  ) async {
    await NormalStorageKit()
        .setValue(jsonEncode(websiteItems), type.storageKey);
  }

  Future addWebsiteItemToStorage(
      DappBrowserWebsiteItemType type, DappBrowserWebsiteItem item) async {
    List<DappBrowserWebsiteItem> websiteItems =
        await getWebsiteItemsFromStorage(type);
    if (websiteItems.contains(item)) websiteItems.remove(item);
    websiteItems.insert(0, item);
    await _setWebsitesItemsToStorage(type, websiteItems);
  }

  Future removeWebsiteItemFromStorage(
      DappBrowserWebsiteItemType type, DappBrowserWebsiteItem item) async {
    List<DappBrowserWebsiteItem> websiteItems =
        await getWebsiteItemsFromStorage(type);
    if (websiteItems.contains(item)) websiteItems.remove(item);
    await _setWebsitesItemsToStorage(type, websiteItems);
  }

  Future removeAllWebsiteItemFromStorage(
      DappBrowserWebsiteItemType type, DappBrowserWebsiteItem item) async {
    await _setWebsitesItemsToStorage(type, []);
  }

  Widget getRemoveStorageItemDialog(BuildContext context,
      DappBrowserWebsiteItemType type, DappBrowserWebsiteItem item) {
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
              onPressed: () async {
                await DappBrowserHelper().removeWebsiteItemFromStorage(
                  type,
                  item,
                );
                Navigator.of(context).pop();
              },
              style: ButtonStyle(
                textStyle: MaterialStateProperty.all<TextStyle>(
                    FXUI.titleTextStyle.copyWith(
                  fontSize: 16,
                  fontWeight: FontWeight.w500,
                )),
                foregroundColor: MaterialStateProperty.all<Color>(Colors.white),
                backgroundColor: MaterialStateProperty.all<Color>(
                    common.getBackGroundColor()),
                shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                    RoundedRectangleBorder(
                  borderRadius: FXUI.cricleRadius,
                )),
                fixedSize: MaterialStateProperty.all<Size>(Size.fromHeight(44)),
              ),
              child: Text('COMMON.CONFIRM'.tr()),
            ),
            SizedBox(height: 14),
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              style: ButtonStyle(
                textStyle: MaterialStateProperty.all<TextStyle>(
                    FXUI.titleTextStyle.copyWith(
                  fontSize: 16,
                  fontWeight: FontWeight.w500,
                )),
                foregroundColor:
                    MaterialStateProperty.all<Color>(FXColor.mainBlueColor),
                shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                    RoundedRectangleBorder(
                  side: BorderSide(
                    color: FXColor.mainBlueColor,
                  ),
                  borderRadius: FXUI.cricleRadius,
                )),
                fixedSize: MaterialStateProperty.all<Size>(Size.fromHeight(44)),
              ),
              child: Text('COMMON.CANCEL'.tr()),
            ),
          ],
        ),
      ),
    );
  }
}
