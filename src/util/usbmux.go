package util

import (
	"encoding/json"
	"os"
	"os/signal"

	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
)

func UsbmuxListen(cb func(gidevice *giDevice.Device, device *entity.Device, e error, cancelFunc func())) (func(), error) {
	usbMuxClient, err := giDevice.NewUsbmux()
	if err != nil {
		return nil, NewErrorPrint(ErrConnect, "usbMux", err)
	}
	usbmuxInput := make(chan giDevice.Device)
	shutdownUsbmuxFun, err2 := usbMuxClient.Listen(usbmuxInput)
	if err2 != nil {
		return nil, NewErrorPrint(ErrSendCommand, "listen", err2)
	}
	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, os.Interrupt, os.Kill)
	go func() {
		for {
			select {
			case d, ok := <-usbmuxInput:
				if !ok { // usbmux channel closed
					close(shutDown)
					return
				}
				if d == nil {
					continue
				}
				deviceByte, _ := json.Marshal(d.Properties())
				device := &entity.Device{}
				errDecode := json.Unmarshal(deviceByte, device)
				if errDecode != nil {
					cb(nil, nil, errDecode, shutdownUsbmuxFun)
					continue
				}
				device.Status = device.GetStatus()
				cb(&d, device, nil, shutdownUsbmuxFun)
			case <-shutDown:
				shutdownUsbmuxFun()
				return
			}
		}
	}()
	return shutdownUsbmuxFun, nil
}
