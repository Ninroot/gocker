package image

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ninroot/gocker/config"
	"github.com/ninroot/gocker/pkg"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Args:  cobra.ExactArgs(0),
	Short: "list images",
	Run: func(cmd *cobra.Command, args []string) {
		store := pkg.NewImageStore(pkg.EnsureDir(config.DefaultImageStoreRootDir))
		images, err := store.ListImages()
		if err != nil {
			log.Fatal(err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "NAME\tTAG\tDIGEST\t\n")
		for _, img := range images {
			fmt.Fprintf(w, "%s\t%s\t%s\n", img.Name, img.Tag, img.Digest)
		}
		w.Flush()
	},
}