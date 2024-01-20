function(param) {
  targets: [
    {
      data_files: [
        // Glob is available
        // e.g. *.json, **/*.json
        'foo.json',
      ],
      lint_files: [
        // Glob is available
        'hello.jsonnet',
      ],
    },
  ],
}
