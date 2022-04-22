import 'package:easy_localization/easy_localization.dart';
import '../common/commonMethod.dart';
import '../pages/assets_detail_page.dart';
import 'package:decimal/decimal.dart';

class TransactionListItem extends StatelessWidget {
  const TransactionListItem(
      {Key? key,
      this.chain,
      this.eurusTxType,
      this.adminFee,
      this.txFrom,
      this.gasPrice,
      this.gasUsed,
      this.confirmTimestamp,
      this.decodedInputRecipientAddress,
      this.transactionHash,
      this.decodedInputAmount,
      this.hvAdminFee: false,
      this.hvGasFee: false,
      this.isAllocate: false,
      this.blockChainType,
      this.toBlockChainType,
      this.itemBuilderContext,
      this.symbol,
      this.onTap})
      : super(key: key);

  final chain;
  final eurusTxType;
  final adminFee;
  final txFrom;
  final gasPrice;
  final gasUsed;
  final confirmTimestamp;
  final decodedInputRecipientAddress;
  final transactionHash;
  final decodedInputAmount;
  final bool hvAdminFee;
  final bool hvGasFee;
  final bool isAllocate;
  final blockChainType;
  final toBlockChainType;
  final itemBuilderContext;
  final symbol;
  final Function? onTap;

  @override
  Widget build(BuildContext context) {
    final calculatedGasFee = gasPrice != null &&
            gasUsed != null &&
            gasUsed != "null" &&
            gasPrice != "null" &&
            gasPrice != "0" &&
            gasUsed != "0" &&
            gasUsed is String
        ? '${((BigInt.tryParse(gasPrice) ?? BigInt.from(0)) * (BigInt.tryParse(gasUsed) ?? BigInt.from(0)) / BigInt.from(pow(10, 9)) / pow(10, 9)).toStringAsFixed(8)}'
        : null;
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 20),
      // color:
      //     i % 2 > 0 ? FXColor.lightWhiteColor.withOpacity(0.7) : Colors.white,
      child: GestureDetector(
        onTap: () async {
          if (onTap != null) onTap!();
          String gasFee = (isEmptyString(string: gasPrice) ||
                  isEmptyString(string: gasUsed))
              ? ""
              : "${(((BigInt.tryParse(gasPrice ?? '') ?? BigInt.zero) * (BigInt.tryParse(gasUsed ?? '') ?? BigInt.zero)) / BigInt.from(pow(10, 9)) / pow(10, 9)).toStringAsFixed(10).trim()} ${chain.toUpperCase()} ${'TX_PAGE.GAS_FEE'.tr()}";
          await canNavigateToTransactionHistoryDetailPage(
            // if not withdrawal but normal transfer
            // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
            centerMessage: (confirmTimestamp != null &&
                    confirmTimestamp is String &&
                    confirmTimestamp != ''
                ? isAllocate
                    ? 'TX_PAGE.TRANSFER_SUCCESS_DESC'.tr()
                    : ''
                : 'COMMON.PENDING'.tr()),
            // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
            date:
                "${confirmTimestamp != null && confirmTimestamp is String && confirmTimestamp != '' ? DateFormat('dd-MM-yyyy HH:mm:ss').format(DateTime.fromMillisecondsSinceEpoch(int.tryParse(confirmTimestamp) ?? 0)) : 'Pending Transaction'}",
            fromAddress: "$txFrom",
            // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
            // gasFeeString: "... ${cryptoCurrencyModelMap['address'] == cryptoCurrencyModelMap['address${TEST_NET??'Ethereum'}'] ? 'ETH' : 'EUN'} Gas Fee",
            gasFeeString: hvGasFee ? gasFee : '',
            adminFeeString: hvAdminFee ? adminFee : '',
            toAddress: "$decodedInputRecipientAddress",
            // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
            transferAmount:
                "${"$txFrom".toLowerCase() == common.currentAddress?.toLowerCase() ? '' : ''} ${Decimal.parse(decodedInputAmount.toString())} $symbol",
            txId: "$transactionHash",
            navigatorContext: itemBuilderContext,
            // TODO: if not withdrawal but normal transfer
            isAssetAllocation: isAllocate,
            // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
            // TODO: TxHistoryRecordItemViewModel e.g. isSuccess
            shouldSkipPendingFetch:
                "$txFrom".toLowerCase() == common.currentAddress?.toLowerCase()
                    ? confirmTimestamp != null &&
                        confirmTimestamp is String &&
                        confirmTimestamp != ''
                    : true,
            blockChainType: blockChainType,
            toBlockChainType: toBlockChainType,
            eurusTxType: eurusTxType,
          );
        },
        child: Container(
          padding: EdgeInsets.symmetric(vertical: 16),
          decoration: BoxDecoration(
            border: Border(
              bottom: BorderSide(
                width: 1,
                color: FXColor.lightBlackColor.withOpacity(0.2),
              ),
            ),
          ),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
              Container(
                margin: EdgeInsets.only(right: 12),
                width: 21,
                child: isAllocate
                    ? Image.asset(
                        eurusTxType == '7'
                            ? 'images/icn_fuel.png'
                            : eurusTxType == '6'
                                ? 'images/icn_reward2.png'
                                : 'images/icnAllocate.png',
                        package: 'euruswallet',
                      )
                    : "$txFrom".toLowerCase() ==
                            common.currentAddress?.toLowerCase()
                        ? Image.asset('images/icnSend.png',
                            package: 'euruswallet')
                        : Image.asset('images/icnReceive.png',
                            package: 'euruswallet'),
              ),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Text(
                      isAllocate
                          ? isCentralized()
                              ? eurusTxType == '2'
                                  ? 'ASSET_TX_HISTORY.DEPOSIT'.tr()
                                  : eurusTxType == '6'
                                      ? 'ASSET_TX_HISTORY.REWARD'.tr()
                                      : eurusTxType == '7'
                                          ? 'REFUEL.TITLE'.tr()
                                          : 'ASSET_TX_HISTORY.WITHDRAWAL'.tr()
                              : 'ASSET_TX_HISTORY.ALLOCATE_TO'.tr(args: [
                                  eurusTxType == '2' ? 'Ethereum' : 'Eurus',
                                  eurusTxType == '2' ? 'Eurus' : 'Ethereum'
                                ])
                          : "$txFrom".toLowerCase() ==
                                  common.currentAddress?.toLowerCase()
                              ? "$decodedInputRecipientAddress"
                              : "$txFrom",
                      style: FXUI.normalTextStyle.copyWith(
                          fontWeight: FontWeight.w500,
                          color: FXColor.textGray,
                          fontSize: 12),
                    ),
                    Padding(
                      padding: EdgeInsets.only(top: 10, bottom: 3),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text('ASSET_TX_HISTORY.TX_ID'.tr(),
                              style: FXUI.normalTextStyle.copyWith(
                                  fontWeight: FontWeight.w700, fontSize: 12)),
                          Text(
                              "${confirmTimestamp != null && confirmTimestamp is String && confirmTimestamp != '' ? DateFormat('dd-MM-yyyy HH:mm:ss').format(DateTime.fromMillisecondsSinceEpoch(int.tryParse(confirmTimestamp) ?? 0)) : 'Pending Transaction'}",
                              style:
                                  FXUI.normalTextStyle.copyWith(fontSize: 10))
                        ],
                      ),
                    ),
                    Text.rich(
                      TextSpan(text: "$transactionHash"),
                      textWidthBasis: TextWidthBasis.longestLine,
                      textScaleFactor: .85,
                      style: FXUI.normalTextStyle.copyWith(
                        color: FXColor.textGray,
                      ),
                    ),
                    Padding(
                      padding: EdgeInsets.only(top: 15),
                    ),
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.stretch,
                            mainAxisAlignment: MainAxisAlignment.start,
                            children: [
                              Text(
                                'ASSET_TX_HISTORY.AMOUNT'.tr(),
                                style: FXUI.normalTextStyle.copyWith(
                                  fontWeight: FontWeight.w700,
                                  fontSize: 12,
                                  color: FXColor.purpleColor,
                                ),
                              ),
                              Padding(
                                padding: EdgeInsets.symmetric(vertical: 2),
                                child: Text(
                                  '${Decimal.parse(decodedInputAmount.toString())}',
                                  style: FXUI.normalTextStyle.copyWith(
                                    fontSize: 12,
                                    color: FXColor.textGray,
                                  ),
                                ),
                              ),
                              Text(
                                symbol,
                                style: FXUI.normalTextStyle.copyWith(
                                  fontSize: 12,
                                  color: FXColor.textGray,
                                ),
                              ),
                            ],
                          ),
                        ),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.stretch,
                            mainAxisAlignment: MainAxisAlignment.start,
                            children: [
                              Text(
                                'TX_PAGE.GAS_FEE'.tr(),
                                style: FXUI.normalTextStyle.copyWith(
                                  fontWeight: FontWeight.w700,
                                  fontSize: 12,
                                  color: FXColor.purpleColor,
                                ),
                              ),
                              Padding(
                                padding: EdgeInsets.symmetric(vertical: 2),
                                child: Text(
                                  '${hvGasFee && !isEmptyString(string: calculatedGasFee) ? calculatedGasFee : '-'}',
                                  style: FXUI.normalTextStyle.copyWith(
                                    fontSize: 12,
                                    color: FXColor.textGray,
                                  ),
                                ),
                              ),
                              Text(
                                '${hvGasFee && !isEmptyString(string: calculatedGasFee) ? chain.toUpperCase() : '-'}',
                                style: FXUI.normalTextStyle.copyWith(
                                  fontSize: 12,
                                  color: FXColor.textGray,
                                ),
                              )
                            ],
                          ),
                        ),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.stretch,
                            mainAxisAlignment: MainAxisAlignment.start,
                            children: [
                              Text(
                                "ADMIN_FEE.NAME".tr(),
                                style: FXUI.normalTextStyle.copyWith(
                                  fontWeight: FontWeight.w700,
                                  fontSize: 12,
                                  color: FXColor.purpleColor,
                                ),
                              ),
                              Padding(
                                padding: EdgeInsets.symmetric(vertical: 2),
                                child: Text(
                                  '${!hvAdminFee ? '-' : common.numberFormat(number: adminFee)}',
                                  style: FXUI.normalTextStyle.copyWith(
                                    fontSize: 12,
                                    color: FXColor.textGray,
                                  ),
                                ),
                              ),
                              Text(
                                '${!hvAdminFee ? '-' : symbol}',
                                style: FXUI.normalTextStyle.copyWith(
                                  fontSize: 12,
                                  color: FXColor.textGray,
                                ),
                              )
                            ],
                          ),
                        ),
                        Icon(Icons.arrow_forward_ios_rounded,
                            color: FXColor.textGray, size: 16),
                      ],
                    ),
                    // TODO: TxHistoryRecordItemViewModel e.g. isSender, isOutbound, isInbound etc
                    // Column( crossAxisAlignment: CrossAxisAlignment.end, mainAxisAlignment: MainAxisAlignment.end, children: [
                    //   Text("${"${txFrom}".toLowerCase() == common.address?.toLowerCase() ? '-' : ''} ${s.data[i]['decodedInputAmount']} ${cryptoCurrencyModelMap['symbol']}", style: TextStyle(color: "${txFrom}".toLowerCase() == common.address?.toLowerCase() ? Color(0xffbd1d00) : Color(0xff00bd6d)),),
                    //   "${txFrom}".toLowerCase() == common.address?.toLowerCase() ? Text(gasPrice != null && gasUsed != null ? '${(BigInt.tryParse(gasPrice) * BigInt.tryParse(gasUsed) / BigInt.from(pow(10, 9)) / pow(10,9)).toStringAsFixed(10)} ${chain == 'eth' ? 'ETH' : 'EUN'}' : '', style: TextStyle(color: Color(0xff00B3F3)),) : Container(),
                    // ],),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
