class KYCCountryList {
  int returnCode;
  String message;
  String nonce;
  List<KYCCountry> data;

  KYCCountryList(this.returnCode, this.message, this.nonce, this.data);

  KYCCountryList.fromJson(Map<String, dynamic> json)
      : returnCode = json['returnCode'],
        message = json['message'],
        nonce = json['nonce'],
        data = json['data']
            .map<KYCCountry>((v) => new KYCCountry.fromJson(v))
            .toList();

  Map<String, dynamic> toJson() => {
        'returnCode': this.returnCode,
        'message': this.message,
        'nonce': this.nonce,
        if (this.data.isNotEmpty)
          'data': this.data.map((v) => v.toJson()).toList(),
      };
}

class KYCCountry {
  String code;
  String name;
  String fullName;
  String iso3;
  String number;
  String continentCode;

  KYCCountry(this.code, this.name, this.fullName, this.iso3, this.number,
      this.continentCode);

  KYCCountry.fromJson(Map<String, dynamic> json)
      : code = json['code'],
        name = json['name'],
        fullName = json['fullName'],
        iso3 = json['iso3'],
        number = json['number'],
        continentCode = json['continentCode'];

  Map<String, dynamic> toJson() => {
        'code': this.code,
        'name': this.name,
        'fullName': this.fullName,
        'iso3': this.iso3,
        'number': this.number,
        'continentCode': this.continentCode,
      };
}
