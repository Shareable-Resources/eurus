import 'package:euruswallet/common/commonMethod.dart';
import 'package:flutter/cupertino.dart';
import 'package:easy_localization/easy_localization.dart';

class Constants {
  Constants._();
  static const double padding = 28;
}

class CustomConfirmDialog extends StatelessWidget {
  final String title, descriptions, buttonText;
  final Widget? icon;
  final Color? btnColor;
  final Function? btnHandler;
  final bool dismissIble;
  final bool showCancelButton;
  final bool useSubmitButton;

  CustomConfirmDialog({
    Key? key,
    this.title: '',
    this.descriptions: '',
    this.buttonText: '',
    this.icon,
    this.btnColor,
    this.btnHandler,
    this.dismissIble: true,
    this.showCancelButton: false,
    this.useSubmitButton: true,
  }) : super(key: key);

  @override
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

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
              isEmptyString(string: title)
                  ? Container()
                  : Text(
                      title,
                      style: FXUI.normalTextStyle
                          .copyWith(fontSize: 18, fontWeight: FontWeight.bold),
                    ),
              icon != null
                  ? Padding(
                      padding:
                          EdgeInsets.symmetric(vertical: Constants.padding),
                      child: icon,
                    )
                  : Padding(padding: EdgeInsets.only(top: Constants.padding)),
              isEmptyString(string: descriptions)
                  ? Container()
                  : Padding(
                      padding: EdgeInsets.only(bottom: Constants.padding),
                      child: Text(
                        descriptions,
                        style: FXUI.normalTextStyle
                            .copyWith(fontSize: 14, color: FXColor.textGray),
                        textAlign: TextAlign.center,
                      ),
                    ),
              useSubmitButton == true
                  ? SubmitButton(
                      btnController: btnController,
                      label: buttonText,
                      buttonBGColor: btnColor ?? common.getBackGroundColor(),
                      onPressed: () {
                        if (dismissIble) {
                          Navigator.pop(context);
                        }
                        if (btnHandler != null) btnHandler!();
                        btnController.reset();
                      },
                    )
                  : CustomTextButton(
                      buttonBGColor: btnColor ?? common.getBackGroundColor(),
                      text: buttonText,
                      onPressed: () {
                        if (dismissIble) {
                          Navigator.pop(context);
                        }
                        if (btnHandler != null) btnHandler!();
                      }),
              showCancelButton == true
                  ? Padding(
                      padding: const EdgeInsets.only(top: 16.0),
                      child: CustomTextButton(
                          buttonBGColor: FXColor.cancelGrayButton,
                          text: 'COMMON.CANCEL'.tr(),
                          onPressed: () {
                            Navigator.pop(context);
                          }),
                    )
                  : Container(),
            ],
          )),
    );
  }
}
