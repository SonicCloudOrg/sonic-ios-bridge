package devmode

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

var devmodeRevealCmd = &cobra.Command{
	Use:   "reveal",
	Short: "Reveal the Developer Mode menu on the device",
	Long:  "Reveal the Developer Mode menu on the device",
	RunE: func(cmd *cobra.Command, args []string) error {
		if bCan, eCan := canToggleDevMode(udid); eCan != nil {
			strErrMsg := fmt.Sprintf("Failed to check device %s iOS version", udid)
			logrus.Warn(strErrMsg)
			return xerrors.New(strErrMsg)
		} else if !bCan {
			strErrMsg := fmt.Sprintf("Device %s iOS version below 16", udid)
			logrus.Warn(strErrMsg)
			return xerrors.New(strErrMsg)
		}
		amfi, errAmfi := getAmfiServer()
		if errAmfi != nil {
			return errAmfi
		}
		res, errReveal := amfi.DevModeReveal()
		if errReveal != nil {
			return errReveal
		}
		if res == http.StatusOK {
			logrus.Infof("Developer Mode menu revealed successfully.")
			return nil
		} else {
			strErrMsg := fmt.Sprintf("Failed to reveal Developer Mode menu (%d).", res)
			logrus.Warn(strErrMsg)
			return xerrors.New(strErrMsg)
		}
	},
}

func initDevModeRevealCmd() {
	devmodeRootCMD.AddCommand(devmodeRevealCmd)
	devmodeRevealCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	devmodeRevealCmd.MarkFlagRequired("udid")
}

/*
References:
https://github.com/libimobiledevice/libimobiledevice/blob/master/tools/idevicedevmodectl.c#L440
*/
