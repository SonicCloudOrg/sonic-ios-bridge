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
