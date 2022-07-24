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

var afcMkDirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "create a directory",
	Long:  "create a directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer := getAFCServer()
		err := (afcServer).Mkdir(mkDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("mkdir success")
		return nil
	},
}

var mkDir string

func initMkDir() {
	afcRootCMD.AddCommand(afcMkDirCmd)
	afcMkDirCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcMkDirCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "app bundleId")
	afcMkDirCmd.Flags().StringVarP(&mkDir, "folder", "f", "", "mkdir directory path")
	afcMkDirCmd.MarkFlagRequired("folder")
}
