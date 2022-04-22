import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/dapp_browser_website_item.dart';
import 'package:favicon/favicon.dart' as favicon;
import 'package:flutter/widgets.dart';
import 'package:flutter_neumorphic/flutter_neumorphic.dart';
import 'package:flutter_svg/svg.dart';

import 'address_bar.dart';
import 'dapp_browser_helper.dart';
import 'web3_browser.dart';

class DappBrowserHomePage extends StatefulWidget {
  const DappBrowserHomePage({Key? key}) : super(key: key);

  @override
  _DappBrowserHomePageState createState() => _DappBrowserHomePageState();
}

class _DappBrowserHomePageState extends State<DappBrowserHomePage> {
  bool shouldShowWeb3Browser = false;
  bool hasFocus = false;

  String? url;
  TextEditingController urlController = TextEditingController();
  List<DappBrowserWebsiteItem> favoriteWebsiteItems = [];

  AddressBar get _addressBar {
    return AddressBar(
      url: url ?? '',
      urlController: urlController,
      handler: (uri) {
        setState(() {
          this.url = uri.toString();
          this.urlController.text = uri.toString();
          shouldShowWeb3Browser = true;
        });
      },
      onFocusChange: (_hasFocus) {
        setState(() {
          hasFocus = _hasFocus;
        });
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);

    return shouldShowWeb3Browser
        ? FutureBuilder<String>(
            future: rootBundle.loadString(
                'packages/euruswallet/lib/pages/dappBrowser/web3_ethereum_setup.js'),
            builder: (BuildContext context, AsyncSnapshot<String> snapshot) {
              if (snapshot.hasData) {
                final source = !isEmptyString(string: snapshot.data)
                    ? (snapshot.data ?? '') +
                        (isCentralized()
                            ? "\nethereum.ownerWalletAddress = '${common.ownerWalletAddress ?? ''}';"
                            : '')
                    : null;
                return Web3Browser(
                  initialUrl: url,
                  initialInjectJs: source,
                  bookmarkButtonHandler: () {
                    setState(() {
                      shouldShowWeb3Browser = false;
                      hasFocus = false;
                    });
                  },
                  searchButtonHandler: () {
                    setState(() {
                      shouldShowWeb3Browser = false;
                      hasFocus = true;
                    });
                  },
                );
              } else {
                return Container();
              }
            })
        : Scaffold(
            body: SafeArea(
              child: Column(
                children: [
                  SizedBox(height: 16),
                  Padding(
                    padding: EdgeInsets.symmetric(vertical: 4, horizontal: 25),
                    child: Row(
                      mainAxisAlignment: hasFocus
                          ? MainAxisAlignment.center
                          : MainAxisAlignment.start,
                      children: [
                        Image.asset(
                          'images/imgLogoH.png',
                          package: 'euruswallet',
                        ),
                      ],
                    ),
                  ),
                  _addressBar,
                  Expanded(
                    child: hasFocus
                        ? getSearchHistorySection()
                        : getFavoritesSection(),
                  ),
                ],
              ),
            ),
          );
  }

  Widget _getSectionTitleRow(String? iconPath, String title) {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: 6, horizontal: 35),
      child: Row(
        children: [
          if (!isEmptyString(string: iconPath))
            Image.asset(
              iconPath ?? '',
              package: 'euruswallet',
              width: 10,
              color: FXColor.middleBlack,
            ),
          SizedBox(width: 8),
          Text(
            title,
            style: FXUI.titleTextStyle.copyWith(
              color: FXColor.middleBlack,
              fontSize: 13,
            ),
          ),
        ],
      ),
    );
  }

  Widget getFavoritesSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _getSectionTitleRow(
          'images/icn_star.png',
          'DAPP_BROWSER.RECOMMENDED'.tr(),
        ),
        SingleChildScrollView(
          padding: const EdgeInsets.symmetric(horizontal: 16.0),
          scrollDirection: Axis.horizontal,
          child: Row(
            children: DefaultDappType.values
                .map((e) => getRecommendedDappWidget(e))
                .toList(),
          ),
        ),
        _getSectionTitleRow(
          'images/icn_star.png',
          'DAPP_BROWSER.FAVOURITE'.tr(),
        ),
        SizedBox(
          height: 4,
        ),
        Expanded(
          child: FutureBuilder<List<DappBrowserWebsiteItem>?>(
            future: DappBrowserHelper().getWebsiteItemsFromStorage(
                DappBrowserWebsiteItemType.favorite),
            builder: (context, snapshot) {
              if (snapshot.connectionState != ConnectionState.done ||
                  !snapshot.hasData ||
                  snapshot.hasError) return Container();

              final items = snapshot.data ?? [];
              this.favoriteWebsiteItems = items;
              return ListView.builder(
                itemCount: items.length,
                itemBuilder: (context, index) {
                  final item = items[index];
                  return getFavoritesItem(item);
                },
              );
            },
          ),
        ),
      ],
    );
  }

  Widget getRecommendedDappWidget(DefaultDappType dappType) {
    return GestureDetector(
      onTap: () => setState(() {
        shouldShowWeb3Browser = true;
        url = dappType.url;
      }),
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Container(
          width: 150,
          height: 150,
          decoration: BoxDecoration(
            color: FXColor.middleBlack,
            borderRadius: BorderRadius.circular(6),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Image.asset(
                  dappType.logoPath,
                  package: 'euruswallet',
                  height: 25,
                ),
                SizedBox(height: 10),
                Expanded(
                  child: Text(
                    dappType.description,
                    style: FXUI.normalTextStyle.copyWith(
                      color: FXColor.veryLightGray,
                      fontSize: 11,
                      fontWeight: FontWeight.normal,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget getFavoritesItem(DappBrowserWebsiteItem item) {
    return Neumorphic(
      margin: EdgeInsets.symmetric(vertical: 9, horizontal: 29),
      style: NeumorphicStyle(color: Colors.white),
      child: ListTile(
        contentPadding: EdgeInsets.only(top: 5, left: 18, right: 5, bottom: 8),
        leading: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              (item.title ?? '').isNotEmpty
                  ? (item.title ?? '').substring(0, 1)
                  : '',
              style: FXUI.inputStyle
                  .copyWith(fontSize: 32, fontWeight: FontWeight.w500),
            ),
          ],
        ),
        title: Text(item.title ?? ''),
        subtitle: Text(item.url),
        trailing: Column(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            IconButton(
              iconSize: 14,
              alignment: Alignment.topRight,
              icon: Icon(Icons.close),
              onPressed: () async {
                await showDialog(
                  context: context,
                  builder: (dialogContext) =>
                      DappBrowserHelper().getRemoveStorageItemDialog(
                    dialogContext,
                    DappBrowserWebsiteItemType.favorite,
                    item,
                  ),
                );
                setState(() {});
              },
              highlightColor: Colors.transparent,
              splashColor: Colors.transparent,
            ),
          ],
        ),
        onTap: () {
          setState(() {
            url = item.url;
            urlController.text = item.url;
            shouldShowWeb3Browser = true;
          });
        },
      ),
    );
  }

  Widget getSearchHistorySection() {
    return FutureBuilder<List<DappBrowserWebsiteItem>>(
      future: DappBrowserHelper()
          .getWebsiteItemsFromStorage(DappBrowserWebsiteItemType.searchHistory),
      builder: (context, snapshot) {
        if (snapshot.connectionState != ConnectionState.done ||
            !snapshot.hasData ||
            snapshot.hasError) return Container();

        final items = snapshot.data ?? [];
        return Column(
          children: [
            _getSectionTitleRow(
              'images/icn_search_history.png',
              'DAPP_BROWSER.SEARCH_HISTORY'.tr(),
            ),
            SizedBox(height: 4),
            Expanded(
              child: ListView.builder(
                itemCount: items.length,
                itemBuilder: (context, index) {
                  final item = items[index];
                  return FutureBuilder<favicon.Icon?>(
                    future: favicon.Favicon.getBest(item.url),
                    builder: (context, snapshot) {
                      return ListTile(
                        leading: snapshot.hasData && snapshot.data != null
                            ? Image.network(
                                snapshot.data!.url,
                                width: 14,
                                errorBuilder: (_, __, ___) =>
                                    snapshot.data!.url.endsWith('svg')
                                        ? SvgPicture.network(
                                            snapshot.data!.url,
                                            width: 14,
                                          )
                                        : Text(''),
                              )
                            : Text(''),
                        title: Text(item.title ?? item.url),
                        trailing: IconButton(
                          iconSize: 14,
                          alignment: Alignment.topRight,
                          icon: Icon(Icons.close),
                          onPressed: () async {
                            await showDialog(
                              context: context,
                              builder: (dialogContext) => DappBrowserHelper()
                                  .getRemoveStorageItemDialog(
                                dialogContext,
                                DappBrowserWebsiteItemType.searchHistory,
                                item,
                              ),
                            );
                            setState(() {});
                          },
                          highlightColor: Colors.transparent,
                          splashColor: Colors.transparent,
                        ),
                        onTap: () {
                          setState(() {
                            url = item.url;
                            urlController.text = item.url;
                            _addressBar.submit(item.url);
                          });
                        },
                      );
                    },
                  );
                },
              ),
            ),
          ],
        );
      },
    );
  }
}
