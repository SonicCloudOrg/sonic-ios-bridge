package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var xcuitestCmd = &cobra.Command{
	Use:   "xcuitest",
	Short: "Run XCUITest on your devices",
	Long:  `Run XCUITest on your devices`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(autoProxy)
		fmt.Println("xcuitest called")
	},
}

var (
	wdaBundleID string
	wdaPort     int
	mjpegPort   int
	autoProxy   bool
)

func init() {
	rootCmd.AddCommand(xcuitestCmd)
	xcuitestCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	xcuitestCmd.MarkFlagRequired("udid")
	xcuitestCmd.Flags().StringVarP(&wdaBundleID, "bundleId", "b", "com.facebook.WebDriverAgentRunner.xctrunner", "WebDriverAgentRunner bundleId")
	xcuitestCmd.Flags().BoolVarP(&autoProxy, "auto-proxy", "", false, "auto proxy ports to local with same port number")
	xcuitestCmd.Flags().IntVarP(&wdaPort, "wda-port", "", 8100, "WebDriverAgentRunner test port")
	xcuitestCmd.Flags().IntVarP(&mjpegPort, "mjpeg-port", "", 9100, "mjpeg-server port")
}
