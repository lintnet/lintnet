package domain

import "github.com/lintnet/lintnet/pkg/config"

type Target struct {
	ID        string
	LintFiles []*config.LintFile
	DataFiles Paths
}

type Paths []*Path

func (ps Paths) Raw() []string {
	arr := make([]string, len(ps))
	for i, p := range ps {
		arr[i] = p.Raw
	}
	return arr
}

type Path struct {
	Raw string
	Abs string
}
