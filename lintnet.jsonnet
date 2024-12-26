function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github_archive/github.com/lintnet-modules/ghalint/workflow/**/main.jsonnet@12aac7476916a42e9de8646ac75c98274cfe8521:v0.3.2',
        'github_archive/github.com/lintnet-modules/github-actions/workflow/**/main_combine.jsonnet@eb941dd42ce4ec800588fb2b4d822c591dd54364:v0.2.0',
      ],
    },
  ],
}
