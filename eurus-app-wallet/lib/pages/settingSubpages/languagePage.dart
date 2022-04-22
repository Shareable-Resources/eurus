import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/pages/settingSubpages/cardContainer.dart';
import 'package:euruswallet/pages/settingSubpages/settingAppBar.dart';

class LanguagePage extends StatefulWidget {
  LanguagePage({Key? key}) : super(key: key);

  _LanguagePageState createState() => _LanguagePageState();
}

class _LanguagePageState extends State<LanguagePage> {
  String? _activeLocale;
  Color get themeColor => common.getBackGroundColor();

  @override
  void initState() {
    _getAppLang();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: FXColor.veryLightGreyTextColor,
      appBar: SettingAppBar(true),
      body: Container(
        child: SingleChildScrollView(
          child: SafeArea(
            child: CardContainer(
              'LANGUAGE_PAGE.TITLE'.tr(),
              Container(
                  child: Column(
                children: [
                  _langOptBtn('English', 'en_US', borderBtm: true),
                  _langOptBtn('简体中文', 'zh_Hans'),
                  _langOptBtn('繁體中文', 'zh_Hant'),
                ],
              )),
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _getAppLang() async {
    String? appLang = await NormalStorageKit().readValue('APP_LANGUAGE') ??
        context.locale.toString();
    setState(() {
      _activeLocale = appLang;
    });
  }

  Widget _langOptBtn(
    String title,
    String? tag, {
    bool borderBtm = false,
  }) {
    print('$tag :: $_activeLocale');
    bool isActive =
        tag != null && _activeLocale != null && tag == _activeLocale;

    return InkWell(
      onTap: () async {
        Locale langToSet = tag == 'en_US'
            ? Locale('en', 'US')
            : tag == 'zh_Hant'
                ? Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hant')
                : Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hans');
        await NormalStorageKit().setValue(langToSet.toString(), 'APP_LANGUAGE');
        context.setLocale(langToSet);
        api.postUserStorage(
          currentLanguage: tag == 'zh_Hant'
              ? 'tw'
              : tag == 'zh_Hans'
                  ? 'cn'
                  : 'en',
        );
        setState(() {
          _activeLocale = langToSet.toString();
        });
      },
      child: Container(
        width: double.infinity,
        decoration: BoxDecoration(
          border: Border(
            bottom: borderBtm == true
                ? BorderSide(width: 1, color: FXColor.verylightBlack)
                : BorderSide.none,
          ),
        ),
        padding: EdgeInsets.symmetric(vertical: 25, horizontal: 20),
        child: Text(
          title,
          style: FXUI.normalTextStyle.copyWith(
            color: isActive ? themeColor : Colors.black,
            fontWeight: isActive ? FontWeight.bold : FontWeight.normal,
          ),
        ),
      ),
    );
  }
}
