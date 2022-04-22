import 'dart:async';
import 'dart:io';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'lifecycle_event_handler.dart';
import 'qr_scanner_overlay_shape.dart';
import 'types/barcode.dart';
import 'types/barcode_format.dart';
import 'types/camera.dart';
import 'types/camera_exception.dart';
import 'types/features.dart';

typedef QRViewCreatedCallback = void Function(QRViewController);
typedef PermissionSetCallback = void Function(QRViewController, bool);

/// The [QRView] is the view where the camera
/// and the barcode scanner gets displayed.
class QRView extends StatefulWidget {
  const QRView({
    required Key key,
    required this.onQRViewCreated,
    required this.pgTitle,
    this.overlay,
    this.overlayMargin = EdgeInsets.zero,
    this.cameraFacing = CameraFacing.back,
    this.onPermissionSet,
    this.formatsAllowed,
  }) : super(key: key);

  /// [onQRViewCreated] gets called when the view is created
  final QRViewCreatedCallback onQRViewCreated;

  /// Use [overlay] to provide an overlay for the view.
  /// This can be used to create a certain scan area.
  final QrScannerOverlayShape? overlay;

  /// Use [overlayMargin] to provide a margin to [overlay]
  final EdgeInsetsGeometry overlayMargin;

  /// Set which camera to use on startup.
  ///
  /// [cameraFacing] can either be CameraFacing.front or CameraFacing.back.
  /// Defaults to CameraFacing.back
  final CameraFacing cameraFacing;

  /// Calls the provided [onPermissionSet] callback when the permission is set.
  final PermissionSetCallback? onPermissionSet;

  /// Use [formatsAllowed] to specify which formats needs to be scanned.
  final List<BarcodeFormat>? formatsAllowed;

  /// Scanner Page App Bar Title
  final String pgTitle;

  @override
  State<StatefulWidget> createState() => _QRViewState();
}

class _QRViewState extends State<QRView> with SingleTickerProviderStateMixin {
  MethodChannel? _channel;
  late LifecycleEventHandler _observer;

  late AnimationController _aniCon;
  bool lightOn = false;

  @override
  void initState() {
    super.initState();
    _aniCon = AnimationController(vsync: this, duration: Duration(seconds: 3))
      ..repeat();
    _observer = LifecycleEventHandler(
        resumeCallBack: () async => {
              if (_channel != null)
                {
                  QRViewController.updateDimensions(
                      widget.key as GlobalKey, _channel!,
                      overlay: widget.overlay)
                }
            });
    WidgetsBinding.instance?.addObserver(_observer);
  }

  @override
  Widget build(BuildContext context) {
    return NotificationListener(
      onNotification: onNotification,
      child: SizeChangedLayoutNotifier(
        child: (widget.overlay != null)
            ? _getPlatformQrViewWithOverlay()
            : _getPlatformQrView(),
      ),
    );
  }

  @override
  void dispose() {
    _aniCon.dispose();
    super.dispose();
    WidgetsBinding.instance?.removeObserver(_observer);
  }

  bool onNotification(notification) {
    Future.microtask(() => {
          if (_channel != null)
            QRViewController.updateDimensions(
              widget.key as GlobalKey,
              _channel!,
              overlay: widget.overlay,
            )
        });

    return false;
  }

  Widget _getPlatformQrViewWithOverlay() {
    return Stack(
      children: [
        _loadingScreen(),
        _getPlatformQrView(),
        Container(
          padding: widget.overlayMargin,
          decoration: ShapeDecoration(
            shape: widget.overlay!,
          ),
        ),
        _getScanAnimation(),
      ],
    );
  }

  Widget _loadingScreen() {
    return Container(
      alignment: Alignment(0, 0),
      color: Colors.black,
      child: Icon(
        Icons.qr_code_rounded,
        color: Colors.white,
        size: widget.overlay!.cutOutSize * 0.7,
      ),
    );
  }

  Widget _getPlatformQrView() {
    Widget _platformQrView;
    switch (defaultTargetPlatform) {
      case TargetPlatform.android:
        _platformQrView = AndroidView(
          viewType: 'net.touchcapture.qr.flutterqr/qrview',
          onPlatformViewCreated: _onPlatformViewCreated,
          creationParams:
              _QrCameraSettings(cameraFacing: widget.cameraFacing).toMap(),
          creationParamsCodec: StandardMessageCodec(),
        );
        break;
      case TargetPlatform.iOS:
        _platformQrView = UiKitView(
          viewType: 'net.touchcapture.qr.flutterqr/qrview',
          onPlatformViewCreated: _onPlatformViewCreated,
          creationParams:
              _QrCameraSettings(cameraFacing: widget.cameraFacing).toMap(),
          creationParamsCodec: StandardMessageCodec(),
        );
        break;
      default:
        throw UnsupportedError(
            "Trying to use the default webview implementation for $defaultTargetPlatform but there isn't a default one");
    }
    return _platformQrView;
  }

  Widget _getScanAnimation() {
    return Center(
      child: Container(
        width: widget.overlay!.cutOutSize,
        height: widget.overlay!.cutOutSize,
        alignment: Alignment(0, -1),
        child: AnimatedBuilder(
          animation: _aniCon.view,
          builder: (_, __) {
            var scanAniRatio = 0.25;
            var scanAniSize = widget.overlay!.cutOutSize * scanAniRatio;
            var yOffset =
                ((_aniCon.value - scanAniRatio) / (1 - scanAniRatio)) *
                    widget.overlay!.cutOutSize;

            return Transform.translate(
              offset: Offset(0, _aniCon.value < scanAniRatio ? 0 : yOffset),
              child: Container(
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.topCenter,
                    end: Alignment.bottomCenter,
                    colors: [
                      widget.overlay!.borderColor.withOpacity(
                          _aniCon.value >= scanAniRatio
                              ? 0
                              : (1 - (_aniCon.value * 4))),
                      widget.overlay!.borderColor
                    ],
                  ),
                ),
                height: _aniCon.value < scanAniRatio
                    ? (_aniCon.value * widget.overlay!.cutOutSize)
                    : _aniCon.value > (1 - scanAniRatio) &&
                            (widget.overlay!.cutOutSize - yOffset) < scanAniSize
                        ? widget.overlay!.cutOutSize - yOffset
                        : scanAniSize,
                width: widget.overlay!.cutOutSize,
                child: Padding(
                  padding: EdgeInsets.zero,
                ),
              ),
            );
          },
        ),
      ),
    );
  }

  void _onPlatformViewCreated(int id) {
    final channel = MethodChannel('net.touchcapture.qr.flutterqr/qrview_$id');
    _channel = channel;

    // Start scan after creation of the view
    final controller = QRViewController._(channel, widget.key as GlobalKey,
        widget.onPermissionSet, widget.cameraFacing)
      .._startScan(
          widget.key as GlobalKey, widget.overlay, widget.formatsAllowed);

    // Initialize the controller for controlling the QRView
    widget.onQRViewCreated(controller);
  }
}

class _QrCameraSettings {
  _QrCameraSettings({
    required this.cameraFacing,
  });

  final CameraFacing cameraFacing;

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'cameraFacing': cameraFacing.index,
    };
  }
}

class QRViewController {
  QRViewController._(MethodChannel channel, GlobalKey qrKey,
      PermissionSetCallback? onPermissionSet, CameraFacing cameraFacing)
      : _channel = channel,
        _cameraFacing = cameraFacing {
    _channel.setMethodCallHandler((call) async {
      switch (call.method) {
        case 'onRecognizeQR':
          if (call.arguments != null) {
            final args = call.arguments as Map;
            final code = args['code'] as String;
            final rawType = args['type'] as String;
            // Raw bytes are only supported by Android.
            final rawBytes = args['rawBytes'] is List<int>
                ? args['rawBytes'] as List<int>
                : null;
            final format = BarcodeTypesExtension.fromString(rawType);
            if (format != null) {
              final barcode = Barcode(code, format, rawBytes ?? []);
              _scanUpdateController.sink.add(barcode);
            } else {
              throw Exception('Unexpected barcode type $rawType');
            }
          }
          break;
        case 'onPermissionSet':
          await getSystemFeatures(); // if we have no permission all features will not be avaible
          if (call.arguments != null) {
            if (call.arguments as bool) {
              _hasPermissions = true;
            } else {
              _hasPermissions = false;
            }
            if (onPermissionSet != null) {
              onPermissionSet(this, call.arguments as bool);
            }
          }
          break;
      }
    });
  }

  final MethodChannel _channel;
  final CameraFacing _cameraFacing;
  final StreamController<Barcode> _scanUpdateController =
      StreamController<Barcode>();

  Stream<Barcode> get scannedDataStream => _scanUpdateController.stream;

  late bool _hasPermissions;

  bool get hasPermissions => _hasPermissions;

  /// Starts the barcode scanner
  Future<void> _startScan(
    GlobalKey key,
    QrScannerOverlayShape? overlay,
    List<BarcodeFormat>? barcodeFormats,
  ) async {
    // We need to update the dimension before the scan is started.
    try {
      await QRViewController.updateDimensions(key, _channel, overlay: overlay);
      return await _channel.invokeMethod(
          'startScan', barcodeFormats?.map((e) => e.asInt()).toList() ?? []);
    } on PlatformException catch (e) {
      throw CameraException(e.code, e.message);
    }
  }

  /// Gets information about which camera is active.
  Future<CameraFacing> getCameraInfo() async {
    try {
      var cameraFacing = await _channel.invokeMethod('getCameraInfo') as int;
      if (cameraFacing == -1) return _cameraFacing;
      return CameraFacing
          .values[await _channel.invokeMethod('getCameraInfo') as int];
    } on PlatformException catch (e) {
      throw CameraException(e.code, e.message);
    }
  }

  /// Flips the camera between available modes
  Future<CameraFacing> flipCamera() async {
    try {
      return CameraFacing
          .values[await _channel.invokeMethod('flipCamera') as int];
    } on PlatformException catch (e) {
      throw CameraException(e.code, e.message);
    }
  }

  /// Get flashlight status
  Future<bool> getFlashStatus() async {
    try {
      return await _channel.invokeMethod('getFlashInfo');
    } on PlatformException catch (e) {
      throw CameraException(e.code, e.message);
    }
  }

  /// Toggles the flashlight between available modes
  Future<void> toggleFlash() async {
    try {
      await _channel.invokeMethod('toggleFlash') as bool;
    } on PlatformException catch (e) {
      throw CameraException(e.code, e.message);
    }
  }

  /// Pauses the camera and barcode scanning
  Future<void> pauseCamera() async {
    try {
      await _channel.invokeMethod('pauseCamera');
    } on PlatformException catch (e) {
      throw CameraException(e.code, e.message);
    }
  }

  /// Stops barcode scanning and the camera
  Future<void> stopCamera() async {
    try {
      await _channel.invokeMethod('stopCamera');
    } on PlatformException catch (e) {
      throw CameraException(e.code, e.message);
    }
  }

  /// Resumes barcode scanning
  Future<void> resumeCamera() async {
    try {
      await _channel.invokeMethod('resumeCamera');
    } on PlatformException catch (e) {
      throw CameraException(e.code, e.message);
    }
  }

  /// Returns which features are available on device.
  Future<SystemFeatures> getSystemFeatures() async {
    try {
      var features =
          await _channel.invokeMapMethod<String, dynamic>('getSystemFeatures');
      return SystemFeatures.fromJson(features);
    } on PlatformException catch (e) {
      throw CameraException(e.code, e.message);
    }
  }

  /// Stops the camera and disposes the barcode stream.
  void dispose() {
    if (Platform.isIOS) stopCamera();
    _scanUpdateController.close();
  }

  /// Updates the view dimensions for iOS.
  static Future<void> updateDimensions(
    GlobalKey key,
    MethodChannel channel, {
    QrScannerOverlayShape? overlay,
  }) async {
    if (defaultTargetPlatform == TargetPlatform.iOS) {
      // Add small delay to ensure the renderbox is loaded
      await Future.delayed(Duration(milliseconds: 100));
      final renderBox = key.currentContext?.findRenderObject() as RenderBox;
      try {
        await channel.invokeMethod('setDimensions', {
          'width': renderBox.size.width,
          'height': renderBox.size.height,
          'scanArea': overlay?.cutOutSize ?? 0,
          'scanAreaOffset': overlay?.cutOutBottomOffset ?? 0
        });
      } on PlatformException catch (e) {
        throw CameraException(e.code, e.message);
      }
    }
  }
}
