import 'package:euruswallet/common/commonMethod.dart';
import 'package:flutter/cupertino.dart';

class Constants {
  Constants._();
  static const double padding = 28;
}

class CustomDialogBox extends StatefulWidget {
  final String title, descriptions, buttonText;
  final Widget? titleIcon;
  final Color? btnColor;
  final Function? btnHandler;
  final bool dismissIble;

  const CustomDialogBox({
    Key? key,
    this.title: '',
    this.descriptions: '',
    this.buttonText: '',
    this.titleIcon,
    this.btnColor,
    this.btnHandler,
    this.dismissIble: true,
  }) : super(key: key);

  @override
  _CustomDialogBoxState createState() => _CustomDialogBoxState();
}

class _CustomDialogBoxState extends State<CustomDialogBox> {
  @override
  RoundedLoadingButtonController btnController = RoundedLoadingButtonController();

  Widget build(BuildContext context) {
    return Dialog(
      shape: RoundedRectangleBorder(
        borderRadius: FXUI.cricleRadius,
      ),
      insetPadding: EdgeInsets.all(24),
      elevation: 0,
      backgroundColor: Colors.transparent,
      child: contentBox(context),
    );
  }

  contentBox(context) {
    return Container(
      padding: EdgeInsets.only(
          left: 0, top: Constants.padding, right: 0, bottom: Constants.padding),
      decoration: BoxDecoration(
        shape: BoxShape.rectangle,
        color: Colors.white,
        borderRadius: FXUI.cricleRadius,
      ),
      child: Padding(
          padding: EdgeInsets.only(left: 27, right: 27),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: <Widget>[
              Row(
                mainAxisSize: MainAxisSize.min,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  widget.titleIcon ?? Container(),
                  widget.titleIcon != null
                      ? Padding(
                          padding: EdgeInsets.only(right: 7),
                        )
                      : Container(),
                  isEmptyString(string: widget.title)
                      ? Container()
                      : Text(
                          widget.title,
                          style: FXUI.normalTextStyle.copyWith(
                              fontSize: 18, fontWeight: FontWeight.bold),
                        ),
                ],
              ),
              widget.titleIcon != null && isEmptyString(string: widget.title)
                  ? Padding(
                      padding: EdgeInsets.only(bottom: Constants.padding),
                    )
                  : Container(),
              isEmptyString(string: widget.descriptions)
                  ? Container()
                  : Padding(
                      padding: EdgeInsets.only(bottom: Constants.padding),
                      child: Text(
                        widget.descriptions,
                        style: FXUI.normalTextStyle.copyWith(fontSize: 14, color: FXColor.textGray),
                        textAlign: TextAlign.center,
                      ),
                    ),
              SubmitButton(
                btnController: btnController,
                label: widget.buttonText,
                buttonBGColor: widget.btnColor ?? common.getBackGroundColor(),
                onPressed: () {
                  if (widget.dismissIble) {
                    Navigator.pop(context);
                  }
                  if (widget.btnHandler != null) widget.btnHandler!();
                  btnController.reset();
                },
              ),
            ],
          )),
    );
  }
}
