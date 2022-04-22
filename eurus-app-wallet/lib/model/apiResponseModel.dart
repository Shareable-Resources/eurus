abstract class ApiResponseModel {
  final int returnCode;
  final String message;
  final Map<String, dynamic> data;

  ApiResponseModel(this.returnCode, this.message, this.data);

  ApiResponseModel.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        data = json['data'];
}
