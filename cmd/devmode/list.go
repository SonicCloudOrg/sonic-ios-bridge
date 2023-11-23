package devmode

import (
	"encoding/json"
	"fmt"

	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var devmodeListCmd = &cobra.Command{
	Use:   "list",
	Short: "Print the Developer Mode status of connected devices",
	Long:  "Print the Developer Mode status of connected devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		util.InitLogger()
		usbMuxClient, err := giDevice.NewUsbmux()
		if err != nil {
			return util.NewErrorPrint(util.ErrConnect, "usbMux", err)
		}
		allErrors := []error{}
		localList, errReadLocal := usbMuxClient.Devices()
		if errReadLocal != nil {
			allErrors = append(allErrors, errReadLocal)
		}
		remoteList, _ := util.ReadRemote()
		devices := []entity.DevMode{}
		if len(localList) != 0 {
			for _, d := range localList {
				deviceByte, _ := json.Marshal(d.Properties())
				device := &entity.Device{}
				json.Unmarshal(deviceByte, device)
				if len(udid) > 0 && device.SerialNumber != udid {
					continue
				}
				detail, err2 := entity.GetDetail(d)
				if err2 != nil {
					allErrors = append(allErrors, err2)
				} else {
					device.DeviceDetail = *detail
				}
				devmode := entity.DevMode{Device: *device}
				devices = append(devices, devmode)
			}
		}
		if len(remoteList) != 0 {
			for _, d := range remoteList {
				deviceByte, _ := json.Marshal(d.Properties())
				device := &entity.Device{}
				json.Unmarshal(deviceByte, device)
				if len(udid) > 0 && device.SerialNumber != udid {
					continue
				}
				detail, err2 := entity.GetDetail(d)
				if err2 != nil {
					allErrors = append(allErrors, err2)
				} else {
					device.DeviceDetail = *detail
				}
				devmode := entity.DevMode{Device: *device}
				devices = append(devices, devmode)
			}
		}
		jsonResults := []map[string]any{}
		for _, devmode := range devices {
			jsonObj := map[string]any{"udid": devmode.SerialNumber, "status": "N/A"}
			bCanCheck, errChk := devmode.CanCheck()
			if errChk != nil {
				if bIsOutputJson {
					jsonObj["error"] = errChk.Error()
				} else {
					fmt.Printf("%s\t%s\n", devmode.SerialNumber, "Error: "+errChk.Error())
				}
			} else if !bCanCheck {
				if bIsOutputJson {
					jsonObj["error"] = "iOS version below 16"
				} else {
					fmt.Printf("%s\tN/A\n", devmode.SerialNumber)
				}
			} else {
				device := util.GetDeviceByUdId(devmode.SerialNumber)
				if device == nil {
					continue
				}
				interResult, errInfo := device.GetValue("com.apple.security.mac.amfi", "DeveloperModeStatus")
				if errInfo != nil {
					if bIsOutputJson {
						jsonObj["error"] = errInfo.Error()
					} else {
						fmt.Printf("%s\t%s\n", devmode.SerialNumber, "Error: "+errInfo.Error())
					}
				} else {
					strDevModeStatus := "N/A"
					switch interResult := interResult.(type) {
					case bool:
						if interResult {
							strDevModeStatus = "enabled"
						} else {
							strDevModeStatus = "disabled"
						}
					}
					if bIsOutputJson {
						jsonObj["status"] = strDevModeStatus
					} else {
						fmt.Printf("%s\t%s\n", devmode.SerialNumber, strDevModeStatus)
					}
				}
			}
			if bIsOutputJson {
				jsonResults = append(jsonResults, jsonObj)
			}
		}
		if len(jsonResults) > 0 && bIsOutputJson {
			if len(udid) > 0 {
				b, _ := json.Marshal(jsonResults[0])
				fmt.Println(string(b))
			} else {
				b, _ := json.Marshal(jsonResults)
				fmt.Println(string(b))
			}
		}
		if len(allErrors) > 0 {
			for _, e := range allErrors {
				logrus.Warnf("%+v\n", e)
			}
		}
		return nil
	},
}

func initDevModeListCmd() {
	devmodeRootCMD.AddCommand(devmodeListCmd)
	devmodeListCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	devmodeListCmd.Flags().BoolVarP(&bIsOutputJson, "json", "j", false, "output in JSON format")
}

/*
References:
https://github.com/libimobiledevice/libimobiledevice/blob/master/tools/idevicedevmodectl.c#L99
*/
