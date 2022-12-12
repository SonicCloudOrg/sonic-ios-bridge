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
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"os"

	"github.com/spf13/cobra"
)

var locationSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set simulate location to your device.",
	Long:  "Set simulate location to your device.",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		err := device.SimulateLocationUpdate(long, lat, giDevice.CoordinateSystemBD09)
		if err != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "location set", err)
		}
		fmt.Println("location set successful!")
		return nil
	},
}

var long, lat float64

func initLocationSet() {
	locationRootCMD.AddCommand(locationSetCmd)
	locationSetCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	locationSetCmd.Flags().Float64Var(&long, "long", 0, "longitude")
	locationSetCmd.Flags().Float64Var(&lat, "lat", 0, "latitude")
	locationSetCmd.MarkFlagRequired("long")
	locationSetCmd.MarkFlagRequired("lat")
}
