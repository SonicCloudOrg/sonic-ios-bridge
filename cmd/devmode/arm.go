package devmode

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

var devmodeArmCmd = &cobra.Command{
	Use:   "arm",
	Short: "Arm the Developer Mode (device will reboot)",
	Long:  "Arm the Developer Mode (device will reboot)",
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
		res, errArm := amfi.DevModeArm()
		if errArm != nil {
			return errArm
		}
		if res == http.StatusOK {
			logrus.Infof("Developer Mode armed.")
			return nil
		} else {
			strErrMsg := fmt.Sprintf("Failed to arm Developer Mode (%d).", res)
			logrus.Warn(strErrMsg)
			return xerrors.New(strErrMsg)
		}
	},
}

func initDevModeArmCmd() {
	devmodeRootCMD.AddCommand(devmodeArmCmd)
	devmodeArmCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	devmodeArmCmd.MarkFlagRequired("udid")
}
