import * as crypto from '../../utils/crypto/Crypto';

export function testCryptoSign() {
  console.log('testCryptoSign');
  let sign: string = crypto.Crypto.generateRSASignFromBase64(
    'MIICXAIBAAKBgQDyWeqeXzTbce2sEDkWTdp7ohs4sQrG4kSybYUyPJ+8OuUdwZzBGJtKVrE7ib/mQ91qi5Mi1y7eD7i+opSAb7riBPU9KqRK/HR0wrS4Ty5pVC1El+L16Bzv3rdiAiauD788OA+6IyYR6ym9/pqjraEBYNkiq6lVrQmr+bVPklL8WQIDAQABAoGBAI0b+7bw+mbXxCpOW973v8IC1lTBCORadONrT4W7tMo1hJg51lq8SBcVuSZQg3LAixhxYl+D6V8Uyl39bfgGlknHw0p+hHeLpGQVwqk7tlgLOe1jdHVNpa3XkYS9XWdMaR3xP3u1UDSNAmmlH9F7M9R/lavpBua2bcN4xUoIDJUBAkEA9muOG2leCxZQ6LBBCJ1AfJk5DiMTP6MpJKvMNbmTg7VaVcNEvAG/ah8onE/RQn8oKCR2FF+Nt3OIWtvytzU5uQJBAPvF3dInXns6CQdpI0XKqXxVHo/EsntRrm7AE2vKxhknoiXwutjE2A3b02IALRSncTnfN8elc5S1aM3yPXvCp6ECQHyWjqffUzQJ2Wh8TLU9RHB6y7URGNdQ+ClZTwtOw32RAJZh/uuLeqr+C/tUA7o6LiR1otnUIgPBvwwpLoNgvskCQDBf5C7swYqLrBBwcwaF2eq5sBWPhXuMRiehYBnl7AoJQ25mwn/D0n9XLVY3EnJVebvU17LUZ1C7SlfHo2iUDgECQCp2bJ3NsYVCMlIaHUND0AHSa5K4HRaUXKsFtPvKO2U/OYukothq1zyf0/so4pwTY7i0gl05jL4/3Ck/inDb0B8=',
    Buffer.from('HelloWorld'),
  );
  console.log(sign);
}
