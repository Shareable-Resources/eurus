import 'package:euruswallet/common/commonMethod.dart';

class CardBtnRow extends StatelessWidget {
  CardBtnRow(
    this.title, {
    this.subTitle,
    this.btnContent,
    this.onPressFnc,
    this.borderBtm = false,
    this.customBtmContent,
  });

  final String title;
  final String? subTitle;
  final Widget? btnContent;
  final Function()? onPressFnc;
  final bool borderBtm;
  final Widget? customBtmContent;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onPressFnc ?? () {},
      child: Container(
        decoration: BoxDecoration(
          border: Border(
            bottom: borderBtm == true
                ? BorderSide(width: 1, color: FXColor.verylightBlack)
                : BorderSide.none,
          ),
        ),
        padding: EdgeInsets.symmetric(vertical: 20),
        child: Column(
          children: [
            Row(
              children: [
                Expanded(
                  flex: 2,
                  child: Container(
                    padding: EdgeInsets.only(left: 30, right: 15),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          title,
                          style:
                          FXUI.normalTextStyle.copyWith(fontSize: 15, color: FXColor.mediumGrayColor),
                        ),
                        Padding(
                          padding: EdgeInsets.only(
                            top: subTitle != null ? 3 : 0,
                          ),
                        ),
                        subTitle != null
                            ? Text(
                                subTitle ?? '',
                                style: FXUI.normalTextStyle.copyWith(
                                  fontSize: 12,
                                  color: FXColor.mediumGrayColor,
                                ),
                              )
                            : Container()
                      ],
                    ),
                  ),
                ),
                Container(
                  padding: EdgeInsets.only(right: 32),
                  child: btnContent != null
                      ? btnContent
                      : onPressFnc != null
                          ? Icon(
                              Icons.navigate_next_outlined,
                              color: FXColor.untabColor,
                            )
                          : Container(),
                )
              ],
            ),
            customBtmContent ?? Container()
          ],
        ),
      ),
    );
  }
}
