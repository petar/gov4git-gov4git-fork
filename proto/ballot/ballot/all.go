package ballot

import (
	"context"
	"fmt"
	"sync"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func TallyAll(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	maxPar int,
) git.Change[form.Map, []common.Tally] {

	base.Infof("fetching and tallying community votes ...")

	govOwner := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := TallyAll_StageOnly(ctx, govAddr, govOwner, maxPar)
	if len(chg.Result) == 0 {
		return chg
	}
	proto.Commit(ctx, govOwner.Public.Tree(), chg)
	govOwner.Public.Push(ctx)
	return chg
}

func TallyAll_StageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govOwner id.OwnerCloned,
	maxPar int,
) git.Change[form.Map, []common.Tally] {

	// list all open ballots
	communityTree := govOwner.Public.Tree()
	ads := common.FilterOpenClosedAds(false, List_Local(ctx, communityTree))

	// compute union of all voter accounts from all open ballots
	participatingVoters := make([]participatingVoters, len(ads))
	allVoters := map[member.User]member.Account{}
	for i, ad := range ads {
		participatingVoters[i] = *loadParticipatingVoters(ctx, communityTree, ad)
		for user, acct := range participatingVoters[i].VoterAccounts {
			allVoters[user] = acct
		}
	}

	// fetch repos of all participating users
	allVoterClones := clonePar(ctx, allVoters, maxPar)

	// populate participating voter clones
	for _, pv := range participatingVoters {
		pv.attachVoterClones(ctx, allVoterClones)
	}

	// perform tallies for all open ballots
	tallyChanges := []git.Change[map[string]form.Form, common.Tally]{}
	tallies := []common.Tally{}
	for _, pv := range participatingVoters {
		if tallyChg, changed := tallyVotersCloned_StageOnly(ctx, govAddr, govOwner, pv.Ad.Name, pv.VoterAccounts, pv.VoterClones); changed {
			tallyChanges = append(tallyChanges, tallyChg)
			tallies = append(tallies, tallyChg.Result)
		}
	}

	return git.NewChange(
		fmt.Sprintf("Tallied votes on all ballots"),
		"ballot_tally_all",
		form.Map{},
		tallies,
		form.ToForms(tallyChanges),
	)
}

func clonePar(ctx context.Context, userAccounts map[member.User]member.Account, maxPar int) map[member.User]git.Cloned {

	must.Assertf(ctx, maxPar > 0, "clone parallelism must be greater than zero")

	var wg sync.WaitGroup
	wg.Add(len(userAccounts))

	sem := make(chan bool, maxPar)

	var allLock sync.Mutex
	allClones := map[member.User]git.Cloned{}

	for u, a := range userAccounts {
		sem <- true
		go func(u member.User, a member.Account) {

			base.Infof("cloning voter %v repository %v", u, a.PublicAddress)
			cloned, err := git.TryCloneOne(ctx, git.Address(a.PublicAddress))
			if err != nil {
				base.Infof("user %v repository %v unresponsive (%v)", u, a.PublicAddress, err)
			} else {
				base.Infof("user %v repository %v cloned successfully (%v)", u, a.PublicAddress, err)
				allLock.Lock()
				allClones[u] = cloned
				allLock.Unlock()
			}

			<-sem
			wg.Done()
		}(u, a)
	}

	wg.Wait()

	return allClones
}
