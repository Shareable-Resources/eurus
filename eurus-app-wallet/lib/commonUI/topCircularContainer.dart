import 'package:flutter/material.dart';

class TopCircularContainer extends StatelessWidget {
  final Widget child;
  final double? height;
  final double? width;

  TopCircularContainer({
    this.height,
    this.width,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      width: width,
      color: Colors.transparent,
      child: new Container(
        decoration: new BoxDecoration(
            color: Colors.white,
            borderRadius: new BorderRadius.only(
              topLeft: const Radius.circular(20.0),
              topRight: const Radius.circular(20.0),
            )),
        child: child,
      ),
    );
  }
}
