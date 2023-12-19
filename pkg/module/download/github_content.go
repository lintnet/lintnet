package download

import (
	"context"

	"github.com/suzuki-shunsuke/lintnet/pkg/github"
)

type GitHubContent struct {
	api GitHubContentAPI
}

type GitHubContentAPI interface {
	GetContents(ctx context.Context, repoOwner, repoName, path string, opt *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error)
}

func (g *GitHubContent) Download(ctx context.Context) (string, []string, error) {
	return "", nil, nil
}
