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
