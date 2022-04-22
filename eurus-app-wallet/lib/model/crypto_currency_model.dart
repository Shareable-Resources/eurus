import 'package:app_crypto_icons/app_crypto_icons.dart';

class CryptoCurrencyModel {
  CryptoCurrencyModel({
    required this.currency,
    required this.symbol,
    required this.showAssets,
    @deprecated this.address,
    this.supported,
    this.addressEthereum,
    this.addressEurus,
    this.addressRinkeby,
    this.addressBSC,
    this.iconSource,
    this.imgUrl,
    this.order,
  });

  final String currency;
  final String symbol;
  String? addressEthereum;
  String? addressEurus;
  String? addressRinkeby;
  String? addressBSC;
  @deprecated
  final String? address;
  IconSourceType? iconSource;
  final bool? supported;
  bool showAssets;
  String? imgUrl;
  int? order;

  TokenViewState get addressEthereumViewState => _getViewState(addressEthereum);
  TokenViewState get addressEurusViewState => _getViewState(addressEurus);
  TokenViewState get addressRinkebyViewState => _getViewState(addressRinkeby);
  TokenViewState get addressBSCViewState => _getViewState(addressBSC);

  static CryptoCurrencyModel clone(CryptoCurrencyModel r) {
    return CryptoCurrencyModel(
      currency: r.currency,
      symbol: r.symbol,
      showAssets: r.showAssets,
      supported: r.supported,
      address: r.address,
      addressEthereum: r.addressEthereum,
      addressEurus: r.addressEurus,
      addressRinkeby: r.addressRinkeby,
      addressBSC: r.addressBSC,
      iconSource: r.iconSource,
      imgUrl: r.imgUrl,
      order: r.order,
    );
  }

  CryptoCurrencyModel.fromJson(Map<String, dynamic> json)
      : currency = json['currency'],
        symbol = json['symbol'],
        addressEthereum = json['addressEthereum'] ?? json['address'],
        addressEurus = json['addressEurus'],
        addressRinkeby = json['addressRinkeby'],
        addressBSC = json['addressBSC'],
        supported = json['supported'],
        showAssets = json['showAssets'],
        imgUrl = json['imgUrl'],
        order = json['order'],
        address = json['address'] ?? json['addressEthereum'];

  Map<String, dynamic> toJson() => {
        'currency': currency,
        'symbol': symbol,
        'addressEthereum': addressEthereum ?? address,
        'addressEurus': addressEurus,
        'addressRinkeby': addressRinkeby,
        'addressBSC': addressBSC,
        'supported': supported,
        'showAssets': showAssets,
        'imgUrl': imgUrl,
        'order': order
      };

  TokenViewState _getViewState(String? t) {
    if (t == '0x0') return TokenViewState.ROOT;
    if (t == null) return TokenViewState.SHOULD_NOT_DISPLAY;
    if (t == '') return TokenViewState.UNSUPPORTED;

    return TokenViewState.NORMAL;
  }
}

enum TokenViewState {
  ROOT, // Token address equals to 0x0
  SHOULD_NOT_DISPLAY, // Token address is null / not set
  UNSUPPORTED, // Token address is empty string ''
  NORMAL, // Any other valid ERC20 smart contract address 0x--40 length--
}
