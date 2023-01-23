package cmd

import (
	"log"

	"github.com/ninroot/gocker/cmd/input"
	"github.com/ninroot/gocker/pkg"
	"github.com/spf13/cobra"
)

var techCommand = &cobra.Command{
	Use:   "tech",
	Short: "Technical subcommand used by gocker itself",
	Long:  "Technical subcommand used by gocker itself. You probably don't want to use it, unless you know what you are doing.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imageName, imageTag := input.Parse(args[0])

		conCmd := []string{}
		if len(args) >= 1 {
			conCmd = args[1:]
		}

		runtime := pkg.NewRuntimeService()
		if err := runtime.InitContainer(imageName, imageTag, containerName, conCmd); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	techCommand.Flags().StringVarP(&containerName, "name", "", "", "Assign a name to the container")
}
