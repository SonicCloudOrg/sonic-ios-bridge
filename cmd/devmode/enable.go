package devmode

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// option bindings
var bWaitReboot, bAutoConfirm bool
var intEnableWaitTimeout int

var devmodeEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable Developer Mode (device will reboot)",
	Long:  "Enable Developer Mode (device will reboot)",
	RunE: func(cmd *cobra.Command, args []string) error {
		errArm := devmodeArmCmd.RunE(cmd, args)
		if errArm != nil {
			return errArm
		}
		if bWaitReboot {
			bIsDeviceOnline := true
			wg := new(sync.WaitGroup)
			wg.Add(1)
			var deviceIdMap = make(map[int]string)
			shutDownFun := util.UsbmuxListen(func(gidevice *giDevice.Device, device *entity.Device, e error, cancelFunc context.CancelFunc) {
				if e != nil {
					logrus.Warnf("Error: %+v", e)
				}
				if device == nil {
					return
				}
				if len(device.SerialNumber) > 0 {
					deviceIdMap[device.DeviceID] = device.SerialNumber
				} else {
					device.SerialNumber = deviceIdMap[device.DeviceID]
					delete(deviceIdMap, device.DeviceID)
				}
				if device.SerialNumber == udid {
					logrus.Infof("Device %s is %s.", device.SerialNumber, device.Status)
					if device.Status == "offline" {
						bIsDeviceOnline = false
					} else if !bIsDeviceOnline && device.Status == "online" { // transition from offline to online
						if cancelFunc != nil {
							cancelFunc()
						}
						bIsDeviceOnline = true
						wg.Done()
						return
					}
				} else {
					logrus.Debugf("Device %s is %s.", device.SerialNumber, device.Status)
				}
			})
			go (func(cancelFunc *context.CancelFunc) { // timer to cancel listening
				time.Sleep(time.Duration(intEnableWaitTimeout) * time.Second)
				logrus.Warnf("Timeout waiting for device %s to reboot.", udid)
				if cancelFunc != nil {
					(*cancelFunc)()
				}
				wg.Done()
			})(&shutDownFun)
			wg.Wait()
			if bIsDeviceOnline && bAutoConfirm {
				bPreCheckIOSVer = false
				devmodeConfirmCmd.RunE(cmd, args)
			} else {
				executable, _ := os.Executable()
				pkgPath := strings.Split(__PACKAGE__(), "/")
				logrus.Infof("Please check the device %s is online and then run '%s %s %s -u %s'.", udid, executable, pkgPath[len(pkgPath)-1], devmodeConfirmCmd.Use, udid)
			}
		}
		return nil
	},
}

func initDevModeEnableCmd() {
	devmodeRootCMD.AddCommand(devmodeEnableCmd)
	devmodeEnableCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	devmodeEnableCmd.MarkFlagRequired("udid")
	devmodeEnableCmd.Flags().BoolVar(&bWaitReboot, "wait", false, "wait for reboot to complete")
	devmodeEnableCmd.Flags().IntVar(&intEnableWaitTimeout, "wait-timeout", 60, "wait timeout in seconds")
	devmodeEnableCmd.Flags().BoolVarP(&bAutoConfirm, "confirm", "y", false, "automatically confirm after reboot")
}
