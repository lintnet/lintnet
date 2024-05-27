// A configuration file of lintnet.
// https://lintnet.github.io/
function(param) {
  targets: [
    {
      data_files: [
        // Relative paths from this configuration file.
        // '.github/workflows/*.yaml',
        // '**/*',
        // '!.gitkeep',
      ],
      // lint_files: [
      //   // Same as data_files.
      //   'examples/lint/github_actions.jsonnet',
      //   {
      //     path: 'examples/lint/filename.jsonnet',
      //     config: {
      //       excluded: ['foo'],
      //     },
      //   },
      // ],
      // modules: [
      //   'github_archive/github.com/lintnet-modules/ghalint/workflow/**/main.jsonnet@0f350f659c7c64c7398249ea0fc23d1cec45c12a:v0.2.0',
      //   {
      //     path: 'github_archive/github.com/lintnet-modules/ghalint@0f350f659c7c64c7398249ea0fc23d1cec45c12a:v0.2.0',
      //     files: [
      //       'workflow/**/main.jsonnet',
      //       '!workflow/action_ref_should_be_full_length_commit_sha/main.jsonnet',
      //       {
      //         path: 'workflow/action_ref_should_be_full_length_commit_sha/main.jsonnet',
      //         config: {
      //           excludes: [
      //             'slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml',
      //           ],
      //         },
      //       },
      //     ],
      //   },
      // ],
    },
  ],
}
