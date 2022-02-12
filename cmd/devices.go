package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get iOS device list",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("devices called")
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
}
