import 'package:euruswallet/common/commonMethod.dart';
import 'package:pin_code_fields/pin_code_fields.dart';

class PinCode extends StatelessWidget {
  final void Function(String) onChanged;
  final void Function(String)? onCompleted;
  final Function? onTap;
  final StreamController<ErrorAnimationType>? errorController;
  final TextEditingController? paymentPasswordController;
  final FocusNode? myFocusNode;
  final bool hasError;

  PinCode({
    required this.onChanged,
    this.onCompleted,
    this.onTap,
    this.errorController,
    this.paymentPasswordController,
    this.myFocusNode,
    this.hasError = false,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
        width: size.blockSizeHorizontal * 100,
        height: 80,
        child: Card(
          elevation: 3,
          shape: RoundedRectangleBorder(
            borderRadius: FXUI.cricleRadius,
          ),
          margin: EdgeInsets.zero,
          color: Colors.white,
          child: Padding(
              padding: EdgeInsets.only(
                  left: size.leftPadding * 2, right: size.leftPadding * 2),
              child: PinCodeTextField(
                inputFormatters: <TextInputFormatter>[
                  FilteringTextInputFormatter.allow(RegExp(r'[0-9]')),
                ],
                keyboardType: TextInputType.numberWithOptions(signed: true, decimal: true),
                onTap: onTap,
                focusNode: myFocusNode,
                length: 6,
                obscureText: false,
                animationType: AnimationType.fade,
                cursorColor: Colors.transparent,
                pinTheme: PinTheme(
                  activeColor: hasError ? Colors.red : Colors.black,
                  selectedColor: Colors.black,
                  inactiveColor: FXColor.lightGray,
                  inactiveFillColor: FXColor.lightGray,
                  selectedFillColor: FXColor.lightGray,
                  disabledColor: FXColor.lightGray,
                  shape: PinCodeFieldShape.underline,
                  fieldHeight: 50,
                  fieldWidth: 20,
                  activeFillColor: hasError ? Colors.red : Colors.white,
                ),
                animationDuration: Duration(milliseconds: 300),
                backgroundColor: Colors.transparent,
                enableActiveFill: false,
                errorAnimationController: errorController,
                controller: paymentPasswordController,
                onCompleted: onCompleted,
                onChanged: onChanged,
                appContext: context,
              )),
        ));
  }
}
