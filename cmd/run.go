package cmd

import (
	"github.com/ninroot/gocker/pkg"
	"github.com/spf13/cobra"
)

var runCommand = &cobra.Command{
	Use:   "run IMAGE",
	Short: "Run a command in a new container",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Run(args)
	},
}
