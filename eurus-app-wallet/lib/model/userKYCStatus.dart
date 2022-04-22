import 'package:euruswallet/common/commonMethod.dart';

class UserKYCStatus {
  int returnCode;
  String message;
  String nonce;
  UserKYCStatusData? data;

  UserKYCStatus(this.returnCode, this.message, this.nonce, this.data);

  UserKYCStatus.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data'] != null
            ? new UserKYCStatusData.fromJson(json['data']!)
            : null;

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data != null) 'data': this.data!.toJson(),
      };
}

class UserKYCStatusData {
  List<KYCStatusData> kycdata;
  int kycLevel;

  UserKYCStatusData(this.kycdata, this.kycLevel);

  UserKYCStatusData.fromJson(Map<String, dynamic> json)
      : kycdata = json['data'] != null
            ? json['data']
                .map<KYCStatusData>((v) => new KYCStatusData.fromJson(v))
                .toList()
            : [],
        kycLevel = json['kycLevel'];

  Map<String, dynamic> toJson() => {
        if (this.kycdata.isNotEmpty)
          'kycdata': this.kycdata.map((v) => v.toJson()).toList(),
        'kycLevel': this.kycLevel,
      };
}

class KYCStatusData {
  int id;
  int userId;
  int kycLevel;
  String? approvalDate;
  String? operatorId;
  int kycRetryCount;
  String kycCountryCode;
  String createdDate;
  String lastModifiedDate;
  int kycStatus;
  int kycDoc;
  List<KYCImage> images;

  KYCStatusData(
      this.id,
      this.userId,
      this.kycLevel,
      this.approvalDate,
      this.operatorId,
      this.kycRetryCount,
      this.kycCountryCode,
      this.createdDate,
      this.lastModifiedDate,
      this.kycStatus,
      this.kycDoc,
      this.images);

  KYCStatusData.fromJson(Map<String, dynamic> json)
      : id = json['id'],
        userId = json['userId'],
        kycLevel = json['kycLevel'],
        approvalDate = json['approvalDate'],
        operatorId = json['operatorId'],
        kycRetryCount = json['kycRetryCount'],
        kycCountryCode = json['kycCountryCode'],
        createdDate = json['createdDate'],
        lastModifiedDate = json['lastModifiedDate'],
        kycStatus = json['kycStatus'],
        kycDoc = json['kycDoc'],
        images = json['images'] != null
            ? json['images']
                .map<KYCImage>((v) => new KYCImage.fromJson(v))
                .toList()
            : [];

  Map<String, dynamic> toJson() => {
        'id': this.id,
        'userId': this.userId,
        'kycLevel': this.kycLevel,
        'approvalDate': this.approvalDate,
        'operatorId': this.operatorId,
        'kycRetryCount': this.kycRetryCount,
        'kycCountryCode': this.kycCountryCode,
        'createdDate': this.createdDate,
        'lastModifiedDate': this.lastModifiedDate,
        'kycStatus': this.kycStatus,
        'kycDoc': this.kycDoc,
        if (this.images.isNotEmpty)
          'images': this.images.map((v) => v.toJson()).toList(),
      };
}

class KYCImage {
  int userKycId;
  int docType;
  int imageSeq;
  int status;
  String imagePath;
  String createdDate;
  String lastModifiedDate;
  RejectReasonData? rejectReason;
  String? operatorId;

  KYCImage(
      this.userKycId,
      this.docType,
      this.imageSeq,
      this.status,
      this.imagePath,
      this.createdDate,
      this.lastModifiedDate,
      this.rejectReason,
      this.operatorId);

  KYCImage.fromJson(Map<String, dynamic> json)
      : userKycId = json['userKYCId'],
        docType = json['docType'],
        imageSeq = json['imageSeq'],
        status = json['status'],
        imagePath = json['imagePath'],
        createdDate = json['createdDate'],
        lastModifiedDate = json['lastModifiedDate'],
        rejectReason = isEmptyString(string: json['rejectReason']) == false
            ? new RejectReasonData.fromJson(jsonDecode(json['rejectReason']!))
            : null,
        operatorId = json['operatorId'];

  Map<String, dynamic> toJson() => {
        'userKycId': this.userKycId,
        'docType': this.docType,
        'imageSeq': this.imageSeq,
        'status': this.status,
        'imagePath': this.imagePath,
        'createdDate': this.createdDate,
        'lastModifiedDate': this.lastModifiedDate,
        if (this.rejectReason != null)
          'rejectReason': this.rejectReason!.toJson(),
        'operatorId': this.operatorId,
      };
}

class RejectReasonData {
  String key;
  String en;
  String zhCN;
  String zhHK;

  RejectReasonData(this.key, this.en, this.zhCN, this.zhHK);

  RejectReasonData.fromJson(Map<String, dynamic> json)
      : key = json['key'],
        en = json['en'],
        zhCN = json.containsKey('zh_cn') ? json['zh_cn'] : json['zh'],
        zhHK = json.containsKey('zh_hk') ? json['zh_hk'] : json['zh'];

  Map<String, dynamic> toJson() => {
        'key': this.key,
        'en': this.en,
        'zh_cn': this.zhCN,
        'zh_hk': this.zhHK,
      };
}
