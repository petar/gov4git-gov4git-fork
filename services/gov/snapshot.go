package gov

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type SnapshotBranchLatestIn struct {
	SourceRepo   string    `json:"source_repo"`   // source repo url
	SourceBranch string    `json:"source_branch"` // source branch
	Community    git.Local `json:"community"`     // community local clone
}

type SnapshotBranchLatestOut struct {
	In           *SnapshotBranchLatestIn
	SourceCommit string `json:"source_commit"` // latest commit found at source
}

func (x SnapshotBranchLatestOut) Human(context.Context) string {
	return fmt.Sprintf("Took a snapshot of %v on branch %v at commit %v, into local clone of community at %v.",
		x.In.SourceRepo, x.In.SourceBranch, x.SourceCommit, x.In.Community.Path)
}

// SnapshotBranchLatest downloads the latest commit on a given branch at a remote source repo, and
// places it into a local community repo.
// SnapshotBranchLatest stages but does not commit the changes made to the local community repo.
func (x GovService) SnapshotBranchLatest(ctx context.Context, in *SnapshotBranchLatestIn) (*SnapshotBranchLatestOut, error) {
	// clone source repo locally at the branch
	source, err := git.MakeLocalInCtx(ctx, "source")
	if err != nil {
		return nil, err
	}
	if err := source.CloneBranch(ctx, in.SourceRepo, in.SourceBranch); err != nil {
		return nil, err
	}

	// get current commit
	latestCommit, err := source.HeadCommitHash(ctx)
	if err != nil {
		return nil, err
	}

	// directory inside community where snapshot lives
	srcPath := proto.SnapshotDir(in.SourceRepo, latestCommit) //XXX: when SourceRepo has directory separators???
	srcParent, _ := filepath.Split(srcPath)

	// if the community repo already has a snapshot of the source commit, remove it.
	if err := in.Community.Dir().RemoveAll(srcPath); err != nil {
		return nil, err
	}

	// prepare the parent directory of the snapshot
	if err := in.Community.Dir().Mkdir(srcParent); err != nil {
		return nil, err
	}

	// remove the .git directory from the source clone
	if err := source.Dir().RemoveAll(".git"); err != nil {
		return nil, err
	}

	// move the source clone to the community clone
	if err := files.Rename(source.Dir(), in.Community.Dir().Subdir(srcPath)); err != nil {
		return nil, err
	}

	// stage changes in community clone
	if err := in.Community.Add(ctx, []string{srcPath}); err != nil {
		return nil, err
	}

	return &SnapshotBranchLatestOut{In: in, SourceCommit: latestCommit}, nil
}
