package cmd

import (
	"log"

	"github.com/ninroot/gocker/pkg"
	"github.com/spf13/cobra"
)

var techCommand = &cobra.Command{
	Use:   "tech",
	Short: "technical subcommand used by gocker itself",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Go tech with args:", args)
		runtime := pkg.NewRuntimeService()
		if err := runtime.InitContainer(args); err != nil {
			log.Fatal(err)
		}
	},
}
