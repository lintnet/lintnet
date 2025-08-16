package domain

type Paths []*Path

func (ps Paths) Raw() []string {
	arr := make([]string, len(ps))
	for i, p := range ps {
		arr[i] = p.Raw
	}
	return arr
}

type Path struct {
	Raw string `json:"raw,omitempty"`
	Abs string `json:"abs,omitempty"`
}

type Data struct {
	Text     string `json:"text"`
	Value    any    `json:"value"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
	JSON     []byte `json:"-"`
}

type TopLevelArgument struct {
	Data         *Data          `json:"data,omitempty"`
	CombinedData []*Data        `json:"combined_data,omitempty"`
	Config       map[string]any `json:"config"`
}

type DataSet struct {
	File  *Path
	Files []*Path
}
