#import "ConfigurePlugin.h"
#if __has_include(<configure/configure-Swift.h>)
#import <configure/configure-Swift.h>
#else
// Support project import fallback if the generated compatibility header
// is not copied when this plugin is created as a library.
// https://forums.swift.org/t/swift-static-libraries-dont-copy-generated-objective-c-header/19816
#import "configure-Swift.h"
#endif

@implementation ConfigurePlugin
+ (void)registerWithRegistrar:(NSObject<FlutterPluginRegistrar>*)registrar {
  [SwiftConfigurePlugin registerWithRegistrar:registrar];
}
@end
