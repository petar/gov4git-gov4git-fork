package arb

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/gov/arb"
	"github.com/gov4git/gov4git/testutil"
)

func TestVote(t *testing.T) {
	// base.LogVerbosely()

	// create test community
	// dir := testutil.MakeStickyTestDir()
	dir := t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 1)
	if err != nil {
		t.Fatal(err)
	}
	ctx := testCommunity.Background()

	// create poll
	arbService := arb.GovArbService{
		GovConfig:      testCommunity.CommunityGovConfig(),
		IdentityConfig: testCommunity.UserIdentityConfig(0),
	}
	pollOut, err := arbService.Poll(ctx,
		&arb.PollIn{
			Path:            "test_poll",
			Choices:         []string{"a", "b", "c"},
			Group:           "participants",
			Strategy:        "prioritize",
			GoverningBranch: proto.MainBranch,
		})
	if err != nil {
		t.Fatal(err)
	}

	// cast a vote
	voteOut, err := arbService.Vote(ctx,
		&arb.VoteIn{
			ReferendumBranch: pollOut.PollBranch,
			ReferendumPath:   "test_poll",
			VoteChoice:       "a",
			VoteStrength:     1.0,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("poll: %v\nvote: %v\n", form.Pretty(pollOut), form.Pretty(voteOut))
}
