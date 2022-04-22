class Ethgasstation {
  int fast;
  int fastest;
  int safeLow;
  int average;

  Ethgasstation(this.fast, this.fastest, this.safeLow, this.average);

  Ethgasstation.fromJson(Map<String, dynamic> json)
      : fast = json['fast'],
        fastest = json['fastest'],
        safeLow = json['safeLow'],
        average = json['average'];

  Map<String, dynamic> toJson() => {
        'fast': this.fast,
        'fastest': this.fastest,
        'safeLow': this.safeLow,
        'average': this.average,
      };
}
