// Package member implements community member management services
package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/kv"
)

const (
	Everybody = Group("everybody")
)

type User string
type Group string

var (
	usersNS = mod.NS("users")
	usersKV = kv.KV[User, git.URL]{}

	groupsNS = mod.NS("groups")
	groupsKV = kv.KV[Group, form.None]{}

	userGroupsNS  = mod.NS("user_groups")
	userGroupsKKV = kv.KKV[User, Group, bool]{}

	groupUsersNS  = mod.NS("group_users")
	groupUsersKKV = kv.KKV[Group, User, bool]{}
)

func AddMember(ctx context.Context, t *git.Tree, user User, group Group) mod.Change[form.None] {
	userGroupsKKV.Set(ctx, userGroupsNS, t, user, group, true)
	groupUsersKKV.Set(ctx, groupUsersNS, t, group, user, true)
	return mod.Change[form.None]{
		Msg: fmt.Sprintf("Added user %v to group %v", user, group),
	}
}

func IsMember(ctx context.Context, t *git.Tree, user User, group Group) bool {
	var userHasGroup, groupHasUser bool
	must.Try0(
		func() { userHasGroup = userGroupsKKV.Get(ctx, userGroupsNS, t, user, group) },
	)
	must.Try0(
		func() { groupHasUser = groupUsersKKV.Get(ctx, groupUsersNS, t, group, user) },
	)
	return userHasGroup && groupHasUser
}

func RemoveMember(ctx context.Context, t *git.Tree, user User, group Group) mod.Change[form.None] {
	userGroupsKKV.Remove(ctx, userGroupsNS, t, user, group)
	groupUsersKKV.Remove(ctx, groupUsersNS, t, group, user)
	return mod.Change[form.None]{
		Msg: fmt.Sprintf("Removed user %v from group %v", user, group),
	}
}

func ListUserGroups(ctx context.Context, t *git.Tree, user User) []Group {
	return userGroupsKKV.ListSecondaryKeys(ctx, userGroupsNS, t, user)
}

func ListGroupUsers(ctx context.Context, t *git.Tree, group Group) []User {
	return groupUsersKKV.ListSecondaryKeys(ctx, groupUsersNS, t, group)
}
