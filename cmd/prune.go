package cmd

import (
	"github.com/ninroot/gocker/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pruneCommand = &cobra.Command{
	Use:   "prune",
	Args:  cobra.ExactArgs(0),
	Short: "Remove unused data",
	Run: func(cmd *cobra.Command, args []string) {
		r := pkg.NewRuntimeService()
		if err := r.Prune(); err != nil {
			logrus.WithError(err).Fatal("Failed to prune")
		}
	},
}
