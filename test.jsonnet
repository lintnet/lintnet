{
  forbid_main: {
    message: 'Action version should not be main',
    data: [
      {
        job_name: job.key,
        uses: step.uses,
      }
      for job in std.objectKeysValues(std.extVar('input').jobs)
      if std.objectHas(job.value, 'steps')
      for step in job.value.steps
      if std.objectHas(step, 'uses') && std.endsWith(step.uses, '@main')
    ],
  },
}
