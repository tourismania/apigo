// Package cli holds cobra command definitions. They share the same
// application use-cases as the HTTP layer — no duplicated orchestration.
package cli

import (
	"context"
	"errors"
	"fmt"

	createuser "api/internal/application/command/create_user"

	"github.com/spf13/cobra"
)

// NewCreateUserCommand returns the `create-user` cobra command.
//
// Usage:
//
//	app create-user <firstName> <lastName> <email> <password>
func NewCreateUserCommand(uc createuser.UseCase) *cobra.Command {
	return &cobra.Command{
		Use:   "create-user <firstName> <lastName> <email> <password>",
		Short: "Create a new user",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			firstName, lastName, email, password := args[0], args[1], args[2], args[3]
			if email == "" || password == "" {
				return errors.New("email and password are required")
			}
			res, err := uc.Handle(context.Background(), createuser.Command{
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Password:  password,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "User successfully generated! id=%d\n", res.ID)
			return nil
		},
	}
}
