import 'package:web3dart/web3dart.dart';

extension EthereumAddressExtension on EthereumAddress {
  String get eip55TruncatedString {
    final hexEip55String = this.hexEip55;
    return hexEip55String.replaceRange(
      6,
      hexEip55String.length - 4,
      '...',
    );
  }
}
