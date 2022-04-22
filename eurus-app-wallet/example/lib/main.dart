import 'package:euruswallet/pages/mainApp.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:firebase_core/firebase_core.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await EasyLocalization.ensureInitialized();
  runApp(EasyLocalization(
    supportedLocales: [
      Locale('en', 'US'),
      Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hant'),
      Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hans')
    ],
    path: 'packages/euruswallet/i18n',
    child: MaterialApp(
      home: (MyApp()),
    ),
  ));
}

class MyApp extends StatefulWidget {
  @override
  _MyAppState createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  @override
  Widget build(BuildContext context) {
    NormalStorageKit().readValue('APP_LANGUAGE').then((val) {
      if (val != null) {
        Locale langToSet = val == 'en_US'
            ? Locale('en', 'US')
            : val == 'zh_Hant'
                ? Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hant')
                : Locale.fromSubtags(languageCode: 'zh', scriptCode: 'Hans');
        context.setLocale(langToSet);
      }
    });
    return MainApp();
  }
}
