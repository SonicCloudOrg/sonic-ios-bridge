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
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get iOS device list",
	Long:  "Get iOS device list",
	RunE: func(cmd *cobra.Command, args []string) error {
		usbMuxClient, err := giDevice.NewUsbmux()
		if err != nil {
			return util.NewErrorPrint(util.ErrConnect, "usbMux", err)
		}
		list, err1 := usbMuxClient.Devices()
		remoteList, err2 := util.ReadRemote()

		if err1 != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "listDevices", err1)
		}
		if len(list) != 0 || len(remoteList) != 0 {
			if len(list) == 0 {
				list = []giDevice.Device{}
			}
			var deviceList entity.DeviceList
			for _, d := range list {
				deviceByte, _ := json.Marshal(d.Properties())
				device := &entity.Device{}
				if isDetail {
					detail, err2 := entity.GetDetail(d)
					if err2 != nil {
						return err2
					}
					device.DeviceDetail = *detail
				}
				json.Unmarshal(deviceByte, device)
				device.Status = device.GetStatus()
				device.RemoteAddr = "localhost"
				deviceList.DeviceList = append(deviceList.DeviceList, *device)
			}
			if err2 == nil {
				for k, dev := range remoteList {
					deviceByte, _ := json.Marshal(dev.Properties())
					device := &entity.Device{}
					if isDetail {
						detail, err2 := entity.GetDetail(dev)
						if err2 != nil {
							return err2
						}
						device.DeviceDetail = *detail
					}
					json.Unmarshal(deviceByte, device)
					device.Status = device.GetStatus()
					device.RemoteAddr = k
					deviceList.DeviceList = append(deviceList.DeviceList, *device)
				}
			}
			data := util.ResultData(deviceList)
			fmt.Println(util.Format(data, isFormat, isDetail))
		} else {
			for _, v := range remoteList {
				list = append(list, v)
			}
			if len(list) != 0 {
				device := &entity.Device{}
				for _, d := range list {
					if d.Properties().SerialNumber == udid {
						deviceByte, _ := json.Marshal(d.Properties())
						if isDetail {
							detail, err2 := entity.GetDetail(d)
							if err2 != nil {
								return err2
							}
							device.DeviceDetail = *detail
						}
						json.Unmarshal(deviceByte, device)
						device.Status = device.GetStatus()
						break
					}
				}
				if device.SerialNumber != "" {
					data := util.ResultData(device)
					fmt.Println(util.Format(data, isFormat, isDetail))
				} else {
					fmt.Println("device no found")
					os.Exit(0)
				}
			} else {
				fmt.Println("no device connected")
				os.Exit(0)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	devicesCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
	devicesCmd.Flags().BoolVarP(&isDetail, "detail", "d", false, "output every device's detail")
}
