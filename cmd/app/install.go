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
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
	"os"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install App in your device",
	Long:  "Install App in your device",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		errInstall := device.AppInstall(path)
		if errInstall != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "install", errInstall)
		}
		fmt.Println("install successful")
		return nil
	},
}

var path string

func initAppInstall() {
	appRootCMD.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	installCmd.Flags().StringVarP(&path, "path", "p", "", "path of ipa file")
	installCmd.MarkFlagRequired("path")
}
