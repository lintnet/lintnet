package github

import (
	"context"
	"errors"
	"net/url"
)

var errGetTar = errors.New("failed to get tar")

type MockRepositoriesService struct {
	URL *url.URL
}

func (m *MockRepositoriesService) GetArchiveLink(ctx context.Context, owner, repo string, archiveformat ArchiveFormat, opts *RepositoryContentGetOptions, maxRedirects int) (*url.URL, *Response, error) {
	if m.URL == nil {
		return nil, nil, errGetTar
	}
	return m.URL, nil, nil
}
