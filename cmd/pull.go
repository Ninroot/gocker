package cmd

import (
	"log"

	"github.com/ninroot/gocker/config"
	"github.com/ninroot/gocker/pkg"
	"github.com/ninroot/gocker/pkg/storage"
	"github.com/ninroot/gocker/pkg/util"
	"github.com/spf13/cobra"
)

var pullCommand = &cobra.Command{
	Use:   "pull IMAGE",
	Args:  cobra.ExactArgs(1),
	Short: "Pull container image",
	Run: func(cmd *cobra.Command, args []string) {
		regSvc := pkg.NewRegistryService(
			storage.NewImageStore(util.EnsureDir(config.DefaultImageStoreRootDir)),
		)
		img := args[0]
		if err := regSvc.Pull(img); err != nil {
			log.Println("Failed to pull the image ", img, ": ", err)
		}
	},
}
