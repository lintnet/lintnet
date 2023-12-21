local contains = std.native('strings.Contains');

function(data) [{
  message: 'Use underscores "_" rather than dash in file names',
  failed: contains(data.file_path, '-'),
}]
