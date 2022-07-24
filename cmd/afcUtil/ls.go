package afcUtil

import (
	"fmt"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"os"
	gPath "path"
)

var afcLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "ls to view the directory",
	Long:  "ls to view the directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer:=getAFCServer()
		lsShow(afcServer, lsDirPath)
		return nil
	},
}

var lsDirPath string

func initLs() {
	afcRootCMD.AddCommand(afcLsCmd)
	afcLsCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcLsCmd.Flags().StringVarP(&bundleID, "bundleId", "b", "", "app bundleId")
	afcLsCmd.Flags().StringVarP(&lsDirPath,"folder",  "f","", "ls folder path")
	afcLsCmd.MarkFlagRequired("folder")
}

func lsShow(afc giDevice.Afc, filePath string) {
	fileNames, err := afc.ReadDir(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	for _, fileName := range fileNames {
		if fileName == "." || fileName == ".." {
			continue
		}
		info, err := afc.Stat(gPath.Join(filePath, fileName))

		if err != nil {
			os.Exit(0)
		}
		if info.IsDir() {
			fmt.Println(fileName + "/")
		} else {
			fmt.Println(fmt.Sprintf("- %s %d", fileName, info.Size()))
		}
	}
}