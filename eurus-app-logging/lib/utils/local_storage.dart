import 'dart:io';

/// Check folder exists
Future<bool> ckPathExists(String path) async {
  return await Directory(path).exists();
}

/// Create folder
Future<String> createFolder(String path) async {
  final Directory directory = await Directory(path).create(recursive: true);

  return directory.path;
}

/// Check File exists
Future<bool> ckFileExists(String path, String filename) async {
  return await File("$path/$filename").exists();
}

/// Write [String] to file
Future<bool> writeFile(String path, String filename, String log) async {
  File file = _getFile(path, filename);

  await file.writeAsString(log, mode: FileMode.append);
  return true;
}

/// Read log file and return logs in [List<String>]
Future<List<String>> readFile(String path, String filename) async {
  bool exists = await ckFileExists(path, filename);
  if (exists) {
    File file = _getFile(path, filename);
    String rawLog = await file.readAsString();

    List<String> logs = rawLog.split("\n");
    return logs;
  }
  return null;
}

/// Delete file
Future<bool> deleteFile(String path, String filename) async {
  bool fileExists = await ckFileExists(path, filename);

  /// Return [true] if file does not exists, treade it as deleted
  if (!fileExists) return true;

  /// Try to delete file
  try {
    File file = _getFile(path, filename);
    await file.delete();
  } catch (e) {
    throw ('Unable to delete file');
  }

  return true;
}

/// Create file
File _getFile(String path, String filename) {
  return File("$path/$filename");
}
