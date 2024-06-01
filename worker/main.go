package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/khuedoan/yeet"
)

func main() {
	// Create a Temporal Client
	// A Temporal Client is a heavyweight object that should be created just once per process.
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	// Create a new Worker.
	yeetWorker := worker.New(temporalClient, "yeet-task-queue", worker.Options{})

	// Workflow
	yeetWorker.RegisterWorkflow(yeet.YeetStandard)

	// Activitis
	message := "This could be a connection string or endpoint details"
	number := 100
	yeetWorker.RegisterActivity(&yeet.Build{
		Message: &message,
		Number:  &number,
	})
	yeetWorker.RegisterActivity(yeet.GitClone)

	err = yeetWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}
