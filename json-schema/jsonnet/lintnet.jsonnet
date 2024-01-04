{
  '$schema': 'https://json-schema.org/draft/2020-12/schema',
  additionalProperties: false,
  type: 'object',
  properties: {
    targets: {
      description: 'targets',
      type: 'array',
      items: {
        type: 'object',
        additionalProperties: false,
        description: 'target',
        properties: {
          data_files: {
            type: 'array',
            description: 'data files',
            items: {
              type: 'string',
              description: 'file path to data files. Glob is available',
            },
          },
          lint_files: {
            type: 'array',
            description: 'lint files',
            items: {
              description: 'lint files',
              anyOf: [
                {
                  type: 'string',
                  description: 'file path to lint files. Glob is available',
                },
                {
                  type: 'object',
                  required: [
                    'path',
                  ],
                  properties: {
                    path: {
                      type: 'string',
                      description: 'file path to lint files. Glob is available',
                    },
                    config: {
                      type: 'object',
                      description: 'configuration of the lint files',
                    },
                  },
                },
              ],
            },
          },
          modules: {
            type: 'array',
            description: 'modules',
            items: {
              description: 'modules',
              anyOf: [
                {
                  type: 'string',
                  description: 'file path to lint files. Glob is available',
                },
                {
                  type: 'object',
                  required: [
                    'path',
                  ],
                  properties: {
                    path: {
                      type: 'string',
                      description: 'file path to lint files. Glob is available',
                    },
                    config: {
                      type: 'object',
                      description: 'configuration of the lint files',
                    },
                  },
                },
              ],
            },
          },
        },
      },
    },
    outputs: {
      description: 'outputs',
      type: 'array',
      items: {
        type: 'object',
        additionalProperties: false,
        required: [
          'id',
          'renderer',
          'template',
        ],
        properties: {
          id: {
            type: 'string',
            description: 'output id',
          },
          renderer: {
            type: 'string',
            description: 'renderer',
            enum: [
              'jsonnet',
              'text/template',
              'html/template',
            ],
          },
          template: {
            type: 'string',
            description: 'file path to template',
          },
          transform: {
            type: 'string',
            description: 'file path to Jsonnet to transform results',
          },
          config: {
            type: 'object',
            description: 'configuration of transform and output',
          },
        },
      },
    },
  },
}
