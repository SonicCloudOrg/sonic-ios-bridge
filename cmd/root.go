package cmd

import (
	"github.com/spf13/cobra/doc"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var isJson, isDetail, isFormat bool

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
	err1 := doc.GenMarkdownTree(rootCmd, "/")
	if err1 != nil {
		log.Fatal(err1)
	}
	if err != nil {
		os.Exit(1)
	}
}
