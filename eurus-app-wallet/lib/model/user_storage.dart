import 'dart:convert';

class UserStorage {
  String? currentLanguage;

  UserStorage({
    this.currentLanguage,
  });

  factory UserStorage.fromEncodedJson(String json) => UserStorage(
        currentLanguage: jsonDecode(json)['currentLanguage'] as String?,
      );

  String toEncodedJson() => jsonEncode({
        'currentLanguage': currentLanguage,
      });
}
