local contains = std.native('strings.Contains');

function(param)
  if contains(param.data.file_path, '-') then
    [{
      message: 'Use underscores "_" rather than dash in file names',
    }] else []
