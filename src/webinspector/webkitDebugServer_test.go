package webinspector

import (
	"encoding/json"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"testing"
	"time"
)

var device giDevice.Device

func setupDeviceSrv() {
	device = util.GetDeviceByUdId("")
}

func TestWebkitDebugService(t *testing.T) {
	setupDeviceSrv()

	webkitDebug := NewWebkitDebugServer(&device)
	SetRPCDebug(true)
	// init
	err := webkitDebug.ConnectInspector()
	if err != nil {
		panic(err)
	}
	time.Sleep(6 * time.Second)
	// get all page
	pages, err := webkitDebug.GetOpenPages()
	if err != nil {
		panic(err)
	}
	time.Sleep(6 * time.Second)
	if arr, err1 := json.MarshalIndent(pages, "", "\t"); err1 != nil {
		panic(err1)
	} else {
		fmt.Println(string(arr))
	}
	time.Sleep(6 * time.Second)
	// find page
	app, page, err := webkitDebug.FindPagesByID("1")
	if err != nil {
		panic(err)
	}
	time.Sleep(6 * time.Second)
	// start cdp
	err = webkitDebug.StartCDP(app.ApplicationID, page.PageID)
	if err != nil {
		panic(err)
	}
	time.Sleep(40 * time.Second)
}
