package main

import (
	"github.com/lintnet/lintnet/pkg/cli"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

var version = ""

func main() {
	urfave.Main("lintnet", version, cli.Run)
}
