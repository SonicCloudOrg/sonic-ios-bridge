package util

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-gidevice/pkg/libimobiledevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/cenkalti/backoff/v4"
	"github.com/quintans/toolkit/latch"
	"github.com/sirupsen/logrus"
)

func UsbmuxListen(cbOnData func(gidevice *giDevice.Device, device *entity.Device, e error, cancelFunc context.CancelFunc)) context.CancelFunc {
	usbmuxInput := make(chan giDevice.Device)
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, os.Interrupt, os.Kill)
	var funcCancelListen context.CancelFunc
	healthCheck := make(chan bool)
	mylatch := latch.NewCountDownLatch()
	mylatch.Add(2)
	go (func() {
		backoffAlgorithm := backoff.NewConstantBackOff(30 * time.Second)
		bIsOk := true
		backoff.RetryNotify(func() error {
			if mylatch.Counter() <= 0 {
				return nil
			}
			usbMuxClient, err := giDevice.NewUsbmux()
			if err != nil {
				return NewErrorPrint(ErrConnect, "usbMux", err)
			}
			for {
				if mylatch.Counter() <= 0 {
					return nil
				}
				time.Sleep(5 * time.Second)
				if _, errBuid := usbMuxClient.ReadBUID(); errBuid != nil {
					return errBuid
				}
				logrus.Debug("usbmux health check success")
				if !bIsOk { // transition from not OK to OK
					logrus.Trace("Reset health check backoff algorithm")
					backoffAlgorithm.Reset()
					bIsOk = true
				}
			}
		}, backoffAlgorithm, func(err error, d time.Duration) {
			bIsOk = false
			logrus.Warnf("usbmux health check error: %+v", err)
			healthCheck <- false
		})
		logrus.Trace("end health check")
	})()
	go (func(funcStop *context.CancelFunc) {
		backoffAlgorithm := backoff.NewConstantBackOff(30 * time.Second)
		bIsOk := true
		backoff.RetryNotify(func() error {
			if mylatch.Counter() <= 1 { // 'read channel input' go routine is stopped
				return nil
			}
			usbMuxClient, err := giDevice.NewUsbmux()
			if err != nil {
				return NewErrorPrint(ErrConnect, "usbMux", err)
			}
			shutDownFun, errListen := usbMuxClient.Listen(usbmuxInput)
			(*funcStop) = shutDownFun
			if errListen != nil {
				return NewErrorPrint(ErrSendCommand, string(libimobiledevice.MessageTypeListen), errListen)
			}
			logrus.Info("Start listening...")
			<-healthCheck // empty out the channel
			if !bIsOk {   // transition from not OK to OK
				logrus.Trace("Reset usbmux listen backoff algorithm")
				backoffAlgorithm.Reset()
				bIsOk = true
			}
			for range healthCheck {
				logrus.Info("Cancel listening")
				(*funcStop) = nil
				return fmt.Errorf("usbmux listening is cancelled")
			}
			return nil
		}, backoffAlgorithm, func(err error, d time.Duration) {
			bIsOk = false
			logrus.Warnf("usbmux listening error: %+v", err)
		})
		mylatch.Done()
		logrus.Trace("end usbmux listen")
	})(&funcCancelListen)
	go (func(funcStop *context.CancelFunc) {
		numOnlineDevices := 0
	loopRead:
		for {
			select {
			case d, ok := <-usbmuxInput:
				if !ok { // channel is closed
					logrus.Info("usbmux input channel is closed")
					break loopRead
				}
				if d == nil {
					continue
				}
				deviceByte, _ := json.Marshal(d.Properties())
				var device entity.Device
				errDec := json.Unmarshal(deviceByte, &device)
				var ptrEntityDevice *entity.Device = &device
				if errDec == nil {
					device.Status = device.GetStatus()
					if device.Status == "online" {
						numOnlineDevices += 1
					} else {
						numOnlineDevices -= 1
					}
				} else {
					ptrEntityDevice = nil
				}
				var _fStop context.CancelFunc
				if funcStop != nil {
					_fStop = *funcStop
				}
				if cbOnData != nil {
					cbOnData(&d, ptrEntityDevice, errDec, _fStop)
				}
				logrus.Debugf("Number of online devices= %d", numOnlineDevices)
				if numOnlineDevices <= 0 {
					logrus.Info("No devices are online")
				}
			case <-sigTerm:
				logrus.Info("Stop listening")
				if funcStop != nil {
					(*funcStop)()
				}
				os.Exit(128 + int(syscall.SIGTERM)) // https://itsfoss.com/linux-exit-codes/#code-143-or-sigterm
				break loopRead
			}
		}
		mylatch.Done()
		logrus.Trace("end reading channel input")
	})(&funcCancelListen)
	return funcCancelListen
}
