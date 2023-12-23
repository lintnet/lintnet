function(param) if param.data.file_type != 'hcl2' then null else std.mapWithIndex(function(idx, elem) {
  message: 'service must not be public',
  failed: std.get(elem.value, 'public', false),
  location: '%s[%d]' % [elem.key, idx],
}, [
  {
    key: '%s.%s.%s' % [x.key, y.key, z.key],
    value: elem,
  }
  for x in std.objectKeysValues(param.data.value)
  for y in std.objectKeysValues(x.value)
  for z in std.objectKeysValues(y.value)
  for elem in z.value
])
