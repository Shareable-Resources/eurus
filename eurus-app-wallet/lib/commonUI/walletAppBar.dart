import 'package:euruswallet/commonUI/constant.dart';
import 'package:flutter/material.dart';

class WalletAppBar extends StatefulWidget implements PreferredSizeWidget {
  WalletAppBar({
    Key? key,
    required this.title,
    bool? backButton,
    this.rightWidget,
    this.function,
  })  : this.backButton = backButton ?? true,
        this.preferredSize = Size.fromHeight(kToolbarHeight),
        super(key: key);

  @override
  final Size preferredSize; // default is 56.0
  final String title;
  final bool backButton;
  final Widget? rightWidget;
  final VoidCallback? function;

  @override
  _WalletAppBarState createState() => _WalletAppBarState();
}

class _WalletAppBarState extends State<WalletAppBar> {
  @override
  Widget build(BuildContext context) {
    return AppBar(
      automaticallyImplyLeading: widget.backButton,
      centerTitle: true,
      backgroundColor: Colors.transparent,
      elevation: 0.0,
      title: Text(
        widget.title,
        style: FXUI.normalTextStyle.copyWith(color: Colors.white),
      ),
      actions: <Widget>[
        TextButton(
          style: TextButton.styleFrom(
            primary: Colors.white,
            minimumSize: Size(88, 36),
            padding: EdgeInsets.symmetric(horizontal: 16.0),
            shape: CircleBorder(
              side: BorderSide(
                color: Colors.transparent,
              ),
            ),
          ),
          onPressed: widget.function,
          child: widget.rightWidget ?? Container(),
        ),
      ],
    );
  }
}
