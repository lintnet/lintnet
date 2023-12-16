local fileType = std.extVar('file_type');
local input = std.extVar('input');

if fileType != 'csv' then null else std.mapWithIndex(function(idx, line) {
  message: 'age must be greater or equal than 18',
  failed: std.parseInt(line[1]) < 18,
  level: 'warn',
  location: {
    index: idx,
    line: std.join(',', line),
  },
}, input)
