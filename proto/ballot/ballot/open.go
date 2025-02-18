package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Open(
	ctx context.Context,
	strat common.Strategy,
	govAddr gov.GovAddress,
	name ns.NS,
	title string,
	description string,
	choices []string,
	participants member.Group,
) git.Change[form.Map, common.BallotAddress] {

	govCloned := git.CloneOne(ctx, git.Address(govAddr))
	chg := Open_StageOnly(ctx, strat, govAddr, govCloned, name, title, description, choices, participants)
	proto.Commit(ctx, govCloned.Tree(), chg)
	govCloned.Push(ctx)
	return chg
}

func Open_StageOnly(
	ctx context.Context,
	strat common.Strategy,
	govAddr gov.GovAddress,
	govCloned git.Cloned,
	name ns.NS,
	title string,
	description string,
	choices []string,
	participants member.Group,
) git.Change[form.Map, common.BallotAddress] {

	// check no open ballots by the same name
	openAdNS := common.BallotPath(name).Append(common.AdFilebase)
	if _, err := git.TreeStat(ctx, govCloned.Tree(), openAdNS); err == nil {
		must.Errorf(ctx, "ballot already exists: %v", openAdNS.GitPath())
	}

	// verify group exists
	if !member.IsGroup_Local(ctx, govCloned.Tree(), participants) {
		must.Errorf(ctx, "participant group %v does not exist", participants)
	}

	// write ad
	ad := common.Advertisement{
		Gov:          govAddr,
		Name:         name,
		Title:        title,
		Description:  description,
		Choices:      choices,
		Strategy:     strat.Name(),
		Participants: participants,
		Frozen:       false,
		Closed:       false,
		Cancelled:    false,
		ParentCommit: git.Head(ctx, govCloned.Repo()),
	}
	git.ToFileStage(ctx, govCloned.Tree(), openAdNS, ad)

	// write initial tally
	tally := common.Tally{
		Ad:            ad,
		Scores:        map[string]float64{},
		VotesByUser:   map[member.User]map[string]common.StrengthAndScore{},
		AcceptedVotes: map[member.User]common.AcceptedElections{},
		RejectedVotes: map[member.User]common.RejectedElections{},
		Charges:       map[member.User]float64{},
	}
	openTallyNS := common.BallotPath(name).Append(common.TallyFilebase)
	git.ToFileStage(ctx, govCloned.Tree(), openTallyNS, tally)

	// write strategy
	openStratNS := common.BallotPath(name).Append(common.StrategyFilebase)
	git.ToFileStage(ctx, govCloned.Tree(), openStratNS, strat)

	return git.NewChange(
		fmt.Sprintf("Create ballot of type %v", strat.Name()),
		"ballot_open",
		form.Map{
			"strategy":     strat,
			"name":         name,
			"participants": participants,
		},
		common.BallotAddress{Gov: govAddr, Name: name},
		nil,
	)
}
