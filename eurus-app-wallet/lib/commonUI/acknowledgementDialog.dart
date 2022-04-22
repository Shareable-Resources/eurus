import 'package:euruswallet/common/commonMethod.dart';
import 'package:flutter/cupertino.dart';

class AcknowledgementDialog extends StatefulWidget {
  const AcknowledgementDialog({
    Key? key,
    required this.statement,
    this.dontAskAgainText,
    required this.buttonText,
    this.buttonHandler,
    this.mainIcon,
  }) : super(key: key);

  final Widget? mainIcon;
  final String statement, buttonText;
  final Function? buttonHandler;
  final String? dontAskAgainText;

  @override
  _AcknowledgementDialogState createState() => _AcknowledgementDialogState();
}

class _AcknowledgementDialogState extends State<AcknowledgementDialog> {
  bool dontAskAgain = false;
  RoundedLoadingButtonController btnController =
      RoundedLoadingButtonController();

  @override
  Widget build(BuildContext context) {
    return Dialog(
      shape: RoundedRectangleBorder(
        borderRadius: FXUI.cricleRadius,
      ),
      insetPadding: EdgeInsets.all(8),
      elevation: 0,
      backgroundColor: Colors.transparent,
      child: _contentBox(context),
    );
  }

  Widget _contentBox(context) {
    return Container(
      decoration: BoxDecoration(
        shape: BoxShape.rectangle,
        color: Colors.white,
        borderRadius: FXUI.cricleRadius,
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              GestureDetector(
                onTap: () => Navigator.of(context).pop(dontAskAgain),
                child: Padding(
                  padding: EdgeInsets.only(top: 20, right: 20),
                  child: Icon(Icons.close, color: FXColor.textGray),
                ),
              )
            ],
          ),
          Padding(
            padding: EdgeInsets.only(left: 28, bottom: 28, right: 28),
            child: Column(
              children: [
                widget.mainIcon != null
                    ? Padding(
                        padding: EdgeInsets.only(bottom: 12),
                        child: widget.mainIcon)
                    : Container(),
                Padding(
                  padding: EdgeInsets.only(top: 24, bottom: 28),
                  child: Text(
                    widget.statement,
                    style:
                        FXUI.normalTextStyle.copyWith(color: FXColor.textGray),
                    textAlign: TextAlign.justify,
                  ),
                ),
                SubmitButton(
                  btnController: btnController,
                  label: widget.buttonText,
                  loadingSecond: 4,
                  onPressed: () {
                    btnController.reset();
                    Navigator.of(context).pop(dontAskAgain);
                    if (widget.buttonHandler != null) widget.buttonHandler!();
                  },
                ),
              ],
            ),
          ),
          if (widget.dontAskAgainText != null &&
              !isEmptyString(string: widget.dontAskAgainText))
            Padding(
              padding: EdgeInsets.only(bottom: 18, left: 14),
              child: Row(
                children: [
                  Checkbox(
                    activeColor: common.getBackGroundColor(),
                    value: dontAskAgain,
                    onChanged: (v) => setState(() {
                      if (v != null) dontAskAgain = v;
                    }),
                  ),
                  GestureDetector(
                    onTap: () => setState(() {
                      dontAskAgain = !dontAskAgain;
                    }),
                    child: Text(widget.dontAskAgainText ?? '',
                        style: FXUI.normalTextStyle
                            .copyWith(color: FXColor.textGray)),
                  ),
                ],
              ),
            ),
        ],
      ),
    );
  }
}
