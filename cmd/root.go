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
	"os"

	"github.com/spf13/cobra"
)

var isJson, isDetail, isFormat bool
var udid string

var rootCmd = &cobra.Command{
	Use:   "sib",
	Short: "Bridge of iOS Devices by usbmuxd",
	Long:`
   ▄▄▄▄      ▄▄▄▄    ▄▄▄   ▄▄   ▄▄▄▄▄▄      ▄▄▄▄
 ▄█▀▀▀▀█    ██▀▀██   ███   ██   ▀▀██▀▀    ██▀▀▀▀█
 ██▄       ██    ██  ██▀█  ██     ██     ██▀
  ▀████▄   ██    ██  ██ ██ ██     ██     ██
      ▀██  ██    ██  ██  █▄██     ██     ██▄
 █▄▄▄▄▄█▀   ██▄▄██   ██   ███   ▄▄██▄▄    ██▄▄▄▄█
  ▀▀▀▀▀      ▀▀▀▀    ▀▀   ▀▀▀   ▀▀▀▀▀▀      ▀▀▀▀

          Author: SonicCloudOrg
https://github.com/SonicCloudOrg/sonic-ios-bridge
`,
}

// Execute error
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	//err1 := doc.GenMarkdownTree(rootCmd, "doc")
	//if err1 != nil {
	//	log.Fatal(err1)
	//}
}
