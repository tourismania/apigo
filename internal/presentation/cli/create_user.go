// Package cli holds cobra command definitions. They reuse the same
// application bus the HTTP layer uses — there's no duplicated
// orchestration.
package cli

import (
	"context"
	"errors"
	"fmt"

	"api/internal/application/bus"
	createusercmd "api/internal/application/command/create_user"

	"github.com/spf13/cobra"
)

// NewCreateUserCommand returns the `create-user` cobra command.
//
// Usage:
//
//	app create-user <firstName> <lastName> <email> <password>
func NewCreateUserCommand(b bus.CommandBus) *cobra.Command {
	return &cobra.Command{
		Use:   "create-user <firstName> <lastName> <email> <password>",
		Short: "Create a new user via the command bus",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			firstName, lastName, email, password := args[0], args[1], args[2], args[3]
			if email == "" || password == "" {
				return errors.New("email and password are required")
			}
			raw, err := b.Dispatch(context.Background(), createusercmd.Command{
				FirstName: firstName,
				LastName:  lastName,
				Email:     email,
				Password:  password,
			})
			if err != nil {
				return err
			}
			res, ok := raw.(createusercmd.Result)
			if !ok {
				return fmt.Errorf("unexpected handler result %T", raw)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "User successfully generated! id=%d\n", res.ID)
			return nil
		},
	}
}
