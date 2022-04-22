import 'package:euruswallet/common/commonMethod.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/pages/settingSubpages/KYC/kycStatus.dart';
import 'package:euruswallet/pages/settingSubpages/KYC/uploadPhotoPage.dart';
import 'package:euruswallet/commonUI/searchbarList.dart';
import 'package:euruswallet/model/kYCCountryList.dart';

class SelectDocumentType extends StatefulWidget {
  const SelectDocumentType({Key? key}) : super(key: key);
  @override
  _SelectDocumentType createState() => _SelectDocumentType();
}

class _SelectDocumentType extends State<SelectDocumentType> {
  KYCCountry? country;
  String? searchText;
  List<dynamic> list = [];

  @override
  void initState() {
    getCountryList();
    super.initState();
  }

  Future<void> getCountryList() async {
    KYCCountryList? kYCCountryList = common.kYCCountryList;
    if (kYCCountryList == null) {
      kYCCountryList = await api.getKYCCountryList();
    }
    setState(() {
      list = kYCCountryList!.data.toList();
    });
  }

  List<dynamic> filterCountryList() {
    List<dynamic> filterList = common.kYCCountryList!.data.toList();
    if (isEmptyString(string: searchText) == false) {
      filterList = filterList
          .where((e) =>
              e.name.toLowerCase().contains(searchText?.toLowerCase()) ||
              e.code.toLowerCase().contains(searchText?.toLowerCase()) ||
              'COUNTRY.${e.code}'.tr().contains(searchText!))
          .toList();
    }
    return filterList;
  }

  String? getCountryName(KYCCountry? country) {
    if (country == null) return null;
    final localeCode = tr('COUNTRY.${country.code}');
    final name = localeCode.contains('.') ? country.name : localeCode;

    return name;
  }

  @override
  Widget build(BuildContext context) {
    double listHeight = size.heightWithoutAppBar / 4;
    return Column(children: [
      Padding(
        padding: const EdgeInsets.symmetric(horizontal: 60.0),
        child: Text(
          'KYC.SELECT_TYPE'.tr(),
          style: FXUI.normalTextStyle.copyWith(fontSize: 14),
          textAlign: TextAlign.center,
        ),
      ),
      Padding(
        padding: const EdgeInsets.only(top: 17.0),
      ),
      Column(
        children: <Widget>[
          SearchbarList(
              list: list,
              value: getCountryName(country),
              listMaxHeight: listHeight,
              hintText: 'KYC.SELECT_COUNTRY'.tr(),
              onChanged: (value) => setState(() {
                    searchText = value;
                    list = filterCountryList();
                  }),
              onSelect: (index) => setState(() {
                    country = list[index];
                  }),
              buildListItem: (index, isSelected) => Container(
                      child: Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.only(right: 10.0),
                        child: Text(
                            list[index].code.toUpperCase().replaceAllMapped(
                                RegExp(r'[A-Z]'),
                                (match) => String.fromCharCode(
                                    match.group(0)!.codeUnitAt(0) + 127397)),
                            style: TextStyle(fontSize: 22)),
                      ),
                      Flexible(
                          child: Text(
                        getCountryName(list[index]) ?? '',
                        style: FXUI.normalTextStyle.copyWith(
                            fontSize: 18,
                            color: isSelected == true
                                ? FXColor.mainBlueColor
                                : Colors.black),
                      )),
                    ],
                  ))),
          GenDocTypeButtons(
              disabled: country == null, countryCode: country?.code ?? ''),
        ],
      ),
    ]);
  }
}

class GenDocTypeButtons extends StatelessWidget {
  const GenDocTypeButtons(
      {Key? key, this.disabled: false, this.countryCode: ''})
      : super(key: key);
  final bool disabled;
  final String countryCode;

  void onTap(BuildContext context, IdentityType type) {
    CommonMethod().pushPage(
      page: UploadPhotoPage(identityType: type, countryCode: countryCode),
      context: context,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
      Padding(
        padding: const EdgeInsets.only(top: 30.0),
      ),
      GestureDetector(
          onTap: disabled
              ? null
              : () {
                  onTap(context, IdentityType.Passport);
                },
          child: selectDocmentType('images/passportIcon.png',
              'KYC.PASSPORT'.tr(), 'KYC.FACE_PHOTO_PAGE'.tr(), disabled)),
      (countryCode.toUpperCase() == "HK" || countryCode.toUpperCase() == "CN")
          ? Column(children: [
              Padding(
                padding: const EdgeInsets.only(bottom: 27.0),
              ),
              GestureDetector(
                  onTap: disabled
                      ? null
                      : () {
                          onTap(context, IdentityType.IdentityCard);
                        },
                  child: selectDocmentType('images/idCardIcon.png',
                      'KYC.ID_CARD'.tr(), 'KYC.FRONT_AND_BACK'.tr(), disabled))
            ])
          : Container(),
    ]);
  }

  Widget selectDocmentType(
      String imgPath, String title, String desc, bool disabled) {
    final color = disabled ? FXColor.lightGray : FXColor.mainBlueColor;
    return Container(
      width: double.infinity,
      padding: getEdgeInsetsSymmetric(),
      decoration: FXUI.boxDecorationWithShadow,
      child: Row(
        children: [
          Expanded(
              flex: 2,
              child: Container(
                  alignment: Alignment.centerRight,
                  margin: const EdgeInsets.only(right: 12.0),
                  child: Image.asset(imgPath,
                      package: 'euruswallet',
                      width: 30,
                      height: 30,
                      color: color,
                      fit: BoxFit.contain))),
          Expanded(
              flex: 3,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: FXUI.normalTextStyle.copyWith(color: color),
                  ),
                  Text(
                    desc,
                    style: FXUI.normalTextStyle
                        .copyWith(fontSize: 12, color: color),
                  )
                ],
              )),
        ],
      ),
    );
  }
}
