import 'package:app_crypto_icons/app_crypto_icons.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/reorderable_list_simple.dart';
import 'package:euruswallet/model/crypto_currency_model.dart';
import 'package:euruswallet/pages/decentralized/add_token.dart';
import 'package:flutter/cupertino.dart';

class EditAssetsListPage extends StatefulWidget {
  EditAssetsListPage({
    required this.userSuffix,
    Key? key,
  }) : super(key: key);

  final String userSuffix;

  _EditAssetsListPageState createState() => _EditAssetsListPageState();
}

class _EditAssetsListPageState extends State<EditAssetsListPage> {
  List<CryptoCurrencyModel> _cryptos = [];

  bool get _cryptosInited => _cryptos.length > 0;
  Color get _themeColor => !isCentralized()
      ? FXColor.assetsListPurpleColor
      : FXColor.assetsListBlueColor;

  @override
  void initState() {
    _initAssetsList();

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: FXColor.veryLightGreyTextColor,
      appBar: AppBar(
        centerTitle: true,
        title: Text(
          'EDIT_ASSETS_LIST_PAGE.MAIN_TITLE'.tr(),
          style: FXUI.normalTextStyle
              .copyWith(color: Colors.black, fontWeight: FontWeight.bold),
        ),
        leading: IconButton(
          icon:
              Icon(Icons.arrow_back_ios_outlined, color: FXColor.deepGrayColor),
          onPressed: () {
            Navigator.of(context).pop();
          },
        ),
        backgroundColor: Colors.transparent,
        shadowColor: Colors.transparent,
        actions: [
          IconButton(
            icon: Icon(Icons.add, color: _themeColor),
            onPressed: () async {
              common
                  .pushPage(
                      page: AddTokenPage(userSuffix: widget.userSuffix),
                      context: context)
                  .then((value) => _initAssetsList());
            },
          ),
        ],
      ),
      body: Stack(
        children: [
          Container(
            child: _cryptosInited
                ? Container()
                : Center(
                    child: SizedBox(
                      width: 45,
                      height: 45,
                      child: Center(
                        child: CircularProgressIndicator(
                          valueColor: AlwaysStoppedAnimation<Color>(
                            FXColor.mainDeepBlueColor,
                          ),
                        ),
                      ),
                    ),
                  ),
          ),
          Container(
            child: ReorderableListSimple(
              handleIcon: _customDragIcon(),
              handleSide: ReorderableListSimpleSide.Left,
              onReorder: (oldIdx, newIdx) {
                final CryptoCurrencyModel item = _cryptos[oldIdx];

                setState(() {
                  _cryptos.remove(item);
                  _cryptos.insert(newIdx, item);
                });

                _updateAssetsList();
              },
              children: _cryptos
                  .map(
                    (item) => ListTile(
                      key: Key(item.hashCode.toString()),
                      // leading: item.iconWidget,
                      leading: common.getIcon(
                        item.symbol,
                        50,
                        source: item.iconSource ?? IconSourceType.svg,
                        imgUrl: item.imgUrl,
                      ),
                      title: Text('${item.symbol}'),
                      subtitle: Text('${item.currency}'),
                      trailing: isAssetEditable(item)
                          ? StatefulBuilder(
                              builder: (context, setState) => CupertinoSwitch(
                                value: item.showAssets,
                                onChanged: (v) {
                                  setState(() {
                                    item.showAssets = v;
                                  });
                                  _updateAssetsList();
                                },
                                activeColor: _themeColor,
                              ),
                            )
                          : null,
                      contentPadding: EdgeInsets.only(right: 19, left: 8),
                    ),
                  )
                  .toList(),
            ),
          ),
        ],
      ),
    );
  }

  void _initAssetsList() async {
    String? alistString =
        await NormalStorageKit().readValue('assetsList_${widget.userSuffix}');
    List<CryptoCurrencyModel> _finalList = isEmptyString(string: alistString)
        ? []
        : (jsonDecode(alistString ?? '') as List)
            .map((e) => CryptoCurrencyModel.fromJson(e as Map<String, dynamic>))
            .toList();

    setState(() {
      _cryptos = _finalList;
    });
  }

  void _updateAssetsList() async {
    List<CryptoCurrencyModel> _orderedList = [];

    int idx = 0;
    _cryptos.forEach((e) {
      _orderedList.add(e..order = idx);
      idx++;
    });
    String rdString = jsonEncode(_cryptos);

    await NormalStorageKit()
        .setValue(rdString, 'assetsList_${widget.userSuffix}');
  }

  Widget _customDragIcon() {
    var _greyLine = DecoratedBox(
      decoration: BoxDecoration(color: FXColor.greyTextColor),
      child: SizedBox(
        width: double.infinity,
        height: 2,
      ),
    );

    return Container(
      margin: EdgeInsets.only(right: 8),
      color: Colors.transparent,
      width: 20,
      height: 16,
      child: Column(
        mainAxisSize: MainAxisSize.max,
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [_greyLine, _greyLine, _greyLine],
      ),
    );
  }
}
