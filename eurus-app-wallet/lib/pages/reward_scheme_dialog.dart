import 'package:euruswallet/common/commonMethod.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:easy_localization/easy_localization.dart';

class RewardSchemeDialog extends StatelessWidget {
  const RewardSchemeDialog({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
      onWillPop: () async => false,
      child: Dialog(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Stack(
              children: [
                getRewardSchemeImage(context: context),
                Positioned(
                  top: 11.0,
                  right: 8.0,
                  child: GestureDetector(
                    onTap: () => Navigator.of(context).pop(),
                    child: CircleAvatar(
                      radius: 14.0,
                      backgroundColor: Colors.black54,
                      child: Icon(Icons.close, color: Colors.white),
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Image getRewardSchemeImage({required BuildContext context}) {
    final locale = context.locale.toString();
    final imageName = locale == 'zh_Hant'
        ? 'reward_scheme_type_0_tw.png'
        : locale == 'zh_Hans'
            ? 'reward_scheme_type_0_cn.png'
            : 'reward_scheme_type_0_en.png';
    return Image.asset(
      'images/$imageName',
      package: 'euruswallet',
    );
  }
}
