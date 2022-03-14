package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version code of sib",
	Long:  "Version code of sib",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.0.2")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
