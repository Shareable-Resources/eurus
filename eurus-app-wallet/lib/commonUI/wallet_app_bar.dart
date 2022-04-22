import 'package:euruswallet/common/commonMethod.dart';

class WalletAppBarV2 extends StatefulWidget implements PreferredSizeWidget {
  WalletAppBarV2({
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
  _WalletAppBarV2State createState() => _WalletAppBarV2State();
}

Color? getBackButtonColor() {
  if (common.currentUserProfile?.userType == CurrentUserType.centralized)
    return FXColor.mainBlueColor;
  else if (common.currentUserProfile?.userType == CurrentUserType.decentralized)
    return FXColor.mainDeepBlueColor;
  else
    return null;
}

class _WalletAppBarV2State extends State<WalletAppBarV2> {
  @override
  Widget build(BuildContext context) {
    return AppBar(
      automaticallyImplyLeading: widget.backButton,
      centerTitle: true,
      backgroundColor: Colors.transparent,
      elevation: 0.0,
      leading: new IconButton(
        icon: Image.asset(
          'images/icon_back.png',
          width: 13,
          height: 24,
          package: 'euruswallet',
          color: getBackButtonColor(),
        ),
        onPressed: () => Navigator.pop(context),
      ),
      title: Text(widget.title, style: TextStyle(
          fontWeight: FontWeight.bold,
          fontSize: 16,
          fontFamily: 'packages/euruswallet/SFProDisplay',
          color: Color(0xFF4A4A4A))),
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
