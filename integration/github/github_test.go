// Some parts of this file are adopted from github.com/google/go-github

//go:build integration
// +build integration

package github_test

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v55/github"
	govgh "github.com/gov4git/gov4git/github"
)

var (
	client *github.Client
	// auth indicates whether tests are being run with an OAuth token.
	// Tests can use this flag to skip certain tests when run without auth.
	auth bool
)

func init() {
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		print("no oauth token (that's ok)\n\n")
		client = github.NewClient(nil)
	} else {
		client = github.NewTokenClient(context.Background(), token)
		auth = true
	}
}

func checkAuth(name string) bool {
	if !auth {
		fmt.Printf("no auth, skipping portions of %v\n", name)
	}
	return auth
}

var (
	TestRepo = govgh.GithubRepo{Owner: "gov4git", Name: "testing.project"}
)
