function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github.com/lintnet/modules/modules/ghalint/**/main.jsonnet@c3bbbab20cc7f94304796e8ca1e6056f4a06fe0a',
      ],
    },
    {
      combine: true,
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github.com/lintnet/modules/modules/github_actions/**/main.jsonnet@c3bbbab20cc7f94304796e8ca1e6056f4a06fe0a',
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
