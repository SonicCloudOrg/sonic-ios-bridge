package afcUtil

import (
	"fmt"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var afcCatCmd = &cobra.Command{
	Use:   "cat",
	Short: "cat to view files",
	Long:  "cat to view files",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer:=getAFCServer()
		catFile(afcServer, catFilePath)
		return nil
	},
}

var catFilePath string

func initCat() {
	afcRootCMD.AddCommand(afcCatCmd)
	afcCatCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcCatCmd.Flags().StringVarP(&bundleID, "bundleId", "b", "", "app bundleId")
	afcCatCmd.Flags().StringVarP(&catFilePath, "file",  "f","", "cat file path")
	afcCatCmd.MarkFlagRequired("file")
}

func catFile(afc giDevice.Afc, filePath string) {
	fileInfo, err := afc.Stat(filePath)
	if err != nil {
		fmt.Println("file path is null")
		os.Exit(0)
	}
	p := make([]byte, fileInfo.Size())
	afcFile, err := afc.Open(filePath, giDevice.AfcFileModeRdOnly)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		_ = afcFile.Close()
	}()
	n, err := afcFile.Read(p)
	if err == io.EOF {
		fmt.Println(err)
		os.Exit(0)
		//break
	}
	fmt.Print(string(p[:n]))
}
