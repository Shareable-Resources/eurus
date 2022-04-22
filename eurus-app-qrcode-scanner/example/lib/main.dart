import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:app_qrcode_scanner/app_qrcode_scanner.dart';
import 'package:permission_handler/permission_handler.dart';

void main() => runApp(MaterialApp(home: QRViewExample()));

class QRViewExample extends StatefulWidget {
  const QRViewExample({
    Key? key,
  }) : super(key: key);

  @override
  State<StatefulWidget> createState() => _QRViewExampleState();
}

class _QRViewExampleState extends State<QRViewExample> {
  final GlobalKey qrKey = GlobalKey(debugLabel: 'QR');

  String result = '';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            TextButton(
              onPressed: () async {
                await scanQRCode(context);
              },
              child: Padding(
                padding: EdgeInsets.all(20),
                child: Text('Show QRCode Scanner'),
              ),
            ),
            Text(
              result,
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Future scanQRCode(BuildContext _) async {
    var response = await QRCodeScanner().tryOpenScanner(_);

    setState(() {
      result = response;
    });
  }
}

class QRCodeScanner extends AppQRCodeScanner {
  @override
  Future<bool?> ckCameraPermission() async {
    var camStatus = await Permission.camera.status;
    bool? camPerm;
    if (camStatus.isDenied) {
      camPerm = null;
    } else {
      if (camStatus.isPermanentlyDenied) {
        camPerm = false;
      } else {
        camPerm = camStatus.isGranted;
      }
    }

    return camPerm;
  }

  @override
  Future<bool?> ckPhotoPermission() async {
    var photoStatus = await Permission.photos.status;
    bool? photoPerm;
    if (photoStatus.isDenied) {
      photoPerm = null;
    } else {
      if (photoStatus.isPermanentlyDenied) {
        photoPerm = false;
      } else {
        photoPerm = photoStatus.isGranted;
      }
    }

    return photoPerm;
  }
}
