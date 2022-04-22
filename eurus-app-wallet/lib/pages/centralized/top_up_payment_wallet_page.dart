import 'dart:ui' as ui;

import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/customDialogBox.dart';
import 'package:euruswallet/commonUI/walletLockerPWDialog.dart';
import 'package:image/image.dart' as image;

export 'package:web3dart/web3dart.dart';

class TopUpPaymentWalletPage extends StatefulWidget {
  const TopUpPaymentWalletPage({
    Key? key,
  }) : super(key: key);

  @override
  _TopUpPaymentWalletPageState createState() => _TopUpPaymentWalletPageState();
}

class _TopUpPaymentWalletPageState extends State<TopUpPaymentWalletPage> {
  double _maxTopUpGasAmount = 0.0;
  double _currentGasBalance = 0.0;
  double _currentGasBalanceInPercentage = 0.0;

  double _topUpGasAmount = 0.0;
  double _topUpGasAmountInPercentage = 0.0;

  double get _topUpGasBalanceInEun {
    return _topUpGasAmount * web3dart.eurusGasPrice / pow(10, 18);
  }

  ui.Image? _sliderMarkerImage;

  RoundedLoadingButtonController _submitButtonController =
      RoundedLoadingButtonController();

  @override
  void initState() {
    super.initState();

    _loadMarkerImage();
    _fetchGasBalance();
  }

  _loadMarkerImage() async {
    final ByteData data = await rootBundle
        .load('packages/euruswallet/images/icn_fuel_slider_indicator.png');
    final baseSizeImage = image.decodeImage(data.buffer.asUint8List());
    if (baseSizeImage == null) return;
    final resizeImage = image.copyResize(baseSizeImage, height: 50, width: 40);
    ui.Codec codec = await ui.instantiateImageCodec(
        Uint8List.fromList(image.encodePng(resizeImage)));
    ui.FrameInfo frameInfo = await codec.getNextFrame();

    setState(() {
      _sliderMarkerImage = frameInfo.image;
    });
  }

  _fetchGasBalance() {
    Future.wait([
      web3dart.getMaxTopUpGasAmount(),
      web3dart.eurusEthClient
          .getBalance(EthereumAddress.fromHex(common.ownerWalletAddress ?? '')),
    ]).then((value) {
      setState(() {
        _maxTopUpGasAmount = value.first != null && value.first is double
            ? value.first as double
            : 0.0;
        _currentGasBalance = value[1] is EtherAmount
            ? (value[1] as EtherAmount).getInWei.toDouble() /
                web3dart.eurusGasPrice
            : 0.0;
        _updateTopUpGasAmount(_topUpGasAmount);
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        title: Text(
          'REFUEL.TITLE'.tr(),
          style: FXUI.inputStyle.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        leading: IconButton(
          icon: Icon(Icons.arrow_back_ios_outlined,
              color: common.getBackGroundColor()),
          onPressed: () => Navigator.of(context).pop(),
        ),
        backgroundColor: Colors.transparent,
        shadowColor: Colors.transparent,
      ),
      body: SafeArea(
        bottom: false,
        child: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.all(24.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                _getTopUpInputWidget(),
                SizedBox(height: 15),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _getTopUpInputWidget() {
    return _getShadowContainer(
      child: Padding(
        padding: const EdgeInsets.all(14.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            _getTitleRow(),
            _getGasInputSlider(),
            Container(
              padding: EdgeInsets.fromLTRB(12.0, 12.0, 12.0, 6),
              decoration: BoxDecoration(
                border: Border.all(color: common.getBackGroundColor()),
                borderRadius: BorderRadius.circular(14),
              ),
              child: Column(
                children: [
                  Text(
                    _formatToString(_topUpGasAmount / pow(10, 6)) + 'M',
                    style: FXUI.titleTextStyle.copyWith(
                      fontSize: 24,
                      fontWeight: FontWeight.normal,
                      color: common.getBackGroundColor(),
                    ),
                  ),
                  SizedBox(height: 4),
                  Text(
                    'REFUEL.REFUEL_GAS_VOLUME'.tr(),
                    style: FXUI.titleTextStyle.copyWith(
                      fontSize: 13,
                      color: FXColor.darkBlack,
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: EdgeInsets.only(top: 8),
              child: Center(
                child: Text(
                  'REFUEL.SPEND_EUN'
                      .tr(args: [_formatToString(_topUpGasBalanceInEun)]),
                  style: FXUI.hintStyle.copyWith(
                    fontSize: 11,
                  ),
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: Icon(
                Icons.arrow_downward,
                color: common.getBackGroundColor(),
              ),
            ),
            Container(
              padding: EdgeInsets.fromLTRB(12.0, 12.0, 12.0, 6),
              decoration: BoxDecoration(
                border: Border.all(color: common.getBackGroundColor()),
                borderRadius: BorderRadius.circular(14),
              ),
              child: Column(
                children: [
                  Text(
                    (_formatToString((_currentGasBalance + _topUpGasAmount) /
                            pow(10, 6))) +
                        'M',
                    style: FXUI.titleTextStyle.copyWith(
                      fontSize: 24,
                      fontWeight: FontWeight.normal,
                      color: FXColor.darkBlack,
                    ),
                  ),
                  SizedBox(height: 4),
                  Text(
                    'REFUEL.TOTAL_GAS_VOLUME'.tr(),
                    style: FXUI.titleTextStyle.copyWith(
                      fontSize: 13,
                      color: FXColor.darkBlack,
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: EdgeInsets.symmetric(vertical: 20.0),
              child: SubmitButton(
                buttonBGColor: common.getBackGroundColor(),
                onPressed: _topUpGasBalanceInEun.toDouble() == 0.0
                    ? null
                    : () async {
                        RoundedLoadingButtonController();
                        Future<String?> Function(TextEditingController)?
                            _submit = (_textEditingController) async {
                          final serverAddressPair = await common.getAddressPair(
                              email: common.email,
                              password: _textEditingController.text,
                              mnemonic: common.serverMnemonic,
                              addressPairType: AddressPairType.paymentPw);
                          if (serverAddressPair.address.toLowerCase() ==
                              common.ownerWalletAddress?.toLowerCase()) {
                            common.cenSignKey = serverAddressPair.privateKey;
                          } else {
                            common.cenSignKey = null;
                          }
                          return common.cenSignKey;
                        };

                        String? result = await showGeneralDialog(
                          context: context,
                          barrierColor: Colors.black45,
                          pageBuilder: (
                            BuildContext context,
                            Animation<double> animation,
                            Animation<double> secondaryAnimation,
                          ) =>
                              WalletLockerPWDialog(
                            themeColor: common.getBackGroundColor(),
                            cenSubmitFnc: _submit,
                          ),
                        );

                        if (!isEmptyString(string: result)) {
                          try {
                            await api.topUpPaymentWallet(
                              targetGasAmount:
                                  (_currentGasBalance + _topUpGasAmount)
                                      .toInt(),
                            );

                            await showDialog(
                              context: context,
                              builder: (BuildContext context) {
                                return CustomDialogBox(
                                  btnColor: common.getBackGroundColor(),
                                  descriptions: "COMMON.SUCCESS".tr(),
                                  buttonText: "COMMON.OK".tr(),
                                );
                              },
                            );

                            _fetchGasBalance();
                          } catch (e) {
                            await showDialog(
                              context: context,
                              builder: (BuildContext context) {
                                return CustomDialogBox(
                                  btnColor: common.getBackGroundColor(),
                                  descriptions: e.toString(),
                                  buttonText: "COMMON.OK".tr(),
                                );
                              },
                            );
                          }
                        }
                        _submitButtonController.reset();
                      },
                label: 'COMMON.CONFIRM'.tr(),
                btnController: _submitButtonController,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Container _getShadowContainer({required Widget child}) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
              color: Colors.black.withOpacity(0.16),
              offset: Offset(0, 3),
              blurRadius: 12)
        ],
        borderRadius: BorderRadius.circular(19),
      ),
      child: child,
    );
  }

  Padding _getTitleRow() {
    return Padding(
      padding: const EdgeInsets.only(top: 6.0),
      child: Row(
        children: [
          Expanded(
            child: Center(
              child: Text(
                'REFUEL.SLIDER.SLIDE_TO_REFUEL'.tr(),
                style: FXUI.hintStyle.copyWith(
                  color: FXColor.darkBlack,
                ),
              ),
            ),
          ),
          SizedBox(width: 4),
          Icon(
            Icons.info,
            size: 24,
            color: common.getBackGroundColor(),
          ),
        ],
      ),
    );
  }

  Container _getGasInputSlider() {
    return Container(
      decoration: BoxDecoration(
        color: FXColor.lightWhiteColor,
        borderRadius: BorderRadius.circular(12),
      ),
      padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      margin: EdgeInsets.symmetric(vertical: 20),
      child: Column(
        children: [
          Text(
            "REFUEL.SLIDER.CURRENT_GAS_VOLUME".tr(
                args: [_formatToString(_currentGasBalanceInPercentage * 100)]),
            style: FXUI.hintStyle.copyWith(color: common.getBackGroundColor()),
          ),
          Stack(
            children: [
              Image.asset(
                'images/icn_fuel_slider_background.png',
                package: 'euruswallet',
              ),
              Positioned.fill(
                left: 16,
                top: 13,
                right: _maxTopUpGasAmount == 0.0
                    ? double.infinity
                    : 16 +
                        (MediaQuery.of(context).size.width - 2 * 70) *
                            (1 - _currentGasBalanceInPercentage),
                bottom: 20,
                child: Container(
                  decoration: BoxDecoration(
                    color: common.getBackGroundColor(),
                    borderRadius: BorderRadius.only(
                      topRight: Radius.circular(8),
                      bottomRight: Radius.circular(8),
                    ),
                  ),
                ),
              ),
              Positioned.fill(
                left: _currentGasBalanceInPercentage > 0
                    ? ((MediaQuery.of(context).size.width - 2 * 70) *
                                _currentGasBalanceInPercentage >=
                            8
                        ? (8 +
                            (MediaQuery.of(context).size.width - 2 * 70) *
                                _currentGasBalanceInPercentage)
                        : 16)
                    : 16,
                top: 13,
                right: _maxTopUpGasAmount == 0.0
                    ? double.infinity
                    : 16 +
                        (MediaQuery.of(context).size.width - 2 * 70) *
                            (1 -
                                _currentGasBalanceInPercentage -
                                _topUpGasAmountInPercentage),
                bottom: 20,
                child: Container(
                  decoration: BoxDecoration(
                    color: common.getBackGroundColor().withOpacity(0.25),
                    borderRadius: BorderRadius.only(
                      topRight: Radius.circular(8),
                      bottomRight: Radius.circular(8),
                    ),
                  ),
                ),
              )
            ],
          ),
          Stack(
            children: [
              SliderTheme(
                data: SliderTheme.of(context).copyWith(
                  trackHeight: 0,
                  trackShape: TopUpPaymentSliderTrackShape(),
                  overlayColor: Colors.transparent,
                  thumbShape: _sliderMarkerImage != null
                      ? TopUpPaymentSliderThumbShape(image: _sliderMarkerImage!)
                      : null,
                  showValueIndicator: ShowValueIndicator.never,
                ),
                child: Slider(
                  value: _currentGasBalance + _topUpGasAmount,
                  min: 0,
                  max: _maxTopUpGasAmount,
                  onChanged: (value) {
                    setState(() {
                      _updateTopUpGasAmount(value);
                    });
                  },
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  void _updateTopUpGasAmount(double value) {
    if (value <= _currentGasBalance) {
      _topUpGasAmount = 0;
    } else {
      List<double> marks = [
        0,
        _currentGasBalance,
        _maxTopUpGasAmount * 1 / 4,
        _maxTopUpGasAmount * 2 / 4,
        _maxTopUpGasAmount * 3 / 4,
        _maxTopUpGasAmount * 4 / 4,
      ]..sort();

      List<double> marksMidPoints = [];
      marks.forEach((e) {
        final index = marks.indexOf(e);
        if (index < marks.length - 1) {
          marksMidPoints.add((marks[index] + marks[index + 1]) / 2);
        }
      });

      for (var element in marksMidPoints) {
        final index = marksMidPoints.indexOf(element);
        if (value >= marksMidPoints.last) {
          _topUpGasAmount = _maxTopUpGasAmount - _currentGasBalance;
          break;
        } else if (value <= element) {
          _topUpGasAmount = marks[index] - _currentGasBalance;
          if (_topUpGasAmount <= 0) _topUpGasAmount = 0;
          break;
        }
      }
    }
    _topUpGasAmountInPercentage = _topUpGasAmount / _maxTopUpGasAmount;
    _currentGasBalanceInPercentage = _currentGasBalance / _maxTopUpGasAmount;
  }

  String _formatToString(double value) {
    return value.toStringAsFixed(2);
  }
}

class TopUpPaymentSliderThumbShape extends SliderComponentShape {
  const TopUpPaymentSliderThumbShape({
    required this.image,
  });

  final ui.Image image;

  @override
  Size getPreferredSize(bool isEnabled, bool isDiscrete) {
    return Size.zero;
  }

  @override
  void paint(
    PaintingContext context,
    ui.Offset center, {
    required Animation<double> activationAnimation,
    required Animation<double> enableAnimation,
    required bool isDiscrete,
    required TextPainter labelPainter,
    required RenderBox parentBox,
    required SliderThemeData sliderTheme,
    required ui.TextDirection textDirection,
    required double value,
    required double textScaleFactor,
    required ui.Size sizeWithOverflow,
  }) {
    final Canvas canvas = context.canvas;
    canvas.drawImage(
      image,
      Offset(center.dx - image.width / 2, -center.dy + image.height / 2),
      Paint(),
    );
  }
}

class TopUpPaymentSliderTrackShape extends RoundedRectSliderTrackShape {
  @override
  ui.Rect getPreferredRect({
    required RenderBox parentBox,
    ui.Offset offset = Offset.zero,
    required SliderThemeData sliderTheme,
    bool isEnabled = false,
    bool isDiscrete = false,
  }) {
    final double trackHeight = sliderTheme.trackHeight ?? 0.0;
    final double trackLeft = offset.dx + 15;
    final double trackTop =
        offset.dy + (parentBox.size.height - trackHeight) / 2;
    final double trackWidth = parentBox.size.width - 34;
    return Rect.fromLTWH(trackLeft, trackTop, trackWidth, trackHeight);
  }
}
