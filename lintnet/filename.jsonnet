local fileType = std.extVar('file_type');
local input = std.extVar('input');
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
        {
          message: '',
          job_name: job.key,
          uses: step.uses,
        }
        for job in std.objectKeysValues(input.jobs)
        if std.objectHas(job.value, 'steps')
        for step in job.value.steps
        if fileType == 'yaml' && std.objectHas(step, 'uses') && std.endsWith(step.uses, '@main')
      ],
    },
  ],
}
