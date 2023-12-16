local fileType = std.extVar('file_type');
local input = std.extVar('input');

std.mapWithIndex(function(idx, line) {
  message: 'age must be greater or equal than 18',
  failed: std.parseInt(line[1]) < 18,
  location: {
    index: idx,
    line: std.join(',', line),
  },
}, input)
