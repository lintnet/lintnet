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

type Data struct {
	Text     string `json:"text"`
	Value    any    `json:"value"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
	JSON     []byte `json:"-"`
}

type TopLevelArgment struct {
	Data         *Data          `json:"data,omitempty"`
	CombinedData []*Data        `json:"combined_data,omitempty"`
	Config       map[string]any `json:"config"`
}

type DataSet struct {
	File  *Path
	Files []*Path
}
