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
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"os"

	"github.com/spf13/cobra"
)

var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot or Shutdown device",
	Long:  "Reboot or Shutdown device",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		var errReboot error
		if isShutdown {
			errReboot = device.Shutdown()
		} else {
			errReboot = device.Reboot()
		}
		if errReboot != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "reboot", errReboot)
		}
		fmt.Println("reboot successful")
		return nil
	},
}

var isShutdown bool

func init() {
	rootCmd.AddCommand(rebootCmd)
	rebootCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	rebootCmd.Flags().BoolVarP(&isShutdown, "shutdown", "s", false, "shutdown your device")
}
