import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/common/web3dart.dart';
import 'package:flutter/cupertino.dart';

typedef void OnSegmentChosen(BlockChainType type);

enum TopSelectBarType { enable, disable }

class TopSelectBlockChainBar extends StatefulWidget {
  TopSelectBlockChainBar({
    Key? key,
    this.onSegmentChosen,
    required this.currentSelection,
    this.disableSelectBlockchain: false,
    this.topBarType,
    this.dropDownList: const [
      BlockChainType.Eurus,
      BlockChainType.Ethereum,
      // BlockChainType.BinanceCoin
    ],
  });

  final OnSegmentChosen? onSegmentChosen;
  BlockChainType currentSelection;
  final bool disableSelectBlockchain;
  final TopSelectBarType? topBarType;
  final List<BlockChainType> dropDownList;

  @override
  _TopSelectBlockChainBarState createState() => _TopSelectBlockChainBarState();
}

class _TopSelectBlockChainBarState extends State<TopSelectBlockChainBar> {
  @override
  Widget build(BuildContext context) {
    if (widget.dropDownList.contains(widget.currentSelection) == false) {
      widget.currentSelection = BlockChainType.Eurus;
    }

    return widget.disableSelectBlockchain
        ? Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              _tabStyle(widget.dropDownList.firstWhere(
                  (element) => element == widget.currentSelection)),
            ],
          )
        : Container(
            padding: EdgeInsets.all(4),
            child: DropdownButton<BlockChainType>(
              value: widget.currentSelection,
              icon: const Icon(Icons.arrow_drop_down),
              iconSize: 24,
              elevation: 16,
              style: const TextStyle(color: Colors.deepPurple),
              underline: Container(),
              onChanged: (BlockChainType? newValue) {
                if (newValue != null) {
                  common.currentBlockchainSelection = newValue;
                  common.selectTokenSymbol = null;
                  if (widget.topBarType != TopSelectBarType.enable ||
                      newValue == widget.currentSelection) return;
                  if (widget.onSegmentChosen != null)
                    widget.onSegmentChosen!(newValue);
                }
              },
              items: widget.dropDownList.map((data) {
                return DropdownMenuItem<BlockChainType>(
                  value: data,
                  child: _tabStyle(data),
                );
              }).toList(),
            ),
          );
  }

  Widget _tabStyle(BlockChainType chainType) {
    final bool isActive = widget.currentSelection == chainType;
    return Container(
      decoration: BoxDecoration(
        borderRadius: FXUI.cricleRadius,
        color: isActive ? common.getBackGroundColor() : Colors.transparent,
      ),
      padding: EdgeInsets.fromLTRB(12, 5, 12, 5),
      child: Row(
        children: [
          if (getIcon(chainType, isActive) != null)
            getIcon(chainType, isActive)!,
          SizedBox(
            width: 3,
          ),
          Text(
            getBlockChainName(chainType),
            style: Theme.of(context)
                .textTheme
                .bodyText1
                ?.apply(
                  color: isActive ? Colors.white : FXColor.textGray,
                  fontWeightDelta: 2,
                )
                .copyWith(fontSize: 13),
          ),
        ],
      ),
      //   ),
    );
  }

  Widget? getIcon(
    BlockChainType blockChainType,
    bool isActive,
  ) {
    String imagePath = "";

    switch (blockChainType) {
      case BlockChainType.Eurus:
        imagePath = "images/icn_eun.png";
        break;
      case BlockChainType.Ethereum:
        imagePath = "images/icn_eth.png";
        break;
      // case BlockChainType.BinanceCoin:
      //   imagePath = "images/icn_eth.png";
      //   break;
      default:
        return null;
    }

    return Image.asset(
      imagePath,
      package: 'euruswallet',
      width: 16,
      height: 16,
      color: isActive ? Colors.white : FXColor.textGray,
    );
  }
}
