//go:build integration
// +build integration

package github_test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-github/v55/github"
	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/testutil"
	"github.com/gov4git/lib4git/util"
)

func TestGithubIssueStructure(t *testing.T) {
	issues, _, err := client.Issues.ListByRepo(context.Background(), TestRepo.Owner, TestRepo.Name, &github.IssueListByRepoOptions{State: "all"})
	if err != nil {
		t.Fatalf("Issues.ListByRepo returned error: %v", err)
	}
	issuesByNumber(issues).Sort()
	// fmt.Println("ISSUES", form.SprintJSON(issues))

	if len(issues) < 4 {
		t.Fatalf("Expected at least 4 issues, got %v", len(issues))
	}

	// test issue 1
	if issues[0].GetNumber() != 1 {
		t.Fatalf("Expected issue number 1, got %v", issues[0].GetNumber())
	}
	if issues[0].GetState() != "open" {
		t.Fatalf("Expected issue state 'open', got %v", issues[0].GetState())
	}
	if !util.IsIn(govgh.PrioritizeIssueByGovernanceLabel, govgh.LabelsToStrings(issues[0].Labels)...) {
		t.Fatalf("Expected issue to be prioritized")
	}

	// test issue 2
	if issues[1].GetNumber() != 2 {
		t.Fatalf("Expected issue number 2, got %v", issues[1].GetNumber())
	}
	if !issues[1].GetLocked() {
		t.Fatalf("Expected issue to be locked")
	}

	// test issue 3
	if issues[2].GetNumber() != 3 {
		t.Fatalf("Expected issue number 3, got %v", issues[2].GetNumber())
	}
	if issues[2].GetState() != "closed" {
		t.Fatalf("Expected issue state 'open', got %v", issues[2].GetState())
	}

	// test issue 5 (pull request)
	if issues[4].GetNumber() != 5 {
		t.Fatalf("Expected issue number 5, got %v", issues[4].GetNumber())
	}
	if issues[4].GetState() != "open" {
		t.Fatalf("Expected issue state 'open', got %v", issues[4].GetState())
	}
	if issues[4].GetPullRequestLinks() == nil {
		t.Fatalf("Expected issue to be a pull request")
	}

	fmt.Println(form.SprintJSON(issues))
}

type issuesByNumber []*github.Issue

func (x issuesByNumber) Sort() {
	sort.Sort(x)
}

func (x issuesByNumber) Len() int {
	return len(x)
}

func (x issuesByNumber) Less(i, j int) bool {
	return x[i].GetNumber() < x[j].GetNumber()
}

func (x issuesByNumber) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func TestImportIssuesForPrioritization(t *testing.T) {

	ghRepo := TestRepo
	ghClient := client

	// init governance
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	// import issues
	govgh.ImportIssuesForPrioritization(ctx, ghRepo, ghClient, cty.Organizer())

	// list ballots
	ads1 := ballot.List(ctx, cty.Gov())
	if len(ads1) < 4 {
		t.Errorf("expecting at least 3, got %v,\n%v", len(ads1), form.SprintJSON(ads1))
	}
	// issue-1: open, not-frozen
	if ads1[0].Name.Path() != "github/issues/1" {
		t.Errorf("expecting github/issues/1, got %v", ads1[0].Name.Path())
	}
	if ads1[0].Closed || ads1[0].Frozen {
		t.Errorf("expecting open, not-frozen")
	}
	// issue-2: open, frozen
	if ads1[1].Name.Path() != "github/issues/2" {
		t.Errorf("expecting github/issues/2, got %v", ads1[1].Name.Path())
	}
	if ads1[1].Closed || !ads1[1].Frozen {
		t.Errorf("expecting open, frozen")
	}
	// issue-3: closed, frozen
	if ads1[2].Name.Path() != "github/issues/3" {
		t.Errorf("expecting github/issues/3, got %v", ads1[2].Name.Path())
	}
	if !ads1[2].Closed || !ads1[2].Frozen {
		t.Errorf("expecting closed, frozen")
	}
}
