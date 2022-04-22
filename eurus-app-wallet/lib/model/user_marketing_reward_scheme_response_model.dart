class UserMarketingRewardSchemeResponseModel {
  final int returnCode;
  final String message;
  final String? nonce;
  final List<Map<String, dynamic>> data;

  List<UserMarketingRewardScheme> get schemes =>
      data.map((e) => UserMarketingRewardScheme.fromJson(e)).toList();

  UserMarketingRewardSchemeResponseModel(
      this.returnCode, this.message, this.nonce, this.data);

  factory UserMarketingRewardSchemeResponseModel.fromJson(
          Map<String, dynamic> json) =>
      UserMarketingRewardSchemeResponseModel(
        json['returnCode'] as int,
        json['message'] as String,
        json['nonce'] as String?,
        json['data'] is List<dynamic>
            ? List.from((json['data'] as List<dynamic>)
                .where((element) => element is Map<String, dynamic>))
            : [],
      );
}

class UserMarketingRewardScheme {
  final int? schemeId;
  final UserMarketingRewardSchemeContent? en;
  final UserMarketingRewardSchemeContent? zhTw;
  final UserMarketingRewardSchemeContent? zhCn;

  UserMarketingRewardScheme(this.schemeId, this.en, this.zhTw, this.zhCn);

  factory UserMarketingRewardScheme.fromJson(Map<String, dynamic> json) =>
      UserMarketingRewardScheme(
        json['schemeId'] as int?,
        (json['en'] as Map<String, dynamic>?) != null
            ? UserMarketingRewardSchemeContent.fromJson(
                json['en'] as Map<String, dynamic>)
            : null,
        (json['zh-tw'] as Map<String, dynamic>?) != null
            ? UserMarketingRewardSchemeContent.fromJson(
                json['zh-tw'] as Map<String, dynamic>)
            : null,
        (json['zh-cn'] as Map<String, dynamic>?) != null
            ? UserMarketingRewardSchemeContent.fromJson(
                json['zh-cn'] as Map<String, dynamic>)
            : null,
      );
}

class UserMarketingRewardSchemeContent {
  final String? title;
  final String? details;

  UserMarketingRewardSchemeContent(this.title, this.details);

  factory UserMarketingRewardSchemeContent.fromJson(
          Map<String, dynamic> json) =>
      UserMarketingRewardSchemeContent(
        json['title'] as String?,
        json['details'] as String?,
      );
}
