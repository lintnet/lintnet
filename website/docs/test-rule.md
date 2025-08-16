---
sidebar_position: 400
---

# Test lint rules

To test lint file `x.jsonnet`, you need to create a test file `x_test.jsonnet` on the same directory with `x.jsonnet`.

## Schema

JSON Schema: https://github.com/lintnet/lintnet/blob/main/json-schema/test-result.json

A test file is a function returning a list of test case.

```jsonnet
function(param) [
  {
    name: 'Test case name',
    // ...
  }
  // ...
]
```

A test case is a pair of test data and expected result.

```jsonnet
{
  name: 'Test case name',
  // test data
  // ...
  param: {
    // config is configuration passed to the lint file
    // config is optional.
    config: {},
  },
  result: [
    // expected return value of the lint file
  ],
}
```

About test data, you can refer to test data files with `data_file` and `data_files` fields.

```jsonnet
function(param) [
  {
    data_file: 'testdata/pass.json', // relative path from the test file
  },
  // ...
]
```

`fake_data_file` is useful to disguise a data file path for testing.
If you set `fake_data_file`, a test data is read from `data_file` but the data file path passed to the lint file is disguised as `fake_data_file`.

```jsonnet
data_file: 'testdata/pass.yaml',
fake_data_file: '.github/workflows/pass.yaml',
```

If [a lint file lints across multiple files](/docs/guides/lint-across-files/), `data_files` is used instead of `data_file`.

```jsonnet
{
  data_files: [
    // a list of data files.
    // The element is either a string or an object.
    'testdata/pass.json',
    {
      path: 'testdata/foo.json',
      fake_path: '/etc/app/foo.json',
    },
  ],
},
```

Instead of `data_file` and `data_files`, you can also define data in a test file directly using `param`.
But we recommend `data_file` and `data_files` because they are more maintainable.

```jsonnet
{
  param: {
    // These fields are optional.
    // You only have to set fields used in the lint file.
    data: {
      file_path: 'foo.json',
      file_type: 'json',
      text: '', // raw text
      value: {
        // parsed data
      },
    },
  },
  // ...
},
```

About expected result, the format depends on the lint rule.

```jsonnet
{
  result: [ // expected return value of the lint file
    {
      name: 'age must be greater or equal than 18',
      level: 'error',
      location: {
        index: 0,
        line: 'mike,1',
      },
    },
  ],
},
```

## Scaffold a test file

```sh
lintnet new [<lint file name | main.jsonnet>]
```

## Run test

```sh
lintnet test [<lint file, test file, or directory> ...]
```

If you run `lintnet test` without any argument, lintnet searches lint files using a configuration file and tests all lint files having test files.
lint files without test files are ignored.
You can test only specific files by specifying files as command line arguments.
If you specify files explicitly, a configuration file is unnecessary.
This means when you develop modules, you don't have to prepare a configuration file.

If you specify directories, lint files in those directories and subdirectories are tested.
For example, `lintnet test .` searches files matching the glob pattern `**/*.jsonnet` and `lintnet test foo` search files matching `foo/**/*.jsonnet`.

If a configuration file isn't specified and isn't found, `lintnet test` works as `lintnet test .`.

## Normalization of evaluation result

The evaluation result of lint file is normalized before it is compared with `result`.

- `description` and `excluded` are removed
- Array elements whose `excluded` is `true` are excluded

For example, If the evaluation result of lint file is as the following,

```json
[
  {
    "name": "foo",
    "description": "Hello, lintnet!",
    "excluded": true
  },
  {
    "name": "foo",
    "description": "Hello, lintnet!",
    "excluded": false
  }
]
```

`result` must be as the following.

```json
[
  {
    "name": "foo"
  }
]
```

This normalization simplifies `result`.
