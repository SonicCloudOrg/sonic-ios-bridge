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
	"encoding/json"
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listener for devices status",
	Long:  "Listener for devices status",
	RunE: func(cmd *cobra.Command, args []string) error {
		usbMuxClient, err := giDevice.NewUsbmux()
		if err != nil {
			return util.NewErrorPrint(util.ErrConnect, "usbMux", err)
		}
		model := make(chan giDevice.Device)
		shutDownFun, err2 := usbMuxClient.Listen(model)
		if err2 != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "listen", err2)
		}

		shutDown := make(chan os.Signal, 1)
		signal.Notify(shutDown, os.Interrupt, os.Kill)
		var deviceIdMap = make(map[int]string)
		for {
			select {
			case d := <-model:
				deviceByte, _ := json.Marshal(d.Properties())
				device := &entity.Device{}
				json.Unmarshal(deviceByte, device)
				if len(device.SerialNumber) > 0 {
					deviceIdMap[device.DeviceID] = device.SerialNumber
				} else {
					device.SerialNumber = deviceIdMap[device.DeviceID]
					delete(deviceIdMap, device.DeviceID)
				}
				device.Status = device.GetStatus()
				if device.Status == "online" && isDetail {
					detail, err1 := entity.GetDetail(d)
					if err1 != nil {
						continue
					}
					device.DeviceDetail = *detail
				}
				data := util.ResultData(device)
				fmt.Println(util.Format(data, isFormat, isDetail))
			case <-shutDown:
				shutDownFun()
				return nil
			}
		}
		return nil
	},
}

func init() {
	devicesCmd.AddCommand(listenCmd)
	listenCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
	listenCmd.Flags().BoolVarP(&isDetail, "detail", "d", false, "output every device's detail")
}
