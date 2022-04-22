import 'package:flutter/material.dart';

import '../template/prem_template.dart';

class CameraPermModal extends PermModalTemplate {
  CameraPermModal({
    this.disabled,
    this.openPhotoAction,
    this.themeColor,
  }) : super(
          title: 'Camera${disabled == true ? ' is Disabled' : ''}',
          desc: disabled == true
              ? 'Please allow us to access your Camera in your System Setting'
              : 'Get Started by allowing us to use your camera',
          color: themeColor,
          icon: disabled == true
              ? Icons.warning_rounded
              : Icons.camera_alt_rounded,
          iconColor: disabled == true ? Color(0xffFB4245) : themeColor,
          hideDecline: disabled ?? false,
          acceptText: disabled == true ? 'Ok' : null,
          otherAction: openPhotoAction,
        );

  final bool? disabled;
  final Widget? openPhotoAction;
  final Color? themeColor;
}
