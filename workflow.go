package yeet

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type YeetStandardParam struct {
	WorkflowParamX string
	WorkflowParamY int
}

type YeetStandardResultObject struct {
	success bool
}

func YeetStandard(ctx workflow.Context, param YeetStandardParam) (*YeetStandardResultObject, error) {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	activityParam := BuildParam{
		ActivityParamX: param.WorkflowParamX,
		ActivityParamY: param.WorkflowParamY,
	}

	err := workflow.ExecuteActivity(ctx, GitClone, activityParam).Get(ctx, nil)

	var build *Build
	var activityResult BuildResult
	err = workflow.ExecuteActivity(ctx, build.Buildpacks, activityParam).Get(ctx, &activityResult)
	if err != nil {
		return nil, err
	}

	// TODO defer?
	err = workflow.ExecuteActivity(ctx, CleanUp, activityParam).Get(ctx, nil)

	workflowResult := &YeetStandardResultObject{
		success: true,
	}

	return workflowResult, nil
}
