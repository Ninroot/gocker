package cmd

import (
	"log"

	"github.com/ninroot/gocker/cmd/input"
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
			storage.NewImageStore(util.EnsureDir(config.DefaultImageStoreRootDir), storage.Btrfs{}),
		)
		name, tag := input.Parse(args[0])
		if err := regSvc.Pull(name, tag); err != nil {
			log.Fatalf("Failed to pull the image: %s\n", err)
		}
	},
}
