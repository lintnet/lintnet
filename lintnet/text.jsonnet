local input = std.extVar('input');
local text = std.extVar('file_text');

{
  name: 'Never give up!',
  failed: std.native('strings.Contains')(text, 'Give up'),
}
