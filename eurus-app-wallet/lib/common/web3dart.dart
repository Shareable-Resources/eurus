import 'package:apihandler/apiHandler.dart';
import 'package:convert/convert.dart';
import 'package:euruswallet/common/web3dartTransaction.dart';
import 'package:http/http.dart';
import 'package:web3dart/credentials.dart';
import 'package:web3dart/crypto.dart' as _web3dartCrypto;
import 'package:web3dart/crypto.dart';
import 'package:web3dart/src/utils/typed_data.dart';

import 'commonMethod.dart';
import 'logging_client.dart';

export 'package:web3dart/web3dart.dart';

enum BlockChainType { Ethereum, Eurus, BinanceCoin }
enum EnvType { Dev, Staging, Testnet, Production }
EnvType envType = EnvType.Dev;

class Web3dart {
  static final Web3dart _instance = Web3dart._internal();
  Client httpClient = new LoggingClient(Client());
  late Web3Client eurusEthClient;
  late Web3Client mainNetEthClient;
  late Web3Client mainNetBscClient;
  Credentials? credentials;
  late String blockchainRpcUrl;
  late String eurusRPCUrl;
  late String bscRpcUrl;
  int eurusChainId = 2021;
  int mainNetChainId = 4;
  int bscChainId = 97; //Mainnet= 56, Testnet= 97
  EthereumAddress? myEthereumAddress;
  int? estimateMaxGas;
  String? ethBalanceFromEthereum;
  String? erc20TokenBalanceFromEthereum;
  String? ethBalanceFromEurus;
  String? erc20TokenBalanceFromEurus;
  String? lastTxId;
  late DeployedContract erc20ContractFromEthereum;
  late DeployedContract erc20ContractFromEurus;
  late DeployedContract erc20ContractBSC;
  List<dynamic> tokenList = [];
  Map<String, String> tokenListMap = Map<String, String>();
  Map<String, String> bscTokenListMap = {
    'BUSD': '0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee',
    'ETH': '0xd66c6b4f0be8ce5b39d52e0fd1344c389929b378',
  };
  Future<Credentials?> Function() get canGetCredentialsHandler => () async {
        final privateKey = await canGetPrivateKeyHandler();
        if (isEmptyString(string: privateKey)) return null;
        Credentials credentials =
            await mainNetEthClient.credentialsFromPrivateKey(privateKey);
        print("credentials Address:${await credentials.extractAddress()}");
        return credentials;
      };
  late Future<String> Function() canGetPrivateKeyHandler;
  double ethereumGasPrice = 200000000000;
  double bscGasPrice = 20000000000;
  double eurusGasPrice = envType == EnvType.Staging
      ? 15000
      : envType == EnvType.Dev
          ? 2400000000
          : 2400000000;
  double transactionSpeed = 1;
  late String externalSmartContractConfigAddress;
  late String eurusInternalConfigAddress;

  /// init method
  Web3dart._internal();

  factory Web3dart() {
    return _instance;
  }

  /// initEthClient
  Future<bool> initEthClient({
    String? privateKey,
    String? publicAddress,
    Future<String> Function()? canGetPrivateKeyHandler,
  }) async {
    credentials = privateKey != null
        ? await mainNetEthClient.credentialsFromPrivateKey(privateKey)
        : null;
    myEthereumAddress = publicAddress != null
        ? EthereumAddress.fromHex(publicAddress)
        : credentials != null
            ? await credentials!.extractAddress()
            : null;
    print("ethereumAddress:${myEthereumAddress.toString()}");
    // canGetCredentialsHandler = canGetPrivateKeyHandler != null ? () async => await mainNetEthClient.credentialsFromPrivateKey(await canGetPrivateKeyHandler()) : this.canGetCredentialsHandler;
    return true;
  }

  /// setErc20Contract
  DeployedContract setErc20Contract({
    required String contractAddress,
    required BlockChainType blockChainType,
  }) {
    DeployedContract deployedContract;
    if (blockChainType == BlockChainType.Ethereum) {
      erc20ContractFromEthereum =
          getEthereumERC20Contract(contractAddress: contractAddress);
      deployedContract = erc20ContractFromEthereum;
    } else if (blockChainType == BlockChainType.BinanceCoin) {
      erc20ContractBSC =
          getEthereumERC20Contract(contractAddress: contractAddress);
      deployedContract = erc20ContractBSC;
    } else {
      erc20ContractFromEurus =
          getEurusERC20Contract(contractAddress: contractAddress);
      deployedContract = erc20ContractFromEurus;
    }
    return deployedContract;
  }

  /// getBalance
  Future<bool> getErc20Balance({
    required BlockChainType type,
    required bool isEthOrEun,
  }) async {
    if (type == BlockChainType.Ethereum) {
      if (!isEthOrEun) {
        web3dart.erc20TokenBalanceFromEthereum = await web3dart.getERC20Balance(
            blockChainType: BlockChainType.Ethereum,
            deployedContract: web3dart.erc20ContractFromEthereum);
      } else {
        web3dart.ethBalanceFromEthereum = await web3dart.getETHBalance(
            blockChainType: BlockChainType.Ethereum);
      }
    }
    if (type == BlockChainType.Eurus) {
      if (!isEthOrEun) {
        web3dart.erc20TokenBalanceFromEurus = await web3dart.getERC20Balance(
            blockChainType: BlockChainType.Eurus,
            deployedContract: web3dart.erc20ContractFromEurus);
      } else {
        web3dart.ethBalanceFromEurus =
            await web3dart.getETHBalance(blockChainType: BlockChainType.Eurus);
      }
    }
    return true;
  }

  /// estimateGas
  Future<BigInt> estimateGas({
    required BlockChainType blockChainType,
    required Transaction transaction,
  }) async {
    final client = web3dart.getCurrentClient(blockChainType: blockChainType);

    try {
      final estimateGas = (await client.estimateGas(
        sender: transaction.from ??
            (isCentralized()
                ? EthereumAddress.fromHex(common.ownerWalletAddress ?? '')
                : myEthereumAddress),
        to: transaction.to,
        data: transaction.data,
        value: transaction.value,
        gasPrice: web3dart.getGasPrice(blockChainType: blockChainType),
      ));
      print("estimateGas:$estimateGas");
      return estimateGas;
    } catch (e) {
      print("estimateGas:${e.toString()}");
      return BigInt.from(1000000);
    }
  }

  /// getTransactionFromCallContract
  Transaction getTransactionFromCallContract({
    required DeployedContract deployedContract,
    required BigInt amount,
    required String toAddress,
    required BlockChainType blockChainType,
    int? maxGas,
  }) {
    ContractFunction transferEvent = deployedContract.function('transfer');
    EthereumAddress? toETHAddress = EthereumAddress.fromHex(toAddress);
    Transaction transaction = Transaction.callContract(
      maxGas: maxGas,
      gasPrice: getGasPrice(blockChainType: blockChainType),
      contract: deployedContract,
      function: transferEvent,
      parameters: [toETHAddress, amount],
    );
    return transaction;
  }

  /// getEurusInternalConfig
  DeployedContract getEurusInternalConfig() {
    final EthereumAddress contractAddr =
        EthereumAddress.fromHex(eurusInternalConfigAddress);
    print('internalConfigAddress:$eurusInternalConfigAddress');
    String abiCode =
        '''[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"","type":"string"}],"name":"Event","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"}],"name":"OwnerAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"}],"name":"OwnerRemoved","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"addOwner","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"addressList","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true},{"inputs":[],"name":"eurusUserDepositAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true},{"inputs":[],"name":"getOwnerCount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function","constant":true},{"inputs":[],"name":"getOwners","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function","constant":true},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"isOwner","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function","constant":true},{"inputs":[],"name":"platformWalletAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"removeOwner","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_currencyAddr","type":"address"},{"internalType":"string","name":"asset","type":"string"}],"name":"addCurrencyInfo","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"asset","type":"string"}],"name":"removeCurrencyInfo","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"asset","type":"string"}],"name":"getErc20SmartContractAddrByAssetName","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function","constant":true},{"inputs":[{"internalType":"address","name":"_currencyAddr","type":"address"}],"name":"getErc20SmartContractByAddr","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function","constant":true},{"inputs":[{"internalType":"address","name":"coldWalletAddr","type":"address"}],"name":"setPlatformWalletAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"userDepositAddr","type":"address"}],"name":"setEurusUserDepositAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getAssetAddress","outputs":[{"internalType":"string[]","name":"","type":"string[]"},{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function","constant":true}]''';
    DeployedContract contract = DeployedContract(
        ContractAbi.fromJson(abiCode, 'EurusInternalConfig'), contractAddr);
    return contract;
  }

  /// getExternalSmartContractConfig
  DeployedContract getExternalSmartContractConfig() {
    final EthereumAddress contractAddr =
        EthereumAddress.fromHex(externalSmartContractConfigAddress);
    print(
        'externalSmartContractConfigAddress:$externalSmartContractConfigAddress');
    String abiCode = '''[{
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "name": "Event",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "OwnerAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "OwnerRemoved",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferred",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "reader",
          "type": "address"
        }
      ],
      "name": "ReaderAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "reader",
          "type": "address"
        }
      ],
      "name": "ReaderRemoved",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "writer",
          "type": "address"
        }
      ],
      "name": "WriterAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "writer",
          "type": "address"
        }
      ],
      "name": "WriterRemoved",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "addOwner",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newReader",
          "type": "address"
        }
      ],
      "name": "addReader",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newWriter",
          "type": "address"
        }
      ],
      "name": "addWriter",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "addressList",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "currencyList",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getOwnerCount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getOwners",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getReaderList",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getWriterList",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "isOwner",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "isWriter",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "removeOwner",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "existingReader",
          "type": "address"
        }
      ],
      "name": "removerReader",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "existingWriter",
          "type": "address"
        }
      ],
      "name": "removerWriter",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "renounceOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_currencyAddr",
          "type": "address"
        },
        {
          "internalType": "string",
          "name": "asset",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "decimal",
          "type": "uint256"
        },
        {
          "internalType": "string",
          "name": "id",
          "type": "string"
        }
      ],
      "name": "addCurrencyInfo",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "asset",
          "type": "string"
        }
      ],
      "name": "removeCurrencyInfo",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "asset",
          "type": "string"
        }
      ],
      "name": "getErc20SmartContractAddrByAssetName",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_currencyAddr",
          "type": "address"
        }
      ],
      "name": "getErc20SmartContractByAddr",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getAssetAddress",
      "outputs": [
        {
          "internalType": "string[]",
          "name": "",
          "type": "string[]"
        },
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "asset",
          "type": "string"
        }
      ],
      "name": "getAssetDecimal",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "asset",
          "type": "string"
        }
      ],
      "name": "getAssetListID",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "ethFee",
          "type": "uint256"
        },
        {
          "internalType": "string[]",
          "name": "asset",
          "type": "string[]"
        },
        {
          "internalType": "uint256[]",
          "name": "amount",
          "type": "uint256[]"
        }
      ],
      "name": "setETHFee",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "asset",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "setAdminFee",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "asset",
          "type": "string"
        }
      ],
      "name": "getAdminFee",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "asset",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "kycLevel",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "limit",
          "type": "uint256"
        }
      ],
      "name": "setKycLimit",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "symbol",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "kycLevel",
          "type": "string"
        }
      ],
      "name": "getCurrencyKycLimit",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "gasPriceWei",
          "type": "uint256"
        }
      ],
      "name": "setEurusGasPrice",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getEurusGasPrice",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "gasAmount",
          "type": "uint256"
        }
      ],
      "name": "setMaxTopUpGasAmount",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getMaxTopUpGasAmount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }]''';
    DeployedContract contract = DeployedContract(
        ContractAbi.fromJson(abiCode, 'ExternalSmartContractConfig'),
        contractAddr);
    return contract;
  }

  /// getEthereumERC20Contract
  DeployedContract getEthereumERC20Contract({
    required String contractAddress,
  }) {
    final EthereumAddress contractAddr =
        EthereumAddress.fromHex(contractAddress);
    String abiCode =
        '''[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"tokens","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokens","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"tokenOwner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"acceptOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"drip","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"tokens","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"tokens","type":"uint256"},{"name":"data","type":"bytes"}],"name":"approveAndCall","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"newOwner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"tokenAddress","type":"address"},{"name":"tokens","type":"uint256"}],"name":"transferAnyERC20Token","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"tokenOwner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"tokens","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"tokenOwner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"tokens","type":"uint256"}],"name":"Approval","type":"event"}]''';
    DeployedContract contract = DeployedContract(
        ContractAbi.fromJson(abiCode, 'TestingCoin'), contractAddr);
    return contract;
  }

  /// getEurusERC20Contract
  DeployedContract getEurusERC20Contract({
    required String contractAddress,
  }) {
    final EthereumAddress contractAddr =
        EthereumAddress.fromHex(contractAddress);
    String abiCode = '''[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Approval",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "name": "Event",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "OwnerAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "OwnerRemoved",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferred",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "reader",
          "type": "address"
        }
      ],
      "name": "ReaderAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "reader",
          "type": "address"
        }
      ],
      "name": "ReaderRemoved",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "writer",
          "type": "address"
        }
      ],
      "name": "WriterAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "writer",
          "type": "address"
        }
      ],
      "name": "WriterRemoved",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "addOwner",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newReader",
          "type": "address"
        }
      ],
      "name": "addReader",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newWriter",
          "type": "address"
        }
      ],
      "name": "addWriter",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        }
      ],
      "name": "allowance",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "blackListDestAddress",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "blackListDestAddressMap",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "decimals",
      "outputs": [
        {
          "internalType": "uint8",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "subtractedValue",
          "type": "uint256"
        }
      ],
      "name": "decreaseAllowance",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getOwnerCount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getOwners",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getReaderList",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getWriterList",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "addedValue",
          "type": "uint256"
        }
      ],
      "name": "increaseAllowance",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "isOwner",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "isWriter",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "name",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "removeOwner",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "existingReader",
          "type": "address"
        }
      ],
      "name": "removerReader",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "existingWriter",
          "type": "address"
        }
      ],
      "name": "removerWriter",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "renounceOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "symbol",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "totalSupply",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "recipient",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "transfer",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "sender",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "recipient",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "transferFrom",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "name_",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "symbol_",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "totalSupply_",
          "type": "uint256"
        },
        {
          "internalType": "uint8",
          "name": "decimals_",
          "type": "uint8"
        }
      ],
      "name": "init",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "internalSCAddr",
          "type": "address"
        },
        {
          "internalType": "string",
          "name": "name_",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "symbol_",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "totalSupply_",
          "type": "uint256"
        },
        {
          "internalType": "uint8",
          "name": "decimals_",
          "type": "uint8"
        }
      ],
      "name": "init",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "mint",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "burn",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "dest",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "withdrawAmount",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "amountWithFee",
          "type": "uint256"
        }
      ],
      "name": "submitWithdraw",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "setInternalSCConfigAddress",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getInternalSCConfigAddress",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "addBlackListDestAddress",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "removeBlackListDestAddress",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]''';
    DeployedContract contract = DeployedContract(
        ContractAbi.fromJson(abiCode, 'EurusERC20'), contractAddr);
    return contract;
  }

  Future<DeployedContract> getBSCTokenContract({
    required String contractAddress,
  }) async {
    final EthereumAddress contractAddr =
        EthereumAddress.fromHex(contractAddress);
    final result = await apiHandler
        .get(
            "https://api-testnet.bscscan.com/api?module=contract&action=getabi&address=$contractAddress&apikey=YourApiKeyToken")
        .then((value) => value)
        .catchError((e) => e);
    String abiCode = result['result'];
    if (result['status'] == "0") {
      abiCode =
          '''[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"tokens","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokens","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"tokenOwner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"acceptOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"drip","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"tokens","type":"uint256"}],"name":"transfer","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"tokens","type":"uint256"},{"name":"data","type":"bytes"}],"name":"approveAndCall","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"newOwner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"tokenAddress","type":"address"},{"name":"tokens","type":"uint256"}],"name":"transferAnyERC20Token","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"tokenOwner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"tokens","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"tokenOwner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"tokens","type":"uint256"}],"name":"Approval","type":"event"}]''';
    }
    DeployedContract contract = DeployedContract(
        ContractAbi.fromJson(abiCode, 'BSCContract'), contractAddr);
    return contract;
  }

  /// getCurrentClient
  Web3Client getCurrentClient({
    required BlockChainType blockChainType,
  }) {
    Map<BlockChainType, Web3Client> list = {
      BlockChainType.Eurus: eurusEthClient,
      BlockChainType.Ethereum: mainNetEthClient,
      BlockChainType.BinanceCoin: mainNetBscClient,
    };
    return list[blockChainType] ?? eurusEthClient;
  }

  double getBlockChainGasPrice(BlockChainType blockChainType) {
    Map<BlockChainType, double> list = {
      BlockChainType.Eurus: eurusGasPrice,
      BlockChainType.Ethereum: ethereumGasPrice,
      BlockChainType.BinanceCoin: bscGasPrice,
    };
    return list[blockChainType] ?? eurusGasPrice;
  }

  DeployedContract getERC20Contract(BlockChainType blockChainType) {
    Map<BlockChainType, DeployedContract> list = {
      BlockChainType.Eurus: erc20ContractFromEurus,
      BlockChainType.Ethereum: erc20ContractFromEthereum,
      BlockChainType.BinanceCoin: erc20ContractBSC,
    };
    return list[blockChainType] ?? erc20ContractFromEurus;
  }

  String getRpcUrl(BlockChainType blockChainType) {
    Map<BlockChainType, String> list = {
      BlockChainType.Eurus: eurusRPCUrl,
      BlockChainType.Ethereum: blockchainRpcUrl,
      BlockChainType.BinanceCoin: bscRpcUrl,
    };
    return list[blockChainType] ?? blockchainRpcUrl;
  }

  int getChainId(BlockChainType blockChainType) {
    Map<BlockChainType, int> list = {
      BlockChainType.Eurus: eurusChainId,
      BlockChainType.Ethereum: mainNetChainId,
      BlockChainType.BinanceCoin: bscChainId,
    };
    return list[blockChainType] ?? mainNetChainId;
  }

  /// get getETHBalance
  Future<String> getETHBalance({
    required BlockChainType blockChainType,
  }) async {
    if (myEthereumAddress == null) return '';

    Web3Client client = getCurrentClient(blockChainType: blockChainType);
    EtherAmount balance = await client.getBalance(myEthereumAddress!);
    int decimals = 18;
    double balanceInEther = balance.getInWei.toDouble() / pow(10, decimals);
    print("getETHBalance: $balanceInEther");
    String balanceString = balanceInEther.toString();
    int decimalsDifference = decimals -
        balanceString.substring(balanceString.indexOf('.') + 1).length;
    if (decimalsDifference > 0) {
      balanceString = balanceString.padRight(
          balanceString.length + decimalsDifference, '0');
    }

    String trimmedBalanceString = balanceString.replaceRange(
        balanceString.indexOf('.') + (decimals >= 8 ? 8 : decimals) + 1,
        null,
        '');
    return trimmedBalanceString;
  }

  /// get getERC20Balance
  Future<String?> getERC20Balance({
    required DeployedContract deployedContract,
    required BlockChainType blockChainType,
  }) async {
    Web3Client client = getCurrentClient(blockChainType: blockChainType);
    ContractFunction getBalance = deployedContract.function('balanceOf');
    List balanceList = await client.call(
        contract: deployedContract,
        function: getBalance,
        params: [myEthereumAddress]);
    if (balanceList.isEmpty) return null;
    double? balanceInDecimal = (balanceList.first as BigInt?)?.toDouble();
    int? decimals = (await getContractDecimal(
            deployedContract: deployedContract, blockChainType: blockChainType))
        ?.toInt();
    if (balanceInDecimal == null || decimals == null) return null;
    double balance = balanceInDecimal / pow(10, decimals);
    print("getERC20Balance: $balance");
    String balanceString = balance.toString();
    int decimalsDifference = decimals -
        balanceString.substring(balanceString.indexOf('.') + 1).length;
    if (decimalsDifference > 0) {
      balanceString = balanceString.padRight(
          balanceString.length + decimalsDifference, '0');
    }

    String trimmedBalanceString = balanceString.replaceRange(
        balanceString.indexOf('.') + (decimals >= 8 ? 8 : decimals) + 1,
        null,
        '');
    return trimmedBalanceString;
  }

  /// get getETHBalance
  Future<List<dynamic>> getERC20TokenList({
    required BlockChainType blockChainType,
  }) async {
    Web3Client client = getCurrentClient(blockChainType: blockChainType);
    DeployedContract deployedContract;
    if (blockChainType == BlockChainType.Eurus) {
      deployedContract = getExternalSmartContractConfig();
    } else {
      deployedContract = getEurusInternalConfig();
    }
    ContractFunction getAssetAddress =
        deployedContract.function('getAssetAddress');
    print("0xdE9c12961680811aa7d068EB727ef4017BA94929");
    tokenList = await client.call(
        contract: deployedContract, function: getAssetAddress, params: []);
    print("tokenList:$tokenList");
    tokenListMap = new Map<String, String>();
    if (tokenList.isNotEmpty && tokenList[0] != null) {
      for (var i = 0; i < tokenList[0].length; i++) {
        String tokenName = tokenList[0][i];
        EthereumAddress tokenAddress = tokenList[1][i];
        tokenListMap[tokenName] = tokenAddress.toString();
      }
    }
    print('tokenListMap$tokenListMap');
    return tokenList;
  }

  /// getContractDecimal
  Future<BigInt?> getContractDecimal({
    required DeployedContract deployedContract,
    required BlockChainType blockChainType,
  }) async {
    Web3Client client = getCurrentClient(blockChainType: blockChainType);
    ContractFunction getDecimals = deployedContract.function('decimals');
    List decimalsNumber = await client
        .call(contract: deployedContract, function: getDecimals, params: []);

    if (decimalsNumber.isEmpty) return null;
    BigInt decimalsBalance = decimalsNumber.first;
    print("decimalsBalance$decimalsBalance");
    return decimalsBalance;
  }

  /// getEurusUserDepositAddress
  Future<String?> getEurusUserDepositAddress() async {
    Web3Client client =
        getCurrentClient(blockChainType: BlockChainType.Ethereum);
    DeployedContract deployedContract = getEurusInternalConfig();
    ContractFunction eurusUserDepositAddress =
        deployedContract.function('eurusUserDepositAddress');
    var address = await client.call(
        contract: deployedContract,
        function: eurusUserDepositAddress,
        params: []);
    print('getEurusUserDepositAddress$address');

    if (address.isEmpty) return null;
    EthereumAddress ethereumAddress = address.first;
    return ethereumAddress.toString();
  }

  /// getTokenSymbol
  Future<String?> getTokenSymbol({
    required DeployedContract deployedContract,
    required BlockChainType blockChainType,
  }) async {
    Web3Client client = getCurrentClient(blockChainType: blockChainType);
    ContractFunction getDecimals = deployedContract.function('symbol');
    List name = await client
        .call(contract: deployedContract, function: getDecimals, params: []);
    if (name.isEmpty) return null;
    String tokenName = name.first;
    print("tokenName:$tokenName");
    return tokenName;
  }

  EtherAmount getGasPrice({
    required BlockChainType blockChainType,
  }) {
    return EtherAmount.inWei(
        BigInt.from(getBlockChainGasPrice(blockChainType) * transactionSpeed));
  }

  /// sendETH
  Future<Transaction?> sendETH({
    required double enterAmount,
    required String toAddress,
    required BlockChainType type,
    int? maxGas,
  }) async {
    BigInt amount = BigInt.from(1000000000000000000 * enterAmount);
    EthereumAddress toETHAddress = EthereumAddress.fromHex(toAddress);
    return Transaction(
      to: toETHAddress,
      maxGas: maxGas,
      gasPrice: getGasPrice(blockChainType: type),
      value: EtherAmount.inWei(amount),
    );
  }

  /// sendERC20
  Future<Transaction?> sendERC20({
    required DeployedContract deployedContract,
    required double enterAmount,
    required String toAddress,
    required BlockChainType blockChainType,
    int? maxGas,
  }) async {
    BigInt? decimalsBalance = await getContractDecimal(
        deployedContract: deployedContract, blockChainType: blockChainType);
    if (decimalsBalance == null) return null;
    String decimalsString = "1".padRight(decimalsBalance.toInt() + 1, "0");
    BigInt amount = BigInt.from(double.parse(decimalsString) * enterAmount);
    print("BigIntamount:$amount");
    return getTransactionFromCallContract(
      deployedContract: deployedContract,
      maxGas: maxGas,
      amount: amount,
      toAddress: toAddress,
      blockChainType: blockChainType,
    );
  }

  /// submitWithdrawERC20
  Future<Transaction?> submitWithdrawERC20({
    required DeployedContract deployedContract,
    required double enterAmount,
    required double enterAmountWithFee,
    required String toAddress,
    int? maxGas,
  }) async {
    BigInt? decimalsBalance = await getContractDecimal(
        deployedContract: deployedContract,
        blockChainType: BlockChainType.Eurus);
    if (decimalsBalance == null) return null;
    String decimalsString = "1".padRight(decimalsBalance.toInt() + 1, "0");
    BigInt amount = BigInt.from(double.parse(decimalsString) * enterAmount);
    BigInt amountWithFee =
        BigInt.from(double.parse(decimalsString) * enterAmountWithFee);
    print("BigIntamount:$amount");
    ContractFunction transferEvent =
        deployedContract.function('submitWithdraw');
    EthereumAddress toETHAddress = EthereumAddress.fromHex(toAddress);
    return Transaction.callContract(
      maxGas: maxGas,
      gasPrice: getGasPrice(blockChainType: BlockChainType.Eurus),
      contract: deployedContract,
      function: transferEvent,
      parameters: [toETHAddress, amount, amountWithFee],
    );
  }

  /// BroadCastTranscation
  Future<String?> broadCastTranscation({
    required Credentials credentials,
    required Transaction transaction,
    required BlockChainType blockChainType,
  }) async {
    String transactionResult;
    Web3Client client = getCurrentClient(blockChainType: blockChainType);
    transactionResult = await client.sendTransaction(
      credentials,
      transaction,
      chainId: getChainId(blockChainType),
    );
    print("sendERC20 result:$transactionResult");
    return transactionResult;
  }

  /// InvokeSmartContract
  Future<String?> invokeSmartContract({
    required Credentials credentials,
    required EthereumAddress? scAddr,
    required BigInt? eun,
    required Uint8List? data,
    required BlockChainType blockChainType,
  }) async {
    String transactionResult;
    Web3Client client = getCurrentClient(blockChainType: blockChainType);
    DeployedContract deployedContract =
        getUserWallet(contractAddrString: common.cenUserWalletAddress ?? '');
    ContractFunction invokeSmartContractFunction =
        deployedContract.function('invokeSmartContract');
    Transaction transaction = Transaction.callContract(
      gasPrice: getGasPrice(blockChainType: blockChainType),
      contract: deployedContract,
      function: invokeSmartContractFunction,
      parameters: [
        scAddr,
        eun,
        data,
      ],
    );
    transactionResult = await client.sendTransaction(
      credentials,
      transaction,
      chainId: getChainId(blockChainType),
    );
    print("sendERC20 result:$transactionResult");
    return transactionResult;
  }

  /// getETHClientDetail
  Future<Web3Client> getETHClientDetail({
    required BlockChainType blockChainType,
  }) async {
    Web3Client client = getCurrentClient(blockChainType: blockChainType);
    print("---------------------- getETHClientDetail ----------------------");
    print("getClientVersion:${await client.getClientVersion()}");
    print("getBlockNumber:${await client.getBlockNumber()}");
    print("getGasPrice:${await client.getGasPrice()}");
    print("getEtherProtocolVersion:${await client.getEtherProtocolVersion()}");
    print("getMiningHashrate:${await client.getMiningHashrate()}");
    print("getNetworkId:${await client.getNetworkId()}");
    print("getPeerCount:${await client.getPeerCount()}");

    return client;
  }

  /// getAddressDetail
  Future<Web3Client> getAddressDetail({
    required BlockChainType blockChainType,
  }) async {
    Web3Client client = getCurrentClient(blockChainType: blockChainType);
    print("---------------------- getAddressDetail ----------------------");
    print("getBalance:${await client.getBalance(
      EthereumAddress.fromHex('0x44f426bc9ac7a83521EA140Aeb70523C0a85945a'),
    )}");
    print("etTransactionCount:${await client.getTransactionCount(
      EthereumAddress.fromHex('0x44f426bc9ac7a83521EA140Aeb70523C0a85945a'),
    )}");
    TransactionReceipt? transactionReceipt = await client.getTransactionReceipt(
        "0xfa0a7ed6a87b655f2302ce2d88d1d051c4eeef2af6e82de9850f3527a8106744");
    print("---------------------- hash data ----------------------");
    print(
        "transactionReceipt.contractAddress:${transactionReceipt?.contractAddress}");
    print("transactionReceipt.gasUsed:${transactionReceipt?.gasUsed}");
    print("transactionReceipt.from:${transactionReceipt?.from}");
    print("transactionReceipt.to:${transactionReceipt?.to}");
    return client;
  }

  /// initNewWallet
  void initNewWallet() async {
    var rng = new Random.secure();
    EthPrivateKey random = EthPrivateKey.createRandom(rng);

    var address = await random.extractAddress();
    print("extract address: ${address.hex}");

    Wallet wallet = Wallet.createNew(
      random,
      'password',
      Random(),
      scryptN: pow(2, 8) as int,
    );
    print("wallet json ${wallet.toJson()}");
  }

  ///cen part
  /// transferRequest
  // Future<String> transferRequest(
  //     {String selectTokenSymbol,String userWalletAddress,DeployedContract deployedContract,
  //       double enterAmount,
  //       String toAddress,
  //       BlockChainType blockChainType}) async {
  //   String transactionResult;
  //   Web3Client client = getCurrentClient(blockChainType:blockChainType);
  //   BigInt decimalsBalance = await getContractDecimal(deployedContract: deployedContract,blockChainType: blockChainType);
  //   String decimalsString = "1".padRight(decimalsBalance.toInt()+1,"0");
  //   decimalsString = (selectTokenSymbol == "EUN") && common.currentUserType == CurrentUserType.centralized ? "1000000000000000000" : decimalsString;
  //   BigInt amount = BigInt.from(double.parse(decimalsString) * enterAmount);
  //   print("BigIntamount:$amount");
  //
  //   DeployedContract userWalletContract = getUserWallet(contractAddrString:userWalletAddress);
  //     ContractFunction transferEvent = userWalletContract.function('transferRequest');
  //     EthereumAddress toETHAddress = EthereumAddress.fromHex(toAddress);
  //     Transaction transaction = Transaction.callContract(
  //       gasPrice: getGasPrice(blockChainType: blockChainType),
  //       contract: userWalletContract,
  //       function: transferEvent,
  //       parameters: [toETHAddress,selectTokenSymbol,amount],
  //     );
  //
  //   transactionResult = await client.sendTransaction(
  //       await eurusEthClient.credentialsFromPrivateKey(common.cenSignKey),
  //       transaction,
  //       chainId: chainId);
  //   print("sendERC20 result:$transactionResult");
  //   return transactionResult;
  // }

  Future<Transaction> requestTransfer({
    required String selectTokenSymbol,
    required String userWalletAddress,
    required DeployedContract deployedContract,
    required double enterAmount,
    required String toAddress,
    required BlockChainType blockChainType,
    int? maxGas,
  }) async {
    BigInt decimalsBalance;
    String decimalsString;
    if ((selectTokenSymbol == "EUN") && isCentralized()) {
      decimalsString = "1000000000000000000";
    } else {
      decimalsBalance = await getContractDecimal(
              deployedContract: deployedContract,
              blockChainType: blockChainType) ??
          BigInt.from(0);
      decimalsString = "1".padRight(decimalsBalance.toInt() + 1, "0");
    }
    BigInt amount = BigInt.from(double.parse(decimalsString) * enterAmount);
    print("BigIntamount:$amount");

    DeployedContract userWalletContract =
        getUserWallet(contractAddrString: userWalletAddress);
    ContractFunction transferEvent =
        userWalletContract.function('requestTransferV1');
    EthereumAddress toETHAddress = EthereumAddress.fromHex(toAddress);

    Transaction transaction = Transaction.callContract(
      gasPrice: getGasPrice(blockChainType: blockChainType),
      contract: userWalletContract,
      function: transferEvent,
      maxGas: maxGas,
      parameters: [toETHAddress, selectTokenSymbol, amount],
    );

    return transaction;
  }

  EthereumAddress personalEcRecover(Uint8List message, Uint8List signature) {
    final prefix = Uint8List.fromList(utf8.encode(
        '\u0019Ethereum Signed Message:\n${message.length.toString()}'));

    final r = bytesToUnsignedInt(signature.sublist(0, 32));
    final s = bytesToUnsignedInt(signature.sublist(32, 64));
    final v = signature.elementAt(64);
    final signatureData = MsgSignature(r, s, v);

    return EthereumAddress.fromHex(
      bytesToHex(
        publicKeyToAddress(ecRecover(
          keccak256(Uint8List.fromList(prefix + message)),
          signatureData,
        )),
        include0x: true,
      ),
    );
  }

  Uint8List parameterSignature({
    required String functionName,
    required DeployedContract userWalletContract,
    required List<dynamic> params,
    List<String>? paramsName,
  }) {
    if (common.cenSignKey == null) return Uint8List(0);

    ContractFunction transferFunction =
        userWalletContract.function(functionName);
    Uint8List transferEventData = transferFunction
        .encodeCallWithKeccak256(params, paramsName: paramsName);
    print("transferEventData:$transferEventData");
    final Uint8List _privateKeyBytes =
        _web3dartCrypto.hexToBytes(common.cenSignKey!);
    //print("common.cenSignKey:${common.cenSignKey}");
    final _web3dartCrypto.MsgSignature _msgSignature =
        _web3dartCrypto.sign(transferEventData, _privateKeyBytes);
    final Uint8List _rBytes =
        padUint8ListTo32(_web3dartCrypto.unsignedIntToBytes(_msgSignature.r));
    final Uint8List _sBytes =
        padUint8ListTo32(_web3dartCrypto.unsignedIntToBytes(_msgSignature.s));
    final Uint8List _vBytes =
        _web3dartCrypto.unsignedIntToBytes(BigInt.from(_msgSignature.v));
    final Uint8List _concatenatedInto65Bytes =
        uint8ListFromList(_rBytes + _sBytes + _vBytes);
    print("_concatenatedInto65Bytes:$_concatenatedInto65Bytes");
    return _concatenatedInto65Bytes;
  }

  /// transferRequest
  Future<Transaction> cenSubmitWithdraw({
    required String selectTokenSymbol,
    required String userWalletAddress,
    required DeployedContract deployedContract,
    required double enterAmount,
    required String toAddress,
    required double enterAmountWithFee,
    int? maxGas,
  }) async {
    BigInt? decimalsBalance = await getContractDecimal(
        deployedContract: deployedContract,
        blockChainType: BlockChainType.Eurus);
    // if (decimalsBalance == null) return null;
    String decimalsString = "1".padRight(decimalsBalance!.toInt() + 1, "0");
    BigInt amount = BigInt.from(double.parse(decimalsString) * enterAmount);
    BigInt amountWithFee =
        BigInt.from(double.parse(decimalsString) * enterAmountWithFee);
    print("BigIntamount:$amount");

    DeployedContract userWalletContract =
        getUserWallet(contractAddrString: userWalletAddress);
    ContractFunction transferEvent =
        userWalletContract.function('submitWithdrawV1');
    EthereumAddress toETHAddress = EthereumAddress.fromHex(toAddress);

    // Uint8List signature = parameterSignature(
    //     functionName: 'directSubmitWithdraw',
    //     userWalletContract: userWalletContract,
    //     params: [toETHAddress, amount, amountWithFee, selectTokenSymbol]);
    Transaction transaction = Transaction.callContract(
      maxGas: maxGas,
      gasPrice: getGasPrice(blockChainType: BlockChainType.Eurus),
      contract: userWalletContract,
      function: transferEvent,
      parameters: [
        toETHAddress,
        amount,
        amountWithFee,
        selectTokenSymbol,
        // signature
      ],
    );

    return transaction;
  }

  DeployedContract getUserWallet({
    required String contractAddrString,
  }) {
    final EthereumAddress contractAddr =
        EthereumAddress.fromHex(contractAddrString);
    String abiCode = '''[
    {
      "inputs": [],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "sender",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "Confirmation",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "sender",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Deposit",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "name": "Event",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "Execution",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "ExecutionFailure",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "OwnerAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "OwnerRemoved",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferred",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "reader",
          "type": "address"
        }
      ],
      "name": "ReaderAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "reader",
          "type": "address"
        }
      ],
      "name": "ReaderRemoved",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "sender",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "Rejection",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "required",
          "type": "uint256"
        }
      ],
      "name": "RequirementChange",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "sender",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "Revocation",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "Submission",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "dest",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "targetGasWei",
          "type": "uint256"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "gasTransferred",
          "type": "uint256"
        }
      ],
      "name": "TopUpPaymentWalletEvent",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "dest",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "string",
          "name": "assetName",
          "type": "string"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "TransferEvent",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "dest",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "string",
          "name": "assetName",
          "type": "string"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "TransferRequestEvent",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "dest",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "withdrawAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "string",
          "name": "assetName",
          "type": "string"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amountWithFee",
          "type": "uint256"
        }
      ],
      "name": "WithdrawRequestEvent",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "writer",
          "type": "address"
        }
      ],
      "name": "WriterAdded",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "writer",
          "type": "address"
        }
      ],
      "name": "WriterRemoved",
      "type": "event"
    },
    {
      "inputs": [],
      "name": "MAX_OWNER_COUNT",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "TranList",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "addOwner",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newReader",
          "type": "address"
        }
      ],
      "name": "addReader",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "operatorAddr",
          "type": "address"
        }
      ],
      "name": "addWalletOperator",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "newWriter",
          "type": "address"
        }
      ],
      "name": "addWriter",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "_required",
          "type": "uint256"
        }
      ],
      "name": "changeRequirement",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "confirmTransaction",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "confirmations",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "getConfirmationCount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "count",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "getConfirmations",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "_confirmations",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getOwnerCount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getOwners",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getReaderList",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bool",
          "name": "pending",
          "type": "bool"
        },
        {
          "internalType": "bool",
          "name": "executed",
          "type": "bool"
        }
      ],
      "name": "getTransactionCount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "count",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "from",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "to",
          "type": "uint256"
        },
        {
          "internalType": "bool",
          "name": "pending",
          "type": "bool"
        },
        {
          "internalType": "bool",
          "name": "executed",
          "type": "bool"
        }
      ],
      "name": "getTransactionIds",
      "outputs": [
        {
          "internalType": "uint256[]",
          "name": "_transactionIds",
          "type": "uint256[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getWalletOperatorList",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getWalletOwner",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getWriterList",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "isConfirmed",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "isOwner",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "transId",
          "type": "uint256"
        }
      ],
      "name": "isRejected",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "isWriter",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        },
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "name": "miscellaneousData",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "transId",
          "type": "uint256"
        }
      ],
      "name": "rejectTransaction",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "removeOwner",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "operatorAddr",
          "type": "address"
        }
      ],
      "name": "removeWalletOperator",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "existingReader",
          "type": "address"
        }
      ],
      "name": "removerReader",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "existingWriter",
          "type": "address"
        }
      ],
      "name": "removerWriter",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "renounceOwnership",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "required",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "transactionId",
          "type": "uint256"
        }
      ],
      "name": "revokeConfirmation",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "setWalletOwner",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "a",
          "type": "address"
        }
      ],
      "name": "toBytes",
      "outputs": [
        {
          "internalType": "bytes",
          "name": "",
          "type": "bytes"
        }
      ],
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "transactionCount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "transactions",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "transId",
          "type": "uint256"
        },
        {
          "internalType": "bool",
          "name": "isDirectInvokeData",
          "type": "bool"
        },
        {
          "internalType": "address",
          "name": "destination",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "data",
          "type": "bytes"
        },
        {
          "internalType": "uint256",
          "name": "timestamp",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "blockNumber",
          "type": "uint256"
        },
        {
          "internalType": "bool",
          "name": "executed",
          "type": "bool"
        },
        {
          "internalType": "bool",
          "name": "rejected",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "walletOperatorList",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "walletOwner",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "stateMutability": "payable",
      "type": "receive"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "addr",
          "type": "address"
        }
      ],
      "name": "setInternalSmartContractConfig",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "dest",
          "type": "address"
        },
        {
          "internalType": "string",
          "name": "assetName",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "requestTransferV1",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getGasFeeWalletAddress",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "hash",
          "type": "bytes32"
        },
        {
          "internalType": "bytes",
          "name": "signature",
          "type": "bytes"
        }
      ],
      "name": "verifySignature",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "dest",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "withdrawAmount",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "amountWithFee",
          "type": "uint256"
        },
        {
          "internalType": "string",
          "name": "assetName",
          "type": "string"
        }
      ],
      "name": "submitWithdrawV1",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "targetGasWei",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "gasLimit",
          "type": "uint256"
        }
      ],
      "name": "directTopUpPaymentWallet",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "paymentWalletAddr",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "targetGasWei",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "signature",
          "type": "bytes"
        }
      ],
      "name": "topUpPaymentWallet",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getWalletOwnerBalance",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "scAddr",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "eun",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "inputArg",
          "type": "bytes"
        }
      ],
      "name": "invokeSmartContract",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]''';
    DeployedContract contract = DeployedContract(
        ContractAbi.fromJson(abiCode, 'UserWallet'), contractAddr);
    return contract;
  }

  Future<String> sendRawTransaction({
    required String signedTransactionString,
  }) async {
    signedTransactionString = signedTransactionString.replaceAll("0x", '');
    print("signedTransactionString:$signedTransactionString");
    var signedTransaction = hex.decode(signedTransactionString);
    Web3Client client = getCurrentClient(blockChainType: BlockChainType.Eurus);
    String lastTxId =
        await client.sendRawTransaction(Uint8List.fromList(signedTransaction));
    print("sendRawTransactionlastTxId:$lastTxId");
    return lastTxId;
  }

  Future<double?> getMaxTopUpGasAmount() async {
    DeployedContract deployedContract = getExternalSmartContractConfig();
    ContractFunction getMaxTopUpGasAmount =
        deployedContract.function('getMaxTopUpGasAmount');
    final response = await eurusEthClient.call(
        contract: deployedContract, function: getMaxTopUpGasAmount, params: []);
    if (response.isEmpty) return null;
    BigInt maxGasAmount = response.first;
    return maxGasAmount.toDouble();
  }
}

/// you can use web3dart
Web3dart web3dart = Web3dart();
