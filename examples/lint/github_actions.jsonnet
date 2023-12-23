function(param)
  if param.data.file_type == 'yaml' && std.startsWith(param.data.file_path, '.github/workflows/') then [
    {
      name: 'uses_must_not_be_main',
      message: 'actions reference must not be main',
      location: {
        job_name: job.key,
        uses: step.uses,
      },
    }
    for job in std.objectKeysValues(param.data.value.jobs)
    if std.objectHas(job.value, 'steps')
    for step in job.value.steps
    if std.objectHas(step, 'uses') && std.endsWith(step.uses, '@main')
  ] else []
