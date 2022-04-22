import 'package:app_qrcode_scanner/app_qrcode_scanner.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/commonUI/cameraPermissionPopup.dart';
import 'package:euruswallet/commonUI/photoPermissionPopup.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:easy_localization/easy_localization.dart';

class QRCodeScanner extends AppQRCodeScanner {
  @override
  Color? get themeColor => _getThemeColor();

  @override
  String? camtPgTitle = 'SCANNER_PAGE.CAMERA_VIEW.TITLE'.tr();
  @override
  String? scanningText = 'SCANNER_PAGE.CAMERA_VIEW.SCANNING'.tr();
  @override
  String? flashOffText = 'SCANNER_PAGE.CAMERA_VIEW.LIGHTS_ON'.tr();
  @override
  String? flashOnText = 'SCANNER_PAGE.CAMERA_VIEW.LIGHTS_OFF'.tr();

  @override
  Future<bool> ckCameraPermission() async {
    final cameraStatus = await Permission.camera.status;
    return  cameraStatus.isGranted || cameraStatus.isLimited || cameraStatus.isDenied;
  }

  @override
  Future<bool> ckPhotoPermission() async {
    final photoStatus = await Permission.photos.status;
    return photoStatus.isGranted || photoStatus.isLimited || photoStatus.isDenied;
  }

  @override
  Widget? imgPickerBtn(
    BuildContext _,
    bool? hvPhotoPerm,
    CustomModal photoPermModal, {
    Color? themeColor,
  }) {
    return hvPhotoPerm != false
        ? Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Text('-- ${'COMMON.OR'.tr()} --'),
              TextButton(
                onPressed: () async {
                  var result = await tryOpenImgPicker(
                    _,
                    hvPhotoPerm,
                    photoPermModal,
                    themeColor: themeColor,
                  );
                  Navigator.of(_).pop(result);
                },
                child:
                    Text('SCANNER_PAGE.CAMERA_PERM_POPUP.SCAN_FROM_IMG'.tr()),
              )
            ],
          )
        : null;
  }

  @override
  CustomModal genCameraPermModal(
    bool? disabled,
    Widget? openPhotoAction, {
    Color? themeColor,
  }) {
    return CameraPermissionPopup(
      disabled: disabled,
      openPhotoAction: openPhotoAction,
      themeColor: themeColor,
    );
  }

  @override
  CustomModal genPhotoPermModal(bool? disabled, {Color? themeColor}) {
    return PhotoLibraryPermPopup(disabled: disabled, themeColor: themeColor);
  }

  Color _getThemeColor() {
    return common.getBackGroundColor();
  }
}
