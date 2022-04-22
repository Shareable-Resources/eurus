import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/crypto_currency_model.dart';
import 'package:easy_localization/easy_localization.dart';

class SelectCurrencyPage extends StatefulWidget {
  final String titleName;
  final String? userWalletAccountAvailableAssetsList;

  SelectCurrencyPage({
    Key? key,
    required this.titleName,
    this.userWalletAccountAvailableAssetsList,
  }) : super(key: key);

  @override
  _SelectCurrencyPageState createState() => _SelectCurrencyPageState();
}

class _SelectCurrencyPageState extends State<SelectCurrencyPage> {
  List<Widget>? widgetList = [];

  bool get _isEth => common.fromBlockChainType == BlockChainType.Ethereum;
  bool get _isEur => common.fromBlockChainType == BlockChainType.Eurus;
  Map addressBalanceMap = Map<String, String>();

  @override
  void initState() {
    super.initState();
    // getERC20TokenList();
    _getWidgetList();
  }

  @override
  void dispose() {
    super.dispose();
  }

  Widget currencyRowUI({
    required String currencyName,
    String? currencyAmount,
    String? contractAddress,
  }) {
    return InkWell(
      onTap: () {
        common.selectTokenSymbol = currencyName;
        if (common.isEthOrEun()) {
        } else if (contractAddress != null) {
          // String contractAddress =
          //     web3dart.tokenListMap[common.selectTokenSymbol];
          web3dart.setErc20Contract(
            contractAddress: contractAddress,
            blockChainType: common.fromBlockChainType,
          );
        }
        if ((common.selectTokenSymbol == "EUN") && common.isCenWithdraw) {
          common.showPopUpError(
              context: context,
              descriptions: "COMMON_ERROR.EUN_WITHDRAW_ERROR".tr());
        } else {
          Navigator.pop(context);
        }
      },
      child: Container(
        height: 54,
        width: size.blockSizeHorizontal * 100,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(currencyName, style: FXUI.normalTextStyle),
                Row(
                  children: [
                    Text(
                      isEmptyString(string: contractAddress)
                          ? common.numberFormat(
                              maxDecimal: 12,
                              number: currencyAmount,
                              minDecimal: 8)
                          : double.tryParse(
                                      addressBalanceMap[contractAddress]) !=
                                  null
                              ? common.numberFormat(
                                  maxDecimal: 12,
                                  number: addressBalanceMap[contractAddress],
                                  minDecimal: 8)
                              : "-",
                      style: FXUI.normalTextStyle.copyWith(
                        color: FXColor.lightGray,
                      ),
                    ),
                    Padding(
                      padding: EdgeInsets.only(left: 10),
                      child: Image.asset(
                        "images/rightArrow.png",
                        package: 'euruswallet',
                        width: 10,
                        height: 16,
                        fit: BoxFit.cover,
                      ),
                    ),
                  ],
                )
              ],
            ),
            Divider()
          ], //
        ),
      ),
    );
  }

  // Future<List<Widget>> getWidgetList({List<dynamic> tokenList}) async {
  //   List<Widget> widgetList = [];
  //
  //   widgetList.add(Padding(
  //     padding: EdgeInsets.only(top: 20),
  //     child:
  //         Text("Choose currency", style: TextStyle(color: FXColor.lightGray)),
  //   ));
  //
  //   if (tokenList != null && tokenList[0] != null) {
  //     widgetList.add(currencyRowUI(
  //         currencyName: common.fromBlockChainType == BlockChainType.Ethereum
  //             ? "ETH"
  //             : "EUN",
  //         currencyAmount: common.fromBlockChainType == BlockChainType.Ethereum
  //             ? web3dart.ethBalanceFromEthereum ?? ""
  //             : web3dart.ethBalanceFromEurus ?? ""));
  //     for (var i = 0; i < tokenList[0].length; i++) {
  //       String tokenName = tokenList[0][i];
  //       EthereumAddress tokenAddress = tokenList[1][i];
  //       print("tokenName:$tokenName");
  //       print("tokenAddress:$tokenAddress");
  //       DeployedContract deployedContract = web3dart.getEurusERC20Contract(
  //           contractAddress: tokenAddress.toString());
  //       String balance = await web3dart.getERC20Balance(
  //           blockChainType: common.fromBlockChainType,
  //           deployedContract: deployedContract);
  //       widgetList.add(
  //           currencyRowUI(currencyName: tokenName, currencyAmount: balance));
  //     }
  //   }
  //
  //   return widgetList;
  // }

  // Future<void> loadBalance() async {
  //   if (widget.userWalletAccountAvailableAssetsList != null) {
  //     List aJsonList = jsonDecode(widget.userWalletAccountAvailableAssetsList);
  //     for (var i = 0; i < aJsonList.length; i++) {
  //       CryptoCurrencyModel token = CryptoCurrencyModel.fromJson(aJsonList[i]);
  //
  //       bool shownInEth =
  //           common.fromBlockChainType == BlockChainType.Ethereum &&
  //               token.addressEthereumViewState == TokenViewState.NORMAL;
  //       bool shownInEur = common.fromBlockChainType == BlockChainType.Eurus &&
  //           token.addressEurusViewState == TokenViewState.NORMAL;
  //
  //       if (((shownInEth && _isEth) || (shownInEur && _isEur)) &&
  //           token.showAssets == true) {
  //         String tokenAddress =
  //             common.fromBlockChainType == BlockChainType.Ethereum
  //                 ? token.addressEthereum
  //                 : token.addressEurus;
  //
  //         DeployedContract deployedContract = web3dart.getEurusERC20Contract(
  //           contractAddress: tokenAddress,
  //         );
  //
  //         String balance = await web3dart.getERC20Balance(
  //           blockChainType: common.fromBlockChainType,
  //           deployedContract: deployedContract,
  //         );
  //         setState(() {
  //           addressBalanceMap.update(tokenAddress, (value) => balance);
  //         });
  //       }
  //     }
  //   }
  // }

  Future<bool> _getWidgetList() async {
    // web3dart.tokenList =
    //     await web3dart.getERC20TokenList(blockChainType: common.fromBlockChainType);
    List<Widget> tokenList = [];

    if (widget.userWalletAccountAvailableAssetsList != null) {
      if (!common.isCenWithdraw) {
        tokenList.add(currencyRowUI(
            currencyName: getSymbolByBlockChainType(common.fromBlockChainType),
            currencyAmount: await web3dart.getETHBalance(
                blockChainType: common.fromBlockChainType)));
      }

      List aJsonList = jsonDecode(widget.userWalletAccountAvailableAssetsList!);
      for (var i = 0; i < aJsonList.length; i++) {
        CryptoCurrencyModel token = CryptoCurrencyModel.fromJson(aJsonList[i]);

        bool shownInEth =
            common.fromBlockChainType == BlockChainType.Ethereum &&
                token.addressEthereumViewState == TokenViewState.NORMAL;
        bool shownInEur = common.fromBlockChainType == BlockChainType.Eurus &&
            token.addressEurusViewState == TokenViewState.NORMAL;
        bool shownInBsc =
            common.fromBlockChainType == BlockChainType.BinanceCoin &&
                token.addressBSCViewState == TokenViewState.NORMAL;

        if (((shownInEth && _isEth) || (shownInEur && _isEur) || shownInBsc) &&
            token.showAssets == true) {
          String tokenName = token.symbol;

          String? tokenAddress = "";
          if (common.fromBlockChainType == BlockChainType.Ethereum) {
            tokenAddress = token.addressEthereum;
          } else if (common.fromBlockChainType == BlockChainType.BinanceCoin) {
            tokenAddress = token.addressBSC;
          } else {
            tokenAddress = token.addressEurus;
          }

          if (tokenAddress == null) return false;

          DeployedContract deployedContract = web3dart.getEurusERC20Contract(
            contractAddress: tokenAddress,
          );

          //String balance = "-";
          addressBalanceMap[tokenAddress] = '-';
          Widget currencyUI = currencyRowUI(
            currencyName: tokenName,
            contractAddress: tokenAddress,
          );
          tokenList.add(currencyUI);
          Future.delayed(const Duration(milliseconds: 200), () async {
            String? balance = await web3dart.getERC20Balance(
                blockChainType: common.fromBlockChainType,
                deployedContract: deployedContract);
            setState(() {
              addressBalanceMap[tokenAddress] = balance;
              Widget newCurrencyUI = currencyRowUI(
                currencyName: tokenName,
                contractAddress: tokenAddress,
              );
              int? index = widgetList?.indexOf(currencyUI);
              if (index != null)
                widgetList?.replaceRange(index, index + 1, [newCurrencyUI]);
            });
          });
        }
        // setState(() {
        //   widgetList = tokenList;
        // });

      }
    }

    setState(() {
      widgetList = tokenList;
    });
    // loadBalance();
    return true;
  }

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    return BackGroundImage(
      child: Scaffold(
        backgroundColor: Colors.transparent,
        appBar: WalletAppBar(title: widget.titleName),
        body: TopCircularContainer(
          height: size.heightWithoutAppBar,
          width: size.blockSizeHorizontal * 100,
          child: Padding(
            padding: EdgeInsets.only(
                left: size.leftPadding,
                right: size.leftPadding,
                top: size.leftPadding),
            child: widgetList == null
                ? Center(child: CircularProgressIndicator())
                : SingleChildScrollView(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: widgetList!,
                    ),
                  ),
          ),
        ),
      ),
    );
  }
}
