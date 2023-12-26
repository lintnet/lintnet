function(param) [
  {
    name: 'fail',
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
        level: 'error',
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
        excluded: true,
        level: 'error',
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
  {
    name: 'data_file',
    data_file: '../data/hello.csv',
    result: [
      {
        message: 'age must be greater or equal than 18',
        excluded: false,
        level: 'error',
        location: {
          index: 0,
          line: 'mike,10',
        },
      },
    ],
  },
]
