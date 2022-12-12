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
package location

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"os"

	"github.com/spf13/cobra"
)

var locationUnsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset simulate location to your device.",
	Long:  "Unset simulate location to your device.",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		err := device.SimulateLocationRecover()
		if err != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "location unset", err)
		}
		fmt.Println("location unset successful!")
		return nil
	},
}

func initLocationUnset() {
	locationRootCMD.AddCommand(locationUnsetCmd)
	locationUnsetCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
}
