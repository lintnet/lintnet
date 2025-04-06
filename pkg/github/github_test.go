package github_test

import (
	"testing"

	"github.com/lintnet/lintnet/pkg/github"
)

func TestNew(t *testing.T) {
	t.Parallel()
	if client := github.New(t.Context()); client == nil {
		t.Fatal("client must not be nil")
	}
}
