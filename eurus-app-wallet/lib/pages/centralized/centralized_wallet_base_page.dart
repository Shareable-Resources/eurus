import 'package:flutter/material.dart';

class CentralizedWalletBasePage extends StatefulWidget {
  final Text? appBarTitle;
  final Widget? body;

  const CentralizedWalletBasePage({
    Key? key,
    this.appBarTitle,
    this.body,
  }) : super(key: key);

  @override
  _CentralizedWalletBasePageState createState() =>
      _CentralizedWalletBasePageState();
}

class _CentralizedWalletBasePageState extends State<CentralizedWalletBasePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
          title: widget.appBarTitle,
          backgroundColor: Colors.transparent,
          elevation: 0,
          centerTitle: true),
      extendBodyBehindAppBar: true,
      backgroundColor: Colors.transparent,
      body: Builder(builder: (context) => _buildBody(context)),
    );
  }

  Widget _buildBody(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        image: DecorationImage(
          image: AssetImage(
            'images/bgCentralized.png',
            package: 'euruswallet',
          ),
          fit: BoxFit.cover,
          alignment: Alignment.topCenter,
        ),
      ),
      child: Padding(
        padding: EdgeInsets.only(
            top: MediaQuery.of(context).padding.top > 0
                ? MediaQuery.of(context).padding.top
                : 80),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: <Widget>[
            Expanded(
              child: DecoratedBox(
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadiusDirectional.only(
                    topStart: Radius.circular(30),
                    topEnd: Radius.circular(30),
                  ),
                ),
                child: Padding(
                  padding: EdgeInsets.only(top: 36)
                      .add(EdgeInsets.symmetric(horizontal: 16)),
                  child: widget.body,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
