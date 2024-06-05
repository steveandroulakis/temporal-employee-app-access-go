package app

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

type PermissionInput struct {
	EmployeeName    string
	ApplicationName string
	Permission      string
	Expiry          time.Duration
}

func GrantAccess(ctx context.Context, employeeName string, applicationName string, permission string, expiry int64) error {
	// Log statement to indicate activity execution
	activity.GetLogger(ctx).Info("Granting access", "EmployeeName", employeeName,
		"ApplicationName", applicationName, "Permission", permission, "Expiry", expiry)

	// Simulate work by sleeping
	time.Sleep(2 * time.Second)

	// Stubbed out implementation
	return nil
}

func RevokeAccess(ctx context.Context, employeeName string, applicationName string) error {
	// Log statement to indicate activity execution
	activity.GetLogger(ctx).Info("Revoking access", "EmployeeName", employeeName,
		"ApplicationName", applicationName)

	// Simulate work by sleeping
	time.Sleep(2 * time.Second)

	// Stubbed out implementation
	return nil
}
