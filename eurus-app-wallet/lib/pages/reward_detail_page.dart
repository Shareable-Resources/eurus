import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/cta_button.dart';
import 'package:euruswallet/model/user_marketing_reward_list_response_model.dart';
import 'package:euruswallet/model/user_marketing_reward_scheme_response_model.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:flutter_neumorphic/flutter_neumorphic.dart';
import 'package:intl/intl.dart';
import 'package:url_launcher/url_launcher.dart';

class RewardDetailPage extends StatefulWidget {
  const RewardDetailPage({
    Key? key,
    required this.scheme,
  }) : super(key: key);

  final UserMarketingRewardScheme scheme;

  @override
  _RewardDetailPageState createState() => _RewardDetailPageState();
}

class _RewardDetailPageState extends State<RewardDetailPage> {
  UserMarketingRewardSchemeContent? schemeContent;

  UserMarkingRewardList? reward;

  @override
  void initState() {
    super.initState();
    reward = (common.rewardedList ?? []).isNotEmpty
        ? (common.rewardedList ?? []).first
        : null;
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();

    final locale = context.locale.toString();
    if (locale == 'zh_Hant') {
      schemeContent = widget.scheme.zhTw;
    } else if (locale == 'zh_Hans') {
      schemeContent = widget.scheme.zhCn;
    } else {
      schemeContent = widget.scheme.en;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: FXColor.lightWhiteColor,
      appBar: AppBar(
        centerTitle: true,
        title: Text(
          'REWARD.TITLE'.tr(),
          style: FXUI.normalTextStyle.copyWith(
            color: FXColor.middleBlack,
            fontWeight: FontWeight.bold,
          ),
        ),
        leading: IconButton(
          icon: Icon(Icons.arrow_back_ios_outlined,
              color: common.getBackGroundColor()),
          onPressed: () => Navigator.of(context).pop(),
        ),
        backgroundColor: FXColor.lightWhiteColor,
        shadowColor: Colors.transparent,
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 24.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              if (common.getBannerPath(context: context) != null)
                Neumorphic(
                  child: Image.asset(
                    common.getBannerPath(context: context)!,
                    package: 'euruswallet',
                  ),
                  style: FXUI.neumorphicBannerImage,
                ),
              SizedBox(height: 20),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisSize: MainAxisSize.max,
                children: [
                  Padding(
                    padding: const EdgeInsets.only(left: 15.0, right: 113.0),
                    child: Text('REWARD.STATUS'.tr()),
                  ),
                  Expanded(
                    child: CtaButton(
                      type: reward == null
                          ? CtaButtonType.secondary
                          : CtaButtonType.primary,
                      borderRadius: BorderRadius.circular(8.0),
                      text: reward == null
                          ? 'REWARD.ONGOING'.tr()
                          : 'REWARD.COMPLETED'.tr(),
                      onPressed: null,
                    ),
                  ),
                ],
              ),
              SizedBox(height: 13),
              Expanded(
                child: Container(
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(15.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        if (!isEmptyString(string: schemeContent?.title))
                          Padding(
                            padding: const EdgeInsets.only(bottom: 8.0),
                            child: Text(schemeContent?.title ?? ''),
                          ),
                        if (!isEmptyString(string: schemeContent?.details))
                          Padding(
                            padding: const EdgeInsets.only(bottom: 8.0),
                            child: Text(schemeContent?.details ?? ''),
                          ),
                        SizedBox(height: 8.0),
                        Divider(),
                        getInfoItemRow(
                          title: 'REWARD.REWARDED_DATE'.tr(),
                          content: reward?.createdDate != null
                              ? DateFormat('yyyy/MM/dd')
                                  .format((reward!.createdDate)!)
                                  .toString()
                              : '-',
                        ),
                        Divider(),
                        getInfoItemRow(
                          title: 'REWARD.REWARDED_EUN'.tr(),
                          content: reward?.amount == null
                              ? '-'
                              : common.numberFormat(
                                  number: common
                                      .divisionDecimal(18, reward?.amount ?? 0)
                                      .toString()),
                        ),
                        Divider(),
                        getInfoItemRow(
                          title: 'REWARD.REWARDED_TX_ID'.tr(),
                          content: reward?.txHash ?? '-',
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              SizedBox(height: 14),
              CtaButton(
                type: CtaButtonType.primary,
                borderRadius: BorderRadius.circular(14.0),
                text: 'REWARD.LEARN_MORE'.tr(),
                onPressed: () => launch(
                  'REWARD.LEARN_MORE_LINK'.tr(),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget getInfoItemRow({
    required String title,
    String? content,
  }) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: Text(
        title,
        style: FXUI.hintStyle.copyWith(
          fontWeight: FontWeight.normal,
          color: FXColor.placeholderGreyColor,
        ),
      ),
      title: Text(
        content ?? '-',
        textAlign: TextAlign.right,
        style: FXUI.inputStyle.copyWith(
          fontSize: 14,
          fontWeight: FontWeight.normal,
        ),
      ),
    );
  }
}
