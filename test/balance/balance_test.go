package balance

import (
	"testing"

	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/testutil"
)

func TestBalance(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	bal := balance.Balance{"test_balance"}

	// test set/get roundtrip
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), bal, 30.0)
	actual1 := balance.Get(ctx, cty.Gov(), cty.MemberUser(0), bal)
	if actual1 != 30.0 {
		t.Errorf("expecting %v, got %v", 30.0, actual1)
	}

	// test balance transfer
	cloned := gov.Clone(ctx, cty.Gov())
	balance.Transfer_StageOnly(ctx, cloned.Tree(), cty.MemberUser(0), bal, cty.MemberUser(1), bal, 10.0)
	git.Commit(ctx, cloned.Tree(), "test commit")
	cloned.Push(ctx)
	actual2 := balance.Get(ctx, cty.Gov(), cty.MemberUser(1), bal)
	if actual2 != 10.0 {
		t.Errorf("expecting %v, got %v", 10.0, actual2)
	}

	// test cannot set for non-existent user
	err := must.Try(
		func() {
			balance.Set(ctx, cty.Gov(), cty.NonExistentMemberUser(), bal, 30.0)
		},
	)
	if err == nil {
		t.Errorf("must fail")
	}

	// test cannot get for non-existent user
	err = must.Try(
		func() {
			balance.Get(ctx, cty.Gov(), cty.NonExistentMemberUser(), bal)
		},
	)
	if err == nil {
		t.Errorf("must fail")
	}
}
