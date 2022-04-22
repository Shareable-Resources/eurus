import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';

class WalletLockerPWDialog extends StatefulWidget {
  const WalletLockerPWDialog({
    Key? key,
    this.textEditingController,
    this.decenUserkey,
    this.submitFnc,
    this.cenSubmitFnc,
    this.tryBioAuthFnc,
    this.themeColor,
  }) : super(key: key);

  final TextEditingController? textEditingController;
  final String? decenUserkey;
  final Future<String?> Function(String?, TextEditingController)? submitFnc;
  final Future<String?> Function(TextEditingController)? cenSubmitFnc;

  final Future<String?> Function()? tryBioAuthFnc;

  final Color? themeColor;

  @override
  _WalletLockerPWDialogState createState() => _WalletLockerPWDialogState();
}

class _WalletLockerPWDialogState extends State<WalletLockerPWDialog> {
  late TextEditingController _textEditingController;
  bool _isPwMasked = false;
  Color _themeColor = FXColor.mainBlueColor;

  @override
  void initState() {
    super.initState();
    _textEditingController =
        widget.textEditingController ?? TextEditingController();
    if (widget.themeColor != null) _themeColor = widget.themeColor!;
    if (widget.tryBioAuthFnc != null)
      widget.tryBioAuthFnc!().then((value) {
        _textEditingController.text = value ?? '';
        __submit();
      });
  }

  @override
  void dispose() {
    _textEditingController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black87,
      body: Center(
        child: SingleChildScrollView(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Container(
                decoration: BoxDecoration(
                  borderRadius: FXUI.cricleRadius,
                  color: Colors.white.withOpacity(.9),
                ),
                padding: EdgeInsets.all(16),
                child: Column(
                  children: [
                    ListTile(
                      title: Text(
                        isCentralized()
                            ? "WALLET_LOCKER_DIALOG.CEN_TITLE".tr()
                            : "WALLET_LOCKER_DIALOG.TITLE".tr(),
                        style: Theme.of(context)
                            .textTheme
                            .headline6
                            ?.apply(fontWeightDelta: 2),
                        textAlign: TextAlign.center,
                      ),
                    ),
                    ListTile(
                      title: TextField(
                        textInputAction: TextInputAction.go,
                        onSubmitted: (_) {
                          __submit();
                        },
                        autofocus: true,
                        obscureText: !_isPwMasked,
                        controller: _textEditingController,
                        decoration: InputDecoration(
                          border: OutlineInputBorder(
                            borderSide: BorderSide(color: _themeColor),
                            borderRadius: FXUI.cricleRadius,
                          ),
                          focusedBorder: OutlineInputBorder(
                            borderSide:
                                BorderSide(color: _themeColor, width: 2),
                            borderRadius: FXUI.cricleRadius,
                          ),
                          suffixIcon: IconButton(
                            onPressed: () {
                              setState(() {
                                _isPwMasked = !_isPwMasked;
                              });
                            },
                            icon: Image.asset(
                              _isPwMasked
                                  ? 'images/eyeClose.png'
                                  : 'images/eyeOpen.png',
                              package: 'euruswallet',
                              width: 16,
                              height: 16,
                              color: common.getBackGroundColor(),
                            ),
                          ),
                          hintText: 'COMMON.PASSWORD_HINT'.tr(),
                          hintStyle: FXUI.normalTextStyle.copyWith(
                            fontWeight: FontWeight.w600,
                            fontSize: 14,
                            color: Colors.grey.shade400,
                          ),
                        ),
                        cursorColor: _themeColor,
                      ),
                    ),
                    ListTile(
                      title: TextButton(
                        child: Text(
                          "COMMON.SUBMIT".tr(),
                          style: Theme.of(context)
                              .textTheme
                              .caption
                              ?.apply(color: _themeColor, fontSizeDelta: 3),
                        ),
                        onPressed: () {
                          __submit();
                        },
                      ),
                      subtitle: TextButton(
                          child: Text(
                            "COMMON.CANCEL".tr(),
                            style: FXUI.normalTextStyle
                                .copyWith(color: _themeColor),
                          ),
                          onPressed: () {
                            common.cenSignKey = null;
                            Navigator.of(context).pop();
                          }),
                    )
                  ],
                ),
              )
            ],
          ),
        ),
      ),
    );
  }

  Future<void> __submit() async {
    ScaffoldMessenger.of(context).hideCurrentSnackBar();
    String? value;
    if (isCentralized() && widget.cenSubmitFnc != null) {
      value = await widget.cenSubmitFnc!(_textEditingController);
    } else if (widget.submitFnc != null) {
      value = await widget.submitFnc!(
          widget.decenUserkey ?? "", _textEditingController);
    }

    _trySubmitPW(context, value);
  }

  void _trySubmitPW(BuildContext _, String? val) {
    if (val != null) {
      _successAndPop(_, val: val);
    } else {
      _showValidationError(_);
    }
  }

  void _showValidationError(BuildContext _) {
    ScaffoldMessenger.of(_).showSnackBar(
      SnackBar(
        content:
            Text("COMMON_ERROR.AUTH_FAIL".tr(), textAlign: TextAlign.center),
        backgroundColor: Colors.redAccent.shade100.withOpacity(.9),
      ),
    );
  }

  void _successAndPop(BuildContext _, {String? val}) {
    Navigator.of(_).maybePop(val);
  }
}
