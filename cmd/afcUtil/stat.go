package afcUtil

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var afcStatCmd = &cobra.Command{
	Use:   "stat",
	Short: "view file details",
	Long:  "view file details",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer:=getAFCServer()
		info, err := (afcServer).Stat(statPath)
		if err != nil {
			os.Exit(0)
		}
		if info.IsDir() {
			fmt.Println("type:DIR")
		} else {
			fmt.Println("type:FILE")
		}
		fmt.Println("CTime:", info.CreationTime().Format("2006-01-02 15:04:05"))
		fmt.Println("MTime:", info.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Println(fmt.Sprintf("Size:%d", info.Size()))
		return nil
	},
}

var statPath string

func initStat() {
	afcRootCMD.AddCommand(afcStatCmd)
	afcStatCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcStatCmd.Flags().StringVarP(&bundleID, "bundleId", "b", "", "app bundleId")
	afcStatCmd.Flags().StringVarP(&statPath,"path","p","","files or folders for which details need to be viewed")
	afcStatCmd.MarkFlagRequired("path")
}
