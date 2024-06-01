package yeet

import (
	"context"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"go.temporal.io/sdk/activity"
)

func GitClone(ctx context.Context) error {
	repo, err := git.PlainClone("/tmp/yeet/example-service", false, &git.CloneOptions{
		URL: "https://github.com/khuedoan/example-service",
		Depth: 1,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	ref, err := repo.Head()
	if err != nil {
		return err
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return err
	}

	// Print commit list
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})

	return nil
}

type BuildParam struct {
	ActivityParamX string
	ActivityParamY int
}

type BuildResult struct {
	ResultFieldX string
	ResultFieldY int
}

type Build struct {
	Message *string
	Number  *int
}

func (a *Build) Buildpacks(ctx context.Context, param BuildParam) (*BuildResult, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("The message is:", param.ActivityParamX)
	logger.Info("The number is:", param.ActivityParamY)

	result := &BuildResult{
		ResultFieldX: "Success",
		ResultFieldY: 1,
	}

	return result, nil
}
