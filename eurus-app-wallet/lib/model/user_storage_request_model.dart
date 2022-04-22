import 'user_storage.dart';

class UserStorageRequestModel {
  final String? nonce;
  final int? platform;
  final UserStorage? storage;

  UserStorageRequestModel({
    this.nonce,
    this.platform = 0,
    this.storage,
  });

  factory UserStorageRequestModel.fromJson(Map<String, dynamic> json) =>
      UserStorageRequestModel(
        nonce: json['nonce'] as String?,
        platform: json['platform'] as int?,
        storage: (json['storage'] as String?) != null
            ? UserStorage.fromEncodedJson(json['storage'] as String)
            : null,
      );

  Map<String, dynamic> toJson() => {
        'nonce': nonce,
        'platform': platform,
        'storage': storage?.toEncodedJson(),
      };
}
