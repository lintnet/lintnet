function(param)
  local message = 'TOML requires the field "name"';
  [{
    message: message,
    failed: param.data.file_type == 'toml' && !std.objectHas(param.data.value, 'name'),
  }]
