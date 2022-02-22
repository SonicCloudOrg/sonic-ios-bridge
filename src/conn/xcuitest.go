package conn

import (
	"fmt"
	"github.com/Masterminds/semver"
	log "github.com/sirupsen/logrus"
)

func RunWebDriverAgent(device iDevice, webDriverAgentBundleID, webDriverAgentBundleIDTestRunnerBundleID string) error {
	version, err := device.GetSemverProductVersion()
	if err != nil {
		return err
	}
	if version.LessThan(semver.MustParse("14.0")) {
		log.Infof("udId: %s iOSVersion: %s, WebDriverAgent will run without secure...", device.Properties.SerialNumber, version)
	}
	return nil
}

func (device *iDevice) ConnectService(serviceName string) (DeviceConnectInterface, error) {
	//step1 start service , get resp
	startServiceResp, err := device.StartService(serviceName)
	if err != nil {
		return nil, err
	}
	//step2 get pair record
	usbMuxClient, err := NewUsbMuxClient()
	if err != nil {
		return nil, err
	}
	pairRecord, err := usbMuxClient.ReadPair(device.Properties.SerialNumber)
	if err != nil {
		return nil, err
	}
}
