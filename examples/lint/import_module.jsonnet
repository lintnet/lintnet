local hello = import 'github.com/suzuki-shunsuke/example-lintnet-modules/hello.jsonnet@923da72cf3330c7710393b86a0e2f4bea533ff51';

function(param) [{
  message: hello.message,
  failed: true,
}]
