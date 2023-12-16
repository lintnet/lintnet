local fileType = std.extVar('file_type');
local input = std.extVar('input');

local test() = std.filter(
  function(line) line != null,
  std.mapWithIndex(function(idx, line) if std.parseInt(line[1]) >= 18 then null else {
    index: idx,
    line: std.join(',', line),
  }, input)
);

{
  name: 'CSV',
  description: |||
    Lint rules regarding CSV
  |||,
  sub_rules: [
    {
      name: 'age must be greater than 18',
      locations: if fileType == 'csv' || fileType == 'tsv' then test() else null,
      'error': 'age must be greater than 18',
    },
  ],
}
