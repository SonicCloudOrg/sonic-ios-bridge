package afcUtil

import (
	"fmt"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"io"
	"os"
	gPath "path"
)

var afcPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull file or directory from device",
	Long:  "pull file or directory from device",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer:=getAFCServer()
		pullOperate(afcServer, pullDevicePath, pullSaveLocalPath)
		fmt.Println(fmt.Sprintf("success,pull %s --> %s", pullDevicePath, pullSaveLocalPath))
		return nil
	},
}

var pullDevicePath 		string
var pullSaveLocalPath 	string
func initPullCmd() {
	afcRootCMD.AddCommand(afcPullCmd)
	afcPullCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcPullCmd.Flags().StringVarP(&bundleID, "bundleId", "b", "", "app bundleId")
	afcPullCmd.Flags().StringVarP(&pullDevicePath,"devicePath",  "d","", "pull file or directory device path")
	afcPullCmd.Flags().StringVarP(&pullSaveLocalPath,"localPath",  "l","", "pull save file or directory to local path")
	afcPullCmd.MarkFlagRequired("devicePath")
	afcPullCmd.MarkFlagRequired("localPath")
}

func pullOperate(afc giDevice.Afc, devicePath string, localPath string) {
	fileInfo, err := afc.Stat(devicePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	if fileInfo.IsDir() {
		localFile, err := os.ReadDir(localPath)
		if localFile == nil || err != nil {
			mkdirError := os.Mkdir(localPath, os.ModePerm)
			if mkdirError != nil {
				fmt.Println(mkdirError)
				os.Exit(0)
			}
		}
		fileNames, err := afc.ReadDir(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for _, file := range fileNames {
			if file == "." || file == ".." {
				continue
			}
			pullOperate(afc, gPath.Join(devicePath, file), gPath.Join(localPath, file))
		}
	} else {
		pullFile(afc, devicePath, localPath)
	}
}

func pullFile(afc giDevice.Afc, devicePath string, localPath string) {
	afcFile, err := afc.Open(devicePath, giDevice.AfcFileModeRdOnly)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		_ = afcFile.Close()
	}()
	file, err := os.Create(localPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		_ = file.Close()
	}()
	if _, err = io.Copy(file, afcFile); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
