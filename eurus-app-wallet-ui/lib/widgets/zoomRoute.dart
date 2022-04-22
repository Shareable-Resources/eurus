import 'package:flutter/material.dart';

class MyPageRoute<T> extends MaterialPageRoute<T> {
  MyPageRoute(WidgetBuilder builder, [RouteSettings settings]) : super(builder: builder, settings: settings);

  @override
  Widget buildTransitions(BuildContext context, Animation<double> animation, Animation<double> secondaryAnimation, Widget child) {
    
    // return _ZoomPageTransition()
    final route = this;
    return ZoomPageTransitionsBuilder().buildTransitions<T>(route, context, animation, secondaryAnimation, child);
    return super.buildTransitions(context, animation, secondaryAnimation, child);
  }
}