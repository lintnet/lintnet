package module

type Module struct {
	ID     string
	Source string
	// http, github_content, github_archive
	Type      string
	RepoOwner string
	RepoName  string
	Path      string
	Ref       string
}
