function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github_archive/github.com/lintnet-modules/ghalint/workflow/**/main.jsonnet@0f350f659c7c64c7398249ea0fc23d1cec45c12a:v0.2.0',
        'github_archive/github.com/lintnet-modules/github-actions/workflow/**/main_combine.jsonnet@eb941dd42ce4ec800588fb2b4d822c591dd54364:v0.2.0',
      ],
    },
  ],
}
