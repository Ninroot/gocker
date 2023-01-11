package image

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(listCommand)
}

var Command = &cobra.Command{
	Use:   "image COMMAND",
	Short: "Manage image",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing or unknown command")
		cmd.Help()
	},
}
