package afcUtil

import (
	"fmt"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"os"
	gPath "path"
)

var afcRMTreeCmd = &cobra.Command{
	Use:   "rmtree",
	Short: "recursively delete all files in a directory",
	Long:  "recursively delete all files in a directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer:=getAFCServer()
		removeTree(afcServer,rmDir)
		fmt.Println("success")
		return nil
	},
}

var rmDir string

func initRMTree() {
	afcRootCMD.AddCommand(afcRMTreeCmd)
	afcRMTreeCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcRMTreeCmd.Flags().StringVarP(&bundleID, "bundleId", "b", "", "app bundleId")
	afcRMTreeCmd.Flags().StringVarP(&rmDir,"folder","f","","folder address to delete")
	afcRMTreeCmd.MarkFlagRequired("folder")
}

func removeTree(afc giDevice.Afc, devicePath string) {
	fileInfo, err := afc.Stat(devicePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	if fileInfo.IsDir() {
		fileNames, err := afc.ReadDir(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for _, file := range fileNames {
			if file == "." || file == ".." {
				continue
			}
			var childPath string
			childPath = gPath.Join(devicePath, file)

			removeTree(afc, childPath)
		}

		err = afc.Remove(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	} else {
		err := afc.Remove(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
