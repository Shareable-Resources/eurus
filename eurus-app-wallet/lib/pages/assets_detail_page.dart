import 'package:collection/src/iterable_extensions.dart';
import 'package:decimal/decimal.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/model/coinPrice.dart';

import '../common/commonMethod.dart';
import '../commonUI/crypto_asset_list_item_card_tile_widget.dart';
import '../commonUI/replacing_qrcode_widget.dart';
import '../commonUI/topSelectBlockChainBar.dart';
import 'assets_bsc_detail_list.dart';
import 'assets_eurus_detail_list.dart';
import 'transferSuccessfulPage.dart';

class AssetsDetailPage extends StatefulWidget {
  final Map<String, dynamic> cryptoCurrencyModelMap;
  final Map<String, dynamic> cArgs;

  const AssetsDetailPage(
      {Key? key, required this.cryptoCurrencyModelMap, required this.cArgs})
      : super(key: key);

  @override
  _AssetsDetailPageState createState() => _AssetsDetailPageState();
}

class _AssetsDetailPageState extends State<AssetsDetailPage> {
  late Map<String, dynamic> _cArgs;
  List<Map<String, dynamic>> recordList = [];
  @override
  void initState() {
    super.initState();

    final address = widget.cryptoCurrencyModelMap['address'];
    final addressWithTestNet =
        widget.cryptoCurrencyModelMap['address$TEST_NET'];
    _cArgs = {
      ...widget.cArgs,
      ...{
        '${address == addressWithTestNet ? "ethereum" : "eurus"}Erc20ContractAddress':
            address ?? '0x0' ?? '' ?? null,
        // 'fromBlockChainType': address == addressWithTestNet ? 0 : 1,
        'fromBlockChainType': common.topSelectedBlockchainType.index,
        // 'currencyName': cryptoCurrencyModelMap['symbol']??'ETH'??'EUN',
        'disableSelectBlockchain': false,
        'ethereumAddress': common.currentAddress,
        // 'canGetPrivateKeyHandler': _getPrivateKeyHandler,
        // 'navigateToAssetAllocationPage': () {
        //   Navigator.of(context).push
        // },

        // for contextual navigated-scene-aware asset allocation button in function bar item to be shown or not
        '${address == addressWithTestNet ? 'eurus' : 'ethereum'}Erc20ContractAddress':
            widget.cryptoCurrencyModelMap[
                'address${address == addressWithTestNet ? 'Eurus' : TEST_NET}'],
        'symbol': widget.cryptoCurrencyModelMap["symbol"]
      },
    };
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Builder(
        builder: (scaffoldBuilderContext) => Container(
          decoration: BoxDecoration(
            image: DecorationImage(
              image: AssetImage(
                  !isCentralized()
                      ? 'images/bgDecentralized.png'
                      : 'images/bgCentralized.png',
                  package: 'euruswallet'),
              fit: BoxFit.cover,
              alignment: Alignment.topCenter,
            ),
          ),
          child: Center(
            child: Column(
              children: [
                AppBar(
                    centerTitle: true,
                    brightness: Brightness.light,
                    foregroundColor: Colors.black,
                    iconTheme: IconThemeData(color: Colors.white),
                    backgroundColor: Colors.transparent,
                    elevation: 0,
                    title: Text('${widget.cryptoCurrencyModelMap['symbol']}',
                        style: Theme.of(context)
                            .textTheme
                            .headline6
                            ?.copyWith(color: Colors.white))),
                // ListTile(title: Card(child: Text(),))
                // Container(child: ListTile(leading: Icon(Icons.circle, size: 40), title: Text('USDT'), subtitle: Text('Tether USD'), trailing: Text("9876.54321"),), margin: EdgeInsets.all(4), decoration: BoxDecoration(shape: BoxShape.rectangle, boxShadow: kElevationToShadow[3], color: Colors.white, borderRadius: BorderRadius.circular(15),)),
                Padding(
                  padding: EdgeInsets.symmetric(horizontal: 8),
                  child: CryptoAssetListItemCardTileWidget(
                    data: widget.cryptoCurrencyModelMap,
                    chainIndicator: Container(
                      padding: EdgeInsets.only(top: 15),
                      child: TopSelectBlockChainBar(
                        disableSelectBlockchain: true,
                        currentSelection: common.currentBlockchainSelection,
                        // widget.cryptoCurrencyModelMap['address'] ==
                        //         widget.cryptoCurrencyModelMap[
                        //             'address$TEST_NET']
                        //     ? BlockChainType.Ethereum
                        //     : BlockChainType.Eurus
                      ),
                    ),
                    coinPrice: common.coinPriceList?.coinPriceList
                        .firstWhereOrNull((element) =>
                            element.symbol.toUpperCase() ==
                            widget.cryptoCurrencyModelMap['symbol']),
                  ),
                ),
                Padding(
                  padding: EdgeInsets.only(bottom: 13),
                  child: IntrinsicHeight(
                    child: Row(
                      children: common
                          .homeFunctionsBarItemBuilder(scaffoldBuilderContext, {
                            ..._cArgs,
                            ...{
                              if (common.currentUserProfile
                                          ?.seedPraseBackuped !=
                                      true &&
                                  common.currentUserType ==
                                      CurrentUserType.decentralized)
                                'replacingQRCodeWidget':
                                    ReplacingQRCodeWidget(),
                              'navigateToAssetAllocationPage':
                                  (Future<Null> Function(BuildContext,
                                                  Map<String, dynamic>)
                                              __navigateToAssetAllocationTransfer,
                                          [_barBuilderContext]) async =>
                                      await __navigateToAssetAllocationTransfer(
                                          _barBuilderContext,
                                          widget.cryptoCurrencyModelMap)
                            },
                            ...{
                              'btnColor': Colors.white,
                              'enableAllocation':
                                  (widget.cryptoCurrencyModelMap['addressEurus']
                                          is String &&
                                      (widget.cryptoCurrencyModelMap[
                                              'addressEurus'] as String)
                                          .isNotEmpty)
                            }
                          })
                          .map((e) => Expanded(child: e))
                          .toList(),
                      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                    ),
                  ),
                ),
                Expanded(
                  child: Container(
                    decoration: BoxDecoration(
                        color: Colors.white, borderRadius: FXUI.cricleRadius),
                    child: genAssetsList(),
                  ),
                ),
                // itemBuilder: (itemBuilderContext, i) => ListTile(
                //   onTap: () async {await widget.canNavigateToTransactionHistoryDetailPage(
                //     // if not withdrawal but normal transfer
                //     // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
                //     centerMessage: '' ?? (confirmTimestamp != null && confirmTimestamp is String && confirmTimestamp != '' ? 'Success' : 'Pending'),
                //     // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
                //     date: "${confirmTimestamp != null && confirmTimestamp is String && confirmTimestamp != '' ? DateFormat.yMd('en_HK').add_jms().format(DateTime.fromMillisecondsSinceEpoch(int.tryParse(confirmTimestamp)))+' UTC+8' : 'Pending Transaction'}",
                //     fromAddress: "${txFrom}",
                //     // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                //     // gasFeeString: "... ${cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'] ? 'ETH' : 'EUN'} Gas Fee",
                //     gasFeeString: "${txFrom}".toLowerCase() == common.address?.toLowerCase() ? "... ${chain == 'eth' ? 'ETH' : 'EUN'} Gas Fee" : '',
                //     toAddress: "${s.data[i]['decodedInputRecipientAddress']}",
                //     // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                //     transferAmount: "${"${txFrom}".toLowerCase() == common.address?.toLowerCase() ? '-' : ''} ${s.data[i]['decodedInputAmount']} ${cryptoCurrencyModelMap['symbol']}",
                //     txId: "${s.data[i]['transactionHash']}",
                //     navigatorContext: itemBuilderContext,
                //     // TODO: if not withdrawal but normal transfer
                //     isAssetAllocation: null??false,
                //     // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                //     // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
                //     shouldSkipPendingFetch: "${txFrom}".toLowerCase() == common.address?.toLowerCase() ? confirmTimestamp != null && confirmTimestamp is String && confirmTimestamp != '' : true,
                //     blockChainType: chain == 'eth' ? 0 : 1,
                //   );},
                //   // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                //   // leading: Icon("${txFrom}".toLowerCase() == common.address?.toLowerCase() ? Icons.arrow_upward : Icons.arrow_downward),
                //   leading: "${txFrom}".toLowerCase() == common.address?.toLowerCase() ? Image.asset('images/icnSend.png', package: 'euruswallet') : Image.asset('images/icnReceive.png', package: 'euruswallet'),
                //   title: Text("${txFrom}".toLowerCase() == common.address?.toLowerCase() ? "${s.data[i]['decodedInputRecipientAddress']}" : "${txFrom}"),
                //   // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
                //   subtitle: Text.rich(TextSpan(text: "${confirmTimestamp != null && confirmTimestamp is String && confirmTimestamp != '' ? DateFormat.yMd('en_HK').add_jms().format(DateTime.fromMillisecondsSinceEpoch(int.tryParse(confirmTimestamp)))+' UTC+8' : 'Pending Transaction'}\nTx ID: ${s.data[i]['transactionHash']}"), textWidthBasis: TextWidthBasis.longestLine, textScaleFactor: .85,),
                //   // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                //   trailing: Text("${"${txFrom}".toLowerCase() == common.address?.toLowerCase() ? '-' : ''} ${s.data[i]['decodedInputAmount']}", style: TextStyle(color: "${txFrom}".toLowerCase() == common.address?.toLowerCase() ? Colors.black87 : Colors.green),)
                // ))))),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget genAssetsList() {
    if (common.topSelectedBlockchainType == BlockChainType.Eurus ||
        common.topSelectedBlockchainType == BlockChainType.Ethereum) {
      return AssetEurusDetailList(
          cryptoCurrencyModelMap: widget.cryptoCurrencyModelMap);
    } else if (common.topSelectedBlockchainType == BlockChainType.BinanceCoin) {
      return AssetBSCDetailList(
          cryptoCurrencyModelMap: widget.cryptoCurrencyModelMap);
    } else
      return Container();
  }
}

Future<void> canNavigateToTransactionHistoryDetailPage({
  String centerMessage = '',
  String fromAddress = '',
  String toAddress = '',
  String txId = '',
  String date = '',
  String gasFeeString = '',
  String? adminFeeString,
  String transferAmount = '',
  required BuildContext navigatorContext,
  bool isAssetAllocation = false,
  bool shouldSkipPendingFetch = true,
  BlockChainType blockChainType = BlockChainType.Ethereum,
  BlockChainType? toBlockChainType,
  String? eurusTxType,
}) async {
  await common.pushPage(
    page: TransferSuccessfulPage(
      transactionType: isAssetAllocation
          ? shouldSkipPendingFetch
              ? TransactionType.allocationTransactionHistorySuccessfulStatus
              : TransactionType.allocationTransactionHistoryProcessingStatus
          : shouldSkipPendingFetch
              ? TransactionType.sendTransactionHistorySuccessfulStatus
              : TransactionType.sendTransactionHistoryPendingStatus,
      centerMessage: centerMessage,
      fromAddress: fromAddress,
      toAddress: toAddress,
      txId: txId,
      date: date,
      gasFeeString: gasFeeString,
      adminFeeString: adminFeeString,
      transferAmount: transferAmount,
      fromBlockChainType: blockChainType,
      toBlockChainType: toBlockChainType ?? blockChainType,
      eurusTxType: eurusTxType,
    ),
    context: navigatorContext,
  );
}
