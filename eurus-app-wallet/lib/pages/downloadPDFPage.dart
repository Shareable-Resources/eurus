import 'package:euruswallet/common/commonMethod.dart';
import 'package:easy_localization/easy_localization.dart';

enum ReportType {
  summary,
  assetTxDetail,
}

class DownloadPDFPage extends StatefulWidget {
  DownloadPDFPage({
    Key? key,
    required this.reportType,
  }) : super(key: key);

  final ReportType reportType;

  _DownloadPDFPageState createState() => _DownloadPDFPageState();
}

class _DownloadPDFPageState extends State<DownloadPDFPage> {
  Color get themeColor => common.getBackGroundColor();

  bool readyToOpen = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          image: DecorationImage(
            image: AssetImage(
              'images/backgroundImage${!isCentralized() ? '2' : ''}.png',
              package: 'euruswallet',
            ),
            fit: BoxFit.cover,
            alignment: Alignment.topCenter,
          ),
        ),
        child: Column(
          children: [
            AppBar(
              title: Text(
                'DOWNLOAD_STATEMENT.${widget.reportType == ReportType.assetTxDetail ? 'TX_DETAIL' : 'SUMMARY'}.TITLE'
                    .tr(),
              ),
              backgroundColor: Colors.transparent,
              elevation: 0,
            ),
            Expanded(
              flex: 1,
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.only(
                    topLeft: Radius.circular(15),
                    topRight: Radius.circular(15),
                  ),
                ),
                margin: EdgeInsets.only(top: 12),
                child: Container(
                  padding:
                      EdgeInsets.only(top: 35, left: 35, right: 35, bottom: 20),
                  child: Column(
                    children: [
                      readyToOpen ? Text(''.tr()) : Container(),
                      SizedBox(
                        width: MediaQuery.of(context).size.width * 0.27,
                        child: Image.asset('images/report.png',
                            package: 'euruswallet'),
                      ),
                      SizedBox(height: 60),
                      _genReportDateRangeCard(
                        'DOWNLOAD_STATEMENT.CURRENT_MONTH'.tr(),
                        DateFormat('MM/yyyy').format(DateTime.now()),
                        () {},
                      ),
                      _genReportDateRangeCard(
                        'DOWNLOAD_STATEMENT.PREVIOUS_MONTH'.tr(),
                        DateFormat('MM/yyyy').format(DateTime(
                            DateTime.now().year, DateTime.now().month - 1, 1)),
                        () {},
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _genReportDateRangeCard(
      String title, String dateRange, Function exportFnc) {
    return Container(
      margin: EdgeInsets.symmetric(vertical: 10),
      padding: EdgeInsets.symmetric(vertical: 12, horizontal: 15),
      width: double.infinity,
      alignment: Alignment.centerLeft,
      decoration: BoxDecoration(
        borderRadius: FXUI.cricleRadius,
        boxShadow: [BoxShadow(color: FXColor.grey80Color, offset: Offset(1, 2), blurRadius: 8)],
        color: Colors.white,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: EdgeInsets.only(left: 10, right: 10, top: 8, bottom: 20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: FXUI.normalTextStyle.copyWith(fontWeight: FontWeight.w500, fontSize: 16)),
                SizedBox(height: 20),
                Text(dateRange, style: FXUI.normalTextStyle.copyWith(fontSize: 14, color: FXColor.textGray)),
              ],
            ),
          ),
          ElevatedButton(
            onPressed: () => exportFnc,
            style: ButtonStyle(
              shape: MaterialStateProperty.all(RoundedRectangleBorder(borderRadius: FXUI.cricleRadius)),
              backgroundColor: MaterialStateProperty.resolveWith((states) => themeColor),
              elevation: MaterialStateProperty.resolveWith((states) => 0),
            ),
            child: SizedBox(
              width: double.infinity,
              child: Padding(
                padding: EdgeInsets.all(12),
                child: Center(child: Text('DOWNLOAD_STATEMENT.EXPORT_PDF'.tr(), style: FXUI.normalTextStyle.copyWith(fontSize: 16))),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
