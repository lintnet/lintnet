package jsonnet

import (
	"github.com/google/go-jsonnet"
	"github.com/suzuki-shunsuke/go-jsonnet-native-functions/pkg/path/filepath"
	"github.com/suzuki-shunsuke/go-jsonnet-native-functions/pkg/regexp"
	"github.com/suzuki-shunsuke/go-jsonnet-native-functions/pkg/strings"
)

func SetNativeFunctions(vm *jsonnet.VM) {
	vm.NativeFunction(strings.Contains("strings.Contains"))
	vm.NativeFunction(strings.TrimPrefix("strings.TrimPrefix"))
	vm.NativeFunction(strings.TrimSpace("strings.TrimSpace"))
	vm.NativeFunction(regexp.MatchString("regexp.MatchString"))
	vm.NativeFunction(filepath.Base("filepath.Base"))
}

func MakeVM() *jsonnet.VM {
	return jsonnet.MakeVM()
}

func NewVM(param string, importer jsonnet.Importer) *jsonnet.VM {
	vm := jsonnet.MakeVM()
	vm.TLACode("param", param)
	SetNativeFunctions(vm)
	vm.Importer(importer)
	return vm
}
