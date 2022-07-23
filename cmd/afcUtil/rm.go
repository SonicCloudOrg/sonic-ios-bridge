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
		afcServer:=getAFCServer()
		err := (afcServer).Remove(rmFilePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("rm success")
		return nil
	},
}

var rmFilePath string

func initRM() {
	afcRootCMD.AddCommand(afcRMCmd)
	afcRMCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcRMCmd.Flags().StringVarP(&bundleID, "bundleId", "b", "", "app bundleId")
	afcRMCmd.Flags().StringVarP(&rmFilePath,"file","f","","the address of the file to be deleted")
	afcRMCmd.MarkFlagRequired("file")
}
