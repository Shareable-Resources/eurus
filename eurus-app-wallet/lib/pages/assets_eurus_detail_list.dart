import 'package:easy_localization/easy_localization.dart';
import '../common/commonMethod.dart';
import '../commonUI/transaction_list_item.dart';

class AssetEurusDetailList extends StatefulWidget {
  final Map<String, dynamic> cryptoCurrencyModelMap;

  const AssetEurusDetailList({
    Key? key,
    required this.cryptoCurrencyModelMap,
  }) : super(key: key);

  @override
  _AssetEurusDetailListState createState() => _AssetEurusDetailListState();
}

class _AssetEurusDetailListState extends State<AssetEurusDetailList> {
  late Future<List<Map<String, dynamic>>> _fetchTransactionHistoryFuture;
  List<Map<String, dynamic>> recordList = [];

  @override
  void initState() {
    _fetchTransactionHistoryFuture =
        Future.delayed(Duration(milliseconds: 100), () async {
      common.selectTokenSymbol = widget.cryptoCurrencyModelMap['symbol'];
      return await common.fetchTransactionHistory(
          common.currentAddress ?? '',
          widget.cryptoCurrencyModelMap['address'] ==
                  widget.cryptoCurrencyModelMap['address$TEST_NET']
              ? 0
              : 1,
          widget.cryptoCurrencyModelMap['address']);
    });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
        future: _fetchTransactionHistoryFuture,
        // builder: (_, AsyncSnapshot<List<Map<String, Object>>> s) => s.connectionState != ConnectionState.done ? Center(child: CircularProgressIndicator()) : ListView.builder(itemBuilder: (itemBuilderContext, i) => ListTile(leading: Icon(s.data[i]['from'] == common.address ? Icons.arrow_upward : Icons.arrow_downward), title: Text("${s.data[i]['peerAddress']}"), subtitle: Text.rich(TextSpan(text: "${s.data[i]['transactionDateTime']}\nTx ID: ${s.data[i]['transactionHash']}"), textWidthBasis: TextWidthBasis.longestLine, textScaleFactor: .85,), trailing: Text("${s.data[i]['from'] == common.address ? '-' : ''} ${s.data[i]['decodedInputAmount']}", style: TextStyle(color: s.data[i]['from'] != common.address ? Colors.green : Colors.black87),)), itemCount: s.data.length),)),
        builder: (_, AsyncSnapshot<List<Map<String, dynamic>>> s) => s
                    .connectionState !=
                ConnectionState.done
            ? Center(child: CircularProgressIndicator.adaptive())
            : s.hasError
                ? Center(child: Text("ASSET_TX_HISTORY.NO_TX_RECORD".tr()))
                : !s.hasData || s.data is List && s.data!.length == 0
                    // ? cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}']
                    //   ? Column(children: [])
                    //   : Align(alignment: Alignment.topCenter, child: Column(children: []))
                    ? Container(
                        child: Center(
                          child: Text("ASSET_TX_HISTORY.NO_TX_RECORD".tr()),
                        ),
                      )
                    : transactionListView(s.data!));
  }

  Widget transactionListView(List<Map<String, dynamic>> data) {
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
              /// Mainnet
              ///   Send:           ✓ Gas Fee, ✗ Admin Fee
              ///   Receive:        ✗ Gas Fee, ✗ Admin Fee
              ///   to Sidechain:   ✓ Gas Fee, ✗ Admin Fee
              /// Sidechain
              ///   Send:           ✓ Gas Fee, ✗ Admin Fee
              ///   Receive:        ✗ Gas Fee, ✗ Admin Fee
              ///   to Mainnet:     ✓ Gas Fee, ✓ Admin Fee
              final chain = data[i]['chain'];
              final adminFee = data[i]['adminFee'];
              final eurusTxType = data[i]['eurusTxType'];
              final isSend = data[i]['isSend'] as bool?;
              bool hvGasFee = eurusTxType == '1'
                  ? isSend ?? false
                  : eurusTxType == '2'
                      ? chain == 'eth'
                      : eurusTxType == '3' || eurusTxType == '7'
                          ? chain == 'eun'
                          : false;

              recordList = data;
              return TransactionListItem(
                itemBuilderContext: itemBuilderContext,
                chain: chain,
                eurusTxType: eurusTxType,
                adminFee: adminFee,
                txFrom: data[i]['txFrom'],
                gasPrice: data[i]['gasPrice'],
                gasUsed: data[i]['gasUsed'],
                confirmTimestamp: data[i]['confirmTimestamp'],
                decodedInputRecipientAddress: data[i]
                    ['decodedInputRecipientAddress'],
                transactionHash: data[i]['transactionHash'],
                decodedInputAmount: eurusTxType == '7'
                    ? data[i]['transferGas'] *
                        web3dart
                            .getGasPrice(blockChainType: BlockChainType.Eurus)
                            .getInWei
                            .toDouble() /
                        pow(10, 18)
                    : data[i]['decodedInputAmount'],
                hvAdminFee: chain == 'eun' &&
                    eurusTxType != null &&
                    eurusTxType == '3' &&
                    adminFee != null,
                hvGasFee: hvGasFee,
                isAllocate: eurusTxType != null && eurusTxType != '1',
                symbol: widget.cryptoCurrencyModelMap['symbol'],
                blockChainType: chain == 'eun'
                    ? (eurusTxType == '2'
                        ? BlockChainType.Ethereum
                        : BlockChainType.Eurus)
                    : (eurusTxType == '3'
                        ? BlockChainType.Eurus
                        : BlockChainType.Ethereum),
                toBlockChainType: chain == 'eun'
                    ? (eurusTxType == '3'
                        ? BlockChainType.Ethereum
                        : BlockChainType.Eurus)
                    : (eurusTxType == '2'
                        ? BlockChainType.Eurus
                        : BlockChainType.Ethereum),
                onTap: () =>
                    {common.eurusTXType = recordList[i]["eurusTxType"]},
              );
            },
          ),
        )
      ],
    );
  }
}
