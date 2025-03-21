package bureau

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/mail"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Process(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	group member.Group,
) git.ChangeNoResult {

	base.Infof("fetching service requests from the community ...")

	govOwner := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg, changed := Process_StageOnly(ctx, govAddr, govOwner, group)
	if changed {
		proto.Commit(ctx, govOwner.Public.Tree(), chg)
		govOwner.Public.Push(ctx)
	}
	return chg
}

func Process_StageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govOwner id.OwnerCloned,
	group member.Group,
) (change git.ChangeNoResult, changed bool) {

	communityTree := govOwner.Public.Tree()

	// list participating users
	users := member.ListGroupUsers_Local(ctx, communityTree, group)

	// get user accounts
	accounts := make([]member.Account, len(users))
	for i, user := range users {
		accounts[i] = member.GetUser_Local(ctx, communityTree, user)
	}

	// fetch user requests
	var fetchedReqs FetchedRequests
	for i, account := range accounts {
		if fetched, err := fetchUserRequests(ctx, govAddr, govOwner, users[i], account); err != nil {
			base.Infof("fetching bureau requests for user %v (%v)", users[i], err)
		} else {
			fetchedReqs = append(fetchedReqs, fetched.Result...)
		}
	}

	// process requests
	for _, fetched := range fetchedReqs {
		nOK, nErr := processRequest_StageOnly(ctx, govAddr, govOwner, fetched)
		if nOK+nErr > 0 {
			changed = true
		}
	}

	return git.NewChangeNoResult(
		fmt.Sprintf("Process bureau requests of users in group %v", group),
		"bureau_process",
	), changed
}

func processRequest_StageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govOwner id.OwnerCloned,
	fetched FetchedRequest,
) (numOK int, numErr int) {
	for _, req := range fetched.Requests {
		if req.Transfer == nil {
			numErr++
			continue
		}
		if req.Transfer.FromUser != fetched.User {
			base.Infof("bureau: invalid transfer request from user %v; origin of transfer is not the requesting user", fetched.User)
			numErr++
			continue
		}
		err := must.Try(func() {
			balance.Transfer_StageOnly(
				ctx,
				govOwner.Public.Tree(),
				req.Transfer.FromUser, req.Transfer.FromBalance,
				req.Transfer.ToUser, req.Transfer.ToBalance,
				req.Transfer.Amount,
			)
		})
		if err != nil {
			base.Infof("bureau: transfer error (%v)", err)
			numErr++
			continue
		}
		numOK++
		base.Infof("bureau: transferred %v from %v:%v to %v:%v",
			req.Transfer.Amount,
			req.Transfer.FromUser, req.Transfer.FromBalance,
			req.Transfer.ToUser, req.Transfer.ToBalance,
		)
	}
	return
}

func fetchUserRequests(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govOwner id.OwnerCloned,
	user member.User,
	account member.Account,
) (git.Change[form.Map, FetchedRequests], error) {

	fetched := FetchedRequests{}
	var respond mail.Responder[Request, Request] = func(
		ctx context.Context,
		_ mail.SeqNo,
		req Request,
	) (resp Request, err error) {
		fetched = append(fetched,
			FetchedRequest{
				User:     user,
				Address:  account.PublicAddress,
				Requests: Requests{req},
			})
		return req, nil
	}

	userPublic, err := git.TryCloneOne(ctx, git.Address(account.PublicAddress))
	if err != nil {
		return git.Change[form.Map, FetchedRequests]{}, err
	}

	recvOnly := mail.Respond_StageOnly[Request, Request](
		ctx,
		govOwner,
		account.PublicAddress,
		userPublic.Tree(),
		BureauTopic,
		respond,
	)

	return git.NewChange(
		fmt.Sprintf("Fetched requests from user %v", user),
		"bureau_fetch_user_requests",
		form.Map{"user": user, "account": account},
		fetched,
		form.Forms{recvOnly},
	), nil
}
