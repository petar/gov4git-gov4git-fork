package cron

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v55/github"
	govgh "github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

var CronNS = ns.NS{"cron", "cron.json"}

func Cron(
	ctx context.Context,
	repo govgh.GithubRepo,
	ghc *github.Client,
	govAddr gov.OrganizerAddress,
	//
	githubFreq time.Duration, // frequency of importing from github
	communityFreq time.Duration, // frequency of fetching community votes and service requests
	//
	maxPar int, // parallelism for fetching community votes
) form.Map {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	govTree := govCloned.Public.Tree()

	// use a separate branch for cron logs
	cronAddr := git.Address(govAddr.Public)
	cronAddr.Branch = cronAddr.Branch + ".cron"
	cronCloned := git.CloneOne(ctx, cronAddr)
	cronTree := cronCloned.Tree()

	// read cron state
	state, err := git.TryFromFile[CronState](ctx, cronTree, CronNS)
	must.Assertf(ctx, err == nil || err == os.ErrNotExist, "opening cron state (%v)", err)

	now := time.Now()
	shouldSyncGithub := now.Sub(state.LastGithubImport) > githubFreq
	shouldSyncCommunity := now.Sub(state.LastCommunityTally) > communityFreq

	report := form.Map{}

	// import from github
	if shouldSyncGithub {

		// fetch repo maintainers
		maintainers := govgh.FetchRepoMaintainers(ctx, repo, ghc)
		base.Infof("maintainers for %v are %v", repo, form.SprintJSON(maintainers))

		// process issues and pull requests
		report["processed_issues"] = govgh.ImportIssuesForPrioritization_StageOnly(ctx, repo, ghc, govAddr, govCloned)

		// process joins
		report["processed_joins"] = govgh.ProcessJoinRequestIssues_StageOnly(ctx, repo, ghc, govAddr, govCloned, maintainers)

		// process directives
		report["processed_directives"] = govgh.ProcessDirectiveIssues_StageOnly(ctx, repo, ghc, govAddr, govCloned, maintainers)

		state.LastGithubImport = time.Now()
	}

	// sync community
	if shouldSyncCommunity {

		// tally votes for all ballots from all community members
		report["tally"] = ballot.TallyAll_StageOnly(ctx, govAddr, govCloned, maxPar).Result

		state.LastCommunityTally = time.Now()
	}

	report["cron"] = state
	cronChg := git.NewChange[form.Map, form.Map](
		fmt.Sprintf("Cron job."),
		"cron",
		form.Map{"time": now},
		report,
		nil,
	)

	// push gov state
	govStatus, err := govTree.Status()
	must.NoError(ctx, err)
	if !govStatus.IsClean() {
		proto.Commit(ctx, govCloned.Public.Tree(), cronChg)
		govCloned.Public.Push(ctx)
	}

	// push cron state
	git.ToFileStage(ctx, cronTree, CronNS, state)
	proto.Commit(ctx, cronTree, cronChg)
	cronCloned.Push(ctx)

	return report
}

type CronState struct {
	LastGithubImport   time.Time `json:"last_github_import"`
	LastCommunityTally time.Time `json:"last_community_tally"`
}
