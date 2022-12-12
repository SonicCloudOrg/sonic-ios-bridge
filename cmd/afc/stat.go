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
package afc

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var afcStatCmd = &cobra.Command{
	Use:   "stat",
	Short: "view file details",
	Long:  "view file details",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer := getAFCServer()
		info, err := (afcServer).Stat(statPath)
		if err != nil {
			os.Exit(0)
		}
		if info.IsDir() {
			fmt.Println("type:DIR")
		} else {
			fmt.Println("type:FILE")
		}
		fmt.Println("CTime:", info.CreationTime().Format("2006-01-02 15:04:05"))
		fmt.Println("MTime:", info.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Println(fmt.Sprintf("Size:%d", info.Size()))
		return nil
	},
}

var statPath string

func initStat() {
	afcRootCMD.AddCommand(afcStatCmd)
	afcStatCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcStatCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "app bundleId")
	afcStatCmd.Flags().StringVarP(&statPath, "path", "p", "", "files or folders for which details need to be viewed")
	afcStatCmd.MarkFlagRequired("path")
}
