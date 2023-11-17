package usb

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/SonicCloudOrg/sonic-ios-bridge/src/errorcodes"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/google/gousb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// option bindings
var intBus, intDevAddr int
var intHubPortIndex uint16
var bPowerOn, bPowerOff bool

var usbHubPortPowerCmd = &cobra.Command{
	Use:   "power",
	Short: "Turn on/off power to a USB hub port",
	Long:  "Turn on/off power to a USB hub port",
	RunE: func(cmd *cobra.Command, args []string) error {
		util.InitLogger()
		devices := findUsbDevice(func(desc *gousb.DeviceDesc) bool {
			return (desc.Class == gousb.Class(gousb.ClassHub) && desc.Address == intDevAddr && desc.Bus == intBus)
		}, isMySerial(udid))
		if len(devices) <= 0 {
			os.Exit(errorcodes.ERROR_DEV_NOT_EXIST)
		}
		hub := devices[0]
		defer hub.Close()
		if bPowerOn == bPowerOff { // show status
			descHub, errDescHub := getHubDescriptor(hub)
			if errDescHub != nil {
				logrus.Warn(errDescHub)
				os.Exit(errorcodes.ERROR_READ_FAULT)
			}
			if len(descHub) < 3 {
				logrus.Warnf("Unable to determine number of ports: descriptor= %+v", descHub)
				os.Exit(errorcodes.ERROR_BAD_LENGTH)
			}
			numberOfPorts := uint16(descHub[2])
			if intHubPortIndex > numberOfPorts {
				logrus.Warnf("Port index out of range (must be 1 ~ %d)", numberOfPorts)
				os.Exit(errorcodes.ERROR_INVALID_PARAMETER)
			}
			results := []map[string]any{}
			for i := uint16(1); i <= numberOfPorts; i++ {
				if intHubPortIndex > 0 && i != intHubPortIndex {
					continue
				}
				descHubPort, errDescPort := getHubPortDescriptor(hub, uint16(i))
				resultItem := map[string]any{"index": i}
				if errDescPort != nil {
					if bIsOutputJson {
						resultItem["error"] = errDescPort.Error()
					} else {
						fmt.Printf("Port #%d: %+v\n", i, errDescPort)
					}
				} else {
					state := getHubPortStatus(descHubPort)
					if bIsOutputJson {
						resultItem["status"] = state
					} else {
						fmt.Printf("Port #%d: %+v\n", i, strings.Join(state, ", "))
					}
				}
				results = append(results, resultItem)
			}
			if bIsOutputJson {
				b, _ := json.Marshal(results)
				fmt.Println(string(b))
			}
		} else {
			if intHubPortIndex <= 0 {
				logrus.Debug("Port index ('--port') must be greater than 0")
				os.Exit(errorcodes.ERROR_INVALID_PARAMETER)
			}
			s := true
			if bPowerOff {
				s = false
			}
			toggleUsbHubPortPower(hub, intHubPortIndex, s)
		}
		return nil
	},
}

func initUsbHubPortPowerCmd() {
	usbRootCMD.AddCommand(usbHubPortPowerCmd)
	usbHubPortPowerCmd.Flags().IntVarP(&intBus, "bus", "b", 0, "bus")
	usbHubPortPowerCmd.MarkFlagRequired("bus")
	usbHubPortPowerCmd.Flags().IntVarP(&intDevAddr, "device", "d", 0, "device address")
	usbHubPortPowerCmd.MarkFlagRequired("device")
	usbHubPortPowerCmd.Flags().Uint16VarP(&intHubPortIndex, "port", "p", 0, "port index")
	usbHubPortPowerCmd.Flags().BoolVarP(&bPowerOn, "on", "1", false, "power ON")
	usbHubPortPowerCmd.Flags().BoolVarP(&bPowerOff, "off", "0", false, "power OFF")
	usbHubPortPowerCmd.Flags().BoolVarP(&bIsOutputJson, "json", "j", false, "output in JSON format")
}
