package image

import (
	"github.com/spf13/cobra"
)

func init() {
	Command.AddCommand(listCommand)
	Command.AddCommand(removeCommand)
}

var Command = &cobra.Command{
	Use:   "image COMMAND",
	Short: "Manage image",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
