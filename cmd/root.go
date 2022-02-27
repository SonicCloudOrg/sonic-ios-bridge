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

      Bridge of iOS Devices by usbmuxd.
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
}
