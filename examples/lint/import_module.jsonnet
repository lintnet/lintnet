local hello = import 'github.com/suzuki-shunsuke/example-lintnet-modules/hello.jsonnet@03d2ecad06b8c7a980e677ce81387f0c3fe6461b:v0.1.1';

function(param) [{
  message: hello.message,
}]
