package usb

import (
	"os"
	"time"

	"github.com/SonicCloudOrg/sonic-ios-bridge/src/errorcodes"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// option bindings
var intReconnectWaitTime int

var usbReconnectCmd = &cobra.Command{
	Use:   "reconnect",
	Short: "Reconnct a USB device",
	Long:  "Reconnct a USB device",
	RunE: func(cmd *cobra.Command, args []string) error {
		util.InitLogger()
		devices := findUsbDevice(isAppleDevice, isMySerial(udid))
		if len(devices) <= 0 {
			os.Exit(errorcodes.ERROR_DEV_NOT_EXIST)
		}
		defer devices[0].Close()
		lineage := getUsbDeviceLineage(devices[0], 1)
		if len(lineage) <= 0 {
			os.Exit(errorcodes.ERROR_BAD_UNIT)
		}
		parent := lineage[0]
		defer parent.Close()
		logrus.Infof("Connected to hub: %+v,port=%d", parent, devices[0].Desc.Port)
		portIndex := uint16(devices[0].Desc.Port)
		toggleUsbHubPortPower(parent, portIndex, false)
		if intReconnectWaitTime > 0 {
			time.Sleep(time.Duration(intReconnectWaitTime) * time.Second)
		}
		toggleUsbHubPortPower(parent, portIndex, true)
		return nil
	},
}

func initUsbReconnectCmd() {
	usbRootCMD.AddCommand(usbReconnectCmd)
	usbReconnectCmd.Flags().StringVarP(&udid, "udid", "u", "", "target specific device by UDID")
	usbReconnectCmd.MarkFlagRequired("udid")
	usbReconnectCmd.Flags().IntVar(&intReconnectWaitTime, "wait", 3, "wait time in seconds")
}

/*
References:
https://www.gniibe.org/development/ac-power-control-by-USB-hub/
https://git.gniibe.org/cgit/gnuk/gnuk.git/tree/tool/hub_ctrl.py
*/
