/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
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
		if err1 != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "listDevices", err1)
		}
		if len(udid) == 0 {
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
				deviceList.DeviceList = append(deviceList.DeviceList, *device)
			}
			data := util.ResultData(deviceList)
			fmt.Println(util.Format(data, isFormat, isDetail))
		} else {
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
