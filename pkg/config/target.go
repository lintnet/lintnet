package config

import "strings"

type Target struct {
	ID             string                    `json:"id,omitempty"`
	LintFiles      []*LintGlob               `json:"lint_files,omitempty"`
	Modules        []*ModuleGlob             `json:"modules,omitempty"`
	ModuleArchives map[string]*ModuleArchive `json:"module_archives,omitempty"`
	DataFiles      []*DataFile               `json:"data_files,omitempty"`
}

type RawTarget struct {
	ID        string       `json:"id,omitempty"`
	LintGlobs []*LintGlob  `json:"lint_files"`
	Modules   []*RawModule `json:"-"`
	DataFiles []string     `json:"data_files"`
}

func (rt *RawTarget) Parse() (*Target, error) {
	modules := make([]*RawModule, 0, len(rt.LintGlobs))
	lintGlobs := make([]*LintGlob, 0, len(rt.LintGlobs))
	for _, lintGlob := range rt.LintGlobs {
		switch {
		case strings.HasPrefix(lintGlob.Glob, "!module:"):
			modules = append(modules, &RawModule{
				Glob:   "!" + strings.TrimPrefix(lintGlob.Glob, "!module:"),
				Files:  lintGlob.Files,
				Config: lintGlob.Config,
			})
		case strings.HasPrefix(lintGlob.Glob, "module:"):
			modules = append(modules, &RawModule{
				Glob:   strings.TrimPrefix(lintGlob.Glob, "module:"),
				Files:  lintGlob.Files,
				Config: lintGlob.Config,
			})
		default:
			lintGlob.Clean()
			lintGlobs = append(lintGlobs, lintGlob)
		}
	}
	rt.LintGlobs = lintGlobs
	rt.Modules = modules
	dataFiles := make([]*DataFile, len(rt.DataFiles))
	for i, dataFile := range rt.DataFiles {
		dataFiles[i] = NewDataFile(dataFile)
	}
	target := &Target{
		ID:        rt.ID,
		LintFiles: rt.LintGlobs,
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
