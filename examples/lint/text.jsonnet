function(param) if std.native('strings.Contains')(param.data.text, 'Give up') then
  [{
    name: 'Never give up!',
  }] else []
