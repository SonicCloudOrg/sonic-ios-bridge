package devmode

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

var devmodeConfirmCmd = &cobra.Command{
	Use:   "confirm",
	Short: "Confirm enabling of Developer Mode",
	Long:  "Confirm enabling of Developer Mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		if bPreCheckIOSVer {
			if bCan, eCan := canToggleDevMode(udid); eCan != nil {
				strErrMsg := fmt.Sprintf("Failed to check device %s iOS version", udid)
				logrus.Warn(strErrMsg)
				return xerrors.New(strErrMsg)
			} else if !bCan {
				strErrMsg := fmt.Sprintf("Device %s iOS version below 16", udid)
				logrus.Warn(strErrMsg)
				return xerrors.New(strErrMsg)
			}
		}
		amfi, errAmfi := getAmfiServer()
		if errAmfi != nil {
			return errAmfi
		}
		res, errReveal := amfi.DevModeEnable()
		if errReveal != nil {
			return errReveal
		}
		if res == http.StatusOK {
			logrus.Infof("Developer Mode menu enabled successfully.")
			return nil
		} else {
			strErrMsg := fmt.Sprintf("Failed to enable Developer Mode menu (%d).", res)
			logrus.Warn(strErrMsg)
			return xerrors.New(strErrMsg)
		}
	},
}

func initDevModeConfirmCmd() {
	devmodeRootCMD.AddCommand(devmodeConfirmCmd)
	devmodeConfirmCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	devmodeConfirmCmd.MarkFlagRequired("udid")
}
