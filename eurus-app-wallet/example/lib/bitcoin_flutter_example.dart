import 'package:bip39/bip39.dart' as bip39;
import 'package:euruswallet/bitcoinLibrary/bitcoin_flutter.dart';
import 'package:euruswallet/bitcoinLibrary/src/utils/magic_hash.dart';
import 'package:hex/hex.dart';

main() {
  var seed = bip39.mnemonicToSeed(
      'carbon shuffle shoot knock alter bottom polar maple husband poet match spring');
  var hdWallet = new HDWallet.fromSeed(seed,network: testnet);
  print('hdWallet.address:${hdWallet.address}');
  // => 12eUJoaWBENQ3tNZE52ZQaHqr3v4tTX4os
  print('hdWallet.pubKey:${hdWallet.pubKey}');
  // => 0360729fb3c4733e43bf91e5208b0d240f8d8de239cff3f2ebd616b94faa0007f4
  print('hdWallet.privKey:${hdWallet.privKey}');
  // => 01304181d699cd89db7de6337d597adf5f78dc1f0784c400e41a3bd829a5a226
  print('hdWallet.wif:${hdWallet.wif}');
  // => KwG2BU1ERd3ndbFUrdpR7ymLZbsd7xZpPKxsgJzUf76A4q9CkBpY

  var wallet =
      Wallet.fromWIF('cNR1r2WgKiTdTd1VwhvBbYbJcsERDYyYXzA2FxWtGNxcNWjjZV2k',testnet);
  print('wallet.address:${wallet.address}');
  // => 19AAjaTUbRjQCMuVczepkoPswiZRhjtg31
  print('wallet.pubKey:${wallet.pubKey}');
  // => 03aea0dfd576151cb399347aa6732f8fdf027b9ea3ea2e65fb754803f776e0a509
  print('wallet.privKey:${wallet.privKey}');
  // => 3095cb26affefcaaa835ff968d60437c7c764da40cdd1a1b497406c7902a8ac9
  print('wallet.pubKey:${wallet.wif}');
  // => Kxr9tQED9H44gCmp6HAdmemAzU3n84H3dGkuWTKvE23JgHMW8gct


  final alice = ECPair.fromWIF(
      'cNR1r2WgKiTdTd1VwhvBbYbJcsERDYyYXzA2FxWtGNxcNWjjZV2k',network: testnet);


  dynamic sighString1 = signString(alice: alice,toSignString: '528987d88e56d9246912febc127829f3c3c6ecf288ebea05e1dcfa14eb488a8c');
  dynamic sighString2 = signString(alice: alice,toSignString: '187005b2370ceeeabe4f618fb2c979afa06d0ed427b5380e58801a93c30eb0e1');
  dynamic sighString3 = signString(alice: alice,toSignString: '1fc144a7a8cae278ba889427de3fef97357f3deca363035aa59f0ba36a15c071');

  final txb = new TransactionBuilder(network: testnet);

  txb.setVersion(1);
  txb.addInput(
      'de6d41e1ee517e3f37e92c31ab447c6547309e8ce156b8ee77cace6bd8a03c08',
      0); // Alice's previous transaction output, has 15000 satoshis
  txb.addOutput(
      'mjaufUtcBUdZpx3AFxwjNX2fW3vQqhy5wD',
      93000);
  txb.addOutput('mhARbrfUzFoepzrAwdzwEVWAi3Wmq22FmR', 200);
  // (in)15000 - (out)12000 = (fee)3000, this is the miner fee

  txb.sign(vin: 0, keyPair: alice);
  String transcationData = txb.build().toHex();
  print("txb.build().toHex():${transcationData}");
  if("01000000014eb24795c8b6ec53757023cc37619a9beaacc077b5eb047989936bebe96946d8000000006a47304402203b080cbb2d10d93f4230c38f040483ef16df6c4280dbee5d6323b78ba295bb4d0220334b7296b5185fbb4dc3fe0657d8433098fb673a46fd70d5cb9f72326771b7ad01210360729fb3c4733e43bf91e5208b0d240f8d8de239cff3f2ebd616b94faa0007f4ffffffff0200770100000000001976a9142c9ff89a399da4d2297e587ebbb8e49f0adad70788ac90010000000000001976a914120e4fa1d95fa9f6a7581fe2e267715020bc2f7f88ac00000000" == txb.build().toHex()){
    print("true");
  } else {
    print("false");
  }

}

String signString({required ECPair alice,required String toSignString}){

  dynamic messageHash = magicHash(toSignString, testnet);
  dynamic sighStringByte1 = alice.sign(messageHash);
  String publicKey = HEX.encode(alice.publicKey!);
  print('publicKey:$publicKey');
  dynamic sighString =    HEX.encode(sighStringByte1);
  print('sighString:$sighString');
  return sighString;
}
