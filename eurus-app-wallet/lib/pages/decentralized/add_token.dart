import 'package:apihandler/apiHandler.dart';
import 'package:app_crypto_icons/app_crypto_icons.dart';
import 'package:collection/src/iterable_extensions.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/asset_row.dart';
import 'package:euruswallet/model/crypto_currency_model.dart';
import 'package:euruswallet/model/coinPrice.dart';

class AddTokenPage extends StatefulWidget {
  AddTokenPage({
    required this.userSuffix,
    Key? key,
  }) : super(key: key);

  final String userSuffix;

  _AddTokenPageState createState() => _AddTokenPageState();
}

class _AddTokenPageState extends State<AddTokenPage> {
  final _searchFormKey = GlobalKey<FormState>();
  final TextEditingController _searchTc = TextEditingController();

  /// Tokens that are already added into asset list
  Map<String, CryptoCurrencyModel> _curAssetList = {};

  /// Tokens to be shown initially, (Token added and 30 supported token)
  List<CryptoCurrencyModel> _displayList = [];

  /// Search result
  Map<String, CryptoCurrencyModel> _searchResultList = {};

  /// Fetched supported list
  List<dynamic> _fetchedSupportedTokens = [];

  // String _searchKey = '0x70c621f949b6556c4545707a2d5d73A776b98359';
  // String _searchKey = '0x4C0fBE1BB46612915E7967d2C3213cd4d87257AD';
  String _searchKey = '';
  bool _onSearchToken = false;
  bool _emptyResult = false;
  bool _invalidAddress = false;

  Timer? _debounceTimer;
  Timer? _searchBtnTimer;

  bool get _supportedTokensInited => _fetchedSupportedTokens.length > 0;

  @override
  void initState() {
    _getSupportedList();
    _initCurrentAssetList().then((value) => _initDisplayList());

    super.initState();
  }

  @override
  void dispose() {
    _debounceTimer?.cancel();
    _searchBtnTimer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    List<Widget> _tokenList = [];

    if (_supportedTokensInited) {
      _displayList.forEach((e) {
        if (!isEmptyString(string: _searchKey)) {
          bool matchC =
              e.currency.toLowerCase().contains(_searchKey.toLowerCase());
          bool matchS =
              e.symbol.toLowerCase().contains(_searchKey.toLowerCase());
          bool matchEthA = e.addressEthereum != null &&
              e.addressEthereum!
                  .toLowerCase()
                  .contains(_searchKey.toLowerCase());
          bool matchEruA = e.addressEurus != null &&
              e.addressEurus!.toLowerCase().contains(_searchKey.toLowerCase());
          bool matchRinkebyA = e.addressRinkeby != null &&
              e.addressRinkeby!
                  .toLowerCase()
                  .contains(_searchKey.toLowerCase());

          if (!matchC && !matchS && !matchEthA && !matchEruA && !matchRinkebyA)
            return;
        }
        _tokenList.add(_genTokenRow(e));
      });

      if (_tokenList.length == 0 && _searchResultList.length == 0) {
        _tokenList.add(_genSearchOtherRow());
      }

      _searchResultList.values.toList().forEach((e) {
        _tokenList.add(_genTokenRow(e));
      });
    } else {
      _tokenList.add(
        Padding(
          padding: EdgeInsets.symmetric(vertical: 8),
          child: AssetRow(
            vPadding: 18,
            child: ListTile(
              contentPadding: EdgeInsets.symmetric(horizontal: 25),
              title: Text('COMMON.LOADING_W_DOT'.tr()),
              leading: SizedBox(
                width: 45,
                height: 45,
                child: Center(
                  child: CircularProgressIndicator(
                    valueColor: AlwaysStoppedAnimation<Color>(
                      FXColor.mainDeepBlueColor,
                    ),
                  ),
                ),
              ),
            ),
          ),
        ),
      );
    }

    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          image: DecorationImage(
            image: AssetImage(
              isCentralized()
                  ? "images/backgroundImage.png"
                  : "images/backgroundImage2.png",
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
              title: Text('ADD_TOKEN_PAGE.MAIN_TITLE'.tr()),
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
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Padding(
                      padding: EdgeInsets.only(
                          top: 35, left: 35, right: 35, bottom: 20),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          Text(
                            'ADD_TOKEN_PAGE.SEARCH_BOX.LABEL'.tr(),
                            style: FXUI.normalTextStyle.copyWith(
                              color: FXColor.lightGray,
                              fontSize: 14,
                            ),
                          ),
                          Padding(
                            padding: EdgeInsets.only(top: 16),
                          ),
                          Form(
                            key: _searchFormKey,
                            child: TextFormField(
                                controller: _searchTc,
                                onChanged: (v) {
                                  _onSearchKeyChange(v);
                                },
                                decoration: InputDecoration(
                                  filled: true,
                                  fillColor: FXColor.veryLightGray,
                                  border: OutlineInputBorder(
                                    borderSide: BorderSide.none,
                                    borderRadius: FXUI.cricleRadius,
                                  ),
                                  errorBorder: OutlineInputBorder(
                                    borderSide:
                                        BorderSide(width: 1, color: Colors.red),
                                    borderRadius: FXUI.cricleRadius,
                                  ),
                                  hintText:
                                      'ADD_TOKEN_PAGE.SEARCH_BOX.PLACEHOLDER'
                                          .tr(),
                                  hintStyle: FXUI.normalTextStyle.copyWith(
                                      color: FXColor.centralizedGrayTextColor),
                                  suffixIcon: IconButton(
                                    icon: Icon(Icons.close),
                                    onPressed: () {
                                      _searchTc.text = '';
                                      setState(() {
                                        _searchKey = '';
                                        _searchResultList = {};
                                        _emptyResult = false;
                                        _invalidAddress = false;
                                      });
                                    },
                                  ),
                                ),
                                maxLines: 1,
                                validator: (s) {
                                  if (_invalidAddress) {
                                    return "ADD_TOKEN_PAGE.ERROR.INVALID_ADDRESS"
                                        .tr();
                                  }

                                  return null;
                                },
                                autovalidateMode: AutovalidateMode.always),
                          )
                        ],
                      ),
                    ),
                    Expanded(
                      flex: 1,
                      child: SingleChildScrollView(
                        padding:
                            EdgeInsets.only(left: 23, right: 23, bottom: 35),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.stretch,
                          children: _tokenList,
                        ),
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

  Future<Null> _initDisplayList() async {
    List<CryptoCurrencyModel> tokenListToDisplay = [];

    List<CryptoCurrencyModel> curList = _curAssetList.values.toList();

    List<CryptoCurrencyModel> supportedTokens = await _getSupportedList();
    CoinPriceList? coinPriceList =
        await common.getCoingekoCryptoImage(supportedTokens);

    for (var i = 0; i < supportedTokens.length; i++) {
      CryptoCurrencyModel token = CryptoCurrencyModel.clone(supportedTokens[i]);
      CoinPrice? coin = coinPriceList?.coinPriceList.firstWhereOrNull(
          (element) => element.symbol.toUpperCase() == token.symbol);
      if (coin != null) token.imgUrl = coin.image;
      token.iconSource =
          await AppCryptoIcons.ckIconSource(token.symbol, imgUrl: coin?.image);
      tokenListToDisplay.add(token);
    }

    for (var i = 0; i < curList.length; i++) {
      if (!(curList[i].supported ?? false)) {
        CryptoCurrencyModel token = CryptoCurrencyModel(
          currency: curList[i].currency,
          symbol: curList[i].symbol,
          address: curList[i].address,
          supported: false,
          iconSource: IconSourceType.url,
          showAssets: curList[i].showAssets,
          imgUrl: curList[i].imgUrl,
        );
        tokenListToDisplay.add(token);
      }
    }

    setState(() {
      _displayList.addAll(tokenListToDisplay);
    });
  }

  Future<Null> _initCurrentAssetList() async {
    String? assetsList =
        await NormalStorageKit().readValue("assetsList_${widget.userSuffix}");

    List aJsonList = assetsList != null ? jsonDecode(assetsList) : [];
    Map<String, CryptoCurrencyModel> _finalList = {};
    for (var i = 0; i < aJsonList.length; i++) {
      CryptoCurrencyModel token = CryptoCurrencyModel.fromJson(aJsonList[i]);
      _finalList.addAll({token.symbol: token});
    }

    setState(() {
      _curAssetList = _finalList;
    });
  }

  Widget _genTokenRow(CryptoCurrencyModel t) {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: 8),
      child: AssetRow(
        vPadding: 18,
        child: ListTile(
          contentPadding: EdgeInsets.symmetric(horizontal: 25),
          title: Text(t.symbol),
          subtitle: Text(t.currency),
          leading: common.getIcon(t.symbol, 50,
              source: t.iconSource ?? IconSourceType.none, imgUrl: t.imgUrl),
          trailing: isAssetEditable(t)
              ? StatefulBuilder(
                  builder: (_, setState) => IconButton(
                    onPressed: () => setState(() {
                      if (!_ckIsAdded(t)) {
                        _curAssetList.addAll({t.symbol: t});
                        if (!(t.supported ?? false)) {
                          _displayList.add(t);
                        }
                      } else {
                        _curAssetList.remove(t.symbol);
                      }
                      _updateAssetsList();
                    }),
                    icon: Icon(
                      // Icons.favorite_rounded,
                      _ckIsAdded(t) ? Icons.remove_circle : Icons.add,
                      size: 30,
                      // color: _ckIsAdded(t) ? Colors.redAccent : Colors.black12,
                      color: _ckIsAdded(t) ? Colors.redAccent : Colors.green,
                    ),
                  ),
                )
              : null,
        ),
      ),
    );
  }

  Widget _genSearchOtherRow() {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: 8),
      child: FlatButton(
        padding: EdgeInsets.zero,
        shape: RoundedRectangleBorder(borderRadius: FXUI.cricleRadius),
        onPressed: () {
          _fetchByAddress();
        },
        child: AssetRow(
          vPadding: 18,
          child: ListTile(
            contentPadding: EdgeInsets.symmetric(horizontal: 25),
            title: Text(_onSearchToken
                ? 'ADD_TOKEN_PAGE.SEARCHING'.tr()
                : _emptyResult
                    ? 'ADD_TOKEN_PAGE.ERROR.TOKEN_NOT_FOUND.TITLE'.tr()
                    : 'ADD_TOKEN_PAGE.SEARCH_TOKEN_BTN'.tr()),
            subtitle: Text(_onSearchToken
                ? 'ADD_TOKEN_PAGE.SEARCHING_ADDRESS'.tr(args: [_searchKey])
                : _emptyResult
                    ? 'ADD_TOKEN_PAGE.ERROR.TOKEN_NOT_FOUND.DESC'.tr()
                    : 'ADD_TOKEN_PAGE.ADDRESS_PREFIX'.tr(args: [_searchKey])),
            leading: Container(
              child: _onSearchToken
                  ? SizedBox(
                      width: 45,
                      height: 45,
                      child: Center(
                          child: CircularProgressIndicator(
                        valueColor: AlwaysStoppedAnimation<Color>(
                          FXColor.mainDeepBlueColor,
                        ),
                      )),
                    )
                  : Container(
                      width: 55,
                      height: 55,
                      decoration: BoxDecoration(
                          color: FXColor.mainDeepBlueColor.withOpacity(0.15),
                          borderRadius: BorderRadius.circular(45)),
                      child: Icon(
                        _emptyResult ? Icons.search_off : Icons.search_rounded,
                        size: 35,
                        color: FXColor.mainDeepBlueColor,
                      ),
                    ),
            ),
          ),
        ),
      ),
    );
  }

  void _onSearchKeyChange(String k) {
    if (_debounceTimer?.isActive ?? false) _debounceTimer?.cancel();

    _debounceTimer = Timer(Duration(milliseconds: 500), () {
      if (_searchBtnTimer?.isActive ?? false) _searchBtnTimer?.cancel();

      setState(() {
        _searchKey = k;
        _searchResultList = {};
        _emptyResult = false;
        _invalidAddress = false;
      });
    });
  }

  bool _ckIsAdded(CryptoCurrencyModel t) {
    if (_curAssetList.containsKey(t.symbol)) return true;

    return false;
  }

  void _updateAssetsList() async {
    if (_curAssetList.length <= 0) {
      await NormalStorageKit().deleteValue('assetsList_${widget.userSuffix}');
      return;
    }

    List<CryptoCurrencyModel> _orderedList = [];

    int idx = 0;
    _curAssetList.values.toList().forEach((e) {
      _orderedList.add(e..order = idx);
      idx++;
    });
    String rdString = jsonEncode(_orderedList);

    await NormalStorageKit()
        .setValue(rdString, 'assetsList_${widget.userSuffix}');
  }

  void _fetchByAddress() async {
    setState(() {
      _onSearchToken = true;
    });
    if (!RegExp(r"^0[xX][a-fA-F0-9]{40}$").hasMatch(_searchKey)) {
      setState(() {
        _onSearchToken = false;
        _invalidAddress = true;
      });
      return;
    }
    if (_searchFormKey.currentState != null &&
        !_searchFormKey.currentState!.validate()) {
      setState(() {
        _onSearchToken = false;
      });
      return;
    }

    apiHandler
        .get(
      "https://api.coingecko.com/api/v3/coins/ethereum/contract/$_searchKey",
    )
        .then((value) {
      final symbol = value['symbol'] is String
          ? (value['symbol'] as String).toUpperCase()
          : '';
      CryptoCurrencyModel result = CryptoCurrencyModel(
        currency: value['name'],
        symbol: symbol,
        address: _searchKey,
        addressEthereum: _searchKey,
        addressEurus: '',
        addressRinkeby: _searchKey,
        iconSource: IconSourceType.url,
        imgUrl: value['image']['large'],
        showAssets: true,
        supported: true,
      );
      setState(() {
        _searchResultList.addAll({symbol: result});
      });
    }).catchError((e) {
      setState(() {
        _emptyResult = true;
      });
    }).whenComplete(
      () => setState(() {
        _onSearchToken = false;
      }),
    );
  }

  Future<List<CryptoCurrencyModel>> _getSupportedList() async {
    Map<String, CryptoCurrencyModel> tokensInMap =
        await CommonMethod().getSupportedTokens();

    setState(() {
      _fetchedSupportedTokens = tokensInMap.values.toList();
    });

    return tokensInMap.values.toList();
  }
}
