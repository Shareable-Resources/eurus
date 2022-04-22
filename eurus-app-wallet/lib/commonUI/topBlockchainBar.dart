import 'package:euruswallet/common/commonMethod.dart';

class TopBlockchainBar extends StatelessWidget {
  final TransactionType transactionType;
  BlockChainType? fromBlockChainType;
  BlockChainType? toBlockChainType;
  TopBlockchainBar({
    this.transactionType: TransactionType.successful,
    this.fromBlockChainType,
    this.toBlockChainType,
  });

  String fromBlockchain() {
    String fromBlockchain;
    fromBlockChainType = fromBlockChainType == null
        ? common.fromBlockChainType
        : fromBlockChainType;
    if (common.isTransactionHistory(transactionType: transactionType)) {
      fromBlockchain = getBlockChainName(fromBlockChainType!).toUpperCase();
    } else {
      fromBlockchain = getBlockChainName(fromBlockChainType!).toUpperCase();
    }
    return fromBlockchain;
  }

  String toBlockchain() {
    String toBlockchain;
    toBlockChainType =
        toBlockChainType == null ? common.fromBlockChainType : toBlockChainType;
    if (common.isTransactionHistory(transactionType: transactionType)) {
      toBlockchain = getBlockChainName(toBlockChainType!).toUpperCase();
    } else {
      if (common.transferToMySelf == true) {
        toBlockchain = common.fromBlockChainType == BlockChainType.Ethereum
            ? "EURUS"
            : "ETHEREUM";
      } else {
        toBlockchain = getBlockChainName(toBlockChainType!).toUpperCase();
      }
    }
    return toBlockchain;
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(fromBlockchain(),
              style: FXUI.titleTextStyle
                  .copyWith(color: FXColor.lightGray, fontSize: 15)),
          Padding(
            padding: EdgeInsets.only(left: 10, right: 10),
            child: SizedBox(
              width: 16,
              height: 16,
              child: Image.asset(
                'images/arrowRight.png',
                package: 'euruswallet',
                color: FXColor.lightGray,
              ),
            ),
          ),
          Text(toBlockchain(),
              style: FXUI.titleTextStyle
                  .copyWith(color: FXColor.lightGray, fontSize: 15)),
        ],
      ),
    );
  }
}
