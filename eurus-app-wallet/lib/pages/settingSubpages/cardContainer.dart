import 'package:euruswallet/common/commonMethod.dart';

class CardContainer extends StatelessWidget {
  CardContainer(
    this.title,
    this.content, {
    this.themeColor,
        this.padding,
        this.titleWidget
  });

  final String title;
  final Widget content;
  final Color? themeColor;
  final EdgeInsets? padding;
  final Widget? titleWidget;

  @override
  Widget build(BuildContext context) {
    var _themeColor = themeColor != null
        ? themeColor
        : !isCentralized()
            ? FXColor.mainDeepBlueColor
            : FXColor.mainBlueColor;

    return Container(
      margin: padding ?? EdgeInsets.all(12),
      child: Column(
        children: [
          titleWidget ?? Container(
            width: double.infinity,
            padding: EdgeInsets.symmetric(horizontal: 21, vertical: 18),
            child: Text(
              title,
              textAlign: TextAlign.left,
              style: FXUI.normalTextStyle.copyWith(
                fontWeight: FontWeight.w600,
                color: _themeColor,
                fontSize: 16,
              ),
            ),
          ),
          Container(
            width: double.infinity,
            decoration: BoxDecoration(
              borderRadius: FXUI.cricleRadius,
              boxShadow: [
                BoxShadow(
                  color: FXColor.grey80Color,
                  offset: Offset(1, 2),
                  blurRadius: 8,
                )
              ],
              color: Colors.white,
            ),
            child: content,
          )
        ],
      ),
    );
  }
}
