function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github.com/lintnet/modules/modules/ghalint/**/main.jsonnet@805119063d195ffbafb3b0509704e5239741f86c:v0.1.1',
        'github.com/lintnet/modules/modules/github_actions/**/main_combine.jsonnet@805119063d195ffbafb3b0509704e5239741f86c:v0.1.1',
      ],
    },
  ],
}
