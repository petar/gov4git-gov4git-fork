package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Unfreeze(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	ballotName ns.NS,
) git.ChangeNoResult {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := Unfreeze_StageOnly(ctx, govAddr, govCloned, ballotName)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Unfreeze_StageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	ballotName ns.NS,
) git.ChangeNoResult {

	govTree := govCloned.Public.Tree()

	// verify ad and strategy are present
	ad, _ := load.LoadStrategy(ctx, govTree, ballotName)

	must.Assertf(ctx, !ad.Closed, "ballot is closed")
	must.Assertf(ctx, ad.Frozen, "ballot is not frozen")

	ad.Frozen = false

	// write updated ad
	adNS := common.BallotPath(ballotName).Append(common.AdFilebase)
	git.ToFileStage(ctx, govTree, adNS, ad)

	return git.NewChangeNoResult(fmt.Sprintf("Unfreeze ballot %v", ballotName), "ballot_unfreeze")
}
