package github

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v74/github"
	"golang.org/x/oauth2"
)

type (
	ArchiveFormat               = github.ArchiveFormat
	ListOptions                 = github.ListOptions
	ReleaseAsset                = github.ReleaseAsset
	RepositoriesService         = github.RepositoriesService
	Repository                  = github.Repository
	RepositoryContentGetOptions = github.RepositoryContentGetOptions
	RepositoryRelease           = github.RepositoryRelease
	RepositoryTag               = github.RepositoryTag
	Response                    = github.Response
)

const Tarball = github.Tarball

func New(ctx context.Context) *RepositoriesService {
	return github.NewClient(getHTTPClientForGitHub(ctx, getGitHubToken())).Repositories
}

func getGitHubToken() string {
	if token := os.Getenv("LINTNET_GITHUB_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITHUB_TOKEN")
}

func getHTTPClientForGitHub(ctx context.Context, token string) *http.Client {
	if token == "" {
		return http.DefaultClient
	}
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
}
