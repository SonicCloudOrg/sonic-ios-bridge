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

func init() {
	rootCmd.AddCommand(xcuitestCmd)
}
