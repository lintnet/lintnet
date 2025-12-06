package lintfile

import (
	"strings"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type Parser struct {
	fs afero.Fs
}

func NewParser(fs afero.Fs) *Parser {
	return &Parser{
		fs: fs,
	}
}

func (p *Parser) Parse(lintFile *config.LintFile) (*domain.Node, error) {
	node, err := jsonnet.ReadToNode(p.fs, lintFile.Path)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return &domain.Node{
		Node:    node,
		Key:     lintFile.ID,
		Config:  lintFile.Config,
		Link:    lintFile.Link,
		Combine: strings.HasSuffix(lintFile.Path, "_combine.jsonnet"),
	}, nil
}

func (p *Parser) Parses(lintFiles []*config.LintFile) ([]*domain.Node, error) {
	nodes := make([]*domain.Node, 0, len(lintFiles))
	for _, lintFile := range lintFiles {
		node, err := p.Parse(lintFile)
		if err != nil {
			return nil, slogerr.With(err, "file_path", lintFile.Path) //nolint:wrapcheck
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
