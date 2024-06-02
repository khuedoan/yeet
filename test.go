package yeet

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	wfParam := YeetStandardParam{
		Repository: "Hello World!",
		Revision:   "100",
	}
	activityResult := BuildResult{
		ResultFieldX: "Message",
		ResultFieldY: "1",
	}
	var activities *Build
	env.OnActivity(activities.Buildpacks, mock.Anything, mock.Anything).Return(&activityResult, nil)
	env.ExecuteWorkflow(YeetStandard, wfParam)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
	var result YeetStandardResultObject
	require.NoError(t, env.GetWorkflowResult(&result))
}

func Test_Activity(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()
	activityParam := BuildParam{
		ActivityParamX: "Message",
		ActivityParamY: "1",
	}
	var activities Build
	message := "No messages!"
	counter := "0"
	activities.Message = &message
	activities.Number = &counter
	env.RegisterActivity(activities.Buildpacks)
	val, err := env.ExecuteActivity(activities.Buildpacks, activityParam)
	require.NoError(t, err)
	var res BuildResult
	require.NoError(t, val.Get(&res))
	require.Equal(t, "Success", res.ResultFieldX)
}
