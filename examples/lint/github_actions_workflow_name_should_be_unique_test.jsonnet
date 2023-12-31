function(param) [
  {
    name: 'fail',
    param: {
      combined_data: [
        {
          value: {
            name: 'build',
          },
          file_path: '.github/workflows/build.yaml',
        },
        {
          value: {
            name: 'build',
          },
          file_path: '.github/workflows/test.yaml',
        },
        {
          value: {
            name: 'release',
          },
          file_path: '.github/workflows/release.yaml',
        },
      ],
    },
    result: [
      {
        name: 'GitHub Actions workflow name must be unique',
        location: {
          workflow_name: 'build',
          files: [
            '.github/workflows/build.yaml',
            '.github/workflows/test.yaml',
          ],
        },
      },
    ],
  },
  {
    name: 'succees',
    param: {
      combined_data: [
        {
          value: {
            name: 'build',
          },
          file_path: '.github/workflows/build.yaml',
        },
        {
          value: {
            name: 'test',
          },
          file_path: '.github/workflows/test.yaml',
        },
        {
          value: {
            name: 'release',
          },
          file_path: '.github/workflows/release.yaml',
        },
      ],
    },
    result: [],
  },
]
