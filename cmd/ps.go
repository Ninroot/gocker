package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/ninroot/gocker/pkg"
	"github.com/spf13/cobra"
)

var psCommand = &cobra.Command{
	Use:   "ps",
	Args:  cobra.ExactArgs(0),
	Short: "List containers",
	Run: func(cmd *cobra.Command, args []string) {
		runtime := pkg.NewRuntimeService()
		conts, err := runtime.ListContainers()
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "CONTAINER ID\tIMAGE\tCREATED\tNAME\n")
		for _, c := range *conts {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", c.ID, c.Image.Name, c.CreatedAt.Format(time.RFC3339), c.Name)
		}
		w.Flush()
	},
}
