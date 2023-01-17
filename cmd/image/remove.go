package image

import (
	"fmt"
	"log"

	"github.com/ninroot/gocker/config"
	"github.com/ninroot/gocker/pkg/image"
	"github.com/ninroot/gocker/pkg/storage"
	"github.com/ninroot/gocker/pkg/util"
	"github.com/spf13/cobra"
)

var removeCommand = &cobra.Command{
	Use:   "rm IMAGE",
	Args:  cobra.ExactArgs(1),
	Short: "remove image",
	Run: func(cmd *cobra.Command, args []string) {
		img, err := image.Parse(args[0])
		if err != nil {
			log.Fatal(err)
		}
		store := storage.NewImageStore(util.EnsureDir(config.DefaultImageStoreRootDir))
		if err := store.RemoveImage(img.Digest); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Image deleted: ", img.Digest)
		}
	},
}
