---
sidebar_position: 530
---

# Learn Jsonnet for lintnet

First, you should read Jsonnet's official documents.

- [Tutorial](https://jsonnet.org/learning/tutorial.html)
- [Editor plugin, linter, formatter, and so on](https://jsonnet.org/learning/tools.html)
- [Standard library](https://jsonnet.org/ref/stdlib.html)
- [Language reference](https://jsonnet.org/ref/language.html)
   - [Equivalence and Equality](https://jsonnet.org/ref/language.html#equivalence-and-equality)

## Check Jsonnet behaviour

When you write lintnet rules, sometimes you want to check how Jsonnet works.
In that case, you can check this in some ways.

1. Web editor
1. Jsonnet CLI

### 1. Web editor

You can try Jsonnet instantly using the web editor on https://jsonnet.org/learning/tutorial.html .

![image](https://github.com/lintnet/lintnet/assets/13323303/408441c3-9c1d-4ff9-99a7-f37a17d0e297)

### 2. Jsonnet CLI

You can also execute Jsonnet with Jsonnet CLI.
There are two implementation, C++ implementation and Go implementation ([go-jsonnet](https://github.com/google/go-jsonnet)).
lintnet uses go-jsonnet, so go-jsonnet would be better.

You can install go-jsonnet with Homebrew and [aqua](https://aquaproj.github.io/).

```sh
brew install go-jsonnet
```

```sh
aqua g -i google/go-jsonnet
```

You can evaluate Jsonnet by jsonnet command.

```sh
jsonnet hello.jsonnet
```

## Standard library

https://jsonnet.org/ref/stdlib.html

We pick out some functions that we often use.

- std.type(x): type check
- std.length(x): check the size of array and object
- std.get(o, f, default=null, inc_hidden=true): get the object attribute with a default value
- std.objectHas(o, f): Check if the attribute exists
- std.objectKeysValues(o): convert object to array
- std.startsWith(a, b), std.endsWith(a, b)
- std.map(func, arr), std.mapWithIndex(func, arr), std.filterMap(filter_func, map_func, arr), std.filter(func, arr)
- std.set(arr, keyF=id)
- std.sort(arr, keyF=id)

## local values

```jsonnet
function(param) 
  local foo = 'foo';
  {
    // Define local values in object definitions
    local factor = if large then 2 else 1,
    // Add attributes to objects conditionally
    [if salted then 'garnish']: 'Salt',
  }
```

## Parameterize Entire Config

https://jsonnet.org/learning/tutorial.html#parameterize-entire-config

lintnet uses Top-level arguments, not External variables.
