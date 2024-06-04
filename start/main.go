package main

import (
	"context"
	"log"
	"math/rand"
	"os"

	"go.temporal.io/sdk/client"

	"temporal-employee-app-access-go/app"
)

func generateRandomNumber(n int) string {
	var letterRunes = []rune("1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	clientOptions, err := app.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	employeeInput := app.EmployeeInput{
		EmployeeName: "John Doe",
	}

	// make workflow ID "employee" + a 6 digit random number
	workflowID := "employee" + generateRandomNumber(6)

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: app.EmployeeAccessTaskQueueName,
	}

	log.Printf("Starting workflow for employee: %s", employeeInput.EmployeeName)

	we, err := c.ExecuteWorkflow(context.Background(), options, app.EmployeeAppAccessWorkflow, employeeInput)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result string

	err = we.Get(context.Background(), &result)

	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}

	log.Println(result)
}
