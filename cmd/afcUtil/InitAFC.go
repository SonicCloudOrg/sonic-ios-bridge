package afcUtil

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"os"
)

var afcServer 	giDevice.Afc
var afcRootCMD 	*cobra.Command
var udid 		string
var bundleID 	string
// InitAfc 用于初始化，在上层中调用这个方法，否则不会正常进行初始化
func InitAfc(afcCMD *cobra.Command,pUdid string,pBundleID string){
	udid = pUdid
	bundleID = pBundleID
	afcRootCMD = afcCMD

	initMkDir()

	initTree()
	initCat()
	initLs()
	initStat()

	initPullCmd()
	initPush()

	initRM()
	initRMTree()
}

func getAFCServer()()  {
	device := util.GetDeviceByUdId(udid)
	if device == nil {
		os.Exit(0)
	}
	var err error
	if bundleID != "" {
		var houseArrestSrv giDevice.HouseArrest
		houseArrestSrv, err = device.HouseArrestService()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		afcServer, err = houseArrestSrv.Documents(bundleID)
	} else {
		afcServer, err = device.AfcService()
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}