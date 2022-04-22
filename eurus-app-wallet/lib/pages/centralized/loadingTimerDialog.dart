import 'package:euruswallet/common/commonMethod.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter/material.dart';
import 'package:easy_localization/easy_localization.dart';
import 'dart:math';

class LoadingTimerDialog extends StatefulWidget {
  const LoadingTimerDialog({Key? key}) : super(key: key);

  @override
  LoadingTimerDialogState createState() => LoadingTimerDialogState();
}

class LoadingTimerDialogState extends State<LoadingTimerDialog>
    with SingleTickerProviderStateMixin {
  late AnimationController _rotationController;
  late Timer _timer;
  double _seconds = 0;
  bool isCompleted = false;
  static const int defaultLoadingTimer = 60;
  final Tween<double> turnsTween = Tween<double>(
    begin: 1,
    end: 0,
  );

  @override
  void initState() {
    _rotationController = AnimationController(
        duration: const Duration(milliseconds: 1500), vsync: this)
      ..repeat();
    startTimer();
    super.initState();
  }

  @override
  void dispose() {
    _timer.cancel();
    _rotationController.dispose();
    super.dispose();
  }

  void startTimer() {
    const oneSec = const Duration(milliseconds: 100);
    _timer = new Timer.periodic(
      oneSec,
      (Timer timer) {
        if (_seconds >= defaultLoadingTimer) {
          setState(() {
            timer.cancel();
          });
        } else {
          setState(() {
            _seconds = _seconds + 100 / 1000;
          });
        }
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
        onWillPop: () async => false,
        child: Dialog(
            backgroundColor: Colors.transparent,
            insetPadding: EdgeInsets.all(10),
            child: Container(
              width: double.infinity,
              padding: EdgeInsets.fromLTRB(0, 26, 0, 120),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: FXUI.cricleRadius,
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    "LOADING_DIALOG.PLEASE_WAIT".tr(),
                    style: FXUI.titleTextStyle.copyWith(
                      fontSize: 18,
                    ),
                  ),
                  Padding(padding: EdgeInsets.only(bottom: 60)),
                  RotationTransition(
                    turns: turnsTween.animate(_rotationController),
                    child: Image.asset('images/loading_timer_icon.png',
                        width: 95, height: 95, package: 'euruswallet'),
                  ),
                  Padding(padding: EdgeInsets.only(bottom: 10)),
                  Text(
                    isCompleted == true
                        ? '100%'
                        : '${min(99, (_seconds / defaultLoadingTimer) * 100).round()}%',
                    style: FXUI.titleTextStyle.copyWith(
                      fontSize: 24,
                    ),
                  ),
                  Padding(padding: EdgeInsets.only(bottom: 8)),
                  Text(
                    'LOADING_DIALOG.CREATING_WALLET'.tr(),
                    style: FXUI.normalTextStyle.copyWith(
                      fontSize: 14,
                      color: FXColor.textGray,
                    ),
                  ),
                  // TextButton(
                  //   onPressed: () => Navigator.pop(context),
                  //   child: const Text('Cancel'),
                  // ),
                ],
              ),
            )));
  }

  onCompleted() async {
    setState(() {
      isCompleted = true;
    });
    await Future.delayed(const Duration(milliseconds: 500), () {
      Navigator.pop(context);
    });
  }
}
