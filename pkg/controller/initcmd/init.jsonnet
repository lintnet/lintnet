function(param) {
  error_level: 'error',
  targets: [
    {
      data_files: [
        'examples/data/hello.csv',
      ],
      lint_files: [
        'examples/lint/csv.jsonnet',
      ],
      modules: [
        'github.com/lintnet/lintnet/examples/lint/csv.jsonnet@07e8eebe7886562380615b663c52007fb8342b51',
      ],
    },
  ],
}
