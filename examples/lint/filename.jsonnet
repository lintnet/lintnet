local contains = std.native('strings.Contains');

function(param) [{
  message: 'Use underscores "_" rather than dash in file names',
  failed: contains(param.data.file_path, '-'),
}]
