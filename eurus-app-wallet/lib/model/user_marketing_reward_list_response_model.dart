class UserMarketingRewardListResponseModel {
  final int returnCode;
  final String message;
  final String? nonce;
  final List<Map<String, dynamic>> data;

  List<UserMarkingRewardList> get list =>
      data.map((e) => UserMarkingRewardList.fromJson(e)).toList();

  UserMarketingRewardListResponseModel(
    this.returnCode,
    this.message,
    this.nonce,
    this.data,
  );

  factory UserMarketingRewardListResponseModel.fromJson(
          Map<String, dynamic> json) =>
      UserMarketingRewardListResponseModel(
        json['returnCode'] as int,
        json['message'] as String,
        json['nonce'] as String?,
        json['data'] is List<dynamic>
            ? List.from((json['data'] as List<dynamic>)
                .where((element) => element is Map<String, dynamic>))
            : [],
      );
}

class UserMarkingRewardList {
  final int? rewardType;
  final String? assetName;
  final int? amount;
  final String? txHash;
  final DateTime? createdDate;

  UserMarkingRewardList(
    this.rewardType,
    this.assetName,
    this.amount,
    this.txHash,
    this.createdDate,
  );

  factory UserMarkingRewardList.fromJson(Map<String, dynamic> json) =>
      UserMarkingRewardList(
        json['rewardType'] as int?,
        json['assetName'] as String?,
        json['amount'] as int?,
        json['txHash'] as String?,
        json['createdDate'] is String
            ? DateTime.tryParse((json['createdDate'] as String))
            : null,
      );
}
