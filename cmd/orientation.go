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
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"os"
	"os/signal"
	"syscall"
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
			shutDown := make(chan os.Signal, syscall.SIGTERM)
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
