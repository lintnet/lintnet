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
      result: common.lint_result,
    },
  },
}
