package image

import (
	"fmt"
	"log"

	"github.com/ninroot/gocker/config"
	"github.com/ninroot/gocker/pkg"
	"github.com/spf13/cobra"
)

var removeCommand = &cobra.Command{
	Use:   "rm IMAGE",
	Args:  cobra.ExactArgs(1),
	Short: "remove image",
	Run: func(cmd *cobra.Command, args []string) {
		img, err := pkg.Parse(args[0])
		if err != nil {
			log.Fatal(err)
		}
		store := pkg.NewImageStore(pkg.EnsureDir(config.DefaultImageStoreRootDir))
		if del, err := store.RemoveImage(&img); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Image deleted: ", del.Digest)
		}
	},
}
