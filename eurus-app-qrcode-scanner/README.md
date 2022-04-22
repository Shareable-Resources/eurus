# app_qrcode_scanner

app_qrcode_scanner is a plugin for application to scan QRCode using camera or images from their device photo library. 

This plugin is forked from [qr_code_scanner](https://github.com/juliuscanute/qr_code_scanner)

## Installation
### iOS
Camera and Photo library usage description area required.
Add the following lines into your Info.plist file to do so.
```plist
<key>io.flutter.embedded_views_preview</key>
<true/>
<key>NSCameraUsageDescription</key>
<string>This app needs camera access to scan QR codes</string>
<key>NSPhotoLibraryUsageDescription</key>
<string>This app needs photo library access to scan QR codes from image</string>
```
### Android
For API 29+
Add the following lines into AndroidManifest.xml, which is required in Android Q since this attribute default value is false
```xml
android:requestLegacyExternalStorage="true"
```
## Usage
Call tryOpenScanner with context will prompt a camera modal for QRCode scanning 
```dart
import 'package:app_qrcode_scanner/app_qrcode_scanner.dart';

Widget build(BuildContext context) {
  return Scaffold(
    body: Column(
      children: [
        onPressed: () async {
          /// await QRCode content as [String]
          String result = await tryOpenScanner(context);
        }
      ]
    )
  );
}
```

### Permission
You may want to check the camera and photo library usage permission before launching the QRCode scanner
You can customize the soft launch permission modal by create custom class extends [CustomModal]
```dart
import 'package:app_qrcode_scanner/app_qrcode_scanner.dart';

/// Get permission using another plugin: permission_handler
var camStatus = await Permission.camera.status;
var photoStatus = await Permission.photos.status;

var camPerm = camStatus.isUndetermined ? null : camStatus.isGranted;
var photoPerm = photoStatus.isUndetermined ? null : photoStatus.isGranted;

/// await QRCode content as [String]
String result = await tryOpenScanner(
  context, 
  hvCameraPerm: camPerm, 
  cameraPermModal: CameraPermissionModal(), 
  hvPhotoPerm: photoPerm,
);

class CameraPermissionModal extends CustomModal {
  @override
  Widget buildPage(
    BuildContext context,
    Animation<double> animation,
    Animation<double> secondaryAnimation,
  ) {
    return Material(
      type: MaterialType.transparency,
      child: Container() // Your own design
    )
  }
}
```