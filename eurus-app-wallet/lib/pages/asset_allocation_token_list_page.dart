import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
import '../common/commonMethod.dart';
import '../commonUI/crypto_asset_list_widget.dart';
import '../commonUI/topSelectBlockChainBar.dart';
import '../pages/decentralized/edit_assets_list.dart';
import 'decentralized/decentralized_wallet_base_page.dart';

class AssetAllocationTokenListPage extends StatefulWidget {
  const AssetAllocationTokenListPage({Key? key}) : super(key: key);

  @override
  _AssetAllocationTokenListPageState createState() =>
      _AssetAllocationTokenListPageState();
}

class _AssetAllocationTokenListPageState
    extends State<AssetAllocationTokenListPage> {
  late Future<String?> _getAssetsListFuture;
  int _viewSelectedChainCode =
      common.currentBlockchainSelection == BlockChainType.Ethereum ? 0 : 1;

  @override
  void initState() {
    _getAssetsListFuture =
        NormalStorageKit().readValue('assetsList_${common.encryptedAddress}');
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return DecentralizedWalletBasePage(
      appBarTitle: Text("ASSET_ALLOCATION_PAGE.TITLE".tr()),
      body: Padding(
        padding: EdgeInsets.only(top: 0),
        child: Column(
          children: [
            ListTile(
              title: Text(
                  "${_viewSelectedChainCode <= 0 ? 'ASSET_ALLOCATION_PAGE.DEPOSIT_TITLE'.tr() : 'ASSET_ALLOCATION_PAGE.WITHDRAWAL_TITLE'.tr()}",
                  textAlign: TextAlign.center,
                  style: Theme.of(context)
                      .textTheme
                      .headline6
                      ?.apply(fontWeightDelta: 1, color: Colors.black)),
              dense: true,
            ),
            Expanded(
              child: Padding(
                padding: EdgeInsets.only(top: 10 + 16.0),
                child: Container(
                  clipBehavior: Clip.hardEdge,
                  decoration: BoxDecoration(
                      borderRadius: BorderRadiusDirectional.only(
                          topStart: Radius.circular(30),
                          topEnd: Radius.circular(30)),
                      color: Colors.white,
                      boxShadow: kElevationToShadow[3]),
                  child: Padding(
                    padding: EdgeInsets.only(top: 0, left: 0, right: 0),
                    child: Column(
                      children: [
                        Expanded(
                          child: Stack(
                            children: [
                              // asset list section content body
                              // SafeArea(child:
                              FutureBuilder(
                                future: _getAssetsListFuture,
                                builder:
                                    (_, AsyncSnapshot<String?> s) =>
                                        s.connectionState !=
                                                ConnectionState.done
                                            ? Center(
                                                child: CircularProgressIndicator
                                                    .adaptive())
                                            : !s.hasData ||
                                                    !(s.data is String) ||
                                                    s.data is String &&
                                                        isEmptyString(
                                                            string: s.data
                                                                as String)
                                                ? ListView(
                                                    padding: EdgeInsets.only(
                                                        top: 57,
                                                        left: 8,
                                                        right: 8,
                                                        bottom:
                                                            SCREEN_WITH_BOTTOM_NAV_TAB_BAR_SAFE_AREA_BOTTOM_PADDING_CONTENT_INSET),
                                                    children: [
                                                      Container(
                                                          child: ListTile(
                                                              title:
                                                                  CupertinoButton(
                                                                      onPressed:
                                                                          () async {
                                                                        await Navigator.of(_)
                                                                            .push(MaterialPageRoute(
                                                                          builder: (_) =>
                                                                              EditAssetsListPage(userSuffix: common.encryptedAddress ?? ''),
                                                                        ));
                                                                        setState(
                                                                            () {});
                                                                      },
                                                                      child: Text(
                                                                          "COMMON.ADD_TOKEN"
                                                                              .tr(),
                                                                          style: Theme.of(_)
                                                                              .textTheme
                                                                              .button
                                                                              ?.apply(color: Theme.of(_).colorScheme.primary)))),
                                                          margin: EdgeInsets.all(4),
                                                          decoration: BoxDecoration(
                                                            shape: BoxShape
                                                                .rectangle,
                                                            boxShadow:
                                                                kElevationToShadow[
                                                                    3],
                                                            color: Colors.white,
                                                            borderRadius:
                                                                BorderRadius
                                                                    .circular(
                                                                        15),
                                                          )),
                                                    ],
                                                    // );}) : Builder(builder: (_) {final data = (jsonDecode(s.data) as List).where((e) => e['showAssets'] && e['address${_viewSelectedChainCode <= 0 ? TEST_NET??'Ethereum' : 'Eurus'}'] != null).map((e) {e['address'] = e['address${_viewSelectedChainCode <= 0 ? TEST_NET??'Ethereum' : 'Eurus'}']??'0xb9f5377a3f1c4748c71bd4ce2a16f13eecc4fec4'; return e;}).toList(); return ListView.builder(
                                                  )
                                                : CryptoAssetListWidget(
                                                    data:
                                                        jsonDecode(s.data ?? '')
                                                            as List,
                                                    addressSuffix:
                                                        _viewSelectedChainCode <=
                                                                0
                                                            ? TEST_NET
                                                            : 'Eurus',
                                                    additionalListViewBottomPaddingContentInset:
                                                        0,
                                                    onTapHandler: (_data) async =>
                                                        await common
                                                            .navigateToAssetAllocationTransfer(
                                                                context, _data),
                                                    additionalWhereClauseFilterTest:
                                                        (e) =>
                                                            e['address$TEST_NET']
                                                                is String &&
                                                            (e['address$TEST_NET']
                                                                    as String)
                                                                .isNotEmpty &&
                                                            e['addressEurus']
                                                                is String &&
                                                            (e['addressEurus']
                                                                    as String)
                                                                .isNotEmpty,
                                                  ),
                              ),
                              // bottom: false ,),
                              // asset list section header
                              Align(
                                alignment: Alignment.topCenter,
                                child: Container(
                                  decoration: BoxDecoration(
                                    gradient: LinearGradient(
                                      begin: Alignment.center,
                                      end: Alignment.bottomCenter,
                                      colors: [Colors.white, Colors.white],
                                    ),
                                  ),
                                  child: Padding(
                                    padding: EdgeInsets.symmetric(horizontal: 8)
                                        .add(EdgeInsets.only(top: 5)),
                                    child: IntrinsicHeight(
                                      child: Row(
                                        children: [
                                          TopSelectBlockChainBar(
                                              dropDownList: [
                                                BlockChainType.Eurus,
                                                BlockChainType.Ethereum
                                              ],
                                              onSegmentChosen:
                                                  (BlockChainType type) {
                                                setState(() {
                                                  if (type ==
                                                      BlockChainType.Eurus) {
                                                    _viewSelectedChainCode = 1;
                                                  } else if (type ==
                                                          BlockChainType
                                                              .Ethereum ||
                                                      type ==
                                                          BlockChainType
                                                              .BinanceCoin) {
                                                    _viewSelectedChainCode =
                                                        !isEmptyString(
                                                                string:
                                                                    TEST_NET)
                                                            ? -1
                                                            : 0;
                                                  }
                                                  common.topSelectedBlockchainType =
                                                      type;
                                                });
                                              },
                                              currentSelection: common
                                                  .currentBlockchainSelection,
                                              topBarType:
                                                  TopSelectBarType.enable),
                                          SizedBox(
                                            height: 48,
                                          ),
                                        ],
                                      ),
                                    ),
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
