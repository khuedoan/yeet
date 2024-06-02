package yeet

import (
	"context"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"go.temporal.io/sdk/activity"
)

type Git struct {
	path string
}

type GitParam struct {
	Repository string
	Revision   string
}

type GitResult struct {
	Path string
}

func (a *Git) Clone(ctx context.Context, param GitParam) (*GitResult, error) {
    tempDir, err := os.MkdirTemp("", "yeet-")
	if err != nil {
        return nil, err
	}

	a.path = tempDir
	_, err = git.PlainClone(a.path, false, &git.CloneOptions{
		URL:           param.Repository,
		ReferenceName: plumbing.ReferenceName(param.Revision),
		Depth:         1,
	})
	if err != nil {
		return nil, err
	}

	result := &GitResult{
		Path: a.path,
	}
	return result, nil
}

func (a *Git) Clean(ctx context.Context, param GitParam) error {
	os.RemoveAll(a.path)

	return nil
}

type BuildParam struct {
	ActivityParamX string
	ActivityParamY string
}

type BuildResult struct {
	ResultFieldX string
	ResultFieldY string
}

type Build struct {
	Message *string
	Number  *string
}

func (a *Build) Buildpacks(ctx context.Context, param BuildParam) (*BuildResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("The message is:", param.ActivityParamX)
	logger.Info("The number is:", param.ActivityParamY)

	result := &BuildResult{
		ResultFieldX: "Success",
		ResultFieldY: "1",
	}

	return result, nil
}
