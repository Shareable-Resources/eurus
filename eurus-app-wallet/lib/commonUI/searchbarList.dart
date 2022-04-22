import 'package:euruswallet/common/commonMethod.dart';
import 'package:easy_localization/easy_localization.dart';

class SearchbarList extends StatefulWidget {
  const SearchbarList(
      {Key? key,
      required this.list,
      required this.value,
      required this.onChanged,
      required this.onSelect,
      this.buildListItem,
      this.listMaxHeight,
      this.hintText: ''})
      : super(key: key);

  final List<dynamic> list;
  final dynamic value;
  final void Function(String) onChanged;
  final void Function(int) onSelect;
  final Widget Function(int, bool)? buildListItem;
  final double? listMaxHeight;
  final String hintText;

  @override
  _SearchbarListState createState() => _SearchbarListState();
}

class _SearchbarListState extends State<SearchbarList> {
  bool onFocus = true;
  final FocusNode _focusNode = new FocusNode();
  final TextEditingController _textEditingController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _focusNode.addListener(() {
      onFocus = _focusNode.hasFocus;
    });
  }

  @override
  void dispose() {
    _textEditingController.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (widget.value != null && onFocus == false) {
      _textEditingController.text = widget.value;
    }
    return Container(
      width: double.infinity,
      padding: onFocus == true
          ? const EdgeInsets.all(13.0)
          : const EdgeInsets.all(0.0),
      decoration: onFocus == true
          ? FXUI.boxDecorationWithShadow
          : FXUI.circleBoxDecoration,
      child: Column(
        children: [
          TextField(
            controller: _textEditingController,
            decoration: FXUI.defaultTextFieldInputDecoration.copyWith(
                hintText: widget.hintText, prefixIcon: Icon(Icons.search)),
            maxLines: 1,
            onChanged: (String value) => widget.onChanged(value),
            focusNode: _focusNode,
            autofocus: true,
          ),
          (onFocus == true && widget.list.isNotEmpty
              ? Container(
                  constraints: BoxConstraints(
                      maxHeight: widget.listMaxHeight ?? double.infinity),
                  child: ListView.builder(
                      itemCount: widget.list.length,
                      itemBuilder: (BuildContext context, int index) {
                        bool isSelected = widget.list[index] == widget.value;
                        return GestureDetector(
                            onTap: () {
                              widget.onSelect(index);
                              _focusNode.unfocus();
                            },
                            child: Container(
                              color: isSelected == true
                                  ? FXColor.listItemBlueColor
                                  : Colors.white,
                              padding: const EdgeInsets.all(8.0),
                              child: (widget.buildListItem != null
                                  ? widget.buildListItem!(index, isSelected)
                                  : Flexible(
                                      child: Text(widget.list[index],
                                          style: FXUI.normalTextStyle.copyWith(
                                              color: isSelected == true
                                                  ? FXColor.mainBlueColor
                                                  : Colors.black)))),
                            ));
                      }))
              : Container()),
        ],
      ),
    );
  }
}
