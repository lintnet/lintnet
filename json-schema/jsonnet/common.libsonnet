{
  lint_tla: {
    type: 'object',
    additionalProperties: false,
    description: 'Top level arguments',
    required: [
      'data',
      'config',
    ],
    properties: {
      data: {
        type: 'object',
        description: 'data file',
        additionalProperties: false,
        required: [
          'file_path',
          'text',
          'value',
          'file_type',
        ],
        properties: {
          file_path: {
            type: 'string',
            description: 'data file path',
          },
          text: {
            type: 'string',
            description: 'data file content',
          },
          value: {
            description: 'data file content',
          },
          file_type: {
            type: 'string',
            description: 'data file type',
            enum: [
              'csv',
              'hcl2',
              'json',
              'plain_text',
              'toml',
              'tsv',
              'yaml',
            ],
          },
        },
      },
      config: {
        type: 'object',
        description: 'configuration',
      },
    },
  },
  lint_result: {
    type: 'array',
    description: 'results',
    items: {
      type: 'object',
      additionalProperties: false,
      description: 'result',
      required: [
        'name',
      ],
      properties: {
        name: {
          type: 'string',
          description: 'rule name',
        },
        message: {
          type: 'string',
          description: 'error message',
        },
        description: {
          type: 'string',
          description: 'rule description',
        },
        level: {
          type: 'string',
          description: 'error level',
          enum: [
            'debug',
            'info',
            'warn',
            'error',
          ],
        },
        location: {
          type: 'object',
          description: 'Location where errors occur',
        },
        excluded: {
          type: 'boolean',
          description: 'Whether the result is excluded',
        },
        custom: {
          type: 'object',
          description: 'Custom fields that users can set freely',
        },
      },
    },
  },
}
