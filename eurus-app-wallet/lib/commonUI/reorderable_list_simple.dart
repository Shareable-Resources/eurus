import 'package:flutter/material.dart' hide ReorderableList;
import 'package:flutter_reorderable_list/flutter_reorderable_list.dart';
import 'asset_row.dart';

enum ReorderableListSimpleSide { Right, Left }

class ReorderableListSimple extends StatefulWidget {
  ReorderableListSimple({
    required this.children,
    required this.onReorder,
    this.allowReordering = true,
    this.childrenAlreadyHaveListener = false,
    this.handleSide = ReorderableListSimpleSide.Right,
    this.handleIcon,
    this.padding,
  });

  final bool allowReordering;
  final bool childrenAlreadyHaveListener;
  final ReorderableListSimpleSide handleSide;
  final Widget? handleIcon;
  final List<Widget> children;
  final ReorderCallback onReorder;
  final EdgeInsets? padding;

  @override
  State<ReorderableListSimple> createState() =>
      new _ReorderableListSimpleState();
}

class _ReorderableListSimpleState extends State<ReorderableListSimple> {
  int? _newIndex;
  late List<Widget> _children;

  @override
  void initState() {
    super.initState();
    _children = List<Widget>.from(widget.children);
  }

  @override
  didUpdateWidget(ReorderableListSimple oldWidget) {
    super.didUpdateWidget(oldWidget);
    _children = List<Widget>.from(widget.children);
  }

  int _oldIndexOfKey(Key key) {
    return widget.children
        .indexWhere((Widget w) => Key(w.hashCode.toString()) == key);
  }

  int _indexOfKey(Key key) {
    return _children
        .indexWhere((Widget w) => Key(w.hashCode.toString()) == key);
  }

  Widget _buildReorderableItem(BuildContext context, int index) {
    return ReorderableItemSimple(
      key: Key(_children[index].hashCode.toString()),
      innerItem: _children[index],
      allowReordering: widget.allowReordering,
      childrenAlreadyHaveListener: widget.childrenAlreadyHaveListener,
      handleSide: widget.handleSide,
      handleIcon: widget.handleIcon,
    );
  }

  @override
  Widget build(BuildContext context) {
    return ReorderableList(
      decoratePlaceholder: (w, d) => DecoratedPlaceholder(offset: 0, widget: w),
      onReorder: (Key draggedItem, Key newPosition) {
        int draggingIndex = _indexOfKey(draggedItem);
        int newPositionIndex = _indexOfKey(newPosition);

        final item = _children[draggingIndex];
        setState(() {
          _newIndex = newPositionIndex;

          _children.removeAt(draggingIndex);
          _children.insert(newPositionIndex, item);
        });

        return true;
      },
      onReorderDone: (Key draggedItem) {
        int oldIndex = _oldIndexOfKey(draggedItem);
        if (_newIndex != null) widget.onReorder(oldIndex, _newIndex!);
        _newIndex = null;
      },
      child: ListView.builder(
        padding: widget.padding,
        itemCount: _children.length,
        itemBuilder: (BuildContext context, int index) {
          return _buildReorderableItem(context, index);
        },
      ),
    );
  }
}

class ReorderableItemSimple extends StatelessWidget {
  ReorderableItemSimple({
    required Key key,
    required this.innerItem,
    this.allowReordering = true,
    this.childrenAlreadyHaveListener = false,
    this.handleSide = ReorderableListSimpleSide.Right,
    this.handleIcon,
  }) : super(key: key);

  final bool allowReordering;
  final bool childrenAlreadyHaveListener;
  final ReorderableListSimpleSide handleSide;
  final Widget? handleIcon;
  final Widget innerItem;

  Widget _buildInnerItem(BuildContext context) {
    if ((!allowReordering) || childrenAlreadyHaveListener) return innerItem;

    Widget? icon = handleIcon;
    if (icon == null) icon = Icon(Icons.menu);

    var item = Expanded(child: innerItem);
    List<Widget> children = <Widget>[];

    if (handleSide == ReorderableListSimpleSide.Right) children.add(item);
    children.add(ReorderableListener(
      child: Container(
        padding: EdgeInsets.only(left: 25, top: 25, right: 10, bottom: 25),
        alignment: Alignment.centerLeft,
        child: icon,
        color: Colors.transparent,
      ),
    ));
    if (handleSide == ReorderableListSimpleSide.Left) children.add(item);

    final Row row = Row(
      mainAxisAlignment: handleSide == ReorderableListSimpleSide.Right
          ? MainAxisAlignment.end
          : MainAxisAlignment.start,
      mainAxisSize: MainAxisSize.min,
      children: children,
    );

    return row;
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 20, vertical: 8),
      child: ReorderableItem(
        key: key!,
        childBuilder: (BuildContext context, ReorderableItemState state) {
          return AssetRow(child: _buildInnerItem(context), state: state);
        },
      ),
    );
  }
}
