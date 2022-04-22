#import "FlutterQrPlugin.h"
#import "app_qrcode_scanner/app_qrcode_scanner-Swift.h"
@implementation FlutterQrPlugin
+ (void)registerWithRegistrar:(NSObject<FlutterPluginRegistrar>*)registrar {
    [SwiftFlutterQrPlugin registerWithRegistrar:registrar];
}
@end
