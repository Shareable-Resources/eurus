import 'package:euruswallet/common/commonMethod.dart';

class BackGroundImage extends StatelessWidget {
  final Widget child;
  final CurrentUserType? currentUserType;
  BackGroundImage({
    required this.child,
    currentUserType,
  }) : this.currentUserType =
            currentUserType == null ? common.currentUserType : currentUserType;

  @override
  Widget build(BuildContext context) {
    return Stack(children: <Widget>[
      Image.asset(
        currentUserType == CurrentUserType.centralized
            ? "images/backgroundImage.png"
            : "images/backgroundImage2.png",
        package: 'euruswallet',
        height: MediaQuery.of(context).size.height,
        width: MediaQuery.of(context).size.width,
        fit: BoxFit.cover,
      ),
      child
    ]);
  }
}
