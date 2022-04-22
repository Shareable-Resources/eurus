import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/topSelectBlockChainBar.dart';
import 'package:flutter/cupertino.dart';
import 'package:path_provider/path_provider.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:screenshot/screenshot.dart';
import 'package:share/share.dart';

class ReceivePage extends StatefulWidget {
  ReceivePage({
    Key? key,
    this.errorPromptBuilder,
    String? errorText,
    bool? ethereumErrorPopUp,
    bool? eurusErrorPopUp,
    BlockChainType? blockChainType,
    this.replacingQRCodeWidget,
    String? ethereumAddress,
    bool? disableSelectBlockchain,
  })  : this.errorText = errorText ?? '',
        this.ethereumErrorPopUp = ethereumErrorPopUp ?? false,
        this.eurusErrorPopUp = ethereumErrorPopUp ?? false,
        this.blockChainType = blockChainType ?? BlockChainType.Eurus,
        this.disableSelectBlockchain = disableSelectBlockchain ?? false,
        this.ethereumAddress = (isEmptyString(string: ethereumAddress)
                ? web3dart.myEthereumAddress.toString()
                : ethereumAddress) ??
            '',
        this.backUpEurusAddress = ethereumAddress ?? '',
        super(key: key);
  final Widget Function(BuildContext)? errorPromptBuilder;
  final String errorText;
  final bool ethereumErrorPopUp;
  final bool eurusErrorPopUp;
  final Widget? replacingQRCodeWidget;
  String ethereumAddress;
  final BlockChainType blockChainType;
  final bool disableSelectBlockchain;
  final String backUpEurusAddress;

  @override
  _ReceivePageState createState() => _ReceivePageState();
}

class _ReceivePageState extends State<ReceivePage> {
  ScreenshotController screenshotController = ScreenshotController();

  BlockChainType? currentNetworkSelection;

  @override
  void initState() {
    super.initState();
    common.fromBlockChainType = widget.blockChainType;
    // common.currentBlockchainSelection = widget.blockChainType;
  }

  @override
  void dispose() {
    super.dispose();
  }

  void displayDialog() {
    if (widget.errorPromptBuilder != null &&
        ((common.fromBlockChainType == BlockChainType.Ethereum &&
                widget.ethereumErrorPopUp) ||
            (common.fromBlockChainType == BlockChainType.Eurus &&
                widget.eurusErrorPopUp))) {
      showDialog(context: context, builder: widget.errorPromptBuilder!);
    }
  }

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);

    double qrcodeSize = MediaQuery.of(context).size.width * 0.6;
    print('receivePage: ${common.fromBlockChainType}');
    print(
        'receivePage: ${common.fromBlockChainType == BlockChainType.Ethereum}');

    return Screenshot(
      controller: screenshotController,
      child: BackGroundImage(
        child: Stack(
          children: <Widget>[
            Scaffold(
              backgroundColor: Colors.transparent,
              appBar: WalletAppBar(
                title: "RECEIVE_PAGE.MAIN_TITLE".tr(),
                rightWidget:
                    (isCentralized() && currentNetworkSelection != null) ||
                            (!isCentralized() &&
                                widget.replacingQRCodeWidget == null)
                        ? Icon(Icons.ios_share)
                        : null,
                function: () {
                  if (widget.replacingQRCodeWidget == null) {
                    screenshotController
                        .capture(pixelRatio: 5)
                        .then((Uint8List? image) async {
                      if (image != null) {
                        final directory =
                            await getApplicationDocumentsDirectory();
                        final imagePath =
                            await File('${directory.path}/image.png').create();
                        await imagePath.writeAsBytes(image);

                        /// Share Plugin
                        String shareText = 'RECEIVE_PAGE.SHARE_TEXT'.tr(args: [
                          getBlockChainName(common.currentBlockchainSelection)
                              .toLowerCase(),
                          widget.ethereumAddress
                        ]);
                        await Share.shareFiles([imagePath.path],
                            text: shareText);
                      }
                    }).catchError((onError) {
                      print(onError);
                    });
                  }
                },
              ),
              body: SingleChildScrollView(
                child: Center(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      Padding(
                        padding: EdgeInsets.symmetric(vertical: 20),
                        child: Container(
                          decoration: FXUI.circleBoxDecoration,
                          width: size.blockSizeHorizontal * 100 - 46,
                          child: Column(
                            children: [
                              Padding(
                                padding: EdgeInsets.only(top: 15),
                                child: isCentralized()
                                    ? Padding(
                                        padding: EdgeInsets.symmetric(
                                            vertical:
                                                currentNetworkSelection == null
                                                    ? 20.0
                                                    : 10.0),
                                        child: DropdownButtonHideUnderline(
                                          child: ButtonTheme(
                                            alignedDropdown: true,
                                            child:
                                                DropdownButton<BlockChainType>(
                                              alignment: Alignment.center,
                                              dropdownColor: Colors.white,
                                              borderRadius:
                                                  BorderRadius.circular(12),
                                              hint: Text(
                                                'RECEIVE_PAGE.SELECT_NETWORK'
                                                    .tr(),
                                                style: FXUI.titleTextStyle
                                                    .copyWith(
                                                  fontSize: 16,
                                                  color: FXColor.deepGreyColor,
                                                ),
                                              ),
                                              value: currentNetworkSelection,
                                              icon: currentNetworkSelection ==
                                                      null
                                                  ? Padding(
                                                      padding:
                                                          const EdgeInsets.all(
                                                              8.0),
                                                      child: Image.asset(
                                                        'images/icn_arrow_down.png',
                                                        package: 'euruswallet',
                                                        width: 16,
                                                        color: FXColor
                                                            .deepGreyColor,
                                                      ),
                                                    )
                                                  : Container(),
                                              items: [
                                                BlockChainType.Eurus,
                                                BlockChainType.Ethereum,
                                              ]
                                                  .map(
                                                    (e) => DropdownMenuItem<
                                                        BlockChainType>(
                                                      value: e,
                                                      child: Container(
                                                        decoration:
                                                            BoxDecoration(
                                                          borderRadius:
                                                              FXUI.cricleRadius,
                                                          color: currentNetworkSelection !=
                                                                  null
                                                              ? Colors.white
                                                              : common
                                                                  .getBackGroundColor(),
                                                        ),
                                                        padding:
                                                            EdgeInsets.fromLTRB(
                                                                12, 5, 12, 5),
                                                        child: Row(
                                                          children: [
                                                            Image.asset(
                                                              'images/${e == BlockChainType.Eurus ? 'icn_eun.png' : 'icn_eth.png'}',
                                                              package:
                                                                  'euruswallet',
                                                              width:
                                                                  currentNetworkSelection !=
                                                                          null
                                                                      ? 26
                                                                      : 16,
                                                              height:
                                                                  currentNetworkSelection !=
                                                                          null
                                                                      ? 26
                                                                      : 16,
                                                              color: currentNetworkSelection !=
                                                                      null
                                                                  ? common
                                                                      .getBackGroundColor()
                                                                  : Colors
                                                                      .white,
                                                            ),
                                                            SizedBox(
                                                              width: 3,
                                                            ),
                                                            Text(
                                                              getBlockChainName(
                                                                  e),
                                                              style: Theme.of(
                                                                      context)
                                                                  .textTheme
                                                                  .bodyText1
                                                                  ?.apply(
                                                                    color: currentNetworkSelection !=
                                                                            null
                                                                        ? common
                                                                            .getBackGroundColor()
                                                                        : Colors
                                                                            .white,
                                                                    fontWeightDelta:
                                                                        2,
                                                                  )
                                                                  .copyWith(
                                                                      fontSize: currentNetworkSelection !=
                                                                              null
                                                                          ? 26
                                                                          : 13),
                                                            ),
                                                          ],
                                                        ),
                                                      ),
                                                    ),
                                                  )
                                                  .toList(),
                                              onChanged:
                                                  currentNetworkSelection ==
                                                          null
                                                      ? (value) {
                                                          if (value != null) {
                                                            if (value ==
                                                                currentNetworkSelection)
                                                              return;

                                                            setState(() {
                                                              currentNetworkSelection =
                                                                  value;
                                                              common.currentBlockchainSelection =
                                                                  value;
                                                              common.selectTokenSymbol =
                                                                  null;
                                                              if (isCentralized()) {
                                                                widget
                                                                    .ethereumAddress = (common.currentBlockchainSelection ==
                                                                            BlockChainType
                                                                                .Ethereum
                                                                        ? common
                                                                            .cenMainNetWalletAddress
                                                                        : web3dart
                                                                            .myEthereumAddress
                                                                            .toString()) ??
                                                                    '';
                                                              }
                                                              displayDialog();
                                                            });
                                                          }
                                                        }
                                                      : null,
                                            ),
                                          ),
                                        ),
                                      )
                                    : Container(
                                        child: TopSelectBlockChainBar(
                                          topBarType: TopSelectBarType.enable,
                                          disableSelectBlockchain:
                                              widget.disableSelectBlockchain,
                                          currentSelection:
                                              common.currentBlockchainSelection,
                                          onSegmentChosen: (blockChainType) {
                                            setState(() {
                                              // common.currentBlockchainSelection =
                                              //     blockChainType;
                                              // common.fromBlockChainType =
                                              //     common.currentBlockchainSelection;

                                              if (isCentralized()) {
                                                setState(() {
                                                  widget
                                                      .ethereumAddress = (common
                                                                  .currentBlockchainSelection ==
                                                              BlockChainType
                                                                  .Ethereum
                                                          ? common
                                                              .cenMainNetWalletAddress
                                                          : web3dart
                                                              .myEthereumAddress
                                                              .toString()) ??
                                                      '';
                                                });
                                              }
                                              displayDialog();
                                            });
                                          },
                                        ),
                                      ),
                              ),
                              if (isCentralized() &&
                                  currentNetworkSelection == null)
                                Padding(
                                  padding: const EdgeInsets.only(
                                    left: 24.0,
                                    right: 24.0,
                                    bottom: 24.0,
                                  ),
                                  child: Text(
                                    'RECEIVE_PAGE.SELECT_NETWORK_STATEMENT'
                                        .tr(),
                                    style: FXUI.hintStyle.copyWith(
                                      color: FXColor.grey44,
                                      fontSize: 16,
                                      fontWeight: FontWeight.normal,
                                    ),
                                    textAlign: TextAlign.center,
                                  ),
                                ),
                              // Text("Powered By Eurus",
                              if ((isCentralized() &&
                                      currentNetworkSelection != null) ||
                                  !isCentralized())
                                Container(
                                  margin: EdgeInsets.only(bottom: 30),
                                  padding: EdgeInsets.all(qrcodeSize * 0.035),
                                  decoration: new BoxDecoration(
                                    image: new DecorationImage(
                                        image: new AssetImage(
                                            "images/qrcodeBorder.png",
                                            package: 'euruswallet'),
                                        fit: BoxFit.cover),
                                  ),
                                  child: QrImage(
                                    data:
                                        '${common.currentBlockchainSelection == BlockChainType.Eurus ? getBlockChainName(common.currentBlockchainSelection).toLowerCase() + ':' : ''}${widget.ethereumAddress}',
                                    version: QrVersions.auto,
                                    embeddedImage: AssetImage(
                                      'images/${common.currentBlockchainSelection == BlockChainType.Ethereum ? 'ETH' : !isCentralized() ? 'Eurus_Violet' : 'Eurus_Blue'}.png',
                                      package: 'euruswallet',
                                    ),
                                    foregroundColor: FXColor.blackColor,
                                    size: qrcodeSize,
                                    errorCorrectionLevel: QrErrorCorrectLevel.M,
                                  ),
                                ),
                              if (isCentralized() &&
                                  currentNetworkSelection != null)
                                Container(
                                  decoration: BoxDecoration(
                                    borderRadius: BorderRadius.circular(12),
                                    color: FXColor.dividerColor,
                                  ),
                                  margin: const EdgeInsets.only(
                                    left: 16.0,
                                    right: 16.0,
                                    bottom: 16.0,
                                  ),
                                  padding: const EdgeInsets.all(16.0),
                                  child: Column(
                                    children: [
                                      TextButton.icon(
                                        onPressed: null,
                                        icon: Icon(
                                          Icons.info,
                                          color: common.getBackGroundColor(),
                                        ),
                                        label: Text(
                                          'COMMON.REMIND'.tr(),
                                          style: FXUI.inputStyle.copyWith(
                                            fontWeight: FontWeight.w500,
                                          ),
                                        ),
                                        style: TextButton.styleFrom(
                                          padding: EdgeInsets.zero,
                                          minimumSize: Size.zero,
                                          tapTargetSize:
                                              MaterialTapTargetSize.shrinkWrap,
                                        ),
                                      ),
                                      Padding(
                                        padding:
                                            const EdgeInsets.only(top: 10.0),
                                        child: Text(
                                          currentNetworkSelection ==
                                                  BlockChainType.Eurus
                                              ? 'RECEIVE_PAGE.SELECT_NETWORK_REMINDER'
                                                  .tr()
                                              : 'RECEIVE_PAGE.DISCLAIMER_DECEN'
                                                  .tr(),
                                          style: FXUI.inputStyle.copyWith(
                                            fontSize: 13,
                                            fontWeight: FontWeight.normal,
                                          ),
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                            ],
                          ),
                        ),
                      ),
                      (!isCentralized() &&
                                  widget.replacingQRCodeWidget == null) ||
                              (isCentralized() &&
                                  currentNetworkSelection != null)
                          ? Container(
                              width: size.blockSizeHorizontal * 100 - 46,
                              decoration: FXUI.circleBoxDecoration,
                              padding: EdgeInsets.symmetric(
                                vertical: 18,
                                horizontal: 15,
                              ),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Row(
                                    children: [
                                      Expanded(
                                        flex: 4,
                                        child: Column(
                                          crossAxisAlignment:
                                              CrossAxisAlignment.stretch,
                                          children: [
                                            Text(
                                              "RECEIVE_PAGE.MY_WALLET_ADDRESS"
                                                  .tr(),
                                              style:
                                                  FXUI.normalTextStyle.copyWith(
                                                color: FXColor.deepGreyColor,
                                                fontWeight: FontWeight.w600,
                                              ),
                                            ),
                                            Padding(padding: EdgeInsets.all(2)),
                                            Text(
                                              EthereumAddress.fromHex(
                                                      widget.ethereumAddress)
                                                  .hexEip55,
                                              style:
                                                  FXUI.normalTextStyle.copyWith(
                                                color: FXColor.lightGray,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      Container(
                                          child: new Material(
                                              child: InkWell(
                                                onTap: () {
                                                  Clipboard.setData(ClipboardData(
                                                      text: EthereumAddress
                                                              .fromHex(widget
                                                                  .ethereumAddress)
                                                          .hexEip55));
                                                  common
                                                      .showCopiedToClipboardSnackBar(
                                                    context,
                                                  );
                                                },
                                                child: Padding(
                                                  padding: EdgeInsets.only(
                                                      left: 10,
                                                      top: 10,
                                                      bottom: 10),
                                                  child: Image.asset(
                                                    "images/paste.png",
                                                    package: 'euruswallet',
                                                    height: 20,
                                                    width: 20,
                                                    color: common
                                                        .getBackGroundColor(),
                                                    fit: BoxFit.contain,
                                                  ),
                                                ),
                                              ),
                                              color: Colors.white)),
                                    ],
                                  ),
                                ],
                              ),
                            )
                          : Container(),
                      Padding(
                        padding: EdgeInsets.only(left: 23, right: 23, top: 23),
                        child: Text(
                            !isCentralized()
                                ? 'RECEIVE_PAGE.DISCLAIMER_DECEN'.tr()
                                : '',
                            style: FXUI.normalTextStyle
                                .copyWith(color: Colors.white, fontSize: 14)),
                      )
                    ],
                  ),
                ),
              ),
            ),
            if (widget.replacingQRCodeWidget != null)
              widget.replacingQRCodeWidget!
          ],
        ),
      ),
    );
  }
}
