import 'package:flutter/material.dart';
import 'package:meta/meta.dart';
import 'src/nav_button.dart';
import 'src/nav_custom_painter.dart';

typedef _LetIndexPage = bool Function(int value);

class CurvedNavigationBar extends StatefulWidget {
  final List<Widget> items;
  final int index;
  final Color color;
  final Color buttonBackgroundColor;
  final Color backgroundColor;
  final ValueChanged<int> onTap;
  final _LetIndexPage letIndexChange;
  final Curve animationCurve;
  final Duration animationDuration;
  final double height;
  final Gradient backgroundGradient;
  final Gradient gradient;
  final Shader shader;
  final BoxDecoration decoration;
  final dynamic Function(int) onAnimationCompleted;

  CurvedNavigationBar({
    Key key,
    @required this.items,
    this.index = 0,
    this.color = Colors.white,
    this.buttonBackgroundColor,
    this.backgroundColor = Colors.blueAccent,
    this.onTap,
    _LetIndexPage letIndexChange,
    this.animationCurve = Curves.easeOut,
    this.animationDuration = const Duration(milliseconds: 600),
    this.height = 75.0,
    this.backgroundGradient,
    this.gradient,
    this.shader,
    this.decoration,
    this.onAnimationCompleted,
  })  : letIndexChange = letIndexChange ?? ((_) => true),
        assert(items != null),
        assert(items.length >= 1),
        assert(0 <= index && index < items.length),
        assert(0 <= height && height <= 75.0),
        super(key: key);

  @override
  CurvedNavigationBarState createState() => CurvedNavigationBarState();
}

class CurvedNavigationBarState extends State<CurvedNavigationBar>
    with SingleTickerProviderStateMixin {
  double _startingPos;
  int _endingIndex = 0;
  double _pos;
  double _buttonHide = 0;
  Widget _icon;
  AnimationController _animationController;
  int _length;

  @override
  void initState() {
    super.initState();
    _icon = widget.items[widget.index];
    _length = widget.items.length;
    _pos = widget.index / _length;
    _startingPos = widget.index / _length;
    _animationController = AnimationController(vsync: this, value: _pos);
    _animationController.addListener(() {
      print('${_animationController.value}');
      _pos = _animationController.value;
      final endingPos = _endingIndex / widget.items.length;
      final middle = (endingPos + _startingPos) / 2;
      Future.delayed(Duration.zero, () => 
        setState(() {
          // _pos = _animationController.value;
          // final endingPos = _endingIndex / widget.items.length;
          // final middle = (endingPos + _startingPos) / 2;
          if ((endingPos - _pos).abs() < (_startingPos - _pos).abs()) {
            _icon = widget.items[_endingIndex];
          }
          _buttonHide =
              (1 - ((middle - _pos) / (_startingPos - middle)).abs()).abs();
          // _buttonHide =
          //     (1 - (((((_endingIndex / widget.items.length) + _startingPos) / 2) - _pos) / (_startingPos - (((_endingIndex / widget.items.length) + _startingPos) / 2))).abs()).abs();
        })
      );
    });
  }

  @override
  void didUpdateWidget(CurvedNavigationBar oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.index != widget.index) {
      final newPosition = widget.index / _length;
      _startingPos = _pos;
      _endingIndex = widget.index;
      _animationController.animateTo(newPosition,
          duration: widget.animationDuration, curve: widget.animationCurve);
    }
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    Size size = MediaQuery.of(context).size;
    return Container(
      decoration: widget.decoration??BoxDecoration(boxShadow: <BoxShadow>[BoxShadow(color: Colors.transparent)], color: widget.backgroundColor, gradient: widget.backgroundGradient),
      // color: widget.backgroundColor,
      height: widget.height,
      child: Stack(
        clipBehavior: Clip.none,
        alignment: Alignment.bottomCenter,
        children: <Widget>[
          AnimatedBuilder(animation: _animationController, builder: (_, child) => Positioned(
            bottom: -40 - (75.0 - widget.height),
            left: Directionality.of(context) == TextDirection.rtl
                ? null
                : _animationController.value * size.width,
            right: Directionality.of(context) == TextDirection.rtl
                ? _animationController.value * size.width
                : null,
            width: size.width / _length,
            child: child,
          ), child: Center(
            child: AnimatedBuilder(animation: _animationController, builder: (_, child) => Transform.translate(
              offset: Offset(
                0,
                -(1 - _buttonHide) * 80,
                // -(1 - (1 - ((middle - _animationController.value) / (_startingPos - middle)).abs()).abs()) * 80,
                // -(1 - ((1 - (((((_endingIndex / widget.items.length) + _startingPos) / 2) - _animationController.value) / (_startingPos - (((_endingIndex / widget.items.length) + _startingPos) / 2))).abs()).abs())) * 80,
              ),
              child: child
            ), child: Material(
                color: widget.buttonBackgroundColor ?? widget.color,
                type: MaterialType.circle,
                child: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: _icon,
                ),
              ),
            ),
          )),
          Positioned(
            left: 0,
            right: 0,
            bottom: 0 - (75.0 - widget.height),
            child: AnimatedBuilder(animation: _animationController, builder: (_, child) => CustomPaint(
              painter: NavCustomPainter(
                  _animationController.value, _length, widget.color, Directionality.of(context), widget.shader),
              child: child,
            ), child: Container(
                height: 75.0,
            )),
          ),
          Positioned(
            left: 0,
            right: 0,
            bottom: 0 - (75.0 - widget.height),
            child: SizedBox(
                height: 100.0,
                child: Row(
                    children: widget.items.map((item) {
                  return AnimatedBuilder(animation: _animationController, builder: (_, child) => NavButton(
                    onTap: _buttonTap,
                    position: _animationController.value,
                    length: _length,
                    index: widget.items.indexOf(item),
                    child: child,
                  ), child: Center(child: item));
                }).toList())),
          ),
        ],
      ),
    );
  }

  void setPage(int index) {
    _buttonTap(index);
  }

  void _buttonTap(int index) {
    if (!widget.letIndexChange(index)) {
      return;
    }
    if (widget.onTap != null) {
      widget.onTap(index);
    }
    final newPosition = index / _length;
    // setState(() {
      _startingPos = _pos;
      _endingIndex = index;
      _animationController.animateTo(newPosition,
          duration: widget.animationDuration, curve: widget.animationCurve).whenComplete(() => widget.onAnimationCompleted(index));
    // });
  }
}
