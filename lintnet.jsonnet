function(param) {
  targets: [
    {
      data_files: [
        '.github/workflows/*.yml',
        '.github/workflows/*.yaml',
      ],
      modules: [
        'github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/**/*.jsonnet@0ed62adf055a4fbd7ef7ebe304f01794508ed325',
        '!github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/**/*_test.jsonnet@0ed62adf055a4fbd7ef7ebe304f01794508ed325',
      ],
    },
  ],
}
