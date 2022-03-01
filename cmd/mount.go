package cmd

import (
	"github.com/spf13/cobra"
)

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var dir string

func init() {
	rootCmd.AddCommand(mountCmd)
	mountCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	mountCmd.MarkFlagRequired("udid")
}
