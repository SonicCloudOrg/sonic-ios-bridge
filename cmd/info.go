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
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"os"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show domain info in your device",
	Long:  "Show domain info in your device",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		interResult, err := device.GetValue(domain, key)
		if err != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "get value", err)
		}
		if isFormat {
			result, _ := json.MarshalIndent(interResult, "", "\t")
			fmt.Println(string(result))
		} else {
			result, _ := json.Marshal(interResult)
			fmt.Println(string(result))
		}
		return nil
	},
}

var domain, key string

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	infoCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain")
	infoCmd.Flags().StringVarP(&key, "key", "k", "", "Key")
	infoCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
}
