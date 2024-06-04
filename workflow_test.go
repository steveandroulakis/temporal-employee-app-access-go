package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_SuccessfulEmployeeAppAccessWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	employeeInput := EmployeeInput{
		EmployeeName: "John Doe",
	}

	permissionInput := PermissionInput{
		EmployeeName:    "John Doe",
		ApplicationName: "Salesforce",
		Permission:      "Admin",
		Expiry:          0,
	}

	// Mock activity implementation
	env.OnActivity(GrantPermission, mock.Anything, permissionInput).Return(nil)
	env.OnActivity(RevokePermission, mock.Anything, permissionInput).Return(nil)

	env.ExecuteWorkflow(EmployeeAppAccessWorkflow, employeeInput)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func Test_GrantPermissionFailedWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	employeeInput := EmployeeInput{
		EmployeeName: "Jane Doe",
	}

	permissionInput := PermissionInput{
		EmployeeName:    "Jane Doe",
		ApplicationName: "Office365",
		Permission:      "User",
		Expiry:          0,
	}

	// Mock activity implementation
	env.OnActivity(GrantPermission, mock.Anything, permissionInput).Return(errors.New("unable to grant permission"))
	env.OnActivity(RevokePermission, mock.Anything, permissionInput).Return(nil)

	env.ExecuteWorkflow(EmployeeAppAccessWorkflow, employeeInput)
	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}
