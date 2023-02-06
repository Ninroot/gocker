package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ninroot/gocker/pkg"
	"github.com/sirupsen/logrus"
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
			logrus.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "CONTAINER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tNAME\n")
		for _, c := range *conts {
			cmdFmt := fmt.Sprintf("%s %s", c.Command, strings.Join(c.Args, " "))
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", c.ID, c.Image.Name, cmdFmt, c.CreatedAt.Format(time.RFC3339), status(c.ExitCode), c.Name)
		}
		w.Flush()
	},
}

func status(s *int) string {
	if s == nil {
		return "Running"
	}
	return fmt.Sprintf("Exited (%d)", *s)
}
