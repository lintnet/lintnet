function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github_archive/github.com/lintnet-modules/ghalint/workflow/**/main.jsonnet@eef8b404583a671005e9a9997f11f0c91c73e1de:v0.3.0-1',
        'github_archive/github.com/lintnet-modules/github-actions/workflow/**/main_combine.jsonnet@eb941dd42ce4ec800588fb2b4d822c591dd54364:v0.2.0',
      ],
    },
  ],
}
