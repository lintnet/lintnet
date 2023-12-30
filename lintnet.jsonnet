function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/**/*.jsonnet@429ef22b1fe1d5bb85b4f420157c25b41b590b10',
        '!github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/**/*_test.jsonnet@429ef22b1fe1d5bb85b4f420157c25b41b590b10',
      ],
    },
  ],
}
