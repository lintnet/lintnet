function(param) [{
  message: 'Never give up!',
  failed: std.native('strings.Contains')(param.data.text, 'Give up'),
}]
