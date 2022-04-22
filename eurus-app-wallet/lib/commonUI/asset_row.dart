import 'package:flutter/material.dart';
import 'package:flutter_reorderable_list/flutter_reorderable_list.dart';
import 'package:euruswallet/commonUI/constant.dart';

class AssetRow extends StatelessWidget {
  AssetRow({
    required this.child,
    this.state,
    this.vPadding,
  });

  final Widget child;
  final ReorderableItemState? state;
  final double? vPadding;

  @override
  Widget build(BuildContext _) {
    return Container(
      child: Opacity(
        opacity: state == ReorderableItemState.placeholder ? 0.0 : 1.0,
        child: Container(
          decoration: BoxDecoration(
            borderRadius: FXUI.cricleRadius,
            boxShadow: [
              BoxShadow(
                  color: FXColor.grey80Color,
                  offset: Offset(1, 2),
                  blurRadius: 8)
            ],
            color: Colors.white,
          ),
          padding: EdgeInsets.symmetric(vertical: vPadding ?? 22),
          child: child,
        ),
      ),
    );
  }
}
