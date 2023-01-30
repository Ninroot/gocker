package image

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ninroot/gocker/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "ls",
	Args:  cobra.ExactArgs(0),
	Short: "List images",
	Run: func(cmd *cobra.Command, args []string) {
		run := pkg.NewRuntimeService()
		images, err := run.ListImages()
		if err != nil {
			logrus.Fatal(err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "NAME\tTAG\tDIGEST\t\n")
		for _, img := range *images {
			fmt.Fprintf(w, "%s\t%s\t%s\n", img.Name, img.Tag, img.Digest)
		}
		w.Flush()
	},
}
