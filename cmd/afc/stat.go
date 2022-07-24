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
