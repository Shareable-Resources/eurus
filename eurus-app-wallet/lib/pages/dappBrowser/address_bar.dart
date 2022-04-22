import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/dapp_browser_website_item.dart';
import 'package:euruswallet/pages/dappBrowser/dapp_browser_helper.dart';
import 'package:euruswallet/pages/dappBrowser/web3_bridge.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:flutter_neumorphic/flutter_neumorphic.dart';
import 'package:http/http.dart' as http;
import 'package:metadata_fetch/metadata_fetch.dart';

class AddressBar extends StatefulWidget {
  const AddressBar({
    Key? key,
    this.favicon,
    required this.url,
    required this.urlController,
    required this.handler,
    required this.onFocusChange,
    this.refreshHandler,
    this.collapseHandler,
    this.expandHandler,
  }) : super(key: key);

  final Favicon? favicon;
  final String url;
  final TextEditingController urlController;
  final Function handler;
  final Function onFocusChange;
  final Function()? refreshHandler;
  final Function()? collapseHandler;
  final Function()? expandHandler;

  @override
  _AddressBarState createState() => _AddressBarState();

  submit(String url) async {
    final searchPrefix = 'https://www.google.com/search?q=';
    Uri uri = Uri.tryParse(url) ?? Uri.parse('https://google.com');
    DappBrowserWebsiteItem? item;
    if (!url.contains(searchPrefix))
      try {
        final __url = (!uri.hasScheme ? 'http://' : '') + uri.toString();
        final data = await MetadataFetch.extract(__url);
        if (data != null)
          item = DappBrowserWebsiteItem(
            data.url ?? __url,
            title: data.title ?? url,
          );
      } catch (e) {}

    if (item == null) {
      final _url =
          (!url.contains(searchPrefix) ? searchPrefix : '') + uri.toString();
      item = DappBrowserWebsiteItem(
        _url,
        title: _url.replaceAll(searchPrefix, ''),
      );
    }
    DappBrowserHelper().addWebsiteItemToStorage(
        DappBrowserWebsiteItemType.searchHistory, item);
    handler(Uri.parse(item.url));
  }
}

class _AddressBarState extends State<AddressBar> {
  FocusNode _focus = new FocusNode();
  bool isExpanded = true;

  late final _urlController;

  @override
  void initState() {
    super.initState();

    _urlController = widget.urlController;

    _focus.addListener(_onFocusChange);
  }

  @override
  void dispose() {
    _focus.removeListener(_onFocusChange);
    _focus.dispose();

    super.dispose();
  }

  void _onFocusChange() {
    widget.onFocusChange(_focus.hasFocus);
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.only(left: 25, top: 8, right: 16, bottom: 8),
      child: Row(
        children: [
          if (widget.collapseHandler != null && widget.expandHandler != null)
            GestureDetector(
              onTap: () => setState(() {
                if (isExpanded) {
                  (widget.collapseHandler ?? () {})();
                } else {
                  (widget.expandHandler ?? () {})();
                }
                isExpanded = !isExpanded;
              }),
              child: Padding(
                padding: const EdgeInsets.only(left: 8, right: 16),
                child: Image.asset(
                  'images/${isExpanded ? 'icn_collapse' : 'icn_expand'}.png',
                  package: 'euruswallet',
                  width: 14,
                ),
              ),
            ),
          Expanded(
            child: FXUI.neumorphicTextField(
              context,
              padding: EdgeInsets.symmetric(vertical: 8, horizontal: 16),
              keyboardType: TextInputType.url,
              hintText: 'DAPP_BROWSER.ADDRESS_BAR_PLACEHOLDER'.tr(),
              prefixIcon: _getPrefixIcon(),
              suffixIcon: _getSuffixIcon(),
              controller: _urlController,
              focusNode: _focus,
              autocorrect: false,
              onSubmitted: (value) async {
                widget.submit(value);
              },
            ),
          ),
          Theme(
            data: Theme.of(context).copyWith(
              highlightColor: Colors.transparent,
            ),
            child: PopupMenuButton(
              icon: Image.asset(
                'images/icn_more.png',
                package: 'euruswallet',
                width: 16,
              ),
              shape: RoundedRectangleBorder(borderRadius: FXUI.cricleRadius),
              itemBuilder: (BuildContext context) =>
                  <PopupMenuEntry<String>>[
                    if (widget.refreshHandler != null)
                      PopupMenuItem(
                        value: 'DAPP_BROWSER.ADDRESS_BAR_REFRESH'.tr(),
                        child: Row(
                          children: [
                            Padding(
                              padding: const EdgeInsets.all(8),
                              child: Image.asset(
                                'images/icn_browser_refresh.png',
                                package: 'euruswallet',
                                width: 14,
                              ),
                            ),
                            Text('DAPP_BROWSER.ADDRESS_BAR_REFRESH'.tr()),
                          ],
                        ),
                        onTap: widget.refreshHandler,
                      ),
                    if (widget.refreshHandler != null) PopupMenuDivider(),
                  ] +
                  _getNetworkPopupMenuButton(),
            ),
          ),
        ],
      ),
    );
  }

  Widget _getPrefixIcon() {
    return Padding(
      padding: EdgeInsets.only(left: 0, right: 4),
      child: widget.favicon != null
          ? Image.network(
              widget.favicon!.url.toString(),
              width: 16,
              errorBuilder: (BuildContext context, Object exception,
                  StackTrace? stackTrace) {
                return Icon(
                  Icons.search,
                  size: 12,
                );
              },
            )
          : Icon(
              Icons.search,
              size: 12,
            ),
    );
  }

  Widget? _getSuffixIcon() {
    return Uri.parse(widget.url).scheme == "https"
        ? Padding(
            padding: const EdgeInsets.only(left: 8),
            child: Image.asset(
              'images/icn_lock.png',
              package: 'euruswallet',
              width: 12,
            ),
          )
        : null;
  }

  List<PopupMenuItem<String>> _getNetworkPopupMenuButton() {
    return NetworkType.values
        .map((type) => type.isEnabled
            ? PopupMenuItem(
                value: type.displayValue,
                child: Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(11.0),
                      child: Image.asset(
                        'images/icn_browser_network_status.png',
                        package: 'euruswallet',
                        width: 8,
                      ),
                    ),
                    Text(type.displayValue),
                    if (type == Web3Bridge.instance.currentNetworkType)
                      Padding(
                        padding: const EdgeInsets.only(left: 6.0),
                        child: Image.asset(
                          'images/icn_browser_network_tick.png',
                          package: 'euruswallet',
                          width: 15,
                        ),
                      ),
                  ],
                ),
                onTap: () => Web3Bridge.instance.setNetwork(type),
              )
            : null)
        .whereType<PopupMenuItem<String>>()
        .toList();
  }
}
