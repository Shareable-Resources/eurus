import 'package:apihandler/apiHandler.dart';
import 'package:collection/collection.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/coinPrice.dart';
import 'package:flutter/widgets.dart';

import 'crypto_asset_list_item_card_tile_widget.dart';

class CryptoAssetListWidget extends StatelessWidget {
  final List? data;
  final String addressSuffix;
  final double listViewTopPaddingContentInset;
  final double additionalListViewBottomPaddingContentInset;
  final Future Function(dynamic)? onTapHandler;
  final bool Function(dynamic)? additionalWhereClauseFilterTest;
  final Map<String, dynamic>? assetBalanceCache;

  const CryptoAssetListWidget({
    Key? key,
    this.data,
    required this.addressSuffix,
    this.listViewTopPaddingContentInset = 57,
    this.additionalListViewBottomPaddingContentInset = 0,
    this.onTapHandler,
    this.additionalWhereClauseFilterTest,
    this.assetBalanceCache,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final filteredData = data
            ?.where(
                (e) => e['showAssets'] && e['address$addressSuffix'] != null)
            .where(additionalWhereClauseFilterTest ?? (e) => true)
            .map((e) {
          e['address'] = e['address$addressSuffix'];
          e['supported'] = e['address$addressSuffix'] is String &&
              (e['address$addressSuffix'] as String).isNotEmpty;
          return e;
        }).toList() ??
        [];
    final _chainAddressPrefix = "$addressSuffix${':'}";

    return FutureBuilder(
      future: getExchangeData(filteredData),
      builder: (_, AsyncSnapshot<CoinPriceList?> s) => ListView.builder(
        padding: EdgeInsets.only(
            top: listViewTopPaddingContentInset,
            left: 8,
            right: 8,
            bottom: additionalListViewBottomPaddingContentInset +
                MediaQuery.of(context).padding.bottom),
        // padding bottom need to consider SafeArea (i.e. 40 + MediaQuery.of(context).padding.bottom)
        itemBuilder: (itemBuilderContext, i) =>
            CryptoAssetListItemCardTileWidget(
          data: filteredData[i],
          onTap: () async {
            if (onTapHandler != null) await onTapHandler!(filteredData[i]);
          },
          assetBalancePlaceholder: assetBalanceCache?[
              _chainAddressPrefix + filteredData[i]['address']],
          notifyUpdateCacheAssetBalance: (k, v) {
            assetBalanceCache?[_chainAddressPrefix + k] = v;
            print("handler notifyUpdateCacheAssetBalance $assetBalanceCache");
          },
          coinPrice: s.data?.coinPriceList.firstWhereOrNull((element) =>
              element.symbol.toUpperCase() == filteredData[i]['symbol']),
        ),
        itemCount: filteredData.length,
      ),
    );
  }

  Future<CoinPriceList?> getExchangeData(List filteredData) async {
    final timeout = Duration(seconds: 3);
    if (common.coingeckoCoinsList.isEmpty) {
      try {
        final result = await api.get(
          url: 'https://api.coingecko.com/api/v3/coins/list',
          shouldPrintWrapped: false,
          timeout: timeout,
        );
        if (result is List<dynamic>)
          common.coingeckoCoinsList = common.coingeckoCoinsList;
      } catch (e) {}
    }

    final ids = common.coingeckoCoinsList
        .where((element) => filteredData
            .map((e) => e['symbol'])
            .contains((element['symbol'] as String).toUpperCase()))
        .map((e) => e['id'])
        .toList()
        .join(',');
    CoinPriceList? coinPriceList;
    try {
      final result = await api.get(
        url:
            "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=$ids&order=market_cap_desc&sparkline=false",
        shouldPrintWrapped: false,
        timeout: timeout,
      );
      coinPriceList =
          result is List<dynamic> ? CoinPriceList.fromJson(result) : null;
    } catch (e) {
      coinPriceList = common.coinPriceList;
    }

    // FIXME: temp fix eun price
    coinPriceList?.coinPriceList.add(CoinPrice(
      id: 'eurus',
      symbol: 'eun',
      name: 'Eurus',
      image: 'images/Eurus_Blue.png',
      currentPrice: 0.1,
    ));
    common.coinPriceList = coinPriceList;
    return coinPriceList;
  }
}
