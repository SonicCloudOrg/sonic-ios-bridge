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
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/spf13/cobra"
	"os"
	gPath "path"
)

var afcRMTreeCmd = &cobra.Command{
	Use:   "rmtree",
	Short: "recursively delete all files in a directory",
	Long:  "recursively delete all files in a directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer := getAFCServer()
		removeTree(afcServer, rmDir)
		fmt.Println("success")
		return nil
	},
}

var rmDir string

func initRMTree() {
	afcRootCMD.AddCommand(afcRMTreeCmd)
	afcRMTreeCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcRMTreeCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "app bundleId")
	afcRMTreeCmd.Flags().StringVarP(&rmDir, "folder", "f", "", "folder address to delete")
	afcRMTreeCmd.MarkFlagRequired("folder")
}

func removeTree(afc giDevice.Afc, devicePath string) {
	fileInfo, err := afc.Stat(devicePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	if fileInfo.IsDir() {
		fileNames, err := afc.ReadDir(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for _, file := range fileNames {
			if file == "." || file == ".." {
				continue
			}
			var childPath string
			childPath = gPath.Join(devicePath, file)

			removeTree(afc, childPath)
		}

		err = afc.Remove(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	} else {
		err := afc.Remove(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
