import 'package:decimal/decimal.dart';
import 'package:easy_localization/easy_localization.dart';
import '../common/commonMethod.dart';
import '../commonUI/crypto_asset_list_item_card_tile_widget.dart';
import '../commonUI/replacing_qrcode_widget.dart';
import '../commonUI/topSelectBlockChainBar.dart';
import 'transferSuccessfulPage.dart';
import '../commonUI/transaction_list_item.dart';
import 'package:apihandler/apiHandler.dart';

class AssetBSCDetailList extends StatefulWidget {
  final Map<String, dynamic> cryptoCurrencyModelMap;

  const AssetBSCDetailList({
    Key? key,
    required this.cryptoCurrencyModelMap,
  }) : super(key: key);

  @override
  _AssetBSCDetailListState createState() => _AssetBSCDetailListState();
}

class _AssetBSCDetailListState extends State<AssetBSCDetailList> {
  List<Map<String, dynamic>> recordList = [];
  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
        future: Future.delayed(Duration(milliseconds: 100), () async {
          common.selectTokenSymbol = widget.cryptoCurrencyModelMap['symbol'];
          bool isBNB = widget.cryptoCurrencyModelMap['symbol'] == 'BNB';
          if (isBNB) {
            return await apiHandler
                .get(
                    "${api.bscTxAPIUrl}api?module=account&action=txlist&address=${common.currentAddress}&startblock=1&endblock=99999999&sort=desc")
                .then((value) => value)
                .catchError((e) => common.coinPriceList);
          } else {
            //Get a list of "BEP-20 - Token Transfer Events" by Address
            String contractaddress =
                widget.cryptoCurrencyModelMap['addressBSC'];
            return await apiHandler
                .get(
                    "${api.bscTxAPIUrl}api?module=account&action=tokentx&contractaddress=$contractaddress&address=${common.currentAddress}&page=1&offset=100&sort=desc")
                .then((value) => value)
                .catchError((e) => common.coinPriceList);
          }
        }),
        builder: (_, AsyncSnapshot result) {
          if (result.connectionState != ConnectionState.done)
            return Center(child: CircularProgressIndicator.adaptive());
          if (result.data != null &&
              result.data['result'] is List &&
              result.data['result']?.length > 0) {
            return transactionListView(result.data['result']);
          } else
            return Container(
              child: Center(
                child: Text("ASSET_TX_HISTORY.NO_TX_RECORD".tr()),
              ),
            );
        });
  }

  Widget transactionListView(List<dynamic> data) {
    return Column(
      children: [
        Container(
          padding: EdgeInsets.only(top: 5, left: 25, right: 5),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'ASSET_TX_HISTORY.RECENT_TX'.tr(),
                style: FXUI.normalTextStyle.copyWith(
                    fontWeight: FontWeight.w700,
                    color: common.getBackGroundColor()),
              ),
              // IconButton(icon: Icon(Icons.download_sharp, color: common.currentUserType == CurrentUserType.decentralized ? FXColor.mainDeepBlueColor : FXColor.mainBlueColor), onPressed: () => common.pushPage(page: DownloadPDFPage(reportType: ReportType.assetTxDetail), context: context))
            ],
          ),
        ),
        Expanded(
          child: ListView.builder(
            padding: EdgeInsets.only(bottom: 35),
            itemCount: data.length,
            itemBuilder: (itemBuilderContext, i) {
              return TransactionListItem(
                itemBuilderContext: itemBuilderContext,
                chain: "BNB",
                eurusTxType: "",
                adminFee: data[i]['cumulativeGasUsed'],
                txFrom: data[i]['from'],
                gasPrice: data[i]['gasPrice'],
                gasUsed: data[i]['gasUsed'],
                confirmTimestamp: '${int.parse(data[i]['timeStamp']) * 1000}',
                decodedInputRecipientAddress: data[i]['to'],
                transactionHash: data[i]['hash'],
                decodedInputAmount:
                    EtherAmount.inWei(BigInt.parse(data[i]['value']))
                        .getValueInUnit(EtherUnit.ether),
                hvGasFee: false,
                symbol: widget.cryptoCurrencyModelMap['symbol'],
                blockChainType: BlockChainType.BinanceCoin,
                toBlockChainType: BlockChainType.BinanceCoin,
                onTap: () => {},
              );
            },
          ),
        )
      ],
    );
  }
}
