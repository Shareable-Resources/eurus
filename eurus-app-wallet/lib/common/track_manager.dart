import 'package:euruswallet/common/commonMethod.dart';
import 'package:firebase_analytics/firebase_analytics.dart';

class TrackManager {
  final FirebaseAnalytics? _analytics =
      envType == EnvType.Production ? FirebaseAnalytics.instance : null;

  static final TrackManager instance = TrackManager._internal();

  factory TrackManager() {
    return instance;
  }

  TrackManager._internal();

  Future trackSignUp({required CurrentUserType userType}) async {
    _analytics?.logSignUp(
        signUpMethod: userType == CurrentUserType.centralized
            ? 'custodial'
            : 'decentralized');
  }

  Future trackLogin({required CurrentUserType userType}) async {
    _analytics?.logLogin(
        loginMethod: userType == CurrentUserType.centralized
            ? 'custodial'
            : 'decentralized');
  }
}
