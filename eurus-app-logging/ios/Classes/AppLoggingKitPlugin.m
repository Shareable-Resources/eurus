#import "AppLoggingKitPlugin.h"
#if __has_include(<app_logging_kit/app_logging_kit-Swift.h>)
#import <app_logging_kit/app_logging_kit-Swift.h>
#else
// Support project import fallback if the generated compatibility header
// is not copied when this plugin is created as a library.
// https://forums.swift.org/t/swift-static-libraries-dont-copy-generated-objective-c-header/19816
#import "app_logging_kit-Swift.h"
#endif

@implementation AppLoggingKitPlugin
+ (void)registerWithRegistrar:(NSObject<FlutterPluginRegistrar>*)registrar {
  [SwiftAppLoggingKitPlugin registerWithRegistrar:registrar];
}
@end
