package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// xcuitestCmd represents the xcuitest command
var xcuitestCmd = &cobra.Command{
	Use:   "xcuitest",
	Short: "Run XCUITest on your devices",
	Long:  `Run XCUITest on your devices`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("xcuitest called")
	},
}

var (
	wdaBundleID string
	wdaPort int
	mjpegPort int
	autoProxy bool
)

func init() {
	rootCmd.AddCommand(xcuitestCmd)
	xcuitestCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	xcuitestCmd.MarkFlagRequired("udid")
}
