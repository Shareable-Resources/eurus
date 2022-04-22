import 'dart:ui';

import 'package:euruswallet/common/commonMethod.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';

enum CtaButtonType { primary, secondary }

class CtaButton extends StatefulWidget {
  const CtaButton({
    Key? key,
    required this.type,
    required this.onPressed,
    required this.text,
    this.borderRadius,
    this.textStyle,
  }) : super(key: key);

  final CtaButtonType type;
  final void Function()? onPressed;
  final String text;
  final borderRadius;
  final textStyle;

  @override
  _CtaButtonState createState() => _CtaButtonState();
}

class _CtaButtonState extends State<CtaButton> {
  @override
  Widget build(BuildContext context) {
    final foregroundColor = widget.type == CtaButtonType.primary
        ? Colors.white
        : common.getBackGroundColor();
    final backgroundColor = widget.type == CtaButtonType.primary
        ? common.getBackGroundColor()
        : Colors.white;
    return TextButton(
      onPressed: widget.onPressed,
      style: TextButton.styleFrom(
        shape: RoundedRectangleBorder(
          borderRadius: widget.borderRadius ?? BorderRadius.circular(14),
          side: widget.type == CtaButtonType.primary
              ? BorderSide.none
              : BorderSide(color: foregroundColor),
        ),
        textStyle: widget.textStyle,
        padding: EdgeInsets.symmetric(vertical: 12.0, horizontal: 24.0),
        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
      ).copyWith(
        foregroundColor: MaterialStateProperty.all(foregroundColor),
        backgroundColor: MaterialStateProperty.all(backgroundColor),
      ),
      child: Text(widget.text),
    );
  }
}
