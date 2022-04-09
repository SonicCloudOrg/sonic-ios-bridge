package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"os"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall App in your device",
	Long:  "Uninstall App in your device",
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
			if len(udid) != 0 {
				for i, d := range list {
					if d.Properties().SerialNumber == udid {
						device = list[i]
						break
					}
				}
			} else {
				device = list[0]
			}
			if device.Properties().SerialNumber != "" {
				util.CheckMount(device)
				errUninstall := device.AppUninstall(bundleId)
				if errUninstall != nil {
					fmt.Println("uninstall failed")
					os.Exit(0)
				}
			} else {
				fmt.Println("device no found")
				os.Exit(0)
			}
		} else {
			fmt.Println("no device connected")
			os.Exit(0)
		}
		fmt.Println("uninstall successful")
		return nil
	},
}

var bundleId string

func init() {
	appCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	uninstallCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "target bundleId")
	uninstallCmd.MarkFlagRequired("bundleId")
}
