package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/guregu/null.v3"

	v1 "github.com/ChipArtem/k6/api/v1"
	"github.com/ChipArtem/k6/api/v1/client"
	"github.com/ChipArtem/k6/cmd/state"
)

func getCmdPause(gs *state.GlobalState) *cobra.Command {
	// pauseCmd represents the pause command
	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause a running test",
		Long: `Pause a running test.

  Use the global --address flag to specify the URL to the API server.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(gs.Flags.Address)
			if err != nil {
				return err
			}
			status, err := c.SetStatus(gs.Ctx, v1.Status{
				Paused: null.BoolFrom(true),
			})
			if err != nil {
				return err
			}
			return yamlPrint(gs.Stdout, status)
		},
	}
	return pauseCmd
}
