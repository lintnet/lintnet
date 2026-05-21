package github_test

import (
	"testing"

	"github.com/lintnet/lintnet/pkg/github"
)

func TestNew(t *testing.T) {
	t.Parallel()
	client, err := github.New(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("client must not be nil")
	}
}
