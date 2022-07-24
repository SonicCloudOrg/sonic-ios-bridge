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
	afcRMCmd.Flags().StringVarP(&bundleID, "bundleId", "b", "", "app bundleId")
	afcRMCmd.Flags().StringVarP(&rmFilePath, "file", "f", "", "the address of the file to be deleted")
	afcRMCmd.MarkFlagRequired("file")
}
