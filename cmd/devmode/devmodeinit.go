package devmode

import (
	"encoding/json"
	"os"
	"reflect"

	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

var devmodeRootCMD *cobra.Command
var bPreCheckIOSVer bool = true

// option bindings
var udid string
var bIsOutputJson bool

func InitDevmode(devmodeCMD *cobra.Command) {
	devmodeRootCMD = devmodeCMD
	initDevModeListCmd()
	initDevModeArmCmd()
	initDevModeRevealCmd()
	initDevModeEnableCmd()
	initDevModeConfirmCmd()
}

func getAmfiServer() (giDevice.Amfi, error) {
	device := util.GetDeviceByUdId(udid)
	if device == nil {
		os.Exit(55 /* https://git.islam.gov.my/mohdrizal/fabric/-/blob/v2.3.3/vendor/golang.org/x/sys/windows/zerrors_windows.go#L195 */)
	}
	return device.AmfiService()
}

func canToggleDevMode(udid string) (bool, error) {
	gidevice := util.GetDeviceByUdId(udid)
	if gidevice == nil {
		return false, xerrors.Errorf("Device %s not found", udid)
	}
	device := entity.Device{}
	deviceByte, _ := json.Marshal(gidevice.Properties())
	json.Unmarshal(deviceByte, &device)
	detail, err2 := entity.GetDetail(gidevice)
	if err2 != nil {
		return false, err2
	} else {
		device.DeviceDetail = *detail
	}
	devmode := entity.DevMode{Device: device}
	b, e := devmode.CanCheck()
	if e != nil {
		return false, e
	}
	return b, nil
}

func __PACKAGE__() string { // https://www.appsloveworld.com/go/3/how-to-get-name-of-current-package-in-go?expand_article=1
	type dummy struct{}
	return reflect.TypeOf(dummy{}).PkgPath()
}
