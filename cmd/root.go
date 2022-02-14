package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var isJson, isDetail bool

var rootCmd = &cobra.Command{
	Use:   "sib",
	Short: "Bridge of iOS Devices by usbmuxd",
}

// Execute error
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
