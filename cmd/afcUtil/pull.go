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
		if len(args) != 2 {
			fmt.Println("arguments error")
			os.Exit(0)
		}
		if afcServer==nil {
			getAFCServer()
		}
		pullOperate(afcServer, args[0], args[1])
		fmt.Println(fmt.Sprintf("success,pull %s --> %s", args[0], args[1]))
		return nil
	},
}

func initPullCmd() {
	afcRootCMD.AddCommand(afcPullCmd)
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
