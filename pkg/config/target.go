package config

import "path/filepath"

type Target struct {
	ID             string                    `json:"id,omitempty"`
	LintFiles      []*ModuleGlob             `json:"lint_files,omitempty"`
	Modules        []*ModuleGlob             `json:"modules,omitempty"`
	ModuleArchives map[string]*ModuleArchive `json:"module_archives,omitempty"`
	DataFiles      []string                  `json:"data_files,omitempty"`
}

type RawTarget struct {
	ID        string       `json:"id,omitempty"`
	LintGlobs []*LintGlob  `json:"lint_files"`
	Modules   []*RawModule `json:"modules"`
	DataFiles []string     `json:"data_files"`
}

func (rt *RawTarget) Parse() (*Target, error) {
	lintFiles := make([]*ModuleGlob, len(rt.LintGlobs))
	for i, lintGlob := range rt.LintGlobs {
		lintFiles[i] = lintGlob.ToModule()
	}
	dataFiles := make([]string, len(rt.DataFiles))
	for i, dataFile := range rt.DataFiles {
		dataFiles[i] = filepath.Clean(dataFile)
	}
	target := &Target{
		ID:        rt.ID,
		LintFiles: lintFiles,
		Modules:   make([]*ModuleGlob, len(rt.Modules)),
		DataFiles: dataFiles,
	}
	archives := make(map[string]*ModuleArchive, len(rt.Modules))
	for i, m := range rt.Modules {
		a, err := m.Parse()
		if err != nil {
			return nil, err
		}
		target.Modules[i] = a
		archives[a.Archive.String()] = a.Archive
	}
	target.ModuleArchives = archives
	return target, nil
}
