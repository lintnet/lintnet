local input = std.extVar('input');
local fileType = std.extVar('file_type');

{
  name: 'TOML requires the field "name"',
  failed: fileType == 'toml' && !std.objectHas(input, 'name'),
}
