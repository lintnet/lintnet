local input = std.extVar('input');
local text = std.extVar('file_text');

[{
  message: 'Never give up!',
  failed: std.native('strings.Contains')(text, 'Give up'),
}]
