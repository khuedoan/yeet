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

	var workflowParams yeet.YeetStandardParam
	if err := json.NewDecoder(r.Body).Decode(&workflowParams); err != nil {
        http.Error(w, "Bad request: unable to decode JSON", http.StatusBadRequest)
        return
    }

	// TODO any better way to do this?
	if workflowParams.Repository == "" || workflowParams.Revision == "" {
        http.Error(w, "Bad request: missing required fields", http.StatusBadRequest)
        return
    }

	workflowOptions := client.StartWorkflowOptions{
		// TODO deterministic workflow ID?
		// ID: "pipeline/githost/owner/repo/name",
		TaskQueue: "yeet",
	}

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
