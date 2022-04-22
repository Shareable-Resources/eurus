import 'package:flutter/material.dart';
import 'package:app_qrcode_scanner/app_qrcode_scanner.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';

class PhotoLibraryPermPopup extends PermModalTemplate {
  PhotoLibraryPermPopup({
    this.disabled,
    this.themeColor,
  }) : super(
          title: disabled == true
              ? 'SCANNER_PAGE.PHOTO_DISABLED_POPUP.TITLE'.tr()
              : 'SCANNER_PAGE.PHOTO_PERM_POPUP.TITLE'.tr(),
          desc: disabled == true
              ? 'SCANNER_PAGE.PHOTO_DISABLED_POPUP.DESC'.tr()
              : 'SCANNER_PAGE.PHOTO_PERM_POPUP.DESC'.tr(),
          color: themeColor,
          icon: disabled == true
              ? Icons.warning_rounded
              : Icons.image_rounded,
          iconColor: disabled == true ? FXColor.lightRedColor : themeColor,
          hideDecline: disabled ?? false,
          declineText: 'COMMON.CANCEL'.tr(),
          acceptText:
              disabled == true ? 'COMMON.OK'.tr() : 'COMMON.CONTINUE'.tr(),
        );

  final bool? disabled;
  final Color? themeColor;
}
