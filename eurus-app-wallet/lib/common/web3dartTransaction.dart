import 'dart:typed_data';
import 'package:euruswallet/common/commonMethod.dart';
import 'package:web3dart/crypto.dart';
import 'package:web3dart/web3dart.dart';
import 'package:web3dart/src/utils/length_tracking_byte_sink.dart';

extension ContractFunctionEncode on ContractFunction {
  Uint8List encodeCallWithKeccak256(
    List<dynamic> params, {
    List<String>? paramsName,
  }) {
    final sink = LengthTrackingByteSink();
    List<AbiType> parameters = this
        .parameters
        .where((element) {
          if ((paramsName ?? []).isNotEmpty)
            return (paramsName ?? []).contains(element.name);
          return true;
        })
        .map((e) => e.type)
        .toList();

    if (params.length != parameters.length) {
      throw ArgumentError.value(
          params.length, 'params', 'Must match function parameters');
    }

    TupleType(parameters).encode(params, sink);
    Uint8List sinkData = sink.asBytes();
    print("sinkData:$sinkData");
    return keccak256(sinkData);
  }

  List<dynamic> decodeReturnValuesFromParameters(
    Uint8List data,
  ) {
    final tuple = TupleType(this.parameters.map((p) => p.type).toList());
    final buffer = data.buffer;

    final parsedData = tuple.decode(buffer, 0);
    return parsedData.data;
  }
}
