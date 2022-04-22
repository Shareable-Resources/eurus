import 'package:euruswallet/common/commonMethod.dart';
import 'package:web3dart/crypto.dart';

class UserProfile {
  UserProfile({
    required this.userType,
    this.decenUserType,
    required this.address,
    required this.encryptedAddress,
    required this.encryptedPrivateKey,
    this.alias,
    this.email,
    this.lastLoginTime,
    this.bioAuthValidTime,
    this.bioTxAuthValidTime,
    this.enableBiometric = false,
    this.enableTxBiometric = false,
    this.seedPraseBackuped,
    this.mnemonicSeedPhrases,
    this.pendingTransaction,
  });

  final CurrentUserType userType;
  DecenUserType? decenUserType;
  String address;
  String encryptedAddress;
  String? alias;
  String? email;
  String encryptedPrivateKey;
  String? lastLoginTime;
  String? bioAuthValidTime;
  String? bioTxAuthValidTime;
  String? mnemonicSeedPhrases;
  bool? enableBiometric; //login face id
  bool? enableTxBiometric; //transaction face id
  bool? seedPraseBackuped;
  List<Uint8List>? pendingTransaction = [];

  UserProfile.fromJson(Map<String, dynamic> json)
      : userType = json['userType'] == "cen" ||
                json['userType'] == CurrentUserType.centralized
            ? CurrentUserType.centralized
            : CurrentUserType.decentralized,
        decenUserType = json['decenUserType'] == "created" ||
                json['decenUserType'] == DecenUserType.created
            ? DecenUserType.created
            : DecenUserType.imported,
        address = (json['address'] as String?)?.toLowerCase() ?? '',
        encryptedAddress = json['encryptedAddress'],
        alias = json['alias'],
        email = json['email'],
        encryptedPrivateKey = json['encryptedPrivateKey'],
        lastLoginTime = json['lastLoginTime'],
        bioAuthValidTime = json['bioAuthValidTime'],
        bioTxAuthValidTime = json['bioTxAuthValidTime'],
        enableBiometric = json['enableBiometric'],
        enableTxBiometric = json['enableTxBiometric'],
        seedPraseBackuped = json['seedPraseBackuped'],
        mnemonicSeedPhrases = json['mnemonicSeedPhrases'],
        pendingTransaction = (json['pendingTransaction'] as List?)
            ?.map((e) => hexToBytes(e as String? ?? ''))
            .toList();

  Map<String, dynamic> toJson() {
    return {
      "userType":
          this.userType == CurrentUserType.centralized ? "cen" : "decen",
      "decenUserType":
          this.decenUserType == DecenUserType.created ? "created" : "imported",
      "address": this.address.toLowerCase(),
      "encryptedAddress": this.encryptedAddress,
      "alias": this.alias,
      "email": this.email,
      "encryptedPrivateKey": this.encryptedPrivateKey,
      "lastLoginTime": this.lastLoginTime,
      "bioAuthValidTime": this.bioAuthValidTime,
      "bioTxAuthValidTime": this.bioTxAuthValidTime,
      "enableBiometric": this.enableBiometric,
      "enableTxBiometric": this.enableTxBiometric,
      "seedPraseBackuped": this.seedPraseBackuped,
      "mnemonicSeedPhrases": this.mnemonicSeedPhrases,
      "pendingTransaction":
          this.pendingTransaction?.map((e) => bytesToHex(e)).toList(),
    };
  }
}
