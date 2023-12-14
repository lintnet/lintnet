package lint

import (
	"github.com/google/go-jsonnet"
	"github.com/suzuki-shunsuke/go-jsonnet-native-functions/pkg/path/filepath"
	"github.com/suzuki-shunsuke/go-jsonnet-native-functions/pkg/regexp"
	"github.com/suzuki-shunsuke/go-jsonnet-native-functions/pkg/strings"
)

func setNativeFunctions(vm *jsonnet.VM) {
	vm.NativeFunction(strings.Contains("strings.Contains"))
	vm.NativeFunction(strings.TrimPrefix("strings.TrimPrefix"))
	vm.NativeFunction(strings.TrimPrefix("strings.TrimSpace"))
	vm.NativeFunction(regexp.MatchString("regexp.MatchString"))
	vm.NativeFunction(filepath.Base("filepath.Base"))
}
