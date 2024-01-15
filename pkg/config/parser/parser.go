package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/errlevel"
)

func Parse(rc *config.RawConfig) (*config.Config, error) { //nolint:cyclop,funlen
	cfg := &config.Config{
		ErrorLevel:      errlevel.Error,
		ShownErrorLevel: errlevel.Info,
		Targets:         make([]*config.Target, len(rc.Targets)),
		Outputs:         rc.Outputs,
		IgnoredPatterns: getIgnoredPatterns(rc.IgnoredDirs),
	}

	if cfg.IgnoredPatterns == nil {
		cfg.IgnoredPatterns = []string{
			"node_modules",
			".git",
		}
	}

	if rc.ErrorLevel != "" {
		level, err := errlevel.New(rc.ErrorLevel)
		if err != nil {
			return nil, fmt.Errorf("parse the error level: %w", err)
		}
		cfg.ErrorLevel = level
	}

	if rc.ShownErrorLevel != "" {
		level, err := errlevel.New(rc.ShownErrorLevel)
		if err != nil {
			return nil, fmt.Errorf("parse the error level: %w", err)
		}
		cfg.ShownErrorLevel = level
	}

	if cfg.ShownErrorLevel > cfg.ErrorLevel {
		cfg.ShownErrorLevel = cfg.ErrorLevel
	}

	moduleArchives := map[string]*config.ModuleArchive{}
	for i, rt := range rc.Targets {
		target, err := ParseTarget(rt)
		if err != nil {
			return nil, err
		}
		cfg.Targets[i] = target
		for k, ma := range target.ModuleArchives {
			moduleArchives[k] = ma
		}
	}
	for _, output := range rc.Outputs {
		if strings.HasPrefix(output.Template, "github_archive/github.com/") {
			m, err := ParseImport(output.Template)
			if err != nil {
				return nil, fmt.Errorf("parse a module path: %w", err)
			}
			output.TemplateModule = m
			moduleArchives[m.Archive.String()] = m.Archive
		}
		if strings.HasPrefix(output.Transform, "github_archive/github.com/") {
			m, err := ParseImport(output.Transform)
			if err != nil {
				return nil, fmt.Errorf("parse a module path: %w", err)
			}
			output.TransformModule = m
			moduleArchives[m.Archive.String()] = m.Archive
		}
	}
	cfg.ModuleArchives = moduleArchives
	return cfg, nil
}

func getIgnoredPatterns(ignoredDirs []string) []string {
	if ignoredDirs == nil {
		ignoredDirs = []string{
			".git",
			"node_modules",
		}
	}
	ignoredPatterns := make([]string, len(ignoredDirs))
	for i, d := range ignoredDirs {
		ignoredPatterns[i] = fmt.Sprintf("**/%s/**", d)
	}
	return ignoredPatterns
}

func ParseImport(line string) (*config.Module, error) {
	mg, err := ParseModuleLine(line)
	if err != nil {
		return nil, err
	}
	return &config.Module{
		ID:        mg.ID,
		Archive:   mg.Archive,
		SlashPath: mg.SlashPath,
	}, nil
}

func ParseModuleLine(line string) (*config.ModuleGlob, error) {
	// <type>/github.com/<repo owner>/<repo name>/<path>@<commit hash>[:<tag>]
	line = strings.TrimSpace(line)
	excluded := false
	if l := strings.TrimPrefix(line, "!"); l != line {
		excluded = true
		line = strings.TrimSpace(l)
	}
	elems := strings.Split(line, "/")
	if len(elems) < 5 { //nolint:gomnd
		return nil, errors.New("line is invalid")
	}
	if elems[0] != "github_archive" {
		return nil, errors.New("unsupported module type")
	}
	if elems[1] != "github.com" {
		return nil, errors.New("module host must be 'github.com'")
	}
	pathAndRefAndTag := strings.Join(elems[4:], "/")
	path, refAndTag, ok := strings.Cut(pathAndRefAndTag, "@")
	if !ok {
		return nil, errors.New("ref is required")
	}
	ref, tag, _ := strings.Cut(refAndTag, ":")
	if err := validateRef(ref); err != nil {
		return nil, err
	}
	return &config.ModuleGlob{
		ID:        line,
		SlashPath: strings.Join(append(elems[:4], ref, path), "/"),
		Archive: &config.ModuleArchive{
			Type:      "github_archive",
			Host:      "github.com",
			RepoOwner: elems[2],
			RepoName:  elems[3],
			Ref:       ref,
			Tag:       tag,
		},
		Excluded: excluded,
	}, nil
}

var fullCommitHashPattern = regexp.MustCompile("[a-fA-F0-9]{40}")

func validateRef(ref string) error {
	if fullCommitHashPattern.MatchString(ref) {
		return nil
	}
	return errors.New("ref must be full commit hash")
}

func ParseModule(rm *config.RawModule) (*config.ModuleGlob, error) {
	m, err := ParseModuleLine(rm.Glob)
	if err != nil {
		return nil, fmt.Errorf("parse a module path: %w", err)
	}
	m.Config = rm.Config
	return m, nil
}

func ParseTarget(rt *config.RawTarget) (*config.Target, error) {
	lintFiles := make([]*config.ModuleGlob, len(rt.LintGlobs))
	for i, lintGlob := range rt.LintGlobs {
		lintFiles[i] = lintGlob.ToModule()
	}
	target := &config.Target{
		ID:        rt.ID,
		LintFiles: lintFiles,
		Modules:   make([]*config.ModuleGlob, len(rt.Modules)),
		DataFiles: rt.DataFiles,
	}
	archives := make(map[string]*config.ModuleArchive, len(rt.Modules))
	for i, m := range rt.Modules {
		a, err := ParseModule(m)
		if err != nil {
			return nil, err
		}
		target.Modules[i] = a
		archives[a.Archive.String()] = a.Archive
	}
	target.ModuleArchives = archives
	return target, nil
}
