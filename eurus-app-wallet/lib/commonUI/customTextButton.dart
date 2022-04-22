import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';

class CustomTextButton extends StatelessWidget {
  CustomTextButton(
      {Key? key,
      this.text: '',
      this.textColor: Colors.white,
      this.buttonBGColor: FXColor.mainBlueColor,
      this.fontSize: 16,
      this.height: 50,
      this.onPressed})
      : super(key: key);

  final String text;
  final Color? textColor;
  final Color buttonBGColor;
  final Function? onPressed;
  final double fontSize;
  final double height;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
        width: double.infinity,
        height: height,
        child: TextButton(
          style: TextButton.styleFrom(
            padding: EdgeInsets.all(3),
            backgroundColor: buttonBGColor,
            shape: RoundedRectangleBorder(
              borderRadius: FXUI.cricleRadius,
            ),
          ),
          child: Text(text,
              style: FXUI.normalTextStyle.copyWith(
                fontSize: fontSize,
                color: textColor,
              )),
          onPressed: () {
            onPressed!();
          },
        ));
  }
}
