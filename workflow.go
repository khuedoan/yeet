package yeet

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type YeetStandardParam struct {
	Repository string
	Revision   string
}

type YeetStandardResultObject struct {
	success bool
}

func YeetStandard(ctx workflow.Context, param YeetStandardParam) (*YeetStandardResultObject, error) {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var git *Git
	gitParam := GitParam{
		Repository: param.Repository,
		Revision:   param.Revision,
	}
	var gitResult GitResult
	err := workflow.ExecuteActivity(ctx, git.Clone, gitParam).Get(ctx, &gitResult)

	buildParam := BuildParam{
		ActivityParamX: param.Repository,
		ActivityParamY: param.Revision,
	}
	var build *Build
	var activityResult BuildResult
	err = workflow.ExecuteActivity(ctx, build.Buildpacks, buildParam).Get(ctx, &activityResult)
	if err != nil {
		return nil, err
	}

	// TODO defer?
	err = workflow.ExecuteActivity(ctx, git.Clean).Get(ctx, nil)

	workflowResult := &YeetStandardResultObject{
		success: true,
	}

	return workflowResult, nil
}
