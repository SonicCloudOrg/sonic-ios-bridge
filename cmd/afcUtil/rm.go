package afcUtil

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var afcRMCmd = &cobra.Command{
	Use:   "rm",
	Short: "delete file",
	Long:  "delete file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if afcServer==nil {
			getAFCServer()
		}
		err := (afcServer).Remove(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("rm success")
		return nil
	},
}

func initRM() {
	afcRootCMD.AddCommand(afcRMCmd)
}
