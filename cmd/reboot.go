package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"os"

	"github.com/spf13/cobra"
)

var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot device",
	Long:  "Reboot device",
	RunE: func(cmd *cobra.Command, args []string) error {
		usbMuxClient, err := giDevice.NewUsbmux()
		if err != nil {
			return util.NewErrorPrint(util.ErrConnect, "usbMux", err)
		}
		list, err1 := usbMuxClient.Devices()
		if err1 != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "listDevices", err1)
		}
		if len(list) != 0 {
			var device giDevice.Device
			for i, d := range list {
				if d.Properties().SerialNumber == udid {
					device = list[i]
					break
				}
			}
			if device.Properties().SerialNumber != "" {
				//errReboot := device.
				//if errReboot != nil {
				//	fmt.Println("reboot failed")
				//	os.Exit(0)
				//}
			} else {
				fmt.Println("device no found")
				os.Exit(0)
			}
		} else {
			fmt.Println("no device connected")
			os.Exit(0)
		}
		fmt.Println("install successful")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rebootCmd)
	rebootCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	rebootCmd.MarkFlagRequired("udid")
}
