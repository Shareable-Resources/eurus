library EurusWalletHomeUIKit;


import 'dart:convert';
import 'dart:developer';
import 'dart:io';

import 'package:app_authentication_kit/app_authentication_kit.dart';
import 'package:curved_navigation_bar/curved_navigation_bar.dart';
import 'package:eurus/widgets/zoomRoute.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:intl/date_symbol_data_local.dart';
import 'package:intl/intl.dart';
import 'package:keyboard_avoider/keyboard_avoider.dart';
import 'package:pointycastle/api.dart' hide Padding;
import 'package:app_security_kit/app_security_kit.dart';
import 'package:app_storage_kit/app_storage_kit.dart';
import 'package:livechat_inc/livechat_inc.dart';
import 'dart:ui' as ui;
import 'dart:async';
import 'package:http/http.dart' as http;

// void main() {
//   runApp(MyApp());
// }

const TEST_NET = 'Rinkeby' ?? null;

class MyApp extends StatelessWidget {
  final List<Widget> Function(BuildContext, [Map<String, dynamic>]) homeFunctionsBarItemBuilder;
  final String Function(dynamic, dynamic) canEncryptHandler;
  final String Function(dynamic, dynamic) canDecryptHandler;
  final Future<dynamic> Function(dynamic) canPersistSecurelyHandler;
  final Future<dynamic> Function(dynamic) canPersistNormallyHandler;
  final Future<dynamic> Function(dynamic, {bool delete}) canPersistWithBiometricSecurelyHandler;
  final List<Widget Function(BuildContext) Function([Map<String, dynamic>])> navigationBarOnTapTabBuilders;
  final WidgetBuilder Function(String) getEditAssetsListWidgetBuilder; 
  final Future<Widget> Function(String, double, {String imgUrl, Widget placeholder}) canGetCryptoIconHandler; 
  final Future<String> Function(String, String, [int]) canGetERC20BalanceHandler;
  final Future<bool> Function() canSupportsBiometricAuthenticatedHandler;
  // final Future<bool> Function(dynamic) initEthClientWithCanGetPrivateKeyHandler;
  final dynamic Function(dynamic) initEthClientWithCanGetPrivateKeyHandler;
  // final Future<Object> Function(String, int, [String]) canFetchTransactionHistoryHandler;
  final Future<List<Map<String, Object>>> Function(String, int, [String]) canFetchTransactionHistoryHandler;
  final String Function(String, String) canGetSignature;
  final Future<Null> Function({String centerMessage, String fromAddress, String toAddress, String txId, String date, String gasFeeString, String transferAmount, BuildContext navigatorContext, bool isAssetAllocation, bool shouldSkipPendingFetch, int blockChainType}) canNavigateToTransactionHistoryDetailPage;

  const MyApp({Key key, this.homeFunctionsBarItemBuilder, this.canEncryptHandler, this.canDecryptHandler, this.canPersistSecurelyHandler, this.canPersistNormallyHandler, this.canPersistWithBiometricSecurelyHandler, this.navigationBarOnTapTabBuilders, this.getEditAssetsListWidgetBuilder, this.canGetCryptoIconHandler, this.canGetERC20BalanceHandler, this.canSupportsBiometricAuthenticatedHandler, this.initEthClientWithCanGetPrivateKeyHandler, this.canFetchTransactionHistoryHandler, this.canGetSignature, this.canNavigateToTransactionHistoryDetailPage}) : super(key: key);

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    Intl.defaultLocale = 'en_HK';
    initializeDateFormatting('en_HK', null);
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      initialRoute: 'HomePage',
      routes: {
        'HomePage': (_) => EurusWalletHomePage(
          homeFunctionsBarItemBuilder: homeFunctionsBarItemBuilder,
          canEncryptHandler: canEncryptHandler,
          canPersistSecurelyHandler: canPersistSecurelyHandler,
          canPersistNormallyHandler: canPersistNormallyHandler,
          canPersistWithBiometricSecurelyHandler: canPersistWithBiometricSecurelyHandler,
          canDecryptHandler: canDecryptHandler,
          navigationBarOnTapTabBuilders: navigationBarOnTapTabBuilders,
          getEditAssetsListWidgetBuilder: getEditAssetsListWidgetBuilder,
          canGetCryptoIconHandler: canGetCryptoIconHandler,
          canGetERC20BalanceHandler: canGetERC20BalanceHandler,
          canSupportsBiometricAuthenticatedHandler: canSupportsBiometricAuthenticatedHandler,
          initEthClientWithCanGetPrivateKeyHandler: initEthClientWithCanGetPrivateKeyHandler,
          canFetchTransactionHistoryHandler: canFetchTransactionHistoryHandler,
          canGetSignature: canGetSignature,
          canNavigateToTransactionHistoryDetailPage: canNavigateToTransactionHistoryDetailPage,
        ),
      },
    );
  }
}

class EurusWalletHomePage extends StatefulWidget {
  EurusWalletHomePage({
    Key key, 
    // this.homeFunctionsBarOnPressedHandlers, 
    @required this.homeFunctionsBarItemBuilder,
    this.canEncryptHandler,
    this.canDecryptHandler,
    this.canPersistSecurelyHandler,
    this.canPersistNormallyHandler,
    this.canPersistWithBiometricSecurelyHandler,
    this.navigationBarOnTapTabBuilders,
    this.getEditAssetsListWidgetBuilder,
    this.canGetCryptoIconHandler,
    this.canGetERC20BalanceHandler,
    this.canSupportsBiometricAuthenticatedHandler,
    this.initEthClientWithCanGetPrivateKeyHandler,
    this.canFetchTransactionHistoryHandler,
    this.canGetSignature,
    this.canNavigateToTransactionHistoryDetailPage,
  }) : assert(homeFunctionsBarItemBuilder != null), super(key: key);

  String _address;
  String get address => _address;
  void set address(String newValue) {
    _address = newValue;
    // hacky trick to enforce override web3dart previously initialized member variable of the shared instance
    // by invoking getbalance via the web3dart shared instance per each wallet public address once newly set through this address setter
    canGetERC20BalanceHandler(newValue, '0x0', 0);
  }
  String encryptedAddress;
  String createdAccountGeneratedMnemonicSeedPhrase;
  // List<Function> homeFunctionsBarOnPressedHandlers;
  final List<Widget> Function(BuildContext, [Map<String, dynamic>]) homeFunctionsBarItemBuilder;
  final String Function(dynamic, dynamic) canEncryptHandler;
  final String Function(dynamic, dynamic) canDecryptHandler;
  final Future<dynamic> Function(dynamic) canPersistSecurelyHandler;
  final Future<dynamic> Function(dynamic) canPersistNormallyHandler;
  final Future<dynamic> Function(dynamic, {bool delete}) canPersistWithBiometricSecurelyHandler;
  final List<Widget Function(BuildContext) Function([Map<String, dynamic>])> navigationBarOnTapTabBuilders;
  final WidgetBuilder Function(String) getEditAssetsListWidgetBuilder;
  final Future<Widget> Function(String, double, {String imgUrl, Widget placeholder}) canGetCryptoIconHandler;
  final Future<String> Function(String, String, [int]) canGetERC20BalanceHandler;
  final Future<bool> Function() canSupportsBiometricAuthenticatedHandler;
  // final Future<bool> Function(dynamic) initEthClientWithCanGetPrivateKeyHandler;
  final dynamic Function(dynamic) initEthClientWithCanGetPrivateKeyHandler;
  // final Future<Object> Function(String, int, [String]) canFetchTransactionHistoryHandler;
  final Future<List<Map<String, Object>>> Function(String, int, [String]) canFetchTransactionHistoryHandler;
  final String Function(String, String) canGetSignature;
  final Future<Null> Function({String centerMessage, String fromAddress, String toAddress, String txId, String date, String gasFeeString, String transferAmount, BuildContext navigatorContext, bool isAssetAllocation, bool shouldSkipPendingFetch, int blockChainType}) canNavigateToTransactionHistoryDetailPage;

  @override
  _EurusWalletHomePageState createState() => _EurusWalletHomePageState();
}

class _EurusWalletHomePageState extends State<EurusWalletHomePage> {
  int navigationBarTabBuilderSelectedIndex;
  bool showBalance = true;
  String __prefix;

  @override
  void initState() {
    super.initState();
    print('_EurusWalletHomePageState initState()');
    final _getPrefix = () async {
      __prefix = await widget.canPersistNormallyHandler('APP_INSTALL_FIRST_LAUNCH_ENTER_WALLET');
      __prefix = __prefix;
      if (__prefix == null || __prefix.isEmpty) {
        __prefix = DateTime.now().microsecondsSinceEpoch.toString();
        await widget.canPersistNormallyHandler({'APP_INSTALL_FIRST_LAUNCH_ENTER_WALLET': __prefix});
      }
    };
    _getPrefix();
    
    // set Auth UI of private key Handler to Web3 ETH Client
    if (mounted) {
      final Future<String> Function() _ = () async {
        final _textEditingController = TextEditingController();
        // final _submit = () async {
        //   final _accountPrivateKeys = await widget.canPersistSecurelyHandler(__prefix+'accountPrivateKeys');
        final _submit = (_accountPrivateKeys, _accountEncryptedAddress, _ctx) {
          print('initEthClientWithCanGetPrivateKeyHandler _submit _accountEncryptedAddress $_accountEncryptedAddress');
          final _encryptedPrivateKey = Uri(query: _accountPrivateKeys).queryParameters[_accountEncryptedAddress];
          final _privateKey = widget.canDecryptHandler(_textEditingController.text, _encryptedPrivateKey);
          if (_privateKey != null) return Navigator.of(_ctx).maybePop(_privateKey);
          // else, auth failure should be prompted
          Scaffold.of(_ctx).showSnackBar(SnackBar(content: Text("Authentication Failure", textAlign: TextAlign.center,), backgroundColor: Colors.redAccent.shade100.withOpacity(.9),));
        };

        return await Navigator.of(context).push(PageRouteBuilder(fullscreenDialog: true, opaque: false, pageBuilder: (pageBuilderContext, animation, secondaryAnimation) => FutureBuilder(
          future: widget.canPersistSecurelyHandler(__prefix+'accountPrivateKeys').then((value) => value as String),
          builder: (fCtx, s) => s.connectionState != ConnectionState.done ? Center(child: CircularProgressIndicator()) : Scaffold(
            backgroundColor: Colors.black87,
            body: Builder(builder: (_scaffoldInnerContext) {
              final __submit = () => _submit(s.data, widget.encryptedAddress, _scaffoldInnerContext);
              final _tryBioAuth = () => Future.delayed(Duration.zero, () async => await widget.canPersistWithBiometricSecurelyHandler(__prefix+widget.encryptedAddress)).then((value) {if (value != null) {_textEditingController.text = value; __submit();}});
              _tryBioAuth();
              return 
                KeyboardAvoider(autoScroll: true, child: Column(mainAxisAlignment: MainAxisAlignment.center, crossAxisAlignment: CrossAxisAlignment.stretch, children: [Container(decoration: BoxDecoration(borderRadius: BorderRadius.circular(15), color: Colors.white.withOpacity(.9)), padding: EdgeInsets.all(16), child: Column(children: [
                  ListTile(title: Text("Security Authentication", style: Theme.of(_scaffoldInnerContext).textTheme.headline6.apply(fontWeightDelta: 2), textAlign: TextAlign.center)),
                  ListTile(title: TextField(textInputAction: TextInputAction.go, onSubmitted: (_) => __submit(), autofocus: true, obscureText: true, controller: _textEditingController, decoration: InputDecoration(border: OutlineInputBorder(borderRadius: BorderRadius.circular(16)), suffixIcon: null))),
                  ListTile(title: TextButton(child: Text("Submit", style: Theme.of(_scaffoldInnerContext).textTheme.caption.apply(color: Colors.blue)), onPressed: () {
                    __submit();
                  },), subtitle: TextButton(child: Text("Cancel"), onPressed: () => Navigator.of(_scaffoldInnerContext).maybePop()))
                ],),)]));
            })
          )
        ))).then((value) => value as String);
      };
      widget.initEthClientWithCanGetPrivateKeyHandler(_);
    }
  }

  Future<String> _getPrivateKeyHandler() async {
    final _textEditingController = TextEditingController();
    // final _submit = () async {
    //   final _accountPrivateKeys = await widget.canPersistSecurelyHandler(__prefix+'accountPrivateKeys');
    final _submit = (_accountPrivateKeys, _accountEncryptedAddress, _ctx) {
      print('initEthClientWithCanGetPrivateKeyHandler _submit _accountEncryptedAddress $_accountEncryptedAddress');
      final _encryptedPrivateKey = Uri(query: _accountPrivateKeys).queryParameters[_accountEncryptedAddress];
      final _privateKey = widget.canDecryptHandler(_textEditingController.text, _encryptedPrivateKey);
      if (_privateKey != null) return Navigator.of(_ctx).maybePop(_privateKey);
      // else, auth failure should be prompted
      Scaffold.of(_ctx).showSnackBar(SnackBar(content: Text("Authentication Failure", textAlign: TextAlign.center,), backgroundColor: Colors.redAccent.shade100.withOpacity(.9),));
    };

    return await Navigator.of(context).push(MaterialPageRoute(fullscreenDialog: true, builder: (ctx) => FutureBuilder(
      future: widget.canPersistSecurelyHandler(__prefix+'accountPrivateKeys').then((value) => value as String),
      builder: (fCtx, s) => s.connectionState != ConnectionState.done ? Center(child: CircularProgressIndicator()) : Scaffold(
        backgroundColor: Colors.black12,
        body: Center(child: Column(children: [
          Expanded(child: Container()),
          TextField(controller: _textEditingController),
          // CupertinoButton(child: Text("Authenticate"), onPressed: () async {await _submit();}),
          // CupertinoButton(child: Text("Authenticate"), onPressed: () {print('initEthClientWithCanGetPrivateKeyHandler Authenticate onPressed current account encryptedAddress ${widget.encryptedAddress}'); _submit(s.data, widget.encryptedAddress, fCtx);}),
          Builder(builder: (_scaffoldInnerContext) => CupertinoButton(child: Text("Authenticate"), onPressed: () {print('initEthClientWithCanGetPrivateKeyHandler Authenticate onPressed current account encryptedAddress ${widget.encryptedAddress}'); _submit(s.data, widget.encryptedAddress, _scaffoldInnerContext);})), // Biulder Context required // ════════ Exception caught by gesture ═══════════════════════════════════════════ // Scaffold.of() called with a context that does not contain a Scaffold.
        ],),)
      )
    ))).then((value) => value as String);
  }

  @override
  Widget build(BuildContext context) {
    final Widget Function(Widget child) _centralizedPageBoardWallTemplate = (Widget child) {
      return Builder(builder: (_ctx) => DecoratedBox(
        decoration: BoxDecoration(
          image: DecorationImage(image: AssetImage('assets/images/bgCentralized.png', package: 'eurus'), fit: BoxFit.cover, alignment: Alignment.topCenter),
        ),
        child: Padding(
          padding: EdgeInsets.only(top: MediaQuery.of(_ctx).padding.top > 0 ? MediaQuery.of(_ctx).padding.top : 80), 
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: <Widget>[
              Expanded(
                child: DecoratedBox(
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadiusDirectional.only(topStart: Radius.circular(30), topEnd: Radius.circular(30)),
                  ), 
                  child: Padding(padding: EdgeInsets.only(top: 36).add(EdgeInsets.symmetric(horizontal: 16)), child: child),
                )
              )
            ]
          )
        )
      ));
    };

    final Widget Function(Widget child) _decentralizedPageBoardWallTemplate = (Widget child) {
      return Builder(builder: (_ctx) => DecoratedBox(
        decoration: BoxDecoration(
          image: DecorationImage(image: AssetImage('assets/images/bgDecentralized.png', package: 'eurus'), fit: BoxFit.cover, alignment: Alignment.topCenter),
        ),
        child: Padding(
          padding: EdgeInsets.only(top: MediaQuery.of(_ctx).padding.top > 0 ? MediaQuery.of(_ctx).padding.top : 80), 
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: <Widget>[
              Expanded(
                child: DecoratedBox(
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadiusDirectional.only(topStart: Radius.circular(30), topEnd: Radius.circular(30)),
                  ), 
                  child: Padding(padding: EdgeInsets.only(top: 36).add(EdgeInsets.symmetric(horizontal: 16)), child: child),
                )
              )
            ]
          )
        )
      ));
    };

    final Widget Function(BuildContext) Function(Widget appBarTitle, Widget Function(BuildContext) bodyBuilder) _importDecentralizedWalletPageBuilder = (appBarTitle, bodyBuilder) => (context) {
      return Scaffold(
        appBar: AppBar(title: appBarTitle, backgroundColor: Colors.transparent, elevation: 0,centerTitle:true),
        extendBodyBehindAppBar: true,
        backgroundColor: Colors.transparent,
        body: Builder(builder: (_) => _decentralizedPageBoardWallTemplate(bodyBuilder(_))),
        resizeToAvoidBottomInset: true
      );
    };

    final _defaultTextFieldInputDecoration = InputDecoration(
      filled: true, 
      fillColor: Color.fromRGBO(242, 244, 247, .75), //(0xFFf4f2f7),
      hintStyle: Theme.of(context).textTheme.subtitle1.apply(color: Theme.of(context).hintColor),
      border: OutlineInputBorder(borderSide: BorderSide.none, borderRadius: BorderRadius.circular(15)),
      contentPadding: EdgeInsets.all(16),
    );

    final _cryptoAssetListItemCardTileWidget = (data, {onTap}) => Builder(builder: (_) => Container(child: ListTile(
      leading: FutureBuilder(future: widget.canGetCryptoIconHandler(data['symbol'], 40), builder: (_, s) => s.data??Icon(Icons.circle, size: 40)), 
      title: Text(data['symbol']), 
      subtitle: Builder(builder: (__) {bool _showAddress = false; return StatefulBuilder(builder: (statefulBuilderCtx, _setState) => GestureDetector(child: Text(data['currency'] + (_showAddress ? "\n${data['address']}" : "")), onLongPressStart: (d) => _setState(() => _showAddress = !_showAddress), onLongPressEnd: (d) => _setState(() => _showAddress = !_showAddress),));}), 
      trailing: data['address'] is String && (data['address'] as String).isNotEmpty ? showBalance ? FutureBuilder(future: widget.canGetERC20BalanceHandler(widget.address, data['address'??'address${TEST_NET??'Ethereum'??'Eurus'}'], data['address'] == data['address${TEST_NET??'Ethereum'}'] ? 0 : 1), builder: (_, s) => Text("${s.connectionState != ConnectionState.done ? 'xxxx.xxxxxx' : !s.hasData ? '' : s.data}")) : Icon(Icons.remove_red_eye_outlined, color: Color(0xFF2684ff)) : Text('Coming Soon'),
      onTap: onTap,
    ), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, gradient: data['address'] == data['address${TEST_NET??'Ethereum'}'] ? null : LinearGradient(colors: [Colors.white, Color.alphaBlend(Colors.white70.withOpacity(.99), Color(0x9F4a00dd)), Color.alphaBlend(Colors.white70.withOpacity(.8), Color(0x9F4a00dd)), Color.alphaBlend(Colors.white70.withOpacity(.8), Color(0x9F0036dd)), Color.alphaBlend(Colors.white70.withOpacity(.99), Color(0x9F0036dd)), Colors.white]), borderRadius: BorderRadius.circular(15),)),);

    final _navigateToAssetDetails = (cryptoCurrencyModelMap) {
      final kTxType = {
        's': 'txSend', // input recepient address
        'r': 'txReceive', // from address
        'a': 'txAssetAllocation',
      };
      final kTxPeerAddress = 'txPeerAddress';
      final kTxHash = 'txHash';
      final kTxSmartContractAddress = 'txSmartContractAddress';
      final _txHis = [
        {'from': '0x1234abc', 'decodedInputAmount': 9876.54321, 'to': 'smart_contract_address', 'peerAddress': '0xrecepientAddress', 'transactionHash': '0xhash', 'transactionDateTime': '12/02/2021 12:34:56'},
        {'from': widget.address, 'decodedInputAmount': 9876.54321, 'to': 'smart_contract_address', 'peerAddress': '0xrecepientAddress', 'transactionHash': '0xhash', 'transactionDateTime': '02/02/2020 12:34:56'},
      ];

      final _cArgs = {
        '${cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'] ? "ethereum" : "eurus"}Erc20ContractAddress': cryptoCurrencyModelMap['address']??'0x0'??''??null, 
        'fromBlockChainType': cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'] ? 0 : 1,
        // 'currencyName': cryptoCurrencyModelMap['symbol']??'ETH'??'EUN',
        'disableSelectBlockchain': true,
        'ethereumAddress': widget.address,
        // 'canGetPrivateKeyHandler': _getPrivateKeyHandler,
        // 'navigateToAssetAllocationPage': () {
        //   Navigator.of(context).push
        // },
        
        // for contextual navigated-scene-aware asset allocation button in function bar item to be shown or not
        '${cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'] ? 'eurus' : 'ethereum'}Erc20ContractAddress': cryptoCurrencyModelMap['address${cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'] ? 'Eurus' : TEST_NET??'Ethereum'}'], 
      };
      final _assetDetailsPageBuilder = (assetPageBuilderContext) => Scaffold(
        appBar: AppBar(centerTitle:true,brightness: Brightness.light, foregroundColor: Colors.black, iconTheme: IconThemeData(color: Colors.black), backgroundColor: Colors.transparent, elevation: 0, title: Text('${cryptoCurrencyModelMap['symbol']}', style: Theme.of(assetPageBuilderContext).textTheme.headline6)),
        body: Builder(builder: (scaffoldBuilderContext) => Center(child: Column(children: [
          // ListTile(title: Card(child: Text(),))
          // Container(child: ListTile(leading: Icon(Icons.circle, size: 40), title: Text('USDT'), subtitle: Text('Tether USD'), trailing: Text("9876.54321"),), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),
          _cryptoAssetListItemCardTileWidget(cryptoCurrencyModelMap),
          Row(children: widget.homeFunctionsBarItemBuilder(scaffoldBuilderContext, {..._cArgs, ...{'navigateToAssetAllocationPage': (Future<Null> Function(BuildContext, Map<String, dynamic>) __navigateToAssetAllocationTransfer, [_barBuilderContext]) async => await __navigateToAssetAllocationTransfer(_barBuilderContext, {'fromBlockChainType': cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'] ? 0 : 1, 'ethereumErc20ContractAddress': cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'], 'eurusErc20ContractAddress': cryptoCurrencyModelMap['addressEurus'],})}}), mainAxisAlignment: MainAxisAlignment.spaceEvenly, crossAxisAlignment: CrossAxisAlignment.start,),
          Expanded(child: FutureBuilder(
            // future: Future.delayed(Duration(milliseconds: 500), () => _txHis),
            future: Future.delayed(Duration(milliseconds: 100), () async => await widget.canFetchTransactionHistoryHandler(widget.address, cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'] ? 0 : 1, cryptoCurrencyModelMap['address'])),
            // builder: (_, AsyncSnapshot<List<Map<String, Object>>> s) => s.connectionState != ConnectionState.done ? Center(child: CircularProgressIndicator()) : ListView.builder(itemBuilder: (itemBuilderContext, i) => ListTile(leading: Icon(s.data[i]['from'] == widget.address ? Icons.arrow_upward : Icons.arrow_downward), title: Text("${s.data[i]['peerAddress']}"), subtitle: Text.rich(TextSpan(text: "${s.data[i]['transactionDateTime']}\nTx ID: ${s.data[i]['transactionHash']}"), textWidthBasis: TextWidthBasis.longestLine, textScaleFactor: .85,), trailing: Text("${s.data[i]['from'] == widget.address ? '-' : ''} ${s.data[i]['decodedInputAmount']}", style: TextStyle(color: s.data[i]['from'] != widget.address ? Colors.green : Colors.black87),)), itemCount: s.data.length),)),
            builder: (_, AsyncSnapshot<List<Map<String, Object>>> s) => s.connectionState != ConnectionState.done 
              ? Center(child: CircularProgressIndicator.adaptive()) 
              : s.hasError 
                ? Center(child: Text("Ooops Something went wrong")) 
                : !s.hasData || s.data is List && s.data.length == 0 
                  ? cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'] 
                    ? Column(children: [Padding(padding: EdgeInsets.all(16), child: Text("\nAsset Allocation could help you Deposit To Eurus.", textAlign: TextAlign.center, style: Theme.of(scaffoldBuilderContext).textTheme.caption,)), Expanded(child: Center(child: Padding(padding: EdgeInsets.symmetric(horizontal: 16), child: Text("Enjoy benefits of Eurus for Attractively Low Transaction Fee in every Eurus Transaction from time to time.", textAlign: TextAlign.center, softWrap: true, style: Theme.of(scaffoldBuilderContext).textTheme.subtitle2))))]) 
                    : Align(alignment: Alignment.topCenter, child: Column(children: [
                      Padding(padding: EdgeInsets.symmetric(horizontal: 16).add(EdgeInsets.only(top: 16)), child: Text.rich(TextSpan(text: "\nIt is Perfectly Clean Sheet.\n\n\n\n", children: [TextSpan(text: "You're ready to benefit from the offer of Eurus for Low Transaction Fee.", style: Theme.of(scaffoldBuilderContext).textTheme.bodyText2), TextSpan(text: "\n\n\nTop-up your balance using Asset Allocation", style: Theme.of(scaffoldBuilderContext).textTheme.bodyText1.apply(color: Colors.black26)), ]), textAlign: TextAlign.center, style: Theme.of(scaffoldBuilderContext).textTheme.caption,)), 
                      CupertinoButton(padding: EdgeInsets.zero, child: Text("Deposit To Eurus", style: Theme.of(scaffoldBuilderContext).textTheme.button.apply(color: Color(0xFF4a00dd)??Theme.of(scaffoldBuilderContext).colorScheme.primary.withOpacity(.7)), ), onPressed: () {}), 
                    ])) 
                  : ListView.builder(itemBuilder: (itemBuilderContext, i) => ListTile(
                    onTap: () async {await widget.canNavigateToTransactionHistoryDetailPage(
                      // if not withdrawal but normal transfer
                      // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
                      centerMessage: '' ?? (s.data[i]['confirmTimestamp'] != null && s.data[i]['confirmTimestamp'] is String && s.data[i]['confirmTimestamp'] != '' ? 'Success' : 'Pending'),
                      // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
                      date: "${s.data[i]['confirmTimestamp'] != null && s.data[i]['confirmTimestamp'] is String && s.data[i]['confirmTimestamp'] != '' ? DateFormat.yMd('en_HK').add_jms().format(DateTime.fromMillisecondsSinceEpoch(int.tryParse(s.data[i]['confirmTimestamp'])))+' UTC+8' : 'Pending Transaction'}",
                      fromAddress: "${s.data[i]['txFrom']}",
                      // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                      // gasFeeString: "... ${cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'] ? 'ETH' : 'EUN'} Gas Fee",
                      gasFeeString: "${s.data[i]['txFrom']}".toLowerCase() == widget.address.toLowerCase() ? "... ${s.data[i]['chain'] == 'eth' ? 'ETH' : 'EUN'} Gas Fee" : '',
                      toAddress: "${s.data[i]['decodedInputRecipientAddress']}",
                      // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                      transferAmount: "${"${s.data[i]['txFrom']}".toLowerCase() == widget.address.toLowerCase() ? '-' : ''} ${s.data[i]['decodedInputAmount']} ${cryptoCurrencyModelMap['symbol']}",
                      txId: "${s.data[i]['transactionHash']}",
                      navigatorContext: itemBuilderContext,
                      // TODO: if not withdrawal but normal transfer
                      isAssetAllocation: null??false, 
                      // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                      // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
                      shouldSkipPendingFetch: "${s.data[i]['txFrom']}".toLowerCase() == widget.address.toLowerCase() ? s.data[i]['confirmTimestamp'] != null && s.data[i]['confirmTimestamp'] is String && s.data[i]['confirmTimestamp'] != '' : true, 
                      blockChainType: s.data[i]['chain'] == 'eth' ? 0 : 1,
                    );},
                    // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                    leading: Icon("${s.data[i]['txFrom']}".toLowerCase() == widget.address.toLowerCase() ? Icons.arrow_upward : Icons.arrow_downward), 
                    title: Text("${s.data[i]['txFrom']}".toLowerCase() == widget.address.toLowerCase() ? "${s.data[i]['decodedInputRecipientAddress']}" : "${s.data[i]['txFrom']}"), 
                    // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
                    subtitle: Text.rich(TextSpan(text: "${s.data[i]['confirmTimestamp'] != null && s.data[i]['confirmTimestamp'] is String && s.data[i]['confirmTimestamp'] != '' ? DateFormat.yMd('en_HK').add_jms().format(DateTime.fromMillisecondsSinceEpoch(int.tryParse(s.data[i]['confirmTimestamp'])))+' UTC+8' : 'Pending Transaction'}\nTx ID: ${s.data[i]['transactionHash']}"), textWidthBasis: TextWidthBasis.longestLine, textScaleFactor: .85,), 
                    // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                    trailing: Text("${"${s.data[i]['txFrom']}".toLowerCase() == widget.address.toLowerCase() ? '-' : ''} ${s.data[i]['decodedInputAmount']}", style: TextStyle(color: "${s.data[i]['txFrom']}".toLowerCase() == widget.address.toLowerCase() ? Colors.black87 : Colors.green),)
                  ), itemCount: s.data?.length),)),
        ],))),
      );
      
      Navigator.of(context).push(MaterialPageRoute(builder: _assetDetailsPageBuilder));
    };

    final _cryptoAssetListWidget = (AsyncSnapshot<String> s, String addressSuffix, {additionalListViewBottomPaddingContentInset, Future Function(dynamic) onTapHandler, bool Function(dynamic) additionalWhereClauseFilterTest}) => Builder(builder: (_) {final data = (jsonDecode(s.data) as List).where((e) => e['showAssets'] && e['address${addressSuffix}'] != null).where(additionalWhereClauseFilterTest??(e)=>true).map((e) {e['address'] = e['address${addressSuffix}']??'0xb9f5377a3f1c4748c71bd4ce2a16f13eecc4fec4'; return e;}).toList(); return ListView.builder(
      padding: EdgeInsets.only(top: 57, left: 8, right: 8, bottom: additionalListViewBottomPaddingContentInset + MediaQuery.of(_).padding.bottom), // padding bottom need to consider SafeArea (i.e. 40 + MediaQuery.of(context).padding.bottom)
      itemBuilder: (itemBuilderContext, i) => _cryptoAssetListItemCardTileWidget(data[i], onTap: () async {await onTapHandler(data[i]);}),
      itemCount: data.length,
    );});

    final _cArgsMaps = [
      {}, 
      {'walletAccountEncryptedAddress': widget.encryptedAddress},
      {},
    ];
    // final _navigationBarOnTapTabPageWidgets = widget.navigationBarOnTapTabBuilders.map((e) => Builder(builder: e({}))).toList();
    int _viewSelectedChainCode = TEST_NET!=null&&TEST_NET.isNotEmpty?-1:0; // default is Testnet (-1) or Ethereum (0)
    
    final _navigateToAssetAllocationTokenList = (Future<Null> Function(BuildContext, Map<String, dynamic>) __navigateToAssetAllocationTransfer, [_barBuilderContext]) async {
      final ___navigateToAssetAllocationTransfer = (_navigationContext, data) => __navigateToAssetAllocationTransfer(_navigationContext, {'fromBlockChainType': data['address'] == data['address${TEST_NET??'Ethereum'}'] ? 0 : 1, 'ethereumErc20ContractAddress': data['address${TEST_NET??'Ethereum'}'], 'eurusErc20ContractAddress': data['addressEurus'],});

      final _navigateToAssetAllocationTransfer = (data) async {
        final _assetAllocationTransferPageBuilder = (ctx) => Scaffold(body: Center(child: Container()));

        Navigator.of(context).push(MaterialPageRoute(builder: _assetAllocationTransferPageBuilder));
      };
      
      final _assetAllocationTokenListPageBuilder = (assetAllocationPageBuilderContext) => Builder(builder: (context) {
        int _viewSelectedChainCode = TEST_NET!=null&&TEST_NET.isNotEmpty?-1:0; // default is Testnet (-1) or Ethereum (0)
        return 
          Scaffold(
            appBar: AppBar(title: Text("Asset Allocation"), backgroundColor: Colors.transparent, elevation: 0,centerTitle:true),
            extendBodyBehindAppBar: true,
            backgroundColor: Colors.transparent,
            resizeToAvoidBottomInset: true,
            body: _decentralizedPageBoardWallTemplate(StatefulBuilder(builder: (statefulBuilderContext, _setState) => Padding(padding: EdgeInsets.only(top: 0), child: Column(children: [
              ListTile(title: Text("${_viewSelectedChainCode <= 0 ? 'Deposit' : 'Withdrawal'}", textAlign: TextAlign.center, style: Theme.of(context).textTheme.headline1.apply(fontSizeFactor: .5, color: Color(_viewSelectedChainCode <= 0 ? 0x9F4a00dd : 0x9F0036dd)),), dense: true,),
              Expanded(child: Padding(padding: EdgeInsets.only(top: 10 + 16.0), child: Container(
                clipBehavior: Clip.hardEdge,
                decoration: BoxDecoration(
                  borderRadius: BorderRadiusDirectional.only(topStart: Radius.circular(30), topEnd: Radius.circular(30)),
                  color: Colors.white,
                  boxShadow: kElevationToShadow[3]
                ),
                child: Padding(
                  padding: EdgeInsets.only(top: 0, left: 0, right: 0), 
                  child: Column(children: [
                    Expanded(child: Stack(children: [
                      // asset list section content body
                      // SafeArea(child: 
                      FutureBuilder(future: Future.delayed(Duration(milliseconds: 0), () async => await widget.canPersistNormallyHandler('assetsList_${widget.encryptedAddress}') as String), builder: (_, AsyncSnapshot<String> s) => s.connectionState != ConnectionState.done ? Center(child: CircularProgressIndicator()) : !s.hasData || !(s.data is String) || s.data is String && s.data.isEmpty ? Builder(builder: (_) {print('${widget.encryptedAddress} s.data ${s.data}'); return ListView(
                        padding: EdgeInsets.only(top: 57, left: 8, right: 8, bottom: 40),
                        children: [Container(child: ListTile(title: CupertinoButton(onPressed: () {Navigator.of(_).push(MaterialPageRoute(builder: widget.getEditAssetsListWidgetBuilder(widget.encryptedAddress),));}, child: Text("Add Token", style: Theme.of(_).textTheme.button.apply(color: Theme.of(_).colorScheme.primary)))), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),],
                      // );}) : Builder(builder: (_) {final data = (jsonDecode(s.data) as List).where((e) => e['showAssets'] && e['address${_viewSelectedChainCode <= 0 ? TEST_NET??'Ethereum' : 'Eurus'}'] != null).map((e) {e['address'] = e['address${_viewSelectedChainCode <= 0 ? TEST_NET??'Ethereum' : 'Eurus'}']??'0xb9f5377a3f1c4748c71bd4ce2a16f13eecc4fec4'; return e;}).toList(); return ListView.builder(
                      );}) : _cryptoAssetListWidget(s, _viewSelectedChainCode <= 0 ? TEST_NET??'Ethereum' : 'Eurus', additionalListViewBottomPaddingContentInset: 0, onTapHandler: (_data) async => await ___navigateToAssetAllocationTransfer(statefulBuilderContext, _data), additionalWhereClauseFilterTest: (e) => e['address${TEST_NET??'Ethereum'}'] is String && (e['address${TEST_NET??'Ethereum'}'] as String).isNotEmpty && e['addressEurus'] is String && (e['addressEurus'] as String).isNotEmpty, )??Builder(builder: (_) {final data = (jsonDecode(s.data) as List).where((e) => e['showAssets'] && e['address${_viewSelectedChainCode <= 0 ? TEST_NET??'Ethereum' : 'Eurus'}'] != null).map((e) {e['address'] = e['address${_viewSelectedChainCode <= 0 ? TEST_NET??'Ethereum' : 'Eurus'}']??'0xb9f5377a3f1c4748c71bd4ce2a16f13eecc4fec4'; return e;}).toList(); return ListView.builder(
                        padding: EdgeInsets.only(top: 57, left: 8, right: 8, bottom: 40 + MediaQuery.of(context).padding.bottom), // padding bottom need to consider SafeArea (i.e. 40 + MediaQuery.of(context).padding.bottom)
                        // itemBuilder: (itemBuilderContext, i) => _cryptoAssetListItemCardTileWidget(data[i], onTap: () {_navigateToAssetAllocationTransfer(data[i]);}),
                        // itemBuilder: (itemBuilderContext, i) => _cryptoAssetListItemCardTileWidget(data[i], onTap: () {__navigateToAssetAllocationTransfer(statefulBuilderContext, {'fromBlockChainType': data[i]['address'] == data[i]['address${TEST_NET??'Ethereum'}'] ? 0 : 1, 'ethereumErc20ContractAddress': data[i]['address${TEST_NET??'Ethereum'}'], 'eurusErc20ContractAddress': data[i]['addressEurus'],});}),
                        itemBuilder: (itemBuilderContext, i) => _cryptoAssetListItemCardTileWidget(data[i], onTap: () {___navigateToAssetAllocationTransfer(statefulBuilderContext, data[i]);}),
                        itemCount: data.length,
                      );})),
                      // bottom: false ,),
                      // asset list section header
                      Align(alignment: Alignment.topCenter, child: Container(decoration: BoxDecoration(gradient: LinearGradient(begin: Alignment.center, end: Alignment.bottomCenter, colors: [Colors.white, Color(0x00FFFFFF)], ), ), child:
                        Padding(padding: EdgeInsets.symmetric(horizontal: 8), child: IntrinsicHeight(child: Row(children: [
                          CupertinoButton(onPressed: () {_setState(()=>_viewSelectedChainCode=TEST_NET!=null&&TEST_NET.isNotEmpty?-1:0);}, child: IntrinsicWidth(child: Column(children: [Text("Ethereum", style: Theme.of(context).textTheme.headline6.apply(color: Color(0xFF002251), fontWeightDelta: 1,)), FractionallySizedBox(widthFactor: .5, child: Container(child: Divider(color: _viewSelectedChainCode <= 0 ? Color(0xFF0041c4) : Colors.transparent, height: 0, thickness: 1.5, indent: 0, endIndent: 0)))]))),
                          VerticalDivider(width: 0, thickness: 1, color: Colors.black12, indent: 15, endIndent: 15,),
                          CupertinoButton(onPressed: () {_setState(()=>_viewSelectedChainCode=1);}, child: IntrinsicWidth(child: Column(children: [Text("Eurus", style: Theme.of(context).textTheme.headline6.apply(color: Color(0xFF002251), fontWeightDelta: 1,)), FractionallySizedBox(widthFactor: .5, child: Container(child: Divider(color: _viewSelectedChainCode == 1 ? Color(0xFF0041c4) : Colors.transparent, height: 0, thickness: 1.5, indent: 0, endIndent: 0)))]))),
                          Spacer(flex:1),
                        ]))),
                      ),)
                    ]),)
                  ])
                )
              )))
            ]))))
          )
        ;
      });

      Navigator.of(context).push(MaterialPageRoute(builder: _assetAllocationTokenListPageBuilder));
    };
    final _cArgs = {
      'ethereumAddress': widget.address,
      'navigateToAssetAllocationPage': _navigateToAssetAllocationTokenList,
      // if ( (() => widget.createdAccountGeneratedMnemonicSeedPhrase)() is String && (() => widget.createdAccountGeneratedMnemonicSeedPhrase)().isNotEmpty ) 'replacingQRCodeWidget': Builder(builder: (context) => Column(children: [
      //   Text("You are recommended to backup your wallet mnemonic seed phrase before receiving any fund into this wallet from other external funding sources"),
      //   CupertinoButton(child: Text("Backup now")),
      // ],),)
    };

    // final _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder = (_backupWalletSecretMnemonicSeedPhrasesWidgetBuilderContext) => Builder(builder: (context) => Padding(padding: EdgeInsets.all(16), child: Center(child: Container(
    // final _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder = ([State<EurusWalletHomePage> _parent]) => (_backupWalletSecretMnemonicSeedPhrasesWidgetBuilderContext) => Builder(builder: (context) => Padding(padding: EdgeInsets.all(16), child: Center(child: Container(
    final _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder = (Function _parentNotifyOnValueChange) => (_backupWalletSecretMnemonicSeedPhrasesWidgetBuilderContext) => Builder(builder: (context) => Padding(padding: EdgeInsets.all(16), child: Center(child: Container(
      decoration: BoxDecoration(borderRadius: BorderRadius.circular(30), color: Colors.white),
      padding: EdgeInsets.symmetric(horizontal: 8, vertical: 26),
      child: Column(children: [
        ListTile(title: Text("Mnemonic Seed Phrase Backup is Required", textAlign: TextAlign.center, style: Theme.of(context).textTheme.headline6.copyWith(letterSpacing: 2, fontWeight: FontWeight.w800)),),
        // Icon(Icons.info_outline_rounded, size: 100,),
        Expanded(child: Center(child: ListTile(title: Icon(Icons.info_outline_rounded, size: 120,), subtitle: Text("\n\nYou are recommended to backup mnemonic seed phrase of any newly created wallet before receiving any fund into this wallet from other external funding sources", textAlign: TextAlign.justify,)))),
        CupertinoButton(child: Text("Backup now"), onPressed: () async {
          final _popResult = await Navigator.of(context).push(MaterialPageRoute(fullscreenDialog: true, builder: (_materialPageRouteBuilderCtx) => Scaffold(
            appBar: AppBar(title: Text("Backup Wallet Secret Seed"), backgroundColor: Colors.transparent, elevation: 0, centerTitle:true),
            extendBodyBehindAppBar: true,
            backgroundColor: Colors.transparent,
            resizeToAvoidBottomInset: true,
            body: _decentralizedPageBoardWallTemplate(Builder(builder: (_scaffoldBodyTemplateChildBuilderCtx) {
              final _textEditingController0 = TextEditingController();
              final _typeToConfirmReveal = 'Wallet Secret';
              
              // Auth UI
              final Future<String> Function(String) _authUI = (String uriQuery) async {
                final _textEditingController = TextEditingController();
                final _submit = (_uriQuery, _accountEncryptedAddress, _ctx) {
                  print('Mnemonic Seed Phrase Backup _submit _accountEncryptedAddress $_accountEncryptedAddress');
                  final _encryptedValue = Uri(query: _uriQuery).queryParameters[_accountEncryptedAddress];
                  final _decryptedValue = widget.canDecryptHandler(_textEditingController.text, _encryptedValue);
                  if (_decryptedValue != null) return Navigator.of(_ctx).maybePop(_decryptedValue);
                  // else, auth failure should be prompted
                  Scaffold.of(_ctx).showSnackBar(SnackBar(content: Text("Authentication Failure", textAlign: TextAlign.center,), backgroundColor: Colors.redAccent.shade100.withOpacity(.9),));
                };
                
                return await Navigator.of(context).push(PageRouteBuilder(fullscreenDialog: true, opaque: false, pageBuilder: (pageBuilderContext, animation, secondaryAnimation) => Scaffold(
                  backgroundColor: Colors.black87,
                  body: Builder(builder: (_scaffoldInnerContext) {
                    final __submit = () => _submit(uriQuery, widget.encryptedAddress, _scaffoldInnerContext);
                    // final _tryBioAuth = () => Future.delayed(Duration.zero, () async => await widget.canPersistWithBiometricSecurelyHandler(__prefix+widget.encryptedAddress)).then((value) {if (value != null) {_textEditingController.text = value; __submit();}});
                    // _tryBioAuth();
                    return 
                      KeyboardAvoider(autoScroll: true, child: Column(mainAxisAlignment: MainAxisAlignment.center, crossAxisAlignment: CrossAxisAlignment.stretch, children: [Container(decoration: BoxDecoration(borderRadius: BorderRadius.circular(15), color: Colors.white.withOpacity(.9)), padding: EdgeInsets.all(16), child: Column(children: [
                        ListTile(title: Text("Security Authentication", style: Theme.of(_scaffoldInnerContext).textTheme.headline6.apply(fontWeightDelta: 2), textAlign: TextAlign.center)),
                        ListTile(title: TextField(textInputAction: TextInputAction.go, onSubmitted: (_) => __submit(), autofocus: true, obscureText: true, controller: _textEditingController, decoration: InputDecoration(border: OutlineInputBorder(borderRadius: BorderRadius.circular(16)), suffixIcon: null))),
                        ListTile(title: TextButton(child: Text("Submit", style: Theme.of(_scaffoldInnerContext).textTheme.caption.apply(color: Colors.blue)), onPressed: () {
                          __submit();
                        },), subtitle: TextButton(child: Text("Cancel"), onPressed: () => Navigator.of(_scaffoldInnerContext).maybePop()))
                      ],),)]));
                  })
                ))).then((value) => value as String);
              };
              
              final _confirmSeedPhraseWidgetBuilder = (String _mSeedPhrase) => (_materialPageRouteBuilderCtx) => Scaffold(
                appBar: AppBar(title: Text("Backup Wallet Secret Seed"), backgroundColor: Colors.transparent, elevation: 0, centerTitle:true),
                extendBodyBehindAppBar: true,
                backgroundColor: Colors.transparent,
                resizeToAvoidBottomInset: true,
                body: _decentralizedPageBoardWallTemplate(Builder(builder: (_scaffoldBodyTemplateChildBuilderCtx) {
                  final randomlyOrderedPhrases = _mSeedPhrase.split(' ')..shuffle();
                  // List<String> inputPhrases = List<String>(12);
                  List<int> inputPhrasesIndices = List<int>(12);
                  int explicitWordCursor = -1;
                  FocusNode focusNode = FocusNode();
                  // Future.delayed(Duration.zero, () {focusNode.requestFocus();});
                  return StatefulBuilder(builder: (_statefulBuilderCtx, __setState) {
                    Future.delayed(Duration.zero, () => focusNode.requestFocus());
                    return Padding(padding: EdgeInsets.only(bottom: MediaQuery.of(_statefulBuilderCtx).padding.bottom), child: Column(children: [
                      Center(child: ListTile(dense: true, title: Text('Confirm seed phrase', textAlign: TextAlign.center, style: Theme.of(_statefulBuilderCtx).textTheme.headline5.copyWith(color: Color(0xFF002251), fontSize: 24, fontFamily: 'packages/eurus/SFProDisplay', fontWeight: FontWeight.w600), ))),
                      Padding(padding:EdgeInsets.only(bottom: 20),  child:Center(child:Text('Please enter mnemonic phrases in correct order for confirmation.', textAlign: TextAlign.left, style: Theme.of(_statefulBuilderCtx).textTheme.headline5.copyWith(color: Color(0xFFc1c7d0), fontSize: 14, fontFamily: 'packages/eurus/SFProDisplay', fontWeight: FontWeight.normal)))),
                      // ════════ Exception caught by rendering library ═════════════════════════════════
                      // An InputDecorator, which is typically created by a TextField, cannot have an unbounded width.
                      // This happens when the parent widget does not provide a finite width constraint. For example, if the InputDecorator is contained by a Row, then its width must be constrained. An Expanded widget or a SizedBox can be used to constrain the width of the InputDecorator or the TextField that contains it.
                      // 'package:flutter/src/material/input_decorator.dart':
                      // Failed assertion: line 948 pos 7: 'layoutConstraints.maxWidth < double.infinity'
                      // ...List.generate(4, (y) => Row(mainAxisAlignment: MainAxisAlignment.spaceAround, children: List.generate(3, (x) => Expanded(child: Padding(padding: EdgeInsets.all(2), child: TextField(onTap: () {inputPhrases[y*3+x] = ''; __setState(() => explicitWordCursor = inputPhrases.indexWhere((e) => e == null || e.isEmpty) == y*3+x ? -1 : y*3+x);}, focusNode: explicitWordCursor > -1 && explicitWordCursor == y*3+x || explicitWordCursor == -1 && inputPhrases.indexWhere((e) => e == null || e.isEmpty) == y*3+x ? focusNode : null, readOnly: true, enableInteractiveSelection: false, decoration: InputDecoration(border: OutlineInputBorder(borderRadius: BorderRadius.circular(15),), focusedBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(15), borderSide: BorderSide(color: Color(0xFF4406dd)))/*, fillColor: Color(0xFF4406dd), filled: false,*/, labelText: "${y*3+x+1}"), controller: TextEditingController(text: "${inputPhrases[y*3+x]??''}"),)))))),
                      // ...List.generate(4, (y) => Row(mainAxisAlignment: MainAxisAlignment.spaceAround, children: List.generate(3, (x) => OutlineButton(onPressed: () {inputPhrases[explicitWordCursor > -1 ? explicitWordCursor : inputPhrases.indexWhere((e) => e == null || e.isEmpty)] = randomlyOrderedPhrases[y*3+x]; __setState(() => explicitWordCursor = -1);}, color: Color(0xFF4406dd), child: Text("${randomlyOrderedPhrases[y*3+x]}"),)))),
                      // if (_mSeedPhrase == (_mSeedPhrase.split(' ')??inputPhrases).join(' ')) CupertinoButton(onPressed: () => Navigator.popUntil(_statefulBuilderCtx, ModalRoute.withName('HomePage')), child: Text('Finished')),
                      ...List.generate(4, (y) => Row(mainAxisAlignment: MainAxisAlignment.spaceAround, children: List.generate(3, (x) => Expanded(child: Padding(padding: EdgeInsets.all(2), child: TextField(onTap: () {inputPhrasesIndices[y*3+x] = -1; __setState(() => explicitWordCursor = inputPhrasesIndices.indexWhere((e) => e == null || e == -1) == y*3+x ? -1 : y*3+x);}, focusNode: explicitWordCursor > -1 && explicitWordCursor == y*3+x || explicitWordCursor == -1 && inputPhrasesIndices.indexWhere((e) => e == null || e == -1) == y*3+x ? focusNode : null, readOnly: true, enableInteractiveSelection: false, decoration: InputDecoration(border: OutlineInputBorder(borderRadius: BorderRadius.circular(15),), focusedBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(15), borderSide: BorderSide(color: Color(0xFF4406dd), width: 2))/*, fillColor: Color(0xFF4406dd), filled: false,*/, labelText: "${y*3+x+1}"), controller: TextEditingController(text: "${(inputPhrasesIndices[y*3+x]??-1) > -1 ? randomlyOrderedPhrases[inputPhrasesIndices[y*3+x]] : ''}"),)))))),
                      Spacer(flex: 1),
                      ...List.generate(4, (y) => Row(mainAxisAlignment: MainAxisAlignment.spaceAround, children: List.generate(3, (x) => Expanded(child: Padding(padding: EdgeInsets.all(2), child: CupertinoButton(onPressed: inputPhrasesIndices.indexOf(y*3+x) > -1 ? null : () {inputPhrasesIndices[explicitWordCursor > -1 ? explicitWordCursor : inputPhrasesIndices.indexWhere((e) => e == null || e == -1)] = y*3+x; __setState(() => explicitWordCursor = -1);}, color: Color(0xFF4406dd), padding: EdgeInsets.zero, child: Text("${randomlyOrderedPhrases[y*3+x]}", style: Theme.of(_statefulBuilderCtx).textTheme.button.apply(color: Colors.white)),)))))),
                      // CupertinoButton(onPressed: (inputPhrasesIndices.indexWhere((e) => e == null || e == -1) == -1 && _mSeedPhrase == (inputPhrasesIndices.map((e) => e == null || e == -1 ? '' : randomlyOrderedPhrases[e]).toList()).join(' ')) ? () {widget.canPersistNormallyHandler({'mPhraseBackuped': Uri(queryParameters: {widget.encryptedAddress: widget.createdAccountGeneratedMnemonicSeedPhrase}).query}); _setState(() => widget.createdAccountGeneratedMnemonicSeedPhrase = null); Navigator.popUntil(_statefulBuilderCtx, ModalRoute.withName('HomePage'));} : null, child: Text('Finished')),
                      // CupertinoButton(onPressed: (inputPhrasesIndices.indexWhere((e) => e == null || e == -1) == -1 && _mSeedPhrase == (inputPhrasesIndices.map((e) => e == null || e == -1 ? '' : randomlyOrderedPhrases[e]).toList()).join(' ')) ? () {widget.canPersistNormallyHandler({'mPhraseBackuped': Uri(queryParameters: {widget.encryptedAddress: widget.createdAccountGeneratedMnemonicSeedPhrase}).query}); _parent.setState(() => widget.createdAccountGeneratedMnemonicSeedPhrase = null); Navigator.popUntil(_statefulBuilderCtx, ModalRoute.withName('HomePage'));} : null, child: Text('Finished')),
                      CupertinoButton(onPressed: (inputPhrasesIndices.indexWhere((e) => e == null || e == -1) == -1 && _mSeedPhrase == (inputPhrasesIndices.map((e) => e == null || e == -1 ? '' : randomlyOrderedPhrases[e]).toList()).join(' ')) ? () {widget.canPersistNormallyHandler({'mPhraseBackuped': Uri(queryParameters: {widget.encryptedAddress: widget.createdAccountGeneratedMnemonicSeedPhrase}).query}); widget.createdAccountGeneratedMnemonicSeedPhrase = null; (_parentNotifyOnValueChange??() {})(); Navigator.popUntil(_statefulBuilderCtx, ModalRoute.withName('HomePage'));} : null, child: Text('Finished')),
                    ]));
                  });
                }))
              );

              final _navigateToMnemonicSeedPhraseConfirmation = (String _mSeedPhrase) => Navigator.of(context).push(MaterialPageRoute(builder: _confirmSeedPhraseWidgetBuilder(_mSeedPhrase)));

              // final _mSeedPhraseRevealerWidget = (String _mSeedPhrase, [bool _isRevealing = false]) => Column(children: [ListTile(title: Text("Please write down the 12 word seed phrases in a secured place"+"\n\n"+"To protect your wallet secret against accidental leakage, reveal mnemonic seed phrase only by long pressing and holding onto the grey area", style: Theme.of(context).textTheme.bodyText2.apply(color: Color(0xFFc1c7d0)), textAlign: TextAlign.justify,),), ListTile(title: StatefulBuilder(builder: (__statefulBuilderCtx, ___setState) => GestureDetector(child: SizedBox(height: .25*MediaQuery.of(__statefulBuilderCtx).size.height, child: TextField(obscureText: !_isRevealing, enabled: false, readOnly: true, decoration: _defaultTextFieldInputDecoration, maxLines: !_isRevealing ? 1 : null, expands: !_isRevealing ? false : true, textAlignVertical: TextAlignVertical.center, style: Theme.of(context).textTheme.headline5, controller: TextEditingController(text: _mSeedPhrase)))??Text("long-press and hold here"), onLongPressStart: (d) {___setState(() => _isRevealing = !_isRevealing);}, onLongPressEnd: (d) {___setState(() => _isRevealing = !_isRevealing);},)), subtitle: Text(/*""/*)), ListTile(title: Text(*/+"\n\n"+*/"WARNING: Please ensure no video recording, nor cctv surveillance around you. Please ensure there is no spy and watcher behind your screen. Please ensure there is no screen mirroring, remoting, sharing, nor projecting presentation of your screen.", textAlign: TextAlign.justify,),)]);
              final _mSeedPhraseRevealerWidget = (String _mSeedPhrase, [bool _isRevealing = false]) => Column(children: [ListTile(title: Text("Please write down the 12 word seed phrases in a secured place"+"\n\n"+"To protect your wallet secret against accidental leakage, reveal mnemonic seed phrase only by long pressing and holding onto the grey area", style: Theme.of(context).textTheme.bodyText2.apply(color: Color(0xFFc1c7d0)), textAlign: TextAlign.justify,),), ListTile(title: StatefulBuilder(builder: (__statefulBuilderCtx, ___setState) => GestureDetector(child: SizedBox(height: .25*MediaQuery.of(__statefulBuilderCtx).size.height, child: Stack(children: [Align(alignment: Alignment.bottomCenter, child: Icon(Icons.videocam_off, size: .15*MediaQuery.of(__statefulBuilderCtx).size.height, color: Colors.black26,)), TextField(obscureText: !_isRevealing, enabled: false, readOnly: true, decoration: _defaultTextFieldInputDecoration, maxLines: !_isRevealing ? 1 : null, expands: !_isRevealing ? false : true, textAlignVertical: TextAlignVertical.center, style: Theme.of(context).textTheme.headline5, controller: TextEditingController(text: _mSeedPhrase))])??Text("long-press and hold here")), onLongPressStart: (d) {___setState(() => _isRevealing = !_isRevealing);}, onLongPressEnd: (d) {___setState(() => _isRevealing = !_isRevealing);},)), subtitle: Text(/*""/*)), ListTile(title: Text(*/+"\n\n"+*/"WARNING: Please ensure no video recording, nor cctv surveillance around you. Please ensure there is no spy and watcher behind your screen. Please ensure there is no screen mirroring, remoting, sharing, nor projecting presentation of your screen.", textAlign: TextAlign.justify,),), ListTile(title: CupertinoButton(onPressed: () {_navigateToMnemonicSeedPhraseConfirmation(_mSeedPhrase);}, child: Text("Next")))]);

              return StatefulBuilder(builder: (_statefulBuilderCtx, __setState) => Column(children: [
                Center(child: ListTile(dense: true, title: Text('Mnemonic Seed Phrase', textAlign: TextAlign.center, style: Theme.of(_statefulBuilderCtx).textTheme.headline5.copyWith(color: Color(0xFF002251), fontSize: 24, fontFamily: 'packages/eurus/SFProDisplay', fontWeight: FontWeight.w600), ))),
                if (_textEditingController0.text != _typeToConfirmReveal) ListTile(
                  title: Padding(child: Text.rich(TextSpan(text: "This action will reveal Mnemonic Seed Phrase. Please type ", children: [TextSpan(text: "$_typeToConfirmReveal", style: TextStyle(fontWeight: FontWeight.w900, color: Colors.black38)), TextSpan(text: ' to confirm')]), style: Theme.of(_statefulBuilderCtx).textTheme.bodyText2.apply(color: Color(0xFFc1c7d0))), padding: EdgeInsets.symmetric(vertical: 24)), 
                  subtitle: TextFormField(
                    textAlignVertical: TextAlignVertical.center,
                    decoration: _defaultTextFieldInputDecoration.copyWith(
                      hintText: "",
                      hintStyle: TextStyle(color: Color.fromRGBO(122, 134, 154, .5)),
                      suffix: CupertinoButton(child: Text("Understood"), onPressed: () async {
                        
                        /**
                         * // Auth UI method 1
                         * // onPress confirm button to setState to re-render the statefulBuilder and let other siblings widget (adjacent children) FutureBuilder to handle subsequent flow where authenticated access verfification is neccessarily ensured
                         */
                        __setState(() {});
                        
                        /**
                         * // Auth UI method 2
                         */
                        // if (_textEditingController0.text == _typeToConfirmReveal && await _authUI(await widget.canPersistSecurelyHandler(__prefix+'createdAccountGeneratedMnemonicSeedPhrases').then((value) => value as String)) != null) __setState(() {});
                      },),
                    ),
                    controller: _textEditingController0,
                  )
                ),
                /**
                 * // Auth UI method 1
                 */
                // issue: did not handle CircularProgressIndicator keep showing if Auth UI Cancel
                // if (_textEditingController0.text == _typeToConfirmReveal) FutureBuilder(future: Future.delayed(Duration.zero, () async => await _authUI(await widget.canPersistSecurelyHandler(__prefix+'createdAccountGeneratedMnemonicSeedPhrases').then((value) => value as String))), builder: (_fBCtx, _s) => _s.connectionState == ConnectionState.done && !_s.hasError && _s.hasData ? _mSeedPhraseRevealerWidget(_s.data) : CircularProgressIndicator(),),
                
                // resolved: handle by forcefully setState _textEditingController0.text to empty string '' and rebuild 
                if (_textEditingController0.text == _typeToConfirmReveal) FutureBuilder(future: Future.delayed(Duration.zero, () async {final _authUIResult = await _authUI(await widget.canPersistSecurelyHandler(__prefix+'createdAccountGeneratedMnemonicSeedPhrases').then((value) => value as String)); if (_authUIResult != null) return _authUIResult; __setState(() => _textEditingController0.text = ''); }), builder: (_fBCtx, _s) => _s.connectionState == ConnectionState.done && !_s.hasError && _s.hasData ? _mSeedPhraseRevealerWidget(_s.data) : CircularProgressIndicator(),),
                
                /**
                 * // Auth UI method 2
                 */
                // if (_textEditingController0.text == _typeToConfirmReveal) _mSeedPhraseRevealerWidget(widget.mSeedPhrase),

                // depending on navigator of context, shall decide Exit and Close to where
                // 1: Positive Pop result matters
                // Container(child: TextButton(child: Text("Close"), onPressed: () => Navigator.of(_scaffoldBodyTemplateChildBuilderCtx).pop(true))),
                // 2: popUntil to Wallet HomePage while Pop result doesn't matter
                // Container(child: TextButton(child: Text("Close"), onPressed: () => Navigator.of(context).popUntil((route) => route.settings.name == 'HomePage'))),
                // 3: Negative Pop result matters
                // Container(child: TextButton(child: Text("It is not ready to backup secret now.\nBackup Later", textAlign: TextAlign.center,), onPressed: () => Navigator.of(_scaffoldBodyTemplateChildBuilderCtx).pop(false))),
                Container(child: TextButton(child: Text.rich(TextSpan(text: "Not ready yet now? ", children: [TextSpan(text: "Backup Later", style: TextStyle(color: Theme.of(_statefulBuilderCtx).colorScheme.primary))], style: TextStyle(color: Colors.black87)), textAlign: TextAlign.center,), onPressed: () => Navigator.of(_scaffoldBodyTemplateChildBuilderCtx).pop(false))),
              ]));
            }))
          )));
          if (_popResult is bool && _popResult) Navigator.of(context).popUntil((route) => route.settings.name == 'HomePage');
        },),
      ]),
    ))));

    final Widget Function(BuildContext) Function(List<Builder>) _decentralizedWalletHomePageBuilder = (_navigationBarOnTapTabPageWidgets) => (context) => StatefulBuilder(builder: (statefulBuilderContext, _setState) {
      
      return Stack(children: [navigationBarTabBuilderSelectedIndex != null ? /*Builder(builder: widget.navigationBarOnTapTabBuilders[navigationBarTabBuilderSelectedIndex])*/ _navigationBarOnTapTabPageWidgets[navigationBarTabBuilderSelectedIndex] : Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: false,
        brightness: Brightness.light, 
        backgroundColor: Colors.transparent, 
        elevation: 0, 
        shadowColor: Colors.transparent,
        centerTitle: true,
        // leading: Padding(padding: EdgeInsets.only(bottom: 15, left: 20), child: Align(alignment: Alignment.bottomLeft, child: Icon(Icons.notifications, color: Color(0xFF2684ff)))), 
        title: SizedBox(child: Column(children: [
          Image.asset('assets/images/imgLogoH.png', package: 'eurus'), 
          Padding(padding: EdgeInsets.only(left: 3.0*2), child: FlatButton(
            height: 0,
            clipBehavior: Clip.antiAlias,
            padding: EdgeInsets.zero,
            shape: StadiumBorder(),
            colorBrightness: Brightness.dark, 
            onPressed: () {
             // log("${widget.address}"); _setState(() {});
              },
            child: DecoratedBox(
              decoration: BoxDecoration(gradient: LinearGradient(colors: [Color(0xFF4a00dd), Color(0xFF0036dd)]),), 
              child: Padding(padding: EdgeInsets.symmetric(vertical: 1.5, horizontal: 10.5), child: Text.rich(TextSpan(text: "Decentralized", children: [TextSpan(text: " 􀄭", style: Theme.of(context).textTheme.bodyText2.merge(TextStyle(fontFamily: 'packages/eurus/SFCompact', color: Colors.white, fontWeight: FontWeight.w100)))]), ))
            ),
          ))
        ])), 
        toolbarHeight: 75.0+4.0,//2*kToolbarHeight,
      ),
      backgroundColor: Color(0xFFf7f9fb),
      body: Builder(builder: (context) => Padding(padding: EdgeInsets.only(top: 0), child: Column(children: [
        Image.asset('assets/images/banner1.png', package: 'eurus'),
        Row(children: widget.homeFunctionsBarItemBuilder(context, /*{'canGetPrivateKeyHandler': _getPrivateKeyHandler}*/{..._cArgs, ...{
          // 'ethereumErrorPopUp': true,
          // 'eurusErrorPopUp': true,
          // 'errorText': '',
          // if ( widget.createdAccountGeneratedMnemonicSeedPhrase is String && widget.createdAccountGeneratedMnemonicSeedPhrase != null && widget.createdAccountGeneratedMnemonicSeedPhrase.isNotEmpty ) 'replacingQRCodeWidget': Builder(builder: _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder),
          // if ( widget.createdAccountGeneratedMnemonicSeedPhrase is String && widget.createdAccountGeneratedMnemonicSeedPhrase != null && widget.createdAccountGeneratedMnemonicSeedPhrase.isNotEmpty ) 'replacingQRCodeWidget': Builder(builder: _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder(this)),
          if ( widget.createdAccountGeneratedMnemonicSeedPhrase is String && widget.createdAccountGeneratedMnemonicSeedPhrase != null && widget.createdAccountGeneratedMnemonicSeedPhrase.isNotEmpty ) 'replacingQRCodeWidget': Builder(builder: _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder(() => _setState(() {}))),
        }}), mainAxisAlignment: MainAxisAlignment.spaceEvenly, crossAxisAlignment: CrossAxisAlignment.start,),
        Expanded(child: Padding(padding: EdgeInsets.only(top: 10), child: Container(
          clipBehavior: Clip.hardEdge,
          decoration: BoxDecoration(
            borderRadius: BorderRadiusDirectional.only(topStart: Radius.circular(30), topEnd: Radius.circular(30)),
            color: Colors.white,
            boxShadow: kElevationToShadow[3]
          ),
          child: Padding(
            padding: EdgeInsets.only(top: 0, left: 0, right: 0), 
            child: Column(children: [
              Expanded(child: Stack(children: [
                // asset list section content body
                // SafeArea(child: 
                FutureBuilder(future: /* /* // work // */ widget.canPersistNormallyHandler('assetsList_${widget.encryptedAddress}').then((value) => value as String)*/ Future.delayed(Duration(milliseconds: 0), /* // wrong // */ /*() => widget.canPersistNormallyHandler('assetsList_${widget.encryptedAddress}') as Future<String>*/ /* /* // debug found wrong // */ () async {print('assetsList_${widget.encryptedAddress} : ${(await widget.canPersistNormallyHandler('assetsList_${widget.encryptedAddress}'))}'); return widget.canPersistNormallyHandler('assetsList_${widget.encryptedAddress}')  as Future<String>;}*/ () async => await widget.canPersistNormallyHandler('assetsList_${widget.encryptedAddress}') as String), builder: (_, AsyncSnapshot<String> s) => s.connectionState != ConnectionState.done ? Center(child: CircularProgressIndicator()) : !s.hasData || !(s.data is String) || s.data is String && s.data.isEmpty ? /*Center(child: TextButton(onPressed: () {}, child: Text("Add Token"))) : */ Builder(builder: (_) {print('${widget.encryptedAddress} s.data ${s.data}'); return ListView(
                  padding: EdgeInsets.only(top: 57, left: 8, right: 8, bottom: 40),
                  // children: [
                  //   Container(child: ListTile(leading: Icon(Icons.circle, size: 40), title: Text('USDT'), subtitle: Text('Tether USD'), trailing: Text("9876.54321"), onTap: () {_navigateToAssetDetails('USDT', '0xabcd1234');},), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),
                  //   Container(child: ListTile(leading: Icon(Icons.circle, size: 40), title: Text('USDT'), subtitle: Text('Tether USD'), trailing: Text("9876.54321"),), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),
                  //   Container(child: ListTile(leading: Icon(Icons.circle, size: 40), title: Text('USDT'), subtitle: Text('Tether USD'), trailing: Text("9876.54321"),), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),
                  //   Container(child: ListTile(leading: Icon(Icons.circle, size: 40), title: Text('USDT'), subtitle: Text('Tether USD'), trailing: Text("9876.54321"),), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),
                  //   Container(child: ListTile(leading: Icon(Icons.circle, size: 40), title: Text('USDT'), subtitle: Text('Tether USD'), trailing: Text("9876.54321"),), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),
                  //   Container(child: ListTile(leading: Icon(Icons.circle, size: 40), title: Text('USDT'), subtitle: Text('Tether USD'), trailing: Text("9876.54321"),), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),
                  // ],
                  // children: [Container(child: ListTile(leading: Icon(Icons.add_circle_rounded, size: 40), title: Text('Add Token'), onTap: () {_navigateToAssetDetails('USDT', '0xabcd1234');},), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),],
                  children: [Container(child: ListTile(title: CupertinoButton(onPressed: () {Navigator.of(_).push(MaterialPageRoute(builder: widget.getEditAssetsListWidgetBuilder(widget.encryptedAddress),));}, child: Text("Add Token", style: Theme.of(_).textTheme.button.apply(color: Theme.of(_).colorScheme.primary)))), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),],
                );}) : _cryptoAssetListWidget(s, _viewSelectedChainCode <= 0 ? TEST_NET??'Ethereum' : 'Eurus', additionalListViewBottomPaddingContentInset: 40, onTapHandler: (data) async => await _navigateToAssetDetails(data))??Builder(builder: (_) {final data = (jsonDecode(s.data) as List).where((e) => e['showAssets'] && e['address${_viewSelectedChainCode <= 0 ? TEST_NET??'Ethereum' : 'Eurus'}'] != null).map((e) {e['address'] = e['address${_viewSelectedChainCode <= 0 ? TEST_NET??'Ethereum' : 'Eurus'}']??'0xb9f5377a3f1c4748c71bd4ce2a16f13eecc4fec4'; return e;}).toList(); return ListView.builder(
                  padding: EdgeInsets.only(top: 57, left: 8, right: 8, bottom: 40 + MediaQuery.of(context).padding.bottom), // padding bottom need to consider SafeArea (i.e. 40 + MediaQuery.of(context).padding.bottom)
                  itemBuilder: (itemBuilderContext, i) => _cryptoAssetListItemCardTileWidget(data[i], onTap: () {_navigateToAssetDetails(data[i]);}),
                  itemCount: data.length,
                );})),
                // bottom: false ,),
                // asset list section header
                Align(alignment: Alignment.topCenter, child: Container(decoration: BoxDecoration(gradient: LinearGradient(begin: Alignment.center, end: Alignment.bottomCenter, colors: [Colors.white, Color(0x00FFFFFF)], ), ), child:
                  Padding(padding: EdgeInsets.symmetric(horizontal: 8), child: IntrinsicHeight(child: Row(children: [
                    CupertinoButton(onPressed: () {_setState(()=>_viewSelectedChainCode=TEST_NET!=null&&TEST_NET.isNotEmpty?-1:0);}, child: IntrinsicWidth(child: Column(children: [Text("Ethereum", style: Theme.of(context).textTheme.headline6.apply(color: Color(0xFF002251), fontWeightDelta: 1,)), FractionallySizedBox(widthFactor: .5, child: Container(child: Divider(color: _viewSelectedChainCode <= 0 ? Color(0xFF0041c4) : Colors.transparent, height: 0, thickness: 1.5, indent: 0, endIndent: 0)))]))),
                    VerticalDivider(width: 0, thickness: 1, color: Colors.black12, indent: 15, endIndent: 15,),
                    CupertinoButton(onPressed: () {_setState(()=>_viewSelectedChainCode=1);}, child: IntrinsicWidth(child: Column(children: [Text("Eurus", style: Theme.of(context).textTheme.headline6.apply(color: Color(0xFF002251), fontWeightDelta: 1,)), FractionallySizedBox(widthFactor: .5, child: Container(child: Divider(color: _viewSelectedChainCode == 1 ? Color(0xFF0041c4) : Colors.transparent, height: 0, thickness: 1.5, indent: 0, endIndent: 0)))]))),
                    Spacer(flex:1),
                    IconButton(icon: Icon(showBalance ? Icons.remove_red_eye : Icons.remove_red_eye_outlined, color: showBalance ? Colors.black12 : null,), color: Color(0xFF2684ff), onPressed: () {
                      _setState(() {
                        showBalance =!showBalance;
                      });
                    },),  
                    Builder(builder: (_) => IconButton(icon: Icon(Icons.more_horiz), color: Color(0xFF2684ff), onPressed: () {Navigator.of(_).push(MaterialPageRoute(builder: widget.getEditAssetsListWidgetBuilder(widget.encryptedAddress),));},)), // Flutter Material Ripple bug, ripple animating underneath the Rounded Wall Card List BoxDecoration Container 
                  ]))),
                ),)
              ])),
            ])
          )
        )))
      ],),))
    ), SafeArea(child: Align(alignment: Alignment.bottomCenter, child: CurvedNavigationBar(
      shader: LinearGradient(begin: Alignment.centerLeft, end: Alignment.centerRight, colors: [Color(0x9F4a00dd), Color(0x9F0036dd)]).createShader(Rect.fromLTRB(0,0,MediaQuery.of(context).size.width,0)),
      backgroundColor: Colors.transparent,
      color: Color(0xFF4a00dd),
      animationDuration: Duration(milliseconds: 100),
      height: 36,
      items: <Widget>[
        Icon(Icons.home, size: 20, color: Colors.white), // listview padding bottom 50 for icon size 30
        Icon(Icons.trending_up, size: 20, color: Colors.white),
        // Icon(Icons.search, size: 20, color: Colors.white),
        // Icon(Icons.messenger, size: 20, color: Colors.white),
        Icon(Icons.settings, size: 20, color: Colors.white),
      ],
      onTap: (index) {
        //Handle button tap
        log('Bottom Navigation Bar onTap index: ${index}');
      },
      onAnimationCompleted: (index) => Timer(Duration(milliseconds: 0), () {
        if (index == 0) _setState(() {
          navigationBarTabBuilderSelectedIndex = null;
        });
        if (index > 0 && index <= widget.navigationBarOnTapTabBuilders.length) _setState(() {
          navigationBarTabBuilderSelectedIndex = index-1;
        });
      }),
    ),))]);});

    // final __prefix = this.__prefix;

    // Import Wallet API
    Future<http.Response> importWallet(http.Client client, _anAddressPair) async {
      final _deviceId = '', _ts = DateTime.now().millisecondsSinceEpoch, _walletAddress = _anAddressPair.address.replaceFirst(r'0x', ''), _publicKey = _anAddressPair.publicKey.replaceFirst(r'0x', '');
      final _message = "deviceId=$_deviceId&timestamp=$_ts&walletAddress=$_walletAddress";
      // final _messageHash = keccakUtf8(_message);//.replaceFirst(r'0x', '');

      // uint8ListFromList (list..add(signature.v)..add(signature.r)..add(signature.s))
      // padUint8ListTo32(intToBytes(signatureData.r))

      final _ = jsonEncode(<String, Object>{
          'nonce': 'nonce',
          'deviceId': _deviceId,
          'timestamp': _ts,
          'walletAddress': _walletAddress,
          'publicKey': _publicKey,
          'sign': widget.canGetSignature(_message, _anAddressPair.privateKey),
        });
      print("importWallet POST request body json string $_");

      return http.post('http://18.141.43.75:8082/user/importWallet', 
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
        },
        body: jsonEncode(<String, Object>{
          'nonce': 'nonce',
          'deviceId': _deviceId,
          'timestamp': _ts,
          'walletAddress': _walletAddress,
          'publicKey': _publicKey,
          'sign': widget.canGetSignature(_message, _anAddressPair.privateKey).substring(0, 128),
        }),
      ).then((value) {
        final Map<String, dynamic> _bodyMap = jsonDecode(value.body); 
        widget.canPersistNormallyHandler({'apiAccessToken_${null??''}': _bodyMap['token'] as String, 'apiAccessTokenExpiryTime_${null??''}': "${_bodyMap['token']}=${_bodyMap['expiryTime'] as String}"}); 
        return value;
      });
    }
    testImportWallet(_addressPair) async {
      print("testImportWallet  ${(await importWallet(http.Client(), _addressPair)).body}");
      // log(response.body);
    }

    final _navigateToImportDecentralizeKeystoreLockerPasscodeSetup = ([String routeSettingsArgsMSeedPhrase]) {
      final Widget Function(BuildContext) Function(String) _keystoreLockerPasscodeSetupBodyBuilder = ([String _routeSettingsArgsMSeedPhrase]) => (context) {

        final _navigateToDecentralizeHomePage = () async {

          // await Navigator.of(context).push(MaterialPageRoute(builder: _decentralizedWalletHomePageBuilder(), fullscreenDialog: true));
          await Navigator.of(context).push(MyPageRoute(_decentralizedWalletHomePageBuilder(widget.navigationBarOnTapTabBuilders.map((e) => Builder(builder: e({
            'walletAccountEncryptedAddress': widget.encryptedAddress, 
            'walletAccountLogout': () {Navigator.popUntil(context, (route) => route.isFirst); widget.canPersistNormallyHandler({'bioAuthValidUntil': Uri(queryParameters: {widget.encryptedAddress: (-1).toString()}).query}).whenComplete(() => this.setState(() => navigationBarTabBuilderSelectedIndex = null));},
            'walletAccountBackupMnemonicSeedPhrases': () async {
              await Navigator.of(context).push(MaterialPageRoute(builder: (_materialPageRouteBuilderContext) => Scaffold(
                appBar: AppBar(title: Text("Backup Wallet Secret"), backgroundColor: Colors.transparent, elevation: 0,centerTitle:true),
                extendBodyBehindAppBar: true,
                backgroundColor: Colors.transparent,
                resizeToAvoidBottomInset: true,
                // body: _decentralizedPageBoardWallTemplate(Builder(builder: _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder,)),
                // body: _decentralizedPageBoardWallTemplate(Builder(builder: _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder(this),)),
                body: _decentralizedPageBoardWallTemplate(Builder(builder: _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder(() {}),)),
              )));
            },
          }))).toList()), RouteSettings(name: 'HomePage')));
        };

        final _textEditingController0 = TextEditingController();
        final _textEditingController1 = TextEditingController();
        final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
        bool _isFormValid = false;

        bool _willEnableBiometricsAuth = false;
        bool loading = false;

        Widget lockButton() {
          if (!loading) {
            return new Text(
              "Lock");
          } else  {
            return Container(child: CircularProgressIndicator(strokeWidth: 2,valueColor: AlwaysStoppedAnimation<Color>(Color(0xFF4a00dd))),width: 20,height: 20);
          }
        }

        return Column(children: [
          ListTile(dense: true, title: Icon(Icons.lock_outline_rounded, size: 120, color: Color(0xFF002251))),
          StatefulBuilder(builder: (context, _setState) => Expanded(flex: 1, child: KeyboardAvoider(autoScroll: true, child: Form(key: _formKey, onChanged: () {_setState(() {_isFormValid = _formKey.currentState.validate();});}, child: Column(children: [
            ListTile(
              title: Padding(child: Text("New Password", style: Theme.of(context).textTheme.bodyText2.apply(color: Color(0xFFc1c7d0))), padding: EdgeInsets.symmetric(vertical: 16)), 
              subtitle: TextFormField(
                decoration: _defaultTextFieldInputDecoration.copyWith(
                  hintText: "Create New Locker Password",
                  hintStyle: TextStyle(color: Color.fromRGBO(122, 134, 154, .5)),
                ),
                obscureText: true,
                controller: _textEditingController0,
              )
            ),
            ListTile(
              title: Padding(child: Text("Confirm Password", style: Theme.of(context).textTheme.bodyText2.apply(color: Color(0xFFc1c7d0))), padding: EdgeInsets.symmetric(vertical: 16)), 
              subtitle: TextFormField(
                decoration: _defaultTextFieldInputDecoration.copyWith(
                  hintText: "Type Locker Password again",
                  hintStyle: TextStyle(color: Color.fromRGBO(122, 134, 154, .5)),
                ),
                obscureText: true,
                controller: _textEditingController1,
                validator: (value) => value == _textEditingController0.text ? null : value.isNotEmpty ? 'Confirm password failed' : 'Type Locker Password again to confirm',
              )
            ),
            FutureBuilder(future: Future.delayed(Duration(milliseconds: 0), () async => await widget.canSupportsBiometricAuthenticatedHandler()), builder: (_, s) => SwitchListTile.adaptive(value: _willEnableBiometricsAuth, onChanged: s.connectionState != ConnectionState.done || !s.hasData || s.hasData && !s.data ? null : (newValue) {_setState(() => _willEnableBiometricsAuth = newValue);}, title: Text("Enable Biometrics", style: Theme.of(context).textTheme.subtitle1.apply(color: Color(0xFFc1c7d0))), contentPadding: EdgeInsets.all(16),)),
            ListTile(title: OutlinedButton(
              child: lockButton(),
              onPressed: !_isFormValid ? null : () async {
                // Generate root key for application to store in local
                // Uses this key to derivative Address and Private key
                _setState(() {
                   loading = true;
                });

                // _routeSettingsArgsMSeedPhrase ??= await compute(_genMnemonic, null); // this will assign a newly generated seed phrase to _routeSettingsArgsMSeedPhrase if null, then side-effect is that repeated invoking onPressed() won't generate uniquely different mnemonic seed phrase each time; but same seed once generated once will be used throughout lifecycle of this page widget instance in navigation stack
                final _mnemonicSeedPhrase = _routeSettingsArgsMSeedPhrase??await compute(_genMnemonic, 128);
                final anAddressPair = await compute(_genAddressPair, await compute(_genBase58, _routeSettingsArgsMSeedPhrase??_mnemonicSeedPhrase));

                final password = _textEditingController0.text;

                final _securePersist = widget.canPersistSecurelyHandler, _biometricPersist = widget.canPersistWithBiometricSecurelyHandler;
                final _encrypt = (c) => widget.canEncryptHandler(password, c), _pwdEncryptedAddress = _encrypt(anAddressPair.address);
                final k_accounts = __prefix+'accounts', k_accountPrivateKeys = __prefix+'accountPrivateKeys', k_createdAccountGeneratedMnemonicSeedPhrases = __prefix+'createdAccountGeneratedMnemonicSeedPhrases', k_encryptPasswordAddress = __prefix+_pwdEncryptedAddress;
                // await _securePersist(Uri(queryParameters: {_pwdEncryptedAddress: "0"}).query, k_accounts);
                final _error = await _securePersist({k_accounts: Uri(queryParameters: {_pwdEncryptedAddress: "0"}).query});

                // secure persistence of private key
                await _securePersist({k_accountPrivateKeys: Uri(queryParameters: {_pwdEncryptedAddress: _encrypt(anAddressPair.privateKey)}).query});
                // secure persistence of newly generated unique mnemonic seed phrase if HD wallet is newly created
                if (_routeSettingsArgsMSeedPhrase == null) await _securePersist({k_createdAccountGeneratedMnemonicSeedPhrases: Uri(queryParameters: {_pwdEncryptedAddress: _encrypt(_mnemonicSeedPhrase)}).query});
                
                if (_willEnableBiometricsAuth) {
                  // in order to ensure and enforce that iOS LocalAuthentication Prompt for user auth permission
                  if (Platform.isIOS) {
                    final _deleteResult = await _biometricPersist(k_encryptPasswordAddress, delete: true);
                    print("Platform.isIOS SecItemDelete Flutter MethodChannel.invokeMethod StorageCallback result in Future ${_deleteResult}");
                    if (_deleteResult == null) print("iOS KeyChain SecItemDelete errSecItemNotFound");
                    await _biometricPersist({k_encryptPasswordAddress: password});
                  }
                  final _bioAuth = await _biometricPersist({k_encryptPasswordAddress: password});
                  // final _bioAuth = await _biometricPersist({k_encryptPasswordAddress: password}, deleteBeforeWriteTwiceEnforceUpdateForiOSPrompt: Platform.isIOS);
                  // TODO - fork Flutter Biometric_Storage, CRUD call be purely CRUD, not implicitly update / overwrite existing conflicting key item while commanding only to create / insert an item; complex operation command can accept optional args or instruction or handler or callback to do the update / overwrite operation if key already exist and create / insert error
                  // final _bioAuthRead = await _biometricPersist(k_encryptPasswordAddress);
                  // TODO - fork Flutter Biometric_Storage, custom Platform System Biometric / Local Authentication prompts text string
                  // TODO - fork Flutter Biometric_Storage, custom enum LAPolicy.deviceOwnerAuthenticationWithBiometrics
                  print("_willEnableBiometricsAuth $_bioAuth");
                } else if (await widget.canSupportsBiometricAuthenticatedHandler()) {
                  final _deleteResult = await _biometricPersist(k_encryptPasswordAddress, delete: true);
                  print("shall not enable biometrics auth Platform BioStorage Flutter MethodChannel.invokeMethod StorageCallback result in Future ${_deleteResult}");
                  if (_deleteResult == null) print("prior item not found");
                }

                final findDefaultAccountEncryptedAddress = () async => Uri.decodeComponent(RegExp(r'([^=&?]*)=0').firstMatch(await _securePersist(k_accounts)).group(1));
                String decryptedString = widget.canDecryptHandler(password, await findDefaultAccountEncryptedAddress());

                widget.address = decryptedString??'password is not valid or SecureStorageKit().readValue failed';
                widget.encryptedAddress = _encrypt(decryptedString)??'password is not valid or SecureStorageKit().readValue failed';
                // set Auth UI of private key Handler to Web3 ETH Client
                // widget.initEthClientWithCanGetPrivateKeyHandler(_getPrivateKeyHandler);
                widget.createdAccountGeneratedMnemonicSeedPhrase = Uri(query: await widget.canPersistNormallyHandler('mPhraseBackuped')).queryParameters[widget.encryptedAddress] == null ? Uri(query: await widget.canPersistSecurelyHandler(__prefix+'createdAccountGeneratedMnemonicSeedPhrases')).queryParameters[widget.encryptedAddress] : null;

                
                testImportWallet(anAddressPair);
                
                await _navigateToDecentralizeHomePage();
                loading = false;
              },
              style: OutlinedButton.styleFrom(
                primary: Color(0xFF4a00dd), 
                side: BorderSide(color: Color(0xFF4a00dd)), 
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(13))
              ),
            )),
          ], mainAxisSize: MainAxisSize.max, mainAxisAlignment: MainAxisAlignment.start,),))))
        ], mainAxisSize: MainAxisSize.max);
      };
      Navigator.of(context).push(MaterialPageRoute(builder: _importDecentralizedWalletPageBuilder(Text("Create Wallet Locker"), _keystoreLockerPasscodeSetupBodyBuilder(routeSettingsArgsMSeedPhrase))));
    };

    final Widget Function(BuildContext) _importDecentralizedWalletPageBodyBuilder = (BuildContext context) {
      final _textEditingController = TextEditingController();
      final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
      bool _isFormValid = false;
      return StatefulBuilder(builder: (context, _setState) => Form(key: _formKey, onChanged: () {_setState(() {_isFormValid = _formKey.currentState.validate();});}, child: Column(
        children: [
          Center(child: ListTile(dense: true, title: Text('Mnemonic Seed Phrase', textAlign: TextAlign.center, style: Theme.of(context).textTheme.headline5.copyWith(color: Color(0xFF002251), fontSize: 24, fontFamily: 'packages/eurus/SFProDisplay', fontWeight: FontWeight.w600), ))),
          Expanded(flex: 1, child: Padding(
            padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: TextFormField(
              validator: _mnemonicValidator,
              decoration: _defaultTextFieldInputDecoration.copyWith(
                hintText: 'Enter your recovery mnemonic seed phrase to import your decentralized wallet',
                hintMaxLines: 4,
              ),
              maxLines: null,
              expands: true,
              textAlignVertical: TextAlignVertical.center,
              style: Theme.of(context).textTheme.headline5,
              autofocus: true,
              controller: _textEditingController,
            )
          )),
          ListTile(title: OutlinedButton(
            child: Text("Import"),
            onPressed: !_isFormValid ? null : () async {
              _navigateToImportDecentralizeKeystoreLockerPasscodeSetup(_textEditingController.text);
              _textEditingController.clear();
              _textEditingController.clearComposing();
            },
            style: OutlinedButton.styleFrom(
              primary: Color(0xFF4a00dd), 
              side: BorderSide(color: Color(0xFF4a00dd)), 
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(13))
            ),
          )),
        ],
      )));
    };
    
    final _navigateToImportDecentralizeRecoverySeedPhrase = () {
      Navigator.of(context).push(MaterialPageRoute(builder: _importDecentralizedWalletPageBuilder(Text("Import Wallet"), _importDecentralizedWalletPageBodyBuilder)));
    };

    final Widget Function(BuildContext) _selectWalletModePageBodyBuilder = (BuildContext context) {
      return
        Column(
          mainAxisSize: MainAxisSize.max,
          children: <Widget>[
            ListTile(dense: true, title: Text("Select Type", textAlign: TextAlign.center, style: Theme.of(context).textTheme.headline5.apply(color: Color(0xFF002251), fontFamily: 'packages/eurus/SFProDisplay'),)),
            // GestureDetector(
            //   child: ConstrainedBox(child: Stack(
            //     children: <Widget>[
            //       Center(child: Image.asset('assets/images/cardGroupCentralized.png', package: 'eurus')),
            //       // Align(alignment: Alignment(0.9, -0.8), child: IconButton(icon: Text('􀅵', style: TextStyle(fontFamily: 'packages/eurus/SFCompact', color: Color(0xFF00b3f3))), onPressed: () {log("0");},)),
            //     ]
            //   ), constraints: BoxConstraints(maxHeight: 194.0, maxWidth: MediaQuery.of(context).size.width),), 
            //   onTap: () {log("1");}, 
            // ),
            GestureDetector(
              child: ConstrainedBox(child: Stack(
                children: <Widget>[
                  Center(child: Image.asset('assets/images/cardGroupDecentralized2.png', package: 'eurus')),
                  // Align(alignment: Alignment(0.9, -0.8), child: IconButton(icon: Text('􀅵', style: TextStyle(fontFamily: 'packages/eurus/SFCompact', color: Color(0xFF00b3f3))), onPressed: () {log("0");},)),
                ]
              ), constraints: BoxConstraints(maxHeight: 194.0, maxWidth: MediaQuery.of(context).size.width),), 
              onTap: () {
                // TODO - create a new decentralized wallet with new random seed
                _navigateToImportDecentralizeKeystoreLockerPasscodeSetup();
              }, 
            ),
          ],
        )
      ;
    };

    final Widget Function(BuildContext) _selectWalletModePageBuilder = (context) {
      return Scaffold(
        appBar: AppBar(title: Text("Create New Wallet"), backgroundColor: Colors.transparent, elevation: 0,centerTitle:true),
        extendBodyBehindAppBar: true,
        backgroundColor: Colors.transparent,
        body: Builder(builder: (_) => _centralizedPageBoardWallTemplate(_selectWalletModePageBodyBuilder(_)))
      );
    };

    final Widget Function(BuildContext) _welcomePageBodyBuilder = (BuildContext context) {
      return
        DecoratedBox(
          decoration: BoxDecoration(
            color: Colors.white,
            image: DecorationImage(image: AssetImage('assets/images/bgLogin.png', package: 'eurus'), fit: BoxFit.cover),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: <Widget>[
              Expanded(flex: 2, child: GestureDetector(
                child: Container(child: Image.asset('assets/images/imgLogo.png', package: 'eurus')),
              )),
              // Bottom Sheet style bottom aligned half-shown vertical Card
              Expanded(flex: 1, child: DecoratedBox(
                decoration: BoxDecoration(
                  borderRadius: BorderRadiusDirectional.only(topStart: Radius.circular(30), topEnd: Radius.circular(30)),
                  color: Color(0xFF129fdd),
                ),
                child: Column(mainAxisAlignment: MainAxisAlignment.center, crossAxisAlignment: CrossAxisAlignment.stretch, children: <Widget>[
                  ListTile(title: FlatButton(
                    color: Colors.white,
                    child: Text('Create an account'),
                    textColor: Color(0xFF129fdd),
                    onPressed: () {
                      Navigator.of(context).push(MaterialPageRoute(builder: _selectWalletModePageBuilder));
                    },
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.all(Radius.circular(13))),
                  )),
                  ListTile(title: IntrinsicHeight(child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Flexible(/*flex: 1, fit: FlexFit.loose,*/ child: FractionallySizedBox(widthFactor: .67, child: TextButton(
                        onPressed: true ? null : () {
                          
                        },
                        child: Text.rich(TextSpan(text: "login", children: [TextSpan(text: "\nUser Account", style: Theme.of(context).textTheme.overline.apply(color: Colors.white))]), textAlign: TextAlign.end, style: TextStyle(letterSpacing: 1, fontWeight: FontWeight.bold,),),// Text("Login User Account", textAlign: TextAlign.center),
                        style: TextButton.styleFrom(primary: Colors.white, shape: StadiumBorder(), textStyle: Theme.of(context).textTheme.bodyText2,),
                      ))),
                      VerticalDivider(width: 16, thickness: 0, color: Colors.white70),
                      Flexible(/*flex: 1, fit: FlexFit.loose,*/ child: FractionallySizedBox(widthFactor: .67, child: TextButton(
                        onPressed: () {
                          _navigateToImportDecentralizeRecoverySeedPhrase();
                        },
                        child: Text.rich(TextSpan(text: "import", children: [TextSpan(text: "\nExisting Wallet", style: Theme.of(context).textTheme.overline.apply(color: Colors.white))]), textAlign: TextAlign.start, style: TextStyle(letterSpacing: 1, fontWeight: FontWeight.bold,)),//Text("Import Existing Wallet", textAlign: TextAlign.center,),
                        style: TextButton.styleFrom(primary: Colors.white, shape: StadiumBorder(), textStyle: Theme.of(context).textTheme.bodyText2,),
                      ))),
                    ]
                  ))),
                  InkWell(child:
                  Text('Customer Service', style: Theme
                      .of(context)
                      .textTheme
                      .overline
                      .apply(color: Colors.white),textAlign: TextAlign.center),
                      onTap: () {
                        LivechatInc.start_chat(
                            "12610959", "", "guest", "guest@gmail.com");
                      })
                ]),
              ),),
            ],
          )
        )
      ;
    };
    

    // return Builder(builder: _importDecentralizedWalletPageBuilder(Text("Import Wallet"), _importDecentralizedWalletPageBodyBuilder));
    // return Builder(builder: _decentralizedWalletHomePageBuilder());
    
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: FutureBuilder(
        future: Future.delayed(Duration(milliseconds: 150), () async {
          // TODO - here can await a result pop from a navigation push into wallet account selection dialog
          // TODO - iOS Keychain service keep retaining keychain items even after app uninstall and reinstall due to platform behaviour
          final findDefaultAccountEncryptedAddress = () async => Uri.decodeComponent(RegExp(r'([^=&?]*)=0').firstMatch(await widget.canPersistSecurelyHandler(__prefix+'accounts')).group(1));
          final encryptedAddress = await findDefaultAccountEncryptedAddress();
          if (encryptedAddress.isNotEmpty) {

            // TODO - should push an Auth Page onto upfront when AppLifecycleState or LifecycleReactor or LifecycleEventHandler or WidgetsBindingObserver onResumed
            final isAuthValidAuthUIPopResult = await Future.delayed(Duration(milliseconds: 50), () => Navigator.of(context).push(PageRouteBuilder(fullscreenDialog: true, opaque: false, pageBuilder: (pageBuilderContext, animation, secondaryAnimation) => Scaffold(backgroundColor: Colors.black.withOpacity(.8), body: Builder(builder: (builderContext) {
              final _textEditingController = TextEditingController();
              final _submit = () async {
                final decryptedAddress = widget.canDecryptHandler(_textEditingController.text, encryptedAddress);
                if (decryptedAddress != null) return Navigator.of(builderContext).maybePop(AddressPair(decryptedAddress, widget.canDecryptHandler(_textEditingController.text, Uri(query: await widget.canPersistSecurelyHandler(__prefix+'accountPrivateKeys')).queryParameters[encryptedAddress])));
                Scaffold.of(builderContext).showSnackBar(SnackBar(content: Text("Authentication Failure", textAlign: TextAlign.center,), backgroundColor: Colors.redAccent.shade100.withOpacity(.9),));
              };

              final _tryBioAuth = () => Future.delayed(Duration.zero, () async => DateTime.now().millisecondsSinceEpoch >= (int.tryParse(Uri(query: await widget.canPersistNormallyHandler('bioAuthValidUntil')).queryParameters[encryptedAddress]??'')??-1) ? null : await widget.canPersistWithBiometricSecurelyHandler(__prefix+encryptedAddress)).then((value) async {if (value != null) {_textEditingController.text = value; await _submit();}});
              _tryBioAuth();
              
              return
                KeyboardAvoider(autoScroll: true, child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [Container(decoration: BoxDecoration(borderRadius: BorderRadius.circular(15), color: Colors.white.withOpacity(.9)), padding: EdgeInsets.all(16), child: Column(
                    children: [
                      ListTile(title: Image.asset('assets/images/imgLogo.png', package: 'eurus')),
                      ListTile(title: Text("Welcome Back!", style: Theme.of(builderContext).textTheme.headline6.apply(fontWeightDelta: 2), textAlign: TextAlign.center)),
                      ListTile(title: TextField(textInputAction: TextInputAction.go, onSubmitted: (_) async => await _submit(), autofocus: true, obscureText: true, controller: _textEditingController, decoration: InputDecoration(border: OutlineInputBorder(borderRadius: BorderRadius.circular(16)), suffixIcon: null))),
                      // refactored
                      ListTile(title: TextButton(child: Text("LOGIN", style: Theme.of(builderContext).textTheme.caption.apply(color: Colors.blue)), onPressed: () async {
                        await _submit();
                      },), subtitle: TextButton(child: Text("Cancel to Switch Account"), onPressed: () => Navigator.of(builderContext).maybePop()))
                    ]
                  ))]
                ));
            })))));
            
            if (isAuthValidAuthUIPopResult == null) return null;

            widget.address = isAuthValidAuthUIPopResult.address;
            widget.encryptedAddress = encryptedAddress;
            // set Auth UI of private key Handler to Web3 ETH Client
            // widget.initEthClientWithCanGetPrivateKeyHandler(_getPrivateKeyHandler);
            widget.createdAccountGeneratedMnemonicSeedPhrase = Uri(query: await widget.canPersistNormallyHandler('mPhraseBackuped')).queryParameters[widget.encryptedAddress] == null ? Uri(query: await widget.canPersistSecurelyHandler(__prefix+'createdAccountGeneratedMnemonicSeedPhrases')).queryParameters[widget.encryptedAddress] : null;
            await widget.canPersistNormallyHandler({'bioAuthValidUntil': Uri(queryParameters: {widget.encryptedAddress: DateTime.now().add(Duration(days: 7)).millisecondsSinceEpoch.toString()}).query});
            testImportWallet(isAuthValidAuthUIPopResult);
          }
          return encryptedAddress;
        }),
        builder: (ctx, AsyncSnapshot<String> snapshot) => snapshot.connectionState != ConnectionState.done ? Center(child: Image.asset('assets/images/imgLogo.png', package: 'eurus')) : (snapshot.hasError || !snapshot.hasData || snapshot.hasData && snapshot.data.isEmpty) ? Builder(builder: _welcomePageBodyBuilder) : Builder(
          builder: (_) => Builder(builder: _decentralizedWalletHomePageBuilder(widget.navigationBarOnTapTabBuilders.map((e) => Builder(builder: e({
            'walletAccountEncryptedAddress': widget.encryptedAddress, 
            'walletAccountLogout': () {widget.canPersistNormallyHandler({'bioAuthValidUntil': Uri(queryParameters: {widget.encryptedAddress: (-1).toString()}).query}).whenComplete(() => this.setState(() {navigationBarTabBuilderSelectedIndex = null;}));},
            'walletAccountBackupMnemonicSeedPhrases': () async {
              await Navigator.of(context).push(MaterialPageRoute(builder: (_materialPageRouteBuilderContext) => Scaffold(
                appBar: AppBar(title: Text("Backup Wallet Secret"), backgroundColor: Colors.transparent, elevation: 0,centerTitle:true),
                extendBodyBehindAppBar: true,
                backgroundColor: Colors.transparent,
                resizeToAvoidBottomInset: true,
                // body: _decentralizedPageBoardWallTemplate(Builder(builder: _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder,)),
                // body: _decentralizedPageBoardWallTemplate(Builder(builder: _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder(this),)),
                body: _decentralizedPageBoardWallTemplate(Builder(builder: _backupWalletSecretMnemonicSeedPhrasesWidgetBuilder(() {}),)),
              )));
            },
          }))).toList()))
        ),
      )
    );


    // method 1: press to select account to push to password unlock
    // method 2: Stack password unlock widget always on top of / above wallet home page
    // method 3: StreamBuilder or FutureBuilder challenge
    // method 4: wallet home page widget internally build() a password unlock widget before actual wallet home page
    // method 5: uses of Widget class lifecycle like onMount, didMount, initState
    // method 6: for App Lifecycle resumes, WidgetBinding 
    // method 7: wrapper around wallet home page
  }

  int a = 0;
}

  /// Generate 12 / 24 words mnemonic phrase base on strength
  /// 
  /// [strength] = 128 for 12 words
  /// [strength] = 256 for 24 words
  String _genMnemonic([int strength = 128]) {
    return MnemonicKit().genMnemonicPhrase(strength: strength);
  }

  /// Check if mnemonic phrase is valid
  bool _validateMnemonic(String mPhrase) {
    return MnemonicKit().validateMnemonic(mPhrase);
  }

  /// Generate Base58 from mnemonic phrase
  String _genBase58(String mPhrase) {
    return MnemonicKit().mnemonicToBase58(mPhrase);
  }

  /// Generate address and private key from Base58
  AddressPair _genAddressPair(String b58) {
    return MnemonicKit().genAddressPairFromBase58(b58);
  }

  // form texteditingcontroller validator
  String _mnemonicValidator(String mPhrase) {
    // TODO - this mnemonic phrase validator currently keep complaining all inputs even typing is not yet completed nor ready to be dictated, i.e. return red warning messages since typing 1st ~ 11th letter, words, phrases
    return _validateMnemonic(mPhrase) ? null : "Invalid Mnemonic Phrase";
  }

