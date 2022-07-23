package afcUtil

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var afcMkDirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "create a directory",
	Long:  "create a directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		if afcServer==nil {
			getAFCServer()
		}
		err := (afcServer).Mkdir(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("mkdir success")
		return nil
	},
}

func initMkDir() {
	afcRootCMD.AddCommand(afcMkDirCmd)
}
