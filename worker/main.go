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
	yeetWorker := worker.New(temporalClient, "yeet", worker.Options{})

	// Workflow
	yeetWorker.RegisterWorkflow(yeet.YeetStandard)

	// Activitis
	yeetWorker.RegisterActivity(&yeet.Build{})
	yeetWorker.RegisterActivity(&yeet.Git{})

	err = yeetWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}
