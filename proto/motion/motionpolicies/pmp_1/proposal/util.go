package proposal

import (
	"context"

	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func SaveMotionPolicyState_StageOnly(ctx context.Context, t *git.Tree, policyNS ns.NS, state *ProposalState) {
	git.ToFileStage[*ProposalState](ctx, t, policyNS.Append(StateFilebase), state)
}

func LoadMotionPolicyState_Local(ctx context.Context, t *git.Tree, policyNS ns.NS) *ProposalState {
	state := git.FromFile[ProposalState](ctx, t, policyNS.Append(StateFilebase))
	return &state
}
