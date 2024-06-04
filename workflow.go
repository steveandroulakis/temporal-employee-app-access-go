package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type ApplicationPermission struct {
	ApplicationName string
	Permission      string
	Expiry          time.Duration
	Status          string
}

type EmployeeInput struct {
	EmployeeName string
}

func EmployeeAppAccessWorkflow(ctx workflow.Context, input EmployeeInput) error {
	// Workflow options
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// Block waiting for signals
	workflow.GetSignalChannel(ctx, "grantPermission").Receive(ctx, nil)

	return nil
}
