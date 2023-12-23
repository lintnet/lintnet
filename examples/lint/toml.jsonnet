function(data)
  local message = 'TOML requires the field "name"';
  [{
    message: message,
    failed: data.file_type == 'toml' && !std.objectHas(data.value, 'name'),
  }]
