import 'package:flutter/material.dart';

import '../template/prem_template.dart';

class PhotoLibraryPermModal extends PermModalTemplate {
  PhotoLibraryPermModal({
    this.disabled,
    this.themeColor,
  }) : super(
          title: 'Photo Library${disabled == true ? ' is Disabled' : ''}',
          desc: disabled == true
              ? 'Please allow us to access your Photo Library in your System Setting'
              : 'Get Started by allowing us to access your Photo Library',
          color: themeColor,
          icon: disabled == true ? Icons.warning_rounded : Icons.image_rounded,
          iconColor: disabled == true ? Color(0xffFB4245) : themeColor,
          hideDecline: disabled ?? false,
          acceptText: disabled == true ? 'Ok' : null,
        );

  final bool? disabled;
  final Color? themeColor;
}
