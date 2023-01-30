package cmd

import (
	"github.com/ninroot/gocker/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rmCommand = &cobra.Command{
	Use:   "rm CONTAINER",
	Args:  cobra.ExactArgs(1),
	Short: "Remove a container",
	Run: func(cmd *cobra.Command, args []string) {
		runtime := pkg.NewRuntimeService()
		err := runtime.RemoveContainer(args[0])
		if err != nil {
			logrus.Fatal(err)
		}
	},
}
