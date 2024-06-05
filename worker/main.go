package main

import (
	"flag"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"temporal-employee-app-access-go/app"
)

// @@@SNIPSTART temporal-employee-app-access-go-worker
func main() {
	set := flag.NewFlagSet("worker", flag.ExitOnError)
	clientOptions, err := app.ParseClientOptionFlags(set, os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, app.EmployeeAccessTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(app.EmployeeAppAccessWorkflow)
	w.RegisterActivity(app.GrantAccess)
	w.RegisterActivity(app.RevokeAccess)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
