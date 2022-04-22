import 'package:euruswallet/common/commonMethod.dart';
import 'package:easy_localization/easy_localization.dart';

class SettingAppBar extends StatelessWidget implements PreferredSizeWidget {
  SettingAppBar(this.backwardsCompatibility, {Key? key})
      : preferredSize = Size.fromHeight(kToolbarHeight + 13),
        super(key: key);

  final bool backwardsCompatibility;

  @override
  final Size preferredSize;

  @override
  Widget build(BuildContext context) {
    return AppBar(
      centerTitle: true,
      title: Text(
        'SETTING_PAGE.MAIN_TITLE'.tr(),
        style: FXUI.normalTextStyle
            .copyWith(color: Colors.black, fontWeight: FontWeight.bold),
      ),
      backgroundColor: Colors.white,
      shadowColor: Color(0x99b0b0b0),
      leading: !backwardsCompatibility
          ? Container()
          : IconButton(
              icon: Icon(
                Icons.arrow_back_ios_outlined,
                color: FXColor.deepGrayColor,
              ),
              onPressed: () {
                Navigator.of(context).pop();
              },
            ),
      elevation: 5,
    );
  }
}
