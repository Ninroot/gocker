package cmd

import (
	"log"

	"github.com/ninroot/gocker/cmd/image"
	"github.com/ninroot/gocker/config"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var logLevel string

func init() {
	rootCmd.AddCommand(image.Command)
	rootCmd.AddCommand(psCommand)
	rootCmd.AddCommand(pullCommand)
	rootCmd.AddCommand(rmCommand)
	rootCmd.AddCommand(runCommand)
	rootCmd.AddCommand(internalCommand)

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", config.DefaultLogLevel,
		"Set the logging level (\"trace\"|\"debug\"|\"info\"|\"warn\"|\"error\"|\"fatal\"|\"panic\")")
}

var rootCmd = &cobra.Command{
	Use:   "gocker",
	Short: "gocker - dockerlike",
	Long:  "Gocker is a Dockerlike project for educational purposes.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		l, err := logrus.ParseLevel(logLevel)
		if err != nil {
			log.Fatal("Failed to set logger: ", err)
		}
		logrus.SetLevel(l)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
