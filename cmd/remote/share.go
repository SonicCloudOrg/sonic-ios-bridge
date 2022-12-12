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
package remote

import (
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Sharing device to remote",
	Long:  "Sharing device to remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		log.Printf("start sharing, target device will share on port:%d", port)
		err := device.Share(port)
		if err != nil {
			log.Panic(err)
		}
		return nil
	},
}

func shareInit() {
	remoteCmd.AddCommand(shareCmd)
	shareCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	shareCmd.Flags().IntVarP(&port, "port", "p", 9123, "share port")
}
