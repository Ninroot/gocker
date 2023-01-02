package cmd

import (
	"github.com/ninroot/gocker/pkg"
	"github.com/spf13/cobra"
)

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "run a command in a new container",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Run(args)
	},
}
