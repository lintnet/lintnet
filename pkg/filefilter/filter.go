package filefilter

import (
	"path/filepath"

	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/filefind"
)

type Param struct {
	DataRootDir string   `json:"data_root_dir,omitempty"`
	TargetID    string   `json:"target_id,omitempty"`
	FilePaths   []string `json:"file_paths,omitempty"`
	PWD         string   `json:"pwd,omitempty"`
}

func FilterTargetsByFilePaths(param *Param, targets []*filefind.Target) []*filefind.Target {
	for i, filePath := range param.FilePaths {
		if filepath.IsAbs(filePath) {
			continue
		}
		param.FilePaths[i] = filepath.Join(param.PWD, filePath)
	}
	if param.TargetID == "" {
		return filterTargets(targets, param.FilePaths)
	}
	arr := make([]*domain.Path, len(param.FilePaths))
	for i, filePath := range param.FilePaths {
		p := &domain.Path{
			Abs: filePath,
			Raw: filePath,
		}
		if !filepath.IsAbs(filePath) {
			p.Abs = filepath.Join(param.PWD, filePath)
		}
		arr[i] = p
	}
	targets[0].DataFiles = arr
	return targets
}

func filterTargets(targets []*filefind.Target, filePaths []string) []*filefind.Target {
	newTargets := make([]*filefind.Target, 0, len(targets))
	for _, target := range targets {
		newTarget := filterTarget(target, filePaths)
		if len(newTarget.LintFiles) > 0 {
			newTargets = append(newTargets, newTarget)
		}
	}
	return newTargets
}

func filterTarget(target *filefind.Target, filePaths []string) *filefind.Target {
	newTarget := &filefind.Target{}
	for _, lintFile := range target.LintFiles {
		for _, filePath := range filePaths {
			if lintFile.Path == filePath {
				newTarget.LintFiles = append(newTarget.LintFiles, lintFile)
				break
			}
		}
	}
	lintChanged := false
	if len(newTarget.LintFiles) > 0 {
		newTarget.DataFiles = target.DataFiles
		lintChanged = true
	}
	dataChanged := false
	for _, dataFile := range target.DataFiles {
		for _, filePath := range filePaths {
			if dataFile.Abs == filePath {
				dataChanged = true
				if !lintChanged {
					newTarget.DataFiles = append(newTarget.DataFiles, dataFile)
				}
			}
		}
	}
	if dataChanged {
		newTarget.LintFiles = target.LintFiles
	}
	return newTarget
}
