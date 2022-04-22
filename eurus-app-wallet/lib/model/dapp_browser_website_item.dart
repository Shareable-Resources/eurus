enum DappBrowserWebsiteItemType { favorite, searchHistory }

extension DappBrowserWebsiteItemTypeExtension on DappBrowserWebsiteItemType {
  String get storageKey {
    switch (this) {
      case DappBrowserWebsiteItemType.favorite:
        return 'dapp_browser_website_favorite_list';
      case DappBrowserWebsiteItemType.searchHistory:
        return 'dapp_browser_search_history_list';
    }
  }
}

class DappBrowserWebsiteItem {
  final String? title;
  final String url;

  DappBrowserWebsiteItem(this.url, {this.title});

  DappBrowserWebsiteItem.fromJson(Map<String, dynamic> json)
      : title = json['title'],
        url = json['url'];

  Map<String, dynamic> toJson() => {
        'title': title,
        'url': url,
      };

  @override
  bool operator ==(other) {
    return other is DappBrowserWebsiteItem && other.url == url;
  }

  @override
  int get hashCode => super.hashCode;
}
