package main

import (
	"context"
	"flag"
	"log"
	"os"

	"temporal-employee-app-access-go/app"

	"go.temporal.io/sdk/client"
)

func main() {
	set := flag.NewFlagSet("signal", flag.ExitOnError)
	employeeID := set.String("employeeID", "", "Employee ID (workflow ID)")
	action := set.String("action", "", "Action to perform (grant/revoke)")
	applicationName := set.String("application", "", "Name of the application")
	permission := set.String("permission", "", "Permission to grant (required for grant action)")
	expiry := set.Int("expiry", 0, "Expiry time in seconds (required for grant action)")

	clientOptions, err := app.ParseClientOptionFlags(set, os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	set.Parse(os.Args[1:])

	var signalName string
	var signalData interface{}

	switch *action {
	case "grant":
		signalName = "grantPermission"
		signalData = app.GrantSignal{
			ApplicationName: *applicationName,
			Permission:      *permission,
			Expiry:          int64(*expiry),
		}
	case "revoke":
		signalName = "revokePermission"
		signalData = app.RevokeSignal{
			ApplicationName: *applicationName,
		}
	default:
		log.Fatalf("Invalid action: %s", *action)
	}

	err = c.SignalWorkflow(context.Background(), *employeeID, "", signalName, signalData)
	if err != nil {
		log.Fatalf("Error sending the signal: %v", err)
	}

	log.Println("Signal sent successfully")
}
