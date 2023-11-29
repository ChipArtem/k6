package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ChipArtem/k6/api/v1/client"
	"github.com/ChipArtem/k6/cmd/state"
)

func getCmdStatus(gs *state.GlobalState) *cobra.Command {
	// statusCmd represents the status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show test status",
		Long: `Show test status.

  Use the global --address flag to specify the URL to the API server.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(gs.Flags.Address)
			if err != nil {
				return err
			}
			status, err := c.Status(gs.Ctx)
			if err != nil {
				return err
			}

			return yamlPrint(gs.Stdout, status)
		},
	}
	return statusCmd
}