#import "EuruswalletPlugin.h"
#if __has_include(<euruswallet/euruswallet-Swift.h>)
#import <euruswallet/euruswallet-Swift.h>
#else
// Support project import fallback if the generated compatibility header
// is not copied when this plugin is created as a library.
// https://forums.swift.org/t/swift-static-libraries-dont-copy-generated-objective-c-header/19816
#import "euruswallet-Swift.h"
#endif

@implementation EuruswalletPlugin
+ (void)registerWithRegistrar:(NSObject<FlutterPluginRegistrar>*)registrar {
  [SwiftEuruswalletPlugin registerWithRegistrar:registrar];
}
@end
