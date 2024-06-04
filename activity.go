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

func GrantPermission(ctx context.Context, input PermissionInput) error {
	// Log statement to indicate activity execution
	activity.GetLogger(ctx).Info("Granting permission", "EmployeeName",
		input.EmployeeName, "ApplicationName", input.ApplicationName, "Permission",
		input.Permission, "Expiry", input.Expiry)

	// Simulate work by sleeping
	time.Sleep(2 * time.Second)

	// Stubbed out implementation
	return nil
}

func RevokePermission(ctx context.Context, input PermissionInput) error {
	// Log statement to indicate activity execution
	activity.GetLogger(ctx).Info("Revoking permission", "EmployeeName",
		input.EmployeeName, "ApplicationName", input.ApplicationName,
		"Permission", input.Permission)

	// Simulate work by sleeping
	time.Sleep(2 * time.Second)

	// Stubbed out implementation
	return nil
}
