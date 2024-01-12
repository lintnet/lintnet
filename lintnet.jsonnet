function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github.com/lintnet/modules/modules/ghalint/**/main.jsonnet@41ea96238c2455f85796446e4fa77f2716c827db',
        'github.com/lintnet/modules/modules/github_actions/**/main_combine.jsonnet@41ea96238c2455f85796446e4fa77f2716c827db',
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
