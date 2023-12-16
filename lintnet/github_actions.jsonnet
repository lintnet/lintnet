local fileType = std.extVar('file_type');
local filePath = std.extVar('file_path');
local input = std.extVar('input');
{
  name: 'GitHub Actions',
  description: |||
    Lint rules regarding GitHub Actions
  |||,
  sub_rules: [
    {
      name: 'uses_must_not_be_main',
      description: |||
        actions reference must not be main.
      |||,
      locations: if fileType == 'yaml' && std.startsWith(filePath, '.github/workflows/') then [
        {
          job_name: job.key,
          uses: step.uses,
        }
        for job in std.objectKeysValues(input.jobs)
        if std.objectHas(job.value, 'steps')
        for step in job.value.steps
        if std.objectHas(step, 'uses') && std.endsWith(step.uses, '@main')
      ] else null,
    },
  ],
}
