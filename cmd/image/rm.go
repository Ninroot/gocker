package image

import (
	"fmt"
	"log"

	"github.com/ninroot/gocker/pkg"
	"github.com/ninroot/gocker/pkg/image"
	"github.com/spf13/cobra"
)

var removeCommand = &cobra.Command{
	Use:   "rm IMAGE",
	Args:  cobra.ExactArgs(1),
	Short: "Remove image",
	Run: func(cmd *cobra.Command, args []string) {
		img, err := image.Parse(args[0])
		if err != nil {
			log.Fatal(err)
		}
		runtime := pkg.NewRuntimeService()
		if err := runtime.RemoveImage(img.Name, img.Tag); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("Image <%s:%s> deleted\n", img.Name, img.Tag)
		}
	},
}
