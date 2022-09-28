package services

import (
	"context"
	"path/filepath"

	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type GovUserRemoveIn struct {
	Name            string `json:"name"`             // community unique handle for this user
	CommunityBranch string `json:"community_branch"` // branch in community repo where user will be added
}

type GovUserRemoveOut struct{}

func (x GovUserRemoveOut) Human() string {
	return ""
}

func (x GovService) UserRemove(ctx context.Context, in *GovUserRemoveIn) (*GovUserRemoveOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := GovRmUser(ctx, community, in.Name); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovUserRemoveOut{}, nil
}

func GovRmUser(ctx context.Context, community git.Local, name string) error {
	userFile := filepath.Join(proto.GovUsersDir, name, proto.GovUserInfoFilebase)
	// remove user file
	if err := community.Dir().Remove(userFile); err != nil {
		return err
	}
	// stage changes
	if err := community.Remove(ctx, []string{userFile}); err != nil {
		return err
	}
	// commit changes
	if err := community.Commitf(ctx, "gov: remove user %v", name); err != nil {
		return err
	}
	return nil
}
