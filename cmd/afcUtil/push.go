package afcUtil

import (
	"fmt"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"io"
	"os"
	gPath "path"
)

var afcPushCmd = &cobra.Command{
	Use:   "push",
	Short: "push a file or directory to the device",
	Long:  "push a file or directory to the device",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			fmt.Println("parameter error")
			os.Exit(0)
		}
		if afcServer==nil {
			getAFCServer()
		}
		pushOperate(afcServer, args[0], args[1])
		fmt.Println(fmt.Sprintf("success,push %s --> %s", args[0], args[1]))

		return nil
	},
}

func initPush() {
	afcRootCMD.AddCommand(afcPushCmd)
}

func pushFile(afc giDevice.Afc, localPath string, devicePath string) {
	file, err := os.Open(localPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = file.Close()
	}()

	afcFile, err := afc.Open(devicePath, giDevice.AfcFileModeWr)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		_ = afcFile.Close()
	}()
	if _, err = io.Copy(afcFile, file); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

}

func pushOperate(afc giDevice.Afc, localPath string, devicePath string) {
	localFile, err := os.Stat(localPath)
	if err != nil {
		os.Exit(0)
	}
	if localFile.IsDir() {
		aPathInfo, _ := afc.ReadDir(devicePath)
		if aPathInfo == nil {
			mkdirError := afc.Mkdir(devicePath)
			if mkdirError != nil {
				fmt.Println(mkdirError)
				os.Exit(0)
			}
		}
		childFiles, err := os.ReadDir(localPath)
		if err != nil {
			os.Exit(0)
		}
		for _, childFile := range childFiles {
			pushOperate(afc, gPath.Join(localPath, childFile.Name()), gPath.Join(devicePath, childFile.Name()))
		}
	} else {
		pushFile(afc, localPath, devicePath)
	}
}

