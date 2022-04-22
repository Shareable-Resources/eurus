class AdvanceGenAddressPairArgs {
  AdvanceGenAddressPairArgs({
    required this.mnemonic,
    this.account,
    this.change,
    this.accountIdx,
    this.pw,
  });

  final String mnemonic;
  final int? account;
  final int? change;
  final int? accountIdx;
  final String? pw;
}
