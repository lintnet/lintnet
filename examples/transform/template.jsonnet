function(param)
  local rules = std.set([
    {
      name: e.name,
      [if std.objectHas(e, 'description') then 'description']: std.toString(e.description),
      errors: [
        {
          [if std.objectHas(e2, 'message') then 'message']: std.toString(e2.message),
          [if std.objectHas(e2, 'level') then 'level']: std.toString(e2.level),
          [if std.objectHas(e2, 'location') then 'location']: std.toString(e2.location),
          [if std.objectHas(e2, 'custom') then 'custom']: std.toString(e2.custom),
        }
        for e2 in param.errors
        if e2.name == e.name
      ],
    }
    for e in param.errors
  ], function(e) e.name);
  rules
