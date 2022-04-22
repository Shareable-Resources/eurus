import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/pages/decentralized/transferPage.dart';
import 'package:easy_localization/easy_localization.dart';

class SelectTargetDetailPage extends StatefulWidget {
  final String titleName;

  SelectTargetDetailPage({
    Key? key,
    required this.titleName,
  }) : super(key: key);

  @override
  _SelectTargetDetailPageState createState() => _SelectTargetDetailPageState();
}

class _SelectTargetDetailPageState extends State<SelectTargetDetailPage> {
  RoundedLoadingButtonController btnController = RoundedLoadingButtonController();

  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    SizeConfig(context: context);
    return BackGroundImage(
        child: Scaffold(
      backgroundColor: Colors.transparent,
      appBar: WalletAppBar(title: widget.titleName),
      body: SingleChildScrollView(
          child: TopCircularContainer(
              height: size.heightWithoutAppBar,
              width: size.blockSizeHorizontal * 100,
              child: Padding(
                padding: EdgeInsets.only(
                    left: size.leftPadding, right: size.leftPadding),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Padding(
                      padding: EdgeInsets.only(top: 36, bottom: 38),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Center(
                            child: Text('SEND_CONFIRM_PAGE.CONFIRM_INFO'.tr(),
                                style:
                                    FXUI.titleTextStyle.copyWith(fontSize: 24)),
                          ),
                          Padding(
                            padding: EdgeInsets.only(top: 10),
                            child: Text('SEND_PAGE.PAGE_DESC'.tr(),
                                style: FXUI.normalTextStyle.copyWith(
                                    fontSize: 14, color: FXColor.lightGray)),
                          ),
                        ],
                      ),
                    ),
                    if (common.findEmailWalletAddress?.data?.userType == 1)
                      Text("SEND_CONFIRM_PAGE.EMAIL".tr(),
                          style: FXUI.normalTextStyle.copyWith(
                              fontSize: 14, color: FXColor.lightGray)),
                    if (common.findEmailWalletAddress?.data?.userType == 1)
                      Padding(
                        padding: EdgeInsets.only(top: 29, left: 16),
                        child: Text(
                            common.findEmailWalletAddress?.data?.email ?? '',
                            style: FXUI.normalTextStyle.copyWith(fontSize: 14)),
                      ),
                    Padding(
                      padding: EdgeInsets.only(top: 35),
                      child: Text("SEND_CONFIRM_PAGE.WALLET_ADDRESS".tr(),
                          style: FXUI.normalTextStyle.copyWith(
                              fontSize: 14, color: FXColor.lightGray)),
                    ),
                    Padding(
                      padding: EdgeInsets.only(top: 29, left: 16),
                      child: Text(
                          common.findEmailWalletAddress?.data?.walletAddress ??
                              common.targetAddress ?? "",
                          style: FXUI.normalTextStyle.copyWith(fontSize: 14)),
                    ),
                    Padding(
                      padding: EdgeInsets.only(top: 47),
                      child: SubmitButton(
                        btnController: btnController,
                        label: 'COMMON.CONFIRM'.tr(),
                        onPressed: () async {
                          btnController.reset();
                          common.pushPage(
                              page: TransferPage(
                                  titleName: "SEND_PAGE.MAIN_TITLE".tr(),
                                  fromBlockChainType:
                                      common.fromBlockChainType),
                              context: context);
                        },
                      ),
                    )
                  ],
                ),
              ))),
    ));
  }
}
