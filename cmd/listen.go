/*
 *   sonic-ios-bridge  Connect to your iOS Devices.
 *   Copyright (C) 2022 SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published
 *   by the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package cmd

import (
	"context"
	"fmt"
	"sync"

	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listener for devices status",
	Long:  "Listener for devices status",
	RunE: func(cmd *cobra.Command, args []string) error {
		var deviceIdMap = make(map[int]string)
		wg := new(sync.WaitGroup)
		wg.Add(1)
		util.UsbmuxListen(func(gidevice *giDevice.Device, device *entity.Device, e error, cancelFunc context.CancelFunc) {
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
			logrus.Debugf("Device %s is %s", device.SerialNumber, device.Status)
			if device.Status == "online" {
				if isDetail {
					detail, err1 := entity.GetDetail(*gidevice)
					if err1 != nil {
						logrus.Warnf("Error: %+v", err1)
					} else {
						device.DeviceDetail = *detail
					}
				}
			}
			data := util.ResultData(device)
			fmt.Println(util.Format(data, isFormat, isDetail))
		}, true)
		wg.Wait()
		return nil
	},
}

func init() {
	devicesCmd.AddCommand(listenCmd)
	listenCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
	listenCmd.Flags().BoolVarP(&isDetail, "detail", "d", false, "output every device's detail")
}
