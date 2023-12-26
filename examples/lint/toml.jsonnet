function(param)
  local name = 'TOML requires the field "name"';
  if param.data.file_type == 'toml' && !std.objectHas(param.data.value, 'name') then
    [{
      name: name,
    }] else []
