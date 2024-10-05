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

func (m *MockRepositoriesService) GetArchiveLink(_ context.Context, _, _ string, _ ArchiveFormat, _ *RepositoryContentGetOptions, _ int) (*url.URL, *Response, error) {
	if m.URL == nil {
		return nil, nil, errGetTar
	}
	return m.URL, nil, nil
}
