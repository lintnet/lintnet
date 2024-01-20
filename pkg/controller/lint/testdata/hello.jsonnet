function(param)
  if std.objectHas(param.data.value, 'description') then [] else [{
    name: 'description is required',
  }]
