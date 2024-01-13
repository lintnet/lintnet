local common = import 'common.libsonnet';

{
  '$schema': 'https://json-schema.org/draft/2020-12/schema',
  additionalProperties: false,
  type: 'array',
  description: 'Test results',
  items: {
    type: 'object',
    description: 'data file',
    required: [
      'name',
      'param',
      'result',
    ],
    properties: {
      name: {
        type: 'string',
        description: 'test name',
      },
      data_file: {
        type: 'string',
        description: "data file path. This overrides param.data's text, value and param.data.text",
      },
      param: common.lint_tla,
      result: {
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
              oneOf: [
                {
                  type: 'object',
                  description: 'Location where errors occur',
                  additionalProperties: true,
                },
                {
                  type: 'string',
                  description: 'Location where errors occur',
                },
              ],
            },
            custom: {
              type: 'object',
              description: 'Custom fields that users can set freely',
              additionalProperties: true,
            },
          },
        },
      },
    },
  },
}
