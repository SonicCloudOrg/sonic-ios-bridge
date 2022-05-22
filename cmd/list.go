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
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "show app list",
	Long:  "show app list",
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
				result, errList := device.InstallationProxyBrowse(
					giDevice.WithApplicationType(giDevice.ApplicationTypeUser),
					giDevice.WithReturnAttributes("CFBundleVersion", "CFBundleDisplayName", "CFBundleIdentifier"))
				if errList != nil {
					return util.NewErrorPrint(util.ErrSendCommand, "appList", errList)
				}
				var appList entity.AppList
				for _, app := range result {
					a := entity.Application{}
					mapstructure.Decode(app, &a)
					if a.CFBundleIdentifier != "" && a.CFBundleDisplayName != "" && a.CFBundleVersion != "" {
						appList.ApplicationList = append(appList.ApplicationList, a)
					}
				}
				data := util.ResultData(appList)
				fmt.Println(util.Format(data, isFormat, isJson))
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
	appCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	listCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
	listCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
}
