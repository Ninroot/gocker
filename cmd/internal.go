package cmd

import (
	"os"

	"github.com/ninroot/gocker/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var internReq pkg.RunRequest

var internalCommand = &cobra.Command{
	Use:   "internal",
	Short: "Internal command for gocker itself",
	Long:  "Internal command for gocker itself. You probably don't want to use it, unless you know what you are doing.",
	Run: func(cmd *cobra.Command, args []string) {
		e, err := pkg.NewRuntimeService().InitContainer(internReq)
		if err != nil {
			logrus.WithError(err).Fatal("Internal: Failed to run container")
		}
		if e != nil {
			os.Exit(*e)
		}
	},
}

func init() {
	internalCommand.Flags().StringVarP(&internReq.ContainerName, "ContainerName", "", "", "")
	internalCommand.Flags().StringVarP(&internReq.ImageName, "ImageName", "", "", "")
	internalCommand.Flags().StringVarP(&internReq.ImageTag, "ImageTag", "", "", "")
	internalCommand.Flags().StringVarP(&internReq.ContainerCommand, "ContainerCommand", "", "", "")
	internalCommand.Flags().StringVarP(&internReq.ContainerID, "ContainerID", "", "", "")
}
