package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"os"

	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch App",
	Long:  "Launch App",
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
				_, errLaunch := device.AppLaunch(bundleId)
				if errLaunch != nil {
					fmt.Println("launch failed")
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
		return nil
	},
}

func init() {
	appCmd.AddCommand(launchCmd)
	launchCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	launchCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "target bundleId")
	launchCmd.MarkFlagRequired("bundleId")
}
