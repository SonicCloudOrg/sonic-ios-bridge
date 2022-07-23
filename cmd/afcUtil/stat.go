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
		if afcServer==nil {
			getAFCServer()
		}
		info, err := (afcServer).Stat(args[0])
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

func initStat() {
	afcRootCMD.AddCommand(afcStatCmd)
}
