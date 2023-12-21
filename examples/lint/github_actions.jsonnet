local fileType = std.extVar('file_type');
local filePath = std.extVar('file_path');
local input = std.extVar('input');

if fileType == 'yaml' && std.startsWith(filePath, '.github/workflows/') then [
  {
    name: 'uses_must_not_be_main',
    message: 'actions reference must not be main',
    location: {
      job_name: job.key,
      uses: step.uses,
    },
  }
  for job in std.objectKeysValues(input.jobs)
  if std.objectHas(job.value, 'steps')
  for step in job.value.steps
  if std.objectHas(step, 'uses') && std.endsWith(step.uses, '@main')
] else []
