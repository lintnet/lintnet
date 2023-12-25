function(param) [
  {
    name: 'fail',
    data_file: 'foo.csv',
    data_text: |||
      mike,1
    |||,
    param: {
      data: {
        file_type: 'csv',
        value: [
          ['mike', '1'],
        ],
      },
    },
    result: [
      {
        message: 'age must be greater or equal than 18',
        failed: true,
        level: 'warn',
        location: {
          index: 0,
          line: 'mike,1',
        },
      },
    ],
  },
  {
    name: 'succees',
    param: {
      data: {
        file_type: 'csv',
        value: [
          ['mike', '20'],
        ],
      },
    },
    result: [
      {
        message: 'age must be greater or equal than 18',
        failed: false,
        level: 'warn',
        location: {
          index: 0,
          line: 'mike,20',
        },
      },
    ],
  },
  {
    name: 'not csv',
    param: {
      data: {
        file_type: 'text',
      },
    },
    result: null,
  },
]
