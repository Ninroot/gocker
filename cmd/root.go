package cmd

import (
	"github.com/ninroot/gocker/cmd/image"
	"github.com/ninroot/gocker/config"
	"github.com/ninroot/gocker/pkg/logging"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var verbosity int

func init() {
	rootCmd.AddCommand(image.Command)
	rootCmd.AddCommand(psCommand)
	rootCmd.AddCommand(pullCommand)
	rootCmd.AddCommand(rmCommand)
	rootCmd.AddCommand(runCommand)
	rootCmd.AddCommand(internalCommand)

	rootCmd.PersistentFlags().IntVarP(&verbosity, "verbose", "v", int(config.DefaultLogLevel), "0 to 6: Trace, Debug, Info, Warn, Error, Fatal, Panic")
	logrus.SetLevel(logging.VerbosityToLogrusLevel(verbosity))
}

var rootCmd = &cobra.Command{
	Use:   "gocker",
	Short: "gocker - dockerlike",
	Long:  "Gocker is a Dockerlike project for educational purposes.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
