import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/model/coinPrice.dart';

import '../common/callApiHandler.dart';
import '../common/commonMethod.dart';
import '../model/faucetRequest.dart';
import 'customDialogBox.dart';

class CryptoAssetListItemCardTileWidget extends StatelessWidget {
  final Map<String, dynamic> data;
  final GestureTapCallback? onTap;
  final Widget? chainIndicator;
  final dynamic assetBalancePlaceholder;
  final Function(dynamic, dynamic)? notifyUpdateCacheAssetBalance;
  final CoinPrice? coinPrice;
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  CryptoAssetListItemCardTileWidget({
    Key? key,
    required this.data,
    this.onTap,
    this.chainIndicator,
    this.assetBalancePlaceholder,
    this.notifyUpdateCacheAssetBalance,
    this.coinPrice,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      child: Column(
        children: [
          haveFreeTokenButton &&
                  chainIndicator != null &&
                  (common.showFaucetBtn ?? false)
              ? Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Container(
                        padding: EdgeInsets.only(top: 15),
                        width: 72,
                        height: 47),
                    chainIndicator ?? Container(),
                    Container(
                      padding: EdgeInsets.only(top: 15),
                      width: 72,
                      height: 47,
                      child: SubmitButton(
                        btnController: btnController,
                        fontSize: 13,
                        label: "FAUCET_FNC.FAUCET".tr(),
                        onPressed: () async {
                          FaucetRequest faucetRequest = await api.faucetRequest(
                              symbol: common.cryptoCurrencyModelMap?['symbol']);
                          String desc = '';
                          if (faucetRequest.data?.status == 1) {
                            desc = 'FAUCET_FNC.DESC'.tr();
                          }
                          if (faucetRequest.data?.status == 2) {
                            desc = 'FAUCET_FNC.DESC2'.tr();
                          }
                          if (faucetRequest.data?.status == 3) {
                            desc = 'FAUCET_FNC.DESC3'.tr();
                          }
                          await showDialog(
                            context: context,
                            builder: (BuildContext context) {
                              return CustomDialogBox(
                                title: 'COMMON.MESSAGE'.tr(),
                                descriptions: desc,
                                buttonText: "COMMON.OK".tr(),
                              );
                            },
                          );
                          btnController.reset();
                        },
                      ),
                    )
                  ],
                )
              : chainIndicator ?? Container(),
          ListTile(
            leading: FutureBuilder<Widget>(
              future: common.getCryptoIcon(
                data['symbol'],
                40,
                imgUrl: data['imgUrl'],
              ),
              builder: (_, s) => s.data ?? Icon(Icons.circle, size: 0),
            ),
            title: Text(data['symbol']),
            subtitle: Builder(
              builder: (_) {
                print(
                    "_cryptoAssetListItemCardTileWidget assetBalancePlaceholder ${assetBalancePlaceholder}");
                bool _showAddress = false;
                return StatefulBuilder(
                    builder: (statefulBuilderCtx, _setState) => GestureDetector(
                          child: Text(data['currency'] +
                              (_showAddress ? "\n${data['address']}" : "")),
                          onLongPressStart: (d) =>
                              _setState(() => _showAddress = !_showAddress),
                          onLongPressEnd: (d) =>
                              _setState(() => _showAddress = !_showAddress),
                        ));
              },
            ),
            trailing: (data['supported'] as bool? ?? false)
                ? common.showBalance
                    ? FutureBuilder(
                        future: common
                            .getERC20BalanceAndInit(
                                common.currentAddress,
                                data['address'],
                                common.topSelectedBlockchainType)
                            .then((value) {
                          if (notifyUpdateCacheAssetBalance != null)
                            notifyUpdateCacheAssetBalance!(
                                data['address'], value);
                          print(
                              "_cryptoAssetListItemCardTileWidget notifyUpdateCacheAssetBalance $data $value");
                          return value;
                        }),
                        builder: (_, s) {
                          final balance;
                          final balanceInUSD;
                          if (s.connectionState != ConnectionState.done ||
                              !s.hasData) {
                            balance = assetBalancePlaceholder ?? '-';
                          } else {
                            balance = s.data;
                          }

                          final balanceWithoutFormat =
                              balance.replaceAll(',', '');
                          final valueWithoutFormat =
                              (double.tryParse(balanceWithoutFormat) ?? 0) *
                                  (coinPrice?.currentPrice ?? 0);
                          balanceInUSD = 'USD ${common.numberFormat(
                            number: valueWithoutFormat.toString(),
                            minDecimal: 2,
                            maxDecimal: 2,
                          )}';

                          return Container(
                            width: 100,
                            child: Column(
                              mainAxisAlignment: MainAxisAlignment.center,
                              crossAxisAlignment: CrossAxisAlignment.end,
                              children: [
                                Text(
                                  balance,
                                  overflow: TextOverflow.ellipsis,
                                  maxLines: 2,
                                  textAlign: TextAlign.right,
                                ),
                                SizedBox(
                                  height: 4,
                                ),
                                if (balanceInUSD != null)
                                  Text(
                                    balanceInUSD,
                                    style: FXUI.subtitleTextStyle,
                                    overflow: TextOverflow.ellipsis,
                                    maxLines: 1,
                                    textAlign: TextAlign.right,
                                  ),
                              ],
                            ),
                          );
                        },
                      )
                    : Text('****.****')
                : Text('COMMON.COMING_SOON'.tr()),
            onTap: onTap,
          )
        ],
      ),
      margin: EdgeInsets.all(4),
      decoration: BoxDecoration(
        shape: BoxShape.rectangle,
        boxShadow: kElevationToShadow[3],
        color: Colors.white,
        borderRadius: FXUI.cricleRadius,
      ),
    );
  }
}
