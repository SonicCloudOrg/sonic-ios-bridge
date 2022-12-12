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

var afcRMCmd = &cobra.Command{
	Use:   "rm",
	Short: "delete file",
	Long:  "delete file",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer := getAFCServer()
		err := (afcServer).Remove(rmFilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("rm success")
		return nil
	},
}

var rmFilePath string

func initRM() {
	afcRootCMD.AddCommand(afcRMCmd)
	afcRMCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcRMCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "app bundleId")
	afcRMCmd.Flags().StringVarP(&rmFilePath, "file", "f", "", "the address of the file to be deleted")
	afcRMCmd.MarkFlagRequired("file")
}
