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
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

var orientationCmd = &cobra.Command{
	Use:   "orientation",
	Short: "Listener for devices orientation",
	Long:  "Listener for devices orientation",
	Run: func(cmd *cobra.Command, args []string) {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		if !isWatchOer {
			o, err := device.GetInterfaceOrientation()
			if err != nil {
				fmt.Println("get orientation failed.")
			}
			fmt.Println(fmt.Sprintf("orientation: %d", o))
		} else {
			shutDown := make(chan os.Signal, 1)
			signal.Notify(shutDown, os.Interrupt, os.Kill)
			go func() {
				var lo giDevice.OrientationState
				for {
					o, err := device.GetInterfaceOrientation()
					if err != nil {
						fmt.Println("get orientation failed.")
					}
					if lo != o {
						lo = o
						fmt.Println(fmt.Sprintf("orientation: %d", o))
					}
					time.Sleep(time.Duration(3) * time.Second)
				}
				shutDown <- os.Interrupt
			}()
			<-shutDown
		}
	},
}

var isWatchOer bool

func init() {
	rootCmd.AddCommand(orientationCmd)
	orientationCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	orientationCmd.Flags().BoolVarP(&isWatchOer, "watch", "w", false, "watch orientation change.")
}
