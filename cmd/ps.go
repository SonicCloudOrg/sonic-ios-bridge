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
	"encoding/json"
	"fmt"
	"os"

	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Show process in your device",
	Long:  "Show process in your device",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		processList, errProcess := device.AppRunningProcesses()
		if errProcess != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "processList", errProcess)
		}
		if isFormat {
			result, _ := json.MarshalIndent(processList, "", "\t")
			fmt.Println(string(result))
		} else if isJson {
			result, _ := json.Marshal(processList)
			fmt.Println(string(result))
		} else {
			for _, process := range processList {
				fmt.Printf("%s %d %s %s\n", process.Name, process.Pid, process.RealAppName, process.StartDate)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	psCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
	psCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
}
