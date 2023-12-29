function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/**/*.jsonnet@764ccddf94a5dcc1c9d619dedaebfc64f0251a04',
        '!github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/**/*_test.jsonnet@764ccddf94a5dcc1c9d619dedaebfc64f0251a04',
      ],
    },
  ],
}
