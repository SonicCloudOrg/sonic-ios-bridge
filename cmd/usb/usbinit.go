package usb

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/google/gousb"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"golang.org/x/xerrors"
)

const (
	APPLE_VENDOR_ID = gousb.ID(0x5ac)

	USB_RT_HUB  = (gousb.ControlClass | gousb.ControlDevice)
	USB_RT_PORT = (gousb.ControlClass | gousb.ControlOther)

	REQUEST_GET_STATUS     = 0x00
	REQUEST_CLEAR_FEATURE  = 0x01
	REQUEST_SET_FEATURE    = 0x03
	REQUEST_GET_DESCRIPTOR = 0x06

	DESCRIPTOR_TYPE_HUB = 41 // https://github.com/pyusb/pyusb/blob/73e87f68fc8ece558057113690f9d2391441b58f/usb/legacy.py#L59

	USB_PORT_FEAT_POWER = 8
)

var usbRootCMD *cobra.Command
var gousbContext *gousb.Context = gousb.NewContext()
var all_devices_opener = func(desc *gousb.DeviceDesc) bool { return true }

// option bindings
var udid string
var bIsOutputJson bool

func InitUsb(usbCMD *cobra.Command) {
	usbRootCMD = usbCMD
	initUsbListCmd()
	initUsbReconnectCmd()
	initUsbHubPortPowerCmd()
}

func __PACKAGE__() string { // https://www.appsloveworld.com/go/3/how-to-get-name-of-current-package-in-go?expand_article=1
	type dummy struct{}
	return reflect.TypeOf(dummy{}).PkgPath()
}

type myUsbDevicePrefilter func(desc *gousb.DeviceDesc) bool
type myUsbDevicePostfilter func(usbdevice *gousb.Device) bool

func isAppleDevice(desc *gousb.DeviceDesc) bool {
	return desc.Vendor == APPLE_VENDOR_ID
}

func trimStringDescriptorValue(s string) string {
	return strings.Trim(s, "\x00")
}

func isMySerial(myUdid string) myUsbDevicePostfilter {
	return func(usbdevice *gousb.Device) bool {
		strSerialNumber, errSerialNumber := usbdevice.SerialNumber()
		if errSerialNumber != nil {
			return false
		}
		return (myUdid == trimStringDescriptorValue(strSerialNumber))
	}
}

func isRootHub(d *gousb.Device) bool {
	return (d.Desc.Class == gousb.ClassHub && len(d.Desc.Path) <= 0)
}

func usbDeviceID(d *gousb.Device) string {
	return fmt.Sprintf("%d-%d", d.Desc.Bus, d.Desc.Address)
}

func findUsbDevice(prefilter myUsbDevicePrefilter, postfilter myUsbDevicePostfilter) []*gousb.Device {
	if prefilter == nil {
		prefilter = all_devices_opener
	}
	d, _ := gousbContext.OpenDevices(prefilter)
	devices := []*gousb.Device{}
	for _, device := range d {
		if postfilter != nil && !postfilter(device) {
			defer device.Close()
			continue
		}
		devices = append(devices, device)
	}
	return devices
}

func getUsbDeviceLineage(device *gousb.Device, levels int) []*gousb.Device {
	lineage := []*gousb.Device{}
	currentDevice := device
	for {
		if len(currentDevice.Desc.Path) <= 0 /* root hub */ || (levels > 0 && len(lineage) >= levels) {
			break
		}
		parentPath := currentDevice.Desc.Path[0 : len(currentDevice.Desc.Path)-1]
		parents := findUsbDevice(func(desc *gousb.DeviceDesc) bool {
			if desc.Bus != currentDevice.Desc.Bus || desc.Class != gousb.ClassHub {
				return false
			}
			return (slices.Compare(parentPath, desc.Path) == 0)
		}, nil)
		if len(parents) <= 0 {
			break
		}
		currentDevice = parents[0]
		lineage = util.Prepend(lineage, currentDevice)
	}
	return lineage
}

func getHubDescriptor(hub *gousb.Device) ([]byte, error) {
	desc := make([]byte, 1024)
	descLength, errDesc := hub.Control((USB_RT_HUB | gousb.ControlIn), REQUEST_GET_DESCRIPTOR, (DESCRIPTOR_TYPE_HUB << 8), 0, desc)
	if errDesc != nil {
		return nil, errDesc
	}
	if descLength <= 0 {
		return nil, xerrors.New("Invalid hub descriptor")
	}
	return desc[0:descLength], nil
}

func getHubPortDescriptor(hub *gousb.Device, index uint16) ([]byte, error) {
	desc := make([]byte, 4)
	descLength, errDesc := hub.Control((USB_RT_PORT | gousb.ControlIn), REQUEST_GET_STATUS, 0, index, desc)
	if errDesc != nil {
		return nil, errDesc
	}
	if descLength <= 0 {
		return nil, xerrors.New("Invalid hub descriptor")
	}
	return desc, nil
}

func getHubPortStatus(desc []byte) []string {
	state := []string{}
	if (desc[1] & 0x10) != 0 {
		state = append(state, "indicator")
	}
	if (desc[1] & 0x08) != 0 {
		state = append(state, "test")
	}
	if (desc[1] & 0x04) != 0 {
		state = append(state, "highspeed")
	}
	if (desc[1] & 0x02) != 0 {
		state = append(state, "lowspeed")
	}
	if (desc[1] & 0x01) != 0 {
		state = append(state, "power")
	}
	if (desc[0] & 0x10) != 0 {
		state = append(state, "RESET")
	}
	if (desc[0] & 0x08) != 0 {
		state = append(state, "oc")
	}
	if (desc[0] & 0x04) != 0 {
		state = append(state, "suspend")
	}
	if (desc[0] & 0x02) != 0 {
		state = append(state, "enable")
	}
	if (desc[0] & 0x01) != 0 {
		state = append(state, "connect")
	}
	return state
}

func toggleUsbHubPortPower(hub *gousb.Device, portIndex uint16, s bool) {
	req := REQUEST_CLEAR_FEATURE
	if s {
		req = REQUEST_SET_FEATURE
	}
	hub.Control(USB_RT_PORT, uint8(req), USB_PORT_FEAT_POWER, portIndex, nil)
}

/*
References:
https://github.com/google/gousb/issues/87
*/
