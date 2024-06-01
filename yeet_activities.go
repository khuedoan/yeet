package yeet

import (
	"context"

	"go.temporal.io/sdk/activity"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

func GitClone(ctx context.Context) error {
	_, _ = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/khuedoan/example-service",
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
