local filePath = std.extVar('file_path');
local nf = {
  contains: std.native('strings.Contains'),
};

{
  name: 'File Name',
  description: |||
    Lint rules regarding file names
  |||,
  sub_rules: [
    {
      name: 'use_underscore_rather_than_dash',
      description: |||
        Use underscores "_" rather than dash.
      |||,
      failed: nf.contains(filePath, '-'),
    },
  ],
}
