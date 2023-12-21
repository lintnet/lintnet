function(data) [{
  message: 'TOML requires the field "name"',
  failed: data.file_type == 'toml' && !std.objectHas(data.value, 'name'),
}]
