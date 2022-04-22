import 'package:euruswallet/common/commonMethod.dart';

import 'apiResponseModel.dart';
import 'user_storage.dart';

class UserStorageResponseModel implements ApiResponseModel {
  int returnCode;
  String message;
  String? nonce;
  Map<String, dynamic> data;

  int? get userId => data['userId'] as int?;
  int? get platform => data['platform'] as int?;
  int? get sequence => data['sequence'] as int?;
  UserStorage? get storage =>
      !isEmptyString(string: (data['storage'] as String?))
          ? UserStorage.fromEncodedJson(data['storage'])
          : null;

  UserStorageResponseModel(
    this.returnCode,
    this.message,
    this.nonce,
    this.data,
  );

  factory UserStorageResponseModel.fromJson(Map<String, dynamic> json) =>
      UserStorageResponseModel(
        json['returnCode'] as int,
        json['message'] as String,
        json['nonce'] as String?,
        (json['data'] as Map<String, dynamic>?) != null
            ? json['data'] as Map<String, dynamic>
            : {},
      );
}
