import 'package:euruswallet/common/commonMethod.dart';
import 'package:flutter/gestures.dart';
import 'package:share/share.dart';
import 'package:webview_flutter/webview_flutter.dart';

class WebViewPage extends StatefulWidget {
  final String? link;
  final String appTitle;

  WebViewPage({
    Key? key,
    this.link,
    this.appTitle = '',
  }) : super(key: key);

  @override
  WebViewPageState createState() => WebViewPageState();
}

class WebViewPageState extends State<WebViewPage> {
  final Set<Factory<OneSequenceGestureRecognizer>> gestureRecognizers =
      [Factory(() => EagerGestureRecognizer())].toSet();
  @override
  void initState() {
    super.initState();
    // Enable hybrid composition.
    if (Platform.isAndroid) WebView.platform = SurfaceAndroidWebView();
  }

  @override
  Widget build(BuildContext context) {

    print("widget.link:${widget.link}");
    return BackGroundImage(
        child: Scaffold(
            backgroundColor: Colors.transparent,
            appBar: WalletAppBar(
                title: widget.appTitle,
                rightWidget: Icon(Icons.ios_share),
                function: () {
                  String shareLink =
                      (common.fromBlockChainType == BlockChainType.Ethereum
                          ? api.mainNetExplorerUrl
                          : api.eurusExplorerUrl);
                  shareLink += web3dart.lastTxId ?? "";
                  Share.share(shareLink);
                }),
            body: WebView(
              javascriptMode: JavascriptMode.unrestricted,
              gestureRecognizers: gestureRecognizers,
              initialUrl: widget.link,
            )));
  }
}
