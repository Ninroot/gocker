package cmd

import (
	"fmt"
	"log"

	"github.com/ninroot/gocker/cmd/image"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pullCommand)
	rootCmd.AddCommand(runCommand)
	rootCmd.AddCommand(rmCommand)
	rootCmd.AddCommand(techCommand)
	rootCmd.AddCommand(image.Command)
}

var rootCmd = &cobra.Command{
	Use:   "gocker",
	Short: "gocker - dockerlike",
	Long:  "gocker is a Dockerlike project for educational purposes.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to gocker")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
