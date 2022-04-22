import 'package:apihandler/apiHandler.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/coinPrice.dart';
import 'package:flutter_neumorphic/flutter_neumorphic.dart';

class CoinPriceRankingPage extends StatefulWidget {
  CoinPriceRankingPage({Key? key}) : super(key: key);

  @override
  _CoinPriceRankingPageState createState() => _CoinPriceRankingPageState();
}

class _CoinPriceRankingPageState extends State<CoinPriceRankingPage> {
  CoinPriceList? coinPriceList;
  late Timer timer;
  bool visible = true;
  int currentSelectButton = 0;

  Future<CoinPriceList?> getExchangeData() async {
    /// getTopErc20CoinPrice
    var result = await api.get(
      url:
          "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=ethereum%2Ctether%2Cchainlink%2Cbinance-usd%2Cusd-coin%2Cyearn-finance%2Cdai%2Comisego%2Cbinancecoin%2Cuniswap%2Cvechain%2Caave%2Chuobi-token%2Csushi%2Ctrue-usd%2Ccdai%2Cswipe%2Cbasic-attention-token%2Cusdk%2Cwrapped-bitcoin%2Czilliqa%2Chavven%2Cokb%2Cband-protocol%2Cmaker%2Chusd%2C0x%2Cpaxos-standard%2Ccompound-ether%2Creserve-rights-token%2Cbalancer&order=market_cap_desc&per_page=100&page=1&sparkline=false",
      shouldPrintWrapped: false,
      timeout: Duration(seconds: 3),
    );

    if (result is! List) return null;
    CoinPriceList coinPriceListData = CoinPriceList.fromJson(result);

    if (this.mounted) {
      setState(() {
        visible = false;
      });
      Future.delayed(const Duration(milliseconds: 300), () {
        setState(() {
          coinPriceList = coinPriceListData;
          visible = true;
        });
      });
    }

    //handle USDM, BTCM, ETHM
    CoinPrice usdm = CoinPrice.fromJson(coinPriceListData.coinPriceList
        .firstWhere((coin) => coin.id == 'usd-coin')
        .toJson());
    usdm.id = "usdm";
    usdm.name = "USDM";
    usdm.image = common.getCustomCryptoSymbol('USDM')!;
    CoinPrice btcm = CoinPrice.fromJson(coinPriceListData.coinPriceList
        .firstWhere((coin) => coin.id == 'wrapped-bitcoin')
        .toJson());
    btcm.id = "btcm";
    btcm.name = "BTCM";
    btcm.image = common.getCustomCryptoSymbol('BTCM')!;
    CoinPrice ethm = CoinPrice.fromJson(coinPriceListData.coinPriceList
        .firstWhere((coin) => coin.id == 'ethereum')
        .toJson());
    ethm.id = "ethm";
    ethm.name = "ETHM";
    ethm.image = common.getCustomCryptoSymbol('ETHM')!;
    coinPriceListData.coinPriceList.insertAll(0, [usdm, btcm, ethm]);

    return coinPriceListData;
  }

  @override
  void initState() {
    super.initState();
    getExchangeData();
    timer = Timer.periodic(Duration(seconds: 5), (timer) {
      getExchangeData();
    });
  }

  @override
  void dispose() {
    timer.cancel();
    super.dispose();
  }

  List<Widget> coinPriceRankingUI({BuildContext? context}) {
    List<Widget> widgetList = <Widget>[];
    if (coinPriceList?.coinPriceList != null) {
      for (var i = 0; i < coinPriceList!.coinPriceList.length; i++) {
        CoinPrice coinPrice = coinPriceList!.coinPriceList[i];
        widgetList.add(
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 8),
            child: Row(children: <Widget>[
              Expanded(
                  flex: 5,
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.only(right: 8),
                        child: Container(
                          height: 20,
                          width: 20,
                          child: Uri.parse(coinPrice.image).isAbsolute
                              ? CachedNetworkImage(
                                  imageUrl: coinPrice.image,
                                  errorWidget: (context, url, error) =>
                                      Icon(Icons.error),
                                )
                              : Image.asset(coinPrice.image,
                                  package: 'euruswallet'),
                        ),
                      ),
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: <Widget>[
                          Padding(
                            padding: EdgeInsets.only(bottom: 2),
                            child: Text(coinPrice.name,
                                style: FXUI.normalTextStyle
                                    .copyWith(fontSize: 16)),
                          ),
                          Text(
                            "${'MARKETS.VOLUME'.tr()} \$ " +
                                common.numberFormat(
                                    minDecimal: 3,
                                    number: coinPrice.totalVolume.toString()),
                            style: TextStyle(fontSize: 12),
                          )
                        ],
                      ),
                    ],
                  )),
              Expanded(
                  flex: 2,
                  child: Column(
                      crossAxisAlignment: CrossAxisAlignment.end,
                      children: [
                        AnimatedOpacity(
                            opacity: visible ? 1.0 : 0.0,
                            duration: Duration(milliseconds: 500),
                            child: Padding(
                              padding: EdgeInsets.only(bottom: 3.0),
                              child: Text(
                                  "\$ " +
                                      common.numberFormat(
                                          minDecimal: 3,
                                          maxDecimal: 8,
                                          number: coinPrice.currentPrice!
                                              .toStringAsFixed(3)),
                                  style: FXUI.normalTextStyle
                                      .copyWith(fontSize: 14),
                                  textAlign: TextAlign.right),
                            )),
                        FittedBox(
                          fit: BoxFit.contain,
                          child: AnimatedOpacity(
                            opacity: visible ? 1.0 : 0.0,
                            duration: Duration(milliseconds: 300),
                            child: Row(
                              children: [
                                SizedBox(
                                    width: 11,
                                    height: 11,
                                    child: Image.asset(
                                        coinPrice.priceChangePercentage24h!
                                                .toStringAsFixed(2)
                                                .contains("-")
                                            ? 'images/priceDownIcon.png'
                                            : 'images/priceUpIcon.png',
                                        package: 'euruswallet')),
                                Padding(
                                  padding: const EdgeInsets.only(left: 5.0),
                                  child: Text(
                                      (coinPrice.priceChangePercentage24h!
                                              .toStringAsFixed(2) +
                                          "%"),
                                      style: FXUI.normalTextStyle.copyWith(
                                          color: coinPrice
                                                  .priceChangePercentage24h!
                                                  .toStringAsFixed(2)
                                                  .contains("-")
                                              ? Colors.red
                                              : Colors.green,
                                          fontSize: 12)),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ])),
            ]),
          ),
        );
        widgetList.add(Divider());
      }
    }
    return widgetList;
  }

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    return SafeArea(
      child: NestedScrollView(
        headerSliverBuilder: (BuildContext context, bool innerBoxIsScrolled) {
          return <Widget>[
            SliverAppBar(
              foregroundColor: Colors.transparent,
              backgroundColor: Colors.transparent,
              expandedHeight: 205.0,
              flexibleSpace: FlexibleSpaceBar(
                background: common.getBannerPath(context: context) != null
                    ? Padding(
                        padding:
                            const EdgeInsets.fromLTRB(20.0, 0.0, 20.0, 10.0),
                        child: GestureDetector(
                          onTap: () =>
                              common.routeToRewardDetailPage(context: context),
                          child: Neumorphic(
                            child: Image.asset(
                              common.getBannerPath(context: context)!,
                              package: 'euruswallet',
                              fit: BoxFit.fill,
                            ),
                            style: FXUI.neumorphicBannerImage,
                          ),
                        ),
                      )
                    : Container(),
              ),
            ),
          ];
        },
        body: coinPriceList?.coinPriceList != null
            ? SingleChildScrollView(
                child: Container(
                  decoration: new BoxDecoration(color: Colors.white),
                  child: Padding(
                      padding: const EdgeInsets.symmetric(
                          vertical: 15.0, horizontal: 24.0),
                      child: Column(children: coinPriceRankingUI())),
                ),
              )
            : Center(child: CircularProgressIndicator()),
      ),
    );
  }
}
