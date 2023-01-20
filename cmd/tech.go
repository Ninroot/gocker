package cmd

import (
	"log"

	"github.com/ninroot/gocker/pkg"
	"github.com/spf13/cobra"
)

var techCommand = &cobra.Command{
	Use:   "tech",
	Short: "Technical subcommand used by gocker itself",
	Long:  "Technical subcommand used by gocker itself. You probably don't want to use it, unless you know what you are doing.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Go tech with args:", args)
		runtime := pkg.NewRuntimeService()
		if err := runtime.InitContainer(args); err != nil {
			log.Fatal(err)
		}
	},
}
