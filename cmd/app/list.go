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
		appType := giDevice.ApplicationTypeUser
		if showSystem {
			appType = giDevice.ApplicationTypeAny
		}
		result, errList := device.InstallationProxyBrowse(
			giDevice.WithApplicationType(appType),
			giDevice.WithReturnAttributes("CFBundleShortVersionString", "CFBundleVersion", "CFBundleDisplayName", "CFBundleIdentifier"))
		if errList != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "app list", errList)
		}
		var appList entity.AppList
		for _, app := range result {
			a := entity.Application{}
			mapstructure.Decode(app, &a)
			if a.CFBundleIdentifier != "" || a.CFBundleDisplayName != "" || a.CFBundleShortVersionString != "" || a.CFBundleVersion != "" {
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

var (
	showSystem bool
	showIcon   bool
)

func initAppList() {
	appRootCMD.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&showIcon, "icon", "i", false, "show app icon")
	listCmd.Flags().BoolVarP(&showSystem, "system", "s", false, "show system apps")
	listCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	listCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
	listCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
}
