package cmd

import (
	"github.com/SonicCloudOrg/sonic-ios-bridge/cmd/devmode"
	"github.com/spf13/cobra"
)

var devmodeCmd = &cobra.Command{
	Use:   "devmode",
	Short: "Enable Developer Mode on iOS 16+ devices or print the current status.",
	Long:  "Enable Developer Mode on iOS 16+ devices or print the current status.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(devmodeCmd)
	devmode.InitDevmode(devmodeCmd)
}
