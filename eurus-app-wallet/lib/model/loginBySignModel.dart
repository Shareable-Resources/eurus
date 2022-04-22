import 'package:app_security_kit/decryption_helper.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/model/apiResponseModel.dart';

class LoginBySignModel implements ApiResponseModel {
  int returnCode;
  String message;
  Map<String, dynamic> data;

  bool get callSuccess => returnCode == 0;
  String? get mainnetWalletAddress =>
      callSuccess ? data['mainnetWalletAddress'] : null;
  String? get token => callSuccess ? data['token'] : null;
  String? get walletAddress => callSuccess ? data['walletAddress'] : null;
  int? get expiryTime => callSuccess ? data['expiryTime'] : null;
  String? get ownerWalletAddress =>
      callSuccess ? data['ownerWalletAddress'] : null;
  String? get encryptedMnemonic => callSuccess ? data['mnemonic'] : null;
  String? get decryptedMnemonic {
    try {
      return (!isEmptyString(string: encryptedMnemonic) &&
              common.rsaPrivateKey != null)
          ? DecryptionHelper(privateKey: common.rsaPrivateKey!)
              .decryptRAEncryption(encryptedMnemonic!)
          : null;
    } catch (e) {
      return null;
    }
  }

  bool? get isMetaMaskUser => callSuccess ? data['isMetaMaskUser'] : null;

  LoginBySignModel(this.returnCode, this.message, this.data);

  LoginBySignModel.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        data = json['data'] != null ? json['data'] : {};
}
