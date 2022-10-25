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
	"os/signal"
	"syscall"
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

		shutDown := make(chan os.Signal, syscall.SIGTERM)
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
