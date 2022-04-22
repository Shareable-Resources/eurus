import 'package:euruswallet/common/commonMethod.dart';
import 'package:rounded_loading_button/rounded_loading_button.dart';

class SubmitButton extends StatefulWidget {
  final Function? onPressed;
  final String label;
  final Color textColor;
  final Color? buttonBGColor;
  final double fontSize;
  bool submitBtnLoading = false;
  int loadingSecond;
  RoundedLoadingButtonController btnController;

  SubmitButton(
      {Key? key,
      this.onPressed,
      required this.label,
      this.textColor: Colors.white,
      this.buttonBGColor,
      this.fontSize: 16,
      this.loadingSecond: 0,
      required this.btnController})
      : super(key: key);

  @override
  SubmitButtonState createState() => SubmitButtonState();
}

class SubmitButtonState extends State<SubmitButton> {
  Timer? _timer;

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    Color? buttonBGColor = widget.buttonBGColor == null
        ? common.getBackGroundColor()
        : widget.buttonBGColor;

    return RoundedLoadingButton(
        borderRadius: 15,
        width: 340,
        color: buttonBGColor,
        onPressed: widget.onPressed == null
            ? null
            : () {
                if (widget.submitBtnLoading) {
                  return;
                }
                widget.submitBtnLoading = true;

                Timer(Duration(seconds: 1), () {
                  widget.submitBtnLoading = false;
                });

                if (widget.loadingSecond != 0) {
                  _timer = Timer(Duration(seconds: widget.loadingSecond), () {
                    widget.btnController.reset();
                  });
                }

                FocusScope.of(context).requestFocus(new FocusNode());
                if (widget.onPressed != null) widget.onPressed!();
              },
        controller: widget.btnController,
        child: Text(widget.label,
            style: FXUI.normalTextStyle
                .copyWith(fontSize: widget.fontSize, color: widget.textColor)));
  }
}
