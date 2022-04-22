import 'dart:collection';

import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/dapp_browser_website_item.dart';
import 'package:euruswallet/pages/dappBrowser/dapp_browser_helper.dart';
import 'package:euruswallet/pages/dappBrowser/web3_bridge.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:metadata_fetch/metadata_fetch.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:web3dart/json_rpc.dart';

import 'address_bar.dart';
import 'dapp_browser_add_favourite_page.dart';

class Web3Browser extends StatefulWidget {
  const Web3Browser({
    Key? key,
    this.web3Bridge,
    this.initialUrl,
    this.initialInjectJs,
    this.bookmarkButtonHandler,
    this.searchButtonHandler,
    this.shouldShowAddressBar = true,
  }) : super(key: key);

  final Web3Bridge? web3Bridge;
  final String? initialUrl;
  final String? initialInjectJs;
  final Function? bookmarkButtonHandler;
  final Function? searchButtonHandler;
  final bool shouldShowAddressBar;

  @override
  _Web3BrowserState createState() => _Web3BrowserState();
}

class _Web3BrowserState extends State<Web3Browser>
    with AutomaticKeepAliveClientMixin {
  late Web3Bridge _web3bridge = widget.web3Bridge ?? Web3Bridge.instance;

  final GlobalKey webViewKey = GlobalKey();

  InAppWebViewController? webViewController;
  InAppWebViewGroupOptions options = InAppWebViewGroupOptions(
      crossPlatform: InAppWebViewOptions(
        useShouldOverrideUrlLoading: true,
        mediaPlaybackRequiresUserGesture: false,
      ),
      android: AndroidInAppWebViewOptions(
        useHybridComposition: true,
      ),
      ios: IOSInAppWebViewOptions(
        allowsInlineMediaPlayback: true,
        disallowOverScroll: true,
      ));

  PullToRefreshController? pullToRefreshController;
  String url = "";
  TextEditingController urlController = TextEditingController();
  double progress = 0;
  Favicon? favicon;
  bool _isFavoriteWebsite = false;
  bool _isButtonBarVisible = true;

  @override
  void initState() {
    super.initState();

    _web3bridge = widget.web3Bridge ?? Web3Bridge.instance;

    url = widget.initialUrl ?? "";

    pullToRefreshController = widget.shouldShowAddressBar
        ? PullToRefreshController(
            options: PullToRefreshOptions(
              color: Colors.blue,
            ),
            onRefresh: () async {
              if (Platform.isAndroid) {
                webViewController?.reload();
              } else if (Platform.isIOS) {
                webViewController
                    ?.loadUrl(
                        urlRequest:
                            URLRequest(url: await webViewController?.getUrl()))
                    .catchError((e) {
                  print(e);
                });
              }
            },
          )
        : null;
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  bool get wantKeepAlive => true;

  @override
  Widget build(BuildContext context) {
    super.build(context);
    return Scaffold(
      body: SafeArea(
        child: Column(
          children: <Widget>[
            if (widget.shouldShowAddressBar && _isButtonBarVisible)
              Padding(
                padding: const EdgeInsets.only(top: 16.0),
                child: ButtonBar(
                  alignment: MainAxisAlignment.spaceAround,
                  children: <Widget>[
                    GestureDetector(
                      child: Image.asset(
                        _isFavoriteWebsite
                            ? 'images/icn_browser_star_active.png'
                            : 'images/icn_browser_star.png',
                        package: 'euruswallet',
                        width: 25,
                        height: 25,
                      ),
                      onTap: () async {
                        if (_isFavoriteWebsite) {
                          await showDialog(
                              context: context,
                              builder: (context) => DappBrowserHelper()
                                  .getRemoveStorageItemDialog(
                                      context,
                                      DappBrowserWebsiteItemType.favorite,
                                      DappBrowserWebsiteItem(url)));
                        } else {
                          EasyLoading.show(status: 'COMMON.LOADING_W_DOT'.tr());
                          final data = await MetadataFetch.extract(url);
                          EasyLoading.dismiss();
                          final item = DappBrowserWebsiteItem(
                            url,
                            title: data?.title ?? '',
                          );
                          await Navigator.of(context).push(MaterialPageRoute(
                              builder: (_) =>
                                  DappBrowserAddFavouritePage(item: item)));
                        }
                        await _checkIsFavoriteWebsite(url);
                        setState(() {});
                      },
                    ),
                    GestureDetector(
                      child: Image.asset(
                        'images/icn_browser_bookmark_active.png',
                        package: 'euruswallet',
                        width: 25,
                        height: 25,
                      ),
                      onTap: () {
                        if (widget.bookmarkButtonHandler != null)
                          widget.bookmarkButtonHandler!();
                      },
                    ),
                    GestureDetector(
                      child: Image.asset(
                        'images/icn_browser_search.png',
                        package: 'euruswallet',
                        width: 25,
                        height: 25,
                      ),
                      onTap: () {
                        if (widget.searchButtonHandler != null)
                          widget.searchButtonHandler!();
                      },
                    ),
                    GestureDetector(
                      child: Image.asset(
                        'images/icn_browser_back.png',
                        package: 'euruswallet',
                        width: 25,
                        height: 25,
                      ),
                      onTap: () {
                        webViewController?.goBack();
                      },
                    ),
                    GestureDetector(
                      child: Image.asset(
                        'images/icn_browser_forward.png',
                        package: 'euruswallet',
                        width: 25,
                        height: 25,
                      ),
                      onTap: () {
                        webViewController?.goForward();
                      },
                    ),
                  ],
                ),
              ),
            if (widget.shouldShowAddressBar)
              AddressBar(
                favicon: favicon,
                url: url,
                urlController: urlController,
                handler: (url) {
                  webViewController?.loadUrl(urlRequest: URLRequest(url: url));
                },
                onFocusChange: (hasFocus) {},
                refreshHandler: () {
                  webViewController?.reload();
                },
                collapseHandler: () {
                  setState(() {
                    _isButtonBarVisible = false;
                  });
                },
                expandHandler: () {
                  setState(() {
                    _isButtonBarVisible = true;
                  });
                },
              ),
            Expanded(
              child: Stack(
                children: [
                  InAppWebView(
                    key: webViewKey,
                    initialUrlRequest: URLRequest(url: Uri.parse(url)),
                    initialOptions: options,
                    initialUserScripts: UnmodifiableListView([
                      if (!isEmptyString(string: widget.initialInjectJs))
                        UserScript(
                            source: widget.initialInjectJs ?? '',
                            injectionTime:
                                UserScriptInjectionTime.AT_DOCUMENT_START),
                    ]),
                    pullToRefreshController: pullToRefreshController,
                    onWebViewCreated: (controller) async {
                      webViewController = controller;

                      _web3bridge.webViewController = controller;
                      _web3bridge.context = context;
                      JsonRpcApiMethod.values.forEach((element) {
                        controller.addJavaScriptHandler(
                            handlerName: element.displayValue,
                            callback: (args) async {
                              try {
                                return await _web3bridge.request(
                                  element.displayValue,
                                  args,
                                );
                              } on RPCError catch (e) {
                                return {
                                  'code': e.errorCode,
                                  'message': e.message,
                                  'data': e.data,
                                };
                              } catch (e) {
                                return {
                                  'code': -10001,
                                  'message': e.toString(),
                                };
                              }
                            });
                      });
                    },
                    onLoadStart: (controller, url) async {
                      if (widget.shouldShowAddressBar) {
                        final _favicons = await controller.getFavicons();
                        this.favicon =
                            _favicons.isNotEmpty ? _favicons.first : null;
                      }

                      setState(() {
                        this.url = url.toString();
                        urlController.text = this.url;
                      });
                    },
                    androidOnPermissionRequest:
                        (controller, origin, resources) async {
                      return PermissionRequestResponse(
                          resources: resources,
                          action: PermissionRequestResponseAction.GRANT);
                    },
                    shouldOverrideUrlLoading:
                        (controller, navigationAction) async {
                      var uri = navigationAction.request.url!;

                      if (![
                        "http",
                        "https",
                        "file",
                        "chrome",
                        "data",
                        "javascript",
                        "about"
                      ].contains(uri.scheme)) {
                        if (await canLaunch(url)) {
                          // Launch the App
                          await launch(
                            url,
                          );
                          // and cancel the request
                          return NavigationActionPolicy.CANCEL;
                        }
                      }

                      return NavigationActionPolicy.ALLOW;
                    },
                    onLoadStop: (controller, url) async {
                      await _checkIsFavoriteWebsite(url.toString());
                      pullToRefreshController?.endRefreshing();
                      setState(() {
                        this.url = url.toString();
                        urlController.text = this.url;
                      });
                    },
                    onLoadError: (controller, url, code, message) {
                      pullToRefreshController?.endRefreshing();
                    },
                    onProgressChanged: (controller, progress) {
                      if (progress == 100) {
                        pullToRefreshController?.endRefreshing();
                      }
                      setState(() {
                        this.progress = progress / 100;
                        urlController.text = this.url;
                      });
                    },
                    onUpdateVisitedHistory: (controller, url, androidIsReload) {
                      setState(() {
                        this.url = url.toString();
                        urlController.text = this.url;
                      });
                    },
                    onConsoleMessage: (controller, consoleMessage) {
                      print(consoleMessage);
                    },
                  ),
                  progress < 1.0
                      ? LinearProgressIndicator(value: progress)
                      : Container(),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future _checkIsFavoriteWebsite(String url) async {
    final item = DappBrowserWebsiteItem(url);
    final items = await DappBrowserHelper()
        .getWebsiteItemsFromStorage(DappBrowserWebsiteItemType.favorite);
    _isFavoriteWebsite = items.contains(item);
  }
}
