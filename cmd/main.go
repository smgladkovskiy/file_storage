package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Start http server",
		RunE:  run,
	}
)

func Execute() {
	runCmd.Flags().IntVarP(
		&storagesAmount,
		"storagesAmount",
		"s",
		5,
		"Number of storages",
	)

	if err := runCmd.MarkFlagRequired("storagesAmount"); err != nil {
		return
	}

	if err := runCmd.ExecuteContext(context.Background()); err != nil {
		return
	}
}
