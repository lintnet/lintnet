local lint_data = {
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
};

{
  lint_tla: {
    type: 'object',
    additionalProperties: false,
    description: 'Top level arguments',
    required: [
      'config',
    ],
    properties: {
      data: lint_data,
      combined_data: {
        type: 'array',
        description: 'A list of data. This is set if combine is true',
        items: lint_data,
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
        links: {
          oneOf: [
            {
              type: 'object',
              description: 'each key is a link title',
              additionalProperties: {
                type: 'string',
                description: 'link',
              },
            },
            {
              type: 'array',
              items: {
                oneOf: [
                  {
                    type: 'string',
                    description: 'link',
                  },
                  {
                    type: 'object',
                    additionalProperties: false,
                    description: 'result',
                    required: [
                      'link',
                    ],
                    properties: {
                      title: {
                        type: 'string',
                        description: 'link title',
                      },
                      link: {
                        type: 'string',
                        description: 'link',
                      },
                    },
                  },
                ],
              },
            },
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
