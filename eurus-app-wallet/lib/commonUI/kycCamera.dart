import 'package:camera/camera.dart';
import 'package:easy_localization/easy_localization.dart';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:euruswallet/pages/settingSubpages/KYC/kycStatus.dart';
import 'package:image/image.dart' as IMG;

enum CameraSide { back, front }

class KYCCamera extends StatefulWidget {
  final void Function(String) onCapture;
  final Widget? imageMask;
  final CameraSide cameraSide;
  final DocType docType;

  const KYCCamera({
    Key? key,
    required this.onCapture,
    this.imageMask,
    this.docType = DocType.Unknown,
    this.cameraSide = CameraSide.back,
  }) : super(key: key);

  @override
  KYCCameraState createState() => KYCCameraState();
}

class KYCCameraState extends State<KYCCamera> {
  List<CameraDescription> cameras = [];
  CameraController? _controller;
  bool lightOn = false;

  @override
  void initState() {
    super.initState();
    _setupCameras(widget.cameraSide);
  }

  @override
  void dispose() {
    _controller?.dispose();
    super.dispose();
  }

  void setCameraFlash() {
    setState(() {
      lightOn = !lightOn;
    });
    _controller
        ?.setFlashMode(lightOn == true ? FlashMode.torch : FlashMode.off);
  }

  Future<void> _setupCameras(CameraSide cameraSide) async {
    try {
      cameras = await availableCameras();
      if (cameras.isNotEmpty) {
        final controller = new CameraController(
            cameras[cameraSide.index], ResolutionPreset.medium);
        _controller = controller;
        await controller.initialize();

        controller.lockCaptureOrientation(DeviceOrientation.portraitUp);
      }
    } on CameraException catch (_) {}
    if (!mounted) return;
    setState(() {});
  }

  static Future cropSquare(String srcFilePath, Size size) async {
    print(srcFilePath);
    var bytes = await File(srcFilePath).readAsBytes();
    IMG.Image? src = IMG.decodeImage(bytes);

    if (src == null) return;

    var cropSize = src.width - src.width ~/ 3;
    int offsetX = (src.width - min(src.width, src.height)) ~/ 2;
    int offsetY = (src.height - min(src.width, src.height)) ~/ 2;

    IMG.Image destImage =
        IMG.copyCrop(src, offsetX, offsetY, cropSize, cropSize);

    var jpg = IMG.encodeJpg(destImage);
    await File(srcFilePath).writeAsBytes(jpg);
  }

  String getTitle() {
    if (widget.cameraSide == CameraSide.front) {
      return 'KYC.TAKE_SELFIE'.tr();
    } else {
      if (widget.docType == DocType.Passport)
        return 'KYC.SCAN_PHOTO_PAGE'.tr();
      else if (widget.docType == DocType.IdFront)
        return 'KYC.ID_CARD_FRONT'.tr();
      else if (widget.docType == DocType.IdBack) return 'KYC.ID_CARD_BACK'.tr();
    }
    return '';
  }

  @override
  Widget build(BuildContext context) {
    final size = MediaQuery.of(context).size;
    final aspectRatio = _controller?.value.aspectRatio ?? 0;

    return Scaffold(
      appBar: WalletAppBar(
        title: getTitle(),
        backButton: true,
      ),
      body: Builder(
        builder: (context) {
          if (_controller != null && _controller!.value.isInitialized) {
            return Stack(children: <Widget>[
              OverflowBox(
                  maxHeight: size.height,
                  maxWidth: size.height * aspectRatio,
                  child: CameraPreview(_controller!)),
              if (widget.imageMask != null)
                Center(
                  child: widget.imageMask,
                ),
            ]);
          } else {
            // return const Center(child: CircularProgressIndicator());
            return Center(
              child: widget.imageMask,
            );
          }
        },
      ),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerFloat,
      floatingActionButton: Stack(children: <Widget>[
        (widget.cameraSide == CameraSide.back
            ? Padding(
                padding: EdgeInsets.only(left: 20),
                child: Align(
                    alignment: Alignment.bottomLeft,
                    child: FloatingActionButton(
                      heroTag: 'flash',
                      backgroundColor: Colors.transparent,
                      elevation: 0.0,
                      onPressed: () {
                        setCameraFlash();
                      },
                      child: Column(children: [
                        Image.asset('images/cameraLightIcon.png',
                            package: 'euruswallet',
                            width: 15,
                            height: 33,
                            fit: BoxFit.contain),
                        Padding(
                            padding: EdgeInsetsDirectional.only(top: 3),
                            child: Text('${lightOn == true ? "On" : "Off"}'))
                      ]),
                    )))
            : Container()),
        Padding(
            padding: EdgeInsets.only(bottom: 50),
            child: Align(
                alignment: Alignment.bottomCenter,
                child: FloatingActionButton(
                  heroTag: 'capture',
                  onPressed: () async {
                    if (_controller != null &&
                        _controller!.value.isInitialized) {
                      try {
                        final image = await _controller!.takePicture();
                        // widget.onCapture(image.path);
                        await cropSquare(image.path, size);
                        widget.onCapture(image.path);
                      } catch (e) {
                        print(e);
                      }
                    }
                  },
                  child: Center(
                      child: Image.asset('images/takePhotoIcon.png',
                          package: 'euruswallet',
                          width: 60,
                          height: 60,
                          fit: BoxFit.contain)),
                )))
      ]),
    );
  }
}
