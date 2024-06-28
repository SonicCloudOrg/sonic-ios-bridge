package util

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"unsafe"

	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/cenkalti/backoff/v4"
	"github.com/quintans/toolkit/latch"
	"github.com/sirupsen/logrus"
)

type UsbmuxListenCallback func(gidevice *giDevice.Device, device *entity.Device, e error, cancelFunc context.CancelFunc)

func CloseUsbmuxClient(usbMuxClient giDevice.Usbmux) error {
	var output_err error = nil
	TryCatch{}.Try(func() {
		// https://stackoverflow.com/a/59196685/12857692
		v := reflect.ValueOf(usbMuxClient).Elem()
		f := v.FieldByName("client")
		rf := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
		rf.MethodByName("Close").Call(nil)
	}).CatchAll(func(err error) {
		output_err = err
	})
	return output_err
} // end CloseUsbmuxClient()

func UsbmuxListen(cbOnData UsbmuxListenCallback, bExitOnCtrlC bool) context.CancelFunc {
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, os.Interrupt, os.Kill)
	var funcCancelListen context.CancelFunc
	healthCheck := make(chan bool)
	mylatch := latch.NewCountDownLatch()
	mylatch.Add(1)
	go (func() {
		logger := GetGoRoutineLogger("healthCheck")
		backoffAlgorithm := backoff.NewConstantBackOff(10 * time.Second)
		bIsOk := true
		backoff.RetryNotify(func() error {
			if mylatch.Counter() <= 0 {
				return nil
			} // end if
			for {
				if mylatch.Counter() <= 0 {
					return nil
				} // end if
				time.Sleep(15 * time.Second)
				logger.Debug("Connecting to usbmux...")
				usbMuxClient, err := giDevice.NewUsbmux()
				if err != nil {
					return err
				} // end if
				if _, errBuid := usbMuxClient.ReadBUID(); errBuid != nil {
					return errBuid
				} // end if
				CloseUsbmuxClient(usbMuxClient)
				logger.Debug("usbmux health check success")
				healthCheck <- true
				if !bIsOk { // transition from not OK to OK
					logger.Trace("Reset health check backoff algorithm")
					backoffAlgorithm.Reset()
					bIsOk = true
				} // end if
			} // end for
		}, backoffAlgorithm, func(err error, d time.Duration) {
			bIsOk = false
			logger.Warn(err)
			healthCheck <- false
			logger.Debugf("Next retry health check in %s", d.String())
		})
		logger.Trace("end health check")
	})()
	go (func(funcStop *context.CancelFunc) {
		logger := GetGoRoutineLogger("signalHandler")
		for range sigTerm {
			logger.Debugf("Stop listening by signal")
			if funcStop != nil && *funcStop != nil {
				(*funcStop)()
			} // end if
			if bExitOnCtrlC {
				os.Exit(128 + int(syscall.SIGTERM)) // https://itsfoss.com/linux-exit-codes/#code-143-or-sigterm
			} // end if
		} // end for
	})(&funcCancelListen)
	go (func(funcStop *context.CancelFunc) {
		logger := GetGoRoutineLogger("usbmuxListen")
		backoffAlgorithm := backoff.NewConstantBackOff(10 * time.Second)
		bIsOk := true
		backoff.RetryNotify(func() error {
			logger.Debug("Connecting to usbmux...")
			usbMuxClient, err := giDevice.NewUsbmux()
			if err != nil {
				return err
			} // end if
			usbmuxInput := make(chan giDevice.Device)
			shutDownFun, errListen := usbMuxClient.Listen(usbmuxInput)
			*funcStop = func() {
				logrus.Debugf("Call usbmux listen shutdown function")
				mylatch.Done()
				if shutDownFun != nil {
					shutDownFun()
				} // end if
			}
			if errListen != nil {
				return errListen
			} // end if
			if !bIsOk { // transition from not OK to OK
				logger.Trace("Reset usbmux listen backoff algorithm")
				backoffAlgorithm.Reset()
				bIsOk = true
			} // end if
			logger.Debugf("Start listening...")
			numOnlineDevices := 0
		loopRead:
			for {
				select {
				case bIsHealthy := <-healthCheck:
					if !bIsHealthy {
						logger.Info("Cancel listening because usbmuxd is unhealthy")
						return fmt.Errorf("usbmux listening is cancelled")
					} // end if
				case d, ok := <-usbmuxInput:
					if !ok { // channel is closed
						logger.Info("usbmux input channel is closed")
						if mylatch.Counter() > 0 {
							return fmt.Errorf("usbmux listening stopped unexpectedly")
						} // end if
						break loopRead
					} // end if
					if d == nil {
						continue
					} // end if
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
						} // end if
					} else {
						ptrEntityDevice = nil
					} // end if
					var _fStop context.CancelFunc
					if funcStop != nil {
						_fStop = *funcStop
					} // end if
					if cbOnData != nil {
						cbOnData(&d, ptrEntityDevice, errDec, _fStop)
					} // end if
					logger.Debugf("Number of online devices= %d", numOnlineDevices)
					if numOnlineDevices <= 0 {
						logger.Info("No devices are online")
					} // end if
				} // end select
			} // end for
			CloseUsbmuxClient(usbMuxClient)
			return nil
		}, backoffAlgorithm, func(err error, d time.Duration) {
			bIsOk = false
			logger.Warn(err)
			<-healthCheck // clear queued signal
			logger.Debugf("Next retry listening in %s", d.String())
		})
		mylatch.Done()
		logger.Trace("end usbmux listen")
	})(&funcCancelListen)
	return funcCancelListen
} // end UsbmuxListen()
