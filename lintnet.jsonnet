function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github_archive/github.com/lintnet-modules/ghalint/workflow/**/main.jsonnet@0d6f9c5dbc856a70fca35511136d4f1c3195c872:v0.3.1',
        'github_archive/github.com/lintnet-modules/github-actions/workflow/**/main_combine.jsonnet@eb941dd42ce4ec800588fb2b4d822c591dd54364:v0.2.0',
      ],
    },
  ],
}
