package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pullCommand)
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
