local filePath = std.extVar('file_path');
local nf = {
  contains: std.native('strings.contains'),
};

{
  name: 'File Name',
  description: |||
    Lint rules regarding file names
  |||,
  rules: [
    {
      name: 'use_underscore_rather_than_dash',
      description: |||
        Use underscores "_" rather than dash.
      |||,
      errors: [
        elem
        for elem in [{}]
        if nf.contains(filePath, '-')
      ],
    },
  ],
}
