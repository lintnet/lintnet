function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/**/*.jsonnet@429ef22b1fe1d5bb85b4f420157c25b41b590b10',
      ],
    },
    {
      // lintnet lint -target-id foo [<data file> ...]
      id: 'foo',
      combine: true,
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      lint_files: [
        'examples/lint/github_actions_workflow_name_should_be_unique.jsonnet',
      ],
    },
  ],
  outputs: [
    {
      id: 'jsonnet',
      renderer: 'jsonnet',
      template: 'examples/output/output.jsonnet',
      // config: {},
    },
    {
      id: 'template',
      renderer: 'text/template',
      template: 'examples/output/output.tpl',
      transform: 'examples/transform/template.jsonnet',
      // config: {},
    },
  ],
}
