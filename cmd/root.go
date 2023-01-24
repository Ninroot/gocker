package cmd

import (
	"log"

	"github.com/ninroot/gocker/cmd/image"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(image.Command)
	rootCmd.AddCommand(psCommand)
	rootCmd.AddCommand(pullCommand)
	rootCmd.AddCommand(rmCommand)
	rootCmd.AddCommand(runCommand)
	rootCmd.AddCommand(internalCommand)
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
		log.Fatal(err)
	}
}
