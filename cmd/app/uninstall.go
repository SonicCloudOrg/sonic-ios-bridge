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

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall App in your device",
	Long:  "Uninstall App in your device",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		errUninstall := device.AppUninstall(bundleId)
		if errUninstall != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "uninstall", errUninstall)
		}
		fmt.Println("uninstall successful")
		return nil
	},
}

func initAppUninstall() {
	appRootCMD.AddCommand(uninstallCmd)
	uninstallCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	uninstallCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "target bundleId")
	uninstallCmd.MarkFlagRequired("bundleId")
}
