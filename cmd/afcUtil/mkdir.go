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
		afcServer:=getAFCServer()
		err := (afcServer).Mkdir(mkDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("mkdir success")
		return nil
	},
}

var mkDir string

func initMkDir() {
	afcRootCMD.AddCommand(afcMkDirCmd)
	afcMkDirCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcMkDirCmd.Flags().StringVarP(&bundleID, "bundleId", "b", "", "app bundleId")
	afcMkDirCmd.Flags().StringVarP(&mkDir,"folder",  "f","", "mkdir directory path")
	afcMkDirCmd.MarkFlagRequired("folder")
}
