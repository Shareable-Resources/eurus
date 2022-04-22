import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:flutter_neumorphic/flutter_neumorphic.dart';

class FXColor {
  static Color centralizedGrayTextColor = Color(0xff7A869A).withOpacity(0.5);
  static Color lightGreyTextColor = Color(0xFFf4f2f7).withOpacity(0.75);
  static const Color verylightBlack = Color(0xfff8f8fb);
  static const Color lightBlack = Color(0xff686C71);
  static const Color gray = Color(0x993C3C43);
  static const Color textGray = Color(0xff7A869A);
  static const Color opacityGray = Color(0x953C3C43);
  static const Color lightGray = Color(0xffC1C7D0);
  static const Color veryLightGray = Color(0xffF2F4F7);
  static const Color cancelGrayButton = Color(0xFF949494);
  static const Color orange = Color(0xffFAA61A);
  static const Color middleBlack = Color(0xff4A4A4A);
  static const Color darkBlack = Color(0xFF596D79);
  static const Color untabColor = Color(0xffCCCCCC);
  static const Color grayLine = Color(0xcc4E5053);
  static const Color inputBorder = Color(0x1A4E5053);
  static const Color mainBlueColor = Color(0xff059FDD);
  static const Color mainDeepBlueColor = Color(0xff390ddd);
  static const Color grayTextColor = Color(0xff687671);
  static const Color blackColor = Color(0xFF002251);
  static const Color alertRedColor = Color(0xffd32f2f);
  static const Color warningRedColor = Color(0xFFC75353);
  static const Color deepBlue = Color(0xFF4406dd);
  static const Color mediumGrayColor = Color(0xff959595);
  static const Color purpleColor = Color(0xff7B61FF);
  static const Color lightWhiteColor = Color(0xfff7f9fb);
  static const Color deepGrayColor = Color(0xff757575);
  static const Color lightBlackColor = Color(0xff4e5053);
  static const Color lightRedColor = Color(0xffFB4245);
  static const Color listItemBlueColor = Color(0xFFDBE8FF);
  static const Color segmentButtonBackgroundColor = Color(0xFFF4F6F7);
  static const Color greyTextColor = Color(0xFFD1D1D1);
  static const Color veryLightGreyTextColor = Color(0xFFf7f9fb);
  static const Color grey80Color = Color(0xffcccccc);
  static const Color lightGreyColor2 = Color(0xffEDEEF0);
  static const Color blueGreenColor = Color(0xff2F3C4E);
  static const Color assetsListPurpleColor = Color(0xff5856d6);
  static const Color assetsListBlueColor = Color(0xff5ac8fa);
  static const Color greenColor = Color(0xff47d187);
  static const Color pinkRedColor = Color(0xffff5e5e);
  static const Color lightPurpleColor = Color(0xFFAA94FC);
  static const Color deepGreyColor = Color(0xff666666);
  static const Color DoveGreyColor = Color(0xff999999);
  static const Color transferSuccessGreyColor = Color(0xffE9EBEF);
  static const Color dividerColor = Color(0xFFF2F5F7);
  static const Color placeholderGreyColor = Color(0xFFB6BDC9);
  static const Color grey44 = Color(0xFF707070);
  static const Color lightGrey = Color(0xFFA4AAB1);
  static const Color dimGrey = Color(0xFF616161);
}

class FXUI {
  static double circular = 15;
  static BorderRadius cricleRadius = BorderRadius.circular(FXUI.circular);
  static TextStyle normalTextStyle = TextStyle(
      color: Colors.black,
      fontSize: 15,
      fontFamily: 'packages/euruswallet/SFProDisplay');
  static TextStyle titleTextStyle = TextStyle(
      color: Colors.black,
      fontSize: 20,
      fontWeight: FontWeight.bold,
      fontFamily: 'packages/euruswallet/SFProDisplay');
  static TextStyle subtitleTextStyle = TextStyle(
      color: FXColor.textGray,
      fontSize: 12,
      fontWeight: FontWeight.w500,
      fontFamily: 'packages/euruswallet/SFProDisplay');
  static InputBorder inputBorder = new OutlineInputBorder(
      borderSide: BorderSide(width: 1, color: Colors.transparent),
      borderRadius: cricleRadius);
  static TextStyle inputStyle = TextStyle(
      color: FXColor.middleBlack,
      fontSize: 16,
      fontFamily: 'packages/euruswallet/SFProDisplay');
  static TextStyle hintStyle = TextStyle(
      color: FXColor.untabColor,
      fontSize: 14,
      fontWeight: FontWeight.w500,
      fontFamily: 'packages/euruswallet/SFProDisplay');
  static InputDecoration inputDecoration = InputDecoration(
    contentPadding: EdgeInsets.symmetric(vertical: 15.5, horizontal: 15.5),
    focusedBorder: FXUI.inputBorder,
    disabledBorder: FXUI.inputBorder,
    enabledBorder: FXUI.inputBorder,
    border: FXUI.inputBorder,
    hintStyle: FXUI.hintStyle,
    hintMaxLines: 4,
    errorMaxLines: 4,
    fillColor: FXColor.veryLightGray,
    filled: true,
    counterText: "", //控制最大字符是否显示
  );
  static InputDecoration defaultTextFieldInputDecoration = InputDecoration(
    filled: true,
    fillColor: FXColor.lightGreyTextColor,
    border: OutlineInputBorder(
        borderSide: BorderSide.none, borderRadius: cricleRadius),
    contentPadding: EdgeInsets.all(16),
  );
  static BoxDecoration circleBoxDecoration = BoxDecoration(
    color: Colors.white,
    borderRadius: cricleRadius,
  );
  static BoxDecoration boxDecorationWithShadow = BoxDecoration(
    color: Colors.white,
    borderRadius: BorderRadius.circular(10),
    boxShadow: [
      BoxShadow(
        color: FXColor.grey80Color,
        offset: Offset(1, 2),
        blurRadius: 8,
      )
    ],
  );
  static DecoratedBox divider = DecoratedBox(
    decoration: BoxDecoration(color: FXColor.dividerColor),
    child: SizedBox(
      width: double.infinity,
      height: 1,
    ),
  );
  static NeumorphicStyle neumorphicBannerImage = NeumorphicStyle(
    lightSource: LightSource.top,
    intensity: 0.9,
    color: Colors.transparent,
  );
  static Neumorphic neumorphicTextField(
    BuildContext context, {
    EdgeInsets padding = const EdgeInsets.all(0),
    NeumorphicBoxShape? shape,
    TextInputAction? textInputAction,
    TextInputType? keyboardType,
    String? hintText,
    String? errorText,
    TextStyle? errorStyle,
    Widget? prefixIcon,
    Widget? suffixIcon,
    TextEditingController? controller,
    FocusNode? focusNode,
    bool obscureText = false,
    bool autocorrect = true,
    void Function(String)? onSubmitted,
  }) =>
      Neumorphic(
        padding: padding,
        style: NeumorphicStyle(
          color: Color(0xfff8f8f8),
          boxShape: shape ?? NeumorphicBoxShape.stadium(),
          depth: NeumorphicTheme.embossDepth(context),
          intensity: 0.65,
          lightSource: LightSource.top,
        ),
        child: Center(
          child: TextField(
            keyboardType: keyboardType,
            textInputAction: textInputAction,
            style: FXUI.inputStyle,
            decoration: InputDecoration.collapsed(
              hintText: hintText,
              hintStyle: FXUI.hintStyle.copyWith(
                color: FXColor.placeholderGreyColor,
                fontSize: 13,
              ),
            ).copyWith(
              errorText: errorText,
              errorStyle: errorStyle,
              prefixIcon: prefixIcon,
              prefixIconConstraints: prefixIcon != null
                  ? BoxConstraints(
                      maxHeight: 100,
                      maxWidth: 100,
                    )
                  : null,
              suffixIcon: suffixIcon,
              suffixIconConstraints: suffixIcon != null
                  ? BoxConstraints(
                      maxHeight: 100,
                      maxWidth: 100,
                    )
                  : null,
            ),
            obscureText: obscureText,
            autocorrect: autocorrect,
            maxLines: 1,
            controller: controller,
            focusNode: focusNode,
            onSubmitted: onSubmitted,
          ),
        ),
      );
}

EdgeInsets getEdgeInsetsSymmetric({double vertical: 16, double horizontal: 0}) {
  return EdgeInsets.symmetric(vertical: 16, horizontal: horizontal);
}
