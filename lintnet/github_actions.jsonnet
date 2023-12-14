local fileType = std.extVar('file_type');
local input = std.extVar('input');
{
  name: 'GitHub Actions',
  description: |||
    Lint rules regarding GitHub Actions
  |||,
  rules: [
    {
      name: 'uses_must_not_be_main',
      description: |||
        actions reference must not be main.
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
