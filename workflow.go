package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// Holds the status of an application permission
type ApplicationPermission struct {
	ApplicationName string
	Permission      string
	Expiry          time.Duration
	Status          string
}

// Workflow input (employee name)
type EmployeeInput struct {
	EmployeeName string
}

// Temporal Signal for granting permission
type GrantSignal struct {
	ApplicationName string
	Permission      string
	Expiry          int64
}

// Temporal Signal for revoking permission
type RevokeSignal struct {
	ApplicationName string
}

// For recording signal events in a history
type SignalEvent struct {
	Type       string
	Permission ApplicationPermission
}

// AwaitSignals is a structure to hold the state of signal handling
type AwaitSignals struct {
	Permissions      map[string]ApplicationPermission
	Timers           map[string]workflow.CancelFunc
	GrantSignalChan  workflow.ReceiveChannel
	RevokeSignalChan workflow.ReceiveChannel
	SignalQueue      workflow.Channel
	History          []SignalEvent
}

// Listen listens to grant and revoke signals
func (a *AwaitSignals) Listen(ctx workflow.Context) {
	log := workflow.GetLogger(ctx)
	for {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(a.GrantSignalChan, func(c workflow.ReceiveChannel, more bool) {
			var signal GrantSignal
			c.Receive(ctx, &signal)
			permission := ApplicationPermission{
				ApplicationName: signal.ApplicationName,
				Permission:      signal.Permission,
				Expiry:          time.Duration(signal.Expiry) * time.Second,
				Status:          "active",
			}
			a.SignalQueue.Send(ctx, SignalEvent{Type: "grant", Permission: permission})
			log.Info("Grant signal received", "ApplicationName", signal.ApplicationName, "Permission", signal.Permission, "Expiry", signal.Expiry)
		})
		selector.AddReceive(a.RevokeSignalChan, func(c workflow.ReceiveChannel, more bool) {
			var signal RevokeSignal
			c.Receive(ctx, &signal)
			permission, ok := a.Permissions[signal.ApplicationName]
			if ok {
				permission.Status = "expired"
				a.SignalQueue.Send(ctx, SignalEvent{Type: "revoke", Permission: permission})
				log.Info("Revoke signal received", "ApplicationName", signal.ApplicationName)
			}
		})
		selector.Select(ctx)
	}
}

// EmployeeAppAccessWorkflow is the main workflow function
func EmployeeAppAccessWorkflow(ctx workflow.Context, input EmployeeInput) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	log := workflow.GetLogger(ctx)
	log.Debug("EmployeeAppAccessWorkflow started", "EmployeeName", input.EmployeeName)

	// Initialize the signal handling structure
	signals := &AwaitSignals{
		Permissions:      make(map[string]ApplicationPermission),
		Timers:           make(map[string]workflow.CancelFunc),
		GrantSignalChan:  workflow.GetSignalChannel(ctx, "grantPermission"),
		RevokeSignalChan: workflow.GetSignalChannel(ctx, "revokePermission"),
		SignalQueue:      workflow.NewBufferedChannel(ctx, 100),
		History:          make([]SignalEvent, 0),
	}

	// Start listening to signals in a separate goroutine
	workflow.Go(ctx, signals.Listen)

	// Setup query handlers
	err := workflow.SetQueryHandler(ctx, "GetAccess", func() (map[string]ApplicationPermission, error) {
		return signals.Permissions, nil
	})
	if err != nil {
		log.Error("SetQueryHandler for GetAccess failed", "Error", err)
		return err
	}

	err = workflow.SetQueryHandler(ctx, "GetHistory", func() ([]SignalEvent, error) {
		return signals.History, nil
	})
	if err != nil {
		log.Error("SetQueryHandler for GetHistory failed", "Error", err)
		return err
	}

	// Main workflow loop
	for {
		// Wait for a signal
		var signalEvent SignalEvent
		signals.SignalQueue.Receive(ctx, &signalEvent)

		signals.History = append(signals.History, signalEvent)

		// Get information from the signal
		permission := signalEvent.Permission
		app := permission.ApplicationName

		// Based on the signal type, grant or revoke access
		if signalEvent.Type == "grant" {
			// If a timer already exists for this application, cancel it without changing the permission status
			if cancel, exists := signals.Timers[app]; exists {
				cancel()
			}

			signals.Permissions[app] = permission
			err := workflow.ExecuteActivity(ctx, GrantAccess, input.EmployeeName, app, permission.Permission, int64(permission.Expiry.Seconds())).Get(ctx, nil)
			if err != nil {
				return err
			}

			// If the permission has an expiry, set a timer to expire it
			if permission.Expiry > 0 {
				expiryTimerCtx, cancel := workflow.WithCancel(ctx)
				signals.Timers[app] = cancel
				workflow.Go(ctx, func(ctx workflow.Context) {
					_ = workflow.NewTimer(expiryTimerCtx, permission.Expiry).Get(ctx, nil)
					signals.SignalQueue.Send(ctx, SignalEvent{
						Type: "expire",
						Permission: ApplicationPermission{
							ApplicationName: app,
							Permission:      permission.Permission,
							Expiry:          0,
							Status:          "expired",
						},
					})
				})
			}
		} else if signalEvent.Type == "revoke" {
			// If a timer exists for this application, cancel it
			if cancel, exists := signals.Timers[app]; exists {
				cancel()
				delete(signals.Timers, app)
			}

			err := workflow.ExecuteActivity(ctx, RevokeAccess, input.EmployeeName, app).Get(ctx, nil)
			if err != nil {
				return err
			}
			delete(signals.Permissions, app)
		} else if signalEvent.Type == "expire" {
			permission.Status = "expired"
			signals.Permissions[app] = permission
			log.Info("[Permission expired]", "ApplicationName", app)
			err := workflow.ExecuteActivity(ctx, RevokeAccess, input.EmployeeName, app).Get(ctx, nil)
			if err != nil {
				return err
			}
		}
	}
}
