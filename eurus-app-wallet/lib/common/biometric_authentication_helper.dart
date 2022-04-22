import 'dart:io';

import 'package:biometric_storage/biometric_storage.dart';
import 'package:flutter/services.dart';

class BiometricAuthenticationHelper {
  Future<String?> canPersistWithBiometricSecurely(dynamic args,
      {bool delete = false}) async {
    if (args is Map<String, String>) {
      List<bool> errorResult = [];
      await Future.forEach(
          args.keys,
          (k) async => await (await _getBiometricStorageFile(
                      k, await canSupportsBiometricAuthenticated()))
                  ?.write(args[k]!)
                  .catchError(
                (error, stackTrace) => null,
                test: (error) {
                  errorResult.add(false);
                  return false;
                },
              ));
      return errorResult.contains(false) ? null : args.values.join(',');
    } else {
      if (!delete) {
        return await (await _getBiometricStorageFile(
                args, await canSupportsBiometricAuthenticated()))
            ?.read();
      } else {
        await (await _getBiometricStorageFile(
                args, await canSupportsBiometricAuthenticated()))
            ?.delete()
            .catchError((error, stackTrace) => null,
                test: (error) => /* // errSecItemNotFound */ error
                        is PlatformException &&
                    error.code.startsWith('SecurityError') &&
                    error.message ==
                        "Error while ${'writing data'}: ${-25300}: ${'The specified item could not be found in the keychain.'}");
      }
    }
  }

  Future<String?> canPersistTxWithBiometricSecurely(dynamic args,
      {bool delete = false}) async {
    if (args is Map<String, String>) {
      await Future.forEach<String>(
          args.keys,
          (k) async => await (await _getTxBiometricStorageFile(
                  k, await canSupportsBiometricAuthenticated()))
              ?.write(args[k]!));
    } else {
      if (!delete) {
        return await (await _getTxBiometricStorageFile(
                args, await canSupportsBiometricAuthenticated()))
            ?.read();
      } else {
        await (await _getTxBiometricStorageFile(
                args, await canSupportsBiometricAuthenticated()))
            ?.delete()
            .catchError((error, stackTrace) => null,
                test: (error) => /* // errSecItemNotFound */ error
                        is PlatformException &&
                    error.code.startsWith('SecurityError') &&
                    error.message ==
                        "Error while ${'writing data'}: ${-25300}: ${'The specified item could not be found in the keychain.'}");
      }
    }
  }

  Future<bool> canSupportsBiometricAuthenticated() async {
    final authenticate = await _checkBiometricCapabilitiesCanAuthenticate();
    bool supportsBiometricAuthenticated = false;
    if (authenticate == CanAuthenticateResponse.success) {
      supportsBiometricAuthenticated = true;
    } else if (authenticate != CanAuthenticateResponse.unsupported) {
      supportsBiometricAuthenticated = false;
    } else {
      print('Unable to use authenticate. Unable to get storage.');
      return false;
    }
    print(
        "_canSupportsBiometricAuthenticated() $supportsBiometricAuthenticated");
    return supportsBiometricAuthenticated;
  }

  Future<CanAuthenticateResponse>
      _checkBiometricCapabilitiesCanAuthenticate() async {
    final response = await BiometricStorage().canAuthenticate();
    print('BIOMETRIC: checked if authentication was possible: $response');
    return response;
  }

  Future<BiometricStorageFile?> _getBiometricStorageFile(
      baseName, supportsAuthenticated) async {
    print("_canDoBio($baseName, $supportsAuthenticated)");
    if (supportsAuthenticated) {
      final _authStorageFile = await BiometricStorage().getStorage(
          '${baseName}_authenticated',
          options: StorageFileInitOptions(
              authenticationValidityDurationSeconds: Platform.isIOS ? 0 : 30));
      print("_canDoBio($baseName, $supportsAuthenticated) $_authStorageFile");
      return _authStorageFile;
    }
    return null;
  }

  Future<BiometricStorageFile?> _getTxBiometricStorageFile(
      baseName, supportsAuthenticated) async {
    print("_canDoBio($baseName, $supportsAuthenticated)");
    if (supportsAuthenticated) {
      final _authStorageFile = await BiometricStorage().getStorage(
          '${baseName}_tx_authenticated',
          options: StorageFileInitOptions(
              authenticationValidityDurationSeconds: Platform.isIOS ? 0 : 30));
      print("_canDoBio($baseName, $supportsAuthenticated) $_authStorageFile");
      return _authStorageFile;
    }
    return null;
  }
}
