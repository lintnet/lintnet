// A lint file of lintnet.
// https://lintnet.github.io/
function(param)
  if std.objectHas(param.data.value, 'description') then [] else [{
    name: 'description is required',
  }]
