function(data) [{
  message: 'Never give up!',
  failed: std.native('strings.Contains')(data.text, 'Give up'),
}]
