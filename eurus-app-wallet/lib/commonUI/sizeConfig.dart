import 'package:flutter/material.dart';

class SizeConfig {
  late MediaQueryData mediaQueryData;
  late double screenWidth;
  late double screenHeight;
  late double blockSizeHorizontal;
  late double blockSizeVertical;

  late double safeAreaHorizontal;
  late double safeAreaVertical;
  late double safeBlockHorizontal;
  late double safeBlockVertical;

  //reponse width rate from device width /iPhone 11 Prox Max width
  //reponse height rate from device height /iPhone 11 Prox Max width
  late double responseWidthRate;
  late double responseHeightRate;
  late double responseFrontSizeRate;
  late double leftPadding;
  late double heightWithoutAppBar;

  static final SizeConfig _instance = SizeConfig._internal();

  SizeConfig._internal() {
    // init things inside this
  }

  factory SizeConfig({BuildContext? context}) {
    if (context == null) {
      return _instance;
    }

    _instance.mediaQueryData = MediaQuery.of(context);
    MediaQueryData mediaQueryData = _instance.mediaQueryData;
    _instance.screenWidth = mediaQueryData.size.width;
    _instance.screenHeight = mediaQueryData.size.height;
    _instance.blockSizeHorizontal = _instance.screenWidth / 100;
    _instance.blockSizeVertical = _instance.screenHeight / 100;

    _instance.safeAreaHorizontal =
        mediaQueryData.padding.left + mediaQueryData.padding.right;
    _instance.safeAreaVertical =
        mediaQueryData.padding.top + mediaQueryData.padding.bottom;
    _instance.safeBlockHorizontal =
        (_instance.screenWidth - _instance.safeAreaHorizontal) / 100;
    _instance.safeBlockVertical =
        (_instance.screenHeight - _instance.safeAreaVertical) / 100;

    //https://www.paintcodeapp.com/news/ultimate-guide-to-iphone-resolutions
//base on iphone 11 Pro Max
    _instance.responseWidthRate =
        ((_instance.screenWidth - _instance.safeAreaHorizontal) / 414);
    _instance.responseHeightRate =
        ((_instance.screenHeight - _instance.safeAreaVertical) / 896);
    _instance.responseFrontSizeRate = _instance.responseWidthRate;
    _instance.leftPadding = 25;
    _instance.heightWithoutAppBar = _instance.screenHeight -
        kToolbarHeight -
        MediaQuery.of(context).padding.top;

//    print("screenWidth: ${_instance.screenWidth}");
//    print("screenHeight: ${_instance.screenHeight}");
//    print("blockSizeHorizontal: ${_instance.blockSizeHorizontal}");
//    print("blockSizeVertical: ${_instance.blockSizeVertical}");
//    print("_safeAreaHorizontal: ${_instance.safeAreaHorizontal}");
//    print("_safeAreaVertical: ${_instance.safeAreaVertical}");
//    print("safeBlockHorizontal: ${_instance.safeBlockHorizontal}");
//    print("safeBlockVertical: ${_instance.safeBlockVertical}");
//    print("responseWidthRate: ${_instance.responseWidthRate}");
//    print("responseHeightRate: ${_instance.responseHeightRate}");
//    print("responseFrontSizeRate: ${_instance.responseFrontSizeRate}");
//    print("fivePercentPadding: ${_instance.fivePercentPadding}");
//    print("insidefivePercentPadding: ${_instance.insidefivePercentPadding}");
    return _instance;
  }
}

SizeConfig size = SizeConfig();
