import 'package:flutter/material.dart';

import 'cus_modal.dart';

class PermModalTemplate extends CustomModal {
  PermModalTemplate({
    required this.icon,
    required this.title,
    required this.desc,
    this.color,
    this.iconColor,
    this.acceptText,
    this.declineText,
    this.hideDecline,
    this.otherAction,
  });

  // Main Color of this modal
  final Color? color;

  // Informations that to be shown on this page
  final IconData icon;
  final String title;
  final String desc;
  final Color? iconColor;
  final String? acceptText;
  final String? declineText;

  // Hide decline for one action btn only
  final bool? hideDecline;

  // Custom actions
  final Widget? otherAction;

  @override
  Widget buildPage(
    BuildContext context,
    Animation<double> animation,
    Animation<double> secondaryAnimation,
  ) {
    return Material(
      type: MaterialType.transparency,
      // make sure that the overlay content is not cut off
      child: Column(
        children: [
          AppBar(
            backgroundColor: Colors.transparent,
            shadowColor: Colors.transparent,
          ),
          Expanded(
            flex: 1,
            child: Container(
              padding: EdgeInsets.all(35),
              decoration: BoxDecoration(
                borderRadius: BorderRadius.only(
                  topLeft: Radius.circular(15),
                  topRight: Radius.circular(15),
                ),
                color: Colors.white,
              ),
              alignment: Alignment.center,
              child: SingleChildScrollView(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(
                      icon,
                      color: iconColor ?? Color.fromRGBO(0, 159, 221, 1),
                      size: MediaQuery.of(context).size.width / 3,
                    ),
                    _cusSpacer(),
                    Text(
                      title,
                      style:
                          TextStyle(fontWeight: FontWeight.bold, fontSize: 21),
                      textAlign: TextAlign.center,
                    ),
                    _cusSpacer(p: 15),
                    Text(desc, textAlign: TextAlign.center),
                    _cusSpacer(p: 35),
                    DecoratedBox(
                      decoration: BoxDecoration(
                        color: color ?? Color.fromRGBO(0, 159, 221, 1),
                        borderRadius: BorderRadius.all(Radius.circular(15)),
                      ),
                      child: TextButton(
                        onPressed: () {
                          Navigator.pop(context, true);
                        },
                        child: Padding(
                          padding: EdgeInsets.all(12),
                          child: Text(acceptText ?? 'OK',
                              style: TextStyle(color: Colors.white)),
                        ),
                      ),
                    ),
                    _cusSpacer(p: 35),
                    if (hideDecline != false)
                      Padding(
                        padding: EdgeInsets.all(12),
                      )
                    else
                      TextButton(
                        onPressed: () {
                          Navigator.pop(context);
                          },
                        child: Padding(
                          padding: EdgeInsets.zero,
                          child: Text(
                            declineText ?? 'Don\'t Allow',
                            style: TextStyle(
                              color: color ?? Color.fromRGBO(0, 159, 221, 1),
                            ),
                          ),
                        ),
                      ),
                    otherAction ?? Padding(padding: EdgeInsets.zero),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _cusSpacer({
    double? p,
  }) {
    return Padding(padding: EdgeInsets.only(top: p ?? 5));
  }
}
