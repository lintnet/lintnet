function(param) {
  targets: [
    {
      data_files: [
        'examples/data/hello.csv',
        'examples/data/hello.tsv',
      ],
      lint_files: [
        'examples/lint/csv.jsonnet',
        'examples/lint/filename.jsonnet',
      ],
      modules: [
        'github.com/lintnet/lintnet/examples/lint/csv.jsonnet@07e8eebe7886562380615b663c52007fb8342b51',
      ],
    },
    {
      data_files: [
        'examples/data/hello.toml',
      ],
      lint_files: [
        'examples/lint/toml.jsonnet',
        'examples/lint/import_module.jsonnet',
      ],
    },
    {
      data_files: [
        'examples/data/hello.txt',
      ],
      lint_files: [
        'examples/lint/text.jsonnet',
      ],
    },
    {
      data_files: [
        'examples/data/test.yaml',
      ],
      lint_files: [
        'examples/lint/github_actions.jsonnet',
      ],
    },
    {
      data_files: [
        'examples/data/hello.hcl',
      ],
      lint_files: [
        'examples/lint/service_forbid_public.jsonnet',
      ],
    },
  ],
}
