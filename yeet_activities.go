package yeet

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

func GitClone(ctx context.Context) error {
	time.Sleep(9 * time.Second)
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
