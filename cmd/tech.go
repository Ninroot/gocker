package cmd

import (
	"github.com/ninroot/gocker/pkg"
	"github.com/spf13/cobra"
)

var techCommand = &cobra.Command{
	Use:   "tech",
	Short: "technical subcommand used by gocker itself",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Chroot(args)
	},
}
