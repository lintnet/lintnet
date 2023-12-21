local filePath = std.extVar('file_path');
local contains = std.native('strings.Contains');

[{
  message: 'Use underscores "_" rather than dash in file names',
  failed: contains(filePath, '-'),
}]
