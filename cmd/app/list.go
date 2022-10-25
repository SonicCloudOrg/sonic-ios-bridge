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
package app

import (
	"encoding/base64"
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "show app list",
	Long:  "show app list",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		result, errList := device.InstallationProxyBrowse(
			giDevice.WithApplicationType(giDevice.ApplicationTypeUser),
			giDevice.WithReturnAttributes("CFBundleVersion", "CFBundleDisplayName", "CFBundleIdentifier"))
		if errList != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "app list", errList)
		}
		var appList entity.AppList
		for _, app := range result {
			a := entity.Application{}
			mapstructure.Decode(app, &a)
			if a.CFBundleIdentifier != "" && a.CFBundleDisplayName != "" && a.CFBundleVersion != "" {
				if showIcon {
					icon, errIcon := device.GetIconPNGData(a.CFBundleIdentifier)
					if errIcon == nil {
						data, _ := ioutil.ReadAll(icon)
						a.IconBase64 = base64.StdEncoding.EncodeToString(data)
					}
				}
				appList.ApplicationList = append(appList.ApplicationList, a)
			}
		}
		data := util.ResultData(appList)
		fmt.Println(util.Format(data, isFormat, isJson))
		return nil
	},
}

var showIcon bool

func initAppList() {
	appRootCMD.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&showIcon, "icon", "i", false, "show app icon")
	listCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	listCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
	listCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
}
