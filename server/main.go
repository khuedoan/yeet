package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/khuedoan/yeet"

	"go.temporal.io/sdk/client"
)

func main() {
	temporalClient, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal Client", err)
	}
	defer temporalClient.Close()

	http.HandleFunc("/v1/yeet", func(w http.ResponseWriter, r *http.Request) {
		startWorkflowHandler(w, r, temporalClient)
	})
	err = http.ListenAndServe(":8091", nil)
	if err != nil {
		log.Fatalln("Unable to run HTTP server", err)
	}
}

func startWorkflowHandler(w http.ResponseWriter, r *http.Request, temporalClient client.Client) {
	// Set the options for the Workflow Execution.
	// A Task Queue must be specified.
	// A custom Workflow Id is highly recommended.
	workflowOptions := client.StartWorkflowOptions{
		//      pipeline/ git host / owner  / repo          / name
		ID:        "yeet/github.com/khuedoan/example-service/master",
		TaskQueue: "yeet-task-queue",
	}
	workflowParams := yeet.YeetStandardParam{
		WorkflowParamX: "Hello World!",
		WorkflowParamY: 999,
	}
	// Call the Temporal Cluster to start the Workflow Execution.
	workflowExecution, err := temporalClient.ExecuteWorkflow(
		context.Background(),
		workflowOptions,
		yeet.YeetStandard,
		workflowParams,
	)
	if err != nil {
		log.Fatalln("Unable to execute the Workflow", err)
	}
	log.Println("Started Workflow!")
	log.Println("WorkflowID:", workflowExecution.GetID())
	log.Println("RunID:", workflowExecution.GetRunID())
	var result yeet.YeetStandardResultObject
	workflowExecution.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}
	b, err := json.Marshal(result)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println(string(b))
}
