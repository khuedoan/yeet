package yeet

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type YeetStandardParam struct {
	Host       string
	Owner      string
	Repository string
	Revision   string
}

type YeetStandardResultObject struct {
	success bool
}

func YeetStandard(ctx workflow.Context, param YeetStandardParam) (*YeetStandardResultObject, error) {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var git *Git
	gitParam := GitParam{
		Host:       param.Host,
		Owner:      param.Owner,
		Repository: param.Repository,
		Revision:   param.Revision,
	}
	var gitResult GitResult
	err := workflow.ExecuteActivity(ctx, git.Clone, gitParam).Get(ctx, &gitResult)
	if err != nil {
		return nil, err
	}

	buildParam := BuildParam{
		Path:  gitResult.Path,
		Image: fmt.Sprintf("%s/%s", param.Owner, param.Repository),
		Tag:   param.Revision,
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
