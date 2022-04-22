import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/extension/ethereum_address_extension.dart';

class TopAppBar extends StatefulWidget implements PreferredSizeWidget {
  TopAppBar({
    Key? key,
    this.rewardButtonOnTap,
    this.refuelButtonOnTap,
  })  : preferredSize = Size.fromHeight(100),
        super(key: key);
  @override
  final Size preferredSize;
  final void Function()? rewardButtonOnTap;
  final void Function()? refuelButtonOnTap;

  @override
  _TopAppBarState createState() => _TopAppBarState();
}

class _TopAppBarState extends State<TopAppBar> {
  bool get _seedPhraseBackuped {
    if (!isCentralized()) {
      if (common.currentUserProfile == null) return false;
      if (common.currentUserProfile!.decenUserType == DecenUserType.created) {
        return common.currentUserProfile!.seedPraseBackuped == true;
      } else
        return true;
    } else {
      return true;
    }
  }

  @override
  void initState() {
    super.initState();
    common.refreshTopBar.listen((bool refresh) {});
  }

  @override
  Widget build(BuildContext context) {
    return AppBar(
      backgroundColor: Colors.white,
      elevation: 0,
      shadowColor: Colors.transparent,
      centerTitle: true,
      automaticallyImplyLeading: false,
      toolbarHeight: widget.preferredSize.height,
      title: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Image.asset(
                'images/imgLogoH.png',
                package: 'euruswallet',
                width: 100,
                height: 28,
              ),
              if (isCentralized())
                Row(
                  children: [
                    Image.asset(
                      'images/icn_fuel.png',
                      package: 'euruswallet',
                      width: 16,
                      height: 16,
                    ),
                    SizedBox(width: 4),
                    FutureBuilder(
                      future: Future.wait([
                        web3dart.getMaxTopUpGasAmount(),
                        web3dart.eurusEthClient.getBalance(
                          EthereumAddress.fromHex(
                              common.ownerWalletAddress ?? ''),
                        ),
                      ]),
                      builder: (_, s) {
                        if (s.data == null || (s.data as List).isEmpty)
                          return Text(
                            '0.00%',
                            style: FXUI.subtitleTextStyle.copyWith(
                              color: common.getBackGroundColor(),
                            ),
                          );
                        final data = s.data as List;
                        final maxTopUpGasAmount = data[0] as double?;
                        final balance = data[1] as EtherAmount;
                        return Text(
                          maxTopUpGasAmount == null
                              ? '0.00'
                              : ((balance.getInWei.toDouble() /
                                              web3dart
                                                  .getGasPrice(
                                                      blockChainType:
                                                          BlockChainType.Eurus)
                                                  .getInWei
                                                  .toDouble() /
                                              maxTopUpGasAmount) *
                                          100)
                                      .toStringAsFixed(2) +
                                  '%',
                          style: FXUI.subtitleTextStyle.copyWith(
                            color: common.getBackGroundColor(),
                          ),
                        );
                      },
                    ),
                    SizedBox(width: 8),
                    TextButton(
                        style: ButtonStyle(
                          minimumSize: MaterialStateProperty.all(Size.zero),
                          fixedSize: MaterialStateProperty.all(
                            Size(double.infinity, 21),
                          ),
                          shape:
                              MaterialStateProperty.all<RoundedRectangleBorder>(
                            RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(6),
                              side: BorderSide(
                                  color: common.getBackGroundColor()),
                            ),
                          ),
                          padding: MaterialStateProperty.all(
                            EdgeInsets.fromLTRB(15, 4, 8, 4),
                          ),
                          textStyle: MaterialStateProperty.all(
                            FXUI.subtitleTextStyle.copyWith(
                              fontSize: 10,
                              color: common.getBackGroundColor(),
                            ),
                          ),
                        ),
                        onPressed: widget.refuelButtonOnTap,
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceAround,
                          children: [
                            Text('REFUEL.TITLE'.tr()),
                            SizedBox(width: 4),
                            Image.asset(
                              'images/rightArrow.png',
                              package: 'euruswallet',
                              width: 8,
                              height: 8,
                              color: common.getBackGroundColor(),
                            ),
                          ],
                        )),
                  ],
                ),
            ],
          ),
          SizedBox(height: 8),
          Container(
            decoration: BoxDecoration(
              color: common.getBackGroundColor(),
              borderRadius: BorderRadius.circular(FXUI.circular),
            ),
            padding: const EdgeInsets.all(16),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Row(
                  children: [
                    Image.asset(
                      'images/icn_key.png',
                      package: 'euruswallet',
                      width: 16,
                      height: 8,
                    ),
                    SizedBox(width: 8),
                    InkWell(
                      onTap: !isCentralized() && _seedPhraseBackuped
                          ? () {
                              Clipboard.setData(ClipboardData(
                                  text: web3dart.myEthereumAddress!.hexEip55
                                      .toString()));
                              common.showCopiedToClipboardSnackBar(
                                context,
                                margin: EdgeInsets.only(bottom: 56),
                                behavior: SnackBarBehavior.floating,
                              );
                            }
                          : null,
                      child: RichText(
                        text: TextSpan(
                          children: [
                            TextSpan(
                              text: _seedPhraseBackuped == true
                                  ? web3dart.myEthereumAddress
                                          ?.eip55TruncatedString ??
                                      ''
                                  : 'HOME.HEADER.SEED_PHRASE_NOT_BACKUP'.tr(),
                              style: FXUI.subtitleTextStyle.copyWith(
                                color: Colors.white,
                                decoration: TextDecoration.underline,
                              ),
                            ),
                            if (!isCentralized() && _seedPhraseBackuped)
                              WidgetSpan(
                                alignment: PlaceholderAlignment.middle,
                                child: Image.asset(
                                  "images/paste.png",
                                  package: 'euruswallet',
                                  color: Colors.white,
                                  width: 12,
                                  height: 12,
                                ),
                              ),
                          ],
                        ),
                      ),
                    ),
                  ],
                ),
                if (common.shouldShowReward ?? false)
                  InkWell(
                    onTap: widget.rewardButtonOnTap,
                    child: Row(
                      children: [
                        Image.asset(
                          "images/icn_reward.png",
                          package: 'euruswallet',
                          color: Colors.white,
                          width: 16,
                          height: 16,
                        ),
                        SizedBox(width: 8),
                        Text(
                          'REWARD.TITLE'.tr(),
                          style: FXUI.subtitleTextStyle.copyWith(
                            fontSize: 10,
                            color: Colors.white,
                          ),
                        ),
                      ],
                    ),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
