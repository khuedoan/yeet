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

	var gitops *Git
	gitopsParam := GitParam{
		Host:       param.Host,
		Owner:      "khuedoan",
		Repository: "timoni-gitops-test",
		Revision:   "master",
	}
	var gitopsResult GitResult
	err = workflow.ExecuteActivity(ctx, gitops.Clone, gitopsParam).Get(ctx, &gitopsResult)
	if err != nil {
		return nil, err
	}

	var deploy *Deploy
	deployParam := DeployParam{
		RepoPath: gitopsResult.Path,
		SubPath:  fmt.Sprintf("%s/%s", param.Owner, param.Repository),
	}
	var deployConfig DeployConfig
	err = workflow.ExecuteActivity(ctx, deploy.GetConfig, deployParam).Get(ctx, &deployConfig)
	if err != nil {
		return nil, err
	}

	for _, stage := range deployConfig[param.Revision].Stages {
		workflow.ExecuteActivity(ctx, deploy.ProcessStage, stage).Get(ctx, nil)
		if stage.Wait != "" {
			wait, err := time.ParseDuration(stage.Wait)
			if err != nil {
				return nil, err
			}
			workflow.Sleep(ctx, wait)
		}
	}

	// TODO defer?
	_ = workflow.ExecuteActivity(ctx, git.Clean).Get(ctx, nil)
	_ = workflow.ExecuteActivity(ctx, gitops.Clean).Get(ctx, nil)

	workflowResult := &YeetStandardResultObject{
		success: true,
	}

	return workflowResult, nil
}
