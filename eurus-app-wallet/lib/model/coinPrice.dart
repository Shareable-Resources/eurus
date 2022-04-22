/// CoinPriceList Class
class CoinPriceList {
  List<CoinPrice> coinPriceList;

  CoinPriceList.fromJson(List<dynamic> json)
      : coinPriceList = CoinPrice.fromJsons(json);
}

/// CoinPrice Class
class CoinPrice {
  String id;
  String symbol;
  String name;
  String image;
  double? currentPrice;
  double? marketCap;
  double? marketCapRank;
  double? fullyDilutedValuation;
  double? totalVolume;
  double? high24h;
  double? low24h;
  double? priceChange24h;
  double? priceChangePercentage24h;
  double? marketCapChange24h;
  double? marketCapChangePercentage24h;
  double? circulatingSupply;
  double? totalSupply;
  double? maxSupply;
  double? ath;
  double? athChangePercentage;
  String? athDate;
  double? atl;
  double? atlChangePercentage;
  String? atlDate;
  String? lastUpdated;

  /// init CoinPrice
  CoinPrice({
    required this.id,
    required this.symbol,
    required this.name,
    required this.image,
    this.currentPrice,
    this.marketCap,
    this.marketCapRank,
    this.fullyDilutedValuation,
    this.totalVolume,
    this.high24h,
    this.low24h,
    this.priceChange24h,
    this.priceChangePercentage24h,
    this.marketCapChange24h,
    this.marketCapChangePercentage24h,
    this.circulatingSupply,
    this.totalSupply,
    this.maxSupply,
    this.ath,
    this.athChangePercentage,
    this.athDate,
    this.atl,
    this.atlChangePercentage,
    this.atlDate,
    this.lastUpdated,
  });

  /// fromJsons to List<CoinPrice>
  static List<CoinPrice> fromJsons(List<dynamic> jsons) {
    // MyLog.d('UserAsset.fromJsons');
    List<CoinPrice> assets = [];

    jsons.forEach((item) {
      CoinPrice ua = CoinPrice.fromJson(item);
      assets.add(ua);
    });

    return assets;
  }

  /// from json to CoinPrice Class
  CoinPrice.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        symbol = json['symbol'],
        name = json['name'],
        image = json['image'],
        currentPrice = json['current_price']?.toDouble(),
        marketCap = json['market_cap']?.toDouble(),
        marketCapRank = json['market_cap_rank']?.toDouble(),
        fullyDilutedValuation = json['fully_diluted_valuation']?.toDouble(),
        totalVolume = json['total_volume']?.toDouble(),
        high24h = json['high_24h']?.toDouble(),
        low24h = json['low_24h']?.toDouble(),
        priceChange24h = json['price_change_24h']?.toDouble(),
        priceChangePercentage24h =
            json['price_change_percentage_24h']?.toDouble(),
        marketCapChange24h = json['market_cap_change_24h']?.toDouble(),
        marketCapChangePercentage24h =
            json['market_cap_change_percentage_24h']?.toDouble(),
        circulatingSupply = json['circulating_supply']?.toDouble(),
        totalSupply = json['total_supply']?.toDouble(),
        maxSupply = json['max_supply']?.toDouble(),
        ath = json['ath']?.toDouble(),
        athChangePercentage = json['ath_change_percentage']?.toDouble(),
        athDate = json['ath_date'],
        atl = json['atl']?.toDouble(),
        atlChangePercentage = json['atl_change_percentage']?.toDouble(),
        atlDate = json['atl_date'],
        lastUpdated = json['last_updated'];

  ///  CoinPrice Class toJson
  Map<String, dynamic> toJson() => {
        'id': this.id,
        'symbol': this.symbol,
        'name': this.name,
        'image': this.image,
        'current_price': this.currentPrice,
        'market_cap': this.marketCap,
        'market_cap_rank': this.marketCapRank,
        'fully_diluted_valuation': this.fullyDilutedValuation,
        'total_volume': this.totalVolume,
        'high_24h': this.high24h,
        'low_24h': this.low24h,
        'price_change_24h': this.priceChange24h,
        'price_change_percentage_24h': this.priceChangePercentage24h,
        'market_cap_change_24h': this.marketCapChange24h,
        'market_cap_change_percentage_24h': this.marketCapChangePercentage24h,
        'circulating_supply': this.circulatingSupply,
        'total_supply': this.totalSupply,
        'max_supply': this.maxSupply,
        'ath': this.ath,
        'ath_change_percentage': this.athChangePercentage,
        'ath_date': this.athDate,
        'atl': this.atl,
        'atl_change_percentage': this.atlChangePercentage,
        'atl_date': this.atlDate,
        'last_updated': this.lastUpdated,
      };
}
